package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	// "net/http"
	// "net/http/httptest"
	"time"
)

const (
	dbNAME     = "healthChecker"
	dbUSER     = "root"
	dbPASSWORD = "ks"
)

var db *sql.DB
var router *gin.Engine

type healthCheckLog struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	URLID     uint
}

type urlRecord struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	URL              string `json:"url"`
	CrawlTimeOut     int    `json:"crawlTimeOut"`
	Frequency        int    `json:"frequency"`
	FailureThreshold int    `json:"failureThreshold"`
	Status           int    `json:"status"`
}

type transformedURLRecord struct {
	URL              string `json:"url"`
	CrawlTimeOut     int    `json:"crawlTimeOut"`
	Frequency        int    `json:"frequency"`
	FailureThreshold int    `json:"failureThreshold"`
	Status           int    `json:"status"`
}

//------------------------------ dbInit funtion --------------------

func dbInit() {
	var err error
	db, err = sql.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("successfully connected")
	}

	stmt, err := db.Prepare(`CREATE Table IF NOT EXISTS urlRecords(
		id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		url varchar(255) UNIQUE NOT NULL,
		crawlTimeOut int NOT NULL,
		frequency int NOT NULL,
		failureThreshold int NOT NULL,
		status int DEFAULT 500,
		created_at DATETIME,
		updated_at DATETIME
		);`)

	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

	stmt, err = db.Prepare(`CREATE Table IF NOT EXISTS healthCheckLogs(
		ID int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		url_id int UNSIGNED NOT NULL,
		trial_number int,
		response int,
		create_at DATETIME,
		FOREIGN KEY fk_url_record(url_id)
		REFERENCES urlRecords(id)
		);`)

	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

}

func main() {
	//--------------------- database initialisation ----------------
	var data []urlRecord
	dbInit()

	//--------------------- unmarshalling json file ----------------

	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		panic("failed to read file")
	}
	// fmt.Println(string(file))
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		panic("failed to unmarshall file")
	}
	// fmt.Printf("%#v", data)

	//----------------- Adding unmarshall data in database ---------

	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range data {
		_, err := db.Exec(`REPLACE INTO urlRecords(
			url ,
			crawlTimeOut ,
			frequency ,
			failureThreshold,
			created_at,
			updated_at)
			VALUES (?,?,?,?,?,?)
			;`,
			v.URL,
			v.CrawlTimeOut,
			v.Frequency,
			v.FailureThreshold,
			time.Now(),
			time.Now())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("inserted record successfully..")
		}
	}

	// ------------- setting up routes using gin -------------------
	router = gin.Default()

	app := router.Group("api/healthcheck")
	{
		app.GET("/check", checkURLHealth)

	}

	router.Run(":80")
}

type url struct {
	Name string `form:"url"`
}

// func checkURLHealth(c *gin.Context) {
// 	// id := c.Query("id")
// 	// pk := c.Param("pk")
// 	// name := c.DefaultQuery("name", "john")
// 	// // fmt.Printf("%v ID %v : corresponds to %v \n", pk, id, name)
// 	// // c.String(http.StatusOK, "%v ID %v : corresponds to %v \n", pk, id, name)
// 	// c.JSON(http.StatusOK, gin.H{
// 	// 	"pk":   pk,
// 	// 	"name": name,
// 	// 	"id":   id,
// 	// })
// 	var obj url
// 	if c.ShouldBindQuery(&obj) == nil {
// 		fmt.Println("====== Only Bind By Query String ======")
// 		fmt.Println(obj)
// 		resp, err := http.Get(obj.Name)
// 		if err != nil {
// 			fmt.Println("error is ", err)
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"status": http.StatusNotFound,
// 				"error":  err,
// 			})
// 			return
// 		}
// 		defer resp.Body.Close()
// 		c.JSON(http.StatusOK, gin.H{
// 			"status": http.StatusOK,
// 			"error":  err,
// 		})
// 	}

// }

func checkURLHealth(c *gin.Context) {
	var urlinfo transformedURLRecord
	results, err := db.Query(`SELECT 
							url,
							crawlTimeOut,
							frequency ,
							failureThreshold
							FROM urlRecords`)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("results fetched successfully")
	}
	defer results.Close()

	for results.Next() {
		err = results.Scan(&urlinfo.URL, &urlinfo.CrawlTimeOut, &urlinfo.Frequency, &urlinfo.FailureThreshold)
		if err != nil {
			fmt.Println("Error while scanning row")
		} else {
			b, err := json.MarshalIndent(urlinfo, "", "   ")
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			fmt.Printf("%s\n", b)
		}
	}
	err = results.Err()
	if err != nil {
		fmt.Println("Error after ending iteration on result set")
	}

}
