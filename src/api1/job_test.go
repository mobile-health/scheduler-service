package api1

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/mobile-health/scheduler-service/src/models"
)

const OverOneSecond = time.Second + 10*time.Millisecond

func wait(wg *sync.WaitGroup) chan bool {
	c := make(chan bool, 1)
	go func() {
		wg.Wait()
		c <- true
	}()
	return c
}

func NewWaitFunc(wg *sync.WaitGroup) func() {
	f := func() {
		wg.Done()
	}
	return f
}

func TestCreateJob(t *testing.T) {
	server := NewTestServer().Run()
	defer server.Stop()
	client := NewClient(server.Srv.Router)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	job := models.Job{
		Name:       "test_" + fake.JobTitle(),
		Expression: "0/1 * * * * * *",
		Args: models.RemoteArgs{
			URL:    server.JobHandler(NewWaitFunc(wg)),
			Method: http.MethodGet,
		},
	}

	job1, err := client.CreateJob(job)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.NewTimer(OverOneSecond).C:
		t.Fatal("job should executed")
	case <-wait(wg):
		// job executed
	}

	assert.Equal(t, job.Name, job1["data"].(map[string]interface{})["name"])
	assert.Equal(t, job.Expression, job1["data"].(map[string]interface{})["expression"])
}

func TestDisableJob(t *testing.T) {
	server := NewTestServer().Run()
	defer server.Stop()
	client := NewClient(server.Srv.Router)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	job := models.Job{
		Name:       "test_" + fake.JobTitle(),
		Expression: "0/1 * * * * * *",
		Args: models.RemoteArgs{
			URL:    server.JobHandler(NewWaitFunc(wg)),
			Method: http.MethodGet,
		},
	}

	job1, err := client.CreateJob(job)
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()

	wg.Add(1)
	if _, err := client.DisableJob(job1["data"].(map[string]interface{})["id"].(string)); err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.NewTimer(OverOneSecond).C:
	case <-wait(wg):
		t.Fatal("job should not executed")
	}
}

func runJob(client *ClientV1, server *TestServer, b *testing.B) {

	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		job := models.Job{
			Name:       "test_" + fake.JobTitle(),
			Expression: "* * * * * * *",
			Args: models.RemoteArgs{
				URL:    server.JobHandler(NewWaitFunc(wg)),
				Method: http.MethodGet,
			},
			MaxSchedule: 1,
		}

		if _, err := client.CreateJob(job); err != nil {
			b.Fatal(err)
		}
	}

	wg.Wait()
}

func BenchmarkDisableJob(b *testing.B) {
	server := NewTestServer().Run()
	defer server.Stop()
	client := NewClient(server.Srv.Router)

	for i := 0; i < b.N; i++ {
		runJob(client, server, b)
	}
}
