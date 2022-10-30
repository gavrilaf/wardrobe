package dbtypes

import "time"

type FO struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	ContentType string     `db:"content_type"`
	Author      string     `db:"author"`
	Source      string     `db:"source"`
	Bucket      string     `db:"bucket"`
	FileName    string     `db:"file_name"`
	Size        int64      `db:"size"`
	Created     time.Time  `db:"created"`
	Uploaded    *time.Time `db:"uploaded"`
}
