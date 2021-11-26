package managers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type JoinRequestManager struct {
	jrRepo           data.JoinRequestCRDRepo
	notificationRepo data.NotificationCRUDRepo
	contestRepo      data.ContestGetterRepo
	teamManager      *TeamManager
}

func NewJoinRequestManager(jrRepo data.JoinRequestCRDRepo, nRepo data.NotificationCRUDRepo, contestRepo data.ContestGetterRepo,
	tm *TeamManager) *JoinRequestManager {
	return &JoinRequestManager{
		jrRepo:           jrRepo,
		notificationRepo: nRepo,
		contestRepo:      contestRepo,
		teamManager:      tm,
	}
}

func (j *JoinRequestManager) CreateRequest(jr models.JoinRequest, cont models.Contestant) error {
	reqMsg := fmt.Sprintf(
		"_REQThe contestant '%s' with the university id '%s' wants to join your team",
		cont.Name, cont.UniversityID)

	if jr.RequestMessage != "" {
		reqMsg += "\nRequest message: " + jr.RequestMessage
	}

	var err error
	jr.RequestedTeam, err = j.teamManager.GetTeam(jr.RequestedTeamID)
	if err != nil {
		return err
	}

	reqMsg += fmt.Sprintf("_IDS%d:%d:%d", cont.ID, jr.RequestedTeamID, jr.RequestedContestID) // a weird way to store ids in the notification text(they won't appear)

	notification := models.Notification{ // send a join request notification to the team leader!
		UserID:  jr.RequestedTeam.LeaderId,
		Content: reqMsg,
	}

	err = j.notificationRepo.Add(&notification)
	if err != nil {
		return err
	}

	jr.NotificationID = notification.ID
	jr.Notification = notification
	jr.RequestMessage = reqMsg[4:]

	err = j.jrRepo.Add(jr)
	if err != nil { // join request didn't go well :(
		_ = j.notificationRepo.Delete(notification)
		return err
	}

	return nil
}

func (j *JoinRequestManager) AcceptJoinRequest(noti models.Notification) error {
	requesterID, teamID, contestID := j.juiceNotification(noti.Content)

	contest, err := j.contestRepo.Get(models.Contest{ID: contestID})
	if err != nil {
		return err
	}

	team, err := j.teamManager.GetTeam(teamID)
	if err != nil {
		return err
	}

	if uint(len(team.Members)) >= contest.ParticipationConditions.MaxTeamMembers {
		err = j.notificationRepo.Delete(noti)
		if err != nil {
			return err
		}

		err = j.notificationRepo.Add(&models.Notification{
			UserID:  requesterID,
			Content: fmt.Sprintf(`your request to the team "%s" was rejected because the team is full!`, team.Name),
		})
		if err != nil {
			return err
		}

		return errors.New("max allowed team members for this contest is exceeded")
	}

	team, err = j.teamManager.AddContestantToTeam(requesterID, teamID)
	if err != nil {
		return err
	}

	err = j.notificationRepo.Add(&models.Notification{
		UserID:  requesterID,
		Content: "Your request to join the team '" + team.Name + "' has been approved",
	})
	if err != nil {
		return err
	}

	return j.DeleteRequests(requesterID, noti.ID) // this will also delete the team's leader notification
}

func (j *JoinRequestManager) RejectJoinRequest(noti models.Notification) error {
	requesterID, teamID, _ := j.juiceNotification(noti.Content)

	team, err := j.teamManager.GetTeam(teamID)
	if err != nil {
		return err
	}

	err = j.notificationRepo.Add(&models.Notification{
		UserID:  requesterID,
		Content: "Your request to join the team '" + team.Name + "' has been denied 😕\nTry an another team or register as teamless!",
	})
	if err != nil {
		return err
	}

	jrs, err := j.jrRepo.GetAll(requesterID)
	if err != nil {
		return err
	}

	for _, jr := range jrs {
		if jr.NotificationID == noti.ID {
			err = j.jrRepo.Delete(jr)
			if err != nil {
				return err
			}
			break
		}
	}

	return j.notificationRepo.Delete(noti)
}

func (j *JoinRequestManager) juiceNotification(notiContent string) (uint, uint, uint) {
	idsStr := notiContent[strings.LastIndex(notiContent, "_IDS")+len("_IDS"):]
	reqID, _ := strconv.Atoi(idsStr[:strings.IndexByte(idsStr, ':')])

	teamContestIDs := idsStr[strings.IndexByte(idsStr, ':')+1:]

	teamID, _ := strconv.Atoi(teamContestIDs[:strings.IndexByte(teamContestIDs, ':')])
	contestID, _ := strconv.Atoi(teamContestIDs[strings.IndexByte(teamContestIDs, ':')+1:])

	return uint(reqID), uint(teamID), uint(contestID)
}

func (j *JoinRequestManager) DeleteRequests(contID, notiID uint) error {
	return j.jrRepo.Delete(models.JoinRequest{
		RequesterID:    contID,
		NotificationID: notiID,
	})
}

func (j *JoinRequestManager) CheckContestantTeamRequests(cont models.Contestant, team models.Team) bool {
	jrs, err := j.jrRepo.GetAll(cont.ID)
	if err != nil {
		return false
	}

	for _, jr := range jrs {
		if jr.RequestedTeamID == team.ID {
			return true
		}
	}

	return false
}
