package dto

import (
	"fmt"

	"github.com/gavrilaf/wardrobe/pkg/utils/timex"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

type InfoObject struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Source    string   `json:"source"`
	Author    string   `json:"author"`
	Published string   `json:"published"`
	Created   string   `json:"created"`
	Finalized string   `json:"finalized"`
	Files     []File   `json:"files"`
	Tags      []string `json:"tags"`
}

func (o InfoObject) ToDBType() (dbtypes.InfoObject, error) {
	tm, err := timex.ParseJsonTime(o.Published)
	if err != nil {
		return dbtypes.InfoObject{}, fmt.Errorf("invalid published time %s (%s, %s, %s), %w", o.Published, o.Name, o.Source, o.Author, err)
	}

	return dbtypes.InfoObject{
		Name:      o.Name,
		Author:    o.Author,
		Source:    o.Source,
		Published: tm,
	}, nil
}

func InfoObjectFromDBType(o dbtypes.InfoObject) InfoObject {
	obj := InfoObject{
		ID:        o.ID,
		Name:      o.Name,
		Source:    o.Source,
		Author:    o.Author,
		Published: timex.TimeToJsonString(o.Published),
		Created:   timex.TimeToJsonString(o.Created),
		Files:     []File{},
		Tags:      []string{},
	}

	if o.Finalized != nil {
		obj.Finalized = timex.TimeToJsonString(*o.Finalized)
	}

	return obj
}
