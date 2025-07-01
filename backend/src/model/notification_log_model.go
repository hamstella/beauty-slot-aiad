package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type NotificationLog struct {
	ID            uuid.UUID          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Type          string             `gorm:"size:20;not null;index" json:"type" validate:"required,oneof=email sms push"`
	Recipient     string             `gorm:"size:255;not null" json:"recipient" validate:"required,max=255"`
	Subject       string             `gorm:"size:255" json:"subject" validate:"omitempty,max=255"`
	Message       string             `gorm:"type:text;not null" json:"message" validate:"required"`
	Status        NotificationStatus `gorm:"size:20;not null;default:pending" json:"status" validate:"required,oneof=pending sent failed"`
	ErrorMessage  string             `gorm:"type:text" json:"error_message"`
	ScheduledAt   *time.Time         `gorm:"index" json:"scheduled_at"`
	SentAt        *time.Time         `gorm:"index" json:"sent_at"`
	CreatedAt     time.Time          `gorm:"autoCreateTime;index:idx_notification_logs_created_at" json:"created_at"`
	UpdatedAt     time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

func (n *NotificationLog) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (n *NotificationLog) TableName() string {
	return "notification_logs"
}