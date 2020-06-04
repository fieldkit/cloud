package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"

	pb "github.com/fieldkit/data-protocol"
)

func isInternalModule(m *pb.ModuleInfo) bool {
	if m.Flags&META_INTERNAL_MASK == META_INTERNAL_MASK {
		return true
	}
	return m.Name == "random" || m.Name == "diagnostics"
}

type VersionRepository struct {
	db *sqlxcache.DB
}

func NewVersionRepository(db *sqlxcache.DB) (rr *VersionRepository, err error) {
	return &VersionRepository{db: db}, nil
}

func (r *VersionRepository) QueryDevice(ctx context.Context, deviceID string, deviceIdBytes []byte, internal bool, pageNumber, pageSize int) (versions []*Version, err error) {
	log := Logger(ctx).Sugar()

	sr, err := NewStationRepository(r.db)
	if err != nil {
		return nil, err
	}

	rr, err := NewRecordRepository(r.db)
	if err != nil {
		return nil, err
	}

	log.Infow("querying", "device_id", deviceID, "page_number", pageNumber, "page_size", pageSize, "internal", internal)

	station, err := sr.QueryStationByDeviceID(ctx, deviceIdBytes)
	if err != nil {
		return nil, err
	}

	page, err := rr.QueryDevice(ctx, deviceID, pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	mf := NewMetaFactory()

	byMeta := make(map[int64][]*data.DataRecord)
	for _, dbDataRecord := range page.Data {
		metaRecordID := dbDataRecord.MetaRecordID
		_, err := mf.Add(ctx, page.Meta[metaRecordID])
		if err != nil {
			return nil, err
		}
		if byMeta[metaRecordID] == nil {
			byMeta[metaRecordID] = make([]*data.DataRecord, 0)
		}
		byMeta[metaRecordID] = append(byMeta[metaRecordID], dbDataRecord)
	}

	log.Infow("querying", "station_id", station.ID, "station_name", station.Name, "data_records", len(page.Data), "meta_records", len(byMeta))

	versions = make([]*Version, 0)
	for _, versionMeta := range mf.InOrder() {
		dataRecords := byMeta[versionMeta.ID]
		rows := make([]*DataRow, 0)
		for _, d := range dataRecords {
			var dataRecord pb.DataRecord
			err := d.Unmarshal(&dataRecord)
			if err != nil {
				return nil, err
			}

			row, err := mf.Resolve(ctx, d, false)
			if err != nil {
				return nil, err
			}

			if row != nil && row.Record != nil {
				rows = append(rows, row.Record.ToDataRow())
			}
		}

		if len(rows) > 0 {
			versions = append(versions, &Version{
				Meta: versionMeta,
				Data: rows,
			})
		} else {
			log.Infow("empty version", "meta_id", versionMeta.ID)
		}
	}

	return
}
