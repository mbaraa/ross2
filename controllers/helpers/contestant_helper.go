package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/multiavatar"
)

type ContestantHelperBuilder struct {
	userRepo         data.UserCRUDRepo
	contestantRepo   data.ContestantCRUDRepo
	contestRepo      data.ContestUpdaterRepo
	notificationRepo data.NotificationCRUDRepo
	teamMgr          *TeamHelper
	jrMgr            *JoinRequestHelper
}

func NewContestantHelperBuilder() *ContestantHelperBuilder {
	return new(ContestantHelperBuilder)
}

func (b *ContestantHelperBuilder) UserRepo(u data.UserCRUDRepo) *ContestantHelperBuilder {
	b.userRepo = u
	return b
}

func (b *ContestantHelperBuilder) ContestantRepo(c data.ContestantCRUDRepo) *ContestantHelperBuilder {
	b.contestantRepo = c
	return b
}

func (b *ContestantHelperBuilder) ContestRepo(c data.ContestUpdaterRepo) *ContestantHelperBuilder {
	b.contestRepo = c
	return b
}

func (b *ContestantHelperBuilder) NotificationRepo(n data.NotificationCRUDRepo) *ContestantHelperBuilder {
	b.notificationRepo = n
	return b
}

func (b *ContestantHelperBuilder) TeamMgr(t *TeamHelper) *ContestantHelperBuilder {
	b.teamMgr = t
	return b
}

func (b *ContestantHelperBuilder) JoinRequestMgr(j *JoinRequestHelper) *ContestantHelperBuilder {
	b.jrMgr = j
	return b
}

func (b *ContestantHelperBuilder) verify() bool {
	sb := new(strings.Builder)

	if b.userRepo == nil {
		sb.WriteString("Contestant Helper Builder: missing user repo!")
	}
	if b.contestantRepo == nil {
		sb.WriteString("Contestant Helper Builder: missing contestant repo!")
	}
	if b.contestRepo == nil {
		sb.WriteString("Contestant Helper Builder: missing contest repo!")
	}
	if b.notificationRepo == nil {
		sb.WriteString("Contestant Helper Builder: missing notification repo!")
	}
	if b.teamMgr == nil {
		sb.WriteString("Contestant Helper Builder: missing team repo!")
	}
	if b.jrMgr == nil {
		sb.WriteString("Contestant Helper Builder: missing join request repo!")
	}

	if sb.Len() != 0 {
		fmt.Println(sb.String())
		return false
	}

	return true
}

func (b *ContestantHelperBuilder) GetContestantManager() *ContestantHelper {
	return NewContestantHelper(b)
}

// ContestantHelper holds contestants underlying operations
type ContestantHelper struct {
	repo             data.ContestantCRUDRepo
	userRepo         data.UserUpdaterRepo
	contestRepo      data.ContestUpdaterRepo
	teamMgr          *TeamHelper
	jrMgr            *JoinRequestHelper
	notificationRepo data.NotificationCRUDRepo
}

// NewContestantHelper returns a new ContestantHelper instance
func NewContestantHelper(b *ContestantHelperBuilder) *ContestantHelper {
	return &ContestantHelper{
		repo:             b.contestantRepo,
		userRepo:         b.userRepo,
		contestRepo:      b.contestRepo,
		teamMgr:          b.teamMgr,
		jrMgr:            b.jrMgr,
		notificationRepo: b.notificationRepo,
	}
}

// Register creates a new contestant
func (c *ContestantHelper) Register(cont models.Contestant) error {
	if cont.User.UserType&enums.UserTypeContestant == 0 {
		cont.User.UserType |= enums.UserTypeContestant
		if cont.User.UserType&1 != 0 {
			cont.User.UserType--
		}

		err := c.userRepo.Update(&cont.User)
		if err != nil {
			return err
		}

		cont.User.AvatarURL = multiavatar.GetAvatarURL()
		cont.TeamID = 1
		cont.UserID = cont.User.ID
		cont.User.ProfileStatus |= enums.ProfileStatusContestantFinished

		return c.repo.Add(&cont)
	}
	return errors.New("user is already a contestant")
}

// GetProfile returns contestant's for the given user
func (c *ContestantHelper) GetProfile(user models.User) (models.Contestant, error) {
	return c.repo.Get(models.Contestant{User: user})
}

// CreateTeam creates a team and adds the given contestant to it as its leader
func (c *ContestantHelper) CreateTeam(contestant models.Contestant, team models.Team) error {
	if contestant.TeamlessContestID != 0 {
		contestant.TeamlessContestID = 0
		contestant.TeamlessedAt = contestant.CreatedAt
	}

	err := c.teamMgr.CreateTeam(contestant, &team)
	if err != nil {
		return err
	}

	return c.jrMgr.DeleteRequests(contestant.User.ID, 0)
}

// DeleteTeam kicks every member out of it and deletes it
func (c *ContestantHelper) DeleteTeam(team models.Team) error {
	return c.teamMgr.DeleteTeam(team)
}

// RequestJoinTeam sends a notification to the requested team leader
func (c *ContestantHelper) RequestJoinTeam(jr models.JoinRequest, cont models.Contestant) error {
	return c.jrMgr.RequestJoinTeam(jr, cont)
}

// AcceptJoinRequest adds the requester to the requested team and deletes the other requests & notifications
// and sends a success notification to the requester
func (c *ContestantHelper) AcceptJoinRequest(notification models.Notification) error {
	return c.jrMgr.AcceptJoinRequest(notification)
}

// RejectJoinRequest rejects the requester to join the team and deletes the leader's notification
func (c *ContestantHelper) RejectJoinRequest(noti models.Notification) error {
	return c.jrMgr.RejectJoinRequest(noti)
}

// LeaveTeam removes the given contestant from their team in a super safe way
func (c *ContestantHelper) LeaveTeam(contestant models.Contestant) error {
	return c.teamMgr.LeaveTeam(contestant)
}

// RegisterAsTeamless adds the given contestant as teamless for the given contest
func (c *ContestantHelper) RegisterAsTeamless(contestant models.Contestant, contest models.Contest) error {
	contestant.TeamlessedAt = time.Now()
	contestant.TeamlessContestID = contest.ID

	err := c.repo.Update(&contestant)
	if err != nil {
		return err
	}

	contest.TeamlessContestants = append(contest.TeamlessContestants, contestant)
	return c.contestRepo.Update(contest)
}

// CheckJoinedTeam reports whether the given contestant is in the given team, or any team at all
func (c *ContestantHelper) CheckJoinedTeam(cont models.Contestant, team models.Team) bool {
	return cont.TeamID > 1 || cont.TeamID == team.ID || c.jrMgr.CheckContestantTeamRequests(cont, team)
}

// GetTeam returns a team using the given id
func (c *ContestantHelper) GetTeam(contestant models.Contestant) (models.Team, error) {
	return c.teamMgr.GetTeam(contestant.TeamID)
}
