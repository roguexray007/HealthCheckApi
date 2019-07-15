package service

import (
	"sync"
	"time"
)

const (
	// REFRESHTIME Refresh Time for cron scheduler
	REFRESHTIME = "*/20 * * * *"
)

// Wg Wait group variable
var Wg sync.WaitGroup

type healthCheckLog struct {
	ID          uint      `json:"id"`
	URLID       uint      `json:"url_id"`
	TrialNumber int       `json:"trial_number"`
	Response    int       `json:"response"`
	CreatedAt   time.Time `json:"created_at"`
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
