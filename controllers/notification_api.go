package controllers

import (
	"encoding/json"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
	"net/http"
	"strings"
)

type NotificationAPI struct {
	endPoints map[string]http.HandlerFunc

	notiRepo data.NotificationCRUDRepo
	sessMgr  *managers.SessionManager
	contMgr  *managers.ContestantManager
}

func NewNotificationAPI(notificationRepo data.NotificationCRUDRepo, sessionManager *managers.SessionManager,
	contestantManager *managers.ContestantManager) *NotificationAPI {
	return (&NotificationAPI{
		notiRepo: notificationRepo,
		sessMgr:  sessionManager,
		contMgr:  contestantManager,
	}).initEndPoints()
}

func (n *NotificationAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/notification"), n.endPoints)
}

func (n *NotificationAPI) initEndPoints() *NotificationAPI {
	n.endPoints = map[string]http.HandlerFunc{
		"GET /all/": n.authenticateHandler(n.handleGetNotifications),
	}
	return n
}

func (n *NotificationAPI) authenticateHandler(h HandlerFuncWithSession) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, ok := n.sessMgr.CheckSessionFromRequest(req)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(res, req, session)
	}
}

// GET /notification/all/
func (n *NotificationAPI) handleGetNotifications(res http.ResponseWriter, req *http.Request, s models.Session) {
	cont, err := n.contMgr.GetContestant(s.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	notifications, err := n.notiRepo.GetAllForUser(cont.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(res).Encode(notifications)
}
