package data

type Input struct {
	ID           int32  `db:"id,omitempty"`
	ExpeditionID int32  `db:"expedition_id"`
	Type         string `db:"type"`
	Name         string `db:"name"`
	Active       bool   `db:"active"`
}
