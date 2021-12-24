package controllers

import (
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/context"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/models/enums"
)

// AdminAPI manages the operations done by a admin
type AdminAPI struct {
	helper        *helpers.AdminHelper
	authenticator *auth.HandlerAuthenticator
	endPoints     map[string]http.HandlerFunc
}

// NewAdminAPI returns a new AdminAPI instance
func NewAdminAPI(helper *helpers.AdminHelper, authenticator *auth.HandlerAuthenticator) *AdminAPI {
	return (&AdminAPI{
		helper:        helper,
		authenticator: authenticator,
	}).initEndPoints()
}

// ServeHTTP is the API's magical port :), since it makes it implement http.Handler
func (a *AdminAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/admin"), a.endPoints)
}

func (a *AdminAPI) initEndPoints() *AdminAPI {
	a.endPoints = a.authenticator.AuthenticateHandlers(map[string]auth.HandlerFunc{
		"GET /profile/": a.handleGetProfile,

		"POST /add-director/":    a.handleAddDirector,
		"POST /delete-director/": a.handleDeleteDirector,

		"GET /get-directors/": a.handleGetDirectors,
	})
	return a
}

// GET /admin/profile/
func (a *AdminAPI) handleGetProfile(ctx context.HandlerContext) {
	admin, err := a.helper.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (admin.User.UserType&enums.UserTypeAdmin) == 0 {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}

	_ = ctx.WriteJSON(admin, 0)
}

// POST /admin/add-director/
func (a *AdminAPI) handleAddDirector(ctx context.HandlerContext) {
	admin, err := a.helper.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (admin.User.UserType&enums.UserTypeAdmin) == 0 {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}

	var director models.Organizer
	if ctx.ReadJSON(&director) != nil {
		return
	}

	baseUser, err := a.helper.GetUserProfileUsingEmail(director.User.Email)
	if err != nil {
		http.Error(ctx.Res, "user doesn't exist!", http.StatusNotFound)
		return
	}

	err = a.helper.AddDirector(director, baseUser)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}
}

// POST /admin/delete-director/
func (a *AdminAPI) handleDeleteDirector(ctx context.HandlerContext) {
	admin, err := a.helper.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (admin.User.UserType&enums.UserTypeAdmin) == 0 {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}

	var director models.Organizer
	if ctx.ReadJSON(&director) != nil {
		return
	}

	err = a.helper.DeleteDirector(director)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}
}

// GET /admin/get-directors/
func (a *AdminAPI) handleGetDirectors(ctx context.HandlerContext) {
	admin, err := a.helper.GetProfile(models.User{ID: ctx.Sess.UserID})
	if err != nil || (admin.User.UserType&enums.UserTypeAdmin) == 0 {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}

	dirs, err := a.helper.GetDirectors()
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = ctx.WriteJSON(dirs, 0)
}
