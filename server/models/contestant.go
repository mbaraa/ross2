package models

import (
	"time"

	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

// Contestant represents a contestant's fields
type Contestant struct {
	gorm.Model `json:"-"`
	User       User `gorm:"foreignkey:UserID" json:"user"`
	UserID     uint `gorm:"column:user_id" json:"user_id"`

	Teams     []Team      `gorm:"many2many:register_teams" json:"-"`
	Major     enums.Major `gorm:"column:major;type:uint" json:"major"`
	MajorName string      `gorm:"-" json:"major_name"`

	TeamlessedAt      time.Time `gorm:"column:teamlessed_at" json:"teamlessed_at"`
	TeamlessContestID uint      `gorm:"column:teamless_contest_id" json:"teamless_contest_id"`

	Gender                     bool `gorm:"column:gender" json:"gender"`
	ParticipateWithOtherGender bool `gorm:"column:participate_with_other" json:"participate_with_other"`
}

func (c *Contestant) AfterFind(db *gorm.DB) error {
	c.MajorName = c.Major.String()

	err := db.
		Model(new(User)).
		First(&c.User, "id = ?", c.UserID).
		Error

	if err != nil {
		return nil
	}

	return db.
		First(&c.User.ContactInfo, "id = ?", c.User.ContactInfoID).
		Error
}
