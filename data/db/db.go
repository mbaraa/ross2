package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mbaraa/ross2/config"
	"github.com/mbaraa/ross2/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbMan *dbManager = nil

type dbManager struct {
	mysqlConn *gorm.DB
	conf      *config.Config
}

func GetDBManagerInstance() *dbManager {
	if dbMan == nil {
		dbMan = &dbManager{
			conf:      config.GetInstance(),
			mysqlConn: nil,
		}
	}

	return dbMan
}

// GetMySQLConn returns a singleton MySQL connection instance
func (db *dbManager) GetMySQLConn() *gorm.DB {
	if db.mysqlConn == nil {
		var err error

		db.mysqlConn, err = gorm.Open(mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN: fmt.Sprintf("%s:%s@tcp(%s)/ross2?parseTime=True&loc=Local",
				config.GetInstance().DBUser, config.GetInstance().DBPassword, config.GetInstance().DBHost),
		}))
		if err != nil {
			panic(err)
		}
	}

	return db.mysqlConn
}

// InitTables creates the db's schemas w/o messing w/ the existing ones
// also also this method is only used with a fresh copy of ross 2, or when updated (LOL I don't think so)
func (db *dbManager) InitTables() {
	if db.mysqlConn != nil {
		err := db.mysqlConn.Debug().AutoMigrate(
			new(models.Contest),
			new(models.ParticipationConditions),
			new(models.Team),
			new(models.ContactInfo),
			new(models.User),
			new(models.Contestant),
			new(models.Organizer),
			new(models.Session),
			new(models.Notification),
			new(models.JoinRequest),
		)
		if err != nil {
			panic(err)
		}

		err = createNoTeam(db.mysqlConn)
		if err != nil {
			panic(err)
		}
	}
}

// needed for new registered contestants :)
func createNoTeam(db *gorm.DB) error {
	return db.
		Exec(`
INSERT INTO ross2.teams (name, leader_id)
SELECT *
FROM (SELECT 'no_team', 0) AS team
WHERE NOT EXISTS(
        SELECT name FROM ross2.teams WHERE name = 'no_team'
    )
LIMIT 1`).
		Error
}
