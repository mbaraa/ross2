package db

import (
	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{db: db}
}

// CREATOR REPO

func (u *UserDB) Add(user *models.User) error {
	return u.db.
		Create(user).
		Error
}

// GETTER REPO

func (u *UserDB) Exists(user models.User) (bool, error) {
	_, err := u.Get(user)
	return err == nil, err
}

func (u *UserDB) Get(user models.User) (fetchedUser models.User, err error) {
	err = u.db.
		First(&fetchedUser, "id = ?", user.ID).
		Error

	return
}

func (u *UserDB) GetByEmail(email string) (fetchedUser models.User, err error) {
	err = u.db.
		First(&fetchedUser, "email = ?", email).
		Error

	return
}

func (u *UserDB) GetAll() (users []models.User, err error) {
	err = u.db.
		Find(&users).
		Error
	return
}

func (u *UserDB) Count() (count int64, err error) {
	err = u.db.
		Model(new(models.User)).
		Count(&count).
		Error

	return
}

// UPDATER REPO

func (u *UserDB) Update(user *models.User) (err error) {
	err = u.db.
		Model(new(models.ContactInfo)).
		Where("id = ?", user.ContactInfoID).
		Updates(&user.ContactInfo).
		Error
	if err != nil {
		return
	}

	return u.db.
		Model(new(models.User)).
		Where("id = ?", user.ID).
		Updates(user).
		Error
}

// DELETER REPO

func (u *UserDB) Delete(user models.User) error {
	return u.db.
		Where("id = ?", user.ID).
		Delete(&user).
		Error
}

func (u *UserDB) DeleteAll() error {
	return u.db.
		Where("true").
		Delete(new(models.Organizer)).
		Error
}
