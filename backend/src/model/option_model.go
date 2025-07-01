package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Option struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Description string    `gorm:"type:text" json:"description"`
	Duration    int       `gorm:"not null" json:"duration" validate:"required,min=0,max=120"` // additional minutes
	Price       int       `gorm:"not null" json:"price" validate:"required,min=0"`             // additional yen
	Category    string    `gorm:"size:50" json:"category" validate:"omitempty,max=50"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	ReservationOptions []ReservationOption `gorm:"foreignKey:OptionID" json:"reservation_options,omitempty"`
}

func (o *Option) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (o *Option) TableName() string {
	return "options"
}