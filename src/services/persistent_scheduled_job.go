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
	ParentJob *PersistentJob
}

func (scheduledJob *PersistentScheduledJob) ScheduledAt() time.Time {
	return scheduledJob.ScheduledJob.ScheduledAt
}

func (scheduledJob *PersistentScheduledJob) Save() {
	if len(scheduledJob.ScheduledJob.ID) != 26 {
		if apperr := scheduledJob.Store.ScheduledJob().Insert(scheduledJob.ScheduledJob); apperr != nil {
			log4go.Error(apperr.Message)
		}
	} else {
		if apperr := scheduledJob.Store.ScheduledJob().Update(scheduledJob.ScheduledJob); apperr != nil {
			log4go.Error(apperr.Message)
		}
	}
}

func (scheduledJob *PersistentScheduledJob) onProcessing() {
	log4go.Info("The scheduled job %s is in processing", scheduledJob.ID)

	scheduledJob.Status = models.JobProcessing
	scheduledJob.RanAt = Now()
	scheduledJob.Save()

	scheduledJob.ParentJob.JobStats.LastRanAt = Now()
	scheduledJob.ParentJob.Save()
}

func (scheduledJob *PersistentScheduledJob) onFailed(err error) {
	log4go.Info("The scheduled job %s failed", scheduledJob.ID)

	scheduledJob.Status = models.JobFailed
	scheduledJob.Error = err
	scheduledJob.Save()

	scheduledJob.ParentJob.JobStats.ErrorCount++
	scheduledJob.ParentJob.JobStats.LastError = err.Error()
	scheduledJob.ParentJob.JobStats.LastErrorAt = Now()
	scheduledJob.ParentJob.Save()
}

func (scheduledJob *PersistentScheduledJob) onSucceeded() {
	log4go.Info("The scheduled job %s has been succeeded", scheduledJob.ID)

	scheduledJob.Status = models.JobSucceeded
	scheduledJob.Save()

	scheduledJob.ParentJob.JobStats.SuccessCount++
	scheduledJob.ParentJob.JobStats.LastSuccededAt = Now()
	scheduledJob.ParentJob.Save()
}

func (scheduledJob *PersistentScheduledJob) Run() error {
	log4go.Info("Process scheduled job %s", scheduledJob.ID)

	scheduledJob.onProcessing()

	req, err := http.NewRequest(scheduledJob.ParentJob.Args.Method, scheduledJob.ParentJob.Args.URL, strings.NewReader(scheduledJob.ParentJob.Args.Body))
	if err != nil {
		if !scheduledJob.ParentJob.IsAsync {
			scheduledJob.onFailed(err)
		}
		return err
	}

	for value, key := range scheduledJob.ParentJob.Args.Headers {
		req.Header.Add(value, key)
	}

	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		if !scheduledJob.ParentJob.IsAsync {
			scheduledJob.onFailed(err)
		}
		return err
	}

	if res.StatusCode <= 299 {
		err = errors.New(res.Status)
		if !scheduledJob.ParentJob.IsAsync {
			scheduledJob.onFailed(err)
		}
		return err
	}

	if !scheduledJob.ParentJob.IsAsync {
		scheduledJob.onSucceeded()
	}

	return nil
}

func NewPersistentScheduledJob(store stores.Store, parentJob *PersistentJob) *PersistentScheduledJob {
	scheduledJob := models.ScheduledJob{
		JobID:       parentJob.ID,
		ScheduledAt: parentJob.NextRunAt,
		Status:      models.JobPending,
	}
	persistentJob := &PersistentScheduledJob{
		ScheduledJob: &scheduledJob,
		Store:        store,
		ParentJob:    parentJob,
	}
	return persistentJob
}
