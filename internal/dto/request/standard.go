package request

import "io"

type ImageUpload struct {
	FileName string
	Size     int64
	Content  io.Reader
}
