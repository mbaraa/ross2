package models

import (
	"gorm.io/gorm"
)

// Organizer represents a contest's organizer
type Organizer struct {
	gorm.Model
	ID     uint `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	User   User `gorm:"foreignkey:UserID" json:"user"`
	UserID uint `gorm:"column:user_id" json:"user_id"`

	DirectorID uint       `gorm:"column:director_id"`
	Director   *Organizer `gorm:"foreignkey:DirectorID" json:"director"`

	Contests []Contest `gorm:"many2many:organize_contests;" json:"contests"`
}

func (o *Organizer) AfterFind(db *gorm.DB) error {
	err := db.
		Model(new(User)).
		First(&o.User, "id = ?", o.UserID).
		Error

	if err != nil {
		return err
	}

	return db.
		Model(new(ContactInfo)).
		First(&o.User.ContactInfo, "id = ?", o.User.ContactInfoID).
		Error
}
