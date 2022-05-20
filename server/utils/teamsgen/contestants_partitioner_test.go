package teamsgen

import (
	"testing"

	"github.com/mbaraa/ross2/models"
)

func TestSplitContestantsByGender(t *testing.T) {
	sConts := splitContestantsByGender(
		createRandomContestants(1000000))

	if !checkSameGender(sConts.females) {
		t.Errorf("failed names: %v\n", sConts.females)
	}

	if !checkSameGender(sConts.males) {

		t.Errorf("failed names: %v\n", sConts.males)
	}

	if !checkSameGender(sConts.regular) {
		t.Errorf("failed names: %v\n", sConts.regular)
	}

}

func checkSameGender(conts []models.Contestant) bool {
	gender := conts[0].Gender
	partWithOther := conts[0].ParticipateWithOtherGender

	for _, cont := range conts {
		if cont.Gender != gender && cont.ParticipateWithOtherGender != partWithOther {
			return false
		}
	}

	return true
}
