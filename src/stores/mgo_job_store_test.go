package stores

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/mobile-health/scheduler-service/src/models"
	"github.com/stretchr/testify/assert"
)

func TestJobInsert(t *testing.T) {
	job := models.Job{
		Name:       fake.Brand(),
		Expression: "0 0 29 2 *",
		Args: models.RemoteArgs{
			URL:    "https://api.coinmarketcap.com/v1/ticker/",
			Method: "GET",
		},
	}

	if apperr := testStore.Job().Insert(&job); apperr != nil {
		t.Fatal(apperr)
	}

	assert.Len(t, job.ID, 26)
	assert.False(t, job.CreatedAt.IsZero())
	assert.False(t, job.UpdatedAt.IsZero())
}

func TestJobUpdate(t *testing.T) {

	job := models.Job{
		Name:       fake.Brand(),
		Expression: "0 0 29 2 *",
		Args: models.RemoteArgs{
			URL:    "https://api.coinmarketcap.com/v1/ticker/",
			Method: "GET",
		},
	}

	if apperr := testStore.Job().Insert(&job); apperr != nil {
		t.Fatal(apperr)
	}

	job.Name = fake.Brand()
	if apperr := testStore.Job().Update(&job); apperr != nil {
		t.Fatal(apperr)
	}
}

func TestJobGet(t *testing.T) {
	job := models.Job{
		Name:       fake.Brand(),
		Expression: "0 0 29 2 *",
		Args: models.RemoteArgs{
			URL:    "https://api.coinmarketcap.com/v1/ticker/",
			Method: "GET",
		},
	}

	if apperr := testStore.Job().Insert(&job); apperr != nil {
		t.Fatal(apperr)
	}

	job1, apperr := testStore.Job().Get(job.ID)
	if apperr != nil {
		t.Fatal(apperr)
	}
	assert.True(t, models.Equal(job1, job))
}

func TestJobFindNotDoneYet(t *testing.T) {
	job1 := models.Job{
		Name:       fake.Brand(),
		Expression: "0 0 29 2 *",
		IsDone:     false,
		Args: models.RemoteArgs{
			URL:    "https://api.coinmarketcap.com/v1/ticker/",
			Method: "GET",
		},
	}

	job2 := job1
	job2.IsDone = true

	if apperr := testStore.Job().Insert(&job1); apperr != nil {
		t.Fatal(apperr)
	}

	if apperr := testStore.Job().Insert(&job2); apperr != nil {
		t.Fatal(apperr)
	}

	jobs, apperr := testStore.Job().FindNotDoneYet()
	if apperr != nil {
		t.Fatal(apperr)
	}

	found := false
	notfound := false

	for _, job := range jobs {
		if models.Equal(job, job1) {
			found = true
			continue
		}

		if models.Equal(job, job2) {
			notfound = true
		}
	}
	assert.True(t, found)
	assert.False(t, notfound)
}
