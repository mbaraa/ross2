package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

// MicrosoftLoginAPI holds microsoft login handlers
// TODO
// exchange id token w/ microsoft for safer login 🙂
type MicrosoftLoginAPI struct {
	contRepo    data.ContestantCRUDRepo
	orgRepo     data.OrganizerGetterRepo
	sessManager *managers.SessionManager

	endPoints map[string]http.HandlerFunc
}

// NewMicrosoftLoginAPI returns a new MicrosoftLoginAPI instance
func NewMicrosoftLoginAPI(sessManager *managers.SessionManager, contRepo data.ContestantCRUDRepo, orgRepo data.OrganizerGetterRepo) *MicrosoftLoginAPI {
	return (&MicrosoftLoginAPI{
		sessManager: sessManager,
		contRepo:    contRepo,
		orgRepo:     orgRepo,
	}).
		initEndPoints()
}

func (m *MicrosoftLoginAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	controllers.GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/msauth"), m.endPoints)
}

func (m *MicrosoftLoginAPI) initEndPoints() *MicrosoftLoginAPI {
	m.endPoints = map[string]http.HandlerFunc{
		"POST /cont-login/": m.handleContestantLogin,
		"POST /org-login/":  m.handleOrganizerLogin,
	}
	return m
}

func (m *MicrosoftLoginAPI) finishLogin(userID uint, res http.ResponseWriter) {
	sess, err := m.sessManager.CreateSession(userID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = res.Write([]byte(`{"token" : "` + sess.ID + `" }`))
}

// contestant stuff

func (m *MicrosoftLoginAPI) handleContestantLogin(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")

	var cont models.Contestant
	err := json.NewDecoder(req.Body).Decode(&cont)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	m.checkMicrosoftJWTToken(token, res)
	m.finishContestantLogin(cont, res)
}

func (m *MicrosoftLoginAPI) finishContestantLogin(cont0 models.Contestant, res http.ResponseWriter) {
	cont, err := m.contRepo.GetByEmail(cont0.Email)
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
		err = m.contRepo.Add(&cont)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	m.finishLogin(cont.ID, res)
}

// organizer stuff

func (m *MicrosoftLoginAPI) handleOrganizerLogin(res http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Authorization")

	var org models.Organizer
	err := json.NewDecoder(req.Body).Decode(&org)
	_ = req.Body.Close()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if m.checkMicrosoftJWTToken(token, res) {
		m.finishOrganizerLogin(org, res)
	}
}

func (m *MicrosoftLoginAPI) finishOrganizerLogin(org0 models.Organizer, res http.ResponseWriter) {
	org, err := m.orgRepo.GetByEmail(org0.Email)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	m.finishLogin(org.ID, res)
}

// Microsoft JWT validation stuff
// 😉e10

func (m *MicrosoftLoginAPI) checkMicrosoftJWTToken(token string, res http.ResponseWriter) bool {
	return true
}

func (m *MicrosoftLoginAPI) validateMicrosoftJWT(token string) error {
	return nil
}
