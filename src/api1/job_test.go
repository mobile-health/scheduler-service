package api1

import (
	"net/http"
	"sync"
	"testing"

	"github.com/icrowley/fake"

	"github.com/mobile-health/scheduler-service/src/models"
)

func TestCreateJob(t *testing.T) {
	server := NewTestServer().Run()
	defer server.Stop()
	client := NewClient(server.Srv.Router)

	job := models.Job{
		Name:       "test_" + fake.JobTitle(),
		Expression: "0/1 * * * * * *",
		Args: models.RemoteArgs{
			URL:    server.ExternalSrv.URL,
			Method: http.MethodGet,
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	server.SetExternalJobFunc(f)

	_, err := client.CreateJob(job)
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}
