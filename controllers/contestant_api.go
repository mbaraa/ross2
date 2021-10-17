package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/config"
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
		"GET /login/":     c.authenticateHandler(c.handleLogin),
		"POST /signup/":   c.handleSignup,
		"DELETE /logout/": c.authenticateHandler(c.handleLogout),
		"DELETE /delete/": c.authenticateHandler(c.handleDelete),
		"GET /profile/":   c.authenticateHandler(c.handleGetProfile),

		"POST /create-team/":         c.authenticateHandler(c.handleCreateTeam),
		"DELETE /delete-team/":       c.authenticateHandler(c.handleDeleteTeam),
		"POST /req-join-team/":       c.authenticateHandler(c.handleRequestJoinTeam),
		"POST /accept-join-request/": c.authenticateHandler(c.handleAcceptJoinRequest),
		"POST /reject-join-request/": c.authenticateHandler(c.handleRejectJoinRequest),
		"DELETE /leave-team/":        c.authenticateHandler(c.handleLeaveTeam),

		"POST /register-as-teamless/": c.authenticateHandler(c.handleRegisterAsTeamless),
		//"POST /invite-teamless":       nil,
	}
	return c
}

type HandlerFuncWithSession func(http.ResponseWriter, *http.Request, models.Session)

func (c *ContestantAPI) authenticateHandler(h HandlerFuncWithSession) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, _ := c.sessManager.CheckSessionFromRequest(req)
		h(res, req, session)
	}
}

// GET /contestant/login/
func (c *ContestantAPI) handleLogin(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		cont, err := c.contManager.GetContestant(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Authorization", session.ID)
		_ = json.NewEncoder(res).Encode(cont)
	} else {
		http.Redirect(res, req,
			config.GetInstance().MachineAddress+"/gauth/cont-login/", http.StatusPermanentRedirect)
	}
}

// POST /contestant/signup
func (c *ContestantAPI) handleSignup(res http.ResponseWriter, req *http.Request) {
	var contestant models.Contestant
	err := json.NewDecoder(req.Body).Decode(&contestant)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.contManager.CreateUser(contestant)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	// login after account creation is done
	http.Redirect(res, req, config.GetInstance().MachineAddress+"/contestant/login/", http.StatusPermanentRedirect)
}

// DELETE /contestant/logout/
func (c *ContestantAPI) handleLogout(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		err := c.sessManager.DeleteSession(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// DELETE /contestant/delete/
func (c *ContestantAPI) handleDelete(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// GET /contestant/profile/
func (c *ContestantAPI) handleGetProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		cont, err := c.contManager.GetContestant(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(res).Encode(cont)
	}
}

// POST /contestant/create-team/
func (c *ContestantAPI) handleCreateTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
		err = c.teamManager.CreateTeam(team)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

// DELETE /contestant/delete-team/
func (c *ContestantAPI) handleDeleteTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// POST /contestant/req-join-team/
func (c *ContestantAPI) handleRequestJoinTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// POST /contestant/accept-join-request/
func (c *ContestantAPI) handleAcceptJoinRequest(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// POST /contestant/reject-join-request/
func (c *ContestantAPI) handleRejectJoinRequest(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// DELETE /contestant/leave-team/
func (c *ContestantAPI) handleLeaveTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}

// POST /contestant/register-as-teamless/
func (c *ContestantAPI) handleRegisterAsTeamless(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
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
}
