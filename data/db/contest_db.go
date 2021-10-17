package db

import (
	"errors"
	"github.com/mbaraa/ross2/data"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// ContestDB represents a CRUD db repo for contests
type ContestDB struct {
	db       *gorm.DB
	teamRepo data.TeamGetterRepo
	tlRepo   data.TeamlessCRUDRepo
}

// NewContestDB returns a new ContestDB instance
func NewContestDB(db *gorm.DB, teamRepo data.TeamGetterRepo, tl data.TeamlessCRUDRepo) *ContestDB {
	return &ContestDB{
		db:       db,
		teamRepo: teamRepo,
		tlRepo:   tl,
	}
}

// CREATOR REPO

func (c *ContestDB) Add(contest models.Contest) error {
	return c.db.
		Create(&contest).
		Error
}

// GETTER REPO

func (c *ContestDB) Exists(contest models.Contest) (bool, error) {
	res := c.db.First(&contest)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (c *ContestDB) Get(contest models.Contest) (models.Contest, error) {
	var (
		fetchedContest models.Contest
		err            error
	)

	err = c.db.
		First(&fetchedContest, "id = ?", contest.ID).
		Error
	if err != nil {
		return models.Contest{}, err
	}

	fetchedContest.Teams, err = c.teamRepo.GetAllByContest(contest)
	if err != nil {
		return models.Contest{}, err
	}

	fetchedContest.TeamlessContestants, err = c.tlRepo.GetAllTeamLess(contest)

	return fetchedContest, err
}

func (c *ContestDB) GetAll() ([]models.Contest, error) {
	count, err := c.Count()
	if err != nil {
		return nil, err
	}

	contests := make([]models.Contest, count)
	err = c.db.Find(&contests).Error

	return contests, err
}

func (c *ContestDB) GetAllByOrganizer(org models.Organizer) ([]models.Contest, error) {
	contests := make([]models.Contest, 0)
	err := c.db.
		Model(&org).
		Association("Contests").
		Find(&contests)

	return contests, err
}

func (c *ContestDB) Count() (int64, error) {
	var count int64
	err := c.db.
		Model(new(models.Contest)).
		Count(&count).
		Error

	return count, err
}

// The Updater & Deleter Repos' Methods doesn't modify data of the teams table :)

// UPDATER REPO

func (c *ContestDB) Update(contest models.Contest) error {
	err := c.db.
		Model(new(models.Contest)).
		Where("id = ?", contest.ID).
		Updates(&contest).
		Error
	if err != nil {
		return err
	}

	return c.db.
		Model(new(models.ParticipationConditions)).
		Where("id = ?", contest.PCsID).
		Updates(&contest.ParticipationConditions).
		Error // error handling goes brr
}

// DELETER REPO

func (c *ContestDB) Delete(contest models.Contest) error {
	err := c.db.
		Model(new(models.Contest)).
		Where("id = ?", contest.ID).
		Delete(&contest).
		Error
	if err != nil {
		return err
	}

	return c.db.
		Model(new(models.ParticipationConditions)).
		Where("id = ?", contest.PCsID).
		Delete(&contest.ParticipationConditions).
		Error // error handling goes brr
}

func (c *ContestDB) DeleteAll() error {
	err := c.db.
		Where("true").
		Delete(new(models.Contest)).
		Error
	if err != nil {
		return err
	}

	return c.db.
		Where("true").
		Delete(new(models.ParticipationConditions)).
		Error // error handling goes brr
}
