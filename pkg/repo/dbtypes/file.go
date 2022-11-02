package dbtypes

import "time"

type File struct {
	ID           int       `db:"id"`
	InfoObjectID int       `db:"info_object_id"`
	Bucket       string    `db:"bucket"`
	Name         string    `db:"name"`
	ContentType  string    `db:"content_type"`
	Size         int64     `db:"size"`
	Uploaded     time.Time `db:"uploaded"`
}
