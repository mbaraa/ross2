package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a session's fields
type Session struct {
	ID        string `gorm:"column:id;primaryKey"`
	CreatedAt time.Time
	UserID    uint `gorm:"column:user_id"`
}

func (s *Session) BeforeCreate(db *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
