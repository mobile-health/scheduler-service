package services

import (
	"fmt"
	"sync"
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
	mux          *sync.Mutex
}

func NewPersistentJob(store stores.Store, job *models.Job) *PersistentJob {
	persistentJob := PersistentJob{
		Store: store,
		Job:   job,
		mux:   &sync.Mutex{},
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

	if job.MaxSchedule > 0 && job.JobStats.SuccessCount >= job.MaxSchedule {
		log4go.Info("The job %s is reached out of the maximum of execution time", job.ID)
		return fmt.Errorf("The job %s reached out of the maximum of execution time", job.ID)
	}

	expr := cronexpr.MustParse(job.Expression)
	next := expr.Next(now)
	if next.IsZero() {
		log4go.Warn("The job %s has already been done", job.ID)
		return fmt.Errorf("The job %s has already been done", job.ID)
	}

	job.NextRunAt = expr.Next(now)
	job.Save()
	job.scheduledJob = NewPersistentScheduledJob(job.Store, job)

	log4go.Debug("Job %s scheduled next run at %s", job.ID, job.NextRunAt)

	return nil
}

func (job *PersistentJob) GetID() string {
	return job.ID
}

func (job *PersistentJob) Disable() {

	job.IsDisabled = true
	job.Save()
}

func (job *PersistentJob) Finish() {

	job.IsDone = true
	job.Save()
}

func (job *PersistentJob) Save() {
	job.mux.Lock()
	defer job.mux.Unlock()

	if apperr := job.Store.Job().Update(job.Job); apperr != nil {
		log4go.Error(apperr)
	}
}
