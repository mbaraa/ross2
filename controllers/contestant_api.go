package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/models"
)

type ContestantAPIBuilder struct {
	contMgr    *managers.ContestantManager
	sessMgr    *managers.SessionManager
	teamMgr    *managers.TeamManager
	joinReqMgr *managers.JoinRequestManager
}

func NewContestantAPIBuilder() *ContestantAPIBuilder {
	return new(ContestantAPIBuilder)
}

func (b *ContestantAPIBuilder) ContestantMgr(c *managers.ContestantManager) *ContestantAPIBuilder {
	b.contMgr = c
	return b
}

func (b *ContestantAPIBuilder) SessionMgr(s *managers.SessionManager) *ContestantAPIBuilder {
	b.sessMgr = s
	return b
}

func (b *ContestantAPIBuilder) TeamMgr(t *managers.TeamManager) *ContestantAPIBuilder {
	b.teamMgr = t
	return b
}

func (b *ContestantAPIBuilder) JoinReqMgr(j *managers.JoinRequestManager) *ContestantAPIBuilder {
	b.joinReqMgr = j
	return b
}

func (b *ContestantAPIBuilder) verify() bool {
	if b.contMgr == nil {
		fmt.Println("Contestant API Builder: missing contestant manager!")
	}
	if b.sessMgr == nil {
		fmt.Println("Contestant API Builder: missing session manager!")
	}
	if b.teamMgr == nil {
		fmt.Println("Contestant API Builder: missing team manager!")
	}
	if b.joinReqMgr == nil {
		fmt.Println("Contestant API Builder: missing join request manager!")
	}

	return b.contMgr != nil && b.sessMgr != nil &&
		b.teamMgr != nil && b.joinReqMgr != nil
}

func (b *ContestantAPIBuilder) GetContestantAPI() *ContestantAPI {
	if !b.verify() {
		return nil
	}
	return NewContestantAPI(b)
}

type ContestantAPI struct {
	endPoints map[string]http.HandlerFunc

	contMgr    *managers.ContestantManager
	sessMgr    *managers.SessionManager
	teamMgr    *managers.TeamManager
	joinReqMgr *managers.JoinRequestManager
}

func NewContestantAPI(b *ContestantAPIBuilder) *ContestantAPI {
	return (&ContestantAPI{
		contMgr:    b.contMgr,
		sessMgr:    b.sessMgr,
		teamMgr:    b.teamMgr,
		joinReqMgr: b.joinReqMgr,
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
		session, ok := c.sessMgr.CheckSessionFromRequest(req)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(res, req, session)
	}
}

// GET /contestant/login/
func (c *ContestantAPI) handleLogin(res http.ResponseWriter, req *http.Request) {
	if session, ok := c.sessMgr.CheckSessionFromRequest(req); ok { // I kinda fucked up a bit 😅
		cont, err := c.contMgr.GetContestant(session.ID)
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

	err = c.contMgr.FinishUser(contestant)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/logout/
func (c *ContestantAPI) handleLogout(res http.ResponseWriter, req *http.Request, session models.Session) {
	err := c.sessMgr.DeleteSession(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/delete/
func (c *ContestantAPI) handleDelete(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.contMgr.DeleteUser(cont)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/profile/
func (c *ContestantAPI) handleGetProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contMgr.GetContestant(session.ID)
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
	leader, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	team.LeaderId = leader.ID
	team.Members = append(team.Members, leader)

	err = c.teamMgr.CreateTeam(team)
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

	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil || team.LeaderId != cont.ID {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.teamMgr.DeleteTeam(team)
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

	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	jr.RequesterID = cont.ID
	jr.Requester = cont

	err = c.joinReqMgr.CreateRequest(jr, cont)
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

	err = c.joinReqMgr.AcceptJoinRequest(noti)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte(fmt.Sprintf(`{"err": "%s"}`, err.Error())))
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

	err = c.joinReqMgr.RejectJoinRequest(noti)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/leave-team/
func (c *ContestantAPI) handleLeaveTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.teamMgr.DeleteContestantFromTeam(cont)
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

	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = c.contMgr.RegisterAsTeamless(cont, contest)
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

	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	inTeam := cont.TeamID > 1 || cont.TeamID == team.ID ||
		c.joinReqMgr.CheckContestantTeamRequests(cont, team)

	_, _ = res.Write([]byte(fmt.Sprintf(`{"team_status" : %v}`, inTeam)))
}

// GET /contestant/get-team/
func (c *ContestantAPI) handleGetTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	cont, err := c.contMgr.GetContestant(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	team, err := c.teamMgr.GetTeam(cont.TeamID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(team)
}
