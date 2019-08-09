package data

type Station struct {
	ID          int32  `db:"id,omitempty"`
	Name        string `db:"name"`
}
