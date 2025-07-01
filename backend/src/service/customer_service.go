package service

import (
	"app/src/model"
	"app/src/utils"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		db:        db,
		validator: validator.New(),
	}
}

func (s *CustomerService) GetCustomers(page, limit int) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	// Check if database is available
	if s.db == nil {
		return customers, 0, errors.New("database connection not available")
	}

	offset := (page - 1) * limit

	// Count total records
	if err := s.db.Model(&model.Customer{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		utils.Log.Errorf("Failed to count customers: %v", err)
		return nil, 0, err
	}

	// Get paginated records
	if err := s.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&customers).Error; err != nil {
		utils.Log.Errorf("Failed to get customers: %v", err)
		return nil, 0, err
	}

	return customers, total, nil
}

func (s *CustomerService) GetCustomerByID(id uuid.UUID) (*model.Customer, error) {
	var customer model.Customer
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		utils.Log.Errorf("Failed to get customer: %v", err)
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) CreateCustomer(customer *model.Customer) (*model.Customer, error) {
	// Validate input
	if err := s.validator.Struct(customer); err != nil {
		utils.Log.Errorf("Customer validation failed: %v", err)
		return nil, err
	}

	// Check if phone already exists
	var existingCustomer model.Customer
	if err := s.db.Where("phone = ? AND is_active = ?", customer.Phone, true).First(&existingCustomer).Error; err == nil {
		return nil, errors.New("phone number already exists")
	}

	// Check if email already exists (if provided)
	if customer.Email != "" {
		if err := s.db.Where("email = ? AND is_active = ?", customer.Email, true).First(&existingCustomer).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	if err := s.db.Create(customer).Error; err != nil {
		utils.Log.Errorf("Failed to create customer: %v", err)
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) UpdateCustomer(customer *model.Customer) (*model.Customer, error) {
	// Validate input
	if err := s.validator.Struct(customer); err != nil {
		utils.Log.Errorf("Customer validation failed: %v", err)
		return nil, err
	}

	// Check if customer exists
	var existingCustomer model.Customer
	if err := s.db.Where("id = ? AND is_active = ?", customer.ID, true).First(&existingCustomer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	// Check if phone already exists for other customers
	var phoneCheck model.Customer
	if err := s.db.Where("phone = ? AND id != ? AND is_active = ?", customer.Phone, customer.ID, true).First(&phoneCheck).Error; err == nil {
		return nil, errors.New("phone number already exists")
	}

	// Check if email already exists for other customers (if provided)
	if customer.Email != "" {
		var emailCheck model.Customer
		if err := s.db.Where("email = ? AND id != ? AND is_active = ?", customer.Email, customer.ID, true).First(&emailCheck).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	if err := s.db.Save(customer).Error; err != nil {
		utils.Log.Errorf("Failed to update customer: %v", err)
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) DeleteCustomer(id uuid.UUID) error {
	// Soft delete by setting is_active to false
	result := s.db.Model(&model.Customer{}).Where("id = ? AND is_active = ?", id, true).Update("is_active", false)
	if result.Error != nil {
		utils.Log.Errorf("Failed to delete customer: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}

	return nil
}