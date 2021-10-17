package data

import "github.com/mbaraa/ross2/models"

// JoinRequestCreatorRepo is an interface that allows creation of a rj
// into a certain data source :)
type JoinRequestCreatorRepo interface {
	Add(jr models.JoinRequest) error
}

// JoinRequestDeleterRepo is an interface that allows deleting values of a rj
// from a certain data source :)
type JoinRequestDeleterRepo interface {
	Delete(rj models.JoinRequest) error
	DeleteAll() error
}

// JoinRequestCDRepo is an interface that allows full CRUD operations of a rj
// on a certain data source :)
type JoinRequestCDRepo interface {
	JoinRequestCreatorRepo
	JoinRequestDeleterRepo
}
