package services

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/mobile-health/go-api-boilerplate/src/models"
)

var (
	ImageMimes = map[string]interface{}{
		"image/gif":  nil,
		"image/jpeg": nil,
		"image/png":  nil,
		"image/bmp":  nil,
		"image/tiff": nil,
	}
	OtherMimes = map[string]interface{}{
		"text/plain":                                                              nil,
		"text/csv":                                                                nil,
		"application/pdf":                                                         nil,
		"application/msword":                                                      nil,
		"application/vnd.ms-excel":                                                nil,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": nil,
	}
)

type FileBaseInfo struct {
	content  []byte
	fileName string
	mimeType string
}

func (c *Context) Upload(fileHeader *multipart.FileHeader, comment string) RenderFunc {
	return c.JSON(200, models.FileInfo{})
}

func (c *Context) parseFile(fileHeader *multipart.FileHeader) (*FileBaseInfo, *models.Error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, models.NewErrorUnexpected(err, http.StatusInternalServerError)
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	io.Copy(buf, file)
	mime := http.DetectContentType(buf.Bytes()[:512])

	return &FileBaseInfo{
		content:  buf.Bytes(),
		mimeType: mime,
		fileName: fileHeader.Filename,
	}, nil
}

func (c *Context) uploadOriginal() {
}

func (c *Context) uploadPreview() {
}

func (c *Context) uploadThumb() {
}

func (c *Context) genPreviewImage() {
}

func (c *Context) genThumb() {
}

func IsAcceptMime(mime string) bool {
	if IsImageMime(mime) {
		return true
	}
	_, ok := OtherMimes[mime]
	return ok
}

func IsImageMime(mime string) bool {
	_, ok := ImageMimes[mime]
	return ok
}
