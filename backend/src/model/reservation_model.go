package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCompleted ReservationStatus = "completed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusNoShow    ReservationStatus = "no_show"
)

type Reservation struct {
	ID                  uuid.UUID         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID          uuid.UUID         `gorm:"type:uuid;not null;index" json:"customer_id" validate:"required"`
	StaffID             uuid.UUID         `gorm:"type:uuid;not null;index" json:"staff_id" validate:"required"`
	ReservationDate     time.Time         `gorm:"not null;index" json:"reservation_date" validate:"required"`
	StartTime           time.Time         `gorm:"not null" json:"start_time" validate:"required"`
	EndTime             time.Time         `gorm:"not null" json:"end_time" validate:"required"`
	Status              ReservationStatus `gorm:"size:20;not null;default:pending" json:"status" validate:"required,oneof=pending confirmed completed cancelled no_show"`
	TotalDuration       int               `gorm:"not null" json:"total_duration"`        // minutes
	TotalPrice          int               `gorm:"not null" json:"total_price"`           // yen
	Notes               string            `gorm:"type:text" json:"notes"`
	CancellationReason  string            `gorm:"type:text" json:"cancellation_reason"`
	CreatedAt           time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relations
	Customer            Customer            `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Staff               Staff               `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	ReservationMenus    []ReservationMenu   `gorm:"foreignKey:ReservationID" json:"reservation_menus,omitempty"`
	ReservationOptions  []ReservationOption `gorm:"foreignKey:ReservationID" json:"reservation_options,omitempty"`
}

type ReservationMenu struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ReservationID uuid.UUID `gorm:"type:uuid;not null;index" json:"reservation_id"`
	MenuID        uuid.UUID `gorm:"type:uuid;not null;index" json:"menu_id"`
	Quantity      int       `gorm:"not null;default:1" json:"quantity" validate:"required,min=1"`
	UnitPrice     int       `gorm:"not null" json:"unit_price"`
	TotalPrice    int       `gorm:"not null" json:"total_price"`
	
	// Relations
	Reservation   Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Menu          Menu        `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}

type ReservationOption struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ReservationID uuid.UUID `gorm:"type:uuid;not null;index" json:"reservation_id"`
	OptionID      uuid.UUID `gorm:"type:uuid;not null;index" json:"option_id"`
	Quantity      int       `gorm:"not null;default:1" json:"quantity" validate:"required,min=1"`
	UnitPrice     int       `gorm:"not null" json:"unit_price"`
	TotalPrice    int       `gorm:"not null" json:"total_price"`
	
	// Relations
	Reservation   Reservation `gorm:"foreignKey:ReservationID" json:"reservation,omitempty"`
	Option        Option      `gorm:"foreignKey:OptionID" json:"option,omitempty"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (rm *ReservationMenu) BeforeCreate(tx *gorm.DB) error {
	if rm.ID == uuid.Nil {
		rm.ID = uuid.New()
	}
	return nil
}

func (ro *ReservationOption) BeforeCreate(tx *gorm.DB) error {
	if ro.ID == uuid.Nil {
		ro.ID = uuid.New()
	}
	return nil
}

func (r *Reservation) TableName() string {
	return "reservations"
}

func (rm *ReservationMenu) TableName() string {
	return "reservation_menus"
}

func (ro *ReservationOption) TableName() string {
	return "reservation_options"
}