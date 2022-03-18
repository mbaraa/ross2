package main

import (
	"github.com/mbaraa/ross2/models"
	"github.com/mbaraa/ross2/data/db"
	"github.com/mbaraa/ross2/router"
	"github.com/mbaraa/ross2/utils/admin"
)

func main() {
	var (
		mysqlDB = db.
			GetDBManagerInstance().
			GetMySQLConn()
		contestantRepo   = db.NewContestantDB[models.Contestant](mysqlDB)
		teamRepo         = db.NewTeamDB[models.Team, any](mysqlDB)
		contestRepo      = db.NewContestDB[models.Contest, any](mysqlDB, teamRepo)
		organizerRepo    = db.NewOrganizerDB(mysqlDB)
		sessionRepo      = db.NewSessionDB(mysqlDB)
		notificationRepo = db.NewNotificationDB(mysqlDB)
		joinReqRepo      = db.NewJoinRequestDB(mysqlDB)
		userRepo         = db.NewUserDB(mysqlDB)
		adminRepo        = db.NewAdminDB[models.Admin](mysqlDB)
		ocRepo           = db.NewOrganizeOrganizeContestDB(mysqlDB)
	)
	db.GetDBManagerInstance().InitTables()

	admin.CreateAdmin(adminRepo)

	r := router.NewRouterBuilder().
		ContestRepo(contestRepo).
		ContestantRepo(contestantRepo).
		SessionRepo(sessionRepo).
		TeamRepo(teamRepo).
		OrganizerRepo(organizerRepo).
		JoinReqRepo(joinReqRepo).
		NotificationRepo(notificationRepo).
		UserRepo(userRepo).
		AdminRepo(adminRepo).
		OrganizeContestRepo(ocRepo).
		GetRouter()

	if r != nil {
		r.Start()
	}
}
