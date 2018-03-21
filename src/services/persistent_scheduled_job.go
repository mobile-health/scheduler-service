package services

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/stores"
)

func Now() *time.Time {
	u := time.Now().UTC()
	return &u
}

type PersistentScheduledJob struct {
	stores.Store
	*models.ScheduledJob
	ParentJob *models.Job
}

func (scheduledJob *PersistentScheduledJob) ScheduledAt() time.Time {
	return scheduledJob.ScheduledJob.ScheduledAt
}

func (scheduledJob *PersistentScheduledJob) Save() {
	if apperr := scheduledJob.Store.ScheduledJob().Insert(scheduledJob.ScheduledJob); apperr != nil {
		log4go.Error(apperr.Message)
	}
}

func (scheduledJob *PersistentScheduledJob) onProcessing() {
	log4go.Info("The scheduled job %s is in processing", scheduledJob.ID)

	scheduledJob.Status = models.JobProcessing
	scheduledJob.RanAt = Now()
	if apperr := scheduledJob.Store.ScheduledJob().Update(scheduledJob.ScheduledJob); apperr != nil {
		log4go.Error(apperr.Message)
	}

	scheduledJob.ParentJob.JobStats.LastRanAt = Now()
	if apperr := scheduledJob.Store.Job().Update(scheduledJob.ParentJob); apperr != nil {
		log4go.Error(apperr.Message)
	}
}

func (scheduledJob *PersistentScheduledJob) onFailed(err error) {
	log4go.Info("The scheduled job %s failed", scheduledJob.ID)

	scheduledJob.Status = models.JobFailed
	scheduledJob.Error = err
	if apperr := scheduledJob.Store.ScheduledJob().Update(scheduledJob.ScheduledJob); apperr != nil {
		log4go.Error(apperr.Message)
	}

	scheduledJob.ParentJob.JobStats.ErrorCount++
	scheduledJob.ParentJob.JobStats.LastError = err.Error()
	scheduledJob.ParentJob.JobStats.LastErrorAt = Now()
	if apperr := scheduledJob.Store.Job().Update(scheduledJob.ParentJob); apperr != nil {
		log4go.Error(apperr.Message)
	}
}

func (scheduledJob *PersistentScheduledJob) onSucceeded() {
	log4go.Info("The scheduled job %s has been succeeded", scheduledJob.ID)

	scheduledJob.Status = models.JobSucceeded
	scheduledJob.ParentJob.JobStats.SuccessCount++
	scheduledJob.ParentJob.JobStats.LastSuccededAt = Now()

	if apperr := scheduledJob.Store.ScheduledJob().Update(scheduledJob.ScheduledJob); apperr != nil {
		log4go.Error(apperr.Message)
	}
}

func (scheduledJob *PersistentScheduledJob) Run() (err error) {
	log4go.Info("Process scheduled job %s", scheduledJob.ID)

	scheduledJob.onProcessing()

	if !scheduledJob.ParentJob.IsAsync {
		defer func(err error) {
			if err != nil {
				scheduledJob.onFailed(err)
			}
		}(err)
	}

	req, err := http.NewRequest(scheduledJob.ParentJob.Args.Method, scheduledJob.ParentJob.Args.URL, strings.NewReader(scheduledJob.ParentJob.Args.Body))
	if err != nil {
		return err
	}

	for value, key := range scheduledJob.ParentJob.Args.Headers {
		req.Header.Add(value, key)
	}

	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	if !scheduledJob.ParentJob.IsAsync {
		scheduledJob.onSucceeded()
	}

	return err
}

func NewPersistentScheduledJob(store stores.Store, job *models.Job) *PersistentScheduledJob {
	scheduledJob := models.ScheduledJob{
		JobID:       job.ID,
		ScheduledAt: job.NextRunAt,
		Status:      models.JobPending,
	}
	persistentJob := &PersistentScheduledJob{
		ScheduledJob: &scheduledJob,
		Store:        store,
		ParentJob:    job,
	}
	return persistentJob
}
