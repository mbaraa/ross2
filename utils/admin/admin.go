package admin

import (
	"os"

	"github.com/mbaraa/ross2/data"
	"github.com/mbaraa/ross2/models"
)

func CreateAdmin(repo data.AdminCRUDRepo) {
	args := os.Args
	if len(args) >= 4 {
		if args[1] == "admin" {
			admin := models.Admin{
				User: models.User{
					Email: args[3],
				},
			}

			var err error
			switch args[2] {
			case "add":
				err = repo.Add(&admin)
			case "del":
				err = repo.Delete(admin)
			}

			if err != nil {
				panic(err)
			}
		}
		os.Exit(0)
	}
}
