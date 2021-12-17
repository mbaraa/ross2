package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/namesgetter"
	"github.com/mbaraa/ross2/utils/partsexport"
	"github.com/mbaraa/ross2/utils/postsgen"
	"github.com/mbaraa/ross2/utils/sheevhelper"
	"github.com/mbaraa/ross2/utils/teamsgen"
)

type OrganizerAPIBuilder struct {
	orgMgr          *managers.OrganizerManager
	sessMgr         *managers.SessionManager
	teamMgr         *managers.TeamManager
	contMgr         *managers.ContestantManager
	contestRepo     data.ContestCRUDRepo
	notificationMgr *managers.NotificationManager
}

func NewOrganizerAPIBuilder() *OrganizerAPIBuilder {
	return new(OrganizerAPIBuilder)
}

func (b *OrganizerAPIBuilder) OrganizerMgr(o *managers.OrganizerManager) *OrganizerAPIBuilder {
	b.orgMgr = o
	return b
}

func (b *OrganizerAPIBuilder) SessionMgr(s *managers.SessionManager) *OrganizerAPIBuilder {
	b.sessMgr = s
	return b
}

func (b *OrganizerAPIBuilder) TeamMgr(t *managers.TeamManager) *OrganizerAPIBuilder {
	b.teamMgr = t
	return b
}

func (b *OrganizerAPIBuilder) ContestantMgr(c *managers.ContestantManager) *OrganizerAPIBuilder {
	b.contMgr = c
	return b
}

func (b *OrganizerAPIBuilder) ContestRepo(c data.ContestCRUDRepo) *OrganizerAPIBuilder {
	b.contestRepo = c
	return b
}

func (b *OrganizerAPIBuilder) NotificationMgr(n *managers.NotificationManager) *OrganizerAPIBuilder {
	b.notificationMgr = n
	return b
}

func (b *OrganizerAPIBuilder) verify() bool {
	if b.contestRepo == nil {
		fmt.Println("Organizer API Builder: missing contest repo!")
	}
	if b.orgMgr == nil {
		fmt.Println("Organizer API Builder: missing organizer manager!")
	}
	if b.sessMgr == nil {
		fmt.Println("Organizer API Builder: missing session manager!")
	}
	if b.teamMgr == nil {
		fmt.Println("Organizer API Builder: missing team manager!")
	}
	if b.contMgr == nil {
		fmt.Println("Organizer API Builder: missing contestant manager!")
	}
	if b.notificationMgr == nil {
		fmt.Println("Organizer API Builder: missing notification manager!")
	}

	return b.contestRepo != nil && b.orgMgr != nil &&
		b.sessMgr != nil && b.notificationMgr != nil &&
		b.teamMgr != nil && b.contMgr != nil
}

func (b *OrganizerAPIBuilder) GetOrganizerAPI() *OrganizerAPI {
	if !b.verify() {
		return nil
	}
	return NewOrganizerAPI(b)
}

type OrganizerAPI struct {
	endPoints map[string]http.HandlerFunc

	orgMgr      *managers.OrganizerManager
	sessMgr     *managers.SessionManager
	teamMgr     *managers.TeamManager
	contMgr     *managers.ContestantManager
	contestRepo data.ContestCRUDRepo
	notsMgr     *managers.NotificationManager
}

func NewOrganizerAPI(b *OrganizerAPIBuilder) *OrganizerAPI {
	return (&OrganizerAPI{
		orgMgr:      b.orgMgr,
		sessMgr:     b.sessMgr,
		teamMgr:     b.teamMgr,
		contMgr:     b.contMgr,
		contestRepo: b.contestRepo,
		notsMgr:     b.notificationMgr,
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
		"POST /get-participants/":         o.authenticateHandler(o.handleGetParticipants),
		"POST /generate-teams-posts/":     o.authenticateHandler(o.handleGenerateTeamsPosts),
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
	if err != nil || (org.Roles&enums.RoleDirector == 0 || org.Roles&enums.RoleCoreOrganizer == 0) {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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

	newOrg.User.ContactInfo = models.ContactInfo{
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp := struct {
		Contest models.Contest `json:"contest"`
		Names   []string       `json:"names"`
	}{}
	err = json.NewDecoder(req.Body).Decode(&resp)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	contest, err := o.contestRepo.Get(resp.Contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	teams, leftTeamless :=
		teamsgen.GenerateTeams(contest, // the big ass function that am proud AF from :)
			namesgetter.GetNamesGetter(req.URL.Query().Get("gen-type"), resp.Names...))

	_ = json.NewEncoder(res).Encode(map[string]interface{}{
		"teams":         teams,
		"left_teamless": leftTeamless,
	})
}

// POST /organizer/register-generated-teams/
func (o *OrganizerAPI) handleRegisterGeneratedTeams(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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

// POST /organizer/get-participants/
func (o *OrganizerAPI) handleGetParticipants(res http.ResponseWriter, req *http.Request, session models.Session) {
	org, err := o.orgMgr.GetOrganizer(session.ID)
	if err != nil || org.Roles&enums.RoleDirector == 0 {
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

	_, _ = res.Write([]byte(partsexport.GetParticipants(contest)))
}

// POST /organizer/generate-teams-posts/
func (o *OrganizerAPI) handleGenerateTeamsPosts(res http.ResponseWriter, req *http.Request, session models.Session) {
	// 🙉🙊🙈 if it works it ain't stupid
	var respBody struct {
		Contest        models.Contest            `json:"contest"`
		TeamNameProps  postsgen.TextFieldProps   `json:"teamNameProps"`
		TeamOrderProps postsgen.TextFieldProps   `json:"teamOrderProps"`
		MembersProps   []postsgen.TextFieldProps `json:"membersNamesProps"`
		BaseImage      string                    `json:"baseImage"`
	}
	err := json.NewDecoder(req.Body).Decode(&respBody)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody.Contest, err = o.contestRepo.Get(respBody.Contest)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var postsGen *postsgen.TeamPostsGenerator
	if respBody.BaseImage == "" {
		postsGen, err = postsgen.ThreeMembersPostSamplePostBuilder.
			Teams(respBody.Contest.Teams).
			GetTeamsPostsGenerator()

	} else {
		postsGen, err = postsgen.NewTeamsPostsGeneratorBuilder().
			Teams(respBody.Contest.Teams).
			TeamNameProps(respBody.TeamNameProps).
			TeamOrderProps(respBody.TeamOrderProps).
			MembersNamesProps(respBody.MembersProps).
			B64Image(respBody.BaseImage).
			GetTeamsPostsGenerator()
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	zipFileBytes, err := postsGen.GenerateToZipFileBytes()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = res.Write(zipFileBytes)
}
