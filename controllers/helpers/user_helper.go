package helpers

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

type UserHelper struct {
	userRepo data.UserCRUDRepo
	contRepo data.ContestantCRUDRepo
	sessMgr  *SessionHelper
}

func NewUserHelper(userRepo data.UserCRUDRepo, contRepo data.ContestantCRUDRepo, sessMgr *SessionHelper) *UserHelper {
	return &UserHelper{
		userRepo: userRepo,
		contRepo: contRepo,
		sessMgr:  sessMgr,
	}
}

func (u *UserHelper) Login(user *models.User) (sess models.Session, err error) {
	fetchedUser, err := u.userRepo.GetByEmail(user.Email)
	if err != gorm.ErrRecordNotFound && err != nil {
		return models.Session{}, err
	}

	exists, _ := u.userRepo.Exists(fetchedUser)
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
	return u.userRepo.Add(user)
}

func (u *UserHelper) LoginUsingSession(sessionToken string) (user models.User, err error) {
	sess, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return
	}

	user, err = u.userRepo.Get(models.User{ID: sess.UserID})
	return
}

func (u *UserHelper) Logout(user models.User, sessionToken string) error {
	session, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return err
	}

	user, err = u.userRepo.Get(user)
	if err != nil {
		return err
	}

	if user.ID != session.UserID {
		return errors.New("user id doesn't match user id in session")
	}

	return u.sessMgr.DeleteSession(sessionToken)
}

func (u *UserHelper) UpdateUser(user *models.User) error {
	return u.userRepo.Update(user)
}
