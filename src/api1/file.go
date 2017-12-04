package api1

import (
	"net/http"

	"github.com/mobile-health/go-api-boilerplate/src/config"
	"github.com/mobile-health/go-api-boilerplate/src/models"
	"github.com/mobile-health/go-api-boilerplate/src/services"
)

func uploadFile(c *services.Context) services.RenderFunc {
	err := c.Request.ParseMultipartForm(config.Config().File.MaxSize)
	if err != nil {
		return c.Error(models.NewErrorUnexpected(err, http.StatusBadRequest))
	}

	_, file, err := c.Request.FormFile("file")
	if err != nil {
		return c.Error(models.NewErrorUnexpected(err, http.StatusBadRequest))
	}

	comment := c.Request.PostFormValue("comment")
	return c.Upload(file, comment)
}
