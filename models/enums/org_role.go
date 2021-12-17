package enums

// OrganizerRole represents the organizer's role in a contest
type OrganizerRole uint64

// OrganizerRole constants
const (
	RoleDirector OrganizerRole = 1 << iota
	RoleCoreOrganizer
	RoleChiefJudge
	RoleJudge
	RoleTechnical
	RoleCoordinator
	RoleMedia
	RoleBalloons
	RoleFood
	RoleReceptionist
)

var organizerRoleText = map[OrganizerRole]string{
	RoleDirector:      "Director",
	RoleCoreOrganizer: "Core Organizer",
	RoleChiefJudge:    "Chief Judge",
	RoleJudge:         "Judge",
	RoleMedia:         "Media",
	RoleBalloons:      "Balloons",
	RoleTechnical:     "Technical",
	RoleCoordinator:   "Coordinator",
	RoleFood:          "Food",
	RoleReceptionist:  "Receptionist",
}

// GetRoles returns a string slice that represents the organizer's roles
func (roles OrganizerRole) GetRoles() []string {
	if roles == 1 {
		return []string{organizerRoleText[RoleDirector]}
	}

	rolesTexts := []string{"Organizer"}
	for role := RoleDirector; role <= RoleReceptionist; role <<= 1 {
		if roles&role != 0 {
			rolesTexts = append(rolesTexts, organizerRoleText[role])
		}
	}

	return rolesTexts
}
