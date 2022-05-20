package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// NotificationDB represents a CRUD db repo for notifications
type NotificationDB[T models.Notification] struct {
	db *gorm.DB
}

// NewNotificationDB returns a new NotificationDB instance
func NewNotificationDB(db *gorm.DB) *NotificationDB[models.Notification] {
	return &NotificationDB[models.Notification]{db: db}
}

func (n *NotificationDB[T]) GetDB() *gorm.DB {
	return n.db
}

// CREATOR REPO

func (n *NotificationDB[T]) Add(notification *models.Notification) error {
	return n.db.
		Create(notification).
		Error
}

func (n *NotificationDB[T]) AddMany(notifications []*models.Notification) error {
	return n.db.
		Create(notifications).
		Error
}

// GETTER REPO

func (n *NotificationDB[T]) Exists(nID uint) bool {
	_, err := n.Get(nID)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (n *NotificationDB[T]) Get(nID uint) (notification models.Notification, err error) {
	err = n.db.First(&notification, "id = ?", nID).Error
	return
}

func (n *NotificationDB[T]) GetByConds(conds ...any) (notifications []models.Notification, err error) {
	if (len(conds) < 2) {
        return nil, errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	err = n.db.
		Model(new(models.Notification)).
		Where(conds[0], conds[1:]).
		Find(&notifications).
		Error

	return
}

func (n *NotificationDB[T]) GetAll() ([]models.Notification, error) {
	count, err := n.Count()
	if err != nil {
		return nil, err
	}
	notifications := make([]models.Notification, count)

	res := n.db.Find(&notifications)
	if res.Error != nil {
		return nil, res.Error
	}
	return notifications, nil
}

func (n *NotificationDB[T]) Count() (int64, error) {
	var count int64
	err := n.db.
		Model(new(models.Notification)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (n *NotificationDB[T]) Update(notification *models.Notification, conds ...any) error {
	return n.db.
		Model(new(models.Notification)).
		Where("id = ?", notification.ID).
		Updates(&notification).
		Error
}

func (n *NotificationDB[T]) UpdateAll(notifications []*models.Notification, conds ...any) error {
	return errors.New("not implemented")
}

// DELETER REPO

func (n *NotificationDB[T]) Delete(notification models.Notification, conds ...any) error {
	return n.db.
		Where("id = ?", notification.ID).
		Delete(&notification).
		Error
}

func (n *NotificationDB[T]) DeleteAll(conds ...any) error {
	if (len(conds) < 2) {
        return errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	return n.db.
		Where(conds[0], conds[1:]).
		Delete(new(models.Notification)).
		Error
}

