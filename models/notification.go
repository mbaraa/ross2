package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// Notification represents a notification's fields
type Notification struct {
	gorm.Model
	ID      uint         `gorm:"column:id;primaryKey;autoIncrement"`
	UserID  uint         `gorm:"column:user_id"`
	Content string       `gorm:"column:content"`
	Seen    sql.NullBool `gorm:"column:seen"`
	SeenAt  time.Time    `gorm:"column:seen_at"`
}
