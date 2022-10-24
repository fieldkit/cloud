package messages

import "time"

type RefreshAllMaterializedViews struct {
}

type RefreshMaterializedView struct {
	View  string
	Start time.Time
	End   time.Time
}
