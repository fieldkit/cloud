package data

type Administrator struct {
	ProjectID int32 `db:"project_id"`
	UserID    int32 `db:"user_id"`
}
