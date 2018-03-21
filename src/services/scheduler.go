package services

import (
	"github.com/mobile-health/scheduler-service/src/models"
)

func (c *Context) CreateJob(job *models.Job) *models.Error {

	if apperr := c.Srv.Store.Job().Insert(job); apperr != nil {
		return apperr
	}

	c.Srv.Scheduler.Add(NewPersistentJob(c.Srv.Store, job.Clone()))

	return nil
}

func (c *Context) ReportJob(jobID string, isSuccess bool, err error) *models.Error {
	return nil
}

func (c *Context) DisableJob(jobID string) (*models.Job, *models.Error) {

	job, err := c.Srv.Scheduler.DisableJob(jobID)
	if err != nil {
		return nil, models.NewError("services.scheduler.disable.app_err", map[string]interface{}{"Message": err.Error()}, 500)
	}

	return job.(*PersistentJob).Job, nil
}
