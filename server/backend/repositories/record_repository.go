package repositories

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fieldkit/cloud/server/common/errors"
	"github.com/fieldkit/cloud/server/common/sqlxcache"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/data"
)

type RecordRepository struct {
	db        *sqlxcache.DB
	metaCache map[int64]*data.MetaRecord
	saveData  bool
}

func NewRecordRepository(db *sqlxcache.DB, saveData bool) *RecordRepository {
	return &RecordRepository{
		db:        db,
		metaCache: make(map[int64]*data.MetaRecord),
		saveData:  saveData,
	}
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
	if err := r.db.SelectContext(ctx, &drs, `
	    SELECT
			r.id, r.provision_id, r.time, r.time, r.number, r.meta_record_id, ST_AsBinary(r.location) AS location, r.raw, r.pb
		FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
	    WHERE (p.device_id = $1)
	    ORDER BY r.time DESC LIMIT $2 OFFSET $3`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	mrs := []*data.MetaRecord{}
	if err := r.db.SelectContext(ctx, &mrs, `
	    SELECT
			m.id, m.provision_id, m.time, m.number, m.raw, m.pb
		FROM fieldkit.meta_record AS m WHERE (m.id IN (
	      SELECT DISTINCT q.meta_record_id FROM (
			SELECT r.meta_record_id, r.time FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id) WHERE (p.device_id = $1) ORDER BY r.time DESC LIMIT $2 OFFSET $3
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

func (r *RecordRepository) AddMetaRecord(ctx context.Context, p *data.Provision, i *data.Ingestion, recordNumber int64, dr *pb.DataRecord, pb []byte) (*data.MetaRecord, error) {
	metaTime := i.Time

	metaRecord := &data.MetaRecord{
		ProvisionID: p.ID,
		Time:        metaTime,
		Number:      recordNumber,
		PB:          pb,
	}

	// TODO Sanitize
	if err := metaRecord.SetData(dr); err != nil {
		return nil, fmt.Errorf("error setting meta json: %w", err)
	}

	if err := r.db.NamedGetContext(ctx, metaRecord, `
		INSERT INTO fieldkit.meta_record (provision_id, time, number, raw, pb) VALUES (:provision_id, :time, :number, :raw, :pb)
		ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw, pb = EXCLUDED.pb
		RETURNING id
		`, metaRecord); err != nil {
		return nil, err
	}

	return metaRecord, nil
}

func (r *RecordRepository) AddSignedMetaRecord(ctx context.Context, p *data.Provision, i *data.Ingestion, sr *pb.SignedRecord, dr *pb.DataRecord, pb []byte) (*data.MetaRecord, error) {
	metaTime := i.Time

	if sr.Time > 0 {
		metaTime = time.Unix(sr.Time, 0).UTC()
	}

	metaRecord := &data.MetaRecord{
		ProvisionID: p.ID,
		Time:        metaTime,
		Number:      int64(sr.Record),
		PB:          pb,
	}

	// TODO Sanitize
	if err := metaRecord.SetData(dr); err != nil {
		return nil, fmt.Errorf("error setting meta json: %w", err)
	}

	if err := r.db.NamedGetContext(ctx, metaRecord, `
		INSERT INTO fieldkit.meta_record (provision_id, time, number, raw, pb) VALUES (:provision_id, :time, :number, :raw, :pb)
		ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, raw = EXCLUDED.raw, pb = EXCLUDED.pb
		RETURNING id
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
	if record, ok := r.metaCache[number]; ok {
		return record, nil
	}

	records := []*data.MetaRecord{}
	if err := r.db.SelectContext(ctx, &records, `
		SELECT r.* FROM fieldkit.meta_record AS r WHERE r.provision_id = $1 AND r.number = $2
		`, provisionId, number); err != nil {
		return nil, err
	}

	if len(records) != 1 {
		return nil, nil
	}

	r.metaCache[number] = records[0]

	return records[0], nil
}

func acceptableJson(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func hasUnaccecptableJsonNumbers(dr *pb.DataRecord) bool {
	for _, sg := range dr.Readings.SensorGroups {
		for _, sr := range sg.Readings {
			if !acceptableJson(float64(sr.GetCalibratedValue())) {
				return true
			}
		}
	}

	return false
}

func prepareForMarshalToJson(dr *pb.DataRecord) *pb.DataRecord {
	if !hasUnaccecptableJsonNumbers(dr) {
		return dr
	}

	for _, sg := range dr.Readings.SensorGroups {
		newReadings := make([]*pb.SensorAndValue, 0, len(sg.Readings))
		for _, sr := range sg.Readings {
			if acceptableJson(float64(sr.GetCalibratedValue())) {
				newReadings = append(newReadings, sr)
			}
		}
		sg.Readings = newReadings
	}

	return dr
}

func (r *RecordRepository) AddDataRecord(ctx context.Context, p *data.Provision, i *data.Ingestion, dr *pb.DataRecord, pb []byte) (*data.DataRecord, *data.MetaRecord, error) {
	if dr.Readings == nil {
		return nil, nil, errors.Structured("malformed record")
	}

	metaRecord, err := r.findMeta(ctx, p.ID, int64(dr.Readings.Meta))
	if err != nil {
		return nil, nil, err
	}
	if metaRecord == nil {
		log := Logger(ctx).Sugar()
		log.Warnw("meta-missing", "meta_record_number", dr.Readings.Meta, "data_record_number", dr.Readings.Reading, "provision_id", p.ID)
		return nil, nil, errors.Structured("meta-missing", "meta_record_number", dr.Readings.Meta, "data_record_number", dr.Readings.Reading, "provision_id", p.ID)
	}

	dataTime := i.Time
	if dr.Readings.Time > 0 {
		dataTime = time.Unix(int64(dr.Readings.Time), 0).UTC()
	}

	location, err := r.findLocation(dr)
	if err != nil {
		return nil, nil, err
	}

	dataRecord := &data.DataRecord{
		ProvisionID:  p.ID,
		Time:         dataTime,
		Number:       int64(dr.Readings.Reading),
		MetaRecordID: metaRecord.ID,
		Location:     location,
		PB:           pb,
	}

	if err := dataRecord.SetData(prepareForMarshalToJson(dr)); err != nil {
		return nil, nil, fmt.Errorf("error setting data json: %w", err)
	}

	if r.saveData {
		if err := r.db.NamedGetContext(ctx, dataRecord, `
			INSERT INTO fieldkit.data_record (provision_id, time, number, meta_record_id, location, raw, pb)
			VALUES (:provision_id, :time, :number, :meta_record_id, ST_SetSRID(ST_GeomFromText(:location), 4326), :raw, :pb)
			ON CONFLICT (provision_id, number) DO UPDATE SET number = EXCLUDED.number, time = EXCLUDED.time, location = EXCLUDED.location, raw = EXCLUDED.raw, pb = EXCLUDED.pb
			RETURNING id
			`, dataRecord); err != nil {
			return nil, nil, err
		}
	}

	return dataRecord, metaRecord, nil
}
