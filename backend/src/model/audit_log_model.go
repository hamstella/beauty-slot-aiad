package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Table      string    `gorm:"size:50;not null;index;column:table_name" json:"table_name" validate:"required,max=50"`
	RecordID   uuid.UUID `gorm:"type:uuid;not null;index" json:"record_id" validate:"required"`
	Action     string    `gorm:"size:20;not null;index" json:"action" validate:"required,oneof=CREATE UPDATE DELETE"`
	OldValues  string    `gorm:"type:jsonb" json:"old_values"`
	NewValues  string    `gorm:"type:jsonb" json:"new_values"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	IPAddress  string    `gorm:"size:45" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index:idx_audit_logs_created_at" json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (a *AuditLog) TableName() string {
	return "audit_logs"
}