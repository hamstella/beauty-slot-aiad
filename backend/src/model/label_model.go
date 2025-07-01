package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Label struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:50;not null;uniqueIndex" json:"name" validate:"required,min=1,max=50"`
	Color       string    `gorm:"size:7;default:#000000" json:"color" validate:"omitempty,len=7"` // hex color
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (l *Label) TableName() string {
	return "labels"
}