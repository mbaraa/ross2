package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	err := o.contestRepo.Add(&contest)
	if err != nil {
		return err
	}

	orgContest := models.OrganizeContest{
		ContestID:      contest.ID,
		OrganizerID:    org.ID,
		OrganizerRoles: enums.RoleDirector,
	}

	return o.repo.
		GetDB().
		Model(&orgContest).
		Where("contest_id = ? and organizer_id = ?", contest.ID, org.ID).
		Updates(&orgContest).
		Error
}

// DeleteContest deletes contest, much wow!
func (o *OrganizerHelper) DeleteContest(contest models.Contest) error {
	return o.contestRepo.Delete(contest)
}

// UpdateContest you guessed it, much wow!
func (o *OrganizerHelper) UpdateContest(contest models.Contest) error {
	contest.StartsAt2 = time.UnixMilli(contest.StartsAt)
	contest.RegistrationEnds2 = time.UnixMilli(contest.RegistrationEnds)
	return o.contestRepo.Update(contest)
}

// AddOrganizer adds the given organizer
func (o *OrganizerHelper) AddOrganizer(newOrg, director models.Organizer, baseUser models.User, contest models.Contest, roles enums.OrganizerRole) error {
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

	err = o.repo.Add(&newOrg)
	if err != nil {
		return err
	}

	return o.repo.
		GetDB().
		Model(new(models.OrganizeContest)).
		Where("contest_id = ? and organizer_id = ?", contest.ID, newOrg.ID).
		Updates(models.OrganizeContest{
			ContestID:      contest.ID,
			OrganizerID:    newOrg.ID,
			OrganizerRoles: roles,
		}).Error
}

// UpdateOrganizer updates the given organizer
func (o *OrganizerHelper) UpdateOrganizer(newOrg, director models.Organizer, baseUser models.User, contest models.Contest, roles enums.OrganizerRole) error {
	baseUser.UserType |= enums.UserTypeOrganizer

	err := o.userRepo.Update(&baseUser)
	if err != nil {
		return err
	}

	newOrg.DirectorID = director.User.ID
	newOrg.Director = &director

	newOrg.User = baseUser
	newOrg.UserID = baseUser.ID

	err = o.repo.Update(newOrg)
	if err != nil {
		return err
	}

	return o.repo.
		GetDB().
		Model(new(models.OrganizeContest)).
		Where("contest_id = ? and organizer_id = ?", contest.ID, newOrg.ID).
		Updates(models.OrganizeContest{
			ContestID:      contest.ID,
			OrganizerID:    newOrg.ID,
			OrganizerRoles: roles,
		}).Error
}

// DeleteOrganizer deletes the given organizer
func (o *OrganizerHelper) DeleteOrganizer(org models.Organizer) error {
	if (org.User.UserType & enums.UserTypeOrganizer) != 0 {
		org.User.UserType -= enums.UserTypeOrganizer
		if (org.User.ProfileStatus & enums.ProfileStatusOrganizerFinished) != 0 {
			org.User.ProfileStatus -= enums.ProfileStatusOrganizerFinished
		}

		err := o.userRepo.Update(&org.User)
		if err != nil {
			return err
		}
	}
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
func (o *OrganizerHelper) GetOrganizers(org models.Organizer, contest models.Contest) (orgs []models.Organizer, err error) {
	// return o.repo.GetAllByOrganizer(org)

	var orgsContests []models.OrganizeContest

	err = o.repo.
		GetDB().
		Model(new(models.OrganizeContest)).
		Where("contest_id = ?", contest.ID).
		Find(&orgsContests).
		Error

	if err != nil {
		return
	}

	for _, orgContest := range orgsContests {
		fetchedOrg := models.Organizer{}
		err := o.repo.
			GetDB().
			Model(new(models.Organizer)).
			Where("id = ?", orgContest.OrganizerID).
			First(&fetchedOrg).
			Error

		if err != nil {
			continue
		}

		orgs = append(orgs, fetchedOrg)
	}

	return
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

func (o *OrganizerHelper) GetParticipants(contest models.Contest, org models.Organizer) (parts []models.User, err error) {
	if !o.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) || !o.CheckOrgRole(enums.RoleReceptionist, contest.ID, org.ID) {
		return nil, errors.New("you can't do that :)")
	}

	contest, err = o.contestRepo.Get(contest)
	if err != nil {
		return nil, err
	}

	for _, _org := range contest.Organizers {
		if _org.User.AttendedContestID != contest.ID {
			_org.Contests = nil
			parts = append(parts, _org.User)
		}
	}

	for _, team := range contest.Teams {
		for _, member := range team.Members {
			if member.User.AttendedContestID != contest.ID {
				parts = append(parts, member.User)
			}
		}
	}

	return
}

func (o *OrganizerHelper) MarkAttendance(user models.User, contest models.Contest) (err error) {
	user, err = o.userRepo.GetByEmail(user.Email)
	if err != nil {
		return
	}

	contest, err = o.contestRepo.Get(contest)
	if err != nil {
		return
	}

	user.AttendedContestID = contest.ID
	user.AttendedAt = time.Now()

	return o.userRepo.Update(&user)
}

// CheckOrgRole reports whether the given organizer has the given role over the given contest
func (o *OrganizerHelper) CheckOrgRole(role enums.OrganizerRole, contestID, organizerID uint) bool {
	oc := models.OrganizeContest{}
	err := o.repo.
		GetDB().
		Model(new(models.OrganizeContest)).
		Where("contest_id = ? and organizer_id = ?", contestID, organizerID).
		First(&oc).
		Error

	if err != nil {
		return false
	}

	return (role & oc.OrganizerRoles) != 0
}

func (o *OrganizerHelper) GetOrgRoles(orgID, contestID uint) (roles []string, err error) {
	oc := models.OrganizeContest{}
	err = o.repo.
		GetDB().
		Model(new(models.OrganizeContest)).
		First(&oc, "contest_id = ? and organizer_id = ?", contestID, orgID).
		Error
	if err != nil {
		return
	}

	rolesNum := enums.OrganizerRole(oc.OrganizerRoles)
	roles = rolesNum.GetRoles()

	return
}
