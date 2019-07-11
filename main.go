package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
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
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	Name             string     `json:"name"`
	CrawlTimeOut     int        `json:"crawlTimeOut"`
	Frequency        int        `json:"frequency"`
	FailureThreshold int        `json:"failureThreshold"`
	Status           int        `json:"status"`
}

//------------------------------ dbInit funtion --------------------

func dbInit() {
	var err error
	db, err = gorm.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("successfully connected")
	}
	db.AutoMigrate(&urlRecord{})
}

func main() {
	//--------------------- database initialisation ----------------
	var data []urlRecord
	var tempRec urlRecord
	dbInit()

	//--------------------- unmarshalling json file ----------------

	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic("failed to read file")
	}
	fmt.Println(string(file))
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		panic("failed to unmarshall file")
	}
	// fmt.Printf("%#v", data)

	//----------------- Adding unmarshall data in database ---------
	for _, v := range data {
		db.Where("name = ? ", v.Name).First(&tempRec)
		if tempRec == (urlRecord{}) {
			db.Create(&v)
		}
	}

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
