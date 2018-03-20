package services

import (
	"fmt"
	"time"

	"github.com/canhlinh/log4go"

	"github.com/gorhill/cronexpr"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/schedulers"
	"github.com/mobile-health/scheduler-service/src/stores"
)

type PersistentJob struct {
	*models.Job
	stores.Store
	scheduledJob *PersistentScheduledJob
}

func NewPersistentJob(store stores.Store, job *models.Job) *PersistentJob {
	persistentJob := PersistentJob{
		Store: store,
		Job:   job,
	}
	return &persistentJob
}

func (job *PersistentJob) HasScheduledJob() bool {
	return job.scheduledJob != nil
}

func (job *PersistentJob) ScheduledJob() schedulers.ScheduledJob {
	return job.scheduledJob
}

func (job *PersistentJob) Schedule(now time.Time) error {

	expr, err := cronexpr.Parse(job.Expression)
	if err != nil {
		return err
	}

	next := expr.Next(now)
	if next.IsZero() {
		job.Job.IsDone = true
		if apperr := job.Store.Job().Update(job.Job); apperr != nil {
			log4go.Error(apperr)
		}
		log4go.Warn("The job %s has already been done", job.ID)
		return fmt.Errorf("The job %s has already been done", job.ID)
	}

	job.NextRunAt = expr.Next(now)
	if apperr := job.Store.Job().Update(job.Job); apperr != nil {
		log4go.Error(apperr)
	}
	job.scheduledJob = NewPersistentScheduledJob(job.Store, job.Job)

	log4go.Debug("Job %s scheduled next run at %s", job.ID, job.NextRunAt)

	return nil
}
