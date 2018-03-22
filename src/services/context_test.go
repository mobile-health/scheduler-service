package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/stores"
	"github.com/stretchr/testify/assert"
	"goji.io"
)

func TestNewContext(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	c := NewContext(w, r, NewServer(goji.NewMux(), stores.NewStore()))

	assert.NotNil(t, c, "Context can not be null")
	assert.NotNil(t, c.Writer, "Context can not be null")
	assert.NotNil(t, c.Request, "Context can not be null")
	assert.Len(t, c.RequestID, 26, "Context can not be zero")
}

func TestContextJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	c := NewContext(w, r, NewServer(goji.NewMux(), stores.NewStore()))

	data := models.Job{
		ID: models.NewID(),
	}
	renderFunc := c.JSON(200, data)
	assert.NotNil(t, renderFunc, "renderFunc can not be null")
}

func TestContextError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	c := NewContext(w, r, NewServer(goji.NewMux(), stores.NewStore()))

	apperr := models.NewErrorUnexpected(errors.New(""), 400)
	renderFunc := c.Error(apperr)
	assert.NotNil(t, renderFunc, "renderFunc can not be null")
}
