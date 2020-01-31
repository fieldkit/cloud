package repositories

import (
	"context"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/data"

	pb "github.com/fieldkit/data-protocol"
)

func isInternalModule(m *pb.ModuleInfo) bool {
	return m.Name == "random" || m.Name == "diagnostics"
}

func getLocation(l *pb.DeviceLocation) []float64 {
	if l == nil {
		return nil
	}
	if l.Longitude > 180 || l.Longitude < -180 {
		return nil
	}
	return []float64{
		float64(l.Longitude),
		float64(l.Latitude),
		float64(l.Altitude),
	}
}

type Version struct {
	Meta *VersionMeta
	Data []*DataRow
}

type VersionMeta struct {
	ID      int64
	Station *DataMetaStation
}

type DataMetaSensor struct {
	Name     string
	Key      string
	Units    string
	Internal bool
}

type DataMetaModule struct {
	Manufacturer int
	Kind         int
	Version      int
	Name         string
	ID           string
	Sensors      []*DataMetaSensor
	Internal     bool
}

type DataMetaStation struct {
	ID       string
	Name     string
	Firmware *DataMetaStationFirmware
	Modules  []*DataMetaModule
}

type DataMetaStationFirmware struct {
	Version   string
	Build     string
	Number    string
	Timestamp uint64
	Hash      string
}

type DataRow struct {
	ID       int
	Time     int
	Location []float64
	D        map[string]interface{}
}

type VersionRepository struct {
	Database *sqlxcache.DB
}

func NewVersionRepository(database *sqlxcache.DB) (rr *VersionRepository, err error) {
	return &VersionRepository{Database: database}, nil
}

func (r *VersionRepository) QueryDevice(ctx context.Context, deviceID string, deviceIdBytes []byte, internal bool, pageNumber, pageSize int) (versions []*Version, err error) {
	log := Logger(ctx).Sugar()

	sr, err := NewStationRepository(r.Database)
	if err != nil {
		return nil, err
	}

	rr, err := NewRecordRepository(r.Database)
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

	metaOrder := make([]*data.MetaRecord, 0)
	byMeta := make(map[int64][]*data.DataRecord)
	for _, d := range page.Data {
		if byMeta[d.Meta] == nil {
			byMeta[d.Meta] = make([]*data.DataRecord, 0)
			metaOrder = append(metaOrder, page.Meta[d.Meta])
		}
		byMeta[d.Meta] = append(byMeta[d.Meta], d)
	}

	mf := NewMetaFactory()

	log.Infow("querying", "station_id", station.ID, "station_name", station.Name)

	versions = make([]*Version, 0)
	for _, m := range metaOrder {
		versionMeta, err := mf.Add(m)
		if err != nil {
			return nil, err
		}

		dataRecords := byMeta[m.ID]
		rows := make([]*DataRow, 0)
		for _, d := range dataRecords {
			var dataRecord pb.DataRecord
			err := d.Unmarshal(&dataRecord)
			if err != nil {
				return nil, err
			}

			row, err := mf.Resolve(d)
			if err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}

		if len(rows) > 0 {
			versions = append(versions, &Version{
				Meta: versionMeta,
				Data: rows,
			})
		} else {
			log.Infow("empty version", "meta_id", m.ID)
		}
	}

	return
}
