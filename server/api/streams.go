package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/data"
)

type StreamsControllerOptions struct {
	Session  *session.Session
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type StreamsController struct {
	*goa.Controller
	options StreamsControllerOptions
}

func NewStreamsController(service *goa.Service, options StreamsControllerOptions) *StreamsController {
	return &StreamsController{
		Controller: service.NewController("StreamsController"),
		options:    options,
	}
}

func DeviceStreamSummaryType(s *data.DeviceStream) *app.DeviceStream {
	return &app.DeviceStream{
		DeviceID: s.DeviceID,
		FileID:   s.FileID,
		Firmware: s.Firmware,
		ID:       int(s.ID),
		Meta:     s.Meta.String(),
		Size:     int(s.Size),
		StreamID: s.StreamID,
		Time:     s.Time,
		URL:      s.URL,
	}
}

func DeviceStreamsType(streams []*data.DeviceStream) *app.DeviceStreams {
	summaries := make([]*app.DeviceStream, len(streams))
	for i, summary := range streams {
		summaries[i] = DeviceStreamSummaryType(summary)
	}
	return &app.DeviceStreams{
		Streams: summaries,
	}
}

func (c *StreamsController) ListAll(ctx *app.ListAllStreamsContext) error {
	pageSize := 100
	offset := 0
	if ctx.Page != nil {
		offset = *ctx.Page * pageSize
	}

	streams := []*data.DeviceStream{}
	if err := c.options.Database.SelectContext(ctx, &streams,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE ($1::text IS NULL OR s.file_id = $1) ORDER BY time DESC LIMIT $2 OFFSET $3`, ctx.FileID, pageSize, offset); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	for _, s := range streams {
		signed, err := SignS3URL(svc, s.URL)
		if err != nil {
			return err
		}
		s.URL = signed
	}

	return ctx.OK(DeviceStreamsType(streams))
}

func (c *StreamsController) ListDevice(ctx *app.ListDeviceStreamsContext) error {
	log := Logger(ctx).Sugar()

	log.Infow("Device", "device_id", ctx.DeviceID)

	pageSize := 100
	offset := 0
	if ctx.Page != nil {
		offset = *ctx.Page * pageSize
	}

	streams := []*data.DeviceStream{}
	if err := c.options.Database.SelectContext(ctx, &streams,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE ($1::text IS NULL OR s.file_id = $1) AND (s.device_id = $2) ORDER BY time DESC LIMIT $3 OFFSET $4`, ctx.FileID, ctx.DeviceID, pageSize, offset); err != nil {
		return err
	}

	return ctx.OK(DeviceStreamsType(streams))
}

type Device struct {
	DeviceID        string    `db:"device_id"`
	LastStreamID    string    `db:"last_stream_id"`
	LastStreamTime  time.Time `db:"last_stream_time"`
	NumberOfStreams int       `db:"number_of_streams"`
}

func DeviceSummaryType(s *Device) *app.Device {
	return &app.Device{
		DeviceID:        s.DeviceID,
		LastStreamID:    s.LastStreamID,
		LastStreamTime:  s.LastStreamTime,
		NumberOfStreams: s.NumberOfStreams,
	}
}

func DevicesType(devices []*Device) *app.Devices {
	summaries := make([]*app.Device, len(devices))
	for i, summary := range devices {
		summaries[i] = DeviceSummaryType(summary)
	}
	return &app.Devices{
		Devices: summaries,
	}
}

func (c *StreamsController) ListDevices(ctx *app.ListDevicesStreamsContext) error {
	devices := []*Device{}
	if err := c.options.Database.SelectContext(ctx, &devices,
		`SELECT s.device_id,
		    (SELECT stream_id FROM fieldkit.device_stream AS s2 ORDER BY s2.time DESC LIMIT 1) AS last_stream_id,
		    MAX(s.time) AS last_stream_time,
		    COUNT(s.*) AS number_of_streams
		 FROM fieldkit.device_stream AS s
                 GROUP BY s.device_id ORDER BY last_stream_time DESC`); err != nil {
		return err
	}

	return ctx.OK(DevicesType(devices))
}

func (c *StreamsController) Csv(ctx *app.CsvStreamsContext) error {
	streamer, err := c.LookupStream(ctx, ctx.StreamID)
	if err != nil {
		return err
	}

	exporter := NewSimpleCsvExporter(ctx.ResponseData)

	return ExportAllStreams(ctx, streamer, exporter)
}

func (c *StreamsController) JSON(ctx *app.JSONStreamsContext) error {
	streamer, err := c.LookupStream(ctx, ctx.StreamID)
	if err != nil {
		return err
	}

	exporter := NewSimpleJsonExporter(ctx.ResponseData)

	return ExportAllStreams(ctx, streamer, exporter)
}

func (c *StreamsController) DeviceData(ctx *app.DeviceDataStreamsContext) error {
	streamer, err := c.LookupDeviceStreams(ctx, ctx.DeviceID, "4")
	if err != nil {
		return err
	}

	exporter := NewSimpleCsvExporter(ctx.ResponseData)

	return ExportAllStreams(ctx, streamer, exporter)
}

func (c *StreamsController) DeviceLogs(ctx *app.DeviceLogsStreamsContext) error {
	streamer, err := c.LookupDeviceStreams(ctx, ctx.DeviceID, "2")
	if err != nil {
		return err
	}

	exporter := NewSimpleCsvExporter(ctx.ResponseData)

	return ExportAllStreams(ctx, streamer, exporter)
}

func ExportAllStreams(ctx context.Context, streamer *Streamer, exporter Exporter) error {
	log := Logger(ctx).Sugar()

	for {
		cs, err := streamer.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if cs == nil {
			break
		}

		defer cs.Response.Body.Close()

		binaryReader := backend.NewFkBinaryReader(exporter.ForStream(cs.Stream))

		err = binaryReader.Read(ctx, cs.Response.Body)
		if err != nil {
			log.Infow("Error reading stream", "error", err, "stream_id", cs.Stream.StreamID)
		}
	}

	return nil
}

type Exporter interface {
	ForStream(stream *data.DeviceStream) backend.FormattedMessageReceiver
}

type SimpleJsonExporter struct {
	Stream  *data.DeviceStream
	Writer  io.Writer
	Records int
}

func NewSimpleJsonExporter(writer io.Writer) Exporter {
	return &SimpleJsonExporter{Writer: writer}
}

func (ce *SimpleJsonExporter) ForStream(stream *data.DeviceStream) backend.FormattedMessageReceiver {
	ce.Stream = stream
	return ce
}

func (ce *SimpleJsonExporter) HandleRecord(ctx context.Context, r *pb.DataRecord) error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if ce.Records > 0 {
		fmt.Fprintf(ce.Writer, ",\n")
	} else {
		fmt.Fprintf(ce.Writer, "[\n")
	}

	_, err = ce.Writer.Write(body)
	if err != nil {
		return err
	}

	ce.Records += 1

	return nil
}

func (ce *SimpleJsonExporter) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	return nil, nil
}

func (ce *SimpleJsonExporter) Finish(ctx context.Context) error {
	return nil
}

type SimpleCsvExporter struct {
	Stream        *data.DeviceStream
	Writer        io.Writer
	HeaderWritten bool
}

func NewSimpleCsvExporter(writer io.Writer) Exporter {
	return &SimpleCsvExporter{Writer: writer}
}

func (ce *SimpleCsvExporter) ForStream(stream *data.DeviceStream) backend.FormattedMessageReceiver {
	ce.Stream = stream
	return ce
}

func (ce *SimpleCsvExporter) HandleRecord(ctx context.Context, r *pb.DataRecord) error {
	if r.Log != nil {
		return ce.ExportLog(ctx, r)
	}

	return nil
}

func (ce *SimpleCsvExporter) ExportLog(ctx context.Context, r *pb.DataRecord) error {
	if !ce.HeaderWritten {
		fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v\n", "DeviceID", "StreamID", "StreamID", "Uptime", "Time", "Level", "Facility", "Message")
		ce.HeaderWritten = true
	}

	fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v\n", ce.Stream.DeviceID, ce.Stream.StreamID, ce.Stream.ID, r.Log.Uptime, r.Log.Time, r.Log.Level, r.Log.Facility, strings.TrimSpace(r.Log.Message))

	return nil
}

func (ce *SimpleCsvExporter) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	opaqueKeys := reflect.ValueOf(fm.MapValues).MapKeys()
	keys := make([]string, len(opaqueKeys))
	for i := 0; i < len(opaqueKeys); i++ {
		keys[i] = opaqueKeys[i].String()
	}
	sort.Strings(keys)

	if !ce.HeaderWritten {
		fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v,", "Device", "Stream", "Stream", "Message", "Time", "Longitude", "Latitude", "Fixed")

		for _, key := range keys {
			fmt.Fprintf(ce.Writer, ",%v", key)
		}

		fmt.Fprintf(ce.Writer, "\n")

		ce.HeaderWritten = true
	}

	fmt.Fprintf(ce.Writer, "%v,%v,%v,", ce.Stream.DeviceID, ce.Stream.StreamID, ce.Stream.ID)

	fmt.Fprintf(ce.Writer, "%v,%v,", fm.MessageId, fm.Time)

	if fm.Location != nil && len(fm.Location) >= 2 {
		fmt.Fprintf(ce.Writer, "%v,%v", fm.Location[0], fm.Location[1])
	} else {
		fmt.Fprintf(ce.Writer, "%v,%v", 0.0, 0.0)
	}

	fmt.Fprintf(ce.Writer, ",%v", fm.Fixed)

	for _, key := range keys {
		fmt.Fprintf(ce.Writer, ",%v", fm.MapValues[key])
	}

	fmt.Fprintf(ce.Writer, "\n")

	return nil, nil
}

func (ce *SimpleCsvExporter) Finish(ctx context.Context) error {
	return nil
}

type Streamer struct {
	S3Service *s3.S3
	Offset    int
	Limit     int
	Query     func(s *Streamer) error
	Streams   []*data.DeviceStream
	Index     int
}

type CurrentStream struct {
	Stream    *data.DeviceStream
	SignedURL string
	Response  *http.Response
}

func (c *StreamsController) LookupStream(ctx context.Context, streamID string) (streamer *Streamer, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("Stream", "stream_id", streamID)

	streamer = &Streamer{
		Offset:    0,
		Limit:     0,
		Index:     0,
		S3Service: s3.New(c.options.Session),
		Query: func(s *Streamer) error {
			s.Streams = []*data.DeviceStream{}
			if err := c.options.Database.SelectContext(ctx, &s.Streams, `SELECT s.* FROM fieldkit.device_stream AS s WHERE s.id = $1`, streamID); err != nil {
				return err
			}
			return nil
		},
	}

	return
}

func (c *StreamsController) LookupDeviceStreams(ctx context.Context, deviceID string, fileID string) (streamer *Streamer, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("Stream", "device_id", deviceID, "file_id", fileID)

	streamer = &Streamer{
		Offset:    0,
		Limit:     10,
		Index:     0,
		S3Service: s3.New(c.options.Session),
		Query: func(s *Streamer) error {
			s.Streams = []*data.DeviceStream{}
			if err := c.options.Database.SelectContext(ctx, &s.Streams,
				`SELECT s.* FROM fieldkit.device_stream AS s WHERE ($1::text IS NULL OR s.file_id = $1) AND (s.device_id = $2) ORDER BY time DESC LIMIT $3 OFFSET $4`,
				fileID, deviceID, s.Limit, s.Offset); err != nil {
				return err
			}
			return nil
		},
	}

	return
}

func (streamer *Streamer) Next(ctx context.Context) (cs *CurrentStream, err error) {
	log := Logger(ctx).Sugar()

	if streamer.Streams != nil && streamer.Index >= len(streamer.Streams) {
		streamer.Streams = nil

		// Easy way to handle the single stream per batch.
		if streamer.Limit == 0 {
			log.Infow("No more batches")
			return nil, io.EOF
		}
	}

	if streamer.Streams == nil {
		err = streamer.Query(streamer)
		if err != nil {
			return nil, err
		}

		log.Infow("Queried", "batch_size", len(streamer.Streams))

		if len(streamer.Streams) == 0 {
			log.Infow("No more batches")
			return nil, io.EOF
		}

		streamer.Index = 0
	}

	stream := streamer.Streams[streamer.Index]

	log.Infow("Stream", "stream_id", stream.StreamID, "index", streamer.Index, "size", stream.Size, "url", stream.URL)

	streamer.Index += 1

	signed, err := SignS3URL(streamer.S3Service, stream.URL)
	if err != nil {
		return nil, fmt.Errorf("Error signing stream URL: %v (%v)", stream.URL, err)
	}

	response, err := http.Get(signed)
	if err != nil {
		return nil, fmt.Errorf("Error opening stream URL: %v (%v)", stream.URL, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error retrieving stream: %v (status = %v)", stream.URL, response.StatusCode)
	}

	cs = &CurrentStream{
		Stream:    stream,
		SignedURL: signed,
		Response:  response,
	}

	return
}
