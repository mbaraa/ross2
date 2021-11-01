package models

import (
	"time"

	"gorm.io/gorm"
)

// Contest represents a contest's fields
type Contest struct {
	gorm.Model
	ID          uint          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string        `gorm:"column:name" json:"name"`
	StartsAt    time.Time     `gorm:"column:starts_at" json:"starts_at"`
	Duration    time.Duration `gorm:"column:duration" json:"duration"`
	Location    string        `gorm:"column:location" json:"location"`
	LogoPath    string        `gorm:"column:logo_path" json:"logo_path"`
	CoverPath   string        `gorm:"column:cover_path" json:"cover_path"`
	Description string        `gorm:"column:description" json:"description"`

	ParticipationConditions ParticipationConditions `gorm:"foreignkey:PCsID" json:"participation_conditions"` // the conditions should be a part of the contest instance 🤷‍♂️
	PCsID                   uint                    `gorm:"column:pc_id"`

	Teams []Team `gorm:"many2many:register_teams;" json:"teams"`

	Organizers []Organizer `gorm:"many2many:register_contest;" json:"organizers"`

	TeamlessContestants []Contestant `gorm:"-" json:"teamless_contestants"`
}

func (c *Contest) AfterFind(db *gorm.DB) error {
	err := db.
		First(&c.ParticipationConditions, "id = ?", c.PCsID).
		Error
	if err != nil {
		return err
	}

	err = db.
		Model(new(Contest)).
		Association("Organizers").
		Find(&c.Organizers)

	c.ParticipationConditions.MajorsNames = getMajors(c.ParticipationConditions.Majors)

	return err
}

// this method is left as a run away solution in case that GORM deletes organizers associated with a contest :\
// func (c *Contest) BeforeDelete(db *gorm.DB) error {
// 	for _, org := range c.Organizers {
// 		for i, cont := range org.Contests {
// 			if cont.ID == c.ID {
// 				org.Contests = append(org.Contests[:i], org.Contests[i-1:]...)
// 			}
// 		}
//
// 		db.Model(new(Organizer)).
// 			Where("id = ?", org.ID).
// 			Update(&org)
// 	}
// 	c.Organizers = []Organizer{}
//
// 	return
// }

func getMajors(majors Major) []string {
	majorsTexts := make([]string, 0)
	for i := 0; i <= 63; i++ {
		major := Major(1 << i)
		if majors&major != 0 {
			majorsTexts = append(majorsTexts, majorText[major])
		}
	}

	return majorsTexts
}

// ParticipationConditions represents the conditions needed to participate in a contest
type ParticipationConditions struct {
	gorm.Model
	ID             uint     `gorm:"column:id;primaryKey;autoIncrement"`
	Majors         Major    `gorm:"column:majors;type:uint" json:"majors"`
	MajorsNames    []string `gorm:"-" json:"majors_names"`
	MinTeamMembers uint     `gorm:"column:min_team_members" json:"min_team_members"`
	MaxTeamMembers uint     `gorm:"column:max_team_members" json:"max_team_members"`
}
