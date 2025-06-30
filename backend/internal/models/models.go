package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Customer 顧客
type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	Phone     string    `gorm:"type:varchar(20);not null;unique" json:"phone" validate:"required,e164"`
	Email     string    `gorm:"type:varchar(255);not null;unique" json:"email" validate:"required,email"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Reservations []Reservation `json:"reservations,omitempty"`
}

// Staff スタッフ
type Staff struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Labels       []Label       `gorm:"many2many:staff_labels" json:"labels,omitempty"`
	Shifts       []Shift       `json:"shifts,omitempty"`
	Reservations []Reservation `json:"reservations,omitempty"`
}

// Menu メニュー
type Menu struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"type:varchar(200);not null" json:"name" validate:"required,min=1,max=200"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes" validate:"required,min=1"`
	Price           int       `gorm:"not null" json:"price" validate:"required,min=0"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Labels       []Label       `gorm:"many2many:menu_labels" json:"labels,omitempty"`
	Reservations []Reservation `json:"reservations,omitempty"`
}

// Option オプション
type Option struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name                string    `gorm:"type:varchar(200);not null" json:"name" validate:"required,min=1,max=200"`
	AddDurationMinutes  int       `gorm:"not null" json:"add_duration_minutes" validate:"min=0"`
	AddPrice            int       `gorm:"not null" json:"add_price" validate:"min=0"`
	IsActive            bool      `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt           time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Reservations []Reservation `gorm:"many2many:reservation_options" json:"reservations,omitempty"`
}

// Label ラベル
type Label struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;unique" json:"name" validate:"required,min=1,max=100"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Staff []Staff `gorm:"many2many:staff_labels" json:"staff,omitempty"`
	Menus []Menu  `gorm:"many2many:menu_labels" json:"menus,omitempty"`
}

// Shift シフト
type Shift struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	StaffID   uuid.UUID `gorm:"type:uuid;not null" json:"staff_id" validate:"required"`
	WorkDate  time.Time `gorm:"type:date;not null" json:"work_date" validate:"required"`
	StartTime time.Time `gorm:"type:time;not null" json:"start_time" validate:"required"`
	EndTime   time.Time `gorm:"type:time;not null" json:"end_time" validate:"required"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Staff Staff `json:"staff,omitempty"`
}

// Reservation 予約
type Reservation struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerID           uuid.UUID `gorm:"type:uuid;not null" json:"customer_id" validate:"required"`
	StaffID              uuid.UUID `gorm:"type:uuid;not null" json:"staff_id" validate:"required"`
	MenuID               uuid.UUID `gorm:"type:uuid;not null" json:"menu_id" validate:"required"`
	StartAt              time.Time `gorm:"type:timestamp with time zone;not null" json:"start_at" validate:"required"`
	EndAt                time.Time `gorm:"type:timestamp with time zone;not null" json:"end_at" validate:"required"`
	Status               string    `gorm:"type:varchar(20);default:pending" json:"status" validate:"oneof=pending confirmed in_progress completed cancelled"`
	TotalPrice           int       `gorm:"not null" json:"total_price" validate:"min=0"`
	TotalDurationMinutes int       `gorm:"not null" json:"total_duration_minutes" validate:"min=1"`
	CreatedAt            time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt            time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relations
	Customer Customer `json:"customer,omitempty"`
	Staff    Staff    `json:"staff,omitempty"`
	Menu     Menu     `json:"menu,omitempty"`
	Options  []Option `gorm:"many2many:reservation_options" json:"options,omitempty"`
}

// AuditLog 監査ログ
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	EntityType string    `gorm:"type:varchar(50);not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	Action     string    `gorm:"type:varchar(20);not null" json:"action" validate:"oneof=CREATE UPDATE DELETE"`
	OldValues  *string   `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues  *string   `gorm:"type:jsonb" json:"new_values,omitempty"`
	IPAddress  *string   `gorm:"type:inet" json:"ip_address,omitempty"`
	CreatedAt  time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
}

// NotificationLog 通知ログ
type NotificationLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReservationID    *uuid.UUID `gorm:"type:uuid" json:"reservation_id,omitempty"`
	NotificationType string     `gorm:"type:varchar(50);not null" json:"notification_type"`
	Channel          string     `gorm:"type:varchar(20);not null" json:"channel" validate:"oneof=email sms push"`
	Recipient        string     `gorm:"type:varchar(255);not null" json:"recipient"`
	Status           string     `gorm:"type:varchar(20);default:pending" json:"status" validate:"oneof=pending sent failed delivered"`
	MessageContent   *string    `gorm:"type:text" json:"message_content,omitempty"`
	SentAt           *time.Time `gorm:"type:timestamp with time zone" json:"sent_at,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`

	// Relations
	Reservation *Reservation `json:"reservation,omitempty"`
}

// BeforeCreate フック（UUID自動生成）
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (o *Option) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (s *Shift) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (n *NotificationLog) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}