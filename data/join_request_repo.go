package data

import "github.com/mbaraa/ross2/models"

// JoinRequestCreatorRepo is an interface that allows creation of a rj
// into a certain data source :)
type JoinRequestCreatorRepo interface {
	Add(jr models.JoinRequest) error
}

type JoinRequestGetterRepo interface {
	GetAll(contID uint) ([]models.JoinRequest, error)
}

// JoinRequestDeleterRepo is an interface that allows deleting values of a rj
// from a certain data source :)
type JoinRequestDeleterRepo interface {
	Delete(jr models.JoinRequest) error
	DeleteAll() error
}

// JoinRequestCRDRepo is an interface that allows full CRUD operations of a rj
// on a certain data source :)
type JoinRequestCRDRepo interface {
	JoinRequestCreatorRepo
	JoinRequestGetterRepo
	JoinRequestDeleterRepo
}
