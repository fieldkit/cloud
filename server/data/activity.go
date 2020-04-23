package data

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type StationActivity struct {
	ID        int64     `db:"id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	StationID int64     `db:"station_id" json:"-"`
}

type StationDeployed struct {
	StationActivity
	DeployedAt time.Time `db:"deployed_at"`
	Location   Location  `db:"location"`
}

type StationIngestion struct {
	StationActivity
	UploaderID      int64 `db:"uploader_id"`
	DataIngestionID int64 `db:"data_ingestion_id"`
	DataRecords     int64 `db:"data_records"`
	Errors          bool  `db:"errors"`
}

type StationDeployedWM struct {
	StationActivity
	DeployedAt time.Time `db:"deployed_at" json:"deployed_at"`
	Location   Location  `db:"location" json:"location"`
}

type Uploader struct {
	ID   int64  `db:"uploader_id" json:"id"`
	Name string `db:"uploader_name" json:"name"`
}

type IngestedWM struct {
	ID      int64 `json:"id"`
	Records int64 `json:"records"`
}

type StationIngestionWM struct {
	StationActivity
	Uploader Uploader   `json:"uploader"`
	Data     IngestedWM `json:"data"`
	Errors   bool       `db:"errors" json:"errors"`
}

func ScanStationIngestionWM(rows *sqlx.Rows) (interface{}, error) {
	v := &StationIngestionWM{}
	if err := rows.Scan(&v.ID, &v.CreatedAt, &v.StationID, &v.Data.ID, &v.Data.Records, &v.Errors, &v.Uploader.ID, &v.Uploader.Name); err != nil {
		return nil, err
	}
	return v, nil
}
