package data

type Input struct {
	ID           int32  `db:"id,omitempty"`
	ExpeditionID int32  `db:"expedition_id"`
	TeamID       *int32 `db:"team_id,omitempty"`
	UserID       *int32 `db:"user_id,omitempty"`
}
