package helpers

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/multiavatar"
)

type ContestantHelperBuilder struct {
	userRepo         data.CRUDRepo[models.User]
	contestantRepo   data.CRUDRepo[models.Contestant]
	contestRepo      data.CRUDRepo[models.Contest]
	notificationRepo data.CRUDRepo[models.Notification]
	rtRepo           data.RegisterTeamCRUDRepo
	teamMgr          *TeamHelper
	jrMgr            *JoinRequestHelper
}

func NewContestantHelperBuilder() *ContestantHelperBuilder {
	return new(ContestantHelperBuilder)
}

func (b *ContestantHelperBuilder) UserRepo(u data.CRUDRepo[models.User]) *ContestantHelperBuilder {
	b.userRepo = u
	return b
}

func (b *ContestantHelperBuilder) ContestantRepo(c data.CRUDRepo[models.Contestant]) *ContestantHelperBuilder {
	b.contestantRepo = c
	return b
}

func (b *ContestantHelperBuilder) ContestRepo(c data.CRUDRepo[models.Contest]) *ContestantHelperBuilder {
	b.contestRepo = c
	return b
}

func (b *ContestantHelperBuilder) NotificationRepo(n data.CRUDRepo[models.Notification]) *ContestantHelperBuilder {
	b.notificationRepo = n
	return b
}

func (b *ContestantHelperBuilder) RegisterTeamRepo(rt data.RegisterTeamCRUDRepo) *ContestantHelperBuilder {
	b.rtRepo = rt
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
	if b.rtRepo == nil {
		sb.WriteString("Contestant Helper Builder: missing register team repo!")
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
	repo             data.CRUDRepo[models.Contestant]
	userRepo         data.UpdaterRepo[models.User]
	contestRepo      data.CRUDRepo[models.Contest]
	notificationRepo data.CRUDRepo[models.Notification]
	rtRepo           data.RegisterTeamCRUDRepo
	teamMgr          *TeamHelper
	jrMgr            *JoinRequestHelper
}

// NewContestantHelper returns a new ContestantHelper instance
func NewContestantHelper(b *ContestantHelperBuilder) *ContestantHelper {
	return &ContestantHelper{
		repo:             b.contestantRepo,
		userRepo:         b.userRepo,
		contestRepo:      b.contestRepo,
		notificationRepo: b.notificationRepo,
		rtRepo:           b.rtRepo,
		teamMgr:          b.teamMgr,
		jrMgr:            b.jrMgr,
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
	return c.repo.Get(user.ID)
}

// CreateTeam creates a team and adds the given contestant to it as its leader
func (c *ContestantHelper) CreateTeam(contestant models.Contestant, team models.Team, contest models.Contest) error {
	contestant.TeamlessContestID = math.MaxInt64
	contestant.TeamlessedAt = contestant.CreatedAt

	err := c.teamMgr.CreateTeam(contestant, &team)
	if err != nil {
		return err
	}

	err = c.jrMgr.DeleteRequests(contestant.User.ID, 0)
	if err != nil {
		return err
	}

	return c.rtRepo.Add(&models.RegisterTeam{
		ContestantID: contestant.ID,
		Team:         team,
		TeamID:       team.ID,
		ContestID:    contest.ID,
	})
}

// DeleteTeam kicks every member out of it and deletes it
func (c *ContestantHelper) DeleteTeam(cont models.Contestant, team models.Team) error {
	return c.teamMgr.DeleteTeam(cont, team)
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
func (c *ContestantHelper) LeaveTeam(contestant models.Contestant, team models.Team) error {
	return c.teamMgr.LeaveTeam(contestant, team)
}

// RegisterAsTeamless adds the given contestant as teamless for the given contest
func (c *ContestantHelper) RegisterAsTeamless(contestant models.Contestant, contest models.Contest) (err error) {
	contest, err = c.contestRepo.Get(contest.ID)
	if err != nil {
		return err
	}

	if contest.StartsAt2.Add(contest.Duration).Before(time.Now()) {
		return errors.New("contest registration is over")
	}

	contestant.TeamlessedAt = time.Now()
	contestant.TeamlessContestID = contest.ID

	err = c.repo.Update(&contestant)
	if err != nil {
		return err
	}

	contest.TeamlessContestants = append(contest.TeamlessContestants, contestant)
	return c.contestRepo.Update(&contest)
}

// CheckJoinedTeam reports whether the given contestant is in the given team, or any team at all
func (c *ContestantHelper) CheckJoinedTeam(cont models.Contestant, team models.Team) bool {
	return cont.TeamID > 1 || cont.TeamID == team.ID || c.jrMgr.CheckContestantTeamRequests(cont, team)
}

// GetTeams returns a team using the given id
func (c *ContestantHelper) GetTeams(contestant models.Contestant) ([]models.RegisterTeam, error) {
	return c.rtRepo.Get(models.RegisterTeam{
		ContestantID: contestant.ID,
	}, "contestant_id = ?", contestant.ID)
}

// CheckJoinedContest reports whether the contestant is in the given contest or not
func (c *ContestantHelper) CheckJoinedContest(contest models.Contest, contestant models.Contestant) bool {
	rts, _ := c.rtRepo.Get(models.RegisterTeam{
		ContestID:    contest.ID,
		ContestantID: contestant.ID,
	}, "contest_id = ? and contestant_id = ?", contest.ID, contestant.ID)

	return len(rts) > 0
}

// RegisterInContest adds the given contestant's team to the given contest
func (c *ContestantHelper) RegisterInContest(contest models.Contest, contestant models.Contestant) error {
	team, err := c.teamMgr.GetTeam(contestant.TeamID)
	if err != nil {
		return err
	}

	if team.LeaderId != contestant.UserID {
		return errors.New("only the team's creator can join contests")
	}

	contest, err = c.contestRepo.Get(contest.ID)
	if err != nil {
		return err
	}

	err = c.checkCollidingContests(&team, contest)
	if err != nil {
		return err
	}

	if contest.RegistrationEnds2.Before(time.Now()) {
		return errors.New("registration for this contest is over")
	}

	err = c.contestRepo.GetDB().Model(&team).Association("Contests").Append(&contest)
	if err != nil {
		return err
	}

	return errors.New("all good")
}

func (c *ContestantHelper) checkCollidingContests(team *models.Team, contest models.Contest) error {
	contests := []models.Contest{}

	err := c.contestRepo.
		GetDB().
		Model(&team).
		Association("Contests").
		Find(&contests)

	if err != nil {
		return err
	}

	cStarts := contest.StartsAt
	cEnds := cStarts + int64(contest.Duration*60*1e3)

	for i := 0; i < len(contests); i++ {
		starts := contests[i].StartsAt
		ends := starts + int64(contests[i].Duration*60*1e3)

		if (cStarts <= starts && starts <= cEnds) || (ends <= cEnds && ends >= cStarts) {
			return fmt.Errorf(`can't join the contest "%s", because it collids with the contest "%s"`, contest.Name, contests[i].Name)
		}
	}

	team.Contests = contests

	return nil
}

func (c *ContestantHelper) GetTeamByJoinID(joinID string) (models.Team, error) {
	return c.teamMgr.GetTeamByJoinID(joinID)
}

func (c *ContestantHelper) RegisterInContestUsingTeam(contest models.Contest, team models.Team, cont models.Contestant) error {
	return c.rtRepo.Add(&models.RegisterTeam{
		ContestID:    contest.ID,
		TeamID:       team.ID,
		ContestantID: cont.ID,
	})
}
