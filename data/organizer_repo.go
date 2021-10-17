package data

import "github.com/mbaraa/ross2/models"

// OrganizerCreatorRepo is an interface that allows creation of an organizer
// into a certain data source :)
type OrganizerCreatorRepo interface {
	Add(organizer models.Organizer) error
}

// OrganizerGetterRepo is an interface that allows getting values of an organizer
// from a certain data source :)
type OrganizerGetterRepo interface {
	Exists(organizer models.Organizer) (bool, error)
	Get(organizer models.Organizer) (models.Organizer, error)
	GetByEmail(email string) (models.Organizer, error)
	GetAll() ([]models.Organizer, error)
	Count() (int64, error)
}

// OrganizerUpdaterRepo is an interface that allows updating values of an organizer
// in a certain data source :)
type OrganizerUpdaterRepo interface {
	Update(organizer models.Organizer) error
}

// OrganizerDeleterRepo is an interface that allows deleting values of an organizer
// from a certain data source :)
type OrganizerDeleterRepo interface {
	Delete(organizer models.Organizer) error
	DeleteAll() error
}

// OrganizerCRUDRepo is an interface that allows full CRUD operations of an organizer
// on a certain data source :)
type OrganizerCRUDRepo interface {
	OrganizerCreatorRepo
	OrganizerGetterRepo
	OrganizerUpdaterRepo
	OrganizerDeleterRepo
}
