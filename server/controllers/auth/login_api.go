/*
Package auth holds user login related stuff
*/
package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/models"
)

// OAuthLoginAPI holds google login handlers
type OAuthLoginAPI struct {
	userMgr        *helpers.UserHelper
	tokenValidator JWTTokenValidator
	apiEndPoint    string
	endPoints      map[string]http.HandlerFunc
}

// NewOAuthLoginAPI returns a new OAuthLoginAPI instance
func NewOAuthLoginAPI(userManager *helpers.UserHelper, tokenValidator JWTTokenValidator, apiEndPoint string) *OAuthLoginAPI {
	return (&OAuthLoginAPI{
		userMgr:        userManager,
		tokenValidator: tokenValidator,
		apiEndPoint:    apiEndPoint,
	}).
		initEndPoints()
}

func (l *OAuthLoginAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", config.GetInstance().AllowedClients)
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if handler, exists := l.endPoints[req.Method+" "+strings.TrimPrefix(req.URL.Path, l.apiEndPoint)]; exists {
		handler(res, req)
		return
	}
	if req.Method != http.MethodOptions {
		http.NotFound(res, req)
	}
}

func (l *OAuthLoginAPI) initEndPoints() *OAuthLoginAPI {
	l.endPoints = map[string]http.HandlerFunc{
		"POST /login/":       l.handleLogin,
		"POST /login-token/": l.handleLoginWithToken,
		"POST /logout/":      l.handleLogout,
	}
	return l
}

func (l *OAuthLoginAPI) handleLogin(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = l.tokenValidator.Validate(
		req.Header.Get("Authorization"))
	if err == nil {
		l.finishLogin(res, user)
	} else {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}
}

func (l *OAuthLoginAPI) handleLoginWithToken(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")
	user, err := l.userMgr.LoginUsingSession(token)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(res).Encode(user)
		return
	}

	_ = json.NewEncoder(res).Encode(user)
}

func (l *OAuthLoginAPI) finishLogin(res http.ResponseWriter, user models.User) {
	sess, err := l.userMgr.Login(&user)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = res.Write([]byte(`{"token" : "` + sess.ID + `"}`))
}

func (l *OAuthLoginAPI) handleLogout(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = l.userMgr.Logout(user,
		req.Header.Get("Authorization"))
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
}
