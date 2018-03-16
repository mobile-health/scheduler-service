package services

import (
	"testing"
	"time"

	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/utils"
)

func TestMain(m *testing.M) {
	utils.Init("../../i18n")
	config.Load("../../conf/config.yaml")

	m.Run()
	time.Sleep(1 * time.Second)
	log4go.Close()
}
