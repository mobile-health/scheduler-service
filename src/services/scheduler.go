package services

import (
	"github.com/mobile-health/scheduler-service/src/models"
)

func (c *Context) CreateJob(job *models.Job) *models.Error {

	if apperr := c.Srv.Store.Job().Insert(job); apperr != nil {
		return apperr
	}

	c.Srv.Scheduler.Add(NewPersistentJob(c.Srv.Store, job))

	return nil
}

func (c *Context) ReportJob(jobID string, isSuccess bool, err error) *models.Error {
	return nil
}
