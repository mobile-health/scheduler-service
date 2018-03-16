package stores

import (
	"time"

	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/config"
	"gopkg.in/mgo.v2"
)

func init() {
	time.Local = time.UTC
}

const (
	Database               = "scheduler"
	JobCollection          = "jobs"
	JobScheduledCollection = "scheduled_jobs"
)

func NewMgoSession() (*mgo.Session, error) {
	var retries = 0
	var mgoSession *mgo.Session
	var err error

CONNECT:
	mgoSession, err = mgo.DialWithTimeout(config.GetConfig().Database.MgoDsn, config.GetConfig().Database.ConnectTimeout*time.Second)
	if err != nil {
		if retries < config.GetConfig().Database.Retries {
			time.Sleep(time.Second)
			goto CONNECT
		}
		return nil, err
	}

	if err := mgoSession.Ping(); err != nil {
		log4go.Crash(err)
	}

	mgoSession.SetPoolLimit(1024)
	mgoSession.SetMode(mgo.Strong, true)
	mgoSession.SetSafe(&mgo.Safe{})

	return mgoSession, nil
}

type MgoStore struct {
	db           *mgo.Session
	job          JobStore
	scheduledJob ScheduledJobStore
}

func NewStore() Store {
	mgoSession, err := NewMgoSession()
	if err != nil {
		log4go.Critical("Failed to connect to the mongodb, got error %s", err.Error())
	}

	m := &MgoStore{
		db: mgoSession,
	}

	m.job = NewMgoJobStore(m)
	m.scheduledJob = NewMgoScheduledJobStore(m)
	return m
}

func (m *MgoStore) Job() JobStore {
	return m.job
}

func (m *MgoStore) ScheduledJob() ScheduledJobStore {
	return m.scheduledJob
}
