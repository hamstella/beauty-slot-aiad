package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Staff struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Email       string    `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email,max=255"`
	Phone       string    `gorm:"size:20" json:"phone" validate:"omitempty,min=10,max=20"`
	Position    string    `gorm:"size:50" json:"position" validate:"omitempty,max=50"`
	Specialties string    `gorm:"type:text" json:"specialties"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	Reservations []Reservation `gorm:"foreignKey:StaffID" json:"reservations,omitempty"`
	Shifts       []Shift       `gorm:"foreignKey:StaffID" json:"shifts,omitempty"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (s *Staff) TableName() string {
	return "staff"
}