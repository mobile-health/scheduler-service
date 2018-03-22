package models

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gorhill/cronexpr"
)

// Job present a job
type Job struct {
	ID            string     `json:"id" bson:"_id"`
	FuID          string     `json:"fu_id" bson:"fu_id"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" bson:"updated_at"`
	Name          string     `json:"name" bson:"name"`
	Expression    string     `json:"expression" bson:"expression"`
	IsDone        bool       `json:"is_done" bson:"is_done"`
	NextRunAt     time.Time  `json:"next_run_at" bson:"next_run_at"`
	IsDisabled    bool       `json:"is_disabled" bson:"id_disabled"`
	IsAsync       bool       `json:"is_async" bson:"is_async"`
	Description   string     `json:"description" bson:"description"`
	MaxRetries    uint       `json:"max_retries" bson:"max_retries"`
	MaxExecutions uint       `json:"max_executions" bson:"max_executions"`
	Tags          []string   `json:"tags" bson:"tags"`
	JobStats      JobStats   `json:"job_stats" bson:"job_stats"`
	Args          RemoteArgs `json:"args" bson:"args"`
}

type Jobs []*Job

func (job *Job) PreInsert() *Job {
	job.Args.PreInsert()

	job.ID = NewID()
	job.CreatedAt = Now()
	job.UpdatedAt = Now()

	job.Name = strings.TrimSpace(job.Name)
	job.FuID = strings.TrimSpace(job.FuID)
	job.Expression = strings.TrimSpace(job.Expression)
	if len(job.FuID) > 0 {
		job.FuID = "fu_" + NewID()
	}

	return job
}

func (job *Job) PreUpdate() {
	job.UpdatedAt = Now()
}

func (job *Job) Validate() *Error {
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

	if len(job.Name) == 0 {
		errFields = append(errFields, NewErrorFieldRequired("name"))
	}

	if len(job.Expression) == 0 {
		errFields = append(errFields, NewErrorFieldRequired("expression"))
	}

	if expr, err := cronexpr.Parse(job.Expression); err != nil {
		errFields = append(errFields, NewErrorFieldInvalid("expression"))
	} else {
		nextTime := expr.Next(Now())
		if nextTime.IsZero() {
			errFields = append(errFields, NewErrorFieldInvalid("expression"))
		} else {
			job.NextRunAt = nextTime
		}
	}

	if apperr := job.Args.Validate(); apperr != nil {
		errFields = append(errFields, apperr.Errors...)
	}

	return errFields.GenAppError()
}

func NewJobFromBody(body io.ReadCloser) *Job {
	var job Job
	json.NewDecoder(body).Decode(&job)
	return &job
}

func (job *Job) Clone() *Job {
	clone := *job
	return &clone
}
