package auth

import (
	"net/http"

	"github.com/mbaraa/ross2/controllers/context"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/models"
)

// HandlerFunc is a handler function with an extra parameter(session)
type HandlerFunc func(ctx context.HandlerContext)

// HandlerAuthenticator is responsible for authenticating handlers
type HandlerAuthenticator struct {
	sessMgr *helpers.SessionHelper[models.Session]
}

// NewHandlerAuthenticator returns a new HandlerAuthenticator instance
func NewHandlerAuthenticator(sessionManager *helpers.SessionHelper[models.Session]) *HandlerAuthenticator {
	return &HandlerAuthenticator{sessMgr: sessionManager}
}

// AuthenticateHandler executes the given handler function is authentication is done correctly :)
func (h *HandlerAuthenticator) AuthenticateHandler(handler HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		session, err := h.sessMgr.GetSession(req.Header.Get("Authorization"))
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(context.HandlerContext{Res: res, Req: req, Sess: session})
	}
}

// AuthenticateHandlers returns authenticated handler functions map of the given map
func (h *HandlerAuthenticator) AuthenticateHandlers(handlers map[string]HandlerFunc) map[string]http.HandlerFunc {
	authHandlers := make(map[string]http.HandlerFunc)
	for name, handler := range handlers {
		authHandlers[name] = h.AuthenticateHandler(handler)
	}
	return authHandlers
}
