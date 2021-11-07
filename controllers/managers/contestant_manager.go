package managers

import (
	"time"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type ContestantManager struct {
	contestantRepo data.ContestantCRUDRepo
	sessionManager *SessionManager
	contestRepo    data.ContestCRUDRepo
	teamRepo       data.TeamCRUDRepo
}

func NewContestantManager(contestantRepo data.ContestantCRUDRepo, sessionManager *SessionManager,
	contestRepo data.ContestCRUDRepo, teamRepo data.TeamCRUDRepo) *ContestantManager {
	return &ContestantManager{
		contestantRepo: contestantRepo,
		sessionManager: sessionManager,
		contestRepo:    contestRepo,
		teamRepo:       teamRepo,
	}
}

func (c *ContestantManager) CreateUserSession(email string) error {
	cont, err := c.contestantRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	_, err = c.sessionManager.CreateSession(cont.ID)
	return err
}

func (c *ContestantManager) FinishUser(cont models.Contestant) error {
	return c.contestantRepo.Update(cont)
}

func (c *ContestantManager) RegisterAsTeamless(cont models.Contestant, contest models.Contest) error {
	cont.TeamlessedAt = time.Now()
	cont.TeamlessContestID = contest.ID

	err := c.contestantRepo.Update(cont)
	if err != nil {
		return err
	}

	contest.TeamlessContestants = append(contest.TeamlessContestants, cont)
	return c.contestRepo.Update(contest)
}

func (c *ContestantManager) GetContestant(sessionToken string) (models.Contestant, error) {
	session, err := c.sessionManager.GetSession(sessionToken)
	if err != nil {
		return models.Contestant{}, err
	}

	return c.contestantRepo.Get(models.Contestant{
		ID: session.UserID,
	})
}

func (c *ContestantManager) DeleteUser(cont models.Contestant) error {
	err := c.sessionManager.DeleteAllSessions(cont.ID)
	if err != nil {
		return err
	}

	if cont.TeamID > 0 {
		contTeam, err := c.teamRepo.Get(models.Team{ID: cont.TeamID})
		if err != nil {
			return err
		}

		if contTeam.LeaderId == cont.ID && len(contTeam.Members) > 1 {
			contTeam.LeaderId = contTeam.Members[1].ID
		}

		err = c.teamRepo.Update(contTeam)
		if err != nil {
			return err
		}
	}

	return c.contestantRepo.Delete(cont)
}
