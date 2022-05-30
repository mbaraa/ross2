package data

import "github.com/mbaraa/ross2/models"

type RegisterTeamCreatorRepo interface {
	Add(rt *models.RegisterTeam) error
}

type RegisterTeamGetterRepo interface {
	Exists(rt models.RegisterTeam) (bool, error)
	Get(rt models.RegisterTeam, conds ...any) ([]models.RegisterTeam, error)
	GetTeams(cont models.Contestant) ([]models.RegisterTeam, error)
	GetConts(team models.Team) ([]models.Contestant, error)
	Count() (int64, error)
}

type RegisterTeamUpdaterRepo interface {
	Update(rt *models.RegisterTeam, conds ...any) error
}

type RegisterTeamDeleterRepo interface {
	Delete(rt models.RegisterTeam, conds ...any) error
	DeleteAll() error
}

type RegisterTeamCRUDRepo interface {
	RegisterTeamCreatorRepo
	RegisterTeamGetterRepo
	RegisterTeamUpdaterRepo
	RegisterTeamDeleterRepo
}
