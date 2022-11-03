package dto

import (
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
	"github.com/gavrilaf/wardrobe/pkg/utils"
)

type File struct {
	ID          int    `json:"id"`
	Bucket      string `json:"bucket"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Uploaded    string `json:"uploaded"`
}

func (f File) ToDbType() dbtypes.File {
	return dbtypes.File{
		Bucket:      f.Bucket,
		Name:        f.Name,
		ContentType: f.ContentType,
		Size:        f.Size,
	}
}

func FileFromDBType(f dbtypes.File) File {
	return File{
		ID:          f.ID,
		Bucket:      f.Bucket,
		Name:        f.Name,
		ContentType: f.ContentType,
		Size:        f.Size,
		Uploaded:    utils.TimeToJsonString(f.Uploaded),
	}
}
