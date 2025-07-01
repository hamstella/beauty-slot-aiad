package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Phone       string    `gorm:"size:20;uniqueIndex;not null" json:"phone" validate:"required,min=10,max=20"`
	Email       string    `gorm:"size:255;uniqueIndex" json:"email" validate:"omitempty,email,max=255"`
	Birthday    *time.Time `gorm:"type:date" json:"birthday"`
	Gender      string    `gorm:"size:10" json:"gender" validate:"omitempty,oneof=male female other"`
	Notes       string    `gorm:"type:text" json:"notes"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	Reservations []Reservation `gorm:"foreignKey:CustomerID" json:"reservations,omitempty"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *Customer) TableName() string {
	return "customers"
}