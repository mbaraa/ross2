package helpers

import (
	"math"
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

// TeamHelper manages teams and stuff
type TeamHelper struct {
	repo     data.CRUDRepo[models.Team]
	contRepo data.CRUDRepo[models.Contestant]
	rtRepo   data.RegisterTeamCRUDRepo
}

// NewTeamHelper returns a new TeamHelper instance
func NewTeamHelper(teamRepo data.CRUDRepo[models.Team], contRepo data.CRUDRepo[models.Contestant],
	rtRepo data.RegisterTeamCRUDRepo) *TeamHelper {
	return &TeamHelper{
		repo:     teamRepo,
		contRepo: contRepo,
		rtRepo:   rtRepo,
	}
}

// CreateTeam creates a team and adds the given contestant to it as its leader
func (t *TeamHelper) CreateTeam(contestant models.Contestant, team *models.Team) error {
	// set team's required attributes to be lead by the given contestant
	team.Leader = &contestant
	team.LeaderId = contestant.User.ID

	err := t.repo.Add(team)
	if err != nil {
		return err
	}

	// set contestant's required attributes to join the team
	contestant.TeamID = team.ID
	contestant.Team = *team

	return t.contRepo.Update(&contestant)
}

// CreateTeams creates the given teams :)
func (t *TeamHelper) CreateTeams(teams []*models.Team) error {
	return t.repo.AddMany(teams)
}

// AddContestantToTeam adds the given contestant to the given team
func (t *TeamHelper) AddContestantToTeam(contID, teamID uint) (team models.Team, err error) {
	cont, err := t.contRepo.Get(contID)
	if err != nil {
		return
	}

	team, err = t.repo.Get(teamID) // just to be safe :), Update could remove all team's members :)
	if err != nil {
		return
	}
	team.Members = append(team.Members, cont)

	cont.TeamID = team.ID
	cont.Team = team

	// if cont.TeamlessContestID != 0 {
	cont.TeamlessContestID = math.MaxInt64
	cont.TeamlessedAt = cont.CreatedAt
	// }

	err = t.contRepo.Update(&cont)
	if err != nil {
		return
	}

	return team, t.repo.Update(&team)
}

// UpdateTeam updates the given team after checking that the given organizer is a director on one of the contest that the team is in
func (t *TeamHelper) UpdateTeam(team models.Team, org models.Organizer) error {
	for _, contest := range org.Contests {
		if contest.ID == team.Contests[0].ID { // the team will be stripped from all the other contests :)
			return t.repo.Update(&team)
		}
	}
	return nil
}

func (t *TeamHelper) CreateUpdateTeams(teams []models.Team, removedContestants []models.Contestant, contest models.Contest, org models.Organizer) error {
	for _, cont := range removedContestants {
		cont.Team = models.Team{}
		cont.TeamID = 1
		cont.TeamlessContestID = contest.ID
		cont.TeamlessedAt = time.Now()

		err := t.contRepo.Update(&cont)
		if err != nil {
			return err
		}
	}

	for _, team := range teams {
		// create
		if _, err := t.GetTeam(team.ID); err != nil && (team.Members != nil && len(team.Members) > 0) {
			team.LeaderId = team.Members[0].User.ID
			team.Leader = &team.Members[0]
			team.Contests = []models.Contest{contest}

			err = t.CreateTeam(team.Members[0], &team)
			if err != nil {
				return err
			}
		}

		// delete
		if team.Members == nil || len(team.Members) == 0 {
			err := t.DeleteTeam(models.Contestant{}, team)
			if err != nil {
				return err
			}
			continue
		}

		// update
		for _, cont := range team.Members {
			if cont.TeamID != team.ID {
				cont.Team = team
				cont.TeamID = team.ID
				cont.TeamlessContestID = math.MaxInt
				cont.TeamlessedAt = time.Time{}

				err := t.contRepo.Update(&cont)
				if err != nil {
					return err
				}
			}
		}

		team.Contests = []models.Contest{contest}
		err := t.repo.Update(&team)
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

// GetTeam returns a team using the given id
func (t *TeamHelper) GetTeam(id uint) (models.Team, error) {
	return t.repo.Get(id)
}

// GetTeamByJoinID returns a team using the given join id
func (t *TeamHelper) GetTeamByJoinID(joinID string) (models.Team, error) {
	teams, err := t.repo.GetByConds("join_id = ?", joinID)
	return teams[0], err
}

// LeaveTeam removes the given contestant from their team in a super safe way
func (t *TeamHelper) LeaveTeam(cont models.Contestant, team models.Team) error {
	var err error
	if team.LeaderId == cont.User.ID {
		if len(team.Members) > 1 {
			for _, member := range team.Members { // change leadership!
				if member.User.ID != team.LeaderId {
					team.LeaderId = member.User.ID
					break
				}
			}
			err = t.repo.Update(&team)
		} else {
			err = t.DeleteTeam(cont, team)
		}
		if err != nil {
			return err
		}

		err = t.repo.
			GetDB().
			Model(&cont).
			Association("Teams").
			Delete(&team)

		if err != nil {
			return err
		}
	}

	for memberIndex := range team.Members { // remove the contestant from the team
		if team.Members[memberIndex].User.ID == cont.User.ID {
			team.Members = append(team.Members[:memberIndex], team.Members[memberIndex+1:]...)
			break
		}
	}

	return t.repo.Update(&team)
}

// DeleteTeam kicks every member out of it and deletes it
func (t *TeamHelper) DeleteTeam(cont models.Contestant, team models.Team) error {
	if cont.ID == 0 {
		team, _ = t.GetTeam(team.ID) // better safe than sorry :\
		for _, member := range team.Members {
			member.TeamID = 1
		}

		return t.repo.Delete(team)
	}
	return t.repo.
		GetDB().
		Model(&cont).
		Association("Teams").
		Delete(&team)
}
