package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

type RegisterTeamDB struct {
	db *gorm.DB
}

func NewRegisterTeamDB(db *gorm.DB) *RegisterTeamDB {
	return &RegisterTeamDB{db}
}

// CREATOR REPO

func (r *RegisterTeamDB) Add(rt *models.RegisterTeam) error {
	return r.db.
		Create(rt).
		Error
}

// GETTER REPO

func (r *RegisterTeamDB) Exists(rt models.RegisterTeam) (bool, error) {
	return false, errors.New("not implemented")
}

func (r *RegisterTeamDB) Get(oc models.RegisterTeam, conds ...any) (fetchedRT []models.RegisterTeam, err error) {
	err = r.db.
		Find(&fetchedRT, conds...).
		Error

	return
}

func (r *RegisterTeamDB) GetTeams(cont models.Contestant) (rTeams []models.RegisterTeam, err error) {
	err = r.db.
		Model(new(models.RegisterTeam)).
		Find(&rTeams, "contestant_id = ?", cont.ID).
		Error

	return
}

func (r *RegisterTeamDB) GetConts(team models.Team) (conts []models.Contestant, err error) {
	var rt models.RegisterTeam
	err = r.db.
		Model(new(models.RegisterTeam)).
		First(&rt, "team_id = ?", team.ID).
		Error

	if err != nil {
		return nil, err
	}

	conts = rt.Team.Members[:]

	return
}

func (r *RegisterTeamDB) Count() (int64, error) {
	return 0, errors.New("not implemented")
}

// The Updater & Deleter Repos' Methods doesn't modify data of the teams table :)

// UPDATER REPO

func (r *RegisterTeamDB) Update(rt *models.RegisterTeam, conds ...any) error {
	return r.db.
		Model(rt).
		Where(conds[0], conds...).
		Updates(rt).
		Error
}

// DELETER REPO

func (r *RegisterTeamDB) Delete(rt models.RegisterTeam, conds ...any) error {
	return r.db.
		Model(&rt).
		Where(conds[0], conds...).
		Delete(&rt).
		Error
}

func (r *RegisterTeamDB) DeleteAll() error {
	return r.db.
		Where("true").
		Delete(new(models.RegisterTeam)).
		Error
}
