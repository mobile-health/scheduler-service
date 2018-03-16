package api1

import (
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/services"
	"github.com/mobile-health/scheduler-service/src/utils"
)

func createJob(c *services.Context) utils.Render {
	job := models.NewJobFromBody(c.Request.Body)

	if apperr := c.CreateJob(job); apperr != nil {
		return c.Error(apperr)
	}

	return c.JSON(201, models.NewJsonResponse(job, nil))
}
