package models

import (
	"time"

	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/multiavatar"
	"gorm.io/gorm"
)

// User represents a general user :)
type User struct {
	gorm.Model
	ID            uint                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email         string              `gorm:"column:email;unique" json:"email"`
	Name          string              `gorm:"column:name" json:"name"`
	AvatarURL     string              `gorm:"column:avatar_url" json:"avatar_url"`
	ProfileStatus enums.ProfileStatus `gorm:"profile_status" json:"profile_status"`

	UserType     enums.UserType `gorm:"column:user_type" json:"user_type_base"`
	UserTypeText []string       `gorm:"-" json:"user_type"`

	ContactInfo   ContactInfo `gorm:"foreignkey:ContactInfoID" json:"contact_info"`
	ContactInfoID uint        `gorm:"column:contact_info_id"`

	// used only for organizers and contestants :)
	AttendedAt        time.Time `gorm:"column:attended_at" json:"attended_at"`
	AttendedContestID uint      `gorm:"column:attended_contest_id" json:"attended_contest_id"`
}

func (u *User) AfterFind(db *gorm.DB) error {
	u.UserTypeText = u.UserType.GetTypes()
	return nil
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	u.AvatarURL = multiavatar.GetAvatarURL()
	u.ProfileStatus = enums.ProfileStatusFresh
	u.ContactInfo = ContactInfo{FacebookURL: "/"}

	return nil
}

// ContactInfo represents a user's(any user on Ross) fields
type ContactInfo struct {
	gorm.Model
	FacebookURL         string `gorm:"column:facebook_url" json:"facebook_url"`
	MicrosoftTeamsEmail string `gorm:"column:msteams_email" json:"msteams_email"`
	TelegramNumber      string `gorm:"column:telegram_number" json:"telegram_number"`
}
