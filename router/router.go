package router

import (
	"log"
	"net/http"

	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/controllers"
	"github.com/mbaraa/ross2/controllers/auth"
	"github.com/mbaraa/ross2/controllers/managers"
	"github.com/mbaraa/ross2/data"
)

type Router struct {
	contestAPI     *controllers.ContestAPI
	contestantAPI  *controllers.ContestantAPI
	orgAPI         *controllers.OrganizerAPI
	googleLoginAPI *auth.GoogleLoginAPI
}

func New(contestRepo data.ContestCRUDRepo, sessionRepo data.SessionCRUDRepo, contestantRepo data.ContestantCRUDRepo,
	teamRepo data.TeamCRUDRepo, organizerRepo data.OrganizerCRUDRepo, joinReqRepo data.JoinRequestCDRepo,
	notificationRepo data.NotificationCRUDRepo) *Router {

	var (
		sessionManager    = managers.NewSessionManager(sessionRepo)
		contestantManager = managers.NewContestantManager(contestantRepo, sessionManager, contestRepo)
		teamManager       = managers.NewTeamManager(teamRepo, contestantRepo)
		organizerManager  = managers.NewOrganizerManager(organizerRepo, sessionManager, contestRepo)
		joinReqManager    = managers.NewJoinRequestManager(joinReqRepo, notificationRepo, teamManager)
	)

	return &Router{
		contestAPI:     controllers.NewContestAPI(contestRepo),
		contestantAPI:  controllers.NewContestantAPI(contestantManager, sessionManager, teamManager, joinReqManager),
		orgAPI:         controllers.NewOrganizerAPI(organizerManager, sessionManager, teamManager, contestantManager),
		googleLoginAPI: auth.NewGoogleLoginAPI(sessionManager, contestantRepo, organizerRepo),
	}
}

func (r *Router) getHandler() *http.ServeMux {
	handler := http.NewServeMux()
	handler.Handle("/contest/", r.contestAPI)
	handler.Handle("/contestant/", r.contestantAPI)
	handler.Handle("/organizer/", r.orgAPI)
	handler.Handle("/gauth/", r.googleLoginAPI)

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
