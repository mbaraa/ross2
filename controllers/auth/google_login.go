package auth

import (
	"encoding/json"
	"errors"
	"github.com/mbaraa/ross2/config"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

// GoogleLoginAPI holds google login handlers
type GoogleLoginAPI struct {
	contRepo    data.ContestantCRUDRepo
	orgRepo     data.OrganizerGetterRepo
	sessManager *managers.SessionManager

	endPoints map[string]http.HandlerFunc
}

// NewGoogleLoginAPI returns a new GoogleLoginAPI instance
func NewGoogleLoginAPI(sessManager *managers.SessionManager, contRepo data.ContestantCRUDRepo, orgRepo data.OrganizerGetterRepo) *GoogleLoginAPI {
	return (&GoogleLoginAPI{
		sessManager: sessManager,
		contRepo:    contRepo,
		orgRepo:     orgRepo,
	}).
		initEndPoints()
}

func (g *GoogleLoginAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	controllers.GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/gauth"), g.endPoints)
}

func (g *GoogleLoginAPI) initEndPoints() *GoogleLoginAPI {
	g.endPoints = map[string]http.HandlerFunc{
		"POST /cont-login/": g.handleContestantLogin,
		"POST /org-login/":  g.handleOrganizerLogin,
	}
	return g
}

func (g *GoogleLoginAPI) finishLogin(userID uint, res http.ResponseWriter) {
	sess, err := g.sessManager.CreateSession(userID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = res.Write([]byte(`{"token" : "` + sess.ID + `" }`))
}

// contestant stuff

func (g *GoogleLoginAPI) handleContestantLogin(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")

	var cont models.Contestant
	err := json.NewDecoder(req.Body).Decode(&cont)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	g.checkGoogleJWTToken(token, res)
	g.finishContestantLogin(cont, res)
}

func (g *GoogleLoginAPI) finishContestantLogin(cont0 models.Contestant, res http.ResponseWriter) {
	cont, err := g.contRepo.GetByEmail(cont0.Email)
	if err != nil {
		cont = models.Contestant{
			Name:            cont0.Name,
			Email:           cont0.Email,
			AvatarURL:       cont0.AvatarURL,
			ProfileFinished: false,
			ContactInfo: models.ContactInfo{
				FacebookURL:    "/",
				WhatsappNumber: "/",
				TelegramNumber: "/",
			},
		}
		err = g.contRepo.Add(&cont)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	g.finishLogin(cont.ID, res)
}

// organizer stuff

func (g *GoogleLoginAPI) handleOrganizerLogin(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")

	var org models.Organizer
	err := json.NewDecoder(req.Body).Decode(&org)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if g.checkGoogleJWTToken(token, res) {
		g.finishOrganizerLogin(org, res)
	}
}

func (g *GoogleLoginAPI) finishOrganizerLogin(org0 models.Organizer, res http.ResponseWriter) {
	org, err := g.orgRepo.GetByEmail(org0.Email)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	g.finishLogin(org.ID, res)
}

// Google JWT validation stuff

func (g *GoogleLoginAPI) checkGoogleJWTToken(token string, res http.ResponseWriter) bool {
	_, err := g.validateGoogleJWT(token)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

func (g *GoogleLoginAPI) validateGoogleJWT(tokenString string) (googleClaims, error) {
	gclaims := googleClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &gclaims, func(token *jwt.Token) (interface{}, error) {
		pem, err := g.getGooglePublicKey(token.Header["kid"].(string))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return googleClaims{}, err
	}

	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return googleClaims{}, err
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return googleClaims{}, err
	}

	if claims.Audience != config.GetInstance().GoogleClientID {
		return googleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return googleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func (g *GoogleLoginAPI) getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

type googleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}
