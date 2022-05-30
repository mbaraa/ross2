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
	contestRepo      data.CRUDRepo[models.Contest]
	contestantRepo   data.CRUDRepo[models.Contestant]
	sessionRepo      data.CRUDRepo[models.Session]
	teamRepo         data.CRUDRepo[models.Team]
	organizerRepo    data.CRUDRepo[models.Organizer]
	joinReqRepo      data.CRUDRepo[models.JoinRequest]
	notificationRepo data.CRUDRepo[models.Notification]
	userRepo         data.CRUDRepo[models.User]
	adminRepo        data.CRUDRepo[models.Admin]
	ocRepo           data.OrganizeContestCRUDRepo
	rtRepo           data.RegisterTeamCRUDRepo
}

func NewRouterBuilder() *Builder {
	return new(Builder)
}

func (b *Builder) ContestRepo(c data.CRUDRepo[models.Contest]) *Builder {
	b.contestRepo = c
	return b
}

func (b *Builder) ContestantRepo(c data.CRUDRepo[models.Contestant]) *Builder {
	b.contestantRepo = c
	return b
}

func (b *Builder) SessionRepo(s data.CRUDRepo[models.Session]) *Builder {
	b.sessionRepo = s
	return b
}

func (b *Builder) TeamRepo(t data.CRUDRepo[models.Team]) *Builder {
	b.teamRepo = t
	return b
}

func (b *Builder) OrganizerRepo(o data.CRUDRepo[models.Organizer]) *Builder {
	b.organizerRepo = o
	return b
}

func (b *Builder) JoinReqRepo(j data.CRUDRepo[models.JoinRequest]) *Builder {
	b.joinReqRepo = j
	return b
}

func (b *Builder) NotificationRepo(n data.CRUDRepo[models.Notification]) *Builder {
	b.notificationRepo = n
	return b
}

func (b *Builder) UserRepo(u data.CRUDRepo[models.User]) *Builder {
	b.userRepo = u
	return b
}

func (b *Builder) AdminRepo(a data.CRUDRepo[models.Admin]) *Builder {
	b.adminRepo = a
	return b
}

func (b *Builder) OrganizeContestRepo(c data.OrganizeContestCRUDRepo) *Builder {
	b.ocRepo = c
	return b
}

func (b *Builder) RegisterTeamRepo(rt data.RegisterTeamCRUDRepo) *Builder {
	b.rtRepo = rt
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
		sb.WriteString("Router Builder: missing session repo!\n")
	}
	if b.teamRepo == nil {
		sb.WriteString("Router Builder: missing team repo!\n")
	}
	if b.organizerRepo == nil {
		sb.WriteString("Router: missing organizer repo!\n")
	}
	if b.joinReqRepo == nil {
		sb.WriteString("Router Builder: missing join request repo!\n")
	}
	if b.notificationRepo == nil {
		sb.WriteString("Router Builder: missing notification repo!\n")
	}
	if b.userRepo == nil {
		sb.WriteString("Router Builder: missing user repo!\n")
	}
	if b.adminRepo == nil {
		sb.WriteString("Router Builder: missing admin repo!\n")
	}
	if b.ocRepo == nil {
		sb.WriteString("Router Builder: missing organize contest repo!\n")
	}
	if b.rtRepo == nil {
		sb.WriteString("Router Builer: missing regiser team repo!\n")
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
		teamManager         = helpers.NewTeamHelper(b.teamRepo, b.contestantRepo, b.rtRepo)
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
					RegisterTeamRepo(b.rtRepo).
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
	handler.Handle("/", http.FileServer(http.Dir(config.GetInstance().UploadDirectory)))

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
