package main

import (
	"database/sql"
	"qaa/controllers"
	database "qaa/db"
	"qaa/services/answersService"
	"qaa/services/questionsService"
	"qaa/services/trainingsService"
	"qaa/services/usersService"
	"qaa/templates"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	templates.InitTemplates("./templates")

	db = database.InitDB()
	defer database.DB.Close()
	// Pass the db connection to alertsService
	questionsService.SetDB(db)
    answersService.SetDB(db)
	trainingsService.SetDB(db)
	usersService.SetDB(db)

	// start a new goroutine for the rest api endpoints
	controllers.RestApi()

}

//
//
