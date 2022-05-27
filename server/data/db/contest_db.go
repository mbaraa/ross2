package db

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// ContestDB represents a CRUD db repo for contests
type ContestDB[T models.Contest] struct {
	db       *gorm.DB
	teamRepo data.CRUDRepo[models.Team]
}

// NewContestDB returns a new ContestDB instance
func NewContestDB(db *gorm.DB, teamRepo data.CRUDRepo[models.Team]) *ContestDB[models.Contest] {
	return &ContestDB[models.Contest]{
		db:       db,
		teamRepo: teamRepo,
	}
}

func (c *ContestDB[T]) GetDB() *gorm.DB {
	return c.db
}

// CREATOR REPO

func (c *ContestDB[T]) Add(contest *models.Contest) error {
	return c.db.
		Create(&contest).
		Error
}

func (c *ContestDB[T]) AddMany(contests []*models.Contest) error {
	return errors.New("not implemented")
}

// GETTER REPO

func (c *ContestDB[T]) Exists(id uint) bool {
	res := c.db.First(&models.Contest{ID: id})
	return !errors.Is(res.Error, gorm.ErrRecordNotFound)
}

func (c *ContestDB[T]) Get(id uint) (fetchedContest models.Contest, err error) {
	err = c.db.
		First(&fetchedContest, "id = ?", id).
		Error
	if err != nil {
		return
	}

	err = c.db.
		Model(&fetchedContest).
		Association("Teams").
		Find(&fetchedContest.Teams)
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

	err = c.db.
		Model(new(models.Contestant)).
		Find(&fetchedContest.TeamlessContestants, "teamless_contest_id = ?", fetchedContest.ID).
		Error

	return
}

func (c *ContestDB[T]) GetByConds(conds ...any) ([]models.Contest, error) {
	return nil, errors.New("not implemented")
}

func (c *ContestDB[T]) GetAll() (contests []models.Contest, err error) {
	err = c.db.Find(&contests).Error
	return
}

func (c *ContestDB[T]) Count() (int64, error) {
	var count int64
	err := c.db.
		Model(new(models.Contest)).
		Count(&count).
		Error

	return count, err
}

// The Updater & Deleter Repos' Methods doesn't modify data of the teams table :)

// UPDATER REPO

func (c *ContestDB[T]) Update(contest *models.Contest, conds ...any) error {
	if conds != nil {
		return errors.New("can't use conditions now")
	}

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

func (c *ContestDB[T]) UpdateAll(contests []*models.Contest, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (c *ContestDB[T]) Delete(contest models.Contest, conds ...any) error {
	if conds != nil {
		return errors.New("can't use conditions now")
	}

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

func (c *ContestDB[T]) DeleteAll(conds ...any) error {
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
