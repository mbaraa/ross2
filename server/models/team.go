package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

// Team represents a team's fields
type Team struct {
	gorm.Model `json:"-"`
	ID         uint        `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string      `gorm:"column:name" json:"name"`
	LeaderId   uint        `gorm:"column:leader_id" json:"leader_id"` // not a foreign key to avoid cycling mess :)
	Leader     *Contestant `gorm:"-" json:"leader"`                   // using a pointer to avoid cycling :)
	JoinID     string      `gorm:"column:join_id;unique" json:"join_id"`

	Contests []Contest    `gorm:"many2many:register_teams"`                // a team may register in more than one contest
	Members  []Contestant `gorm:"many2many:register_teams" json:"members"` // each contestant has their team id :)
}

func (t *Team) AfterFind(db *gorm.DB) error {
	if t.ID <= 1 {
		return nil
	}

	err := db.
		Model(new(Contestant)).
		Find(&t.Leader, "user_id = ?", t.LeaderId).
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

func (t *Team) AfterCreate(db *gorm.DB) error {
	for i := range t.Members {
		t.Members[i].TeamID = t.ID
		t.Members[i].TeamlessContestID = math.MaxInt
		t.Members[i].TeamlessedAt = time.Now()

		err := db.
			Model(&t.Members[i]).
			Where("user_id = ?", t.Members[i].User.ID).
			Updates(&t.Members[i]).
			Error

		if err != nil {
			return err
		}
	}

	return nil
}

type JoinRequest struct {
	gorm.Model          `json:"-"`
	ID                  uint         `gorm:"column:id;primaryKey;autoIncrement"`
	RequesterID         uint         `gorm:"column:requester_id" json:"requester_id"`
	RequestedTeamJoinID string       `gorm:"column:requested_team_join_id" json:"requested_team_join_id"`
	RequestedTeamID     uint         `gorm:"column:requested_team_id" json:"requested_team_id"`
	RequestedTeam       Team         `gorm:"foreignkey:RequestedTeamID" json:"requested_team"`
	RequestedContestID  uint         `gorm:"column:req_contest_id" json:"requested_contest_id"`
	RequestMessage      string       `gorm:"column:message" json:"request_message"`
	NotificationID      uint         `gorm:"column:notification_id"`
	Notification        Notification `gorm:"foreignkey:NotificationID"`
}

func (j *JoinRequest) BeforeDelete(db *gorm.DB) error {
	return db.
		Model(new(Notification)).
		Where("id = ?", j.NotificationID).
		Delete(&j.Notification).
		Error
}

type RegisterTeam struct {
	gorm.Model `json:"-"`
	// ID           uint       `gorm:"primaryKey;autoIncrement"`
	ContestID    uint       `gorm:"column:contest_id" json:"-"`
	ContestName  string     `gorm:"column:contest_name" json:"contest_name"`
	TeamID       uint       `gorm:"column:team_id" json:"-"`
	Team         Team       `gorm:"foreignkey:TeamID" json:"team"`
	ContestantID uint       `gorm:"column:contestant_id" json:"-"`
	Contestant   Contestant `gorm:"foreignkey:ContestantID" json:"-"`
}

func (r *RegisterTeam) AfterFind(db *gorm.DB) error {
	var c Contest
	err := db.
		Model(new(Contest)).
		First(&c, "id = ?", r.ContestID).
		Error

	if err != nil {
		return err
	}

	err = db.
		Model(&r.Team).
		Find(&r.Team, "id = ?", r.TeamID).
		Error
	if err != nil {
		return err
	}

	err = db.
		Model(&r.Team).
		Association("Members").
		Find(&r.Team.Members)

	if err != nil {
		return err
	}

	r.ContestName = c.Name

	return nil
}
