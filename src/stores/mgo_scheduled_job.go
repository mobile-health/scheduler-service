package stores

import (
	"github.com/mobile-health/scheduler-service/src/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoScheduledJobStore struct {
	*MgoStore
}

func NewMgoScheduledJobStore(m *MgoStore) ScheduledJobStore {
	store := &MgoScheduledJobStore{m}
	return store
}

func (m *MgoScheduledJobStore) C() *mgo.Collection {
	return m.db.DB(Database).C(JobScheduledCollection)
}

func (m *MgoScheduledJobStore) Insert(scheduledJob *models.ScheduledJob) *models.Error {
	scheduledJob.PreInsert()

	if apperr := scheduledJob.Validate(); apperr != nil {
		return apperr
	}

	if err := m.C().Insert(scheduledJob); err != nil {
		return models.NewError("stores.scheduled_job_store.insert.app_err", nil, 500)
	}

	return nil
}

func (m *MgoScheduledJobStore) Update(scheduledJob *models.ScheduledJob) *models.Error {
	scheduledJob.PreInsert()

	if apperr := scheduledJob.Validate(); apperr != nil {
		return apperr
	}

	if err := m.C().Insert(scheduledJob); err != nil {
		return models.NewError("stores.scheduled_job_store.update.app_err", nil, 500)
	}

	return nil
}

func (m *MgoScheduledJobStore) Get(scheduledJobID string) (*models.ScheduledJob, *models.Error) {
	var scheduledJob models.ScheduledJob

	if err := m.C().FindId(scheduledJobID).One(&scheduledJob); err != nil {
		return nil, models.NewError("stores.scheduled_job_store.get.app_err", nil, 500)
	}

	return &scheduledJob, nil
}

func (m *MgoScheduledJobStore) FindByJobID(jobID string) (models.ScheduledJobs, *models.Error) {
	var scheduledJobs = make(models.ScheduledJobs, 0)

	if err := m.C().Find(bson.M{"job_id": jobID}).All(&scheduledJobs); err != nil {
		return nil, models.NewError("stores.scheduled_job_store.find.app_err", nil, 500)
	}

	return scheduledJobs, nil
}

func (m *MgoScheduledJobStore) FindProcessing() (models.ScheduledJobs, *models.Error) {
	var scheduledJobs = make(models.ScheduledJobs, 0)

	if err := m.C().Find(bson.M{"status": models.JobProcessing}).All(&scheduledJobs); err != nil {
		return nil, models.NewError("stores.scheduled_job_store.find.app_err", nil, 500)
	}

	return scheduledJobs, nil
}
