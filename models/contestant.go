package models

import (
	"time"

	"gorm.io/gorm"
)

// Contestant represents a contestant's fields
type Contestant struct {
	gorm.Model
	ID              uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email           string `gorm:"column:email" json:"email"`
	Name            string `gorm:"column:name" json:"name"`
	AvatarURL       string `gorm:"column:avatar_url" json:"avatar_url"`
	ProfileFinished bool   `gorm:"profile_finished" json:"profile_finished"`

	ContactInfo   ContactInfo `gorm:"foreignkey:ContactInfoID" json:"contact_info"`
	ContactInfoID uint        `gorm:"column:contact_info_id"`

	UniversityID string `gorm:"column:university_id" json:"university_id"`
	Team         Team   `gorm:"foreignkey:TeamID" json:"team"` // big surprise, a contestant gets their contests from here :)
	TeamID       uint   `gorm:"column:team_id" json:"team_id"`
	Major        Major  `gorm:"column:major;type:uint" json:"major"`
	MajorName    string `gorm:"-" json:"major_name"`

	TeamlessedAt      time.Time `gorm:"column:teamlessed_at" json:"teamlessed_at"`
	TeamlessContestID uint      `gorm:"column:teamless_contest_id" json:"teamless_contest_id"`

	Gender                     bool `gorm:"column:gender" json:"gender"`
	ParticipateWithOtherGender bool `gorm:"column:participate_with_other" json:"participate_with_other"`
}

func (c *Contestant) AfterFind(db *gorm.DB) error {
	c.MajorName = majorText[c.Major]
	return nil
}

func (c *Contestant) BeforeCreate(db *gorm.DB) error {
	c.TeamID = 1
	c.Team = Team{ID: 1}
	return nil
}

// ContestantSortable is just a sortable by creation date Contestant slice
type ContestantSortable []Contestant

func (t ContestantSortable) Less(i, j int) bool {
	return t[i].TeamlessedAt.Before(t[j].TeamlessedAt)
}

func (t ContestantSortable) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ContestantSortable) Len() int {
	return len(t)
}
