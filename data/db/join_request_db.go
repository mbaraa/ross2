package db

import (
	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// JoinRequestDB represents a CRUD db repo for jrs
type JoinRequestDB struct {
	db *gorm.DB
}

// NewJoinRequestDB returns a new JoinRequestDB instance
func NewJoinRequestDB(db *gorm.DB) *JoinRequestDB {
	return &JoinRequestDB{db}
}

// CREATOR REPO

func (j *JoinRequestDB) Add(jr models.JoinRequest) error {
	return j.db.
		Create(&jr).
		Error
}

// GETTER REPO

func (j *JoinRequestDB) GetAll(contID uint) (jr []models.JoinRequest, err error) {
	err = j.db.
		Model(new(models.JoinRequest)).
		Find(&jr, "requester_id = ?", contID).
		Error

	return
}

// DELETER REPO

func (j *JoinRequestDB) Delete(jr models.JoinRequest) error {
	return j.db.
		Where("requester_id = ?", jr.RequesterID).
		Delete(&jr).
		Error
}

func (j *JoinRequestDB) DeleteAll() error {
	return j.db.
		Where("true").
		Delete(new(models.JoinRequest)).
		Error
}
