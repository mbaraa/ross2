package managers

import (
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/multiavatar"
)

type OrganizerManager struct {
	orgRepo     data.OrganizerCRUDRepo
	sessMgr     *SessionManager
	contestRepo data.ContestCRUDRepo
}

func NewOrganizerManager(orgRepo data.OrganizerCRUDRepo, sessMgr *SessionManager,
	contestRepo data.ContestCRUDRepo) *OrganizerManager {
	return &OrganizerManager{
		orgRepo:     orgRepo,
		sessMgr:     sessMgr,
		contestRepo: contestRepo,
	}
}

func (o *OrganizerManager) CreateContest(contest models.Contest) error {
	return o.contestRepo.Add(contest)
}

func (o *OrganizerManager) AddOrganizer(org *models.Organizer) error {
	org.User.AvatarURL = multiavatar.GetAvatarURL()
	return o.orgRepo.Add(org)
}

func (o *OrganizerManager) GetOrganizer(sessionToken string) (models.Organizer, error) {
	session, err := o.sessMgr.GetSession(sessionToken)
	if err != nil {
		return models.Organizer{}, err
	}

	return o.orgRepo.Get(models.Organizer{
		User: models.User{ID: session.UserID},
	})
}

func (o *OrganizerManager) GetContests(org models.Organizer) ([]models.Contest, error) {
	return o.contestRepo.GetAllByOrganizer(org)
}

func (o *OrganizerManager) GetOrganizers(org models.Organizer) ([]models.Organizer, error) {
	return o.orgRepo.GetAllByOrganizer(org)
}

func (o *OrganizerManager) UpdateProfile(org models.Organizer) error {
	return o.orgRepo.Update(org)
}

func (o *OrganizerManager) UpdateContest(contest models.Contest) error {
	return o.contestRepo.Update(contest)
}

func (o *OrganizerManager) DeleteOrganizer(org models.Organizer) error {
	return o.orgRepo.Delete(org)
}

func (o *OrganizerManager) DeleteContest(contest models.Contest) error {
	return o.contestRepo.Delete(contest)
}
