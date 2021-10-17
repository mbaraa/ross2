package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleLoginAPI holds google login handlers
type GoogleLoginAPI struct {
	googleOAuthConfig *oauth2.Config
	orgState          string
	contState         string
	config            *config.Config

	contRepo    data.ContestantGetterRepo
	orgRepo     data.OrganizerGetterRepo
	sessManager *managers.SessionManager

	endPoints map[string]http.HandlerFunc
}

// NewGoogleLoginAPI returns a new GoogleLoginAPI instance
func NewGoogleLoginAPI(sessManager *managers.SessionManager, contRepo data.ContestantGetterRepo, orgRepo data.OrganizerGetterRepo) *GoogleLoginAPI {
	return (&GoogleLoginAPI{
		config:      config.GetInstance(),
		sessManager: sessManager,
		contRepo:    contRepo,
		orgRepo:     orgRepo,
	}).
		initEndPoints().
		initOAuthConfig()
}

func (g *GoogleLoginAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	controllers.GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/gauth"), g.endPoints)
}

func (g *GoogleLoginAPI) initEndPoints() *GoogleLoginAPI {
	g.endPoints = map[string]http.HandlerFunc{
		"GET /login-callback/": g.handleCallback,
		"GET /org-login/":      g.handleOrganizerLoginWithGoogle,
		"GET /cont-login/":     g.handleContestantLoginWithGoogle,
	}
	return g
}

func (g *GoogleLoginAPI) initOAuthConfig() *GoogleLoginAPI {
	g.googleOAuthConfig = &oauth2.Config{
		RedirectURL:  g.config.GoogleCallbackHandler,
		ClientID:     g.config.GoogleClientID,
		ClientSecret: g.config.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return g
}

func (g *GoogleLoginAPI) handleContestantLoginWithGoogle(res http.ResponseWriter, req *http.Request) {
	g.contState = uuid.New().String()
	url := g.googleOAuthConfig.AuthCodeURL(g.contState)
	http.Redirect(res, req, url, http.StatusFound)
}

func (g *GoogleLoginAPI) handleOrganizerLoginWithGoogle(res http.ResponseWriter, req *http.Request) {
	g.orgState = uuid.New().String()
	url := g.googleOAuthConfig.AuthCodeURL(g.orgState)
	http.Redirect(res, req, url, http.StatusFound)
}

// handleCallback is called when authenticating with Google
// TODO
// reject un-verified emails
func (g *GoogleLoginAPI) handleCallback(res http.ResponseWriter, req *http.Request) {
	if req.FormValue("state") != g.contState && req.FormValue("state") != g.orgState {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := g.googleOAuthConfig.Exchange(context.Background(), req.FormValue("code"))
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataResponse, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	defer func() {
		err = dataResponse.Body.Close()
	}()
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	userData := make(map[string]interface{})
	err = json.NewDecoder(dataResponse.Body).Decode(&userData)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	g.finishLogin(res, req, userData)
}

func (g *GoogleLoginAPI) finishLogin(res http.ResponseWriter, req *http.Request, userData map[string]interface{}) {
	email := userData["email"].(string)

	org, err := g.orgRepo.GetByEmail(email)
	if err == nil {
		_ = g.sessManager.CreateSession(org.ID)
		http.Redirect(res, req, g.config.MachineAddress+"/organizer/login/", http.StatusPermanentRedirect)
		return
	}

	cont, err := g.contRepo.GetByEmail(email)
	if err == nil {
		_ = g.sessManager.CreateSession(cont.ID)
		http.Redirect(res, req, g.config.MachineAddress+"/contestant/login/", http.StatusPermanentRedirect)
		return
	}

	// when the front end sees the status 202 it redirects to the signup form :)
	res.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(res).Encode(userData)
	return
}
