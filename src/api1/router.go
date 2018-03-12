package api1

import (
	"github.com/mobile-health/go-api-boilerplate/src/services"
	goji "goji.io"
	"goji.io/pat"
)

func Init(srv *services.Srv) {
	apiMux := goji.SubMux()
	srv.Router.Handle(pat.New("/api/v1/*"), apiMux)
}
