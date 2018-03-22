package models

import (
	"time"
)

// JobStats for storing job metrics
type JobStats struct {
	ScheduleCount  uint       `json:"schedule_count"`
	SuccessCount   uint       `json:"success_count" bson:"success_count"`
	ErrorCount     uint       `json:"error_count" bson:"error_count"`
	LastRanAt      *time.Time `json:"last_ran_at" bson:"last_ran_at"`
	LastSuccededAt *time.Time `json:"last_succeded_at" bson:"last_succeded_at"`
	LastError      string     `json:"last_error" bson:"last_error"`
	LastErrorAt    *time.Time `json:"last_error_at" bson:"last_error_at"`
}
