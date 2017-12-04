package services

import (
	"github.com/fvbock/endless"
	"github.com/mobile-health/go-api-boilerplate/src/config"
	"github.com/mobile-health/go-api-boilerplate/src/stores"
	goji "goji.io"
)

type Srv struct {
	Router *goji.Mux
	Store  stores.Store
}

func NewServer(router *goji.Mux, store stores.Store) *Srv {
	return &Srv{
		Router: router,
		Store:  store,
	}
}

func (srv *Srv) Run() {
	endless.ListenAndServe(config.Config().Server.ListenAddress, srv.Router)
}
