package models

import (
	"time"

	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

// Contest represents a contest's fields
type Contest struct {
	gorm.Model
	ID                uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name              string        `gorm:"column:name" json:"name"`
	StartsAt          int64         `gorm:"-" json:"starts_at"`
	StartsAt2         time.Time     `gorm:"column:starts_at"` // weird ain't it? :)
	RegistrationEnds  int64         `gorm:"-" json:"registration_ends"`
	RegistrationEnds2 time.Time     `gorm:"column:registration_ends"` // weird ain't it? :)
	Duration          time.Duration `gorm:"column:duration" json:"duration"`
	Location          string        `gorm:"column:location" json:"location"`
	LogoPath          string        `gorm:"column:logo_path" json:"logo_path"`
	CoverPath         string        `gorm:"column:cover_path" json:"cover_path"`
	Description       string        `gorm:"column:description" json:"description"`
	TeamsHidden       bool          `gorm:"column:teams_hidden" json:"teams_hidden"`

	ParticipationConditions ParticipationConditions `gorm:"foreignkey:PCsID" json:"participation_conditions"` // the conditions should be a part of the contest instance 🤷‍♂️
	PCsID                   uint                    `gorm:"column:pc_id"`

	AllowedContestantsCount uint `gorm:"allowed_contestant_count" json:"allowed_contestant_count"`
	CurrentContestantsCount uint `gorm:"current_contestant_count" json:"current_contestant_count"`

	Teams               []Team       `gorm:"many2many:register_teams;" json:"teams"`
	Organizers          []Organizer  `gorm:"many2many:register_contest;" json:"organizers"`
	TeamlessContestants []Contestant `gorm:"-" json:"teamless_contestants"`
}

func (c *Contest) BeforeCreate(db *gorm.DB) error {
	c.StartsAt2 = time.UnixMilli(c.StartsAt)
	c.RegistrationEnds2 = time.UnixMilli(c.RegistrationEnds)
	c.Duration *= 1e9 * 60
	return nil
}

func (c *Contest) AfterFind(db *gorm.DB) error {
	c.StartsAt = c.StartsAt2.UnixMilli()
	c.RegistrationEnds = c.RegistrationEnds2.UnixMilli()
	c.ParticipationConditions.MajorsNames = getMajors(c.ParticipationConditions.Majors)
	return nil
}

func getMajors(majors enums.Major) []string {
	majorsTexts := make([]string, 0)
	for major := enums.MajorSoftwareEngineering; major <= enums.MajorCyberSecurity; major <<= 1 {
		if majors&major != 0 {
			majorsTexts = append(majorsTexts, major.String())
		}
	}

	return majorsTexts
}

// ParticipationConditions represents the conditions needed to participate in a contest
type ParticipationConditions struct {
	gorm.Model
	ID             uint        `gorm:"column:id;primaryKey;autoIncrement"`
	Majors         enums.Major `gorm:"column:majors;type:uint" json:"majors"`
	MajorsNames    []string    `gorm:"-" json:"majors_names"`
	MinTeamMembers uint        `gorm:"column:min_team_members" json:"min_team_members"`
	MaxTeamMembers uint        `gorm:"column:max_team_members" json:"max_team_members"`
}
