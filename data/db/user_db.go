package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

type UserDB[T models.User] struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB[models.User] {
	return &UserDB[models.User]{db: db}
}

func (u *UserDB[T]) GetDB() *gorm.DB {
	return u.db
}

// CREATOR REPO

func (u *UserDB[T]) Add(user *models.User) error {
	return u.db.
		Create(user).
		Error
}

func (u *UserDB[T]) AddMany(users []*models.User) error {
	return u.db.
		Create(users).
		Error
}

// GETTER REPO

func (u *UserDB[T]) Exists(userID uint) bool {
	_, err := u.Get(userID)
	return err == nil
}

func (u *UserDB[T]) Get(userID uint) (fetchedUser models.User, err error) {
	err = u.db.
		First(&fetchedUser, "id = ?", userID).
		Error

	return
}

func (u *UserDB[T]) GetByConds(conds ...any) (users []models.User, err error) {
	if (len(conds) < 2) {
        return nil, errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	err = u.db.
		Model(new(models.User)).
		First(&users, conds[0], conds[1:]).
		Error

	return
}

func (u *UserDB[T]) GetAll() (users []models.User, err error) {
	err = u.db.
		Find(&users).
		Error
	return
}

func (u *UserDB[T]) Count() (count int64, err error) {
	err = u.db.
		Model(new(models.User)).
		Count(&count).
		Error

	return
}

// UPDATER REPO

func (u *UserDB[T]) Update(user *models.User, conds ...any) (err error) {
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

func (u *UserDB[T]) UpdateAll(users []*models.User, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (u *UserDB[T]) Delete(user models.User, conds ...any) error {
	return u.db.
		Where("id = ?", user.ID).
		Delete(&user).
		Error
}

func (u *UserDB[T]) DeleteAll(conds ...any) error {
	return u.db.
		Where("true").
		Delete(new(models.Organizer)).
		Error
}
