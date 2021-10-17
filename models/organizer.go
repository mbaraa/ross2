package models

import "gorm.io/gorm"

// Organizer represents a contest's organizer
type Organizer struct {
	User
	Contests []Contest `gorm:"many2many:register_contest" json:"contests"`
	// ProfileFinished bool          `gorm:"column:profile_finished" json:"profile_finished"`
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
	if roles == 0 {
		return []string{organizerRoleText[RoleDirector]}
	}

	rolesTexts := make([]string, 0)
	for i := 0; i <= 63; i++ {
		role := OrganizerRole(1 << i)
		if roles&role != 0 {
			rolesTexts = append(rolesTexts, organizerRoleText[role])
		}
	}

	return rolesTexts
}
