package schedulers

import (
	"sync"
	"testing"
	"time"

	"github.com/gorhill/cronexpr"
)

const (
	OneSecond       = 1*time.Second + 10*time.Millisecond
	ExprEverySecond = "0-59 * * * * * *"
	ExprOutdate     = "* * * * * 1970"
)

func stop(runner *Scheduler) chan bool {
	ch := make(chan bool)
	go func() {
		runner.Stop()
		ch <- true
	}()
	return ch
}

type TaskMock struct {
	expr string
	job  *JobMock
}

func NewTaskMock(expr string) *TaskMock {
	j := &TaskMock{
		expr: expr,
	}
	return j
}

func (mock *TaskMock) ScheduledJob() ScheduledJob {
	return mock.job
}

func (mock *TaskMock) Args() interface{} {
	return &sync.WaitGroup{}
}

func (mock *TaskMock) Schedule(now time.Time, args interface{}) error {
	expr := cronexpr.MustParse(mock.expr)
	if nextRunAt := expr.Next(now); !nextRunAt.IsZero() {
		mock.job = NewJobMock(nextRunAt, args.(*sync.WaitGroup))
	}

	return nil
}

type JobMock struct {
	wait  *sync.WaitGroup
	runAt time.Time
}

func NewJobMock(t time.Time, wait *sync.WaitGroup) *JobMock {
	j := &JobMock{
		wait:  &sync.WaitGroup{},
		runAt: t,
	}
	return j
}

func (mock *JobMock) Run() error {

	<-time.NewTimer(time.Second).C
	mock.wait.Done()
	return nil
}

func (mock *JobMock) Save() {
}

func (mock *JobMock) ScheduledAt() time.Time {
	return mock.runAt
}

func TestRunJobSuccess(t *testing.T) {
	scheduler := NewScheduler(1)
	defer scheduler.Stop()

	task1 := NewTaskMock(ExprEverySecond)
	task2 := NewTaskMock(ExprEverySecond)

	scheduler.Add(task1)
	scheduler.Start()
	scheduler.Add(task2)

	task1.Args().(*sync.WaitGroup).Wait()
	task2.Args().(*sync.WaitGroup).Wait()
}

func TestRunOutDateJob(t *testing.T) {
	scheduler := NewScheduler(1)
	defer scheduler.Stop()

	scheduler.Add(NewTaskMock(ExprOutdate))
	scheduler.Start()
	scheduler.Add(NewTaskMock(ExprOutdate))

	select {
	case <-time.After(OneSecond):
		t.Fatal("Expected the runner will be stopped immediately")
	case <-stop(scheduler):
	}
}

func TestNoJobRun(t *testing.T) {
	scheduler := NewScheduler(1).Start()

	select {
	case <-time.After(OneSecond):
		t.Fatal("Expected the runner will be stopped immediately")
	case <-stop(scheduler):
	}
}
