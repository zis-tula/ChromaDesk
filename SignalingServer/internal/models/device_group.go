package models

import "time"

type DeviceGroup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Color       string    `gorm:"size:20" json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DeviceGroup) TableName() string {
	return "device_groups"
}

type DeviceGroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"not null;uniqueIndex:idx_group_device" json:"group_id"`
	DeviceID  string    `gorm:"size:20;not null;uniqueIndex:idx_group_device" json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (DeviceGroupMember) TableName() string {
	return "device_group_members"
}
