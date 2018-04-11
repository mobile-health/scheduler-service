package schedulers

import (
	"sort"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestSortJobs(t *testing.T) {

	job1 := &JobMock{
		ID: fake.CharactersN(10),
		scheduledJob: &ScheduledJobMock{
			runAt: time.Now().Add(time.Minute),
		},
	}

	job2 := &JobMock{
		ID: fake.CharactersN(10),
		scheduledJob: &ScheduledJobMock{
			runAt: time.Now().Add(2 * time.Minute),
		},
	}

	job3 := &JobMock{
		ID: fake.CharactersN(10),
		scheduledJob: &ScheduledJobMock{
			runAt: time.Now().Add(3 * time.Minute),
		},
	}

	jobs := Jobs{job2, job1, job3}

	sort.Sort(jobs)

	assert.Equal(t, job1, jobs[0])
	assert.Equal(t, job2, jobs[1])
	assert.Equal(t, job3, jobs[2])
}
