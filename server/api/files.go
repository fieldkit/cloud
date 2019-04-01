package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type FilesControllerOptions struct {
	Config   *ApiConfiguration
	Session  *session.Session
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type BaseFilesController struct {
	options FilesControllerOptions
	cw      *ConcatenationWorkers
}

type FilesController struct {
	BaseFilesController
	*goa.Controller
}

var (
	ConcatenatedDataSpace  = uuid.MustParse("5554c632-d586-53f2-b2a7-88a63663a3f5")
	ConcatenatedLogsSpace  = uuid.MustParse("1e563c95-7984-4868-abd1-411c4dfc5467")
	ConcatenatedFileSpaces = map[string]uuid.UUID{
		"2": ConcatenatedLogsSpace,
		"4": ConcatenatedDataSpace,
	}
	DataFileTypeIDs = []string{"4"}
	LogFileTypeIDs  = []string{"2", "3"}
	FileTypeNames   = map[string]string{
		"2": "Logs",
		"3": "Logs",
		"4": "Data",
	}
)

func NewFilesController(ctx context.Context, service *goa.Service, options FilesControllerOptions) *FilesController {
	return &FilesController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("FilesController"),
	}
}

func (c *FilesController) listDeviceFiles(ctx context.Context, fileTypeIDs []string, deviceID string, page *int) (*app.DeviceFiles, error) {
	log := Logger(ctx).Sugar()

	log.Infow("Device", "device_id", deviceID, "file_type_ids", fileTypeIDs)

	pageSize := 100
	offset := 0
	if page != nil {
		offset = *page * pageSize
	}

	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.file_id = ANY($1)) AND (s.device_id = $2) ORDER BY time DESC LIMIT $3 OFFSET $4`, pq.StringArray(fileTypeIDs), deviceID, pageSize, offset); err != nil {
		return nil, err
	}

	return DeviceFilesType(c.options.Config, files), nil
}

func (c *FilesController) ListDeviceDataFiles(ctx *app.ListDeviceDataFilesFilesContext) error {
	files, err := c.listDeviceFiles(ctx, DataFileTypeIDs, ctx.DeviceID, ctx.Page)
	if err != nil {
		return err
	}

	return ctx.OK(files)
}

func (c *FilesController) ListDeviceLogFiles(ctx *app.ListDeviceLogFilesFilesContext) error {
	files, err := c.listDeviceFiles(ctx, LogFileTypeIDs, ctx.DeviceID, ctx.Page)
	if err != nil {
		return err
	}

	return ctx.OK(files)
}

func (c *FilesController) ListDevices(ctx *app.ListDevicesFilesContext) error {
	locations := []*data.DeviceStreamLocation{}
	if err := c.options.Database.SelectContext(ctx, &locations,
		`SELECT s.id, s.device_id, s.timestamp, ST_AsBinary(s.location) AS location
		 FROM fieldkit.device_stream_location AS s
                 ORDER BY s.timestamp DESC`); err != nil {
		return err
	}

	devices := []*data.DeviceSummary{}
	if err := c.options.Database.SelectContext(ctx, &devices,
		`SELECT s.device_id,
		    (SELECT stream_id FROM fieldkit.device_stream AS s2 WHERE (s2.device_id = s.device_id) ORDER BY s2.time DESC LIMIT 1) AS last_stream_id,
		    MAX(s.time) AS last_stream_time,
		    COUNT(s.*) AS number_of_files
		 FROM fieldkit.device_stream AS s
                 WHERE s.device_id != ''
                 GROUP BY s.device_id
                 ORDER BY last_stream_time DESC`); err != nil {
		return err
	}

	return ctx.OK(DevicesType(c.options.Config, devices, locations))
}

func (c *FilesController) DeviceInfo(ctx *app.DeviceInfoFilesContext) error {
	ac := c.options.Config

	fr, err := backend.NewFileRepository(c.options.Session, "fk-streams")
	if err != nil {
		return err
	}

	urls := DeviceSummaryUrls(ac, ctx.DeviceID)

	data, err := fr.Info(ctx, urls.Data.ID)
	if err != nil {
		return err
	}

	logs, err := fr.Info(ctx, urls.Logs.ID)
	if err != nil {
		return err
	}

	return ctx.OK(&app.DeviceDetails{
		DeviceID: ctx.DeviceID,
		Urls:     urls,
		Files: &app.ConcatenatedFilesInfo{
			Data: ConcatenatedFileInfo(urls.Data, data),
			Logs: ConcatenatedFileInfo(urls.Logs, logs),
		},
	})
}

func (c *FilesController) File(ctx *app.FileFilesContext) error {
	log := Logger(ctx).Sugar()

	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.stream_id = $1)`, ctx.FileID); err != nil {
		return err
	}

	if len(files) != 1 {
		ac := c.options.Config

		fr, err := backend.NewFileRepository(c.options.Session, "fk-streams")
		if err != nil {
			return err
		}

		fi, err := fr.Info(ctx, ctx.FileID)
		if err != nil {
			return err
		}
		if fi == nil {
			return ctx.NotFound()
		}

		jsonMeta, err := json.Marshal(fi.Meta)
		if err != nil {
			panic(err)
		}

		log.Infow("File", "info", fi)

		return ctx.OK(&app.DeviceFile{
			DeviceID:     fi.DeviceID,
			FileID:       fi.Key,
			FileTypeID:   fi.FileTypeID,
			FileTypeName: FileTypeNames[fi.FileTypeID],
			Time:         fi.LastModified,
			Size:         int(fi.Size),
			URL:          fi.URL,
			Meta:         string(jsonMeta),
			Urls: &app.DeviceFileUrls{
				Csv:  ac.MakeApiUrl("/files/%v/data.csv", fi.Key),
				JSON: ac.MakeApiUrl("/files/%v/data.json", fi.Key),
				Fkpb: ac.MakeApiUrl("/files/%v/data.fkpb", fi.Key),
			},
		})
	}

	return ctx.OK(DeviceFileSummaryType(c.options.Config, files[0]))
}

func (c *FilesController) Csv(ctx *app.CsvFilesContext) error {
	iterator, err := backend.LookupFile(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	exporter := backend.NewSimpleCsvExporter(ctx.ResponseData)

	return backend.ExportAllFiles(ctx, ctx.ResponseData, ctx.Dl, iterator, exporter)
}

func (c *FilesController) JSON(ctx *app.JSONFilesContext) error {
	iterator, err := backend.LookupFile(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	exporter := backend.NewSimpleJsonExporter(ctx.ResponseData)

	return backend.ExportAllFiles(ctx, ctx.ResponseData, ctx.Dl, iterator, exporter)
}

func (c *FilesController) Raw(ctx *app.RawFilesContext) error {
	iterator, err := backend.LookupFile(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	opened, err := iterator.Next(ctx)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.csv", opened.FileID)

	ctx.ResponseData.Header().Set("Content-Type", backend.FkDataBinaryContentType)
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	ctx.ResponseData.WriteHeader(http.StatusOK)

	_, err = io.Copy(ctx.ResponseData, opened.Response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *BaseFilesController) concatenatedDeviceFile(ctx context.Context, responseData *goa.ResponseData, deviceID string, fileTypeIDs []string) (redirection string, err error) {
	space := ConcatenatedFileSpaces[fileTypeIDs[0]]
	deviceStreamID := uuid.NewSHA1(space, []byte(deviceID))

	c.cw.channel <- ConcatenationJob{
		DeviceID:    deviceID,
		FileTypeIDs: fileTypeIDs,
	}

	return c.options.Config.MakeApiUrl("/files/%s", deviceStreamID), nil
}

type DeviceLogsController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceLogsController(ctx context.Context, service *goa.Service, options FilesControllerOptions) *DeviceLogsController {
	cw, err := NewConcatenationWorkers(ctx, options)
	if err != nil {
		panic(err)
	}

	return &DeviceLogsController{
		BaseFilesController: BaseFilesController{
			options: options,
			cw:      cw,
		},
		Controller: service.NewController("DeviceLogsController"),
	}
}

func (c *DeviceLogsController) All(ctx *app.AllDeviceLogsContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, LogFileTypeIDs)
	if err != nil {
		return err
	}

	if url != "" {
		ctx.ResponseData.Header().Set("Location", url)
		return ctx.Found()
	}

	return ctx.NotFound()
}

type DeviceDataController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceDataController(ctx context.Context, service *goa.Service, options FilesControllerOptions) *DeviceDataController {
	cw, err := NewConcatenationWorkers(ctx, options)
	if err != nil {
		panic(err)
	}

	return &DeviceDataController{
		BaseFilesController: BaseFilesController{
			options: options,
			cw:      cw,
		},
		Controller: service.NewController("DeviceDataController"),
	}
}

func (c *DeviceDataController) All(ctx *app.AllDeviceDataContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, DataFileTypeIDs)
	if err != nil {
		return err
	}

	if url != "" {
		ctx.ResponseData.Header().Set("Location", url)
		return ctx.Found()
	}

	return ctx.NotFound()
}

func DeviceFileTypeUrls(ac *ApiConfiguration, kind, deviceID, fileID string) *app.DeviceFileTypeUrls {
	return &app.DeviceFileTypeUrls{
		ID:       fileID,
		Generate: ac.MakeApiUrl("/devices/%v/%s", deviceID, kind),
		Info:     ac.MakeApiUrl("/files/%v", fileID),
		Csv:      ac.MakeApiUrl("/files/%v/data.csv", fileID),
		Fkpb:     ac.MakeApiUrl("/files/%v/data.fkpb", fileID),
	}
}

func DeviceSummaryUrls(ac *ApiConfiguration, deviceID string) *app.DeviceSummaryUrls {
	concatenatedDataID := uuid.NewSHA1(ConcatenatedDataSpace, []byte(deviceID))
	concatenatedLogsID := uuid.NewSHA1(ConcatenatedLogsSpace, []byte(deviceID))

	return &app.DeviceSummaryUrls{
		Details: ac.MakeApiUrl("/devices/%v/data", deviceID),
		Data:    DeviceFileTypeUrls(ac, "data", deviceID, concatenatedDataID.String()),
		Logs:    DeviceFileTypeUrls(ac, "logs", deviceID, concatenatedLogsID.String()),
	}
}

func DeviceSummaryType(ac *ApiConfiguration, s *data.DeviceSummary, entries map[string][]*app.LocationEntry, locations []*data.DeviceStreamLocation) *app.DeviceSummary {
	lh := &app.LocationHistory{}

	if deviceEntries, ok := entries[s.DeviceID]; ok {
		lh.Entries = deviceEntries
	}

	return &app.DeviceSummary{
		DeviceID:      s.DeviceID,
		LastFileID:    s.LastFileID,
		LastFileTime:  s.LastFileTime,
		NumberOfFiles: s.NumberOfFiles,
		Urls:          DeviceSummaryUrls(ac, s.DeviceID),
		Locations:     lh,
	}
}

func DevicesType(ac *ApiConfiguration, devices []*data.DeviceSummary, locations []*data.DeviceStreamLocation) *app.Devices {
	entries := make(map[string][]*app.LocationEntry)
	for _, dsl := range locations {
		if _, ok := entries[dsl.DeviceID]; !ok {
			entries[dsl.DeviceID] = make([]*app.LocationEntry, 0)
		}

		le := &app.LocationEntry{
			Coordinates: dsl.Location.Coordinates(),
			Time:        dsl.Timestamp,
		}

		entries[dsl.DeviceID] = append(entries[dsl.DeviceID], le)
	}

	summaries := make([]*app.DeviceSummary, len(devices))
	for i, summary := range devices {
		summaries[i] = DeviceSummaryType(ac, summary, entries, locations)
	}
	return &app.Devices{
		Devices: summaries,
	}
}

func DeviceFileSummaryType(ac *ApiConfiguration, s *data.DeviceFile) *app.DeviceFile {
	return &app.DeviceFile{
		DeviceID:     s.DeviceID,
		FileID:       s.StreamID,
		Firmware:     s.Firmware,
		ID:           int(s.ID),
		Meta:         s.Meta.String(),
		Size:         int(s.Size),
		FileTypeID:   s.FileID,
		FileTypeName: FileTypeNames[s.FileID],
		Time:         s.Time,
		URL:          s.URL,
		Urls: &app.DeviceFileUrls{
			Csv:  ac.MakeApiUrl("/files/%v/data.csv", s.ID),
			JSON: ac.MakeApiUrl("/files/%v/data.json", s.ID),
			Fkpb: ac.MakeApiUrl("/files/%v/data.fkpb", s.ID),
		},
	}
}

func DeviceFilesType(ac *ApiConfiguration, files []*data.DeviceFile) *app.DeviceFiles {
	summaries := make([]*app.DeviceFile, len(files))
	for i, summary := range files {
		summaries[i] = DeviceFileSummaryType(ac, summary)
	}
	return &app.DeviceFiles{
		Files: summaries,
	}
}

func ConcatenatedFileInfo(urls *app.DeviceFileTypeUrls, fi *backend.FileInfo) *app.ConcatenatedFileInfo {
	if fi == nil {
		return nil
	}
	return &app.ConcatenatedFileInfo{
		Time: fi.LastModified,
		Size: int(fi.Size),
		Csv:  urls.Csv,
	}
}

type ConcatenationJob struct {
	DeviceID    string
	FileTypeIDs []string
}

type ConcatenationWorkers struct {
	options FilesControllerOptions
	channel chan ConcatenationJob
}

func NewConcatenationWorkers(ctx context.Context, options FilesControllerOptions) (cw *ConcatenationWorkers, err error) {
	cw = &ConcatenationWorkers{
		channel: make(chan ConcatenationJob, 100),
		options: options,
	}

	for w := 0; w < 1; w++ {
		go cw.worker(ctx)
	}

	return
}

func (cw *ConcatenationWorkers) worker(ctx context.Context) {
	log := Logger(ctx).Sugar()

	log.Infow("Worker starting")

	for job := range cw.channel {
		log.Infow("Worker", "job", job)

		fileTypeID := job.FileTypeIDs[0]
		space := ConcatenatedFileSpaces[fileTypeID]
		deviceStreamID := uuid.NewSHA1(space, []byte(job.DeviceID))

		fc := &backend.FileConcatenator{
			Session:    cw.options.Session,
			Database:   cw.options.Database,
			FileID:     deviceStreamID.String(),
			FileTypeID: fileTypeID,
			DeviceID:   job.DeviceID,
			TypeIDs:    job.FileTypeIDs,
		}

		fc.Concatenate(ctx)
	}

	log.Infow("Worker exiting")
}
