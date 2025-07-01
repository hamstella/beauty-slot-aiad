package service

import (
	"app/src/model"
	"app/src/utils"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationService struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewReservationService(db *gorm.DB) *ReservationService {
	return &ReservationService{
		db:        db,
		validator: validator.New(),
	}
}

func (s *ReservationService) GetReservations(page, limit int, status, staffID, customerID, dateFrom, dateTo string) ([]model.Reservation, int64, error) {
	var reservations []model.Reservation
	var total int64

	offset := (page - 1) * limit
	query := s.db.Model(&model.Reservation{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if staffID != "" {
		query = query.Where("staff_id = ?", staffID)
	}
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if dateFrom != "" {
		query = query.Where("reservation_date >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("reservation_date <= ?", dateTo)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		utils.Log.Errorf("Failed to count reservations: %v", err)
		return nil, 0, err
	}

	// Get paginated records with preloaded relations
	if err := query.Preload("Customer").Preload("Staff").
		Preload("ReservationMenus.Menu").
		Preload("ReservationOptions.Option").
		Order("reservation_date DESC, start_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&reservations).Error; err != nil {
		utils.Log.Errorf("Failed to get reservations: %v", err)
		return nil, 0, err
	}

	return reservations, total, nil
}

func (s *ReservationService) GetReservationByID(id uuid.UUID) (*model.Reservation, error) {
	var reservation model.Reservation
	if err := s.db.Preload("Customer").Preload("Staff").
		Preload("ReservationMenus.Menu").
		Preload("ReservationOptions.Option").
		Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		utils.Log.Errorf("Failed to get reservation: %v", err)
		return nil, err
	}
	return &reservation, nil
}

func (s *ReservationService) CreateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	// Validate input
	if err := s.validator.Struct(reservation); err != nil {
		utils.Log.Errorf("Reservation validation failed: %v", err)
		return nil, err
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate customer exists
	var customer model.Customer
	if err := tx.Where("id = ? AND is_active = ?", reservation.CustomerID, true).First(&customer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("customer not found")
	}

	// Validate staff exists
	var staff model.Staff
	if err := tx.Where("id = ? AND is_active = ?", reservation.StaffID, true).First(&staff).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("staff not found")
	}

	// Check for time conflicts
	var conflictReservation model.Reservation
	if err := tx.Where("staff_id = ? AND reservation_date = ? AND status NOT IN (?, ?) AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
		reservation.StaffID,
		reservation.ReservationDate.Format("2006-01-02"),
		model.ReservationStatusCancelled,
		model.ReservationStatusNoShow,
		reservation.StartTime,
		reservation.StartTime,
		reservation.EndTime,
		reservation.EndTime).First(&conflictReservation).Error; err == nil {
		tx.Rollback()
		return nil, errors.New("time slot is already booked")
	}

	// Set default status if not provided
	if reservation.Status == "" {
		reservation.Status = model.ReservationStatusPending
	}

	// Create reservation
	if err := tx.Create(reservation).Error; err != nil {
		tx.Rollback()
		utils.Log.Errorf("Failed to create reservation: %v", err)
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.Log.Errorf("Failed to commit reservation transaction: %v", err)
		return nil, err
	}

	// Reload with relations
	return s.GetReservationByID(reservation.ID)
}

func (s *ReservationService) UpdateReservation(reservation *model.Reservation) (*model.Reservation, error) {
	// Validate input
	if err := s.validator.Struct(reservation); err != nil {
		utils.Log.Errorf("Reservation validation failed: %v", err)
		return nil, err
	}

	// Check if reservation exists
	var existingReservation model.Reservation
	if err := s.db.Where("id = ?", reservation.ID).First(&existingReservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	// Don't allow updating cancelled or completed reservations
	if existingReservation.Status == model.ReservationStatusCancelled ||
		existingReservation.Status == model.ReservationStatusCompleted {
		return nil, errors.New("cannot update cancelled or completed reservations")
	}

	if err := s.db.Save(reservation).Error; err != nil {
		utils.Log.Errorf("Failed to update reservation: %v", err)
		return nil, err
	}

	return s.GetReservationByID(reservation.ID)
}

func (s *ReservationService) CancelReservation(id uuid.UUID) error {
	var reservation model.Reservation
	if err := s.db.Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation not found")
		}
		return err
	}

	// Don't allow cancelling already cancelled or completed reservations
	if reservation.Status == model.ReservationStatusCancelled {
		return errors.New("reservation is already cancelled")
	}
	if reservation.Status == model.ReservationStatusCompleted {
		return errors.New("cannot cancel completed reservation")
	}

	// Update status
	reservation.Status = model.ReservationStatusCancelled

	if err := s.db.Save(&reservation).Error; err != nil {
		utils.Log.Errorf("Failed to cancel reservation: %v", err)
		return err
	}

	return nil
}

func (s *ReservationService) CreateReservationFromRequest(customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error) {
	// Parse and validate date
	parsedDate, err := time.Parse("2006-01-02", reservationDate)
	if err != nil {
		return nil, errors.New("無効な日付形式です")
	}

	// Check if date is in the future (at least tomorrow)
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	if parsedDate.Before(tomorrow) {
		return nil, errors.New("予約は翌日以降の日付で設定してください")
	}

	// Check if date is within 90 days
	maxDate := time.Now().AddDate(0, 0, 90)
	if parsedDate.After(maxDate) {
		return nil, errors.New("予約は90日以内の日付で設定してください")
	}

	// Parse start time
	parsedStartTime, err := time.Parse("15:04:05", startTime)
	if err != nil {
		return nil, errors.New("無効な時刻形式です")
	}

	// Calculate total duration from menus
	var totalDuration int
	var totalPrice int
	for _, menuID := range menuIDs {
		var menu model.Menu
		if err := s.db.Where("id = ? AND is_active = ?", menuID, true).First(&menu).Error; err != nil {
			return nil, errors.New("選択されたメニューが見つかりません")
		}
		totalDuration += menu.Duration
		totalPrice += menu.Price
	}

	// Add option prices
	for _, optionID := range optionIDs {
		var option model.Option
		if err := s.db.Where("id = ? AND is_active = ?", optionID, true).First(&option).Error; err != nil {
			return nil, errors.New("選択されたオプションが見つかりません")
		}
		totalPrice += option.Price
	}

	// Calculate end time
	endTime := parsedStartTime.Add(time.Duration(totalDuration) * time.Minute)

	// Create reservation
	reservation := &model.Reservation{
		CustomerID:      customerID,
		StaffID:         staffID,
		ReservationDate: parsedDate,
		StartTime:       time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedStartTime.Hour(), parsedStartTime.Minute(), parsedStartTime.Second(), 0, time.Local),
		EndTime:         time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.Local),
		Status:          model.ReservationStatusConfirmed,
		TotalDuration:   totalDuration,
		TotalPrice:      totalPrice,
		Notes:           notes,
	}

	return s.CreateReservation(reservation)
}

func (s *ReservationService) UpdateReservationFromRequest(id, customerID, staffID uuid.UUID, reservationDate, startTime string, menuIDs, optionIDs []uuid.UUID, notes string) (*model.Reservation, error) {
	// Get existing reservation
	existingReservation, err := s.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	// Check if can be updated (not completed or cancelled)
	if existingReservation.Status == model.ReservationStatusCompleted ||
		existingReservation.Status == model.ReservationStatusCancelled {
		return nil, errors.New("cannot update cancelled or completed reservations")
	}

	// Parse and validate date if provided
	var parsedDate time.Time
	if reservationDate != "" {
		parsedDate, err = time.Parse("2006-01-02", reservationDate)
		if err != nil {
			return nil, errors.New("無効な日付形式です")
		}
	} else {
		parsedDate = existingReservation.ReservationDate
	}

	// Parse start time if provided
	var parsedStartTime time.Time
	if startTime != "" {
		parsedStartTime, err = time.Parse("15:04:05", startTime)
		if err != nil {
			return nil, errors.New("無効な時刻形式です")
		}
	} else {
		parsedStartTime = existingReservation.StartTime
	}

	// Use existing values if not provided
	if customerID == uuid.Nil {
		customerID = existingReservation.CustomerID
	}
	if staffID == uuid.Nil {
		staffID = existingReservation.StaffID
	}

	// Calculate duration and price if menus changed
	var totalDuration int
	var totalPrice int
	if len(menuIDs) > 0 {
		for _, menuID := range menuIDs {
			var menu model.Menu
			if err := s.db.Where("id = ? AND is_active = ?", menuID, true).First(&menu).Error; err != nil {
				return nil, errors.New("選択されたメニューが見つかりません")
			}
			totalDuration += menu.Duration
			totalPrice += menu.Price
		}

		for _, optionID := range optionIDs {
			var option model.Option
			if err := s.db.Where("id = ? AND is_active = ?", optionID, true).First(&option).Error; err != nil {
				return nil, errors.New("選択されたオプションが見つかりません")
			}
			totalPrice += option.Price
		}
	} else {
		totalPrice = existingReservation.TotalPrice
		// Calculate existing duration from current menus
		for _, rm := range existingReservation.ReservationMenus {
			totalDuration += rm.Menu.Duration
		}
	}

	// Calculate end time
	endTime := parsedStartTime.Add(time.Duration(totalDuration) * time.Minute)

	// Update reservation
	existingReservation.CustomerID = customerID
	existingReservation.StaffID = staffID
	existingReservation.ReservationDate = parsedDate
	existingReservation.StartTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), parsedStartTime.Hour(), parsedStartTime.Minute(), parsedStartTime.Second(), 0, time.Local)
	existingReservation.EndTime = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.Local)
	existingReservation.TotalDuration = totalDuration
	existingReservation.TotalPrice = totalPrice
	if notes != "" {
		existingReservation.Notes = notes
	}

	return s.UpdateReservation(existingReservation)
}

func (s *ReservationService) UpdateReservationStatus(id uuid.UUID, status string) (*model.Reservation, error) {
	var reservation model.Reservation
	if err := s.db.Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	// Validate status transition
	validTransitions := map[model.ReservationStatus][]model.ReservationStatus{
		model.ReservationStatusPending:   {model.ReservationStatusConfirmed, model.ReservationStatusCancelled},
		model.ReservationStatusConfirmed: {model.ReservationStatusCompleted, model.ReservationStatusCancelled},
	}

	newStatus := model.ReservationStatus(status)
	if validStatuses, exists := validTransitions[reservation.Status]; exists {
		valid := false
		for _, validStatus := range validStatuses {
			if newStatus == validStatus {
				valid = true
				break
			}
		}
		if !valid {
			return nil, errors.New("invalid status transition")
		}
	} else {
		return nil, errors.New("invalid status transition")
	}

	// Update status
	reservation.Status = newStatus
	if err := s.db.Save(&reservation).Error; err != nil {
		utils.Log.Errorf("Failed to update reservation status: %v", err)
		return nil, err
	}

	return s.GetReservationByID(reservation.ID)
}

func (s *ReservationService) GetAvailability(date, durationStr, staffIDStr, menuIDsStr string) ([]map[string]interface{}, error) {
	// Parse date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("無効な日付形式です")
	}

	// Parse duration
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return nil, errors.New("無効な時間形式です")
	}

	// Get staff list
	var staffList []model.Staff
	query := s.db.Where("is_active = ?", true)
	
	if staffIDStr != "" {
		staffID, err := uuid.Parse(staffIDStr)
		if err != nil {
			return nil, errors.New("無効なスタッフIDです")
		}
		query = query.Where("id = ?", staffID)
	}

	if err := query.Find(&staffList).Error; err != nil {
		return nil, err
	}

	var availableSlots []map[string]interface{}

	for _, staff := range staffList {
		// Check if staff has shift for this date
		var shift model.Shift
		if err := s.db.Where("staff_id = ? AND date = ?", staff.ID, parsedDate.Format("2006-01-02")).First(&shift).Error; err != nil {
			continue // No shift for this staff on this date
		}

		// Get existing reservations for this staff on this date
		var reservations []model.Reservation
		s.db.Where("staff_id = ? AND reservation_date = ? AND status NOT IN (?, ?)",
			staff.ID, parsedDate.Format("2006-01-02"),
			model.ReservationStatusCancelled, model.ReservationStatusNoShow).Find(&reservations)

		// Generate available time slots
		availableTimes := s.calculateAvailableTimes(shift, reservations, duration)

		if len(availableTimes) > 0 {
			availableSlots = append(availableSlots, map[string]interface{}{
				"staff_id":         staff.ID,
				"staff_name":       staff.Name,
				"available_times":  availableTimes,
			})
		}
	}

	return availableSlots, nil
}

func (s *ReservationService) calculateAvailableTimes(shift model.Shift, reservations []model.Reservation, duration int) []map[string]interface{} {
	var availableTimes []map[string]interface{}
	
	// Use shift times directly
	shiftStart := shift.StartTime
	shiftEnd := shift.EndTime
	
	// Create time slots with 15-minute intervals
	current := shiftStart
	slotDuration := 15 * time.Minute
	requiredDuration := time.Duration(duration) * time.Minute
	
	for current.Add(requiredDuration).Before(shiftEnd) || current.Add(requiredDuration).Equal(shiftEnd) {
		slotEnd := current.Add(requiredDuration)
		
		// Check if this slot conflicts with any reservation
		conflicts := false
		for _, reservation := range reservations {
			resStart := reservation.StartTime
			resEnd := reservation.EndTime
			
			// Add 15-minute buffer after each reservation
			resEnd = resEnd.Add(15 * time.Minute)
			
			if (current.Before(resEnd) && slotEnd.After(resStart)) {
				conflicts = true
				break
			}
		}
		
		if !conflicts {
			availableTimes = append(availableTimes, map[string]interface{}{
				"start_time": current.Format("15:04:05"),
				"end_time":   slotEnd.Format("15:04:05"),
			})
		}
		
		current = current.Add(slotDuration)
	}
	
	return availableTimes
}