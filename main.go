package main

import (
	"github.com/mbaraa/ross2/data/db"
	"github.com/mbaraa/ross2/router"
)

func main() {
	var (
		mysqlDB = db.
			GetDBManagerInstance().
			GetMySQLConn()
		contestantRepo   = db.NewContestantDB(mysqlDB)
		teamRepo         = db.NewTeamDB(mysqlDB, contestantRepo)
		contestRepo      = db.NewContestDB(mysqlDB, teamRepo, contestantRepo)
		organizerRepo    = db.NewOrganizerDB(mysqlDB)
		sessionRepo      = db.NewSessionDB(mysqlDB)
		notificationRepo = db.NewNotificationDB(mysqlDB)
		joinReqRepo      = db.NewJoinRequestDB(mysqlDB)
		userRepo         = db.NewUserDB(mysqlDB)
	)

	db.GetDBManagerInstance().InitTables()
	r := router.NewRouterBuilder().
		ContestRepo(contestRepo).
		ContestantRepo(contestantRepo).
		SessionRepo(sessionRepo).
		TeamRepo(teamRepo).
		OrganizerRepo(organizerRepo).
		JoinReqRepo(joinReqRepo).
		NotificationRepo(notificationRepo).
		UserRepo(userRepo).
		GetRouter()

	if r != nil {
		r.Start()
	}
}
