package data

import "github.com/mbaraa/ross2/models"

// ContestantCreatorRepo is an interface that allows creation of a contestant
// into a certain data source :)
type ContestantCreatorRepo interface {
	Add(contestant models.Contestant) error
}

// ContestantGetterRepo is an interface that allows getting values of a contestant
// from a certain data source :)
type ContestantGetterRepo interface {
	Exists(contestant models.Contestant) (bool, error)
	Get(contestant models.Contestant) (models.Contestant, error)
	GetByEmail(email string) (models.Contestant, error)
	GetAll() ([]models.Contestant, error)
	Count() (int64, error)
}

// ContestantUpdaterRepo is an interface that allows updating values of a contestant
// in a certain data source :)
type ContestantUpdaterRepo interface {
	Update(cont models.Contestant) error
}

// ContestantDeleterRepo is an interface that allows deleting values of a contestant
// from a certain data source :)
type ContestantDeleterRepo interface {
	Delete(contestant models.Contestant) error
	DeleteAll() error
}

// TeamlessCRUDRepo is an interface that allows crud operations on the teamless contestants :)
type TeamlessCRUDRepo interface {
	AddTeamLess(contestant models.Contestant, contest models.Contest) error
	GetAllTeamLess(contest models.Contest) ([]models.Contestant, error)
	RegisterInTeam(contestant models.Contestant, team models.Team) error
}

// ContestantCRUDRepo is an interface that allows full CRUD operations of a contestant
// on a certain data source :)
type ContestantCRUDRepo interface {
	ContestantCreatorRepo
	ContestantGetterRepo
	ContestantUpdaterRepo
	ContestantDeleterRepo

	TeamlessCRUDRepo
}
