package data

type Project struct {
	ID          int32  `db:"id,omitempty"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
}
