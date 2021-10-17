package data

import "github.com/mbaraa/ross2/models"

// ContestCreatorRepo is an interface that allows creation of a contest
// into a certain data source :)
type ContestCreatorRepo interface {
	Add(contest models.Contest) error
}

// ContestGetterRepo is an interface that allows getting values of a contest
// from a certain data source :)
type ContestGetterRepo interface {
	Exists(contest models.Contest) (bool, error)
	Get(contest models.Contest) (models.Contest, error)
	GetAll() ([]models.Contest, error)
	GetAllByOrganizer(org models.Organizer) ([]models.Contest, error)
	Count() (int64, error)
}

// ContestUpdaterRepo is an interface that allows updating values of a contest
// in a certain data source :)
type ContestUpdaterRepo interface {
	Update(contest models.Contest) error
}

// ContestDeleterRepo is an interface that allows deleting values of a contest
// from a certain data source :)
type ContestDeleterRepo interface {
	Delete(contest models.Contest) error
	DeleteAll() error
}

// ContestCRUDRepo is an interface that allows full CRUD operations of a contest
// on a certain data source :)
type ContestCRUDRepo interface {
	ContestCreatorRepo
	ContestGetterRepo
	ContestUpdaterRepo
	ContestDeleterRepo
}
