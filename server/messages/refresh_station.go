package messages

type RefreshStation struct {
	StationID  int32
	Completely bool
	Verbose    bool
	UserID     int32
}
