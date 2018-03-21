package api1

import (
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/services"
	"github.com/mobile-health/scheduler-service/src/utils"
	"goji.io/pat"
)

func getJob(c *services.Context) utils.Render {
	jobID := pat.Param(c.Request, "id")

	if job, apperr := c.Srv.Store.Job().Get(jobID); apperr != nil {
		return c.Error(apperr)
	} else {
		return c.JSON(200, job)
	}
}

func createJob(c *services.Context) utils.Render {
	job := models.NewJobFromBody(c.Request.Body)

	if apperr := c.CreateJob(job); apperr != nil {
		return c.Error(apperr)
	}

	return c.JSON(201, models.NewJsonResponse(job, nil))
}

func disableJob(c *services.Context) utils.Render {
	jobID := pat.Param(c.Request, "id")

	if job, apperr := c.Srv.Store.Job().Get(jobID); apperr != nil {
		return c.Error(apperr)
	} else if job.IsDisabled {
		return c.Error(models.NewError("api.job.disable.job_disabled.app_err", nil, 410))
	}

	if job, apperr := c.DisableJob(jobID); apperr != nil {
		return c.Error(apperr)
	} else {
		return c.JSON(200, job)
	}
}
