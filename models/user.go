package models

import (
	"gorm.io/gorm"
)

// this was supposed to be the super type of Organizer and Contestant util it fucked things up with the front end :)
//
//// User represents a general user :)
//type User struct {
//	gorm.Model
//	ID              uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
//	Email           string `gorm:"column:email" json:"email"`
//	Name            string `gorm:"column:name" json:"name"`
//	AvatarURL       string `gorm:"column:avatar_url" json:"avatar_url"`
//	ProfileFinished bool   `gorm:"profile_finished" json:"profile_finished"`
//
//	ContactInfo   ContactInfo `gorm:"foreignkey:ContactInfoID" json:"contact_info"`
//	ContactInfoID uint        `gorm:"column:contact_info_id"`
//}

// ContactInfo represents a user's(any user on Ross) fields
type ContactInfo struct {
	gorm.Model
	FacebookURL    string `gorm:"column:facebook_url" json:"facebook_url"`
	WhatsappNumber string `gorm:"column:whatsapp_number" json:"whatsapp_number"`
	TelegramNumber string `gorm:"column:telegram_number" json:"telegram_number"`
}
