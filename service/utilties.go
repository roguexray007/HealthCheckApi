package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

//AddToDB add data to db
func AddToDB(c *gin.Context) {
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

// CheckHealth fetches records from db and calls CheckURLHealth for each record
func CheckHealth() {
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
			Wg.Add(1)
			go CheckURLHealth(&urlinfo)
			fmt.Println("go function called for ,", urlinfo.URL)
		}
	}
	err = results.Err()
	if err != nil {
		fmt.Println("Error after ending iteration on result set")
	}
	Wg.Wait()
	fmt.Printf("--------------------- HEALTH CHECK COMPLETED ---------------------------\n")

}

// CheckURLHealth checks health status for each url
func CheckURLHealth(urlinfo *urlRecord) {
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
	Wg.Done()
}

// FetchLogs fethces Logs fron healthCheckLog table
func FetchLogs(c *gin.Context) {
	var healthLogs []healthCheckLog

	results, err := db.Query(`SELECT * FROM healthCheckLogs`)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err,
			"status": http.StatusInternalServerError,
		})
		panic(err)
	} else {
		fmt.Println("health Check Logs fetched successfully")
	}
	defer results.Close()

	for results.Next() {
		var healthLogInfo healthCheckLog
		err = results.Scan(&healthLogInfo.ID,
			&healthLogInfo.URLID,
			&healthLogInfo.TrialNumber,
			&healthLogInfo.Response,
			&healthLogInfo.CreatedAt)

		if err != nil {
			fmt.Println("Error while scanning row")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  err,
				"status": http.StatusInternalServerError,
			})
			panic(err)
		} else {
			fmt.Println("result fetched successfully")
			healthLogs = append(healthLogs, healthLogInfo)
		}
	}
	err = results.Err()
	if err != nil {
		fmt.Println("Error after ending iteration on result set")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err,
			"status": http.StatusInternalServerError,
		})
		panic(err)
	} else {

		c.JSON(http.StatusOK, healthLogs)
	}

}
