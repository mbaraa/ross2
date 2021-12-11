package router

import (
	"fmt"
	"log"
	"net/http"

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

func (b *Builder) verify() bool {
	if b.contestRepo == nil {
		fmt.Println("Router Builder: missing contest repo!")
	}
	if b.contestantRepo == nil {
		fmt.Println("Router Builder: missing contestant repo!")
	}
	if b.sessionRepo == nil {
		fmt.Println("Router Builder: missing session repo!")
	}
	if b.teamRepo == nil {
		fmt.Println("Router Builder: missing team repo!")
	}
	if b.organizerRepo == nil {
		fmt.Println("Router: missing organizer repo!")
	}
	if b.joinReqRepo == nil {
		fmt.Println("Router Builder: missing join request repo!")
	}
	if b.notificationRepo == nil {
		fmt.Println("Router Builder: missing notification repo!")
	}

	return b.contestRepo != nil && b.contestantRepo != nil &&
		b.sessionRepo != nil && b.teamRepo != nil &&
		b.organizerRepo != nil && b.joinReqRepo != nil &&
		b.notificationRepo != nil
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
	googleLoginAPI    *auth.GoogleLoginAPI
	microsoftLoginAPI *auth.MicrosoftLoginAPI
}

func (r *Router) verifyAPIs() bool {
	if r.contestAPI == nil {
		fmt.Println("Router: missing contest API!")
	}
	if r.contestantAPI == nil {
		fmt.Println("Router: missing contestant API!")
	}
	if r.orgAPI == nil {
		fmt.Println("Router: missing organizer API!")
	}
	if r.notificationAPI == nil {
		fmt.Println("Router Builder: missing notification API!")
	}
	if r.googleLoginAPI == nil {
		fmt.Println("Router: missing google login API!")
	}
	if r.microsoftLoginAPI == nil {
		fmt.Println("Router: missing microsoft login API!")
	}

	return r.contestAPI != nil && r.contestantAPI != nil &&
		r.orgAPI != nil && r.notificationAPI != nil &&
		r.googleLoginAPI != nil && r.microsoftLoginAPI != nil
}

func NewRouter(b *Builder) *Router {
	var (
		sessionManager      = managers.NewSessionManager(b.sessionRepo)
		contestantManager   = managers.NewContestantManager(b.contestantRepo, sessionManager, b.contestRepo, b.teamRepo)
		teamManager         = managers.NewTeamManager(b.teamRepo, b.contestantRepo)
		organizerManager    = managers.NewOrganizerManager(b.organizerRepo, sessionManager, b.contestRepo)
		joinReqManager      = managers.NewJoinRequestManager(b.joinReqRepo, b.notificationRepo, b.contestRepo, teamManager)
		notificationManager = managers.NewNotificationManager(b.notificationRepo)
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

		googleLoginAPI:    auth.NewGoogleLoginAPI(sessionManager, b.contestantRepo, b.organizerRepo),
		microsoftLoginAPI: auth.NewMicrosoftLoginAPI(sessionManager, b.contestantRepo, b.organizerRepo),
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
