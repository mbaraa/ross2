package db

import (
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

type AdminDB struct {
	db *gorm.DB
}

func NewAdminDB(db *gorm.DB) *AdminDB {
	return &AdminDB{db}
}

// CREATOR REPO

func (a *AdminDB) Add(admin *models.Admin) error {
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

func (a *AdminDB) Get(admin models.Admin) (fetchedAdmin models.Admin, err error) {
	err = a.db.
		Model(new(models.Admin)).
		First(&fetchedAdmin, "user_id = ?", admin.User.ID).
		Error

	err = a.db.
		Model(new(models.User)).
		First(&fetchedAdmin.User, "id = ?", admin.User.ID).
		Error

	return
}

// UPDATER REPO

func (a *AdminDB) Update(admin *models.Admin) error {
	return a.db.
		Model(new(models.Admin)).
		Where("user_id = ?", admin.User.ID).
		Updates(&admin).
		Error
}

// DELETER REPO

func (a *AdminDB) Delete(admin models.Admin) error {
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

func (a *AdminDB) getUser(admin models.Admin) (models.User, error) {
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

func (a *AdminDB) updateUser(user *models.User) error {
	return a.db.
		Model(new(models.User)).
		Where("id = ?", user.ID).
		Updates(user).
		Error
}
