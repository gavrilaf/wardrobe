package dto

import "github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"

type Stat struct {
	ObjectsCount int64 `json:"objects_count"`
	FilesCount   int64 `json:"files_count"`
	TotalSize    int64 `json:"total_size"`
}

func StatFromDDType(s dbtypes.Stat) Stat {
	return Stat{
		ObjectsCount: s.ObjectsCount,
		FilesCount:   s.FilesCount,
		TotalSize:    s.TotalSize,
	}
}
