package managers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type JoinRequestManager struct {
	jrRepo           data.JoinRequestCRDRepo
	notificationRepo data.NotificationCRUDRepo
	teamManager      *TeamManager
}

func NewJoinRequestManager(jrRepo data.JoinRequestCRDRepo, nRepo data.NotificationCRUDRepo, tm *TeamManager) *JoinRequestManager {
	return &JoinRequestManager{
		jrRepo:           jrRepo,
		notificationRepo: nRepo,
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

	reqMsg += fmt.Sprintf("_IDS%d:%d", cont.ID, jr.RequestedTeamID) // a weird way to store ids in the notification text(they won't appear)

	notification := models.Notification{ // send a join request notification to the team leader!
		UserID:  jr.RequestedTeam.LeaderId,
		Content: reqMsg,
	}

	err := j.notificationRepo.Add(&notification)
	if err != nil {
		return err
	}

	err = j.jrRepo.Add(models.JoinRequest{
		RequesterID:     cont.ID,
		RequestedTeamID: jr.RequestedTeamID,
		RequestMessage:  reqMsg[4:],
		NotificationID:  notification.ID,
	})

	if err != nil { // join request didn't go well :(
		_ = j.notificationRepo.Delete(notification)
		return err
	}

	return nil
}

func (j *JoinRequestManager) AcceptJoinRequest(noti models.Notification) error {
	requesterID, teamID := j.getContAndTeamID(noti.Content)

	err := j.teamManager.AddContestantToTeam(requesterID, teamID)
	if err != nil {
		return err
	}

	team, err := j.teamManager.GetTeam(teamID)
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
	requesterID, teamID := j.getContAndTeamID(noti.Content)

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
			break
		}
	}

	return j.notificationRepo.Delete(noti)
}

func (j *JoinRequestManager) getContAndTeamID(notiContent string) (uint, uint) {
	idsStr := notiContent[strings.LastIndex(notiContent, "_IDS")+len("_IDS"):]
	reqID, _ := strconv.Atoi(idsStr[:strings.IndexByte(idsStr, ':')])
	teamID, _ := strconv.Atoi(idsStr[strings.IndexByte(idsStr, ':')+1:])

	return uint(reqID), uint(teamID)
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
