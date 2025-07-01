package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shift struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StaffID   uuid.UUID `gorm:"type:uuid;not null;index" json:"staff_id" validate:"required"`
	Date      time.Time `gorm:"type:date;not null;index" json:"date" validate:"required"`
	StartTime time.Time `gorm:"not null" json:"start_time" validate:"required"`
	EndTime   time.Time `gorm:"not null" json:"end_time" validate:"required"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	Staff Staff `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
}

func (s *Shift) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (s *Shift) TableName() string {
	return "shifts"
}