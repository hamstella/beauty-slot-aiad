package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Description string    `gorm:"type:text" json:"description"`
	Duration    int       `gorm:"not null" json:"duration" validate:"required,min=1,max=600"` // minutes
	Price       int       `gorm:"not null" json:"price" validate:"required,min=0"`            // yen
	Category    string    `gorm:"size:50" json:"category" validate:"omitempty,max=50"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	ReservationMenus []ReservationMenu `gorm:"foreignKey:MenuID" json:"reservation_menus,omitempty"`
}

func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *Menu) TableName() string {
	return "menus"
}