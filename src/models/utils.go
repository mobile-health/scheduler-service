package models

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/oklog/ulid"
)

func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func NewID() string {
	return ulid.MustNew(ulid.Timestamp(Now()), rand.Reader).String()
}

type Pagination struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func Equal(m1, m2 interface{}) bool {
	j1, _ := json.Marshal(m1)
	j2, _ := json.Marshal(m2)

	return bytes.Compare(j1, j2) == 0
}
