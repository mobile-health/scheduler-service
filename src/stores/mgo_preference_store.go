package stores

import (
	"github.com/canhlinh/log4go"
	"github.com/mobile-health/scheduler-service/src/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoPreferenceStore struct {
	*MgoStore
}

func NewMgoPreferenceStore(store *MgoStore) PreferenceStore {
	s := &MgoPreferenceStore{
		MgoStore: store,
	}
	return s
}

func (store *MgoPreferenceStore) C() *mgo.Collection {
	return store.db.DB(Database).C(PreferenceCollection)
}

func (store *MgoPreferenceStore) Save(preference *models.Preference) *models.Error {
	preference.PreSave()

	if apperr := preference.Validate(); apperr != nil {
		return apperr
	}

	if _, err := store.C().UpsertId(preference.ID, preference); err != nil {
		log4go.Error(err)
		return models.NewError("stores.preference.save.app_err", nil, 500)
	}

	return nil
}

func (store *MgoPreferenceStore) Get(name string) (*models.Preference, *models.Error) {
	var preference models.Preference

	if err := store.C().Find(bson.M{"name": name}).One(&preference); err != nil {
		return nil, models.NewError("stores.preference.get.app_err", nil, 500)
	}

	return &preference, nil
}
