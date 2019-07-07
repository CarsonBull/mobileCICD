package models

import (
	gormdb "github.com/revel/modules/orm/gorm/app"
	"github.com/revel/revel"
)

func initializeDB() {

	// TODO: refactor this so we start the db as an env var that goes in the models package not using state in other packages
	// TODO: initialize db with parameters
	gormdb.InitDB()

	//gormdb.DB.Delete(&User{})
	//gormdb.DB.Delete(&Repo{})

	gormdb.DB.AutoMigrate(&User{})
	gormdb.DB.AutoMigrate(&Repo{})
	// TODO: Fix `[2019-07-06 20:40:18]  near "CONSTRAINT": syntax error`
	gormdb.DB.Model(&User{}).AddForeignKey("u_id", "users(user_id)","CASCADE", "CASCADE")
}

func closeDB() {
	gormdb.DB.Close()
}

func init() {
	revel.OnAppStart(initializeDB)
	revel.OnAppStop(closeDB)
}
