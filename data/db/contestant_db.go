package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// ContestantDB represents a CRUD db repo for contestants
type ContestantDB[T models.Contestant] struct {
	db *gorm.DB
}

// NewContestantDB returns a new ContestantDB instance
func NewContestantDB[T models.Contestant](db *gorm.DB) *ContestantDB[T] {
	return &ContestantDB[T]{db: db}
}

func (c *ContestantDB[T]) GetDB() *gorm.DB {
	return c.db
}

// CREATOR REPO

func (c *ContestantDB[T]) Add(contestant *models.Contestant) error {
	return c.db.
		Create(contestant).
		Error
}

func (c *ContestantDB[T]) AddMany(conts []*models.Contestant) error {
	return errors.New("not implemented")
}

// GETTER REPO

func (c *ContestantDB[T]) Exists(userID uint) bool {
	_, err := c.Get(userID)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (c *ContestantDB[T]) Get(userID uint) (fetchedContestant models.Contestant, err error) {
	err = c.db.
		Model(new(models.Contestant)).
		First(&fetchedContestant, "user_id = ?", userID).
		Error

	err = c.db.
		First(&fetchedContestant.User, "id = ?", userID).
		Error

	return
}

func (c *ContestantDB[T]) GetByConds(conds ...any) ([]models.Contestant, error) {
	return nil, errors.New("not implemented")
}

func (c *ContestantDB[T]) GetAll() ([]models.Contestant, error) {
	var (
		count       int64
		err         error
		contestants = make([]models.Contestant, count)
	)

	err = c.db.Find(&contestants).Error

	return contestants, err
}

func (c *ContestantDB[T]) Count() (int64, error) {
	var count int64
	err := c.db.
		Model(new(models.Contestant)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (c *ContestantDB[T]) Update(cont *models.Contestant, conds ...any) error {
	return c.db.
		Model(new(models.Contestant)).
		Where("user_id = ?", cont.User.ID).
		Updates(&cont).
		Error
}

func (c *ContestantDB[T]) UpdateAll(conts []*models.Contestant, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (c *ContestantDB[T]) Delete(contestant models.Contestant, conds ...any) error {
	return c.db.
		Where("user_id = ?", contestant.User.ID).
		Delete(&contestant).
		Error
}

func (c *ContestantDB[T]) DeleteAll(conds ...any) error {
	return c.db.
		Where("true").
		Delete(new(models.Contestant)).
		Error
}

