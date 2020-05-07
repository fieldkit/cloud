package data

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type StationActivity struct {
	ID        int64     `db:"id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	StationID int32     `db:"station_id" json:"-"`
}

type StationDeployed struct {
	StationActivity
	DeployedAt time.Time `db:"deployed_at"`
	Location   *Location `db:"location"`
}

type StationIngestion struct {
	StationActivity
	UploaderID      int32 `db:"uploader_id"`
	DataIngestionID int64 `db:"data_ingestion_id"`
	DataRecords     int64 `db:"data_records"`
	Errors          bool  `db:"errors"`
}

type StationDeployedWM struct {
	StationActivity
	DeployedAt time.Time `db:"deployed_at" json:"deployed_at"`
	Location   Location  `db:"location" json:"location"`
}

type AuthorOrUploaderWM struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type IngestedWM struct {
	ID      int64 `json:"id"`
	Records int64 `json:"records"`
}

type StationIngestionWM struct {
	StationActivity
	Uploader AuthorOrUploaderWM `json:"uploader"`
	Data     IngestedWM         `json:"data"`
	Errors   bool               `json:"errors" db:"errors"`
}

func ScanStationIngestionWM(rows *sqlx.Rows) (interface{}, error) {
	v := &StationIngestionWM{}
	if err := rows.Scan(&v.ID, &v.CreatedAt, &v.StationID, &v.Data.ID, &v.Data.Records, &v.Errors, &v.Uploader.ID, &v.Uploader.Name); err != nil {
		return nil, err
	}

	return v, nil
}

type ProjectActivity struct {
	ID        int64     `db:"id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	ProjectID int32     `db:"project_id" json:"-"`
}

type ProjectUpdate struct {
	ProjectActivity
	AuthorID int32  `db:"author_id" json:"-"`
	Body     string `db:"body" json:"body"`
}

type ProjectStationActivity struct {
	ProjectActivity
	StationActivityID int64 `db:"station_activity_id" json:"-"`
}

type ProjectUpdateWM struct {
	ProjectActivity
	Author AuthorOrUploaderWM `json:"author"`
	Body   string             `db:"body" json:"body"`
}

func ScanProjectUpdateWM(rows *sqlx.Rows) (interface{}, error) {
	v := &ProjectUpdateWM{}
	if err := rows.Scan(&v.ID, &v.CreatedAt, &v.ProjectID, &v.Body, &v.Author.ID, &v.Author.Name); err != nil {
		return nil, err
	}

	return v, nil
}
