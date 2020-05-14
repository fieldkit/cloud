package repositories

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/conservify/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type RecordRepository struct {
	Database *sqlxcache.DB
}

func NewRecordRepository(database *sqlxcache.DB) (rr *RecordRepository, err error) {
	return &RecordRepository{Database: database}, nil
}

type RecordsPage struct {
	Data []*data.DataRecord
	Meta map[int64]*data.MetaRecord
}

func (r *RecordRepository) QueryDevice(ctx context.Context, deviceId string, pageNumber, pageSize int) (page *RecordsPage, err error) {
	log := Logger(ctx).Sugar()

	deviceIdBytes, err := data.DecodeBinaryString(deviceId)
	if err != nil {
		return nil, err
	}

	log.Infow("querying", "device_id", deviceIdBytes, "page_number", pageNumber, "page_size", pageSize)

	drs := []*data.DataRecord{}
	if err := r.Database.SelectContext(ctx, &drs, `
	    SELECT r.id, r.provision_id, r.time, r.time, r.number, r.meta, ST_AsBinary(r.location) AS location, r.raw FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
	    WHERE (p.device_id = $1)
	    ORDER BY r.time DESC LIMIT $2 OFFSET $3`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	mrs := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &mrs, `
	    SELECT m.* FROM fieldkit.meta_record AS m WHERE (m.id IN (
	      SELECT DISTINCT q.meta FROM (
			SELECT r.meta, r.time FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id) WHERE (p.device_id = $1) ORDER BY r.time DESC LIMIT $2 OFFSET $3
	      ) AS q
	    ))`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	metas := make(map[int64]*data.MetaRecord)
	for _, m := range mrs {
		metas[m.ID] = m
	}

	page = &RecordsPage{
		Data: drs,
		Meta: metas,
	}

	return
}

func (r *RecordRepository) AddMetaRecord(ctx context.Context, p *data.Provision, i *data.Ingestion, sr *pb.SignedRecord, dr *pb.DataRecord) (*data.MetaRecord, error) {
	metaTime := i.Time

	if sr.Time > 0 {
		metaTime = time.Unix(sr.Time, 0)
	}

	metaRecord := &data.MetaRecord{
		ProvisionID: p.ID,
		Time:        metaTime,
		Number:      int64(sr.Record),
	}

	if err := metaRecord.SetData(dr); err != nil {
		return nil, fmt.Errorf("error setting meta json: %v", err)
	}

	if err := r.Database.NamedGetContext(ctx, metaRecord, `
		INSERT INTO fieldkit.meta_record (provision_id, time, number, raw) VALUES (:provision_id, :time, :number, :raw)
		ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw
		RETURNING *
		`, metaRecord); err != nil {
		return nil, err
	}

	return metaRecord, nil
}

func (r *RecordRepository) findLocation(dataRecord *pb.DataRecord) (l *data.Location, err error) {
	if dataRecord.Readings == nil || dataRecord.Readings.Location == nil {
		return nil, err
	}
	location := dataRecord.Readings.Location
	lat := float64(location.Latitude)
	lon := float64(location.Longitude)
	altitude := float64(location.Altitude)

	if lat > 90 || lat < -90 || lon > 180 || lon < -180 {
		return nil, err
	}

	if lat == 0 && lon == 0 {
		return nil, err
	}

	l = data.NewLocation([]float64{lon, lat, altitude})
	return
}

func (r *RecordRepository) findMeta(ctx context.Context, provisionId, number int64) (*data.MetaRecord, error) {
	records := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &records, `
		SELECT r.* FROM fieldkit.meta_record AS r WHERE r.provision_id = $1 AND r.number = $2
		`, provisionId, number); err != nil {
		return nil, err
	}

	if len(records) != 1 {
		return nil, nil
	}

	return records[0], nil
}

func hasNaNs(dr *pb.DataRecord) bool {
	for _, sg := range dr.Readings.SensorGroups {
		for _, sr := range sg.Readings {
			if math.IsNaN(float64(sr.Value)) {
				return true
			}

		}
	}
	return false
}

func prepareForMarshalToJson(dr *pb.DataRecord) *pb.DataRecord {
	if !hasNaNs(dr) {
		return dr
	}
	for _, sg := range dr.Readings.SensorGroups {
		newReadings := make([]*pb.SensorAndValue, len(sg.Readings))
		for _, sr := range sg.Readings {
			if !math.IsNaN(float64(sr.Value)) {
				newReadings = append(newReadings, sr)
			}

		}
		sg.Readings = newReadings
	}
	return dr
}

var ErrMalformedRecord = errors.New("malformed record")
var ErrMetaMissing = errors.New("error finding meta record")

func (r *RecordRepository) AddDataRecord(ctx context.Context, p *data.Provision, i *data.Ingestion, dr *pb.DataRecord) (*data.DataRecord, *data.MetaRecord, error) {
	if dr.Readings == nil {
		return nil, nil, ErrMalformedRecord
	}

	metaRecord, err := r.findMeta(ctx, p.ID, int64(dr.Readings.Meta))
	if err != nil {
		return nil, nil, err
	}
	if metaRecord == nil {
		return nil, nil, ErrMetaMissing
	}

	dataTime := i.Time
	if dr.Readings.Time > 0 {
		dataTime = time.Unix(int64(dr.Readings.Time), 0)
	}

	location, err := r.findLocation(dr)
	if err != nil {
		return nil, nil, err
	}

	dataRecord := &data.DataRecord{
		ProvisionID: p.ID,
		Time:        dataTime,
		Number:      int64(dr.Readings.Reading),
		MetaID:      metaRecord.ID,
		Location:    location,
	}

	if err := dataRecord.SetData(prepareForMarshalToJson(dr)); err != nil {
		return nil, nil, fmt.Errorf("error setting data json: %v", err)
	}

	if err := r.Database.NamedGetContext(ctx, dataRecord, `
		INSERT INTO fieldkit.data_record (provision_id, time, number, raw, meta, location)
		VALUES (:provision_id, :time, :number, :raw, :meta, ST_SetSRID(ST_GeomFromText(:location), 4326))
		ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw, location = EXCLUDED.location
		RETURNING id
		`, dataRecord); err != nil {
		return nil, nil, err
	}

	return dataRecord, metaRecord, nil
}
