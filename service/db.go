package service

import (
	"database/sql"
	"fmt"
)

const (
	dbNAME     = "healthChecker"
	dbUSER     = "root"
	dbPASSWORD = "ks"
)

var db *sql.DB

// DBInit initatlises connection and make tables
func DBInit() {
	var err error
	db, err = sql.Open("mysql", dbUSER+":"+dbPASSWORD+"@/"+dbNAME+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	} else {
		fmt.Println("successfully connected")
	}

	// Creation of urlRecords table for storing url info
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

	// Creation of healthCheckLogs table for storing logs

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
