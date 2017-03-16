package data

type Team struct {
	ID           int32  `db:"id,omitempty"`
	ProjectID    int32  `db:"project_id"`
	ExpeditionID int32  `db:"expedition_id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	Description  string `db:"description"`
}
