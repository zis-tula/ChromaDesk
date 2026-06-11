package models

import "time"

type AuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AdminID       uint      `gorm:"index" json:"admin_id"`
	AdminUsername string    `gorm:"size:100;index" json:"admin_username"`
	Action        string    `gorm:"size:50;index" json:"action"`
	ResourceType  string    `gorm:"size:50" json:"resource_type"`
	ResourceID    string    `gorm:"size:100" json:"resource_id"`
	Details       string    `gorm:"type:text" json:"details"`
	IP            string    `gorm:"size:50" json:"ip"`
	CreatedAt     time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
