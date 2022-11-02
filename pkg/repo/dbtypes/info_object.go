package dbtypes

import "time"

type InfoObject struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Author    string     `db:"author"`
	Source    string     `db:"source"`
	Published time.Time  `db:"published"`
	Created   time.Time  `db:"created"`
	Uploaded  *time.Time `db:"uploaded"`
}
