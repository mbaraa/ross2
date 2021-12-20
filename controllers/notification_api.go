package controllers

import (
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/context"
	"github.com/mbaraa/ross2/controllers/helpers"
)

type NotificationAPI struct {
	notificationMgr *helpers.NotificationHelper
	authenticator   *auth.HandlerAuthenticator
	endPoints       map[string]http.HandlerFunc
}

func NewNotificationAPI(notificationsMgr *helpers.NotificationHelper, authenticator *auth.HandlerAuthenticator) *NotificationAPI {
	return (&NotificationAPI{
		notificationMgr: notificationsMgr,
		authenticator:   authenticator,
	}).initEndPoints()
}

func (n *NotificationAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	GetHandlerFromParentPrefix(res, req, strings.TrimPrefix(req.URL.Path, "/notification"), n.endPoints)
}

func (n *NotificationAPI) initEndPoints() *NotificationAPI {
	n.endPoints = n.authenticator.AuthenticateHandlers(map[string]auth.HandlerFunc{
		"GET /all/":   n.handleGetNotifications,
		"GET /check/": n.handleCheckNotifications,
		"GET /clear/": n.handleClearNotifications,
	})
	return n
}

// GET /notification/all/
func (n *NotificationAPI) handleGetNotifications(ctx context.HandlerContext) {
	nots, err := n.notificationMgr.GetNotifications(ctx.Sess)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
		return
	}
	_ = ctx.WriteJSON(nots, 0)
}

// GET /notification/check/
func (n *NotificationAPI) handleCheckNotifications(ctx context.HandlerContext) {
	if !n.notificationMgr.CheckNotifications(ctx.Sess) {
		http.Error(ctx.Res, "unauthorized", http.StatusUnauthorized)
		return
	}
	_ = ctx.WriteJSON(map[string]interface{}{
		"notifications_exists": true,
	}, 0)
}

// GET /notification/clear/
func (n *NotificationAPI) handleClearNotifications(ctx context.HandlerContext) {
	err := n.notificationMgr.ClearNotifications(ctx.Sess)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusUnauthorized)
	}
}
