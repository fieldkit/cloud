package messages

type IngestStation struct {
	StationID int32 `json:"station_id"`
	Verbose   bool  `json:"verbose"`
	UserID    int32 `json:"user_id"`
}

type StationIngested struct {
	StationID int32 `json:"station_id"`
	UserID    int32 `json:"user_id"`
}
