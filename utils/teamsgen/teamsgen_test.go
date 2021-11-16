package teamsgen

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils"
)

func TestBigFatMonstrousCases(t *testing.T) {
	for i := 0; i < 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		var (
			min   = uint(rand.Intn(1000) + 1)
			max   = uint(rand.Intn(1000) + int(min))
			conts = uint(rand.Intn(10000) + 1)
		)

		teams := generalTest(min, max, conts)
		//for _, t := range teams {
		//	fmt.Println(t.Name)
		//}
		if ok, team := checkTeamsMembers(teams, min); !ok {
			t.Errorf("team %s doesn't have enough members!\n", team.Name)
			t.Errorf("team %s has %d members, but wants %d members\n", team.Name, len(team.Members), min)
			t.Errorf("test vars: min %d, max: %d, conts: %d\n\n", min, max, conts)
		}
	}
}

func TestManual(t *testing.T) {
	var (
		min   = uint(rand.Intn(1000) + 1)
		max   = uint(rand.Intn(1000) + int(min))
		conts = uint(rand.Intn(10000) + 1)
	)
	teams := generalTest(min, max, conts)

	if ok, team := checkTeamsMembers(teams, min); !ok {
		t.Errorf("team %s doesn't have enough members!\n", team.Name)
		t.Errorf("team %s has %d members, but wants %d members\n", team.Name, len(team.Members), min)
		t.Errorf("test vars: min %d, max: %d, conts: %d\n\n", min, max, conts)
	}

	logTeams(teams, t)
}

func TestFitGenerateTeams(t *testing.T) {
	minMembers := uint(3)
	teams := generalTest(minMembers, 3, 9)

	if ok, team := checkTeamsMembers(teams, minMembers); !ok {
		t.Errorf("team %s doesn't have enough members!\n", team.Name)
	}

	logTeams(teams, t)
}

func TestFitWithMinGenerateTeams(t *testing.T) {
	minMembers := uint(2)
	teams := generalTest(minMembers, 3, 10)

	if ok, team := checkTeamsMembers(teams, minMembers); !ok {
		t.Errorf("team %s doesn't have enough members!\n", team.Name)
	}

	logTeams(teams, t)

}

func TestFitMinEqMaxGenerateTeams(t *testing.T) {
	minMembers := uint(3)
	teams := generalTest(minMembers, 3, 9)

	if ok, team := checkTeamsMembers(teams, minMembers); !ok {
		t.Errorf("team %s doesn't have enough members!\n", team.Name)
	}

	logTeams(teams, t)
}

func TestNoFitMinEqMaxGenerateTeams(t *testing.T) {
	minMembers := uint(3)
	teams := generalTest(minMembers, 3, 10)

	if ok, team := checkTeamsMembers(teams, minMembers); !ok {
		t.Errorf("team %s doesn't have enough members!\n", team.Name)
	}

	logTeams(teams, t)
}

/////////////////
// TEST UTILS //
////////////////

func generalTest(minMembers, maxMembers, numConts uint) []models.Team {
	teams, _ := GenerateTeams(models.Contest{
		Name:        "Potato Peeling Contest",
		StartsAt2:   time.Now().Add(time.Hour * 300),
		Duration:    time.Minute * 120,
		Location:    "Online",
		Description: "Peel potatoes as fast as you can!",
		ParticipationConditions: models.ParticipationConditions{
			MinTeamMembers: minMembers,
			MaxTeamMembers: maxMembers,
			Majors:         models.MajorAny,
		},
		TeamlessContestants: createRandomContestants(numConts),
	}, utils.NewHardCodeNames())

	return teams
}

func logTeams(teams []models.Team, t *testing.T) {

	for _, team := range teams {
		t.Logf("team %s member's count: %d", team.Name, len(team.Members))
		t.Log("members:")
		for _, cont := range team.Members {
			t.Logf("name: %s, uni_id: %s", cont.Name, cont.UniversityID)
		}
		t.Log()
	}
}

func checkTeamsMembers(teams []models.Team, minMembers uint) (bool, models.Team) {
	for _, team := range teams {
		members := uint(0)
		for _, cont := range team.Members {
			if cont.Name != "" {
				members++
			}
		}
		if members < minMembers {
			return false, team
		}
	}
	return true, models.Team{}
}

func createRandomContestants(numConts uint) []models.Contestant {
	conts := make([]models.Contestant, numConts)

	for contI := range conts {
		rand.Seed(time.Now().UnixNano())
		gender := false
		if rand.Intn(2) == 1 {
			gender = true
		} else {
			gender = false
		}

		partWithOther := false
		if rand.Intn(2) == 1 {
			partWithOther = true
		} else {
			partWithOther = false
		}

		conts[contI] = models.Contestant{
			Name:                       uuid.New().String(),
			UniversityID:               fmt.Sprint(rand.Intn(9000000) + 1000000),
			Major:                      models.MajorAny,
			TeamlessedAt:               time.Now(),
			Gender:                     gender,
			ParticipateWithOtherGender: partWithOther,
		}
	}

	return conts
}
