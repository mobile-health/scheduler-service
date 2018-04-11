package schedulers

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/icrowley/fake"

	"github.com/gorhill/cronexpr"
	"github.com/stretchr/testify/assert"
)

const (
	OneSecond       = 1*time.Second + 10*time.Millisecond
	ExprEverySecond = "0/1 * * * * * *"
	ExprZeroTime    = "* * * * * 1970"
	ExprOutDate     = "* * * * * 2017"
)

func stop(runner *Scheduler) chan bool {
	ch := make(chan bool)
	go func() {
		runner.Stop()
		ch <- true
	}()
	return ch
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func getOneTimeExpr(second int) string {
	n := Now().Add(time.Second * time.Duration(second))

	return fmt.Sprintf("%d %d %d %d %d * %d", n.Second(), n.Minute(), n.Hour(), n.Day(), int(n.Month()), n.Year())
}

type JobMock struct {
	ID           string
	expr         string
	scheduledJob *ScheduledJobMock
	f            func()
}

func NewJobMock(expr string, f func()) *JobMock {
	j := &JobMock{
		ID:   fake.CharactersN(26),
		expr: expr,
		f:    f,
	}
	return j
}

func (mock *JobMock) GetID() string {
	return mock.ID
}

func (mock *JobMock) Disable() {
}

func (mock *JobMock) Finish() {
}

func (mock *JobMock) Save() {
}

func (mock *JobMock) ScheduledJob() ScheduledJob {
	return mock.scheduledJob
}

func (mock *JobMock) HasScheduledJob() bool {
	return mock.scheduledJob != nil
}

func (mock *JobMock) Schedule(now time.Time) error {
	expr := cronexpr.MustParse(mock.expr)

	if nextRunAt := expr.Next(now); !nextRunAt.IsZero() {
		mock.scheduledJob = NewScheduledJobMock(nextRunAt, mock.f)
	} else {
		mock.scheduledJob = nil
		return errors.New("Expired job")
	}

	return nil
}

type ScheduledJobMock struct {
	f     func()
	runAt time.Time
}

func NewScheduledJobMock(scheduledAt time.Time, f func()) *ScheduledJobMock {

	j := &ScheduledJobMock{
		f:     f,
		runAt: scheduledAt,
	}

	return j
}

func (scheduledJob *ScheduledJobMock) Run() error {
	scheduledJob.f()
	return nil
}

func (scheduledJob *ScheduledJobMock) Save() {
}

func (scheduledJob *ScheduledJobMock) ScheduledAt() time.Time {
	return scheduledJob.runAt
}

func TestEverySecodeExpr(t *testing.T) {
	expr := cronexpr.MustParse(ExprEverySecond)

	now := Now()
	next := expr.Next(now)

	assert.True(t, next.Sub(now).Seconds() <= 1)
}

func TestJobNothing(t *testing.T) {
	scheduler := NewScheduler(2, 2).Start()

	assert.Len(t, scheduler.jobs, 0)

	select {
	case <-time.After(OneSecond):
		t.Fatal("Expected the runner will be stopped immediately")
	case <-stop(scheduler):
		// Stopped
	}

}

func TestJobAddBeforeRunning(t *testing.T) {
	scheduler := NewScheduler(2, 2)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	job := NewJobMock(ExprEverySecond, f)
	scheduler.Add(job)
	scheduler.Start()
	defer scheduler.Stop()

	select {
	case <-time.After(OneSecond):
		t.Fatal("expected job runs")
	case <-wait(wg):
		// Job done
	}
}

func TestJobAddWhileRunning(t *testing.T) {
	scheduler := NewScheduler(2, 2).Start()
	defer scheduler.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	job := NewJobMock(ExprEverySecond, f)
	scheduler.Add(job)

	select {
	case <-time.After(OneSecond):
		t.Fatal("expected job runs")
	case <-wait(wg):
		//Job done
	}

}

func TestJobRunningTwice(t *testing.T) {
	scheduler := NewScheduler(2, 2).Start()
	defer scheduler.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	f := func() { wg.Done() }

	job := NewJobMock(ExprEverySecond, f)
	scheduler.Add(job)

	select {
	case <-time.After(2 * OneSecond):
		t.Fatal("expected job fires 2 times")
	case <-wait(wg):
		//Job done
	}

}

func TestJobRunOneTime(t *testing.T) {
	scheduler := NewScheduler(1, 0).Start()
	defer scheduler.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	scheduler.Add(NewJobMock(getOneTimeExpr(1), f))
	wg.Wait()
	wg.Add(1)

	select {
	case <-time.After(OneSecond):
		// Finished
	case <-wait(wg):
		t.Fatal("expected the job not run again")
	}
}

func TestJobRunZeroTime(t *testing.T) {
	scheduler := NewScheduler(1, 0).Start()
	defer scheduler.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	scheduler.Add(NewJobMock(ExprZeroTime, f))

	select {
	case <-time.After(OneSecond):
		// Finished
	case <-wait(wg):
		t.Fatal("expected the job never run")
	}
}

func TestJobExpiredTime(t *testing.T) {
	scheduler := NewScheduler(1, 0).Start()
	defer scheduler.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	f := func() { wg.Done() }

	scheduler.Add(NewJobMock(ExprZeroTime, f))

	select {
	case <-time.After(OneSecond):
		// Finished
	case <-wait(wg):
		t.Fatal("expected the job never run")
	}
}

func TestStopJobWithoutStart(t *testing.T) {
	scheduler := NewScheduler(2, 2)
	scheduler.Stop()
}

func TestStartJobWithZeroStart(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("Must cause a panic")
		}
	}()

	scheduler := NewScheduler(0, 2)
	scheduler.Stop()
}
