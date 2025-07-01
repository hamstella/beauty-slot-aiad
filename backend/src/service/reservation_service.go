package service

import (
	"app/src/model"
	"app/src/utils"
	"errors"

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

func (s *ReservationService) GetReservations(page, limit int, status, date string) ([]model.Reservation, int64, error) {
	var reservations []model.Reservation
	var total int64

	offset := (page - 1) * limit
	query := s.db.Model(&model.Reservation{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if date != "" {
		query = query.Where("DATE(reservation_date) = ?", date)
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

func (s *ReservationService) CancelReservation(id uuid.UUID, reason string) (*model.Reservation, error) {
	var reservation model.Reservation
	if err := s.db.Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	// Don't allow cancelling already cancelled or completed reservations
	if reservation.Status == model.ReservationStatusCancelled {
		return nil, errors.New("reservation is already cancelled")
	}
	if reservation.Status == model.ReservationStatusCompleted {
		return nil, errors.New("cannot cancel completed reservation")
	}

	// Update status and reason
	reservation.Status = model.ReservationStatusCancelled
	reservation.CancellationReason = reason

	if err := s.db.Save(&reservation).Error; err != nil {
		utils.Log.Errorf("Failed to cancel reservation: %v", err)
		return nil, err
	}

	return s.GetReservationByID(reservation.ID)
}