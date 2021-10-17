package db

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
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

func (t *TeamDB) Add(team models.Team) error {
	err := t.db.
		Create(&team).
		Error
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		member.TeamID = team.ID
		_ = t.contRepo.Update(member)
	}

	return nil
}

// AddMany creates multiple teams, this method is only used when a contest's director
// generates teams for the teamless contestant and add them to it,
// on top of that, there's no fucking contest in the entire universe that would
// require any optimization for this code :)
func (t *TeamDB) AddMany(teams []models.Team) error {
	for ti := range teams {
		members := &teams[ti].Members
		for mi := range *members {
			(*members)[mi].TeamID = teams[ti].ID
			_ = t.contRepo.Update((*members)[mi])
		}
	}
	return t.db.
		Create(&teams).
		Error
}

// GETTER REPO

func (t *TeamDB) Exists(team models.Team) (bool, error) {
	res := t.db.First(&team)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (t *TeamDB) Get(team models.Team) (models.Team, error) {
	fetchedTeam := models.Team{}
	err := t.db.First(&fetchedTeam, "id = ?", team.ID).Error

	return fetchedTeam, err
}

func (t *TeamDB) GetAll() ([]models.Team, error) {
	count, err := t.Count()
	if err != nil {
		return nil, err
	}
	teams := make([]models.Team, count)
	err = t.db.Find(&teams).Error

	return teams, err
}

func (t *TeamDB) GetAllByContest(contest models.Contest) ([]models.Team, error) {
	teams := make([]models.Team, 0)

	err := t.db.
		Model(&contest).
		Association("Teams").
		Find(&teams)

	return teams, err
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
