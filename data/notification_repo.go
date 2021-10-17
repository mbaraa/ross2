package data

import "github.com/mbaraa/ross2/models"

// NotificationCreatorRepo is an interface that allows creation of a notification
// into a certain data source :)
type NotificationCreatorRepo interface {
	// Add has pointer to Notification so that the generated id can be used
	Add(notification *models.Notification) error
}

// NotificationGetterRepo is an interface that allows getting values of a notification
// from a certain data source :)
type NotificationGetterRepo interface {
	Exists(notification models.Notification) (bool, error)
	Get(notification models.Notification) (models.Notification, error)
	GetAll() ([]models.Notification, error)
	GetAllForUser(user models.User) ([]models.Notification, error)
	Count() (int64, error)
}

// NotificationUpdaterRepo is an interface that allows updating values of a notification
// in a certain data source :)
type NotificationUpdaterRepo interface {
	Update(src models.Notification, dst models.Notification) error
}

// NotificationDeleterRepo is an interface that allows deleting values of a notification
// from a certain data source :)
type NotificationDeleterRepo interface {
	Delete(notification models.Notification) error
	DeleteAll() error
}

// NotificationCRUDRepo is an interface that allows full CRUD operations of a notification
// on a certain data source :)
type NotificationCRUDRepo interface {
	NotificationCreatorRepo
	NotificationGetterRepo
	NotificationUpdaterRepo
	NotificationDeleterRepo
}
