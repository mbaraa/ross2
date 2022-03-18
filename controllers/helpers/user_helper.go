package helpers

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

type UserHelper struct {
	repo data.UserCRUDRepo
	contRepo data.CRUDRepo[models.Contestant]
	sessMgr  *SessionHelper
}

func NewUserHelper(repo data.UserCRUDRepo, contRepo data.CRUDRepo[models.Contestant], sessMgr *SessionHelper) *UserHelper {
	return &UserHelper{
		repo: repo,
		contRepo: contRepo,
		sessMgr:  sessMgr,
	}
}

func (u *UserHelper) Login(user *models.User) (sess models.Session, err error) {
	fetchedUser, err := u.repo.GetByEmail(user.Email)
	if err != gorm.ErrRecordNotFound && err != nil {
		return models.Session{}, err
	}

	exists, _ := u.repo.Exists(fetchedUser)
	if !exists {
		err = u.Signup(user)
		fetchedUser = *user
	}

	sess, err = u.sessMgr.CreateSession(fetchedUser.ID)
	if err != nil {
		return models.Session{}, err
	}

	return
}

func (u *UserHelper) Signup(user *models.User) error {
	user.UserType |= enums.UserTypeContestant
	return u.contRepo.Add(&models.Contestant{
		User:   *user,
		UserID: user.ID,
		Team:   models.Team{ID: 1},
		TeamID: 1,
	})
	return u.repo.Add(user)
}

func (u *UserHelper) LoginUsingSession(sessionToken string) (user models.User, err error) {
	sess, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return
	}

	user, err = u.repo.Get(models.User{ID: sess.UserID})
	return
}

func (u *UserHelper) Logout(user models.User, sessionToken string) error {
	session, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return err
	}

	user, err = u.repo.Get(user)
	if err != nil {
		return err
	}

	if user.ID != session.UserID {
		return errors.New("user id doesn't match user id in session")
	}

	return u.sessMgr.DeleteSession(sessionToken)
}

func (u *UserHelper) UpdateUser(user *models.User) error {
	return u.repo.Update(user)
}
