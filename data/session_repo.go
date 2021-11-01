package data

import "github.com/mbaraa/ross2/models"

// SessionCreatorRepo is an interface that allows creation of a session
// into a certain data source :)
type SessionCreatorRepo interface {
	Add(session *models.Session) error
}

// SessionGetterRepo is an interface that allows getting values of a session
// from a certain data source :)
type SessionGetterRepo interface {
	Exists(session models.Session) (bool, error)
	Get(session models.Session) (models.Session, error)
}

// SessionDeleterRepo is an interface that allows deleting values of a session
// from a certain data source :)
type SessionDeleterRepo interface {
	Delete(session models.Session) error
	DeleteAllForUser(userID uint) error
	DeleteAll() error
}

// SessionCRUDRepo is an interface that allows full CRUD operations of a session
// on a certain data source :)
type SessionCRUDRepo interface {
	SessionCreatorRepo
	SessionGetterRepo
	SessionDeleterRepo
}
