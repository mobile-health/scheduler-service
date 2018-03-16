package services

import (
	"testing"

	"github.com/mobile-health/scheduler-service/src/stores"
	"github.com/stretchr/testify/assert"
	"goji.io"
)

func TestNewSrv(t *testing.T) {
	srv := NewServer(goji.NewMux(), stores.NewStore())
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.Router, "Router should not be nil")
	assert.NotNil(t, srv.Store, "Store should not be nil")
}
