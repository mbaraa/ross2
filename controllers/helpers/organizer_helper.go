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
	repo            data.CRUDRepo[models.Organizer]
	contestRepo     data.Many2ManyCRUDRepo[models.Contest, any]
	ocRepo          data.OrganizeContestCRUDRepo
	userRepo        data.CRUDRepo[models.User]
	teamMgr         *TeamHelper
	notificationMgr *NotificationHelper
}

func NewOrganizerHelperBuilder() *OrganizerHelperBuilder {
	return new(OrganizerHelperBuilder)
}

func (b *OrganizerHelperBuilder) OrganizerRepo(o data.CRUDRepo[models.Organizer]) *OrganizerHelperBuilder {
	b.repo = o
	return b
}

func (b *OrganizerHelperBuilder) ContestRepo(c data.Many2ManyCRUDRepo[models.Contest, any]) *OrganizerHelperBuilder {
	b.contestRepo = c
	return b
}

func (b *OrganizerHelperBuilder) OrganizeContestRepo(c data.OrganizeContestCRUDRepo) *OrganizerHelperBuilder {
	b.ocRepo = c
	return b
}

func (b *OrganizerHelperBuilder) UserRepo(u data.CRUDRepo[models.User]) *OrganizerHelperBuilder {
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
	if b.ocRepo == nil {
		sb.WriteString("Organizer Helper Builder: missing organize contest repo!")
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
	repo            data.CRUDRepo[models.Organizer]
	contestRepo     data.Many2ManyCRUDRepo[models.Contest, any]
	ocRepo          data.OrganizeContestCRUDRepo
	userRepo        data.CRUDRepo[models.User]
	teamMgr         *TeamHelper
	notificationMgr *NotificationHelper
}

// NewOrganizerHelper returns a new OrganizerHelper instance
func NewOrganizerHelper(b *OrganizerHelperBuilder) *OrganizerHelper {
	return &OrganizerHelper{
		repo:            b.repo,
		contestRepo:     b.contestRepo,
		ocRepo:          b.ocRepo,
		userRepo:        b.userRepo,
		teamMgr:         b.teamMgr,
		notificationMgr: b.notificationMgr,
	}
}

// GetProfile returns organizer's profile for the given user
func (o *OrganizerHelper) GetProfile(user models.User) (models.Organizer, error) {
	return o.repo.Get(user.ID)
}

// GetUserProfileUsingEmail returns user's profile for the given user email
func (o *OrganizerHelper) GetUserProfileUsingEmail(userEmail string) (models.User, error) {
	users, err := o.userRepo.GetByConds("email = ?", userEmail)
	return users[0], err
}

// FinishProfile sets the organizer's profile after the first sign in after the promotion
func (o *OrganizerHelper) FinishProfile(org models.Organizer) error {
	org.User.ProfileStatus |= enums.ProfileStatusOrganizerFinished

	org.User.ContactInfo = o.validateProfile(org.User.ContactInfo)

	err := o.userRepo.Update(&org.User)
	if err != nil {
		return err
	}

	return o.repo.Update(&org)
}

func (o *OrganizerHelper) verifyProfile(ci models.ContactInfo) bool {
	return len(ci.FacebookURL) > len("https://") && len(ci.TelegramNumber) > len("https://") &&
		ci.FacebookURL[:len("https://")] == "https://" && ci.TelegramNumber[:len("https://")] == "https://"
}

func (o *OrganizerHelper) validateProfile(ci models.ContactInfo) models.ContactInfo {
	if !o.verifyProfile(ci) {
		if !strings.Contains(ci.FacebookURL, "http") {
			ci.FacebookURL = "https://" + ci.FacebookURL
		}
		if !strings.Contains(ci.TelegramNumber, "http") {
			ci.TelegramNumber = "https://" + ci.TelegramNumber
		}
	}
	return ci
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

	return o.ocRepo.Update(&orgContest)
}

// DeleteContest deletes contest, much wow!
func (o *OrganizerHelper) DeleteContest(contest models.Contest) error {
	return o.contestRepo.Delete(contest)
}

// UpdateContest you guessed it, much wow!
func (o *OrganizerHelper) UpdateContest(contest models.Contest) error {
	contest.StartsAt2 = time.UnixMilli(contest.StartsAt)
	contest.RegistrationEnds2 = time.UnixMilli(contest.RegistrationEnds)
	return o.contestRepo.Update(&contest)
}

// AddOrganizer adds the given organizer
func (o *OrganizerHelper) AddOrganizer(newOrg, director models.Organizer, baseUser models.User, contest models.Contest, roles enums.OrganizerRole) error {
	org, err := o.repo.GetByConds("email = ?", newOrg.User.Email)
	orgExists := err == nil

	if orgExists {
		newOrg.ID = org[0].ID
		_, err := o.ocRepo.Get(models.OrganizeContest{
			ContestID:   contest.ID,
			OrganizerID: newOrg.ID,
		})

		if err == nil {
			return errors.New("this user is already an organizer on this contest")
		}
	}

	baseUser.UserType |= enums.UserTypeOrganizer

	err = o.userRepo.Update(&baseUser)
	if err != nil {
		return err
	}

	newOrg.DirectorID = director.User.ID
	newOrg.Director = &director

	newOrg.User = baseUser
	newOrg.UserID = baseUser.ID

	if !orgExists {
		err = o.repo.Add(&newOrg)
		if err != nil {
			return err
		}
		return o.ocRepo.Update(&models.OrganizeContest{
			ContestID:      contest.ID,
			OrganizerID:    newOrg.ID,
			OrganizerRoles: roles,
		})
	}

	return o.ocRepo.Add(&models.OrganizeContest{
		ContestID:      contest.ID,
		OrganizerID:    newOrg.ID,
		OrganizerRoles: roles,
	})
}

// UpdateOrganizer updates the given organizer
func (o *OrganizerHelper) UpdateOrganizer(newOrg, director models.Organizer, baseUser models.User, contest models.Contest, roles enums.OrganizerRole) error {
	return o.ocRepo.Update(&models.OrganizeContest{
		ContestID:      contest.ID,
		OrganizerID:    newOrg.ID,
		OrganizerRoles: roles,
	})
}

// DeleteOrganizer deletes the given organizer
func (o *OrganizerHelper) DeleteOrganizer(org models.Organizer, contest models.Contest) error {
	return o.ocRepo.Delete(models.OrganizeContest{
		ContestID:   contest.ID,
		OrganizerID: org.ID,
	})
}

// GenerateTeams generates teams for the teamless contestants of the given contest
func (o *OrganizerHelper) GenerateTeams(contest models.Contest, genType string, names []string) ([]models.Team, []models.Contestant, error) {
	var err error
	contest, err = o.contestRepo.Get(contest.ID)
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

// CreateUpdateTeams creates/updates the given teams :)
func (o *OrganizerHelper) CreateUpdateTeams(teams []models.Team, removedConts []models.Contestant, contest models.Contest, org models.Organizer) error {
	return o.teamMgr.CreateUpdateTeams(teams, removedConts, contest, org)
}

// // UpdateTeams updates the given teams after checking that the given organizer is a director on one of the contest that the teams are in
// func (o *OrganizerHelper) UpdateTeams(teams []models.Team, removedConts []models.Contestant, org models.Organizer) error {
// 	return o.teamMgr.CreateUpdateTeams(teams, removedConts, org)
// }

// GetContests returns the contests of the given organizer, and an occurring error
func (o *OrganizerHelper) GetContests(org models.Organizer) ([]models.Contest, error) {
	return o.contestRepo.GetByAssociation(org)
}

// GetContest well lol
func (o *OrganizerHelper) GetContest(contest models.Contest) (models.Contest, error) {
	return o.contestRepo.Get(contest.ID)
}

// GetOrganizers returns all organizers that are under the given organizer
func (o *OrganizerHelper) GetOrganizers(org models.Organizer, contest models.Contest) ([]models.Organizer, error) {
	orgs, err := o.ocRepo.GetOrgs(contest)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (o *OrganizerHelper) SendSheevNotifications(contest models.Contest) error {
	var err error
	contest, err = o.contestRepo.Get(contest.ID) // lazy loading :)
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
	contest, err := o.contestRepo.Get(contest.ID)
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
	isResp := o.CheckOrgRole(enums.RoleReceptionist, contest.ID, org.ID)
	if !o.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) && !isResp {
		return nil, errors.New("you can't do that :)")
	}

	contest, err = o.contestRepo.Get(contest.ID)
	if err != nil {
		return nil, err
	}

	for _, _org := range contest.Organizers {
		if _org.User.AttendedContestID != contest.ID && isResp {
			_org.Contests = nil
			parts = append(parts, _org.User)
			continue
		}
		_org.Contests = nil
		parts = append(parts, _org.User)
	}

	for _, team := range contest.Teams {
		for _, member := range team.Members {
			if member.User.AttendedContestID != contest.ID && isResp {
				parts = append(parts, member.User)
				continue
			}
			parts = append(parts, member.User)
		}
	}

	return
}

func (o *OrganizerHelper) MarkAttendance(user models.User, contest models.Contest) error {
	users, err := o.userRepo.GetByConds("email = ?",user.Email)
	if err != nil {
		return err
	}

	contest, err = o.contestRepo.Get(contest.ID)
	if err != nil {
		return err
	}

	users[0].AttendedContestID = contest.ID
	users[0].AttendedAt = time.Now()

	return o.userRepo.Update(&users[0])
}

// CheckOrgRole reports whether the given organizer has the given role over the given contest
func (o *OrganizerHelper) CheckOrgRole(role enums.OrganizerRole, contestID, organizerID uint) bool {
	oc, err := o.ocRepo.Get(models.OrganizeContest{
		ContestID:   contestID,
		OrganizerID: organizerID,
	})
	if err != nil {
		return false
	}

	return (role & oc.OrganizerRoles) != 0
}

func (o *OrganizerHelper) GetOrgRoles(orgID, contestID uint) (roles enums.OrganizerRole, rolesNames []string, err error) {
	oc, err := o.ocRepo.Get(models.OrganizeContest{
		ContestID:   contestID,
		OrganizerID: orgID,
	})
	if err != nil {
		return
	}

	roles = enums.OrganizerRole(oc.OrganizerRoles)
	rolesNames = roles.GetRoles()

	return
}

func (o *OrganizerHelper) GetTeamsCSV(contest models.Contest) (string, error) {
	cont, err := o.contestRepo.Get(contest.ID)
	if err != nil {
		return "", err
	}

	csv := new(strings.Builder)

	csv.WriteString("Team Name, ")

	for i := 1; i <= int(cont.ParticipationConditions.MaxTeamMembers); i++ {
		last := ','
		if i == int(cont.ParticipationConditions.MaxTeamMembers) {
			last = '\000'
		}
		csv.WriteString(fmt.Sprintf("Member #%d Name, Member #%d Uni ID%c ", i, i, last))
	}
	csv.WriteString(fmt.Sprintln())

	for _, team := range cont.Teams {
		teamRow := new(strings.Builder)
		teamRow.WriteString(team.Name + ", ")

		for i, member := range team.Members {
			last := ','
			if i == len(team.Members)-1 {
				last = '\000'
			}
			teamRow.WriteString(fmt.Sprintf("%s, %s%c ", member.User.Name, strings.Split(member.User.Email, "@")[0], last))
		}

		if len(team.Members) > 0 {
			csv.WriteString(fmt.Sprintln(teamRow.String()))
		}
	}

	return csv.String(), nil
}
