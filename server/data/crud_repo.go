package data

import "gorm.io/gorm"

type CreatorRepo[T any] interface {
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

type GORMDBGetter interface {
	GetDB() *gorm.DB
}

type CRUDRepo[T any] interface {
	GORMDBGetter
	CreatorRepo[T]
	GetterRepo[T]
	UpdaterRepo[T]
	DeleterRepo[T]
}
