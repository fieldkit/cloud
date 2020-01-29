package repositories

import (
	"context"
	"encoding/hex"
	_ "time"

	"github.com/iancoleman/strcase"

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
	ID      int
	Station *DataMetaStation
}

type DataMetaSensor struct {
	Name  string
	Key   string
	Units string
}

type DataMetaModule struct {
	Manufacturer int
	Kind         int
	Version      int
	Name         string
	ID           string
	Sensors      []*DataMetaSensor
}

type DataMetaStation struct {
	ID       string
	Name     string
	Firmware *DataMetaStationFirmware
	Modules  []*DataMetaModule
}

type DataMetaStationFirmware struct {
	Git   string
	Build string
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

	log.Infow("querying", "station_id", station.ID, "station_name", station.Name)

	versions = make([]*Version, 0)
	for _, m := range metaOrder {
		var metaRecord pb.DataRecord
		err := m.Unmarshal(&metaRecord)
		if err != nil {
			return nil, err
		}

		name := metaRecord.Identity.Name
		if name == "" {
			name = station.Name
		}

		modules := make([]*DataMetaModule, 0)
		for _, module := range metaRecord.Modules {
			if internal || !isInternalModule(module) {
				sensors := make([]*DataMetaSensor, 0)
				for _, sensor := range module.Sensors {
					sensors = append(sensors, &DataMetaSensor{
						Name:  sensor.Name,
						Key:   strcase.ToLowerCamel(sensor.Name),
						Units: sensor.UnitOfMeasure,
					})
				}

				modules = append(modules, &DataMetaModule{
					Name:         module.Name,
					ID:           hex.EncodeToString(module.Id),
					Manufacturer: int(module.Header.Manufacturer),
					Kind:         int(module.Header.Kind),
					Version:      int(module.Header.Version),
					Sensors:      sensors,
				})
			}
		}

		skipped := 0

		dataRecords := byMeta[m.ID]
		rows := make([]*DataRow, 0)
		for _, d := range dataRecords {
			var dataRecord pb.DataRecord
			err := d.Unmarshal(&dataRecord)
			if err != nil {
				return nil, err
			}

			data := make(map[string]interface{})
			for mi, sg := range dataRecord.Readings.SensorGroups {
				if sg.Module == 255 {
					skipped += 1
					continue
				}
				if mi >= len(metaRecord.Modules) {
					log.Warnw("module index range", "module_index", mi, "modules", len(metaRecord.Modules), "meta", metaRecord.Modules)
					continue
				}
				module := metaRecord.Modules[mi]
				if internal || !isInternalModule(module) {
					for si, r := range sg.Readings {
						if r != nil && si < len(module.Sensors) {
							sensor := module.Sensors[si]
							key := strcase.ToLowerCamel(sensor.Name)
							data[key] = r.Value
						}
					}
				}
			}

			rows = append(rows, &DataRow{
				ID:       int(d.ID),
				Time:     int(dataRecord.Readings.Time),
				Location: getLocation(dataRecord.Readings.Location),
				D:        data,
			})
		}

		if len(rows) > 0 {
			versions = append(versions, &Version{
				Meta: &VersionMeta{
					ID: int(m.ID),
					Station: &DataMetaStation{
						ID:      hex.EncodeToString(metaRecord.Metadata.DeviceId),
						Name:    name,
						Modules: modules,
						Firmware: &DataMetaStationFirmware{
							Git:   metaRecord.Metadata.Firmware.Git,
							Build: metaRecord.Metadata.Firmware.Build,
						},
					},
				},
				Data: rows,
			})
		} else {
			log.Infow("empty version", "meta_id", m.ID)
		}
	}

	return
}
