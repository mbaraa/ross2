package db

import (
	"errors"

	"github.com/mbaraa/ross2/models"
	"gorm.io/gorm"
)

// SessionDB represents a CRUD db repo for sessions
type SessionDB[T models.Session] struct {
	db *gorm.DB
}

// NewSessionDB returns a new SessionDB instance
func NewSessionDB[T models.Session](db *gorm.DB) *SessionDB[T] {
	return &SessionDB[T]{db: db}
}

func (s *SessionDB[T]) GetDB() *gorm.DB {
	return s.db
}

// CREATOR REPO

func (s *SessionDB[T]) Add(session *models.Session) error {
	return s.db.
		Create(session).
		Error
}

func (s *SessionDB[T]) AddMany(sessions []*models.Session) error {
	return errors.New("not and will never be implemented :)")
}

// GETTER REPO

func (s *SessionDB[T]) Exists(id uint) bool {
	_, err := s.Get(id)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (s *SessionDB[T]) Get(id uint) (models.Session, error) {
	return models.Session{}, errors.New("not implemented")
}

func (s *SessionDB[T]) GetByConds(conds ...any) (sessions []models.Session, err error) {
	if (len(conds) < 2) {
		return nil, errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	err = s.db.
		Model(new(models.Session)).
		Find(&sessions, conds[0], conds[1:]).
		Error

	return
}

func (s *SessionDB[T]) GetAll() ([]models.Session, error) {
	return nil, errors.New("not implemented")
}

func (s *SessionDB[T]) Count() (int64, error) {
	return 0, errors.New("not implemented")
}

// UPDATER REPO

func (s *SessionDB[T]) Update(session *models.Session, conds ...any) error {
	return errors.New("not and will never be implemented :)")
}

func (s *SessionDB[T]) UpdateAll(sessions []*models.Session, conds ...any) error {
	return errors.New("not and will never be implemented :)")
}

// DELETER REPO

func (s *SessionDB[T]) Delete(session models.Session, conds ...any) error {
	return s.db.
		Where("id = ?", session.ID).
		Delete(&session).
		Error
}

func (s *SessionDB[T]) DeleteAll(conds ...any) error {
	if (len(conds) < 2) {
		return errors.New("conditions should be at least 2, ie condition string and the associated value")
	}

	return s.db.
		Where(conds[0], conds[1:]).
		Delete(new(models.Session)).
		Error
}

