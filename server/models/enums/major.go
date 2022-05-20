package enums

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
	MajorAny                     = Major(0)
	FacultyInformationTechnology = MajorSoftwareEngineering | MajorComputerScience | MajorCyberSecurity
)

var majorText = map[Major]string{
	MajorSoftwareEngineering:     "Software Engineering",
	MajorComputerScience:         "Computer Science",
	MajorCyberSecurity:           "Cyber Security",
	FacultyInformationTechnology: "Information Technology",

	MajorAny: "Any",
}

// String returns the major's text
func (m Major) String() string {
	return majorText[m]
}
