package teamsgen

import "github.com/mbaraa/ross2/models"

type splitTeamless struct {
	males, females, regular []models.Contestant
}

func splitContestantsByGender(conts []models.Contestant) splitTeamless {
	return splitTeamless{
		males:   getMaleContestants(conts),
		females: getFemaleContestants(conts),
		regular: getRegularContestants(conts),
	}
}

func getFemaleContestants(conts []models.Contestant) (females []models.Contestant) {
	for _, cont := range conts {
		if !cont.Gender && !cont.ParticipateWithOtherGender {
			females = append(females, cont)
		}
	}
	return
}

func getMaleContestants(conts []models.Contestant) (males []models.Contestant) {
	for _, cont := range conts {
		if cont.Gender && !cont.ParticipateWithOtherGender {
			males = append(males, cont)
		}
	}
	return
}

func getRegularContestants(conts []models.Contestant) (regulars []models.Contestant) {
	for _, cont := range conts {
		if cont.ParticipateWithOtherGender {
			regulars = append(regulars, cont)
		}
	}
	return
}
