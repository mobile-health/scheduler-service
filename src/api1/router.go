package api1

import (
	"github.com/mobile-health/scheduler-service/src/api"
	"github.com/mobile-health/scheduler-service/src/services"
	goji "goji.io"
	"goji.io/pat"
)

func Init(srv *services.Srv) {
	apiMux := goji.SubMux()
	srv.Router.Handle(pat.New("/api/v1/*"), apiMux)

	api1 := api.NewAPI(srv)
	apiMux.Handle(pat.Post("/"), api1.Handler(version))
	srv.Router.Handle(pat.New("/api/v1"), api1.Handler(version))

	apiMux.Handle(pat.Post("/jobs"), api1.Handler(createJob))
	apiMux.Handle(pat.Get("/jobs/:id"), api1.Handler(getJob))
	apiMux.Handle(pat.Post("/jobs/:id/disable"), api1.Handler(disableJob))
}
