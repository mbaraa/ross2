package models

import (
	"gorm.io/gorm"
	"time"
)

// Contestant represents a contestant's fields
type Contestant struct {
	User
	UniversityID string `gorm:"column:university_id" json:"university_id"`
	Team         Team   `gorm:"foreignkey:TeamID" json:"team"` // big surprise, a contestant gets their contests from here :)
	TeamID       uint   `gorm:"column:team_id"`
	Major        Major  `gorm:"column:major;type:uint" json:"major"`
	MajorName    string `gorm:"-" json:"major_name"`

	TeamlessedAt      time.Time `gorm:"timelessed_at"`
	TeamlessContestID uint      `gorm:"column:teamless_contest_id"`
}

func (c *Contestant) AfterFind(db *gorm.DB) error {
	err := db.
		First(&c.ContactInfo, "id = ?", c.ContactInfoID).
		Error

	c.MajorName = majorText[c.Major]

	return err
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
