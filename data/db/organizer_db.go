package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// OrganizerDB represents a CRUD db repo for organizers
type OrganizerDB struct {
	db *gorm.DB
}

// NewOrganizerDB returns a new OrganizerDB instance
func NewOrganizerDB(db *gorm.DB) *OrganizerDB {
	return &OrganizerDB{db: db}
}

func (o *OrganizerDB) GetDB() *gorm.DB {
	return o.db
}

// CREATOR REPO

func (o *OrganizerDB) Add(organizer *models.Organizer) error {
	return o.db.
		Create(organizer).
		Error
}

// GETTER REPO

func (o *OrganizerDB) Exists(organizer models.Organizer) (bool, error) {
	res := o.db.First(&organizer)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (o *OrganizerDB) Get(organizer models.Organizer) (fetchedOrganizer models.Organizer, err error) {
	err = o.db.
		Model(new(models.Organizer)).
		First(&fetchedOrganizer, "user_id = ?", organizer.User.ID).
		Error

	return
}

func (o *OrganizerDB) GetByEmail(email string) (fetchedOrganizer models.Organizer, err error) {
	var user models.User
	err = o.db.
		Model(new(models.User)).
		First(&user, "email = ?", email).
		Error

	if err != nil {
		return models.Organizer{}, err
	}

	err = o.db.
		Model(new(models.Organizer)).
		First(&fetchedOrganizer, "user_id = ?", user.ID).
		Error

	return
}

func (o *OrganizerDB) GetAllByOrganizer(org models.Organizer) ([]models.Organizer, error) {
	var orgs []models.Organizer

	err := o.db.
		Model(new(models.Organizer)).
		Find(&orgs, "director_id = ?", org.User.ID).
		Error

	return orgs, err
}

func (o *OrganizerDB) GetAll() ([]models.Organizer, error) {
	count, err := o.Count()
	if err != nil {
		return nil, err
	}
	organizers := make([]models.Organizer, count)

	res := o.db.Find(&organizers)
	if res.Error != nil {
		return nil, res.Error
	}

	return organizers, nil
}

func (o *OrganizerDB) Count() (int64, error) {
	var count int64
	err := o.db.
		Model(new(models.Organizer)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (o *OrganizerDB) Update(organizer models.Organizer) error {
	err := o.db.
		Model(new(models.ContactInfo)).
		Where("id = ?", organizer.User.ContactInfoID).
		Updates(&organizer.User.ContactInfo).
		Error

	if err != nil {
		return err
	}

	return o.db.
		Model(new(models.Organizer)).
		Where("user_id = ?", organizer.User.ID).
		Updates(&organizer).
		Error
}

// DELETER REPO

func (o *OrganizerDB) Delete(organizer models.Organizer) error {
	return o.db.
		Where("user_id = ?", organizer.User.ID).
		Delete(&organizer).
		Error
}

func (o *OrganizerDB) DeleteAll() error {
	return o.db.
		Where("true").
		Delete(new(models.Organizer)).
		Error
}
