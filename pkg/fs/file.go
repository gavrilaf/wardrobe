package fs

import "io"

type File struct {
	Bucket      string
	Name        string
	ContentType string
	Size        int64
	Reader      io.ReadCloser
}
