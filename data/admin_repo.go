package data

import "github.com/mbaraa/ross2/models"

// AdminCreatorRepo is an interface that allows creation of an admin
// into a certain data source :)
type AdminCreatorRepo interface {
	Add(admin *models.Admin) error
}

// AdminGetterRepo is an interface that allows getting values of an admin
// from a certain data source :)
type AdminGetterRepo interface {
	Get(admin models.Admin) (models.Admin, error)
}

// AdminUpdaterRepo is an interface that allows updating values of an admin
// in a certain data source :)
type AdminUpdaterRepo interface {
	Update(admin *models.Admin) error
}

// AdminDeleterRepo is an interface that allows deleting values of an admin
// from a certain data source :)
type AdminDeleterRepo interface {
	Delete(admin models.Admin) error
}

// AdminCRUDRepo is an interface that allows full CRUD operations of an admin
// on a certain data source :)
type AdminCRUDRepo interface {
	AdminCreatorRepo
	AdminGetterRepo
	AdminUpdaterRepo
	AdminDeleterRepo
}
