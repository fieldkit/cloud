package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type BaseFilesController struct {
	options *ControllerOptions
}

type FilesController struct {
	BaseFilesController
	*goa.Controller
}

func NewFilesController(ctx context.Context, service *goa.Service, options *ControllerOptions) *FilesController {
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
	files, err := c.listDeviceFiles(ctx, backend.DataFileTypeIDs, ctx.DeviceID, ctx.Page)
	if err != nil {
		return err
	}

	return ctx.OK(files)
}

func (c *FilesController) ListDeviceLogFiles(ctx *app.ListDeviceLogFilesFilesContext) error {
	files, err := c.listDeviceFiles(ctx, backend.LogFileTypeIDs, ctx.DeviceID, ctx.Page)
	if err != nil {
		return err
	}

	return ctx.OK(files)
}

func (c *FilesController) ListDevices(ctx *app.ListDevicesFilesContext) error {
	locations := []*data.DeviceStreamLocationAndPlace{}
	if err := c.options.Database.SelectContext(ctx, &locations, `
	    SELECT
	      s.id, s.device_id, s.timestamp, ST_AsBinary(s.location) AS location,
	      ARRAY(SELECT name FROM fieldkit.countries c WHERE ST_Contains(c.geom, location)) AS places
	    FROM
	    (
	      SELECT
		DISTINCT ON (dsl.location)
		dsl.*,
		ROW_NUMBER() OVER (
		  PARTITION BY device_id
		  ORDER BY dsl.timestamp DESC
		) AS c
	      FROM fieldkit.device_stream_location AS dsl
	    ) AS s
	    WHERE s.c < 100`); err != nil {
		return err
	}

	devices := []*data.DeviceSummary{}
	if err := c.options.Database.SelectContext(ctx, &devices,
		`SELECT s.device_id,
		    (SELECT stream_id FROM fieldkit.device_stream AS s2 WHERE (s2.device_id = s.device_id) ORDER BY s2.time DESC LIMIT 1) AS last_stream_id,
		    (SELECT n.name FROM fieldkit.device_notes AS n WHERE (n.device_id = s.device_id) AND (n.name IS NOT NULL) ORDER BY n.time DESC LIMIT 1) AS name,
		    MAX(s.time) AS last_stream_time,
		    COUNT(s.*) AS number_of_files,
		    COALESCE(SUM(s.size) FILTER (WHERE s.file_id != '4'), 0) AS logs_size,
		    COALESCE(SUM(s.size) FILTER (WHERE s.file_id  = '4'), 0) AS data_size
		 FROM fieldkit.device_stream AS s
                 WHERE s.device_id != '' AND s.file_id != '5' AND s.file_id != '6' AND
                     (s.device_id IN (SELECT key FROM fieldkit.device AS d JOIN fieldkit.source AS s ON (d.source_id = s.id) WHERE s.visible))
                 GROUP BY s.device_id
                 ORDER BY last_stream_time DESC`); err != nil {
		return err
	}

	return ctx.OK(DevicesType(c.options.Config, devices, locations))
}

func DeviceNotesType(ctx context.Context, notes []*data.DeviceNotes) []*app.DeviceNotesEntry {
	entries := make([]*app.DeviceNotesEntry, len(notes))
	for i, note := range notes {
		entries[i] = &app.DeviceNotesEntry{
			Time:  note.Time,
			Name:  note.Name,
			Notes: note.Notes,
		}
	}
	return entries
}

func (c *FilesController) getDeviceDetails(ctx context.Context, deviceID string) (dd *app.DeviceDetails, err error) {
	ac := c.options.Config

	var notes []*data.DeviceNotes
	if err := c.options.Database.SelectContext(ctx, &notes,
		`SELECT n.time, n.name, n.notes
		 FROM fieldkit.device_notes AS n
                 WHERE n.device_id = $1
                 ORDER BY time DESC`, deviceID); err != nil {
		return nil, err
	}

	fr, err := backend.NewFileRepository(c.options.Session, "fk-streams")
	if err != nil {
		return nil, err
	}

	urls := DeviceSummaryUrls(ac, deviceID)

	data, err := fr.Info(ctx, urls.Data.ID)
	if err != nil {
		return nil, err
	}

	logs, err := fr.Info(ctx, urls.Logs.ID)
	if err != nil {
		return nil, err
	}

	notesType := DeviceNotesType(ctx, notes)

	return &app.DeviceDetails{
		DeviceID: deviceID,
		Urls:     urls,
		Notes:    notesType,
		Files: &app.ConcatenatedFilesInfo{
			Data: ConcatenatedFileInfo(urls.Data, data),
			Logs: ConcatenatedFileInfo(urls.Logs, logs),
		},
	}, nil
}

func (c *FilesController) DeviceInfo(ctx *app.DeviceInfoFilesContext) error {
	dd, err := c.getDeviceDetails(ctx, ctx.DeviceID)
	if err != nil {
		return err
	}
	return ctx.OK(dd)
}

func (c *FilesController) GetDeviceLocationHistory(ctx *app.GetDeviceLocationHistoryFilesContext) error {
	pageSize := 10000
	offset := 0
	if ctx.Page != nil {
		offset = *ctx.Page * pageSize
	}
	locations := []*data.DeviceStreamLocationAndPlace{}
	if err := c.options.Database.SelectContext(ctx, &locations, `
	    SELECT dsl.id, dsl.device_id, dsl.timestamp, ST_AsBinary(dsl.location) AS location FROM fieldkit.device_stream_location AS dsl WHERE dsl.device_id = $1 ORDER BY dsl.timestamp DESC LIMIT $2 OFFSET $3
	  `, ctx.DeviceID, pageSize, offset); err != nil {
		return err
	}

	entries := LocationEntries(locations)

	lh := &app.LocationHistory{
		Entries: entries[ctx.DeviceID],
	}

	return ctx.OK(lh)
}

func (c *FilesController) UpdateDeviceInfo(ctx *app.UpdateDeviceInfoFilesContext) error {
	newNote := data.DeviceNotes{
		DeviceID: ctx.Payload.DeviceID,
		Time:     time.Now(),
		Name:     ctx.Payload.Name,
		Notes:    ctx.Payload.Notes,
	}

	err := c.options.Database.NamedGetContext(ctx, &newNote, `
               INSERT INTO fieldkit.device_notes (device_id, time, name, notes)
	       VALUES (:device_id, :time, :name, :notes) RETURNING *`, newNote)
	if err != nil {
		return err
	}

	dd, err := c.getDeviceDetails(ctx, ctx.DeviceID)
	if err != nil {
		return err
	}
	return ctx.OK(dd)
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

		log.Infow("File", "file_key", fi.Key, "file_url", fi.URL, "size", fi.Size)

		return ctx.OK(&app.DeviceFile{
			DeviceID:     fi.DeviceID,
			FileID:       fi.Key,
			FileTypeID:   fi.FileTypeID,
			FileTypeName: backend.FileTypeNames[fi.FileTypeID],
			Time:         fi.LastModified,
			Size:         int(fi.Size),
			URL:          fi.URL,
			Meta:         string(jsonMeta),
			Urls: &app.DeviceFileUrls{
				Csv:  ac.MakeApiUrl("/files/%v/data.csv", fi.Key),
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

func (c *FilesController) Raw(ctx *app.RawFilesContext) error {
	iterator, err := backend.LookupFile(ctx, c.options.Session, c.options.Database, ctx.FileID)
	if err != nil {
		return err
	}

	opened, err := iterator.Next(ctx)
	if err != nil {
		return err
	}

	fileType := backend.FileTypeNames[opened.FileTypeID]

	fileName := fmt.Sprintf("%s-%s-%s.fkpb", opened.StreamID, opened.FileTypeID, fileType)

	ctx.ResponseData.Header().Set("Content-Type", backend.FkDataBinaryContentType)
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	ctx.ResponseData.WriteHeader(http.StatusOK)

	_, err = io.Copy(ctx.ResponseData, opened.Response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *FilesController) Status(ctx *app.StatusFilesContext) error {
	return ctx.OK(&app.FilesStatus{
		Queued: c.options.ConcatWorkers.Length(),
	})
}

func (c *BaseFilesController) concatenatedDeviceFile(ctx context.Context, responseData *goa.ResponseData, deviceID string, fileTypeIDs []string) (redirection string, err error) {
	space := backend.ConcatenatedFileSpaces[fileTypeIDs[0]]
	deviceStreamID := uuid.NewSHA1(space, []byte(deviceID))

	if err := c.options.ConcatWorkers.QueueJob(ctx, deviceID, fileTypeIDs); err != nil {
		return "", err
	}

	return c.options.Config.MakeApiUrl("/files/%s", deviceStreamID), nil
}

type DeviceLogsController struct {
	BaseFilesController
	*goa.Controller
}

func NewDeviceLogsController(ctx context.Context, service *goa.Service, options *ControllerOptions) *DeviceLogsController {
	return &DeviceLogsController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("DeviceLogsController"),
	}
}

func (c *DeviceLogsController) All(ctx *app.AllDeviceLogsContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, backend.LogFileTypeIDs)
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

func NewDeviceDataController(ctx context.Context, service *goa.Service, options *ControllerOptions) *DeviceDataController {
	return &DeviceDataController{
		BaseFilesController: BaseFilesController{
			options: options,
		},
		Controller: service.NewController("DeviceDataController"),
	}
}

func (c *DeviceDataController) All(ctx *app.AllDeviceDataContext) error {
	url, err := c.concatenatedDeviceFile(ctx, ctx.ResponseData, ctx.DeviceID, backend.DataFileTypeIDs)
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
	concatenatedDataID := uuid.NewSHA1(backend.ConcatenatedDataSpace, []byte(deviceID))
	concatenatedLogsID := uuid.NewSHA1(backend.ConcatenatedLogsSpace, []byte(deviceID))

	return &app.DeviceSummaryUrls{
		Details: ac.MakeApiUrl("/devices/%v", deviceID),
		Data:    DeviceFileTypeUrls(ac, "data", deviceID, concatenatedDataID.String()),
		Logs:    DeviceFileTypeUrls(ac, "logs", deviceID, concatenatedLogsID.String()),
	}
}

func DeviceSummaryType(ac *ApiConfiguration, s *data.DeviceSummary, entries map[string][]*app.LocationEntry, locations []*data.DeviceStreamLocationAndPlace) *app.DeviceSummary {
	lh := &app.LocationHistory{}

	if deviceEntries, ok := entries[s.DeviceID]; ok {
		lh.Entries = deviceEntries
	} else {
		lh.Entries = make([]*app.LocationEntry, 0)
	}

	return &app.DeviceSummary{
		DeviceID:      s.DeviceID,
		LastFileID:    s.LastFileID,
		LastFileTime:  s.LastFileTime,
		NumberOfFiles: s.NumberOfFiles,
		LogsSize:      s.LogsSize,
		DataSize:      s.DataSize,
		Name:          s.Name,
		Urls:          DeviceSummaryUrls(ac, s.DeviceID),
		Locations:     lh,
	}
}

func getPlaces(dsl *data.DeviceStreamLocationAndPlace) *string {
	if len(dsl.Places) > 0 {
		places := strings.Join(dsl.Places, ", ")
		return &places
	}
	return nil
}

func LocationEntries(locations []*data.DeviceStreamLocationAndPlace) map[string][]*app.LocationEntry {
	entries := make(map[string][]*app.LocationEntry)
	for _, dsl := range locations {
		if _, ok := entries[dsl.DeviceID]; !ok {
			entries[dsl.DeviceID] = make([]*app.LocationEntry, 0)
		}

		le := &app.LocationEntry{
			Coordinates: dsl.Location.Coordinates(),
			Time:        dsl.Timestamp,
			Places:      getPlaces(dsl),
		}

		entries[dsl.DeviceID] = append(entries[dsl.DeviceID], le)
	}
	return entries
}

func DevicesType(ac *ApiConfiguration, devices []*data.DeviceSummary, locations []*data.DeviceStreamLocationAndPlace) *app.Devices {
	entries := LocationEntries(locations)
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
		FileTypeName: backend.FileTypeNames[s.FileID],
		Time:         s.Time,
		URL:          s.URL,
		Corrupted:    len(s.Flags) > 0,
		Urls: &app.DeviceFileUrls{
			Csv:  ac.MakeApiUrl("/files/%v/data.csv", s.ID),
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
