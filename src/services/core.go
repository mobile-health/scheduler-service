package services

import (
	"github.com/canhlinh/log4go"
	"github.com/fvbock/endless"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/models"
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
		Scheduler: schedulers.New(),
	}
}

func (srv *Srv) Run() {

	srv.WakeupScheduler()                                                       // Start scheduler service
	endless.ListenAndServe(config.GetConfig().Server.ListenAddress, srv.Router) // Start api server
	srv.Scheduler.Stop()
}

func (s *Srv) WakeupScheduler() *models.Error {

	if jobs, apperr := s.Store.Job().FindNotDoneYet(); apperr != nil {
		log4go.Error(apperr)
		return apperr
	} else {
		var persistentJobs schedulers.Jobs
		for _, job := range jobs {
			presistentJob := NewPersistentJob(s.Store, job)
			persistentJobs = append(persistentJobs, presistentJob)
		}

		s.Scheduler.PreLoadExistingJob(persistentJobs)
	}

	s.Scheduler.Start()
	return nil
}
