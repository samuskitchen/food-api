package fileupload

import (
	"github.com/google/uuid"
	"path"
)

func FormatFile(fn string) string {

	ext := path.Ext(fn)
	id := uuid.New()

	newFileName := id.String() + ext

	return newFileName
}
