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
}

type FilesController struct {
	BaseFilesController
	*goa.Controller
}

var (
	ConcatenatedFilesSpace = uuid.MustParse("5554c632-d586-53f2-b2a7-88a63663a3f5")
	DataFileTypeIDs        = []string{"4"}
	LogFileTypeIDs         = []string{"2", "3"}
	FileTypeNames          = map[string]string{
		"2": "Logs",
		"3": "Logs",
		"4": "Data",
	}
)

func NewFilesController(service *goa.Service, options FilesControllerOptions) *FilesController {
	return &FilesController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("FilesController"),
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
			Csv:  fmt.Sprintf("/files/%v/data.csv", s.ID),
			JSON: fmt.Sprintf("/files/%v/data.json", s.ID),
			Raw:  fmt.Sprintf("/files/%v/data.fkpb", s.ID),
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

	return ctx.OK(DevicesType(devices))
}

func (c *FilesController) File(ctx *app.FileFilesContext) error {
	log := Logger(ctx).Sugar()

	files := []*data.DeviceFile{}
	if err := c.options.Database.SelectContext(ctx, &files,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.stream_id = $1)`, ctx.FileID); err != nil {
		return err
	}

	if len(files) != 1 {
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
			FileID: fi.Key,
			Time:   fi.LastModified,
			Size:   int(fi.Size),
			URL:    fi.URL,
			Meta:   string(jsonMeta),
			Urls: &app.DeviceFileUrls{
				Csv:  fmt.Sprintf("/files/%v/data.csv", fi.Key),
				JSON: fmt.Sprintf("/files/%v/data.json", fi.Key),
				Raw:  fmt.Sprintf("/files/%v/data.fkpb", fi.Key),
			},
		})
	}

	return ctx.OK(DeviceFileSummaryType(c.options.Config, files[0]))
}

func (c *FilesController) Csv(ctx *app.CsvFilesContext) error {
	iterator, err := backend.LookupFileOnS3(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	exporter := backend.NewSimpleCsvExporter(ctx.ResponseData)

	return backend.ExportAllFiles(ctx, ctx.ResponseData, ctx.Dl, iterator, exporter)
}

func (c *FilesController) JSON(ctx *app.JSONFilesContext) error {
	iterator, err := backend.LookupFileOnS3(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	exporter := backend.NewSimpleJsonExporter(ctx.ResponseData)

	return backend.ExportAllFiles(ctx, ctx.ResponseData, ctx.Dl, iterator, exporter)
}

func (c *FilesController) Raw(ctx *app.RawFilesContext) error {
	iterator, err := backend.LookupFileOnS3(ctx, c.options.Session, c.options.Database, ctx.FileID)
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
	deviceStreamID := uuid.NewSHA1(ConcatenatedFilesSpace, []byte(deviceID))
	fileTypeID := fileTypeIDs[0]

	fc := &backend.FileConcatenator{
		Session:    c.options.Session,
		Database:   c.options.Database,
		FileID:     deviceStreamID.String(),
		FileTypeID: fileTypeID,
		DeviceID:   deviceID,
		TypeIDs:    fileTypeIDs,
	}

	go fc.Concatenate(ctx)

	return c.options.Config.MakeApiUrl("/files/%s", deviceStreamID), nil
}

type DeviceLogsController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceLogsController(service *goa.Service, options FilesControllerOptions) *DeviceLogsController {
	return &DeviceLogsController{
		BaseFilesController: BaseFilesController{
			options: options,
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

func NewDeviceDataController(service *goa.Service, options FilesControllerOptions) *DeviceDataController {
	return &DeviceDataController{
		BaseFilesController: BaseFilesController{
			options: options,
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

func DeviceSummaryType(s *data.DeviceSummary) *app.Device {
	return &app.Device{
		DeviceID:      s.DeviceID,
		LastFileID:    s.LastFileID,
		LastFileTime:  s.LastFileTime,
		NumberOfFiles: s.NumberOfFiles,
		Urls: &app.DeviceSummaryUrls{
			Data: fmt.Sprintf("/devices/%v/data", s.DeviceID),
			Logs: fmt.Sprintf("/devices/%v/logs", s.DeviceID),
		},
	}
}

func DevicesType(devices []*data.DeviceSummary) *app.Devices {
	summaries := make([]*app.Device, len(devices))
	for i, summary := range devices {
		summaries[i] = DeviceSummaryType(summary)
	}
	return &app.Devices{
		Devices: summaries,
	}
}
