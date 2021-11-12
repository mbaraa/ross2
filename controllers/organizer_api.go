package controllers

import (
	"encoding/json"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils"
	"github.com/mbaraa/ross2/utils/teamsgen"
	"io"
	"net/http"
	"os"
	"strings"
)

type OrganizerAPI struct {
	endPoints map[string]http.HandlerFunc

	orgMgr  *managers.OrganizerManager
	sessMgr *managers.SessionManager
	teamMgr *managers.TeamManager
	contMgr *managers.ContestantManager
}

func NewOrganizerAPI(orgMgr *managers.OrganizerManager, sessMgr *managers.SessionManager,
	teamMgr *managers.TeamManager, contMgr *managers.ContestantManager) *OrganizerAPI {
	return (&OrganizerAPI{
		orgMgr:  orgMgr,
		sessMgr: sessMgr,
		teamMgr: teamMgr,
		contMgr: contMgr,
	}).initEndPoints()
}

func (o *OrganizerAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/organizer"), o.endPoints)
}

func (o *OrganizerAPI) initEndPoints() *OrganizerAPI {
	o.endPoints = map[string]http.HandlerFunc{
		// shared organizer/director operations
		"GET /login/":           o.handleLogin,
		"POST /finish-profile/": o.handleFinishProfile,
		"GET /logout/":          o.authenticateHandler(o.handleLogout),
		"GET /profile/":         o.authenticateHandler(o.handleGetProfile),

		// "GET /get-solved-problems/": nil,
		// "POST /update-solved-problems/": nil,
		// "POST /set-pc-status/": nil,
		// "POST /upload-contest-media/": nil,
		"POST /update-team/": o.authenticateHandler(o.handleUpdateTeam),

		// director operations
		"POST /create-contest/":           o.authenticateHandler(o.handleCreateContest),
		"POST /delete-contest/":           o.authenticateHandler(o.handleDeleteContest),
		"POST /update-contest/":           o.authenticateHandler(o.handleUpdateContest),
		"POST /upload-contest-logo-file/": o.authenticateHandler(o.handleUploadContestLogoFile),
		"POST /add-organizer/":            o.authenticateHandler(o.handleAddOrganizer),
		"GET /delete-organizer/":          o.authenticateHandler(o.handleDeleteOrganizer),
		"POST /update-organizer/":         o.authenticateHandler(o.handleUpdateOrganizer),
		"POST /delete-contestant/":        o.authenticateHandler(o.handleDeleteContestant),
		"GET /get-sub-organizers/":        o.authenticateHandler(o.handleGetSubOrganizers),
		"POST /auto-generate-teams/":      o.authenticateHandler(o.handleAutoGenerateTeams),
		// "POST /man-generate-teams/":         nil,
		"POST /register-generated-teams/": o.authenticateHandler(o.handleRegisterGeneratedTeams),
		// "GET /get-contestants-for-contest/": nil,
		// "GET /get-organizers-for-contest/":  nil,

		"GET /get-contests/": o.authenticateHandler(o.handleGetContests),
	}
	return o
}

func (o *OrganizerAPI) authenticateHandler(h HandlerFuncWithSession) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, ok := o.sessMgr.CheckSessionFromRequest(req)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(res, req, session)
	}
}

// GET /organizer/login/
func (o *OrganizerAPI) handleLogin(res http.ResponseWriter, req *http.Request) {
	if session, ok := o.sessMgr.CheckSessionFromRequest(req); ok { // more fuck-ups :)
		org, err := o.orgMgr.GetOrganizer(session.ID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(res).Encode(org)
	} else {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// POST /organizer/finish-profile/
func (o *OrganizerAPI) handleFinishProfile(res http.ResponseWriter, req *http.Request) {
	orgData := models.Organizer{}
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

// DELETE /organizer/logout/
func (o *OrganizerAPI) handleLogout(res http.ResponseWriter, req *http.Request, session models.Session) {
	err := o.sessMgr.DeleteSession(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GET /organizer/profile/
func (o *OrganizerAPI) handleGetProfile(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(org)
}

// POST /organizer/update-team/
func (o *OrganizerAPI) handleUpdateTeam(res http.ResponseWriter, req *http.Request, session models.Session) {
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

// POST /organizer/create-contest/
func (o *OrganizerAPI) handleCreateContest(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	contest.Organizers = []models.Organizer{org}

	err = o.orgMgr.CreateContest(contest)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// POST /organizer/delete-contest/
func (o *OrganizerAPI) handleDeleteContest(res http.ResponseWriter, req *http.Request, session models.Session) {
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

// POST /organizer/update-contest/
func (o *OrganizerAPI) handleUpdateContest(res http.ResponseWriter, req *http.Request, session models.Session) {
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

// POST /organizer/add-organizer/
func (o *OrganizerAPI) handleAddOrganizer(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	err = o.orgMgr.AddOrganizer(&newOrg)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/delete-organizer/
func (o *OrganizerAPI) handleDeleteOrganizer(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	err = o.orgMgr.DeleteOrganizer(newOrg)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/update-organizer/
func (o *OrganizerAPI) handleUpdateOrganizer(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	err = o.orgMgr.UpdateProfile(newOrg)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/delete-contestant/
func (o *OrganizerAPI) handleDeleteContestant(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil || org.Roles&models.RoleDirector == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var contestant models.Contestant
	err = json.NewDecoder(req.Body).Decode(&contestant)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.contMgr.DeleteUser(contestant)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/auto-generate-teams/
func (o *OrganizerAPI) handleAutoGenerateTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	teams := teamsgen.GenerateTeams(contest, utils.NewHardCodeNames()) // the big ass function that am proud AF from :)
	_ = json.NewEncoder(res).Encode(teams)
}

// POST /organizer/register-generated-teams/
func (o *OrganizerAPI) handleRegisterGeneratedTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
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

// POST /organizer/upload-contest-logo-file/
func (o *OrganizerAPI) handleUploadContestLogoFile(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil || org.Roles&models.RoleDirector == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	maxSize := int64(1024 * 1024)
	req.Body = http.MaxBytesReader(res, req.Body, maxSize)
	err = req.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(res, "file too big blyat!", http.StatusForbidden)
		return
	}

	file, fileHeader, err := req.FormFile("file")
	defer func() { _ = file.Close() }()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !strings.Contains(
		fileHeader.Header.Get("Content-Type"),
		"image") {

		http.Error(res, "file not of an image type!", http.StatusBadRequest)
		return
	}

	newFile, _ := os.Create( /*"./client/dist/" + */ fileHeader.Filename)
	_, _ = io.Copy(newFile, file)
}

// GET /organizer/get-contests/
func (o *OrganizerAPI) handleGetContests(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	contests, err := o.orgMgr.GetContests(org)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(res).Encode(contests)
}

// GET /organizer/get-sub-organizers/
func (o *OrganizerAPI) handleGetSubOrganizers(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	orgs, err := o.orgMgr.GetOrganizers(org)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(res).Encode(orgs)
}
