package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// NotificationDB represents a CRUD db repo for notifications
type NotificationDB struct {
	db *gorm.DB
}

// NewNotificationDB returns a new NotificationDB instance
func NewNotificationDB(db *gorm.DB) *NotificationDB {
	return &NotificationDB{db: db}
}

// CREATOR REPO

func (n *NotificationDB) Add(notification *models.Notification) error {
	return n.db.
		Create(notification).
		Error
}

// GETTER REPO

func (n *NotificationDB) Exists(notification models.Notification) (bool, error) {
	res := n.db.First(&notification)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (n *NotificationDB) Get(notification models.Notification) (models.Notification, error) {
	fetchedNotification := models.Notification{}
	err := n.db.First(&fetchedNotification, "id = ?", notification.ID).Error

	return fetchedNotification, err
}

func (n *NotificationDB) GetAll() ([]models.Notification, error) {
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

func (n *NotificationDB) GetAllForUser(user models.User) ([]models.Notification, error) {
	notifications := make([]models.Notification, 0)

	err := n.db.
		Model(new(models.Notification)).
		Where("user_id = ?", user.ID).
		Find(&notifications).
		Error

	return notifications, err
}

func (n *NotificationDB) Count() (int64, error) {
	var count int64
	err := n.db.
		Model(new(models.Notification)).
		Count(&count).
		Error

	return count, err
}

// UPDATER REPO

func (n *NotificationDB) Update(src models.Notification, dst models.Notification) error {
	return n.db.
		Model(new(models.Notification)).
		Where("id = ?", src.ID).
		Updates(&dst).
		Error
}

// DELETER REPO

func (n *NotificationDB) Delete(notification models.Notification) error {
	return n.db.
		Where("id = ?", notification.ID).
		Delete(&notification).
		Error
}

func (n *NotificationDB) DeleteAll() error {
	return n.db.
		Where("true").
		Delete(new(models.Notification)).
		Error
}
