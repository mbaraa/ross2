package managers

import (
	"net/http"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type SessionManager struct {
	sessionRepo data.SessionCRUDRepo
}

func NewSessionManager(sessionRepo data.SessionCRUDRepo) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionManager) CreateSession(userID uint) error {
	return s.sessionRepo.Add(models.Session{
		UserID: userID,
	})
}

func (s *SessionManager) GetSession(token string) (models.Session, error) {
	return s.sessionRepo.Get(models.Session{ID: token})
}

func (s *SessionManager) CheckSessionFromRequest(req *http.Request) (models.Session, bool) {
	token, exists := req.Header["Authorization"]
	if exists {
		sess, err := s.GetSession(token[0])
		if err != nil {
			return models.Session{}, false
		}
		return sess, true
	}

	return models.Session{}, false
}

func (s *SessionManager) CheckSession(token string) bool {
	exists, _ := s.sessionRepo.Exists(models.Session{ID: token})
	return exists
}

func (s *SessionManager) DeleteSession(token string) error {
	return s.sessionRepo.Delete(models.Session{ID: token})
}

func (s *SessionManager) DeleteAllSessions(userID uint) error {
	return s.sessionRepo.DeleteAllForUser(userID)
}
