package data

import "github.com/mbaraa/ross2/models"

// TeamCreatorRepo is an interface that allows creation of a team
// into a certain data source :)
type TeamCreatorRepo interface {
	Add(team *models.Team) error
	AddMany(teams []*models.Team) error
}

// TeamGetterRepo is an interface that allows getting values of a team
// from a certain data source :)
type TeamGetterRepo interface {
	Exists(team models.Team) (bool, error)
	Get(team models.Team) (models.Team, error)
	GetAll() ([]models.Team, error)
	GetAllByContest(contest models.Contest) ([]models.Team, error)
	Count() (int64, error)
}

// TeamUpdaterRepo is an interface that allows updating values of a team
// in a certain data source :)
type TeamUpdaterRepo interface {
	Update(team models.Team) error
}

// TeamDeleterRepo is an interface that allows deleting values of a team
// from a certain data source :)
type TeamDeleterRepo interface {
	Delete(team models.Team) error
	DeleteAll() error
}

// TeamCRUDRepo is an interface that allows full CRUD operations of a team
// on a certain data source :)
type TeamCRUDRepo interface {
	TeamCreatorRepo
	TeamGetterRepo
	TeamUpdaterRepo
	TeamDeleterRepo
}
