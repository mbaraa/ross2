package controllers

import (
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/context"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/models"
)

// ContestantAPI manages the operations done by a contestant
type ContestantAPI struct {
	contMgr       *helpers.ContestantHelper
	authenticator *auth.HandlerAuthenticator
	endPoints     map[string]http.HandlerFunc
}

// NewContestantAPI returns a new ContestantAPI instance
func NewContestantAPI(contMgr *helpers.ContestantHelper, authenticator *auth.HandlerAuthenticator) *ContestantAPI {
	return (&ContestantAPI{
		contMgr:       contMgr,
		authenticator: authenticator,
	}).initEndPoints()
}

// ServeHTTP is the API's magical port :), since it makes it implement http.Handler
func (c *ContestantAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/contestant"), c.endPoints)
}

func (c *ContestantAPI) initEndPoints() *ContestantAPI {
	c.endPoints = c.authenticator.AuthenticateHandlers(map[string]auth.HandlerFunc{
		"POST /register/": c.handleRegister,
		"POST /profile/":  c.handleGetProfile,

		"POST /create-team/":         c.handleCreateTeam,
		"POST /delete-team/":         c.handleDeleteTeam,
		"POST /req-join-team/":       c.handleRequestJoinTeam,
		"POST /accept-join-request/": c.handleAcceptJoinRequest,
		"POST /reject-join-request/": c.handleRejectJoinRequest,
		"GET /leave-team/":           c.handleLeaveTeam,

		"POST /register-as-teamless/": c.handleRegisterAsTeamless,
		"POST /check-joined-team/":    c.handleCheckJoinedTeam,
		//"POST /invite-teamless/":      nil,

		"GET /get-team/": c.handleGetTeam,
	})
	return c
}

// POST /contestant/register/
// receives models.Contestant in the request body
func (c *ContestantAPI) handleRegister(ctx context.HandlerContext) {
	var cont models.Contestant
	if ctx.ReadJSON(&cont) != nil {
		return
	}

	if c.contMgr.Register(cont) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/profile/
func (c *ContestantAPI) handleGetProfile(ctx context.HandlerContext) {
	var user models.User
	if ctx.ReadJSON(&user) != nil {
		return
	}

	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(cont, 0)
}

// POST /contestant/create-team/
func (c *ContestantAPI) handleCreateTeam(ctx context.HandlerContext) {
	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var team models.Team
	if ctx.ReadJSON(&team) != nil {
		return
	}

	if c.contMgr.CreateTeam(cont, team) != nil {
		return
	}
}

// POST /contestant/delete-team/
func (c *ContestantAPI) handleDeleteTeam(ctx context.HandlerContext) {
	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var team models.Team
	if ctx.ReadJSON(&team) != nil {
		return
	}

	if team.LeaderId != cont.UserID {
		http.Error(ctx.Res, "only team's leader can delete the team", http.StatusUnauthorized)
		return
	}

	if c.contMgr.DeleteTeam(team) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/req-join-team/
func (c *ContestantAPI) handleRequestJoinTeam(ctx context.HandlerContext) {
	var joinRequest models.JoinRequest
	if ctx.ReadJSON(&joinRequest) != nil {
		return
	}

	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if c.contMgr.RequestJoinTeam(joinRequest, cont) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/accept-join-request/
func (c *ContestantAPI) handleAcceptJoinRequest(ctx context.HandlerContext) {
	var notification models.Notification
	if ctx.ReadJSON(&notification) != nil {
		return
	}

	if c.contMgr.AcceptJoinRequest(notification) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/reject-join-request/
func (c *ContestantAPI) handleRejectJoinRequest(ctx context.HandlerContext) {
	var notification models.Notification
	if ctx.ReadJSON(&notification) != nil {
		return
	}

	if c.contMgr.RejectJoinRequest(notification) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /contestant/leave-team/
func (c *ContestantAPI) handleLeaveTeam(ctx context.HandlerContext) {
	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if c.contMgr.LeaveTeam(cont) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/register-as-teamless/
func (c *ContestantAPI) handleRegisterAsTeamless(ctx context.HandlerContext) {
	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if c.contMgr.RegisterAsTeamless(cont, contest) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /contestant/check-joined-team/
func (c *ContestantAPI) handleCheckJoinedTeam(ctx context.HandlerContext) {
	var team models.Team
	if ctx.ReadJSON(&team) != nil {
		return
	}

	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(map[string]interface{}{
		"team_status": c.contMgr.CheckJoinedTeam(cont, team),
	}, 0)
}

// GET /contestant/get-team/
func (c *ContestantAPI) handleGetTeam(ctx context.HandlerContext) {
	cont, err := c.contMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	team, err := c.contMgr.GetTeam(cont)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(team, 0)
}
