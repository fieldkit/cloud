package data

type Station struct {
	ID          int32  `db:"id,omitempty"`
	Name        string `db:"name"`
	UserID 			int32  `db:"user_id,omitempty"`
}

type StationLog struct {
	ID					int32  `db:"id,omitempty"`
	StationId		int32	 `db:"id,omitempty"`
	Body 				string `db:"body"`
	Timestamp		string `db:"timestamp"`
}
