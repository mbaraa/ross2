package enums

// UserType represent user's privileges in Ross
type UserType uint64

// UserType constants
const (
	UserTypeFresh = UserType(1 << iota)
	UserTypeContestant
	UserTypeOrganizer
	UserTypeDirector
	UserTypeAdmin
)

var userTypeText = map[UserType]string{
	UserTypeFresh:      "Fresh",
	UserTypeContestant: "Contestant",
	UserTypeOrganizer:  "Organizer",
	UserTypeDirector:   "Director",
	UserTypeAdmin:      "Admin",
}

// GetTypes returns a string slice that represents the user's types
func (t UserType) GetTypes() (types []string) {
	for typeI := UserTypeFresh; typeI <= UserTypeAdmin; typeI <<= 1 {
		if t&typeI != 0 {
			types = append(types, userTypeText[typeI])
		}
	}

	return
}
