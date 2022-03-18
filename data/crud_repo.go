package data

import "gorm.io/gorm"

type CreatorRepo [T any] interface {
	Add(obj *T) error
	AddMany(objs []*T) error
}

type GetterRepo[T any] interface {
	Exists(id uint) bool
	Get(id uint) (T, error)
	GetByConds(conds ...any) ([]T, error)
	GetAll() ([]T, error)
	Count() (int64, error)
}

type UpdaterRepo[T any] interface {
	Update(obj *T, conds ...any) error
	UpdateAll(objs []*T, conds ...any) error
}

type DeleterRepo[T any] interface {
	Delete(obj T, conds ...any) error
	DeleteAll(conds ...any) error
}

type CRUDRepo[T any] interface {
	CreatorRepo[T]
	GetterRepo[T]
	UpdaterRepo[T]
	DeleterRepo[T]
	GetDB() *gorm.DB
}

type Many2ManyCRUDRepo[T any, T2 any] interface {
	CRUDRepo[T]
	// AppendMany2Many(src T, associated T2) error
	GetByAssociation(associated T2) ([]T, error)
}
