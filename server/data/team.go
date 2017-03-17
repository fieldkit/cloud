package data

type Team struct {
	ID           int32  `db:"id,omitempty"`
	ExpeditionID int32  `db:"expedition_id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	Description  string `db:"description"`
}
