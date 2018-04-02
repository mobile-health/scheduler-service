package api1

import (
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/services"
	"github.com/mobile-health/scheduler-service/src/utils"
)

const (
	ApiVersion = "1.0.0.0"
)

func version(c *services.Context) utils.Render {
	return c.JSON(200, models.ApiVerion{
		BuildCommit: config.BuildCommit,
		BuildDate:   config.BuildCommit,
		Status:      "OK",
		Version:     ApiVersion,
	})
}
