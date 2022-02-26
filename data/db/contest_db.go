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

func (c *ContestDB) Add(contest *models.Contest) error {
	return c.db.
		Create(&contest).
		Error
}

// GETTER REPO

func (c *ContestDB) Exists(contest models.Contest) (bool, error) {
	res := c.db.First(&contest)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (c *ContestDB) Get(contest models.Contest) (fetchedContest models.Contest, err error) {
	err = c.db.
		First(&fetchedContest, "id = ?", contest.ID).
		Error
	if err != nil {
		return
	}

	fetchedContest.Teams, err = c.teamRepo.GetAllByContest(contest)
	if err != nil {
		return
	}

	err = c.db.
		Model(new(models.ParticipationConditions)).
		First(&fetchedContest.ParticipationConditions, "id = ?", fetchedContest.PCsID).
		Error
	if err != nil {
		return
	}

	err = c.db.
		Model(&fetchedContest).
		Association("Organizers").
		Find(&fetchedContest.Organizers)
	if err != nil {
		return
	}

	fetchedContest.TeamlessContestants, err = c.tlRepo.GetAllTeamLess(contest)

	return
}

func (c *ContestDB) GetAll() (contests []models.Contest, err error) {
	err = c.db.Find(&contests).Error
	return
}

func (c *ContestDB) GetAllByOrganizer(org models.Organizer) (contests []models.Contest, err error) {
	err = c.db.
		Model(&org).
		Association("Contests").
		Find(&contests)

	return
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
