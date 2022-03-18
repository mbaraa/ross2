package helpers

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
)

// AdminHelper well it's written on the box :)
type AdminHelper struct {
	repo     data.CRUDRepo[models.Admin]
	orgRepo  data.OrganizerCRUDRepo
	userRepo data.UserCRUDRepo
}

// NewAdminHelper returns a new AdminHelper instance
func NewAdminHelper(repo data.CRUDRepo[models.Admin], orgRepo data.OrganizerCRUDRepo, userRepo data.UserCRUDRepo) *AdminHelper {
	return &AdminHelper{repo, orgRepo, userRepo}
}

// GetProfile returns contestant's for the given user
func (a *AdminHelper) GetProfile(user models.User) (models.Admin, error) {
	return a.repo.Get(user.ID)
}

// AddDirector adds a director, much wow
func (a *AdminHelper) AddDirector(director models.Organizer, baseUser models.User) error {
	if (baseUser.UserType & enums.UserTypeDirector) != 0 {
		return errors.New("user is already a director")
	}

	baseUser.UserType |= enums.UserTypeOrganizer | enums.UserTypeDirector

	err := a.userRepo.Update(&baseUser)
	if err != nil {
		return err
	}

	director.User = baseUser
	director.UserID = baseUser.ID

	return a.orgRepo.Add(&director)
}

// DeleteDirector deletes a director, much wow
func (a *AdminHelper) DeleteDirector(director models.Organizer) error {
	if (director.User.UserType & enums.UserTypeDirector) != 0 {
		director.User.UserType -= enums.UserTypeOrganizer | enums.UserTypeDirector

		if (director.User.ProfileStatus & enums.ProfileStatusOrganizerFinished) != 0 {
			director.User.ProfileStatus -= enums.ProfileStatusOrganizerFinished
		}

		err := a.userRepo.Update(&director.User)
		if err != nil {
			return err
		}
	}

	return a.orgRepo.Delete(director)
}

// GetDirectors returns all directors registered, and an occurring error
func (a *AdminHelper) GetDirectors() (dirs []models.Organizer, err error) {
	orgs, err := a.orgRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		if (org.User.UserType & enums.UserTypeDirector) != 0 {
			dirs = append(dirs, org)
		}
	}
	return
}

// GetUserProfileUsingEmail returns user's profile for the given user email
func (a *AdminHelper) GetUserProfileUsingEmail(userEmail string) (models.User, error) {
	return a.userRepo.GetByEmail(userEmail)
}
