package helpers

import (
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/timefmt"
	"github.com/robfig/cron"
)

// NotificationHelper well it's written on the box :)
type NotificationHelper struct {
	repo data.NotificationCRUDRepo
}

// NewNotificationHelper returns a new NotificationHelper instance
func NewNotificationHelper(repo data.NotificationCRUDRepo) *NotificationHelper {
	return &NotificationHelper{repo}
}

// GetNotifications returns user's notifications based on the given session
func (n *NotificationHelper) GetNotifications(session models.Session) ([]models.Notification, error) {
	return n.repo.GetAllForUser(session.UserID)
}

// CheckNotifications reports whether a user has notifications or not
func (n *NotificationHelper) CheckNotifications(session models.Session) bool {
	nots, err := n.GetNotifications(session)
	return err != nil && len(nots) > 0
}

// ClearNotifications deletes all notifications for user base on the given session
func (n *NotificationHelper) ClearNotifications(session models.Session) error {
	return n.repo.DeleteAllForUser(session.UserID)
}

func (n *NotificationHelper) SendMany(notifications []*models.Notification) error {
	return n.repo.AddMany(notifications)
}

func (n *NotificationHelper) SendManyOnTime(notifications []*models.Notification, t time.Time) error {
	c := cron.New()
	err := c.AddFunc(timefmt.GetCronTime(t), func() {
		_ = n.SendMany(notifications)
		c.Stop()
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
