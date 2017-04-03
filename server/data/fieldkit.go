package data

type FieldkitBinary struct {
	InputID int32    `db:"input_id"`
	ID      int16    `db:"id"`
	Fields  []string `db:"fields"`
}

type FieldkitInput struct {
	Input
}
