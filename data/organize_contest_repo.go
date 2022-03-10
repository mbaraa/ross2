package data

import (
	"github.com/mbaraa/ross2/models"
)

// OrganizeContestCreatorRepo is an interface that allows creation of a organizer on a contest
// into a certain data source :)
type OrganizeContestCreatorRepo interface {
	Add(contest *models.OrganizeContest) error
}

// OrganizeContestGetterRepo is an interface that allows getting values of a organizer on a contest
// from a certain data source :)
type OrganizeContestGetterRepo interface {
	Exists(oc models.OrganizeContest) (bool, error)
	Get(oc models.OrganizeContest) (models.OrganizeContest, error)
	GetOrgs(contest models.Contest) ([]models.Organizer, error)
	GetContests(org models.Organizer) ([]models.Contest, error)
	Count() (int64, error)
}

// OrganizeContestUpdaterRepo is an interface that allows updating values of a organizer on a contest
// in a certain data source :)
type OrganizeContestUpdaterRepo interface {
	Update(oc *models.OrganizeContest) error
}

// OrganizeContestDeleterRepo is an interface that allows deleting values of a organizer on a contest
// from a certain data source :)
type OrganizeContestDeleterRepo interface {
	Delete(oc models.OrganizeContest) error
	DeleteAll() error
}

// OrganizeContestCRUDRepo is an interface that allows full CRUD operations of a organzier on a contest
// on a certain data source :)
type OrganizeContestCRUDRepo interface {
	OrganizeContestCreatorRepo
	OrganizeContestGetterRepo
	OrganizeContestUpdaterRepo
	OrganizeContestDeleterRepo
}
