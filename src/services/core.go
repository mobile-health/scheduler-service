package services

import (
	"github.com/fvbock/endless"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/schedulers"
	"github.com/mobile-health/scheduler-service/src/stores"
	goji "goji.io"
)

type Srv struct {
	Router    *goji.Mux
	Store     stores.Store
	Scheduler *schedulers.Scheduler
}

func NewServer(router *goji.Mux, store stores.Store) *Srv {
	return &Srv{
		Router:    router,
		Store:     store,
		Scheduler: schedulers.NewScheduler(schedulers.DefaultMaxWorker),
	}
}

func (srv *Srv) Run() {
	endless.ListenAndServe(config.GetConfig().Server.ListenAddress, srv.Router)
}
