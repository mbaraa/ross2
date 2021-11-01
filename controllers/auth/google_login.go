package auth

import (
	"encoding/json"
	"errors"
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
		"POST /login/": g.handleLoginWithGoogle,
	}
	return g
}

func (g *GoogleLoginAPI) handleLoginWithGoogle(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")

	var cont models.User
	err := json.NewDecoder(req.Body).Decode(&cont)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = g.validateGoogleJWT(token)
	if err != nil {
		res.WriteHeader(http.StatusForbidden)
		return
	}

	g.finishLogin(res, req, cont)
}

func (g *GoogleLoginAPI) finishLogin(res http.ResponseWriter, req *http.Request, userData models.User) {
	//org, err := g.orgRepo.GetByEmail(userData.Email)
	//if err == nil {
	//sess, _ := g.sessManager.CreateSession(org.ID)
	//req.Header.Set("Authorization", sess.ID)
	//http.Redirect(res, req, g.config.MachineAddress+"/organizer/login/", http.StatusPermanentRedirect)
	//return
	//}

	cont, err := g.contRepo.GetByEmail(userData.Email)
	if err != nil {
		cont = models.Contestant{
			User: models.User{
				Name:            userData.Name,
				Email:           userData.Email,
				AvatarURL:       userData.AvatarURL,
				ProfileFinished: false,
				ContactInfo: models.ContactInfo{
					FacebookURL:    "/",
					WhatsappNumber: "/",
					TelegramNumber: "/",
				},
			},
		}
		err = g.contRepo.Add(&cont)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	sess, err := g.sessManager.CreateSession(cont.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(res).Encode(map[string]interface{}{
		"token": sess.ID,
	})
}

func (g *GoogleLoginAPI) validateGoogleJWT(tokenString string) (googleClaims, error) {
	gclaims := googleClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &gclaims, func(token *jwt.Token) (interface{}, error) {
		pem, err := getGooglePublicKey(token.Header["kid"].(string))
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

	if claims.Audience != "202655727003-gu3umksjmog90n6oonvfeh79msbe1j1e.apps.googleusercontent.com" {
		return googleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return googleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func getGooglePublicKey(keyID string) (string, error) {
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
