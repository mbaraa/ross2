package helpers

import (
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type SessionHelper[T models.Session] struct {
	sessionRepo data.CRUDRepo[models.Session]
}

func NewSessionHelper[T models.Session](sessionRepo data.CRUDRepo[models.Session]) *SessionHelper[T] {
	return &SessionHelper[T]{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionHelper[T]) CreateSession(userID uint) (models.Session, error) {
	sess := models.Session{
		UserID: userID,
	}
	err := s.sessionRepo.Add(&sess)

	return sess, err
}

func (s *SessionHelper[T]) GetSession(token string) (models.Session, error) {
	sess, err := s.sessionRepo.GetByConds("id = ?", token)
	if len(sess) == 0 {
		return models.Session{}, err
	}
	return sess[0], err
}

func (s *SessionHelper[T]) DeleteSession(token string) error {
	return s.sessionRepo.DeleteAll("token = ?", token)
}

func (s *SessionHelper[T]) DeleteAllSessions(userID uint) error {
	return s.sessionRepo.DeleteAll("user_id = ?", userID)
}
