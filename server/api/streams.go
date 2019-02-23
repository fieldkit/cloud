package api

import (
	"context"
	"io"
	"net/http"

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
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE ($1::text IS NULL OR s.file_id = $1) AND (s.device_id = $2) ORDER BY time DESC LIMIT $3 OFFSET $4`, ctx.FileID, ctx.DeviceID, ctx.FileID, pageSize, offset); err != nil {
		return err
	}

	return ctx.OK(DeviceStreamsType(streams))
}

type Streamer struct {
	Stream    *data.DeviceStream
	SignedURL string
}

func (c *StreamsController) lookupDeviceStream(ctx context.Context, streamID string) (*Streamer, error) {
	log := Logger(ctx).Sugar()

	log.Infow("Stream", "stream_id", streamID)

	streams := []*data.DeviceStream{}
	if err := c.options.Database.SelectContext(ctx, &streams,
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE s.id = $1`, streamID); err != nil {
		return nil, err
	}

	if len(streams) != 1 {
		return nil, nil
	}

	stream := streams[0]

	svc := s3.New(c.options.Session)

	signed, err := SignS3URL(svc, stream.URL)
	if err != nil {
		return nil, err
	}

	streamer := &Streamer{
		Stream:    stream,
		SignedURL: signed,
	}

	return streamer, nil
}

func (c *StreamsController) Binary(ctx *app.BinaryStreamsContext) error {
	log := Logger(ctx).Sugar()

	streamer, err := c.lookupDeviceStream(ctx, ctx.StreamID)
	if err != nil {
		return err
	}
	if streamer == nil {
		return ctx.NotFound()
	}

	resp, err := http.Get(streamer.SignedURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.NotFound()
	}

	ctx.ResponseData.WriteHeader(http.StatusOK)

	n, err := io.Copy(ctx.ResponseData, resp.Body)
	if err != nil {
		return err
	}

	log.Infow("Stream sent", "bytes", n, "size", streamer.Stream.Size)

	return nil
}

func (c *StreamsController) Raw(ctx *app.RawStreamsContext) error {
	log := Logger(ctx).Sugar()

	streamer, err := c.lookupDeviceStream(ctx, ctx.StreamID)
	if err != nil {
		return err
	}
	if streamer == nil {
		return ctx.NotFound()
	}

	resp, err := http.Get(streamer.SignedURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.NotFound()
	}

	ctx.ResponseData.WriteHeader(http.StatusOK)

	n, err := io.Copy(ctx.ResponseData, resp.Body)
	if err != nil {
		return err
	}

	log.Infow("Stream sent", "bytes", n, "size", streamer.Stream.Size)

	return nil
}

type TempReceiver struct {
}

func (tr *TempReceiver) HandleRecord(ctx context.Context, r *pb.DataRecord) (error) {
	log := Logger(ctx).Sugar()

	log.Infow("Record", "record", r)

	return nil
}

func (tr *TempReceiver) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	return nil, nil
}

func (c *StreamsController) Csv(ctx *app.CsvStreamsContext) error {
	log := Logger(ctx).Sugar()

	streamer, err := c.lookupDeviceStream(ctx, ctx.StreamID)
	if err != nil {
		return err
	}
	if streamer == nil {
		return ctx.NotFound()
	}

	resp, err := http.Get(streamer.SignedURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	receiver := &TempReceiver{}
	binaryReader := backend.NewFkBinaryReader(receiver)

	err = binaryReader.Read(ctx, resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return ctx.NotFound()
	}

	ctx.ResponseData.WriteHeader(http.StatusOK)

	n, err := io.Copy(ctx.ResponseData, resp.Body)
	if err != nil {
		return err
	}

	log.Infow("Stream sent", "bytes", n, "size", streamer.Stream.Size)

	return nil
}
