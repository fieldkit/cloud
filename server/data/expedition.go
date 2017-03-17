package data

type Expedition struct {
	ID          int32  `db:"id,omitempty"`
	ProjectID   int32  `db:"project_id"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
}
