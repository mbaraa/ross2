package main

import (
	"github.com/mbaraa/ross2/router"

	"github.com/mbaraa/ross2/data/db"
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
	)

	router.New(
		contestRepo, sessionRepo, contestantRepo, teamRepo, organizerRepo, joinReqRepo, notificationRepo,
	).Start()
}
