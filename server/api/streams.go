package api

import (
	_ "encoding/json"
	_ "fmt"
	_ "io"
	_ "net/http"
	_ "net/url"
	_ "strconv"
	_ "strings"
	_ "time"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"

	 "github.com/aws/aws-sdk-go/aws/session"
	 "github.com/aws/aws-sdk-go/service/s3"
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
		DeviceID : s.DeviceID,
		FileID   : s.FileID,
		Firmware : s.Firmware,
		ID       : int(s.ID),
		Meta     : s.Meta.String(),
		Size     : int(s.Size),
		StreamID : s.StreamID,
		Time     : s.Time,
		URL      : s.URL,
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
		`SELECT s.* FROM fieldkit.device_stream AS s ORDER BY time DESC LIMIT $1 OFFSET $2`, pageSize, offset); err != nil {
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
		`SELECT s.* FROM fieldkit.device_stream AS s WHERE s.device_id = $1 ORDER BY time DESC LIMIT $2 OFFSET $3`, ctx.DeviceID, pageSize, offset); err != nil {
		return err
	}

	return ctx.OK(DeviceStreamsType(streams))
}
