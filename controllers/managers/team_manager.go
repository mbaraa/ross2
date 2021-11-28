package managers

import (
	"math"
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type TeamManager struct {
	teamRepo data.TeamCRUDRepo
	contRepo data.ContestantCRUDRepo
}

func NewTeamManager(teamRepo data.TeamCRUDRepo, contRepo data.ContestantCRUDRepo) *TeamManager {
	return &TeamManager{
		teamRepo: teamRepo,
		contRepo: contRepo,
	}
}

func (t *TeamManager) CreateTeam(team models.Team) error {
	return t.teamRepo.Add(team)
}

// TODO
// check for director permissions

func (t *TeamManager) CreateTeams(teams []models.Team) error {
	return t.teamRepo.AddMany(teams)
}

func (t *TeamManager) AddContestantToTeam(contID, teamID uint) (team models.Team, err error) {
	cont, err := t.contRepo.Get(models.Contestant{ID: contID})
	if err != nil {
		return
	}

	team, err = t.teamRepo.Get(models.Team{ID: teamID}) // just to be safe :), Update could remove all team's members :)
	if err != nil {
		return
	}
	team.Members = append(team.Members, cont)

	cont.TeamID = team.ID
	err = t.contRepo.Update(cont)
	if err != nil {
		return
	}

	return team, t.teamRepo.Update(team)
}

func (t *TeamManager) UpdateTeam(team models.Team, org models.Organizer) error {
	for _, contest := range org.Contests {
		if contest.ID == team.Contests[0].ID { // the team will be stripped from all the other contests :)
			return t.teamRepo.Update(team)
		}
	}
	return nil
}

func (t *TeamManager) UpdateTeams(teams []models.Team, removedContestants []models.Contestant, org models.Organizer) error {
	for _, cont := range removedContestants {
		cont.Team = models.Team{}
		cont.TeamID = 1

		err := t.contRepo.Update(cont)
		if err != nil {
			return err
		}
	}

	for _, team := range teams {
		if _, err := t.GetTeam(team.ID); err != nil {
			team.LeaderId = team.Members[0].ID
			team.Leader = &team.Members[0]

			err = t.CreateTeam(team)
			if err != nil {
				return err
			}
		}

		for _, cont := range team.Members {
			if cont.TeamID != team.ID {
				cont.Team = team
				cont.TeamID = team.ID
				cont.TeamlessContestID = math.MaxInt
				cont.TeamlessedAt = time.Time{}

				err := t.contRepo.Update(cont)
				if err != nil {
					return err
				}
			}
		}

		err := t.teamRepo.Update(team)
		if err != nil {
			return err
		}
	}

	for _, team := range teams {
		err := t.UpdateTeam(team, org)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TeamManager) GetTeam(id uint) (models.Team, error) {
	return t.teamRepo.Get(models.Team{ID: id})
}

func (t *TeamManager) DeleteContestantFromTeam(cont models.Contestant) error {
	team, err := t.GetTeam(cont.TeamID)
	if err != nil {
		return err
	}

	if team.LeaderId == cont.ID {
		if len(team.Members) > 1 {
			for _, member := range team.Members { // change leadership!
				if member.ID != team.LeaderId {
					team.LeaderId = member.ID
					break
				}
			}
			err = t.teamRepo.Update(team)
		} else {
			err = t.DeleteTeam(team)
		}

		if err != nil {
			return err
		}
	}

	for memberIndex := range team.Members { // remove the contestant from the team
		if team.Members[memberIndex].ID == cont.ID {
			team.Members = append(team.Members[:memberIndex], team.Members[memberIndex+1:]...)
			break
		}
	}

	err = t.teamRepo.Update(team)
	if err != nil {
		return err
	}

	cont.TeamID = 1 // add to the no_team team
	return t.contRepo.Update(cont)
}

func (t *TeamManager) DeleteTeam(team models.Team) error {
	team, _ = t.GetTeam(team.ID) // better safe than sorry :\
	for _, member := range team.Members {
		member.TeamID = 1
		_ = t.contRepo.Update(member)
	}
	team.Members = []models.Contestant{}

	return t.teamRepo.Delete(team)
}
