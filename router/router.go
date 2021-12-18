package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
)

type Builder struct {
	contestRepo      data.ContestCRUDRepo
	contestantRepo   data.ContestantCRUDRepo
	sessionRepo      data.SessionCRUDRepo
	teamRepo         data.TeamCRUDRepo
	organizerRepo    data.OrganizerCRUDRepo
	joinReqRepo      data.JoinRequestCRDRepo
	notificationRepo data.NotificationCRUDRepo
	userRepo         data.UserCRUDRepo
}

func NewRouterBuilder() *Builder {
	return new(Builder)
}

func (b *Builder) ContestRepo(c data.ContestCRUDRepo) *Builder {
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
		sessionManager      = managers.NewSessionManager(b.sessionRepo)
		contestantManager   = managers.NewContestantManager(b.contestantRepo, sessionManager, b.contestRepo, b.teamRepo)
		teamManager         = managers.NewTeamManager(b.teamRepo, b.contestantRepo)
		organizerManager    = managers.NewOrganizerManager(b.organizerRepo, sessionManager, b.contestRepo)
		joinReqManager      = managers.NewJoinRequestManager(b.joinReqRepo, b.notificationRepo, b.contestRepo, teamManager)
		notificationManager = managers.NewNotificationManager(b.notificationRepo)
		userManager         = managers.NewUserManager(b.userRepo, sessionManager)
	)

	r := &Router{
		contestAPI: controllers.NewContestAPI(b.contestRepo),
		contestantAPI: controllers.NewContestantAPIBuilder().
			ContestantMgr(contestantManager).
			SessionMgr(sessionManager).
			TeamMgr(teamManager).
			JoinReqMgr(joinReqManager).
			GetContestantAPI(),

		orgAPI: controllers.NewOrganizerAPIBuilder().
			OrganizerMgr(organizerManager).
			SessionMgr(sessionManager).
			TeamMgr(teamManager).
			ContestantMgr(contestantManager).
			ContestRepo(b.contestRepo).
			NotificationMgr(notificationManager).
			GetOrganizerAPI(),

		notificationAPI: controllers.NewNotificationAPIBuilder().
			NotificationRepo(b.notificationRepo).
			SessionMgr(sessionManager).
			ContestantMgr(contestantManager).
			GetNotificationAPI(),

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
