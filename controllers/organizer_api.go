package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/utils"
	"github.com/mbaraa/ross2/utils/sheevhelper"
	"github.com/mbaraa/ross2/utils/teamsgen"
)

type OrganizerAPI struct {
	endPoints map[string]http.HandlerFunc

	orgMgr      *managers.OrganizerManager
	sessMgr     *managers.SessionManager
	teamMgr     *managers.TeamManager
	contMgr     *managers.ContestantManager
	contestRepo data.ContestCRUDRepo
	notsMgr     *managers.NotificationManager
}

func NewOrganizerAPI(orgMgr *managers.OrganizerManager, sessMgr *managers.SessionManager, teamMgr *managers.TeamManager,
	contMgr *managers.ContestantManager, contestRepo data.ContestCRUDRepo, notificationMgr *managers.NotificationManager) *OrganizerAPI {
	return (&OrganizerAPI{
		orgMgr:      orgMgr,
		sessMgr:     sessMgr,
		teamMgr:     teamMgr,
		contMgr:     contMgr,
		contestRepo: contestRepo,
		notsMgr:     notificationMgr,
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
		"POST /delete-organizer/":         o.authenticateHandler(o.handleDeleteOrganizer),
		"POST /update-organizer/":         o.authenticateHandler(o.handleUpdateOrganizer),
		"POST /delete-contestant/":        o.authenticateHandler(o.handleDeleteContestant),
		"GET /get-sub-organizers/":        o.authenticateHandler(o.handleGetSubOrganizers),
		"POST /auto-generate-teams/":      o.authenticateHandler(o.handleAutoGenerateTeams),
		// "POST /man-generate-teams/":         nil,
		"POST /register-generated-teams/": o.authenticateHandler(o.handleRegisterGeneratedTeams),
		"POST /update-teams/":             o.authenticateHandler(o.handleUpdateTeams),
		// "GET /get-contestants-for-contest/": nil,
		// "GET /get-organizers-for-contest/":  nil,

		"GET /get-contests/":              o.authenticateHandler(o.handleGetContests),
		"POST /get-contest/":              o.authenticateHandler(o.handleGetContest),
		"POST /send-sheev-notifications/": o.authenticateHandler(o.handleSendSheevNotifications),
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

	newOrg.ContactInfo = models.ContactInfo{
		FacebookURL: "/",
	}
	newOrg.DirectorID = org.ID

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

// POST /organizer/auto-generate-teams/?gen-type="numbered"|"random"
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

	contest, err = o.contestRepo.Get(contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	teams, leftTeamless := teamsgen.GenerateTeams(contest,
		utils.GetNamesGetter(req.URL.Query().Get("gen-type"))) // the big ass function that am proud AF from :)

	_ = json.NewEncoder(res).Encode(map[string]interface{}{
		"teams":         teams,
		"left_teamless": leftTeamless,
	})
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

// POST /organizer/update-teams/
func (o *OrganizerAPI) handleUpdateTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil || org.Roles&models.RoleDirector == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	teamsAndRemovedConts := struct {
		Teams              []models.Team       `json:"teams"`
		RemovedContestants []models.Contestant `json:"removed_contestants"`
	}{} // God will burn me very deep in hell for this :)

	err = json.NewDecoder(req.Body).Decode(&teamsAndRemovedConts)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.teamMgr.UpdateTeams(
		teamsAndRemovedConts.Teams,
		teamsAndRemovedConts.RemovedContestants,
		org)
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

	uploadedFilePath := fileHeader.Filename
	if !config.GetInstance().Development {
		uploadedFilePath = "./client/dist/" + uploadedFilePath
	}
	newFile, _ := os.Create(uploadedFilePath)
	_, _ = io.Copy(newFile, file)
}

// GET /organizer/get-contests/
// 🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂🙂
// will fix later I swear 😉
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

// POST /organizer/get-contest/
// this is needed because the get contest in the ContestAPI doesn't get private teams :)
func (o *OrganizerAPI) handleGetContest(res http.ResponseWriter, req *http.Request, session models.Session) {
	var contest models.Contest
	err := json.NewDecoder(req.Body).Decode(&contest)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	contest, err = o.contestRepo.Get(contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(contest)
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

// POST /organizer/send-sheev-notifications/
// TODO:
// do this using a cron job :)
func (o *OrganizerAPI) handleSendSheevNotifications(res http.ResponseWriter, req *http.Request, session models.Session) {
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

	contest, err = o.contestRepo.Get(contest) // lazy loading :)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	notifications, err := sheevhelper.GetSheevNotifications(contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = o.notsMgr.SendMany(notifications)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
