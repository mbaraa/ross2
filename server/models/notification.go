package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification represents a notification's fields
type Notification struct {
	gorm.Model `json:"-"`
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	Content    string    `gorm:"column:content" json:"content"`
	Seen       bool      `gorm:"column:seen" json:"seen"`
	SeenAt     time.Time `gorm:"column:seen_at" json:"seen_at"`
}
