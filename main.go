package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"sync"
	// "net/http/httptest"
	"github.com/robfig/cron"
	"time"
)

const (
	dbNAME     = "healthChecker"
	dbUSER     = "root"
	dbPASSWORD = "ks"
)

var db *sql.DB
var router *gin.Engine
var wg sync.WaitGroup

type healthCheckLog struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	URLID     uint
}

type urlRecord struct {
	ID               uint      `json:"-"`
	URL              string    `json:"url"`
	CrawlTimeOut     int       `json:"crawlTimeOut"`
	Frequency        int       `json:"frequency"`
	FailureThreshold int       `json:"failureThreshold"`
	Status           int       `json:"status"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

type transformedURLRecord struct {
	ID               uint   `json:"-"`
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
		status int DEFAULT 503,
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
		id int UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		url_id int UNSIGNED NOT NULL,
		trial_number int,
		response int,
		created_at DATETIME,
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
	dbInit()

	sched := cron.New()
	sched.AddFunc("*/2 * * * *", checkHealth)
	sched.Start()

	// ------------- setting up routes using gin -------------------
	router = gin.Default()

	app := router.Group("api/healthcheck")
	{
		app.GET("/addToDB", addToDB)

	}

	router.Run()

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

func addToDB(c *gin.Context) {
	var data []urlRecord

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
		result := db.QueryRow(`SELECT * FROM urlRecords WHERE url = ? limit 1`, v.URL)

		var urlinfo urlRecord
		err = result.Scan(&urlinfo.ID, &urlinfo.URL, &urlinfo.CrawlTimeOut, &urlinfo.Frequency,
			&urlinfo.FailureThreshold, &urlinfo.Status, &urlinfo.CreatedAt,
			&urlinfo.UpdatedAt)
		if err == sql.ErrNoRows {
			fmt.Println("Error no rows were returned .Therefore inserting new record")
			_, err = db.Exec(`INSERT INTO urlRecords(
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
		} else if err == nil {
			_, err = db.Exec(`UPDATE urlRecords SET
				url = ?,
				crawlTimeOut = ?,
				frequency =?,
				failureThreshold = ?,
				updated_at = ?
			WHERE 
				id = ?
				;`,
				v.URL,
				v.CrawlTimeOut,
				v.Frequency,
				v.FailureThreshold,
				time.Now(),
				urlinfo.ID)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("updated record successfully..")
			}
		}
	}

}

func checkHealth() {
	results, err := db.Query(`SELECT * FROM urlRecords`)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("results fetched successfully")
	}
	defer results.Close()

	for results.Next() {
		var urlinfo urlRecord

		err = results.Scan(&urlinfo.ID, &urlinfo.URL, &urlinfo.CrawlTimeOut, &urlinfo.Frequency,
			&urlinfo.FailureThreshold, &urlinfo.Status, &urlinfo.CreatedAt,
			&urlinfo.UpdatedAt)
		if err != nil {
			fmt.Println("Error while scanning row")
		} else {
			urlinfoJSON, err := json.MarshalIndent(urlinfo, "", "   ")
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			fmt.Printf("%s\n", urlinfoJSON)
			fmt.Println(urlinfo)
			wg.Add(1)
			go checkURLHealth(&urlinfo)
			fmt.Println("go function called for ,", urlinfo.URL)
		}
	}
	err = results.Err()
	if err != nil {
		fmt.Println("Error after ending iteration on result set")
	}
	wg.Wait()
	fmt.Printf("--------------------- HEALTH CHECK COMPLETED ---------------------------\n")

}

func checkURLHealth(urlinfo *urlRecord) {
	timeout := time.Duration(time.Duration((*urlinfo).CrawlTimeOut) * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	for trial := 1; trial <= (*urlinfo).FailureThreshold; trial++ {
		_, err := client.Get((*urlinfo).URL)
		if err != nil {
			(*urlinfo).Status = http.StatusServiceUnavailable
			_, err1 := db.Exec(`INSERT INTO healthCheckLogs(
				url_id ,
				trial_number,
				response ,
				created_at )
				VALUES (?,?,?,?)
				;`,
				(*urlinfo).ID,
				trial,
				(*urlinfo).Status,
				time.Now())
			if err1 != nil {
				fmt.Println("error : ", err.Error())
			} else {
				fmt.Println("inserted record successfully..")
			}

			time.Sleep(time.Duration((*urlinfo).Frequency) * time.Millisecond)
		} else {
			(*urlinfo).Status = http.StatusOK
			_, err1 := db.Exec(`INSERT INTO healthCheckLogs(
				url_id ,
				trial_number,
				response ,
				created_at )
				VALUES (?,?,?,?)
				;`,
				(*urlinfo).ID,
				trial,
				(*urlinfo).Status,
				time.Now())
			if err1 != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("inserted record successfully..")
			}
			break
		}
	}
	wg.Done()
}
