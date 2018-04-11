package schedulers

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/canhlinh/log4go"
)

type WorkerPool chan *Worker
type DisableJob struct {
	ID      string
	JobChan chan Job
	ErrChan chan error
}

type Scheduler struct {
	jobs      MapJob
	maxWorker int
	pool      WorkerPool
	isRuning  bool

	addJob        JobChannel
	disableJob    chan *DisableJob
	processJob    ScheduledJobChannel
	stopScheduler StopChannel
	stopProcesser StopChannel
	mutex         *sync.Mutex
}

func Now() time.Time {
	return time.Now().UTC()
}

func New() *Scheduler {
	return NewScheduler(DefaultMaxWorker, DefaultMaxQueue)
}

func NewScheduler(maxWorker, maxQueue int) *Scheduler {
	if maxWorker <= 0 {
		panic("Must set at least one worker")
	}

	r := &Scheduler{
		maxWorker:     maxWorker,
		pool:          make(WorkerPool),
		addJob:        make(JobChannel, maxQueue),
		processJob:    make(ScheduledJobChannel, maxQueue),
		disableJob:    make(chan *DisableJob),
		stopProcesser: make(StopChannel),
		stopScheduler: make(StopChannel),
		mutex:         &sync.Mutex{},
		jobs:          MapJob{},
	}

	return r
}

func (n *Scheduler) startScheduler() {
	for i := 0; i < n.maxWorker; i++ {
		worker := NewWorker(n.pool)
		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-n.processJob:
				w := <-n.pool
				w.jobChannel <- job
			case <-n.stopProcesser:
				for i := 0; i < n.maxWorker; i++ {
					worker := <-n.pool
					worker.Stop()
				}
				n.stopProcesser <- struct{}{}
				return
			}
		}
	}()
}

func (n *Scheduler) Stop() {
	log4go.Info("Stop scheduler")

	if !n.IsRuning() {
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.isRuning = false

	n.stopScheduler <- struct{}{}
	n.stopProcesser <- struct{}{}
	<-n.stopProcesser // wait processer stop
}

func (n *Scheduler) IsRuning() bool {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.isRuning
}

func (n *Scheduler) Start() *Scheduler {
	if n.IsRuning() {
		return n
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.isRuning = true

	go n.run()
	return n
}

func (n *Scheduler) run() {
	log4go.Info("Start scheduler")

	n.startScheduler()
	now := Now()

	for _, job := range n.jobs {
		job.Schedule(now)
	}

	for {
		jobs := n.jobs.Jobs()
		sort.Sort(jobs)

		var timer *time.Timer
		if len(jobs) == 0 || !jobs[0].HasScheduledJob() {
			timer = time.NewTimer(24 * time.Hour)
		} else {
			timer = time.NewTimer(jobs[0].ScheduledJob().ScheduledAt().Sub(now))
		}

		for {
			select {
			case now = <-timer.C:

				for _, job := range jobs {
					if job.ScheduledJob().ScheduledAt().After(now) {
						break
					}

					scheduledJob := job.ScheduledJob()
					scheduledJob.Save()
					n.processJob <- scheduledJob

					if err := job.Schedule(now); err != nil {
						delete(n.jobs, job.GetID())
						job.Finish()
					}
				}

			case newjob := <-n.addJob:
				timer.Stop()
				now = Now()

				if err := newjob.Schedule(now); err == nil {
					n.jobs[newjob.GetID()] = newjob
				}

			case disableJob := <-n.disableJob:

				if job, exist := n.jobs[disableJob.ID]; exist {
					delete(n.jobs, job.GetID())
					job.Disable()
					log4go.Debug("Send disabled job to channel")
					disableJob.JobChan <- job
				} else {
					disableJob.ErrChan <- fmt.Errorf("Job %s not found", disableJob.ID)
				}

			case <-n.stopScheduler:
				return
			}

			break
		}
	}
}

func (n *Scheduler) PreLoadExistingJob(jobs Jobs) {
	n.jobs = NewMapJob(jobs)
}

func (n *Scheduler) ProcessScheduledJob(scheduledJob ScheduledJob) {
	n.processJob <- scheduledJob
}

func (n *Scheduler) Add(newjob Job) {
	if n.IsRuning() {
		n.addJob <- newjob
	} else {
		if err := newjob.Schedule(Now()); err == nil {
			n.jobs[newjob.GetID()] = newjob
		}
	}
}

func (n *Scheduler) DisableJob(jobID string) (Job, error) {
	log4go.Info("Request disable job %s", jobID)

	var disableJob = DisableJob{
		ID:      jobID,
		JobChan: make(JobChannel),
	}
	n.disableJob <- &disableJob

	select {
	case job := <-disableJob.JobChan:
		return job, nil
	case err := <-disableJob.ErrChan:
		return nil, err
	}
}
