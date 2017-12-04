package models

import "time"

type ShareGroup struct {
	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	UserIDs   []string  `bson:"user_ids" json:"user_ids"`
}
