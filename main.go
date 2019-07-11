package main

import (
	"encoding/json"
	"io/ioutil"
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

type urlRecord struct {
	gorm.Model
	urlInfo
}

type urlInfo struct {
	Name             string `json:"name"`
	CrawlTimeOut     int    `json:"crawlTimeOut"`
	Frequency        int    `json:"frequency"`
	FailureThreshold int    `json:"failureThreshold"`
	Status           int    `json:"status"`
}

//------------------------------ dbInit funtion --------------------

// func dbInit() {
// 	var err error
// 	db, err = gorm.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("failed to connect to database")
// 	} else {
// 		fmt.Println("successfully connected")
// 	}
// }

func main() {
	//--------------------- database initialisation ----------------
	// dbInit()

	//--------------------- unmarshalling json file ----------------

	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic("failed to read file")
	}
	fmt.Println(string(file))
	var data []urlInfo
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		panic("failed to unmarshall file")
	}
	fmt.Printf("%#v", data)

	// ------------- setting up routes using gin -------------------
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
