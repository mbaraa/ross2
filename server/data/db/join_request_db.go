package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// JoinRequestDB represents a CRUD db repo for jrs
type JoinRequestDB[T models.JoinRequest] struct {
	db *gorm.DB
}

// NewJoinRequestDB returns a new JoinRequestDB instance
func NewJoinRequestDB(db *gorm.DB) *JoinRequestDB[models.JoinRequest] {
	return &JoinRequestDB[models.JoinRequest]{db}
}

// CREATOR REPO

func (j *JoinRequestDB[T]) Add(jr *models.JoinRequest) error {
	return j.db.
		Create(&jr).
		Error
}

func (j *JoinRequestDB[T]) AddMany(jr []*models.JoinRequest) error {
	return j.db.
		Create(&jr).
		Error
}

func (j *JoinRequestDB[T]) GetDB() *gorm.DB {
	return j.db
}

// GETTER REPO

func (j *JoinRequestDB[T]) Exists(id uint) bool {
	return false
}

func (j *JoinRequestDB[T]) Get(id uint) (jr models.JoinRequest, err error) {
	return models.JoinRequest{}, errors.New("not implemented")
}

func (j *JoinRequestDB[T]) GetByConds(conds ...any) (jrs []models.JoinRequest, err error) {
	err = j.db.
		Model(new(models.JoinRequest)).
		Find(&jrs, conds[0], conds[1:]).
		Error

	return
}

func (j *JoinRequestDB[T]) GetAll() (jr []models.JoinRequest, err error) {
	return nil, errors.New("not implemented")
}

func (j *JoinRequestDB[T]) Count() (count int64, err error) {
	return 0, errors.New("not implemented")
}

// UPDATER REPO

func (j *JoinRequestDB[T]) Update(jr *models.JoinRequest, conds ...any) error {
	return errors.New("not implemented")
}

func (j *JoinRequestDB[T]) UpdateAll(jrs []*models.JoinRequest, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (j *JoinRequestDB[T]) Delete(jr models.JoinRequest, conds ...any) error {
	return j.db.
		Where("requester_id = ?", jr.RequesterID).
		Delete(&jr).
		Error
}

func (j *JoinRequestDB[T]) DeleteAll(conds ...any) error {
	return j.db.
		Where("true").
		Delete(new(models.JoinRequest)).
		Error
}
