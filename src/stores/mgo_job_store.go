package stores

import (
	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoJobStore struct {
	*MgoStore
}

func NewMgoJobStore(m *MgoStore) JobStore {
	store := &MgoJobStore{m}
	return store
}

func (m *MgoJobStore) C() *mgo.Collection {
	return m.db.Clone().DB(Database).C(JobCollection)
}

func (m *MgoJobStore) Insert(job *models.Job) *models.Error {

	job.PreInsert()

	if apperr := job.Validate(); apperr != nil {
		return apperr
	}

	if err := m.C().Insert(&job); err != nil {
		log4go.Error(err)
		return models.NewError("stores.job.insert.app_err", nil, 500)
	}

	return nil
}

func (m *MgoJobStore) Update(job *models.Job) *models.Error {

	job.PreUpdate()

	if apperr := job.Validate(); apperr != nil {
		return apperr
	}

	if err := m.C().UpdateId(job.ID, job); err != nil {
		log4go.Error(err)
		return models.NewError("stores.job.update.app_err", nil, 500)
	}

	return nil
}

func (m *MgoJobStore) Get(jobID string) (*models.Job, *models.Error) {

	var job models.Job

	if err := m.C().FindId(jobID).One(&job); err != nil {
		log4go.Error(err)
		return nil, models.NewError("stores.job.get.app_err", nil, 400)
	}

	return &job, nil
}

func (m *MgoJobStore) FindNotDoneYet() (models.Jobs, *models.Error) {

	var jobs = models.Jobs{}

	if err := m.C().Find(bson.M{"is_done": false}).All(&jobs); err != nil {
		log4go.Error(err)
		return nil, models.NewError("stores.job.find.app_err", nil, 500)
	}

	return jobs, nil
}
