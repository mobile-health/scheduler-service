package schedulers

import (
	"time"
)

type ScheduledJob interface {
	Run() error
	ScheduledAt() time.Time
	Save()
}

type ScheduledJobs []ScheduledJob

type Job interface {
	GetID() string
	HasScheduledJob() bool
	ScheduledJob() ScheduledJob
	Schedule(t time.Time) error
	Disable()
	Finish()
	Save()
}

type Jobs []Job
type MapJob map[string]Job

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

func (jobs *Jobs) Remove(index int) {
	v := *jobs

	if index == 0 && len(v) == 1 {
		*jobs = Jobs{}
	} else if index == len(v)-1 {
		*jobs = v[:index]
	} else {
		*jobs = append(v[:index], v[index+1:]...)
	}
}

func NewMapJob(jobs Jobs) MapJob {
	var mjob = MapJob{}
	for _, job := range jobs {
		mjob[job.GetID()] = job
	}
	return mjob
}

func (mjob MapJob) Jobs() Jobs {
	var jobs Jobs
	for _, job := range mjob {
		jobs = append(jobs, job)
	}
	return jobs
}
