package db

import (
	"errors"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/strutils"
	"gorm.io/gorm"
)

// TeamDB represents a CRUD db repo for teams
type TeamDB struct {
	db       *gorm.DB
	contRepo data.ContestantUpdaterRepo
}

// NewTeamDB returns a new TeamDB instance
func NewTeamDB(db *gorm.DB, contRepo data.ContestantUpdaterRepo) *TeamDB {
	return &TeamDB{db: db, contRepo: contRepo}
}

// CREATOR REPO

func (t *TeamDB) Add(team *models.Team) error {
	return t.createTeam(team)
}

// AddMany creates multiple teams, this method is only used when a contest's director
// generates teams for the teamless contestants.
func (t *TeamDB) AddMany(teams []*models.Team) error {
	for _, team := range teams {
		_ = t.Add(team)
	}
	return nil
}

// createTeam finds a suitable joining id for the team in a very stupid way
func (t *TeamDB) createTeam(team *models.Team) error {
	for {
		team.JoinID = strutils.GetRandomString(3)

		err := t.db.
			Create(team).
			Error

		if err != nil {
			if !strings.Contains(err.Error(), "Duplicate") {
				return err
			}
			team.JoinID = strutils.GetRandomString(3)
		} else {
			return nil
		}
	}
}

// GETTER REPO

func (t *TeamDB) Exists(team models.Team) (bool, error) {
	res := t.db.First(&team)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (t *TeamDB) Get(team models.Team) (fetchedTeam models.Team, err error) {
	err = t.db.First(&fetchedTeam, "id = ?", team.ID).Error
	if err != nil {
		return
	}

	err = t.db.
		Model(fetchedTeam).
		Association("Contests").
		Find(&fetchedTeam.Contests)

	if err != nil {
		return
	}

	return
}

func (t *TeamDB) GetAll() (teams []models.Team, err error) {
	err = t.db.
		Find(&teams).
		Error
	return
}

func (t *TeamDB) GetAllByContest(contest models.Contest) ([]models.Team, error) {
	teams := make([]models.Team, 0)

	err := t.db.
		Model(&contest).
		Association("Teams").
		Find(&teams)

	return teams, err
}

func (t *TeamDB) GetByJoinID(joinID string) (team models.Team, err error) {
	err = t.db.
		Model(new(models.Team)).
		Where("join_id = ?", joinID).
		First(&team).
		Error
	return
}

func (t *TeamDB) Count() (int64, error) {
	var count int64
	err := t.db.
		Model(new(models.Team)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (t *TeamDB) Update(team models.Team) error {
	return t.db.
		Model(new(models.Team)).
		Where("id = ?", team.ID).
		Updates(&team).
		Error
}

// DELETER REPO

func (t *TeamDB) Delete(team models.Team) error {
	return t.db.
		Where("id = ?", team.ID).
		Delete(&team).
		Error
}

func (t *TeamDB) DeleteAll() error {
	return t.db.
		Where("true").
		Delete(new(models.Team)).
		Error
}
