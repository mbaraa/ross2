package controllers

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/context"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
	"github.com/mbaraa/ross2/utils/postsgen"
)

type OrganizerAPI struct {
	orgMgr        *helpers.OrganizerHelper
	authenticator *auth.HandlerAuthenticator
	endPoints     map[string]http.HandlerFunc
}

func NewOrganizerAPI(orgMgr *helpers.OrganizerHelper, authenticator *auth.HandlerAuthenticator) *OrganizerAPI {
	return (&OrganizerAPI{
		orgMgr:        orgMgr,
		authenticator: authenticator,
	}).initEndPoints()
}

func (o *OrganizerAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/organizer"), o.endPoints)
}

func (o *OrganizerAPI) initEndPoints() *OrganizerAPI {
	o.endPoints = o.authenticator.AuthenticateHandlers(map[string]auth.HandlerFunc{
		// shared organizer/director operations
		"GET /profile/": o.handleGetProfile,

		// "GET /get-solved-problems/": nil,
		// "POST /update-solved-problems/": nil,
		// "POST /set-pc-status/": nil,
		// "POST /upload-contest-media/": nil,
		"POST /update-team/": o.handleUpdateTeam,

		// director operations
		"POST /create-contest/":           o.handleCreateContest,
		"POST /delete-contest/":           o.handleDeleteContest,
		"POST /update-contest/":           o.handleUpdateContest,
		"POST /upload-contest-logo-file/": o.handleUploadContestLogoFile,
		"POST /add-organizer/":            o.handleAddOrganizer,
		"POST /delete-organizer/":         o.handleDeleteOrganizer,
		"GET /get-sub-organizers/":        o.handleGetSubOrganizers,
		"POST /generate-teams/":           o.handleGenerateTeams,
		"POST /register-generated-teams/": o.handleRegisterGeneratedTeams,
		"POST /update-teams/":             o.handleUpdateTeams,
		// "GET /get-contestants-for-contest/": nil,
		// "GET /get-organizers-for-contest/":  nil,

		"GET /get-contests/":              o.handleGetContests,
		"POST /get-contest/":              o.handleGetContest,
		"POST /send-sheev-notifications/": o.handleSendSheevNotifications,
		"POST /get-participants-csv/":     o.handleGetParticipantsCSV,
		"POST /generate-teams-posts/":     o.handleGenerateTeamsPosts,
		"GET /get-all-users/":             o.handleGetAllUsers,
	})
	return o
}

// GET /organizer/profile/
func (o *OrganizerAPI) handleGetProfile(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(org, 0)
}

// POST /organizer/update-team/
func (o *OrganizerAPI) handleUpdateTeam(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.User.UserType&enums.UserTypeOrganizer) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var team models.Team
	if ctx.ReadJSON(&team) != nil {
		return
	}

	if o.orgMgr.UpdateTeam(team, org) != nil {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// POST /organizer/create-contest/
func (o *OrganizerAPI) handleCreateContest(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if o.orgMgr.CreateContest(contest, org) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/delete-contest/
func (o *OrganizerAPI) handleDeleteContest(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if o.orgMgr.DeleteContest(contest) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/update-contest/
func (o *OrganizerAPI) handleUpdateContest(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if o.orgMgr.UpdateContest(contest) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/add-organizer/
func (o *OrganizerAPI) handleAddOrganizer(ctx context.HandlerContext) {
	director, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (director.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var newOrg models.Organizer
	if ctx.ReadJSON(&newOrg) != nil {
		return
	}

	if o.orgMgr.AddOrganizer(newOrg, director) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/delete-organizer/
func (o *OrganizerAPI) handleDeleteOrganizer(ctx context.HandlerContext) {
	director, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (director.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var org models.Organizer
	if ctx.ReadJSON(&org) != nil {
		return
	}

	if o.orgMgr.DeleteOrganizer(org) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/generate-teams/?gen-type=numbered|random|given
func (o *OrganizerAPI) handleGenerateTeams(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respBody struct {
		Contest models.Contest `json:"contest"`
		Names   []string       `json:"names"`
	}
	if ctx.ReadJSON(&respBody) != nil {
		return
	}

	teams, leftTeamless, err :=
		o.orgMgr.GenerateTeams(respBody.Contest, ctx.Req.URL.Query().Get("gen-type"), respBody.Names)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(map[string]interface{}{
		"teams":         teams,
		"left_teamless": leftTeamless,
	}, 0)
}

// POST /organizer/register-generated-teams/
func (o *OrganizerAPI) handleRegisterGeneratedTeams(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var teams []*models.Team
	if ctx.ReadJSON(&teams) != nil {
		return
	}

	err = o.orgMgr.CreateTeams(teams)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /organizer/update-teams/
func (o *OrganizerAPI) handleUpdateTeams(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respBody struct {
		Teams              []models.Team       `json:"teams"`
		RemovedContestants []models.Contestant `json:"removed_contestants"`
	}
	if ctx.ReadJSON(&respBody) != nil {
		return
	}

	err = o.orgMgr.UpdateTeams(respBody.Teams, respBody.RemovedContestants, org)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /organizer/upload-contest-logo-file/
func (o *OrganizerAPI) handleUploadContestLogoFile(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	maxSize := int64(1024 * 1024)
	ctx.Req.Body = http.MaxBytesReader(ctx.Res, ctx.Req.Body, maxSize)
	err = ctx.Req.ParseMultipartForm(maxSize)
	if err != nil {
		http.Error(ctx.Res, "file too big (max allowed size is 10MiB)", http.StatusForbidden)
		return
	}

	file, fileHeader, err := ctx.Req.FormFile("file")
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusBadRequest)
		return
	}

	if !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") {
		http.Error(ctx.Res, "only user an image type!", http.StatusBadRequest)
		return
	}

	uploadedFilePath := fileHeader.Filename
	if !config.GetInstance().Development {
		uploadedFilePath = "./client/dist/" + uploadedFilePath
	}
	newFile, _ := os.Create(uploadedFilePath)
	_, _ = io.Copy(newFile, file)
	_ = file.Close()
}

// GET /organizer/get-contests/
func (o *OrganizerAPI) handleGetContests(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	contests, err := o.orgMgr.GetContests(org)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(contests, 0)
}

// POST /organizer/get-contest/
// this is needed because the get contest in the ContestAPI doesn't get private teams :)
func (o *OrganizerAPI) handleGetContest(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	contest, err = o.orgMgr.GetContest(contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(contest, 0)
}

// GET /organizer/get-sub-organizers/
func (o *OrganizerAPI) handleGetSubOrganizers(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	orgs, err := o.orgMgr.GetOrganizers(org)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusNotFound)
		return
	}

	_ = ctx.WriteJSON(orgs, 0)
}

// POST /organizer/send-sheev-notifications/
// TODO:
// do this using a cron job :)
func (o *OrganizerAPI) handleSendSheevNotifications(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	err = o.orgMgr.SendSheevNotifications(contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
	}
}

// POST /organizer/get-participants-csv/
func (o *OrganizerAPI) handleGetParticipantsCSV(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	contsCSV, err := o.orgMgr.GetParticipantsCSV(contest)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = ctx.Res.Write([]byte(contsCSV))
}

// POST /organizer/generate-teams-posts/
func (o *OrganizerAPI) handleGenerateTeamsPosts(ctx context.HandlerContext) {
	// 🙉🙊🙈 if it works it ain't stupid
	var respBody struct {
		Contest        models.Contest            `json:"contest"`
		TeamNameProps  postsgen.TextFieldProps   `json:"teamNameProps"`
		TeamOrderProps postsgen.TextFieldProps   `json:"teamOrderProps"`
		MembersProps   []postsgen.TextFieldProps `json:"membersNamesProps"`
		BaseImage      string                    `json:"baseImage"`
	}
	if ctx.ReadJSON(&respBody) != nil {
		return
	}

	var err error
	respBody.Contest, err = o.orgMgr.GetContest(respBody.Contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
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
		http.Error(ctx.Res, err.Error(), http.StatusBadRequest)
		return
	}

	zipFileBytes, err := postsGen.GenerateToZipFileBytes()
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = ctx.Res.Write(zipFileBytes)
}

// GET /organizer/get-all-users/
// since a director is already approved by the admin they should be able to SEE all users details
// in order to add them as organizers
func (o *OrganizerAPI) handleGetAllUsers(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.Roles&enums.RoleDirector) == 0 {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := o.orgMgr.GetNonOrgUsers()
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = ctx.WriteJSON(users, 0)
}
