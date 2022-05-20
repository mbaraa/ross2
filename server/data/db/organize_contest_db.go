package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

type OrganizeContestDB struct {
	db *gorm.DB
}

func NewOrganizeOrganizeContestDB(db *gorm.DB) *OrganizeContestDB {
	return &OrganizeContestDB{db}
}

// CREATOR REPO

func (c *OrganizeContestDB) Add(oc *models.OrganizeContest) error {
	return c.db.
		Create(oc).
		Error
}

// GETTER REPO

func (c *OrganizeContestDB) Exists(oc models.OrganizeContest) (bool, error) {
	res := c.db.
		First(&oc, "contest_id = ? and organizer_id = ?", oc.ContestID, oc.OrganizerID)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (c *OrganizeContestDB) Get(oc models.OrganizeContest) (fetchedOC models.OrganizeContest, err error) {
	err = c.db.
		First(&fetchedOC, "contest_id = ? and organizer_id = ?", oc.ContestID, oc.OrganizerID).
		Error

	return
}

func (c *OrganizeContestDB) GetOrgs(contest models.Contest) (orgs []models.Organizer, err error) {
	var ocs []models.OrganizeContest

	err = c.db.
		Model(new(models.OrganizeContest)).
		Find(&ocs, "contest_id = ?", contest.ID).
		Error

	if err != nil {
		return nil, err
	}

	for _, oc := range ocs {
		var org models.Organizer
		err = c.db.
			Model(new(models.Organizer)).
			Find(&org, "id = ?", oc.OrganizerID).
			Error

		if err != nil {
			return nil, err
		}

		orgs = append(orgs, org)
	}

	return
}

func (c *OrganizeContestDB) GetContests(org models.Organizer) (contests []models.Contest, err error) {
	var ocs []models.OrganizeContest

	err = c.db.
		Model(new(models.OrganizeContest)).
		Find(&ocs, "contest_id = ?", org.ID).
		Error

	if err != nil {
		return nil, err
	}

	for _, oc := range ocs {
		var contest models.Contest
		err = c.db.
			Model(new(models.Contest)).
			Find(&contest, "id = ?", oc.ContestID).
			Error

		if err != nil {
			return nil, err
		}

		contests = append(contests, contest)
	}

	return
}

func (c *OrganizeContestDB) Count() (int64, error) {
	var count int64
	err := c.db.
		Model(new(models.OrganizeContest)).
		Count(&count).
		Error

	return count, err
}

// The Updater & Deleter Repos' Methods doesn't modify data of the teams table :)

// UPDATER REPO

func (c *OrganizeContestDB) Update(oc *models.OrganizeContest) error {
	return c.db.
		Model(new(models.OrganizeContest)).
		Where("contest_id = ? and organizer_id = ?", oc.ContestID, oc.OrganizerID).
		Updates(oc).
		Error
}

// DELETER REPO

func (c *OrganizeContestDB) Delete(oc models.OrganizeContest) error {
	return c.db.
		Model(new(models.OrganizeContest)).
		Where("contest_id = ? and organizer_id = ?", oc.ContestID, oc.OrganizerID).
		Delete(&oc).
		Error
}

func (c *OrganizeContestDB) DeleteAll() error {
	return c.db.
		Where("true").
		Delete(new(models.OrganizeContest)).
		Error
}
