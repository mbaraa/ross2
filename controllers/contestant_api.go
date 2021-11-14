package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/models"
)

type ContestantAPI struct {
	endPoints map[string]http.HandlerFunc

	contManager    *managers.ContestantManager
	sessManager    *managers.SessionManager
	teamManager    *managers.TeamManager
	joinReqManager *managers.JoinRequestManager
}

func NewContestantAPI(contManager *managers.ContestantManager, sessManager *managers.SessionManager,
	teamManager *managers.TeamManager, joinReqManager *managers.JoinRequestManager) *ContestantAPI {
	return (&ContestantAPI{
		contManager:    contManager,
		sessManager:    sessManager,
		teamManager:    teamManager,
		joinReqManager: joinReqManager,
	}).initEndPoints()
}

func (c *ContestantAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/contestant"), c.endPoints)
}

func (c *ContestantAPI) initEndPoints() *ContestantAPI {
	c.endPoints = map[string]http.HandlerFunc{
		"GET /login/":   c.handleLogin,
		"POST /signup/": c.handleSignup,
		"GET /logout/":  c.authenticateHandler(c.handleLogout),
		"GET /delete/":  c.authenticateHandler(c.handleDelete),
		"GET /profile/": c.authenticateHandler(c.handleGetProfile),

		"POST /create-team/":         c.authenticateHandler(c.handleCreateTeam),
		"POST /delete-team/":         c.authenticateHandler(c.handleDeleteTeam),
		"POST /req-join-team/":       c.authenticateHandler(c.handleRequestJoinTeam),
		"POST /accept-join-request/": c.authenticateHandler(c.handleAcceptJoinRequest),
		"POST /reject-join-request/": c.authenticateHandler(c.handleRejectJoinRequest),
		"GET /leave-team/":           c.authenticateHandler(c.handleLeaveTeam),

		"POST /register-as-teamless/": c.authenticateHandler(c.handleRegisterAsTeamless),
		"POST /check-joined-team/":    c.authenticateHandler(c.handleCheckJoinedTeam),
		//"POST /invite-teamless/":      nil,

		"GET /get-team/": c.authenticateHandler(c.handleGetTeam),
	}
	return c
}

type HandlerFuncWithSession func(http.ResponseWriter, *http.Request, models.Session)

func (c *ContestantAPI) authenticateHandler(h HandlerFuncWithSession) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, ok := c.sessManager.CheckSessionFromRequest(req)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(res, req, session)
	}
}

// GET /contestant/login/
func (c *ContestantAPI) handleLogin(res http.ResponseWriter, req *http.Request) {
	if session, ok := c.sessManager.CheckSessionFromRequest(req); ok { // I kinda fucked up a bit 😅
		cont, err := c.contManager.GetContestant(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(res).Encode(cont)
	} else {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// POST /contestant/signup
func (c *ContestantAPI) handleSignup(res http.ResponseWriter, req *http.Request) {
	contestant := models.Contestant{}
	err := json.NewDecoder(req.Body).Decode(&contestant)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.contManager.FinishUser(contestant)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/logout/
func (c *ContestantAPI) handleLogout(res http.ResponseWriter, req *http.Request, session models.Session) {
	err := c.sessManager.DeleteSession(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/delete/
func (c *ContestantAPI) handleDelete(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.contManager.DeleteUser(cont)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/profile/
func (c *ContestantAPI) handleGetProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(cont)
}

// POST /contestant/create-team/
func (c *ContestantAPI) handleCreateTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	var team models.Team
	err := json.NewDecoder(req.Body).Decode(&team)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	// setting the team leader here is far easier than doing it from the frontend
	leader, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	team.LeaderId = leader.ID
	team.Members = append(team.Members, leader)

	err = c.teamManager.CreateTeam(team)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// POST /contestant/delete-team/
func (c *ContestantAPI) handleDeleteTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	var team models.Team
	err := json.NewDecoder(req.Body).Decode(&team)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil || team.LeaderId != cont.ID {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.teamManager.DeleteTeam(team)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// POST /contestant/req-join-team/
func (c *ContestantAPI) handleRequestJoinTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	var jr models.JoinRequest
	err := json.NewDecoder(req.Body).Decode(&jr)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	jr.RequesterID = cont.ID

	err = c.joinReqManager.CreateRequest(jr, cont)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// POST /contestant/accept-join-request/
func (c *ContestantAPI) handleAcceptJoinRequest(res http.ResponseWriter, req *http.Request, session models.Session) {
	var noti models.Notification // yes I'm gonna squeeze them from here :)
	err := json.NewDecoder(req.Body).Decode(&noti)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.joinReqManager.AcceptJoinRequest(noti)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/reject-join-request/
func (c *ContestantAPI) handleRejectJoinRequest(res http.ResponseWriter, req *http.Request, session models.Session) {
	var noti models.Notification
	err := json.NewDecoder(req.Body).Decode(&noti)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.joinReqManager.RejectJoinRequest(noti)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/leave-team/
func (c *ContestantAPI) handleLeaveTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.teamManager.DeleteContestantFromTeam(cont)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// POST /contestant/register-as-teamless/
func (c *ContestantAPI) handleRegisterAsTeamless(res http.ResponseWriter, req *http.Request, session models.Session) {
	var contest models.Contest
	err := json.NewDecoder(req.Body).Decode(&contest)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.contManager.RegisterAsTeamless(cont, contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/check-joined-team/
func (c *ContestantAPI) handleCheckJoinedTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	var team models.Team
	err := json.NewDecoder(req.Body).Decode(&team)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	inTeam := cont.TeamID == team.ID ||
		c.joinReqManager.CheckContestantTeamRequests(cont, team)

	_, _ = res.Write([]byte(fmt.Sprintf(`{"team_status" : %v}`, inTeam)))
}

// GET /contestant/get-team/
func (c *ContestantAPI) handleGetTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contManager.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	team, err := c.teamManager.GetTeam(cont.TeamID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(team)
}
