package models

import (
	"gorm.io/gorm"
)

// Admin represents an administrator's fields
type Admin struct {
	gorm.Model `json:"-"`
	User       User `gorm:"foreignkey:UserID" json:"user"`
	UserID     uint `gorm:"column:user_id" json:"user_id"`
}

func (a *Admin) AfterFind(db *gorm.DB) error {
	return db.
		First(&a.User.ContactInfo, "id = ?", a.User.ContactInfoID).
		Error
}
