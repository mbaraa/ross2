package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/helpers"
	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

type Builder struct {
	contestRepo      data.Many2ManyCRUDRepo[models.Contest, any]
	contestantRepo   data.ContestantCRUDRepo
	sessionRepo      data.SessionCRUDRepo
	teamRepo         data.TeamCRUDRepo
	organizerRepo    data.OrganizerCRUDRepo
	joinReqRepo      data.JoinRequestCRDRepo
	notificationRepo data.NotificationCRUDRepo
	userRepo         data.UserCRUDRepo
	adminRepo        data.AdminCRUDRepo
	ocRepo           data.OrganizeContestCRUDRepo
}

func NewRouterBuilder() *Builder {
	return new(Builder)
}

func (b *Builder) ContestRepo(c data.Many2ManyCRUDRepo[models.Contest, any]) *Builder {
	b.contestRepo = c
	return b
}

func (b *Builder) ContestantRepo(c data.ContestantCRUDRepo) *Builder {
	b.contestantRepo = c
	return b
}

func (b *Builder) SessionRepo(s data.SessionCRUDRepo) *Builder {
	b.sessionRepo = s
	return b
}

func (b *Builder) TeamRepo(t data.TeamCRUDRepo) *Builder {
	b.teamRepo = t
	return b
}

func (b *Builder) OrganizerRepo(o data.OrganizerCRUDRepo) *Builder {
	b.organizerRepo = o
	return b
}

func (b *Builder) JoinReqRepo(j data.JoinRequestCRDRepo) *Builder {
	b.joinReqRepo = j
	return b
}

func (b *Builder) NotificationRepo(n data.NotificationCRUDRepo) *Builder {
	b.notificationRepo = n
	return b
}

func (b *Builder) UserRepo(u data.UserCRUDRepo) *Builder {
	b.userRepo = u
	return b
}

func (b *Builder) AdminRepo(a data.AdminCRUDRepo) *Builder {
	b.adminRepo = a
	return b
}

func (b *Builder) OrganizeContestRepo(c data.OrganizeContestCRUDRepo) *Builder {
	b.ocRepo = c
	return b
}

func (b *Builder) verify() bool {
	sb := new(strings.Builder)

	if b.contestRepo == nil {
		sb.WriteString("Router Builder: missing contest repo!")
	}
	if b.contestantRepo == nil {
		sb.WriteString("Router Builder: missing contestant repo!")
	}
	if b.sessionRepo == nil {
		sb.WriteString("Router Builder: missing session repo!")
	}
	if b.teamRepo == nil {
		sb.WriteString("Router Builder: missing team repo!")
	}
	if b.organizerRepo == nil {
		sb.WriteString("Router: missing organizer repo!")
	}
	if b.joinReqRepo == nil {
		sb.WriteString("Router Builder: missing join request repo!")
	}
	if b.notificationRepo == nil {
		sb.WriteString("Router Builder: missing notification repo!")
	}
	if b.userRepo == nil {
		sb.WriteString("Router Builder: missing user repo!")
	}
	if b.adminRepo == nil {
		sb.WriteString("Router Builder: missing admin repo!")
	}
	if b.ocRepo == nil {
		sb.WriteString("Router Builder: missing organize contest repo!")
	}

	if sb.Len() != 0 {
		fmt.Println(sb.String())
		return false
	}

	return true
}

func (b *Builder) GetRouter() *Router {
	if !b.verify() {
		return nil
	}
	return NewRouter(b)
}

type Router struct {
	contestAPI        *controllers.ContestAPI
	contestantAPI     *controllers.ContestantAPI
	orgAPI            *controllers.OrganizerAPI
	adminAPI          *controllers.AdminAPI
	notificationAPI   *controllers.NotificationAPI
	googleLoginAPI    *auth.OAuthLoginAPI
	microsoftLoginAPI *auth.OAuthLoginAPI
}

func (r *Router) verifyAPIs() bool {
	sb := new(strings.Builder)
	if r.contestAPI == nil {
		sb.WriteString("Router: missing contest API!")
	}
	if r.contestantAPI == nil {
		sb.WriteString("Router: missing contestant API!")
	}
	if r.orgAPI == nil {
		sb.WriteString("Router: missing organizer API!")
	}
	if r.adminAPI == nil {
		sb.WriteString("Router: missing admin API!")
	}
	if r.notificationAPI == nil {
		sb.WriteString("Router Builder: missing notification API!")
	}
	if r.googleLoginAPI == nil {
		sb.WriteString("Router: missing google login API!")
	}
	if r.microsoftLoginAPI == nil {
		sb.WriteString("Router: missing microsoft login API!")
	}

	if sb.Len() != 0 {
		fmt.Println(sb.String())
		return false
	}

	return true
}

func NewRouter(b *Builder) *Router {
	var (
		teamManager         = helpers.NewTeamHelper(b.teamRepo, b.contestantRepo)
		notificationManager = helpers.NewNotificationHelper(b.notificationRepo)

		organizerManager = helpers.NewOrganizerHelperBuilder().
					OrganizerRepo(b.organizerRepo).
					ContestRepo(b.contestRepo).
					UserRepo(b.userRepo).
					TeamMgr(teamManager).
					NotificationMgr(notificationManager).
					OrganizeContestRepo(b.ocRepo).
					GetOrganizerManager()

		joinReqManager    = helpers.NewJoinRequestHelper(b.joinReqRepo, b.notificationRepo, b.contestRepo, teamManager)
		sessionManager    = helpers.NewSessionHelper(b.sessionRepo)
		userManager       = helpers.NewUserHelper(b.userRepo, b.contestantRepo, sessionManager)
		contestantManager = helpers.NewContestantHelperBuilder().
					UserRepo(b.userRepo).
					ContestantRepo(b.contestantRepo).
					ContestRepo(b.contestRepo).
					NotificationRepo(b.notificationRepo).
					TeamMgr(teamManager).
					JoinRequestMgr(joinReqManager).
					GetContestantManager()

		adminHelper   = helpers.NewAdminHelper(b.adminRepo, b.organizerRepo, b.userRepo)
		authenticator = auth.NewHandlerAuthenticator(sessionManager)
	)

	r := &Router{
		contestAPI:        controllers.NewContestAPI(b.contestRepo),
		contestantAPI:     controllers.NewContestantAPI(contestantManager, authenticator),
		orgAPI:            controllers.NewOrganizerAPI(organizerManager, authenticator),
		adminAPI:          controllers.NewAdminAPI(adminHelper, authenticator),
		notificationAPI:   controllers.NewNotificationAPI(notificationManager, authenticator),
		googleLoginAPI:    auth.NewOAuthLoginAPI(userManager, auth.NewGoogleJWTTokenValidator(), "/gauth"),
		microsoftLoginAPI: auth.NewOAuthLoginAPI(userManager, auth.NewMicrosoftJWTValidator(), "/msauth"),
	}

	if !r.verifyAPIs() {
		return nil
	}
	return r
}

func (r *Router) getHandler() *http.ServeMux {
	handler := http.NewServeMux()
	handler.Handle("/contest/", r.contestAPI)
	handler.Handle("/contestant/", r.contestantAPI)
	handler.Handle("/organizer/", r.orgAPI)
	handler.Handle("/admin/", r.adminAPI)
	handler.Handle("/notification/", r.notificationAPI)
	handler.Handle("/gauth/", r.googleLoginAPI)
	handler.Handle("/msauth/", r.microsoftLoginAPI)
	handler.Handle("/", http.FileServer(http.Dir("./client/dist/")))

	return handler
}

func (r *Router) Start() {
	log.Printf("starting server at:\nhttp://localhost:%s\n%s",
		config.GetInstance().PortNumber, config.GetInstance().MachineAddress)
	log.Fatalln(http.ListenAndServe(
		":"+config.GetInstance().PortNumber,
		r.getHandler(),
	))
}
