package models

import (
	"time"
)

const (
	AzureContainerPrivate = "private-files"
	AzureContainerPublic  = "public-files"
)

type FileInfo struct {
	ID              string                 `bson:"_id" json:"id"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
	CreatorUserID   string                 `bson:"creator_user_id" json:"-"`
	Container       string                 `bson:"container" json:"-"`
	OriginalPath    string                 `bson:"original_path" json:"original_path"`
	PreviewPath     string                 `bson:"preview_path" json:"preview_path"`
	ThumbPath       string                 `bson:"thumb_path" json:"thumb_path"`
	Extension       string                 `bson:"extension" json:"extension"`
	Mime            string                 `bson:"mime" json:"mine"`
	Size            int                    `bson:"size" json:"size"`
	Comment         string                 `bson:"comment" json:"comment"`
	HasPreviewImage bool                   `bson:"has_preview_image" json:"has_preview_image"`
	IsPublic        bool                   `bson:"is_public" json:"is_public"`
	Meta            map[string]interface{} `bson:"meta" json:"meta"`
}

func (f *FileInfo) PreInsert(userID string, container string) {
	f.ID = NewID()
	f.CreatedAt = Now()
	f.UpdatedAt = Now()
	f.CreatorUserID = userID
	f.Container = container
}

func (f *FileInfo) Validate() *Error {
	return nil
}
