package data

type Member struct {
	TeamID int32  `db:"team_id"`
	UserID int32  `db:"user_id"`
	Role   string `db:"role"`
}
