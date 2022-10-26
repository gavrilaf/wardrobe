package dbtypes

type FO struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	ContentType string `db:"content_type"`
	Size        int64  `db:"size"`
}
