package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	// "net/http/httptest"
	"github.com/robfig/cron"
	"healthChecker/service"
)

var router *gin.Engine

//------------------------------ dbInit funtion --------------------

func main() {
	//--------------------- database initialisation ----------------
	service.DBInit()

	sched := cron.New()
	sched.AddFunc(service.REFRESHTIME, service.CheckHealth)
	sched.Start()

	// ------------- setting up routes using gin -------------------
	router = gin.Default()

	app := router.Group("api/healthcheck")
	{
		app.GET("/addToDB", service.AddToDB)
		app.GET("/fetchLogs", service.FetchLogs)
	}

	router.Run()

}
