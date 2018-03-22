package models

import "strings"

const (
	PreferenceDBVersion = "db.version"
)

type Preference struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

func (p *Preference) PreSave() {
	if len(p.ID) != 26 {
		p.ID = NewID()
	}
	p.Name = strings.TrimSpace(p.Name)
	p.Value = strings.TrimSpace(p.Value)
}

func (p *Preference) Validate() *Error {
	var errFields = ErrorFields{}

	if len(p.Name) == 0 {
		errFields = append(errFields, NewErrorFieldInvalid("id"))
	}

	return errFields.GenAppError()
}
