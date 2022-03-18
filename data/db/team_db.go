package db

import (
	"errors"
	"strings"

	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/strutils"
	"gorm.io/gorm"
)

// TeamDB represents a CRUD db repo for teams
type TeamDB[T models.Team, T2 any] struct {
	db *gorm.DB
}

// NewTeamDB returns a new TeamDB instance
func NewTeamDB[T models.Team, T2 any](db *gorm.DB) *TeamDB[T, T2] {
	return &TeamDB[T, T2]{db: db}
}

func (t *TeamDB[T, T2]) GetDB() *gorm.DB {
	return t.db
}

// CREATOR REPO

func (t *TeamDB[T, T2]) Add(team *models.Team) error {
	return t.createTeam(team)
}

// AddMany creates multiple teams, this method is only used when a contest's director
// generates teams for the teamless contestants.
func (t *TeamDB[T, T2]) AddMany(teams []*models.Team) error {
	for _, team := range teams {
		_ = t.Add(team)
	}
	return nil
}

// createTeam finds a suitable joining id for the team in a very stupid way
func (t *TeamDB[T, T2]) createTeam(team *models.Team) error {
	if strutils.IsBadWord(team.Name) {
		return errors.New("team name can't contain bad words 🙂")
	}
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

func (t *TeamDB[T, T2]) Exists(id uint) bool {
	_, err := t.Get(id)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (t *TeamDB[T, T2]) Get(id uint) (fetchedTeam models.Team, err error) {
	err = t.db.First(&fetchedTeam, "id = ?", id).Error
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

func (t *TeamDB[T, T2]) GetByAssociation(contest any) ([]models.Team, error) {
	teams := make([]models.Team, 0)

	err := t.db.
		Model(&contest).
		Association("Teams").
		Find(&teams)

	return teams, err
}

func (t *TeamDB[T, T2]) GetByConds(conds ...any) (teams []models.Team, err error) {
	if (len(conds) < 2) {
		return nil, errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	err = t.db.
		Model(new(models.Team)).
		Where(conds[0], conds[1:]).
		Find(&teams).
		Error
	return
}

func (t *TeamDB[T, T2]) GetAll() (teams []models.Team, err error) {
	return nil, errors.New("not implemented")
}

func (t *TeamDB[T, T2]) Count() (int64, error) {
	var count int64
	err := t.db.
		Model(new(models.Team)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (t *TeamDB[T, T2]) Update(team *models.Team, conds ...any) error {
	return t.db.
		Model(new(models.Team)).
		Where("id = ?", team.ID).
		Updates(&team).
		Error
}

func (t *TeamDB[T, T2]) UpdateAll(teams []*models.Team, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (t *TeamDB[T, T2]) Delete(team models.Team, conds ...any) error {
	return t.db.
		Where("id = ?", team.ID).
		Delete(&team).
		Error
}

func (t *TeamDB[T, T2]) DeleteAll(conds ...any) error {
	return t.db.
		Where("true").
		Delete(new(models.Team)).
		Error
}
