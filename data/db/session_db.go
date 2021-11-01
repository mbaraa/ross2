package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// SessionDB represents a CRUD db repo for sessions
type SessionDB struct {
	db *gorm.DB
}

// NewSessionDB returns a new SessionDB instance
func NewSessionDB(db *gorm.DB) *SessionDB {
	return &SessionDB{db: db}
}

// CREATOR REPO

func (s *SessionDB) Add(session *models.Session) error {
	return s.db.
		Create(session).
		Error
}

// GETTER REPO

func (s *SessionDB) Exists(session models.Session) (bool, error) {
	res := s.db.First(&session)
	return !errors.Is(res.Error, gorm.ErrRecordNotFound), res.Error
}

func (s *SessionDB) Get(session models.Session) (models.Session, error) {
	fetchedSession := models.Session{}
	err := s.db.First(&fetchedSession, "id = ?", session.ID).Error

	return fetchedSession, err
}

// DELETER REPO

func (s *SessionDB) Delete(session models.Session) error {
	return s.db.
		Where("id = ?", session.ID).
		Delete(&session).
		Error
}

func (s *SessionDB) DeleteAll() error {
	return s.db.
		Where("true").
		Delete(new(models.Session)).
		Error
}

func (s *SessionDB) DeleteAllForUser(userID uint) error {
	return s.db.
		Where("user_id = ?", userID).
		Delete(new(models.Session)).
		Error
}
