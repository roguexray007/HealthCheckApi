package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

const (
	dbNAME     = "healthChecker"
	dbUSER     = "root"
	dbPASSWORD = "ks"
)

// func init() {
// 	var err error
// 	db, err = gorm.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("failed to connect to database")
// 	}
// }

func main() {
	router := gin.Default()

	app := router.Group("api/healthcheck")
	{
		app.GET("/", func(c *gin.Context) {

		})
	}

	router.Run()
}
