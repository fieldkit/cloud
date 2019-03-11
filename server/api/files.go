package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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
	Session  *session.Session
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type FilesController struct {
	*goa.Controller
	options FilesControllerOptions
}

var (
	DataFileTypeIDs = []string{"4"}
	LogFileTypeIDs  = []string{"2", "3"}
	FileTypeNames   = map[string]string{
		"2": "Logs",
		"3": "Logs",
		"4": "Data",
	}
)

func NewFilesController(service *goa.Service, options FilesControllerOptions) *FilesController {
	return &FilesController{
		Controller: service.NewController("FilesController"),
		options:    options,
	}
}

func DeviceFileSummaryType(s *data.DeviceFile) *app.DeviceFile {
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
			Csv:  fmt.Sprintf("/files/%v/data.csv", s.ID),
			JSON: fmt.Sprintf("/files/%v/data.json", s.ID),
			Raw:  fmt.Sprintf("/files/%v/data.fkpb", s.ID),
		},
	}
}

func DeviceFilesType(files []*data.DeviceFile) *app.DeviceFiles {
	summaries := make([]*app.DeviceFile, len(files))
	for i, summary := range files {
		summaries[i] = DeviceFileSummaryType(summary)
	}
	return &app.DeviceFiles{
		Files: summaries,
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

	return DeviceFilesType(files), nil
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
	devices := []*DeviceSummary{}
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

	return ctx.OK(DevicesType(devices))
}

func (c *FilesController) File(ctx *app.FileFilesContext) error {
	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.stream_id = $1)`, ctx.FileID); err != nil {
		return err
	}

	if len(files) != 1 {
		return ctx.NotFound()
	}

	return ctx.OK(DeviceFileSummaryType(files[0]))
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

	cs, err := iterator.Next(ctx)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.csv", cs.File.StreamID)

	ctx.ResponseData.Header().Set("Content-Type", backend.FkDataBinaryContentType)
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))

	ctx.ResponseData.WriteHeader(http.StatusOK)

	_, err = io.Copy(ctx.ResponseData, cs.Response.Body)
	if err != nil {
		return err
	}

	return nil
}

type DeviceLogsController struct {
	*goa.Controller
	options FilesControllerOptions
}

func NewDeviceLogsController(service *goa.Service, options FilesControllerOptions) *DeviceLogsController {
	return &DeviceLogsController{
		Controller: service.NewController("DeviceLogsController"),
		options:    options,
	}
}

func (c *DeviceLogsController) All(ctx *app.AllDeviceLogsContext) error {
	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files, `SELECT s.* FROM fieldkit.device_stream AS s WHERE s.device_id = $1 AND s.children IS NOT NULL`, ctx.DeviceID); err != nil {
		return err
	}

	if len(files) == 0 {
		newFileID, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		fc := &backend.FileConcatenator{
			Session:    c.options.Session,
			Database:   c.options.Database,
			FileID:     newFileID.String(),
			FileTypeID: DataFileTypeIDs[0],
			DeviceID:   ctx.DeviceID,
			TypeIDs:    DataFileTypeIDs,
		}

		go fc.Concatenate(ctx)

		if true {
			ctx.ResponseData.Header().Set("Location", fmt.Sprintf("https://api.fieldkit.org/files/%s", newFileID))
			return ctx.Busy()
		}

		return ctx.Busy()
	}

	return ctx.OK([]byte{})
}

type DeviceDataController struct {
	*goa.Controller
	options FilesControllerOptions
}

func NewDeviceDataController(service *goa.Service, options FilesControllerOptions) *DeviceDataController {
	return &DeviceDataController{
		Controller: service.NewController("DeviceDataController"),
		options:    options,
	}
}

func (c *DeviceDataController) All(ctx *app.AllDeviceDataContext) error {
	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files, `SELECT s.* FROM fieldkit.device_stream AS s WHERE s.device_id = $1 AND s.children IS NOT NULL`, ctx.DeviceID); err != nil {
		return err
	}

	if len(files) == 0 {
		newFileID, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		fc := &backend.FileConcatenator{
			Session:    c.options.Session,
			Database:   c.options.Database,
			FileID:     newFileID.String(),
			FileTypeID: DataFileTypeIDs[0],
			DeviceID:   ctx.DeviceID,
			TypeIDs:    DataFileTypeIDs,
		}

		go fc.Concatenate(ctx)

		if true {
			ctx.ResponseData.Header().Set("Location", fmt.Sprintf("https://api.fieldkit.org/files/%s", newFileID))
			return ctx.Busy()
		}

		return ctx.Found()
	}

	return ctx.OK([]byte{})
}

type DeviceSummary struct {
	DeviceID      string    `db:"device_id"`
	LastFileID    string    `db:"last_stream_id"`
	LastFileTime  time.Time `db:"last_stream_time"`
	NumberOfFiles int       `db:"number_of_files"`
}

func DeviceSummaryType(s *DeviceSummary) *app.Device {
	return &app.Device{
		DeviceID:      s.DeviceID,
		LastFileID:    s.LastFileID,
		LastFileTime:  s.LastFileTime,
		NumberOfFiles: s.NumberOfFiles,
	}
}

func DevicesType(devices []*DeviceSummary) *app.Devices {
	summaries := make([]*app.Device, len(devices))
	for i, summary := range devices {
		summaries[i] = DeviceSummaryType(summary)
	}
	return &app.Devices{
		Devices: summaries,
	}
}
