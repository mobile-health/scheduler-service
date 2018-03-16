package schedulers

import (
	"sort"
	"sync"
	"time"
)

type WorkerPool chan *Worker

type Scheduler struct {
	jobs      Jobs
	maxWorker int
	pool      WorkerPool
	isRuning  bool

	addJob     JobChannel
	processJob ScheduledJobChannel
	stop       StopChannel
	mutex      *sync.Mutex
}

func Now() time.Time {
	return time.Now().UTC()
}

func NewScheduler(maxWorker int) *Scheduler {

	r := &Scheduler{
		maxWorker:  maxWorker,
		pool:       make(WorkerPool, DefaultMaxQueue),
		addJob:     make(JobChannel, DefaultMaxQueue),
		processJob: make(ScheduledJobChannel, DefaultMaxQueue),
		stop:       make(StopChannel),
		mutex:      &sync.Mutex{},
	}

	for i := 0; i < maxWorker; i++ {
		worker := NewWorker(r.pool)
		worker.Start()
	}

	return r
}

func (n *Scheduler) startScheduler() {

	go func() {
		for {
			select {
			case job := <-n.processJob:
				w := <-n.pool
				w.Enqueue(job)
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

	if !n.GetRuningState() {
		return
	}

	n.SetRuningState(false)
	n.stop <- struct{}{}
	n.stop <- struct{}{}
}

func (n *Scheduler) GetRuningState() bool {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.isRuning
}

func (n *Scheduler) SetRuningState(b bool) {
	n.mutex.Lock()
	n.isRuning = b
	n.mutex.Unlock()
}

func (n *Scheduler) Start() *Scheduler {

	if n.GetRuningState() {
		return n
	}
	n.SetRuningState(true)
	go n.run()
	return n
}

func (n *Scheduler) run() {

	n.startScheduler()
	now := Now()

	for _, job := range n.jobs {
		job.Schedule(now, job.Args())
	}

	for {

		sort.Sort(n.jobs)

		var timer *time.Timer
		if len(n.jobs) == 0 || n.jobs[0].ScheduledJob() == nil {
			timer = time.NewTimer(1 * time.Hour)
		} else {
			timer = time.NewTimer(n.jobs[0].ScheduledJob().ScheduledAt().Sub(now))
		}

		for {
			select {
			case now := <-timer.C:

				for _, job := range n.jobs {
					if job.ScheduledJob().ScheduledAt().After(now) || job.ScheduledJob().ScheduledAt().IsZero() {
						break
					}
					scheduledJob := job.ScheduledJob()
					scheduledJob.Save()

					n.processJob <- scheduledJob
					job.Schedule(now, job.Args())
				}

			case newjob := <-n.addJob:
				timer.Stop()
				newjob.Schedule(Now(), newjob.Args())

				if newjob.ScheduledJob() == nil {
					n.jobs = append(n.jobs, newjob)
				}

			case <-n.stop:
				return
			}

			break
		}
	}
}

func (n *Scheduler) Add(job Job) {
	if !n.GetRuningState() {
		job.Schedule(Now(), job.Args())
		if job.ScheduledJob() == nil {
			n.jobs = append(n.jobs, job)
		}
	} else {
		n.addJob <- job
	}
}
