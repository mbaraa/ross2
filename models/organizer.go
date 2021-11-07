package models

import "gorm.io/gorm"

// Organizer represents a contest's organizer
type Organizer struct {
	gorm.Model
	ID              uint   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email           string `gorm:"column:email" json:"email"`
	Name            string `gorm:"column:name" json:"name"`
	AvatarURL       string `gorm:"column:avatar_url" json:"avatar_url"`
	ProfileFinished bool   `gorm:"profile_finished" json:"profile_finished"`

	ContactInfo   ContactInfo `gorm:"foreignkey:ContactInfoID" json:"contact_info"`
	ContactInfoID uint        `gorm:"column:contact_info_id"`

	DirectorID uint       `gorm:"column:director_id"`
	Director   *Organizer `gorm:"foreignkey:ContactInfoID" json:"director"`

	Contests   []Contest     `gorm:"many2many:register_contest" json:"contests"`
	Roles      OrganizerRole `gorm:"column:roles;type:uint" json:"roles"`
	RolesNames []string      `gorm:"-" json:"roles_names"`
}

func (o *Organizer) AfterFind(db *gorm.DB) error {
	err := db.
		First(&o.ContactInfo, "id = ?", o.ContactInfoID).
		Error

	o.RolesNames = getRoles(o.Roles)

	return err
}

func getRoles(roles OrganizerRole) []string {
	if roles == 1 {
		return []string{organizerRoleText[RoleDirector]}
	}

	rolesTexts := []string{"Organizer"}
	for i := 0; i <= 63; i++ {
		role := OrganizerRole(1 << i)
		if roles&role != 0 {
			rolesTexts = append(rolesTexts, organizerRoleText[role])
		}
	}

	return rolesTexts
}
