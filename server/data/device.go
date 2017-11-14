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
