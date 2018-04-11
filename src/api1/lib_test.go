package api1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"goji.io/pat"

	"goji.io"

	"github.com/mobile-health/scheduler-service/src/api"
	"github.com/mobile-health/scheduler-service/src/config"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/mobile-health/scheduler-service/src/services"
	"github.com/mobile-health/scheduler-service/src/stores"
	"github.com/mobile-health/scheduler-service/src/utils"
)

type TestServer struct {
	Srv           *services.Srv
	fakeServer    *httptest.Server
	fakeServerMux *goji.Mux
}

func NewTestServer() *TestServer {
	config.Load("../../conf/config.yaml")
	utils.Init("../../i18n")

	srv := services.NewServer(goji.NewMux(), stores.NewStore())
	Init(srv)

	fakeServerMux := goji.NewMux()
	fakeServer := httptest.NewServer(fakeServerMux)

	server := &TestServer{
		Srv:           srv,
		fakeServer:    fakeServer,
		fakeServerMux: fakeServerMux,
	}
	return server
}

func (s *TestServer) JobHandler(f func()) string {

	endpoint := "/" + models.NewID()
	s.fakeServerMux.HandleFunc(pat.Get(endpoint), func(w http.ResponseWriter, r *http.Request) {
		f()
		w.WriteHeader(200)
	})

	return s.fakeServer.URL + endpoint
}

func (s *TestServer) Run() *TestServer {
	s.Srv.WakeupScheduler()
	return s
}

func (s *TestServer) Stop() {
	s.Srv.Scheduler.Stop()
	s.DeleteTestData()
	s.fakeServer.Close()
}

func (s *TestServer) DeleteTestData() {
	s.Srv.Store.Job().DeleteAll("test_")
}

type ClientV1 struct {
	*api.Client
}

func NewClient(mux *goji.Mux) *ClientV1 {
	return &ClientV1{
		Client: &api.Client{
			Mux:        mux,
			BaseApiURL: "/api/v1",
			ApiKey:     config.GetConfig().Auth.ApiToken,
			ApiLogin:   config.GetConfig().Auth.ApiLogin,
		},
	}
}

func (c *ClientV1) CreateJob(job models.Job) (models.MapInterface, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&job); err != nil {
		return nil, err
	}

	return c.DoPost("/jobs", &body)
}

func (c *ClientV1) DisableJob(jobID string) (models.MapInterface, error) {

	endpoint := fmt.Sprintf("/jobs/%s/disable", jobID)
	return c.DoPost(endpoint, nil)
}
