package api

import (
	"net/http"
	"strconv"

	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/services"
)

type Api struct {
	Srv *services.Srv
}

func NewAPI(srv *services.Srv) *Api {
	return &Api{
		Srv: srv,
	}
}

func (api *Api) Handler(f services.HandlerFunc) http.Handler {
	return &handler{
		Srv:        api.Srv,
		handleFunc: f,
	}
}

type handler struct {
	handleFunc        services.HandlerFunc
	Srv               *services.Srv
	basicAuthRequired bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log4go.Error(err)
			http.Error(w, http.StatusText(500), 500)
		}
	}()

	c := services.NewContext(w, r, h.Srv)
	if !h.authenticate(w, r, c) {
		return
	}

	if h.handleFunc == nil {
		return
	}

	render := h.handleFunc(c)
	if render == nil {
		log4go.Error("render can not be null")
		return
	}

	render.Write()
}

func (h *handler) authenticate(w http.ResponseWriter, r *http.Request, c *services.Context) bool {
	realm := "Basic realm=" + strconv.Quote("Authorization Required")
	user, pass, _ := r.BasicAuth()

	if user != config.GetConfig().Auth.ApiToken || pass != config.GetConfig().Auth.ApiLogin {
		// Credentials doesn't match, Kyo return 401 and abort handlers chain.
		w.Header().Add("WWW-Authenticate", realm)
		c.Error(models.NewError("api.unauthorized.app_error", nil, 401)).Write()
		return false
	}

	return true
}
