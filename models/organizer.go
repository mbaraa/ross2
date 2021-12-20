package models

import (
	"github.com/mbaraa/ross2/models/enums"
	"gorm.io/gorm"
)

// Organizer represents a contest's organizer
type Organizer struct {
	gorm.Model
	User   User `gorm:"foreignkey:UserID" json:"user"`
	UserID uint `gorm:"column:user_id" json:"user_id"`

	DirectorID uint       `gorm:"column:director_id"`
	Director   *Organizer `gorm:"foreignkey:DirectorID" json:"director"`

	Contests   []Contest           `gorm:"many2many:register_contest" json:"contests"`
	Roles      enums.OrganizerRole `gorm:"column:roles;type:uint" json:"roles"`
	RolesNames []string            `gorm:"-" json:"roles_names"`
}

func (o *Organizer) AfterFind(db *gorm.DB) error {
	o.RolesNames = o.Roles.GetRoles()

	err := db.
		Model(new(User)).
		First(&o.User, "id = ?", o.UserID).
		Error

	if err != nil {
		return err
	}

	return db.
		Model(new(ContactInfo)).
		First(&o.User.ContactInfo, "id = ?", o.User.ContactInfoID).
		Error
}
