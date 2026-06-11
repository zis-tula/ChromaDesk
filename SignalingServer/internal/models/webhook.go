package models

import "time"

type Webhook struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:100;not null" json:"name"`
	URL           string     `gorm:"size:500;not null" json:"url"`
	Secret        string     `gorm:"size:200" json:"secret,omitempty"`
	Events        string     `gorm:"type:text;not null" json:"events"`
	Enabled       bool       `gorm:"default:true" json:"enabled"`
	LastTriggered *time.Time `json:"last_triggered"`
	LastStatus    int        `json:"last_status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}
