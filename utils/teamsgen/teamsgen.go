package teamsgen

import (
	"math/rand"
	"sort"
	"time"

	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils/namesgetter"
)

var (
	leftMembers = map[uint]bool{} // {id, status}
)

func markMembers(teamless []models.Contestant) {
	for _, cont := range teamless {
		leftMembers[cont.User.ID] = true
	}
}

func getLeftMembers(teamless []models.Contestant) (left []models.Contestant) {
	for _, cont := range teamless {
		if !leftMembers[cont.User.ID] {
			left = append(left, cont)
		}
	}
	return
}

// GenerateTeams generates teams for the sad teamless contestants of the given contest
// also it uses the NamesGetter interface to assign names for the created teams :)
func GenerateTeams(contest models.Contest, names namesgetter.NamesGetter) ([]models.Team, []models.Contestant) {
	sConts := splitContestantsByGender(contest.TeamlessContestants)

	markMembers(contest.TeamlessContestants)

	return combineTeams(
			generateTeams(sConts.males, contest, names),
			generateTeams(sConts.females, contest, names),
			generateTeams(sConts.regular, contest, names)),
		getLeftMembers(contest.TeamlessContestants)
}

func combineTeams(teamsGroups ...[]models.Team) []models.Team {
	combined := append(teamsGroups[0], teamsGroups[1]...)
	return append(combined, teamsGroups[2]...)
}

// generateTeams this is where the fun begins...
func generateTeams(teamless []models.Contestant, contest models.Contest, names namesgetter.NamesGetter) []models.Team {
	minMembers := contest.ParticipationConditions.MinTeamMembers
	maxMembers := contest.ParticipationConditions.MaxTeamMembers

	if !isCompletelyFillable(uint(len(teamless)), minMembers, maxMembers) {
		teamless = deleteLastAddedTeamless(teamless, minMembers)
	}

	return stripContests(
		deleteEmptySlots(
			finalizeTeams(
				fillTeams(teamless,
					generateEmptyTeams(1+(len(teamless)/int(minMembers)), names),
					minMembers, maxMembers),
				contest,
			)))
}

// stripContests removes generated teams' contests to reduce size,
// but it might fuck things up, I'll keep an eye on it 😉
func stripContests(teams []models.Team) []models.Team {
	for i := range teams {
		teams[i].Contests = nil
	}
	return teams
}

// isCompletelyFillable reports whether all the teamless will be filled in teams
func isCompletelyFillable(numConts, min, max uint) bool {
	return numConts%max == 0 ||
		(min == max && numConts%max == 0) ||
		(min < max && numConts%min == 0)
}

// removeLastAddedTeamless removes the last chronologically registered as teamless contestants :)
// fair ain't it?
func deleteLastAddedTeamless(teamless []models.Contestant, min uint) []models.Contestant {
	var (
		conts       = len(teamless)
		leftConts   = conts % int(min)
		sortedConts = models.ContestantSortable(teamless)
	)
	sort.Sort(sortedConts)

	for i := conts - leftConts; i < conts; i++ {
		leftMembers[sortedConts[i].User.ID] = false
	}

	return sortedConts[:conts-leftConts]
}

// deleteEmptySlots deletes empty teams and empty contestants
// also also if you can do an optimized version of it I won't let you down :)
func deleteEmptySlots(teams []models.Team) []models.Team {
	cleanTeams := make([]models.Team, 0)

	for _, team := range teams {
		if team.Members != nil {
			cleanTeams = append(cleanTeams, team)
		}
	}

	return cleanTeams
}

// finalizeTeams adds the generated teams to their contest and sets the first member
// as team leader
func finalizeTeams(teams []models.Team, contest models.Contest) []models.Team {
	for teamIndex := range teams {
		teams[teamIndex].Contests = append(teams[teamIndex].Contests, contest)

		if teams[teamIndex].Members != nil {
			teams[teamIndex].LeaderId = teams[teamIndex].Members[0].User.ID
		}
	}

	return teams
}

// fillTeams adds teamless members to the generated teams
func fillTeams(teamless []models.Contestant, teams []models.Team, minMembers, maxMembers uint) []models.Team {
	teamIndex, teamMembers := uint(0), uint(0)
	teamless = shuffleMembers(teamless)

	for _, cont := range teamless {
		teams[teamIndex].Members = append(teams[teamIndex].Members, cont)
		teamMembers++
		if teamMembers == maxMembers {
			teamIndex++
			teamMembers = 0
		}
	}

	return fillLastTeam(teams, minMembers, maxMembers, teamIndex)
}

// shuffleMembers randomizes the teamless contestants slice so the outcome
// teams' members is not expected by the contestants it's a bit mean I know
func shuffleMembers(teamless []models.Contestant) []models.Contestant {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(teamless), func(i, j int) {
		teamless[i], teamless[j] = teamless[j], teamless[i]
	})

	return teamless
}

// fillLastTeam adds members to the last team
func fillLastTeam(teams []models.Team, minMembers, maxMembers, lastTeamIndex uint) []models.Team {
	for teamIndex := uint(0); len(teams[lastTeamIndex].Members) < int(minMembers) &&
		len(teams[lastTeamIndex].Members) > 0 && teamIndex < lastTeamIndex; {

		teams[lastTeamIndex].Members =
			append(teams[lastTeamIndex].Members, teams[teamIndex].Members[minMembers-1:]...)

		teams[teamIndex].Members = teams[teamIndex].Members[:minMembers]
	}

	return teams
}

// generateEmptyTeams generates teams w/ just names
func generateEmptyTeams(numTeams int, names namesgetter.NamesGetter) []models.Team {
	teams := make([]models.Team, numTeams)
	for teamIndex := range teams {
		teams[teamIndex] = models.Team{
			Name: names.GetName(),
		}
	}

	return teams
}
