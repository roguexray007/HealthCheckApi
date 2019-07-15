package service

import (
	"sync"
	"time"
)

// Wg Wait group variable
var Wg sync.WaitGroup

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
