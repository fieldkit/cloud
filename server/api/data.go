package api

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/jobs"
	"github.com/fieldkit/cloud/server/messages"
)

const (
	MetaTypeName = "meta"
	DataTypeName = "data"
)

type DataControllerOptions struct {
	Config    *ApiConfiguration
	Session   *session.Session
	Database  *sqlxcache.DB
	Publisher jobs.MessagePublisher
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

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	pending, err := ir.QueryPending(ctx)
	if err != nil {
		return err
	}

	for _, i := range pending {
		log.Infow("queueing", "ingestion_id", i.ID)
		c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
			Time: i.Time,
			ID:   i.ID,
			URL:  i.URL,
		})
	}

	return nil
}

func (c *DataController) ProcessIngestion(ctx *app.ProcessIngestionDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("processing", "ingestion_id", ctx.IngestionID)

	i, err := ir.QueryByID(ctx, int64(ctx.IngestionID))
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

	c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
		Time: i.Time,
		ID:   i.ID,
		URL:  i.URL,
	})

	return ctx.OK([]byte("queued"))
}

func (c *DataController) Delete(ctx *app.DeleteDataContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	log := Logger(ctx).Sugar()

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("deleting", "ingestion_id", ctx.IngestionID)

	i, err := ir.QueryByID(ctx, int64(ctx.IngestionID))
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
	if err != nil {
		return err
	}

	if err := ir.Delete(ctx, int64(ctx.IngestionID)); err != nil {
		return err
	}

	return ctx.OK([]byte("deleted"))
}

type ProvisionSummaryRow struct {
	Generation []byte
	Type       string
	Blocks     data.Int64Range
}

type BlocksSummaryRow struct {
	Size   int64
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
	if err := c.options.Database.SelectContext(ctx, &provisions, `SELECT * FROM fieldkit.provision WHERE (device_id = $1) ORDER BY updated DESC`, deviceIdBytes); err != nil {
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
				Size:  int(0),
				First: int(byType[MetaTypeName].Blocks[0]),
				Last:  int(byType[MetaTypeName].Blocks[1]),
			},
			Data: &app.DeviceDataSummary{
				Size:  int(0),
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

	rr, err := repositories.NewRecordRepository(c.options.Database)
	if err != nil {
		return err
	}

	pageNumber := 0
	if ctx.Page != nil {
		pageNumber = *ctx.Page
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

func (c *JSONDataController) Summary(ctx *app.SummaryJSONDataContext) error {
	return ctx.OK(nil)
}

func JSONDataRowsType(dm []*repositories.DataRow) []*app.JSONDataRow {
	wm := make([]*app.JSONDataRow, len(dm))
	for i, r := range dm {
		wm[i] = &app.JSONDataRow{
			ID:       r.ID,
			Time:     r.Time,
			Location: r.Location,
			D:        r.D,
		}
	}
	return wm
}

func JSONDataMetaType(dm *repositories.VersionMeta) *app.JSONDataMeta {
	modules := make([]*app.JSONDataMetaModule, len(dm.Station.Modules))
	for i, m := range dm.Station.Modules {
		sensors := make([]*app.JSONDataMetaSensor, len(m.Sensors))

		for j, s := range m.Sensors {
			sensors[j] = &app.JSONDataMetaSensor{
				Name:  s.Name,
				Key:   s.Key,
				Units: s.Units,
			}
		}

		modules[i] = &app.JSONDataMetaModule{
			Manufacturer: m.Manufacturer,
			Kind:         m.Kind,
			Version:      m.Version,
			ID:           m.ID,
			Name:         m.Name,
			Sensors:      sensors,
		}
	}
	return &app.JSONDataMeta{
		ID: dm.ID,
		Station: &app.JSONDataMetaStation{
			ID:   dm.Station.ID,
			Name: dm.Station.Name,
			Firmware: &app.JSONDataMetaStationFirmware{
				Git:   dm.Station.Firmware.Git,
				Build: dm.Station.Firmware.Build,
			},
			Modules: modules,
		},
	}
}

func JSONDataResponseType(dm []*repositories.Version) []*app.JSONDataVersion {
	wm := make([]*app.JSONDataVersion, len(dm))
	for i, r := range dm {
		wm[i] = &app.JSONDataVersion{
			Meta: JSONDataMetaType(r.Meta),
			Data: JSONDataRowsType(r.Data),
		}
	}
	return nil
}

func (c *JSONDataController) Get(ctx *app.GetJSONDataContext) error {
	log := Logger(ctx).Sugar()

	_ = log

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

	pageNumber := 0
	if ctx.Page != nil {
		pageNumber = *ctx.Page
	}

	pageSize := 200
	if ctx.PageSize != nil {
		pageSize = *ctx.PageSize
	}

	internal := false
	if ctx.Internal != nil {
		internal = *ctx.Internal
	}

	vr, err := repositories.NewVersionRepository(c.options.Database)
	if err != nil {
		return err
	}

	versions, err := vr.QueryDevice(ctx, ctx.DeviceID, deviceIdBytes, internal, pageNumber, pageSize)
	if err != nil {
		return err
	}

	return ctx.OK(&app.JSONDataResponse{
		Versions: JSONDataResponseType(versions),
	})
}

type JSONLine struct {
	Time     int                    `form:"time" json:"time" yaml:"time" xml:"time"`
	Location []float64              `form:"location" json:"location" yaml:"location" xml:"location"`
	ID       int                    `form:"id" json:"id" yaml:"id" xml:"id"`
	Meta     int                    `form:"meta" json:"meta" yaml:"meta" xml:"meta"`
	D        map[string]interface{} `form:"d" json:"d" yaml:"d" xml:"d"`
}

func (c *JSONDataController) GetLines(ctx *app.GetLinesJSONDataContext) error {
	log := Logger(ctx).Sugar()

	_ = log

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

	pageNumber := 0
	if ctx.Page != nil {
		pageNumber = *ctx.Page
	}

	pageSize := 200
	if ctx.PageSize != nil {
		pageSize = *ctx.PageSize
	}

	internal := false
	if ctx.Internal != nil {
		internal = *ctx.Internal
	}

	vr, err := repositories.NewVersionRepository(c.options.Database)
	if err != nil {
		return err
	}

	versions, err := vr.QueryDevice(ctx, ctx.DeviceID, deviceIdBytes, internal, pageNumber, pageSize)
	if err != nil {
		return err
	}

	ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	ctx.ResponseData.WriteHeader(200)

	writer := bufio.NewWriter(ctx.ResponseData)

	for _, version := range versions {
		for _, row := range version.Data {
			line := JSONLine{
				Time:     row.Time,
				Location: row.Location,
				ID:       row.ID,
				Meta:     version.Meta.ID,
				D:        row.D,
			}
			bytes, err := json.Marshal(line)
			if err != nil {
				return err
			}
			writer.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		}
	}

	writer.Flush()

	return nil
}
