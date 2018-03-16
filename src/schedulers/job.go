package schedulers

import (
	"time"
)

type ScheduledJob interface {
	Run() error
	ScheduledAt() time.Time
	Save()
}

type Job interface {
	ScheduledJob() ScheduledJob
	Schedule(t time.Time, args interface{}) error
	Args() interface{}
}

type Jobs []Job

func (jobs Jobs) Len() int      { return len(jobs) }
func (jobs Jobs) Swap(i, j int) { jobs[i], jobs[j] = jobs[j], jobs[i] }
func (jobs Jobs) Less(i, j int) bool {

	if jobs[i].ScheduledJob() == nil {
		return false
	}

	if jobs[j].ScheduledJob() == nil {
		return true
	}

	return jobs[j].ScheduledJob().ScheduledAt().After(jobs[i].ScheduledJob().ScheduledAt())
}
