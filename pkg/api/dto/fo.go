package dto

import (
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
	"github.com/gavrilaf/wardrobe/pkg/utils"
)

type FO struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	ContentType string   `json:"content_type"`
	Source      string   `json:"source"`
	Author      string   `json:"author"`
	Bucket      string   `json:"bucket"`
	FileName    string   `json:"file_name"`
	Tags        []string `json:"tags"`
	Size        int64    `json:"size"`
	Created     string   `json:"created"`
	Uploaded    string   `json:"uploaded"`
}

func (o FO) ToDBType() dbtypes.InfoObject {
	return dbtypes.InfoObject{
		ID:          o.ID,
		Name:        o.Name,
		ContentType: o.ContentType,
		Author:      o.Author,
		Source:      o.Source,
	}
}

func MakeFOFromDBType(o dbtypes.InfoObject) FO {
	fo := FO{
		ID:          o.ID,
		Name:        o.Name,
		ContentType: o.ContentType,
		Source:      o.Source,
		Author:      o.Author,
		Bucket:      o.Bucket,
		FileName:    o.FileName,
		Size:        o.Size,
		Created:     utils.TimeToJsonString(o.Created),
	}

	if o.Uploaded != nil {
		fo.Uploaded = utils.TimeToJsonString(*o.Uploaded)
	}

	return fo
}
