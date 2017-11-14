package data

type Device struct {
	InputID int64  `db:"input_id"`
	Key     string `db:"key"`
	Token   string `db:"token"`
}

type DeviceInput struct {
	Input
	Device
}

type DeviceSchema struct {
	ID       int64  `db:"id"`
	DeviceID int64  `db:"device_id"`
	SchemaID int64  `db:"schema_id"`
	Key      string `db:"key"`
}

type DeviceJSONSchema struct {
	DeviceSchema
	RawSchema
}
