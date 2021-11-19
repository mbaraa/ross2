package sheevhelper

import (
	"errors"
	"fmt"

	"github.com/mbaraa/ross2/models"
)

func GetSheevNotifications(c models.Contest) (notifications []*models.Notification, err error) {
	conts, err := getContestants(c)
	if err != nil {
		return nil, err
	}
	for _, cont := range conts {
		notifications = append(notifications, &models.Notification{
			UserID: cont.ID,
			Content: fmt.Sprintf(
				`Hi %s 👋<br/>
The contest "%s" is over, hope you had fun at it 😉<br/>
if you participated in the contest you can generate socity service form from <a href="%s">this site</a> then print the picture and deliver it to the dean's secretary<br/>
Have a nice day!`, cont.Name, c.Name, "https://mbaraa.fun/sheev"),
		})
	}

	return
}

func getContestants(c models.Contest) (contsIDs []models.Contestant, err error) {
	if len(c.Teams) == 0 {
		return nil, errors.New("no teams were found")
	}

	for _, team := range c.Teams {
		if len(team.Members) == 0 {
			continue
		}
		for _, member := range team.Members {
			contsIDs = append(contsIDs, member)
		}
	}
	return
}
