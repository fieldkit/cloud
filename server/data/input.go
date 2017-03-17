package data

type Input struct {
	ID        int32  `db:"id,omitempty"`
	ProjectID int32  `db:"project_id"`
	Name      string `db:"name"`
	Slug      string `db:"slug"`
}
