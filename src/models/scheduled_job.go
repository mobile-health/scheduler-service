package models

import (
	"time"
)

type JobStatus string

const (
	JobPending    JobStatus = "pending"
	JobSucceeded  JobStatus = "succeeded"
	JobProcessing JobStatus = "processing"
	JobFailed     JobStatus = "failed"
)

// ScheduledJob for storing the information when a job was executed.
type ScheduledJob struct {
	ID           string     `json:"id" bson:"_id"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at"`
	JobID        string     `json:"job_id" bson:"job_id"`
	ScheduledAt  time.Time  `json:"scheduled_at" bson:"scheduled_at"`
	RanAt        *time.Time `json:"ran_at" bson:"ran_at"`
	Error        error      `json:"error" bson:"error"`
	Duration     int        `json:"duration" bson:"duration"`
	NumOfRetries int        `json:"num_of_retries" bson:"num_of_retries"`
	Status       JobStatus  `json:"status" bson:"status"`
}

type ScheduledJobs []*ScheduledJob

func (job *ScheduledJob) PreInsert() {
	job.ID = NewID()
	job.CreatedAt = Now()
	job.UpdatedAt = Now()
	job.Status = JobPending
}

func (job *ScheduledJob) PreUpdate() {
	job.UpdatedAt = Now()
}

func (job *ScheduledJob) Validate() *Error {

	var errFields ErrorFields

	if len(job.ID) != 26 {
		errFields = append(errFields, NewErrorFieldInvalid("id"))
	}

	if job.CreatedAt.IsZero() {
		errFields = append(errFields, NewErrorFieldInvalid("created_at"))
	}

	if job.UpdatedAt.IsZero() {
		errFields = append(errFields, NewErrorFieldInvalid("updated_at"))
	}

	return errFields.GenAppError()
}
