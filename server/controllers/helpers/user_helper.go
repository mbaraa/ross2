package helpers

import (
	"errors"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

type UserHelper struct {
	repo     data.CRUDRepo[models.User]
	contRepo data.CRUDRepo[models.Contestant]
	sessMgr  *SessionHelper[models.Session]
}

func NewUserHelper(repo data.CRUDRepo[models.User], contRepo data.CRUDRepo[models.Contestant], sessMgr *SessionHelper[models.Session]) *UserHelper {
	return &UserHelper{
		repo:     repo,
		contRepo: contRepo,
		sessMgr:  sessMgr,
	}
}

func (u *UserHelper) Login(user *models.User) (sess models.Session, err error) {
	fetchedUsers, err := u.repo.GetByConds("email = ?", user.Email)
	if err != gorm.ErrRecordNotFound && err != nil {
		return
	}

	exists := u.repo.Exists(fetchedUsers[0].ID)
	if !exists {
		err = u.Signup(user)
		fetchedUsers[0] = *user
	}

	sess, err = u.sessMgr.CreateSession(fetchedUsers[0].ID)
	if err != nil {
		return
	}

	return
}

func (u *UserHelper) Signup(user *models.User) error {
	user.UserType |= enums.UserTypeContestant
	return u.contRepo.Add(&models.Contestant{
		User:   *user,
		UserID: user.ID,
	})
}

func (u *UserHelper) LoginUsingSession(sessionToken string) (user models.User, err error) {
	sess, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return
	}

	user, err = u.repo.Get(sess.UserID)
	return
}

func (u *UserHelper) Logout(user models.User, sessionToken string) error {
	session, err := u.sessMgr.GetSession(sessionToken)
	if err != nil {
		return err
	}

	user, err = u.repo.Get(user.ID)
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
