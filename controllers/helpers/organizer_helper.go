package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/namesgetter"
	"github.com/mbaraa/ross2/utils/partsexport"
	"github.com/mbaraa/ross2/utils/sheevhelper"
	"github.com/mbaraa/ross2/utils/teamsgen"
)

type OrganizerHelperBuilder struct {
	repo            data.OrganizerCRUDRepo
	contestRepo     data.ContestCRUDRepo
	userRepo        data.UserCRUDRepo
	teamMgr         *TeamHelper
	notificationMgr *NotificationHelper
}

func NewOrganizerHelperBuilder() *OrganizerHelperBuilder {
	return new(OrganizerHelperBuilder)
}

func (b *OrganizerHelperBuilder) OrganizerRepo(o data.OrganizerCRUDRepo) *OrganizerHelperBuilder {
	b.repo = o
	return b
}

func (b *OrganizerHelperBuilder) ContestRepo(c data.ContestCRUDRepo) *OrganizerHelperBuilder {
	b.contestRepo = c
	return b
}

func (b *OrganizerHelperBuilder) UserRepo(u data.UserCRUDRepo) *OrganizerHelperBuilder {
	b.userRepo = u
	return b
}

func (b *OrganizerHelperBuilder) TeamMgr(t *TeamHelper) *OrganizerHelperBuilder {
	b.teamMgr = t
	return b
}

func (b *OrganizerHelperBuilder) NotificationMgr(n *NotificationHelper) *OrganizerHelperBuilder {
	b.notificationMgr = n
	return b
}

func (b *OrganizerHelperBuilder) verify() bool {
	sb := new(strings.Builder)

	if b.repo == nil {
		sb.WriteString("Organizer Helper Builder: missing organizer repo!")
	}
	if b.contestRepo == nil {
		sb.WriteString("Organizer Helper Builder: missing contest repo!")
	}
	if b.userRepo == nil {
		sb.WriteString("Organizer Helper Builder: missing user repo!")
	}
	if b.teamMgr == nil {
		sb.WriteString("Organizer Helper Builder: missing team repo!")
	}
	if b.notificationMgr == nil {
		sb.WriteString("Organizer Helper Builder: missing notification repo!")
	}

	if sb.Len() != 0 {
		fmt.Println(sb.String())
		return false
	}

	return true
}

func (b *OrganizerHelperBuilder) GetOrganizerManager() *OrganizerHelper {
	return NewOrganizerHelper(b)
}

// OrganizerHelper well hmm
type OrganizerHelper struct {
	repo            data.OrganizerCRUDRepo
	contestRepo     data.ContestCRUDRepo
	userRepo        data.UserCRUDRepo
	teamMgr         *TeamHelper
	notificationMgr *NotificationHelper
}

// NewOrganizerHelper returns a new OrganizerHelper instance
func NewOrganizerHelper(b *OrganizerHelperBuilder) *OrganizerHelper {
	return &OrganizerHelper{
		repo:            b.repo,
		contestRepo:     b.contestRepo,
		userRepo:        b.userRepo,
		teamMgr:         b.teamMgr,
		notificationMgr: b.notificationMgr,
	}
}

// GetProfile returns organizer's profile for the given user
func (o *OrganizerHelper) GetProfile(user models.User) (models.Organizer, error) {
	return o.repo.Get(models.Organizer{User: user})
}

// GetUserProfileUsingEmail returns user's profile for the given user email
func (o *OrganizerHelper) GetUserProfileUsingEmail(userEmail string) (models.User, error) {
	return o.userRepo.GetByEmail(userEmail)
}

// FinishProfile sets the organizer's profile after the first sign in after the promotion
func (o *OrganizerHelper) FinishProfile(org models.Organizer) error {
	org.User.ProfileStatus |= enums.ProfileStatusOrganizerFinished
	err := o.userRepo.Update(&org.User)
	if err != nil {
		return err
	}

	return o.repo.Update(org)
}

// UpdateTeam updates the given team after checking that the given organizer is a director on one of the contest that the team is in
func (o *OrganizerHelper) UpdateTeam(team models.Team, org models.Organizer) error {
	return o.teamMgr.UpdateTeam(team, org)
}

// CreateContest creates a new contest, much wow!
func (o *OrganizerHelper) CreateContest(contest models.Contest, org models.Organizer) error {
	contest.Organizers = []models.Organizer{org}
	return o.contestRepo.Add(contest)
}

// DeleteContest deletes contest, much wow!
func (o *OrganizerHelper) DeleteContest(contest models.Contest) error {
	return o.contestRepo.Delete(contest)
}

// UpdateContest you guessed it, much wow!
func (o *OrganizerHelper) UpdateContest(contest models.Contest) error {
	return o.contestRepo.Update(contest)
}

// AddOrganizer adds the given organizer
func (o *OrganizerHelper) AddOrganizer(newOrg, director models.Organizer, baseUser models.User) error {
	if (baseUser.UserType&enums.UserTypeOrganizer) != 0 ||
		(baseUser.UserType&enums.UserTypeDirector) != 0 {
		return errors.New("user is already an organizer")
	}

	baseUser.UserType |= enums.UserTypeOrganizer

	err := o.userRepo.Update(&baseUser)
	if err != nil {
		return err
	}

	newOrg.DirectorID = director.User.ID
	newOrg.Director = &director

	newOrg.User = baseUser
	newOrg.UserID = baseUser.ID

	return o.repo.Add(&newOrg)
}

// DeleteOrganizer deletes the given organizer
func (o *OrganizerHelper) DeleteOrganizer(org models.Organizer) error {
	return o.repo.Delete(org)
}

// GenerateTeams generates teams for the teamless contestants of the given contest
func (o *OrganizerHelper) GenerateTeams(contest models.Contest, genType string, names []string) ([]models.Team, []models.Contestant, error) {
	var err error
	contest, err = o.contestRepo.Get(contest)
	if err != nil {
		return nil, nil, err
	}

	teams, leftTeamless :=
		teamsgen.GenerateTeams(contest, // the big ass function that am proud AF from :)
			namesgetter.GetNamesGetter(genType, names...))
	if len(teams) == 0 && len(leftTeamless) == 0 {
		return nil, nil, errors.New("no contestants were found")
	}

	return teams, leftTeamless, nil
}

// CreateTeams creates the given teams :)
func (o *OrganizerHelper) CreateTeams(teams []*models.Team) error {
	return o.teamMgr.CreateTeams(teams)
}

// UpdateTeams updates the given teams after checking that the given organizer is a director on one of the contest that the teams are in
func (o *OrganizerHelper) UpdateTeams(teams []models.Team, removedConts []models.Contestant, org models.Organizer) error {
	return o.teamMgr.UpdateTeams(teams, removedConts, org)
}

// GetContests returns the contests of the given organizer, and an occurring error
func (o *OrganizerHelper) GetContests(org models.Organizer) ([]models.Contest, error) {
	return o.contestRepo.GetAllByOrganizer(org)
}

// GetContest well lol
func (o *OrganizerHelper) GetContest(contest models.Contest) (models.Contest, error) {
	return o.contestRepo.Get(contest)
}

// GetOrganizers returns all organizers that are under the given organizer
func (o *OrganizerHelper) GetOrganizers(org models.Organizer) ([]models.Organizer, error) {
	return o.repo.GetAllByOrganizer(org)
}

func (o *OrganizerHelper) SendSheevNotifications(contest models.Contest) error {
	var err error
	contest, err = o.contestRepo.Get(contest) // lazy loading :)
	if err != nil {
		return err
	}

	notifications, err := sheevhelper.GetSheevNotifications(contest)
	if err != nil {
		return err
	}

	return o.notificationMgr.SendMany(notifications)
}

// GetParticipantsCSV returns a string with the participated contestants and organizers
func (o *OrganizerHelper) GetParticipantsCSV(contest models.Contest) (string, error) {
	contest, err := o.contestRepo.Get(contest)
	return partsexport.GetParticipants(contest), err
}

// GetNonOrgUsers returns all contestants or newly signed-up users so that a director can make some of them as organizers
func (o *OrganizerHelper) GetNonOrgUsers() ([]models.User, error) {
	users, err := o.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var nonOrgUsers []models.User
	for _, user := range users {
		if (user.UserType&enums.UserTypeOrganizer) == 0 &&
			(user.UserType&enums.UserTypeDirector) == 0 &&
			(user.UserType&enums.UserTypeAdmin) == 0 {
			nonOrgUsers = append(nonOrgUsers, user)
		}
	}

	return nonOrgUsers, nil
}
