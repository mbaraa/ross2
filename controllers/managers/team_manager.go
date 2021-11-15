package managers

import (
	"errors"

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

func (t *TeamManager) CreateTeams(teams []models.Team) error {
	return t.teamRepo.AddMany(teams)
}

func (t *TeamManager) AddContestantToTeam(contID, teamID uint) error {
	cont, err := t.contRepo.Get(models.Contestant{ID: contID})
	if err != nil {
		return err
	}

	team, err := t.teamRepo.Get(models.Team{ID: teamID}) // just to be safe :), Update could remove all team's members :)
	if err != nil {
		return err
	}
	team.Members = append(team.Members, cont)

	cont.TeamID = team.ID
	err = t.contRepo.Update(cont)
	if err != nil {
		return err
	}

	return t.teamRepo.Update(team)
}

func (t *TeamManager) UpdateTeam(team models.Team, org models.Organizer) error {
	for _, contest := range org.Contests {
		if contest.ID == team.Contests[0].ID { // the team will be stripped from all of the other contests :)
			return t.teamRepo.Update(team)
		}
	}
	return errors.New("not authorized to modify team")
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
			team.LeaderId = team.Members[1].ID
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

	cont.TeamID = 0 // add to the no_team team
	return t.contRepo.Update(cont)
}

func (t *TeamManager) DeleteTeam(team models.Team) error {
	team, _ = t.GetTeam(team.ID) // better safe than sorry :\
	team.Members = []models.Contestant{}

	return t.teamRepo.Delete(team)
}
