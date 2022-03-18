package db

import (
	"errors"
	"strings"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// OrganizerDB represents a CRUD db repo for organizers
type OrganizerDB[T models.Organizer] struct {
	db *gorm.DB
}

// NewOrganizerDB returns a new OrganizerDB instance
func NewOrganizerDB(db *gorm.DB) *OrganizerDB[models.Organizer] {
	return &OrganizerDB[models.Organizer]{db: db}
}

func (o *OrganizerDB[T]) GetDB() *gorm.DB {
	return o.db
}

// CREATOR REPO

func (o *OrganizerDB[T]) Add(organizer *models.Organizer) error {
	return o.db.
		Create(organizer).
		Error
}

func (o *OrganizerDB[T]) AddMany(orgs []*models.Organizer) error {
	return errors.New("not implemented")
}

// GETTER REPO

func (o *OrganizerDB[T]) Exists(userID uint) bool {
	_, err := o.Get(userID)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (o *OrganizerDB[T]) Get(userID uint) (fetchedOrganizer models.Organizer, err error) {
	err = o.db.
		Model(new(models.Organizer)).
		First(&fetchedOrganizer, "user_id = ?", userID).
		Error

	return
}

func (o *OrganizerDB[T]) GetByConds(conds ...any) ([]models.Organizer, error) {
	if (len(conds) < 2) {
		return nil, errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	var orgs []models.Organizer

	if strings.Contains(conds[0].(string), "email") {
		user := models.User{}

		err := o.db.
			Model(new(models.User)).
			First(&user, "email = ?", conds[1]).
			Error

		if err != nil {
			return nil, err
		}

		orgs = make([]models.Organizer, 1)

		err = o.db.
			Model(new(models.Organizer)).
			Find(&orgs, "user_id = ?", user.ID).
			Error

		return orgs, err
	}

	err := o.db.
		Model(new(models.Organizer)).
		Find(&orgs, conds[0], conds[1:]).
		Error

	return orgs, err
}

func (o *OrganizerDB[T]) GetAll() ([]models.Organizer, error) {
	organizers := make([]models.Organizer, 0)

	res := o.db.Find(&organizers)
	if res.Error != nil {
		return nil, res.Error
	}

	return organizers, nil
}

func (o *OrganizerDB[T]) Count() (int64, error) {
	var count int64
	err := o.db.
		Model(new(models.Organizer)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (o *OrganizerDB[T]) Update(organizer *models.Organizer, conds ...any) error {
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

func (o *OrganizerDB[T]) UpdateAll(orgs []*models.Organizer, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (o *OrganizerDB[T]) Delete(organizer models.Organizer, conds ...any) error {
	return o.db.
		Where("user_id = ?", organizer.User.ID).
		Delete(&organizer).
		Error
}

func (o *OrganizerDB[T]) DeleteAll(conds ...any) error {
	return o.db.
		Where("true").
		Delete(new(models.Organizer)).
		Error
}
