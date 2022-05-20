package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

// JoinRequestHelper manages teams join requests
type JoinRequestHelper struct {
	repo             data.CRUDRepo[models.JoinRequest]
	notificationRepo data.CRUDRepo[models.Notification]
	contestRepo      data.GetterRepo[models.Contest]
	teamManager      *TeamHelper
}

// NewJoinRequestHelper returns a new JoinRequestHelper instance
func NewJoinRequestHelper(repo data.CRUDRepo[models.JoinRequest], nRepo data.CRUDRepo[models.Notification],
	contestRepo data.GetterRepo[models.Contest], teamManager *TeamHelper) *JoinRequestHelper {
	return &JoinRequestHelper{
		repo:             repo,
		notificationRepo: nRepo,
		contestRepo:      contestRepo,
		teamManager:      teamManager,
	}
}

// RequestJoinTeam sends a notification to the requested team leader
func (j *JoinRequestHelper) RequestJoinTeam(jr models.JoinRequest, cont models.Contestant) error {
	reqMsg := fmt.Sprintf(
		"_REQThe contestant '%s' with the university id '%s' wants to join your team",
		cont.User.Name, strings.Split(cont.User.Email, "@")[0])

	if jr.RequestMessage != "" {
		reqMsg += "\nRequest message: " + jr.RequestMessage
	}

	var err error
	if jr.RequestedTeamID == 0 {
		jr.RequestedTeam, err = j.teamManager.GetTeamByJoinID(jr.RequestedTeamJoinID)
		if err != nil {
			return err
		}

		contest, _ := j.contestRepo.Get(jr.RequestedContestID)

		if len(jr.RequestedTeam.Members) >= int(contest.ParticipationConditions.MaxTeamMembers) {
			return errors.New("this team is full")
		}
	}

	err = j.checkRequestedTeam(jr, cont)
	if err != nil {
		return err
	}

	reqMsg += fmt.Sprintf("_IDS%d:%d:%d", cont.User.ID, jr.RequestedTeam.ID, jr.RequestedContestID) // a weird way to store ids in the notification text(they won't appear)

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
	jr.RequesterID = cont.User.ID

	err = j.repo.Add(&jr)
	if err != nil { // join request didn't go well :(
		_ = j.notificationRepo.Delete(notification)
		return err
	}

	return nil
}

func (j *JoinRequestHelper) checkRequestedTeam(jr models.JoinRequest, cont models.Contestant) error {
	contest, err := j.contestRepo.Get(jr.RequestedContestID)
	if err != nil {
		return err
	}

	if len(jr.RequestedTeam.Members) >= int(contest.ParticipationConditions.MaxTeamMembers) {
		return errors.New("sorry, this team is full")
	}

	return nil
}

// AcceptJoinRequest adds the requester to the requested team and deletes the other requests & notifications
// and sends a success notification to the requester
func (j *JoinRequestHelper) AcceptJoinRequest(noti models.Notification) error {
	requesterID, teamID, contestID := j.juiceNotification(noti.Content)

	contest, err := j.contestRepo.Get(contestID)
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

// RejectJoinRequest rejects the requester to join the team and deletes the leader's notification
func (j *JoinRequestHelper) RejectJoinRequest(noti models.Notification) error {
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

	jrs, err := j.repo.GetByConds("requester_id = ?", requesterID)
	if err != nil {
		return err
	}

	for _, jr := range jrs {
		if jr.NotificationID == noti.ID {
			err = j.repo.Delete(jr)
			if err != nil {
				return err
			}
			break
		}
	}

	return j.notificationRepo.Delete(noti)
}

func (j *JoinRequestHelper) juiceNotification(notiContent string) (uint, uint, uint) {
	idsStr := notiContent[strings.LastIndex(notiContent, "_IDS")+len("_IDS"):]
	reqID, _ := strconv.Atoi(idsStr[:strings.IndexByte(idsStr, ':')])

	teamContestIDs := idsStr[strings.IndexByte(idsStr, ':')+1:]

	teamID, _ := strconv.Atoi(teamContestIDs[:strings.IndexByte(teamContestIDs, ':')])
	contestID, _ := strconv.Atoi(teamContestIDs[strings.IndexByte(teamContestIDs, ':')+1:])

	return uint(reqID), uint(teamID), uint(contestID)
}

// DeleteRequests deletes all join requests done by a contestant
// used when a contestant is approved to a team, or if a contestant creates a team
func (j *JoinRequestHelper) DeleteRequests(contID, notiID uint) error {
	jrs, err := j.repo.GetByConds("requester_id = ?", contID)
	if err != nil {
		return err
	}

	for _, jr := range jrs {
		_ = j.repo.Delete(jr)
	}
	return nil
	// just in case the upper stuff fails :)
	//return j.repo.Delete(models.JoinRequest{
	//	RequesterID:    contID,
	//	NotificationID: notiID,
	//})
}

// CheckContestantTeamRequests reports whether the given contestant has requested to join the given team
func (j *JoinRequestHelper) CheckContestantTeamRequests(cont models.Contestant, team models.Team) bool {
	jrs, err := j.repo.GetByConds("requester_id = ?", cont.User.ID)
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
