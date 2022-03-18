package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

type AdminDB[T models.Admin] struct {
	db *gorm.DB
}

func NewAdminDB[T models.Admin](db *gorm.DB) *AdminDB[T] {
	return &AdminDB[T]{db}
}

func (a *AdminDB[T]) GetDB() *gorm.DB {
	return a.db
}

// CREATOR REPO

func (a *AdminDB[T]) Add(admin *models.Admin) error {
	var err error
	admin.User, err = a.getUser(*admin)
	if err != nil {
		return err
	}

	admin.User.UserType |= enums.UserTypeAdmin
	admin.UserID = admin.User.ID

	err = a.updateUser(&admin.User)
	if err != nil {
		return err
	}

	return a.db.
		Create(admin).
		Error
}

// GETTER REPO

func (a *AdminDB[T]) Get(id uint) (fetchedAdmin models.Admin, err error) {
	err = a.db.
		Model(new(models.Admin)).
		First(&fetchedAdmin, "user_id = ?", id).
		Error

	err = a.db.
		Model(new(models.User)).
		First(&fetchedAdmin.User, "id = ?", id).
		Error

	return
}

func (a *AdminDB[T]) Exists(id uint) bool {
	return false
}

func (a *AdminDB[T]) GetByConds(conds ...any) ([]models.Admin, error) {
	return nil, errors.New("not implemented")
}

func (a *AdminDB[T]) GetAll() ([]models.Admin, error) {
	return nil, errors.New("not implemented")
}

func (a *AdminDB[T]) Count() (int64, error) {
	return 0, errors.New("not implemented")
}

// UPDATER REPO

func (a *AdminDB[T]) Update(admin *models.Admin, conds ...any) error {
	return a.db.
		Model(new(models.Admin)).
		Where("user_id = ?", admin.User.ID).
		Updates(&admin).
		Error
}

func (a *AdminDB[T]) UpdateAll(admins []*models.Admin, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (a *AdminDB[T]) Delete(admin models.Admin, conds ...any) error {
	var err error
	admin.User, err = a.getUser(admin)
	if err != nil {
		return err
	}

	if (admin.User.UserType & enums.UserTypeAdmin) != 0 {
		admin.User.UserType -= enums.UserTypeAdmin
	}

	err = a.updateUser(&admin.User)
	if err != nil {
		return err
	}

	return a.db.
		Model(new(models.Admin)).
		Where("email = ?", admin.User.Email).
		Delete(&admin).
		Error
}

func (a *AdminDB[T]) DeleteAll(conds ...any) error {
	return errors.New("not implemented")
}

//////////////
// helpers
/////////////

func (a *AdminDB[T]) getUser(admin models.Admin) (models.User, error) {
	// find user with email then add it as admin
	var user models.User
	err := a.db.
		Model(new(models.User)).
		First(&user, "email = ?", admin.User.Email).
		Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a *AdminDB[T]) updateUser(user *models.User) error {
	return a.db.
		Model(new(models.User)).
		Where("id = ?", user.ID).
		Updates(user).
		Error
}
