package helpers

import (
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type SessionHelper struct {
	sessionRepo data.SessionCRUDRepo
}

func NewSessionHelper(sessionRepo data.SessionCRUDRepo) *SessionHelper {
	return &SessionHelper{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionHelper) CreateSession(userID uint) (models.Session, error) {
	sess := models.Session{
		UserID: userID,
	}
	err := s.sessionRepo.Add(&sess)

	return sess, err
}

func (s *SessionHelper) GetSession(token string) (models.Session, error) {
	return s.sessionRepo.Get(models.Session{ID: token})
}

func (s *SessionHelper) DeleteSession(token string) error {
	return s.sessionRepo.Delete(models.Session{ID: token})
}

func (s *SessionHelper) DeleteAllSessions(userID uint) error {
	return s.sessionRepo.DeleteAllForUser(userID)
}
