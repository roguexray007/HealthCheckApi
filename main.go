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

	//--------------------- Cron for scheduling health check -------
	sched := cron.New()
	sched.AddFunc(service.REFRESHTIME, service.CheckHealth)
	sched.Start()

	// ------------- setting up routes using gin --------------------
	router = gin.Default()

	//---------------------- API routes -----------------------------
	app := router.Group("api/healthcheck")
	{
		//---------- add/update data in db --------------------------
		app.GET("/addToDB", service.AddToDB)

		//---------- fetches logs and displays in JSON format -------
		app.GET("/fetchLogs", service.FetchLogs)
	}

	router.Run()

}
