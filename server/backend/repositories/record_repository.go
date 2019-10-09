package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

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
	Meta []*data.MetaRecord
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
	    WHERE (p.device_id = $1) AND r.time < '10/8/2019'
	    ORDER BY r.time DESC LIMIT $2 OFFSET $3`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	mrs := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &mrs, `
	    SELECT m.* FROM fieldkit.meta_record AS m WHERE (m.id IN (
		SELECT DISTINCT r.meta FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id = $1) AND r.time < '10/8/2019' LIMIT $2 OFFSET $3
	    ))`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	page = &RecordsPage{
		Data: drs,
		Meta: mrs,
	}

	return
}
