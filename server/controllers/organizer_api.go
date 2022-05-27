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
		"GET /profile/":         o.handleGetProfile,
		"POST /finish-profile/": o.handleFinishProfile,

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
		"POST /update-organizer/":         o.handleUpdateOrganizer,
		"POST /delete-organizer/":         o.handleDeleteOrganizer,
		"POST /get-sub-organizers/":       o.handleGetSubOrganizers,
		"POST /generate-teams/":           o.handleGenerateTeams,
		"POST /save-teams/":               o.handleSaveTeams,

		"POST /get-contest/":              o.handleGetContest,
		"POST /send-sheev-notifications/": o.handleSendSheevNotifications,
		"POST /get-participants-csv/":     o.handleGetParticipantsCSV,
		"POST /get-teams-csv/":            o.handleGetTeamsCSV,
		"POST /generate-teams-posts/":     o.handleGenerateTeamsPosts,
		"GET /get-all-users/":             o.handleGetAllUsers,
		"POST /check-role/":               o.handleCheckRole,
		"POST /get-org-roles/":            o.handleGetRoles,

		"POST /get-participants/":            o.handleGetParticipants,
		"POST /mark-participant-as-present/": o.markParticipantAsPresent,
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

// POST /finish-profile/
func (o *OrganizerAPI) handleFinishProfile(ctx context.HandlerContext) {
	var org models.Organizer
	if ctx.ReadJSON(&org) != nil {
		return
	}

	if o.orgMgr.FinishProfile(org) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	if err != nil || (org.User.UserType&enums.UserTypeDirector) == 0 {
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
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
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
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
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
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		NewOrg  models.Organizer `json:"organizer"`
		Contest models.Contest   `json:"contest"`
		Roles   float64          `json:"roles"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, director.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	baseUser, err := o.orgMgr.GetUserProfileUsingEmail(reqBody.NewOrg.User.Email)
	if err != nil {
		http.Error(ctx.Res, "user doesn't exist!", http.StatusNotFound)
		return
	}

	err = o.orgMgr.AddOrganizer(reqBody.NewOrg, director, baseUser, reqBody.Contest, enums.OrganizerRole(reqBody.Roles))
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /organizer/update-organizer/
func (o *OrganizerAPI) handleUpdateOrganizer(ctx context.HandlerContext) {
	director, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		UpdatedOrg models.Organizer `json:"organizer"`
		Contest    models.Contest   `json:"contest"`
		Roles      float64          `json:"roles"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, director.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	baseUser, err := o.orgMgr.GetUserProfileUsingEmail(reqBody.UpdatedOrg.User.Email)
	if err != nil {
		http.Error(ctx.Res, "user doesn't exist!", http.StatusNotFound)
		return
	}

	err = o.orgMgr.UpdateOrganizer(reqBody.UpdatedOrg, director, baseUser, reqBody.Contest, enums.OrganizerRole(reqBody.Roles))
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /organizer/delete-organizer/
func (o *OrganizerAPI) handleDeleteOrganizer(ctx context.HandlerContext) {
	director, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		Org     models.Organizer `json:"organizer"`
		Contest models.Contest   `json:"contest"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, director.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if o.orgMgr.DeleteOrganizer(reqBody.Org, reqBody.Contest) != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST /organizer/generate-teams/?gen-type=numbered|random|given
func (o *OrganizerAPI) handleGenerateTeams(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		Contest models.Contest `json:"contest"`
		Names   []string       `json:"names"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	teams, leftTeamless, err :=
		o.orgMgr.GenerateTeams(reqBody.Contest, ctx.Req.URL.Query().Get("gen-type"), reqBody.Names)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(map[string]interface{}{
		"teams":         teams,
		"left_teamless": leftTeamless,
	}, 0)
}

// POST /organizer/save-teams/
func (o *OrganizerAPI) handleSaveTeams(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody struct {
		Teams    []models.Team       `json:"teams"`
		Teamless []models.Contestant `json:"teamless"`
		Contest  models.Contest      `json:"contest"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = o.orgMgr.CreateUpdateTeams(reqBody.Teams, reqBody.Teamless, reqBody.Contest, org)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /organizer/upload-contest-logo-file/
func (o *OrganizerAPI) handleUploadContestLogoFile(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (org.User.UserType&enums.UserTypeDirector) == 0 {
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

	uploadedFilePath := config.GetInstance().UploadDirectory + fileHeader.Filename
	newFile, _ := os.Create(uploadedFilePath)
	_, _ = io.Copy(newFile, file)
	_ = file.Close()
}

// POST /organizer/get-contest/
// this is needed because the get contest in the ContestAPI doesn't get private teams :)
func (o *OrganizerAPI) handleGetContest(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	contest, err = o.orgMgr.GetContest(contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(contest, 0)
}

// POST /organizer/get-sub-organizers/
func (o *OrganizerAPI) handleGetSubOrganizers(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	orgs, err := o.orgMgr.GetOrganizers(org, contest)
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
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
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
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	contsCSV, err := o.orgMgr.GetParticipantsCSV(contest)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = ctx.Res.Write([]byte(contsCSV))
}

// POST /organizer/get-teams-csv/
func (o *OrganizerAPI) handleGetTeamsCSV(ctx context.HandlerContext) {
	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	contsCSV, err := o.orgMgr.GetTeamsCSV(contest)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = ctx.Res.Write([]byte(contsCSV))
}

// POST /organizer/generate-teams-posts/
func (o *OrganizerAPI) handleGenerateTeamsPosts(ctx context.HandlerContext) {
	// 🙉🙊🙈 if it works it ain't stupid
	var reqBody struct {
		Contest        models.Contest            `json:"contest"`
		TeamNameProps  postsgen.TextFieldProps   `json:"teamNameProps"`
		TeamOrderProps postsgen.TextFieldProps   `json:"teamOrderProps"`
		MembersProps   []postsgen.TextFieldProps `json:"membersNamesProps"`
		BaseImage      string                    `json:"baseImage"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	var err error
	reqBody.Contest, err = o.orgMgr.GetContest(reqBody.Contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !o.orgMgr.CheckOrgRole(enums.RoleDirector, reqBody.Contest.ID, org.ID) {
		ctx.Res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var postsGen *postsgen.TeamPostsGenerator
	if reqBody.BaseImage == "" {
		postsGen, err = postsgen.GetThreeMembersPostSamplePostBuilder().
			Teams(reqBody.Contest.Teams).
			GetTeamsPostsGenerator()
	} else {
		postsGen, err = postsgen.NewTeamsPostsGeneratorBuilder().
			Teams(reqBody.Contest.Teams).
			TeamNameProps(reqBody.TeamNameProps).
			TeamOrderProps(reqBody.TeamOrderProps).
			MembersNamesProps(reqBody.MembersProps).
			B64Image(reqBody.BaseImage).
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
	if err != nil || (org.User.UserType&enums.UserTypeDirector) == 0 {
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

// POST /get-participants/
func (o *OrganizerAPI) handleGetParticipants(ctx context.HandlerContext) {
	var contest models.Contest
	if ctx.ReadJSON(&contest) != nil {
		return
	}

	org, err := o.orgMgr.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	participants, err := o.orgMgr.GetParticipants(contest, org)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(participants, 0)
}

// POST /mark-participant-as-present/
func (o *OrganizerAPI) markParticipantAsPresent(ctx context.HandlerContext) {
	var reqBody struct {
		User    models.User    `json:"user"`
		Contest models.Contest `json:"contest"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	err := o.orgMgr.MarkAttendance(reqBody.User, reqBody.Contest)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// POST /check-role/
func (o *OrganizerAPI) handleCheckRole(ctx context.HandlerContext) {
	var reqBody struct {
		ContestID   float64 `json:"contest_id"`
		OrganizerID float64 `json:"organizer_id"`
		Roles       float64 `json:"roles"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	authorized :=
		o.orgMgr.CheckOrgRole(enums.OrganizerRole(reqBody.Roles), uint(reqBody.ContestID), uint(reqBody.OrganizerID))

	_ = ctx.WriteJSON(authorized, 0)
}

// TODO:
// move this to get contest's orgs
// POST /get-org-roles/
func (o *OrganizerAPI) handleGetRoles(ctx context.HandlerContext) {
	var reqBody struct {
		ContestID   float64 `json:"contest_id"`
		OrganizerID float64 `json:"organizer_id"`
	}
	if ctx.ReadJSON(&reqBody) != nil {
		return
	}

	roles, rolesNames, err := o.orgMgr.GetOrgRoles(uint(reqBody.OrganizerID), uint(reqBody.ContestID))
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(map[string]interface{}{
		"roles":       roles,
		"roles_names": rolesNames,
	}, 0)
}
