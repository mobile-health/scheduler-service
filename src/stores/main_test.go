package stores

import (
	"testing"

	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/utils"
)

var testStore Store

func TestMain(m *testing.M) {
	config.Load("../../conf/config.yaml")
	utils.Init("../../i18n")
	testStore = NewStore()
	m.Run()
}
