package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils"
	"github.com/mbaraa/ross2/utils/teamsgen"
)

type OrganizerAPI struct {
	endPoints map[string]http.HandlerFunc

	orgMgr  *managers.OrganizerManager
	sessMgr *managers.SessionManager
	teamMgr *managers.TeamManager
}

func NewOrganizerAPI(orgMgr *managers.OrganizerManager, sessMgr *managers.SessionManager,
	teamMgr *managers.TeamManager) *OrganizerAPI {
	return (&OrganizerAPI{
		orgMgr:  orgMgr,
		sessMgr: sessMgr,
		teamMgr: teamMgr,
	}).initEndPoints()
}

func (c *OrganizerAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/organizer"), c.endPoints)
}

func (o *OrganizerAPI) initEndPoints() *OrganizerAPI {
	o.endPoints = map[string]http.HandlerFunc{
		// shared organizer/director operations
		"GET /login/":           o.authenticateHandler(o.handleLogin),
		"POST /finish-profile/": o.authenticateHandler(o.handleFinishProfile),
		"DELETE /logout/":       o.authenticateHandler(o.handleLogout),
		"GET /profile/":         o.authenticateHandler(o.handleGetProfile),

		// "GET /get-solved-problems/": nil,
		// "POST /update-solved-problems/": nil,
		// "POST /set-pc-status/": nil,
		// "POST /upload-contest-media/": nil,
		"POST /update-team/": o.authenticateHandler(o.handleUpdateTeam),

		// director operations
		"POST /create-contest/":   o.authenticateHandler(o.handleCreateContest),
		"DELETE /delete-contest/": o.authenticateHandler(o.handleDeleteContest),
		"POST /update-contest/":   o.authenticateHandler(o.handleUpdateContest),
		"POST /add-organizer/":    o.authenticateHandler(o.handleAddOrganizer),
		// "DELETE /delete-organizer/":         nil,
		// "POST /update-organizer/":           nil,
		// "DELETE /delete-contestant/":        nil,
		"POST /auto-generate-teams/": o.authenticateHandler(o.handleAutoGenerateTeams),
		// "POST /man-generate-teams/":         nil,
		"POST /register-generated-teams/": o.authenticateHandler(o.handleRegisterGeneratedTeams),
		// "GET /get-contestants-for-contest/": nil,
		// "GET /get-organizers-for-contest/":   nil,
	}
	return o
}

func (o *OrganizerAPI) authenticateHandler(h HandlerFuncWithSession) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, _ := o.sessMgr.CheckSessionFromRequest(req)
		h(res, req, session)
	}
}

// GET /organizer/login/
func (o *OrganizerAPI) handleLogin(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Authorization", session.ID)
		_ = json.NewEncoder(res).Encode(org)
	} else {
		http.Redirect(res, req,
			config.GetInstance().MachineAddress+"/gauth/org-login/", http.StatusPermanentRedirect)
	}
}

// POST /organizer/finish-profile/
func (o *OrganizerAPI) handleFinishProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		var orgData models.Organizer
		err := json.NewDecoder(req.Body).Decode(&orgData)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.orgMgr.UpdateProfile(orgData)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// DELETE /organizer/logout/
func (o *OrganizerAPI) handleLogout(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		err := o.sessMgr.DeleteSession(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GET /organizer/profile/
func (o *OrganizerAPI) handleGetProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(res).Encode(org)
	}
}

// POST /organizer/update-team/
func (o *OrganizerAPI) handleUpdateTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || (org.Roles&models.RoleDirector == 0 || org.Roles&models.RoleCoreOrganizer == 0) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var team models.Team
		err = json.NewDecoder(req.Body).Decode(&team)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.teamMgr.UpdateTeam(team, org)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

// POST /organizer/create-contest/
func (o *OrganizerAPI) handleCreateContest(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var contest models.Contest
		err = json.NewDecoder(req.Body).Decode(&contest)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.orgMgr.CreateContest(contest)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

// DELETE /organizer/delete-contest/
func (o *OrganizerAPI) handleDeleteContest(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var contest models.Contest
		err = json.NewDecoder(req.Body).Decode(&contest)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.orgMgr.DeleteContest(contest)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// POST /organizer/update-contest/
func (o *OrganizerAPI) handleUpdateContest(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var contest models.Contest
		err = json.NewDecoder(req.Body).Decode(&contest)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.orgMgr.UpdateContest(contest)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// POST /organizer/add-organier/
func (o *OrganizerAPI) handleAddOrganizer(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var newOrg models.Organizer
		err = json.NewDecoder(req.Body).Decode(&newOrg)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.orgMgr.AddOrganizer(newOrg)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// POST /organizer/auto-generate-teams/
func (o *OrganizerAPI) handleAutoGenerateTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var contest models.Contest
		err = json.NewDecoder(req.Body).Decode(&contest)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		teams := teamsgen.GenerateTeams(contest, utils.NewHardCodeNames())
		_ = json.NewEncoder(res).Encode(teams)
	}
}

// POST /organizer/register-generated-teams/
func (o *OrganizerAPI) handleRegisterGeneratedTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
	if session.ID != "" {
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil || org.Roles&models.RoleDirector == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		var teams []models.Team
		err = json.NewDecoder(req.Body).Decode(&teams)
		_ = req.Body.Close()
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.teamMgr.CreateTeams(teams)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
