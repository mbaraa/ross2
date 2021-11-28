package partsexport

import (
	"fmt"
	"strings"

	"github.com/mbaraa/ross2/models"
)

func GetParticipants(contest models.Contest) string {
	return getAllCSV(contest)
}

func getAllCSV(contest models.Contest) string {
	sb := new(strings.Builder)
	sb.WriteString("Name, University ID, Role\r\n")
	sb.WriteString(getOrganizersCSV(contest))
	sb.WriteString(getContestantsCSV(contest))

	return sb.String()
}

// getOrganizersCSV returns a string that represents the csv of the given contest's organizers
func getOrganizersCSV(contest models.Contest) string {
	sb := new(strings.Builder)

	for _, org := range contest.Organizers {
		sb.WriteString(fmt.Sprintf("%s, %s, %s\r\n", org.Name, org.Email, "Organizer"))
	}

	return sb.String()
}

// getContestantsCSV returns a string that represents the csv of the given contest's contestants
func getContestantsCSV(contest models.Contest) string {
	var (
		sb    = new(strings.Builder)
		conts = getContestants(contest)
	)

	for _, cont := range conts {
		sb.WriteString(fmt.Sprintf("%s, %s, %s\r\n", cont.Name, cont.UniversityID, "Contestant"))
	}

	return sb.String()
}

func getContestants(contest models.Contest) (conts []models.Contestant) {
	for _, team := range contest.Teams {
		for _, cont := range team.Members {
			conts = append(conts, cont)
		}
	}

	return
}
