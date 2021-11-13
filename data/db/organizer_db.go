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

func (o *OrganizerDB) Get(organizer models.Organizer) (models.Organizer, error) {
	var (
		fetchedOrganizer models.Organizer
		err              error
	)

	err = o.db.
		First(&fetchedOrganizer, "id = ?", organizer.ID).
		Error

	return fetchedOrganizer, err
}

func (o *OrganizerDB) GetByEmail(email string) (models.Organizer, error) {
	var (
		fetchedOrganizer models.Organizer
		err              error
	)

	err = o.db.
		First(&fetchedOrganizer, "email = ?", email).
		Error

	return fetchedOrganizer, err
}

func (o *OrganizerDB) GetAllByOrganizer(org models.Organizer) ([]models.Organizer, error) {
	var orgs []models.Organizer

	err := o.db.
		Model(new(models.Organizer)).
		Find(&orgs, "director_id = ?", org.ID).
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

	for _, organizer := range organizers {
		o.db.First(&organizer.ContactInfo,
			"contact_info_id = ?", organizer.ContactInfoID)
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
		Where("id = ?", organizer.ContactInfoID).
		Updates(&organizer.ContactInfo).
		Error

	if err != nil {
		return err
	}

	return o.db.
		Model(new(models.Organizer)).
		Where("id = ?", organizer.ID).
		Updates(&organizer).
		Error
}

// DELETER REPO

func (o *OrganizerDB) Delete(organizer models.Organizer) error {
	return o.db.
		Where("id = ?", organizer.ID).
		Delete(&organizer).
		Error
}

func (o *OrganizerDB) DeleteAll() error {
	return o.db.
		Where("true").
		Delete(new(models.Organizer)).
		Error
}
