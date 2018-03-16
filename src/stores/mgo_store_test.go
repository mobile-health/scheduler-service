package stores

import (
	"testing"
	"time"

	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMgoSession(t *testing.T) {
	t.Run("ConnectSucess", func(t *testing.T) {
		mgoSession, err := NewMgoSession()
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, mgoSession)
	})

	t.Run("ConnectFailure", func(t *testing.T) {
		defaultConfig := config.GetConfig()
		defer config.SetConfig(&defaultConfig)

		mookConfig := config.Config{
			Database: config.DatabaseSettings{
				MgoDsn:         "localhost:1111",
				Retries:        0,
				ConnectTimeout: 1,
			},
		}
		config.SetConfig(&mookConfig)

		if _, err := NewMgoSession(); err == nil {
			t.Fatal("should not connected to the mongodb")
		}
	})

	t.Run("RetryConnect", func(t *testing.T) {
		defaultConfig := config.GetConfig()
		defer config.SetConfig(&defaultConfig)

		mookConfig := config.Config{
			Database: config.DatabaseSettings{
				MgoDsn:         "localhost:1111",
				Retries:        60,
				ConnectTimeout: 1,
			},
		}
		config.SetConfig(&mookConfig)

		go func() {
			<-time.NewTimer(2 * time.Second).C
			mookConfig.Database.MgoDsn = "localhost"
		}()

		if _, err := NewMgoSession(); err != nil {
			t.Fatal(err)
		}
	})

}

func TestNewMgoStore(t *testing.T) {
	m := NewStore()
	assert.NotNil(t, m)
	assert.NotNil(t, m.Job())
	assert.NotNil(t, m.ScheduledJob())
}
