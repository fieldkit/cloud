package api

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/goadesign/goa"
	"github.com/iancoleman/strcase"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"

	pb "github.com/fieldkit/data-protocol"
)

const (
	MetaTypeName = "meta"
	DataTypeName = "data"
)

type DataControllerOptions struct {
	Config   *ApiConfiguration
	Session  *session.Session
	Database *sqlxcache.DB
}

type DataController struct {
	options DataControllerOptions
	*goa.Controller
}

func NewDataController(ctx context.Context, service *goa.Service, options DataControllerOptions) *DataController {
	return &DataController{
		options:    options,
		Controller: service.NewController("DataController"),
	}
}

func (c *DataController) Process(ctx *app.ProcessDataContext) error {
	log := Logger(ctx).Sugar()

	ir, err := NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	pending, err := ir.QueryPending(ctx)
	if err != nil {
		return err
	}

	recordAdder := backend.NewRecordAdder(c.options.Session, c.options.Database)

	for _, i := range pending {
		log.Infow("pending", "ingestion", i)

		err := recordAdder.WriteRecords(ctx, i)
		if err != nil {
			log.Errorw("error", "error", err)
			err := ir.MarkProcessedHasErrors(ctx, i.ID)
			if err != nil {
				return err
			}
		} else {
			err := ir.MarkProcessedDone(ctx, i.ID)
			if err != nil {
				return err
			}
		}
	}

	log.Infow("done", "elapsed", 0)

	return nil
}

func (c *DataController) Delete(ctx *app.DeleteDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	ir, err := NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("deleting", "ingestion_id", ctx.IngestionID)

	i, err := ir.QueryById(ctx, int64(ctx.IngestionID))
	if err != nil {
		return err
	}
	if i == nil {
		return ctx.NotFound()
	}

	err = p.CanModifyStationByDeviceID(i.DeviceID)
	if err != nil {
		return err
	}

	if err := ir.Delete(ctx, int64(ctx.IngestionID)); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	object, err := backend.GetBucketAndKey(i.URL)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %v", err)
	}

	log.Infow("deleting", "url", i.URL)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(object.Bucket), Key: aws.String(object.Key)})
	if err != nil {
		return fmt.Errorf("Unable to delete object %q from bucket %q, %v", object.Key, object.Bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})

	return ctx.OK([]byte("deleted"))
}

type ProvisionSummaryRow struct {
	Generation []byte
	Type       string
	Blocks     data.Int64Range
}

type BlocksSummaryRow struct {
	Blocks data.Int64Range
}

func (c *DataController) DeviceSummary(ctx *app.DeviceSummaryDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	deviceIdBytes, err := data.DecodeBinaryString(ctx.DeviceID)
	if err != nil {
		return err
	}

	err = p.CanViewStationByDeviceID(deviceIdBytes)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	_ = log

	log.Infow("summary", "device_id", deviceIdBytes)

	provisions := make([]*data.Provision, 0)
	if err := c.options.Database.SelectContext(ctx, &provisions, `SELECT * FROM fieldkit.provision WHERE (device_id = $1)`, deviceIdBytes); err != nil {
		return err
	}

	var rows = make([]ProvisionSummaryRow, 0)
	if err := c.options.Database.SelectContext(ctx, &rows, `SELECT generation, type, range_merge(blocks) AS blocks FROM fieldkit.ingestion WHERE (device_id = $1) GROUP BY type, generation`, deviceIdBytes); err != nil {
		return err
	}

	var blockSummaries = make(map[string]map[string]BlocksSummaryRow)
	for _, p := range provisions {
		blockSummaries[hex.EncodeToString(p.Generation)] = make(map[string]BlocksSummaryRow)
	}

	for _, row := range rows {
		g := hex.EncodeToString(row.Generation)
		blockSummaries[g][row.Type] = BlocksSummaryRow{
			Blocks: row.Blocks,
		}
	}

	provisionVms := make([]*app.DeviceProvisionSummary, len(provisions))
	for i, p := range provisions {
		g := hex.EncodeToString(p.Generation)
		byType := blockSummaries[g]
		provisionVms[i] = &app.DeviceProvisionSummary{
			Generation: g,
			Created:    p.Created,
			Updated:    p.Updated,
			Meta: &app.DeviceMetaSummary{
				First: int(byType[MetaTypeName].Blocks[0]),
				Last:  int(byType[MetaTypeName].Blocks[1]),
			},
			Data: &app.DeviceDataSummary{
				First: int(byType[DataTypeName].Blocks[0]),
				Last:  int(byType[DataTypeName].Blocks[1]),
			},
		}
	}

	return ctx.OK(&app.DeviceDataSummaryResponse{
		Provisions: provisionVms,
	})
}

func (c *DataController) DeviceData(ctx *app.DeviceDataDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	deviceIdBytes, err := data.DecodeBinaryString(ctx.DeviceID)
	if err != nil {
		return err
	}

	err = p.CanViewStationByDeviceID(deviceIdBytes)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	_ = log

	rr, err := NewRecordRepository(c.options.Database)
	if err != nil {
		return err
	}

	pageNumber := 0
	if ctx.PageNumber != nil {
		pageNumber = *ctx.PageNumber
	}

	pageSize := 200
	if ctx.PageSize != nil {
		pageSize = *ctx.PageSize
	}

	page, err := rr.QueryDevice(ctx, ctx.DeviceID, pageNumber, pageSize)
	if err != nil {
		return err
	}

	dataVms := make([]*app.DeviceDataRecord, len(page.Data))
	for i, r := range page.Data {
		data, err := r.GetData()
		if err != nil {
			return err
		}

		coordinates := []float64{}
		if r.Location != nil {
			coordinates = r.Location.Coordinates()
		}

		dataVms[i] = &app.DeviceDataRecord{
			ID:       int(r.ID),
			Time:     r.Time,
			Record:   int(r.Number),
			Meta:     int(r.Meta),
			Location: coordinates,
			Data:     data,
		}
	}

	metaVms := make([]*app.DeviceMetaRecord, len(page.Meta))
	for i, r := range page.Meta {
		data, err := r.GetData()
		if err != nil {
			return err
		}

		metaVms[i] = &app.DeviceMetaRecord{
			ID:     int(r.ID),
			Time:   r.Time,
			Record: int(r.Number),
			Data:   data,
		}
	}

	return ctx.OK(&app.DeviceDataRecordsResponse{
		Meta: metaVms,
		Data: dataVms,
	})
}

type IngestionRepository struct {
	Database *sqlxcache.DB
}

func NewIngestionRepository(database *sqlxcache.DB) (ir *IngestionRepository, err error) {
	return &IngestionRepository{Database: database}, nil
}

func (r *IngestionRepository) QueryById(ctx context.Context, id int64) (i *data.Ingestion, err error) {
	found := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &found, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.id = $1)`, id); err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (r *IngestionRepository) QueryPending(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.errors IS NULL) ORDER BY i.size ASC, i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) QueryAll(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i ORDER BY i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) MarkProcessedHasErrors(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = true, attempted = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) MarkProcessedDone(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = false, completed = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) Delete(ctx context.Context, id int64) (err error) {
	if _, err := r.Database.ExecContext(ctx, `DELETE FROM fieldkit.ingestion WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

type StationRepository struct {
	Database *sqlxcache.DB
}

func NewStationRepository(database *sqlxcache.DB) (rr *StationRepository, err error) {
	return &StationRepository{Database: database}, nil
}

func (r *StationRepository) QueryStationByDeviceID(ctx context.Context, deviceIdBytes []byte) (station *data.Station, err error) {
	station = &data.Station{}
	if err := r.Database.GetContext(ctx, station, "SELECT * FROM fieldkit.station WHERE device_id = $1", deviceIdBytes); err != nil {
		return nil, err
	}
	return station, nil
}

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
	    WHERE (p.device_id = $1)
	    ORDER BY r.time DESC LIMIT $2 OFFSET $3`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	mrs := []*data.MetaRecord{}
	if err := r.Database.SelectContext(ctx, &mrs, `
	    SELECT m.* FROM fieldkit.meta_record AS m WHERE (m.id IN (
		SELECT DISTINCT r.meta FROM fieldkit.data_record AS r JOIN fieldkit.provision AS p ON (r.provision_id = p.id)
		WHERE (p.device_id = $1) LIMIT $2 OFFSET $3
	    ))`, deviceIdBytes, pageSize, pageSize*pageNumber); err != nil {
		return nil, err
	}

	page = &RecordsPage{
		Data: drs,
		Meta: mrs,
	}

	return
}

type JSONDataController struct {
	options DataControllerOptions
	*goa.Controller
}

func NewJSONDataController(ctx context.Context, service *goa.Service, options DataControllerOptions) *JSONDataController {
	return &JSONDataController{
		options:    options,
		Controller: service.NewController("JSONDataController"),
	}
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

func isInternalModule(m *pb.ModuleInfo) bool {
	return m.Name == "random" || m.Name == "diagnostics"
}

func (c *JSONDataController) Get(ctx *app.GetJSONDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	deviceIdBytes, err := data.DecodeBinaryString(ctx.DeviceID)
	if err != nil {
		return err
	}

	err = p.CanViewStationByDeviceID(deviceIdBytes)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	_ = log

	rr, err := NewRecordRepository(c.options.Database)
	if err != nil {
		return err
	}

	pageNumber := 0
	if ctx.PageNumber != nil {
		pageNumber = *ctx.PageNumber
	}

	pageSize := 200
	if ctx.PageSize != nil {
		pageSize = *ctx.PageSize
	}

	page, err := rr.QueryDevice(ctx, ctx.DeviceID, pageNumber, pageSize)
	if err != nil {
		return err
	}

	byMeta := make(map[int64][]*data.DataRecord)
	for _, d := range page.Data {
		if byMeta[d.Meta] == nil {
			byMeta[d.Meta] = make([]*data.DataRecord, 0)
		}
		byMeta[d.Meta] = append(byMeta[d.Meta], d)
	}

	sr, err := NewStationRepository(c.options.Database)
	if err != nil {
		return err
	}

	station, err := sr.QueryStationByDeviceID(ctx, deviceIdBytes)
	if err != nil {
		return err
	}

	versions := make([]*app.JSONDataVersion, 0)
	for _, m := range page.Meta {
		dataRecords := byMeta[m.ID]

		var metaRecord pb.DataRecord
		err := m.Unmarshal(&metaRecord)
		if err != nil {
			return err
		}

		name := metaRecord.Identity.Name
		if name == "" {
			name = station.Name
		}

		modules := make([]*app.JSONDataMetaModule, 0)
		for _, module := range metaRecord.Modules {
			if !isInternalModule(module) {
				sensors := make([]*app.JSONDataMetaSensor, 0)
				for _, sensor := range module.Sensors {
					sensors = append(sensors, &app.JSONDataMetaSensor{
						Name:  sensor.Name,
						Key:   strcase.ToLowerCamel(sensor.Name),
						Units: sensor.UnitOfMeasure,
					})
				}

				modules = append(modules, &app.JSONDataMetaModule{
					Name:         module.Name,
					ID:           hex.EncodeToString(module.Id),
					Manufacturer: int(module.Header.Manufacturer),
					Kind:         int(module.Header.Kind),
					Version:      int(module.Header.Version),
					Sensors:      sensors,
				})
			}
		}

		rows := make([]*app.JSONDataRow, 0)
		for _, d := range dataRecords {
			var dataRecord pb.DataRecord
			err := d.Unmarshal(&dataRecord)
			if err != nil {
				return err
			}

			d := make(map[string]interface{})
			for mi, sg := range dataRecord.Readings.SensorGroups {
				if sg.Module == 255 {
					continue
				}
				if mi >= len(metaRecord.Modules) {
					log.Warnw("module index range", "module_index", mi, "modules", len(metaRecord.Modules), "meta", metaRecord.Modules)
					continue
				}
				module := metaRecord.Modules[mi]
				if !isInternalModule(module) {
					for si, r := range sg.Readings {
						sensor := module.Sensors[si]
						key := strcase.ToLowerCamel(sensor.Name)
						d[key] = r.Value
					}
				}
			}

			l := dataRecord.Readings.Location
			rows = append(rows, &app.JSONDataRow{
				Time:     int(dataRecord.Readings.Time),
				Location: getLocation(l),
				D:        d,
			})
		}

		if len(rows) > 0 {
			versions = append(versions, &app.JSONDataVersion{
				Meta: &app.JSONDataMeta{
					Station: &app.JSONDataMetaStation{
						ID:      hex.EncodeToString(metaRecord.Metadata.DeviceId),
						Name:    name,
						Modules: modules,
						Firmware: &app.JSONDataMetaStationFirmware{
							Git:   metaRecord.Metadata.Firmware.Git,
							Build: metaRecord.Metadata.Firmware.Build,
						},
					},
				},
				Data: rows,
			})
		}
	}

	return ctx.OK(&app.JSONDataResponse{
		Versions: versions,
	})
}
