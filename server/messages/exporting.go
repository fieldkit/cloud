package messages

type ExportData struct {
	ID     int64  `json:"id"`
	UserID int32  `json:"user_id"`
	Token  string `json:"token"`
	Format string `json:"format"`
}
