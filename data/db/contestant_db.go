package db

import (
	"errors"
	"time"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// ContestantDB represents a CRUD db repo for contestants
type ContestantDB struct {
	db *gorm.DB
}

// NewContestantDB returns a new ContestantDB instance
func NewContestantDB(db *gorm.DB) *ContestantDB {
	return &ContestantDB{db: db}
}

// CREATOR REPO

func (c *ContestantDB) Add(contestant *models.Contestant) error {
	return c.db.
		Create(contestant).
		Error
}

// GETTER REPO

func (c *ContestantDB) Exists(contestant models.Contestant) (bool, error) {
	res := c.db.First(&contestant)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

// TODO
// generalize the get methods :_

func (c *ContestantDB) Get(contestant models.Contestant) (fetchedContestant models.Contestant, err error) {
	err = c.db.
		First(&fetchedContestant, "id = ?", contestant.ID).
		Error

	if err != nil {
		return models.Contestant{}, err
	}

	err = c.db.
		First(&fetchedContestant.User.ContactInfo, "id = ?", fetchedContestant.User.ContactInfoID).
		Error

	return
}

func (c *ContestantDB) GetByEmail(email string) (models.Contestant, error) {
	var (
		fetchedContestant models.Contestant
		err               error
	)

	err = c.db.
		First(&fetchedContestant, "email = ?", email).
		Error

	return fetchedContestant, err
}

func (c *ContestantDB) GetAll() ([]models.Contestant, error) {
	var (
		count       int64
		err         error
		contestants = make([]models.Contestant, count)
	)

	err = c.db.Find(&contestants).Error

	return contestants, err
}

func (c *ContestantDB) Count() (int64, error) {
	var count int64
	err := c.db.
		Model(new(models.Contestant)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (c *ContestantDB) Update(cont models.Contestant) error {
	return c.db.
		Model(new(models.Contestant)).
		Where("id = ?", cont.ID).
		Updates(&cont).
		Error
}

// DELETER REPO

func (c *ContestantDB) Delete(contestant models.Contestant) error {
	return c.db.
		Where("id = ?", contestant.ID).
		Delete(&contestant).
		Error
}

func (c *ContestantDB) DeleteAll() error {
	return c.db.
		Where("true").
		Delete(new(models.Contestant)).
		Error
}

// TEAMLESS REPO

func (c *ContestantDB) AddTeamLess(contestant models.Contestant, contest models.Contest) error {
	contestant.TeamlessedAt = time.Now()
	contestant.TeamlessContestID = contest.ID

	return c.Update(contestant)
}

func (c *ContestantDB) GetAllTeamLess(contest models.Contest) ([]models.Contestant, error) {
	tls := make([]models.Contestant, 0)

	err := c.db.
		Model(new(models.Contestant)).
		Find(&tls, "teamless_contest_id = ?", contest.ID).
		Error

	return tls, err
}

func (c *ContestantDB) RegisterInTeam(contestant models.Contestant, team models.Team) error {
	contestant.TeamlessContestID = 0
	contestant.TeamlessedAt = time.Time{}
	contestant.TeamID = team.ID

	return c.Update(contestant)
}
