package data

import (
	"github.com/mbaraa/ross2/models"
)

type UserCreatorRepo interface {
	Add(user *models.User) error
}

type UserGetterRepo interface {
	Exists(user models.User) (bool, error)
	Get(user models.User) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetAll() ([]models.User, error)
	Count() (int64, error)
}

type UserUpdaterRepo interface {
	Update(user *models.User) error
}

type UserDeleterRepo interface {
	Delete(user models.User) error
	DeleteAll() error
}

type UserCRUDRepo interface {
	UserCreatorRepo
	UserGetterRepo
	UserUpdaterRepo
	UserDeleterRepo
}
