package service

import (
	"time"
)

const (
	// REFRESHTIME Refresh Time for cron scheduler
	REFRESHTIME = "*/20 * * * *"
	//JSONDATA input file path for reading json data conatining url info
	JSONDATA = "data.json"
)



// table for storing logs
type healthCheckLog struct {
	ID          uint      `json:"id"`
	URLID       uint      `json:"url_id"`
	TrialNumber int       `json:"trial_number"`
	Response    int       `json:"response"`
	CreatedAt   time.Time `json:"created_at"`
}

// table for storing url info
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
