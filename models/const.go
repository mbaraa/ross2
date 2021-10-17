package models

// Major represents a major enum
type Major uint64

// Major constants
// TODO:
// add more majors :)
const (
	MajorSoftwareEngineering Major = 1 << iota
	MajorComputerScience
	MajorCyberSecurity
)

// Major constants parent faculties
const (
	FacultyInformationTechnology = MajorSoftwareEngineering | MajorComputerScience | MajorCyberSecurity

	MajorAny = FacultyInformationTechnology + 1
)

var majorText = map[Major]string{
	MajorSoftwareEngineering:     "Software Engineering",
	MajorComputerScience:         "Computer Science",
	MajorCyberSecurity:           "Cyber Security",
	FacultyInformationTechnology: "Information Technology",

	MajorAny: "Any",
}

// OrganizerRole represents the organizer's role in a contest
type OrganizerRole uint64

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
