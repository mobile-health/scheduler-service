package schedulers

import (
	"net/http"
)

const (
	DefaultMaxWorker = 1024
	DefaultMaxQueue  = 1024
)

type ScheduledJobChannel chan ScheduledJob
type JobChannel chan Job
type StopChannel chan struct{}

type Worker struct {
	Client     *http.Client
	isRuning   bool
	pool       WorkerPool
	jobChannel ScheduledJobChannel
	stop       StopChannel
}

func NewWorker(pool WorkerPool) *Worker {

	return &Worker{
		Client:     &http.Client{},
		pool:       pool,
		jobChannel: make(ScheduledJobChannel),
		stop:       make(StopChannel),
	}
}

func (w *Worker) Start() {
	w.isRuning = true
	go w.processing()
}

func (w *Worker) processing() {
	for {
		w.pool <- w // I'm ready

		select {
		case job := <-w.jobChannel:
			job.Run()
		case <-w.stop:
			w.isRuning = false
			return
		}
	}
}

func (w *Worker) Enqueue(job ScheduledJob) {

	w.jobChannel <- job
}

func (w *Worker) Stop() {

	if !w.isRuning {
		return
	}

	w.stop <- struct{}{}
	w.isRuning = false
}
