package models

import (
	"gorm.io/gorm"
)

// Team represents a team's fields
type Team struct {
	gorm.Model
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	LeaderId uint   `gorm:"column:leader_id" json:"leader_id"` // not a foreign key to avoid cycling mess :)

	Contests []Contest `gorm:"many2many:register_teams;"` // a team may register in more than one contest

	Members []Contestant `gorm:"-" json:"members"` // each contestant has their team id :)
}

func (t *Team) AfterFind(db *gorm.DB) error {
	err := db.
		Model(new(Contest)).
		Find(&t.Contests).
		Error

	if err != nil {
		return err
	}

	return db.
		Model(new(Contestant)).
		Find(&t.Members, "team_id = ?", t.ID).
		Error
}

func (t *Team) BeforeDelete(db *gorm.DB) error {
	return db.
		Model(new(Contestant)).
		Where("team_id = ?", t.ID).
		Update("team_id", 0).
		Error
}

type JoinRequest struct {
	gorm.Model
	ID              uint         `gorm:"column:id;primaryKey;autoIncrement"`
	RequesterID     uint         `gorm:"column:requester_id"`
	Requester       Contestant   `gorm:"foreignkey:RequesterID"`
	RequestedTeamID uint         `gorm:"column:req_team_id"`
	RequestedTeam   Team         `gorm:"foreignkey:RequestedTeamID"`
	RequestMessage  string       `gorm:"column:message"`
	NotificationID  uint         `gorm:"column:notification_id"`
	Notification    Notification `gorm:"foreignkey:NotificationID"`
}

func (j *JoinRequest) BeforeDelete(db *gorm.DB) error {
	return db.
		Model(new(Notification)).
		Where("id = ?", j.NotificationID).
		Error
}
