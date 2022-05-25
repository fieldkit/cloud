package data

type ProjectAttribute struct {
	ID        int64  `db:"id"`
	ProjectID int32  `db:"project_id"`
	Name      string `db:"name"`
}

type StationProjectAttribute struct {
	ID          int64  `db:"id"`
	StationID   int32  `db:"station_id"`
	ProjectID   int32  `db:"project_id"`
	StringValue string `db:"string_value"`
}

type StationProjectNamedAttribute struct {
	ID          int64  `db:"id"`
	StationID   int32  `db:"station_id"`
	AttributeID int64  `db:"attribute_id"`
	StringValue string `db:"string_value"`
	ProjectID   int32  `db:"project_id"`
	Name        string `db:"name"`
}
