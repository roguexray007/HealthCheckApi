package main

import (
	// "github.com/gin-gonic/gin"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

const (
	dbNAME     = "healthChecker"
	dbUSER     = "root"
	dbPASSWORD = "ks"
)

// type urlRecord struct {
// 	gorm.Model
// 	Name             string `json:"name"`
// 	CrawlTimeOut     int    `json:"crawlTimeOut`
// 	Frequency        int    `json:"frequency`
// 	FailureThreshold int    `json:"failureThreshold`
// 	Status           int    `json:"status"`
// }

func dbInit() {
	var err error
	db, err = gorm.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("successfully connected")
	}
}

func main() {
	dbInit()

	// router := gin.Default()

	// app := router.Group("api/healthcheck")
	// {
	// 	app.GET("/", func(c *gin.Context) {
	// 		c.JSON(200, gin.H{
	// 			"message": "hello world",
	// 		})
	// 	})
	// }

	// router.Run()
}
