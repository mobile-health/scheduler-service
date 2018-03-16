package stores

import (
	"github.com/mobile-health/scheduler-service/src/models"
	"gopkg.in/mgo.v2"
)

type Store interface {
	Job() JobStore
	ScheduledJob() ScheduledJobStore
}

type JobStore interface {
	C() *mgo.Collection
	Insert(job *models.Job) *models.Error
	Update(job *models.Job) *models.Error
	Get(jobID string) (*models.Job, *models.Error)
	FindNotDoneYet() (models.Jobs, *models.Error)
}

type ScheduledJobStore interface {
	C() *mgo.Collection
	Insert(scheduledJob *models.ScheduledJob) *models.Error
	Update(scheduledJob *models.ScheduledJob) *models.Error
	Get(scheduledJobID string) (*models.ScheduledJob, *models.Error)
	FindByJobID(jobID string) (models.ScheduledJobs, *models.Error)
}
