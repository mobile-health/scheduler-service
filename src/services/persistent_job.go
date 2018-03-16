package services

import (
	"time"

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

func (job *PersistentJob) Args() interface{} {
	return job.Job.Args
}

func (job *PersistentJob) ScheduledJob() schedulers.ScheduledJob {
	return job.scheduledJob
}

func (job *PersistentJob) Schedule(now time.Time, args interface{}) error {
	expr := cronexpr.MustParse(job.Expression)

	job.NextRunAt = expr.Next(now)
	job.scheduledJob = NewPersistentScheduledJob(job.Store, job.Job)

	return nil
}
