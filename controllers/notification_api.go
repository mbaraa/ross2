package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type NotificationAPIBuilder struct {
	notiRepo data.NotificationCRUDRepo
	sessMgr  *managers.SessionManager
	contMgr  *managers.ContestantManager
}

func NewNotificationAPIBuilder() *NotificationAPIBuilder {
	return new(NotificationAPIBuilder)
}

func (b *NotificationAPIBuilder) NotificationRepo(n data.NotificationCRUDRepo) *NotificationAPIBuilder {
	b.notiRepo = n
	return b
}

func (b *NotificationAPIBuilder) SessionMgr(s *managers.SessionManager) *NotificationAPIBuilder {
	b.sessMgr = s
	return b
}

func (b *NotificationAPIBuilder) ContestantMgr(c *managers.ContestantManager) *NotificationAPIBuilder {
	b.contMgr = c
	return b
}

func (b *NotificationAPIBuilder) verify() bool {
	if b.contMgr == nil {
		fmt.Println("Notification API Builder: missing contestant manager!")
	}
	if b.sessMgr == nil {
		fmt.Println("Notification API Builder: missing session manager!")
	}
	if b.notiRepo == nil {
		fmt.Println("Notification API Builder: missing notification repo!")
	}

	return b.contMgr != nil && b.sessMgr != nil &&
		b.notiRepo != nil
}

func (b *NotificationAPIBuilder) GetNotificationAPI() *NotificationAPI {
	if !b.verify() {
		return nil
	}
	return NewNotificationAPI(b)
}

type NotificationAPI struct {
	endPoints map[string]http.HandlerFunc

	notiRepo data.NotificationCRUDRepo
	sessMgr  *managers.SessionManager
	contMgr  *managers.ContestantManager
}

func NewNotificationAPI(b *NotificationAPIBuilder) *NotificationAPI {
	return (&NotificationAPI{
		notiRepo: b.notiRepo,
		sessMgr:  b.sessMgr,
		contMgr:  b.contMgr,
	}).initEndPoints()
}

func (n *NotificationAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/notification"), n.endPoints)
}

func (n *NotificationAPI) initEndPoints() *NotificationAPI {
	n.endPoints = map[string]http.HandlerFunc{
		"GET /all/":   n.authenticateHandler(n.handleGetNotifications),
		"GET /check/": n.authenticateHandler(n.handleCheckNotifications),
		"GET /clear/": n.authenticateHandler(n.handleClearNotifications),
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

// GET /notification/check/
func (n *NotificationAPI) handleCheckNotifications(res http.ResponseWriter, req *http.Request, s models.Session) {
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

	_, _ = res.Write([]byte(fmt.Sprint(len(notifications) != 0)))
}

// GET /notification/clear/
func (n *NotificationAPI) handleClearNotifications(res http.ResponseWriter, req *http.Request, s models.Session) {
	cont, err := n.contMgr.GetContestant(s.ID)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = n.notiRepo.DeleteAllForUser(cont.ID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
