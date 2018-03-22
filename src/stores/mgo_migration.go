package stores

import (
	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/models"
	"gopkg.in/mgo.v2"
)

const (
	// DbVersionZero initial db version
	DbVersionZero = "0"
	DbVersion1    = "1"
)

func exitWithError(err error) {
	log4go.Crash(err)
}

func (store *MgoStore) getCurrentDBVerion() (version string) {
	defer func(v string) {
		log4go.Info("Current db version is %s", version)
	}(version)

	if p, apperr := store.Preference().Get(models.PreferenceDBVersion); apperr != nil {
		version = DbVersionZero
		return
	} else {
		version = p.Value
		return
	}
}

func (store *MgoStore) saveDbVersion(version string) {
	preference := models.Preference{
		Name:  models.PreferenceDBVersion,
		Value: version,
	}

	if apperr := store.Preference().Save(&preference); apperr != nil {
		log4go.Error(apperr)
	}

	log4go.Info("Upgrade to db version %s", version)
}

func (store *MgoStore) Upgrade() {
	store.upgradeToV1()
}

func (store *MgoStore) upgradeToV1() {
	if store.getCurrentDBVerion() != DbVersionZero {
		return
	}

	if err := store.db.DB(Database).C(JobCollection).EnsureIndex(mgo.Index{
		Key:    []string{"fu_id"},
		Name:   "udx_fu_id",
		Unique: true,
	}); err != nil {
		exitWithError(err)
	}

	if err := store.db.DB(Database).C(PreferenceCollection).EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Name:   "udx_name",
		Unique: true,
	}); err != nil {
		exitWithError(err)
	}

	store.saveDbVersion(DbVersion1)
}
