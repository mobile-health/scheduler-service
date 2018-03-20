package schedulers

import (
	"sort"
	"sync"
	"time"

	"github.com/canhlinh/log4go"
)

type WorkerPool chan *Worker

type Scheduler struct {
	jobs      Jobs
	maxWorker int
	pool      WorkerPool
	isRuning  bool

	addJob     JobChannel
	disableJob JobChannel
	processJob ScheduledJobChannel
	stop       StopChannel
	mutex      *sync.Mutex
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
		maxWorker:  maxWorker,
		pool:       make(WorkerPool),
		addJob:     make(JobChannel, maxQueue),
		processJob: make(ScheduledJobChannel, maxQueue),
		stop:       make(StopChannel),
		mutex:      &sync.Mutex{},
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
			case <-n.stop:
				for i := 0; i < n.maxWorker; i++ {
					worker := <-n.pool
					worker.Stop()
				}
				return
			}
		}
	}()
}

func (n *Scheduler) Stop() {

	if !n.IsRuning() {
		return
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.isRuning = false

	n.stop <- struct{}{}
	n.stop <- struct{}{}
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

		sort.Sort(n.jobs)

		var timer *time.Timer
		if len(n.jobs) == 0 || !n.jobs[0].HasScheduledJob() {
			timer = time.NewTimer(24 * time.Hour)
		} else {
			timer = time.NewTimer(n.jobs[0].ScheduledJob().ScheduledAt().Sub(now))
		}

		for {
			select {
			case now = <-timer.C:

				for _, job := range n.jobs {
					if job.ScheduledJob().ScheduledAt().After(now) || job.ScheduledJob().ScheduledAt().IsZero() {
						break
					}

					scheduledJob := job.ScheduledJob()
					scheduledJob.Save()
					n.processJob <- scheduledJob

					job.Schedule(now)
				}

			case newjob := <-n.addJob:
				timer.Stop()
				now = Now()

				if err := newjob.Schedule(now); err == nil {
					n.jobs = append(n.jobs, newjob)
				}

			case <-n.stop:
				return
			}

			break
		}
	}
}

func (n *Scheduler) PreLoadExistingJob(jobs Jobs) {
	n.jobs = jobs
}

func (n *Scheduler) ProcessScheduledJob(scheduledJob ScheduledJob) {
	n.processJob <- scheduledJob
}

func (n *Scheduler) Add(newjob Job) {
	if n.IsRuning() {
		n.addJob <- newjob
	} else {
		if err := newjob.Schedule(Now()); err == nil {
			n.jobs = append(n.jobs, newjob)
		}
	}
}
