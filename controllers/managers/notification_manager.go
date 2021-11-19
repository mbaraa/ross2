package managers

import (
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/timefmt"
	"github.com/robfig/cron"
)

type NotificationManager struct {
	repo data.NotificationCRUDRepo
}

func NewNotificationManager(repo data.NotificationCRUDRepo) *NotificationManager {
	return &NotificationManager{repo}
}

func (n *NotificationManager) SendMany(notifications []*models.Notification) error {
	return n.repo.AddMany(notifications)
}

func (n *NotificationManager) SendManyOnTime(notifications []*models.Notification, t time.Time) error {
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
