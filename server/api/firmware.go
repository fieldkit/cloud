package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FirmwareControllerOptions struct {
	Session  *session.Session
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type FirmwareController struct {
	*goa.Controller
	options FirmwareControllerOptions
}

func NewFirmwareController(service *goa.Service, options FirmwareControllerOptions) *FirmwareController {
	return &FirmwareController{
		Controller: service.NewController("FirmwareController"),
		options:    options,
	}
}

func stripQuotes(p *string) string {
	if p == nil {
		return ""
	}
	s := *p
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func isOutgoingFirmwareOlderThanIncoming(compiled string, fw *data.DeviceFirmware) bool {
	unix, err := strconv.Atoi(compiled)
	if err != nil {
		return false
	}
	incoming := time.Unix(int64(unix), 0)
	return incoming.After(fw.Time)
}

func (c *FirmwareController) Check(ctx *app.CheckFirmwareContext) error {
	log := Logger(ctx).Sugar()

	incomingETag := stripQuotes(ctx.IfNoneMatch)
	compiled := ctx.FkCompiled

	log.Infow("Device", "device_id", ctx.DeviceID, "module", ctx.Module, "incoming_etag", incomingETag, "compiled", compiled)

	firmwares := []*data.DeviceFirmware{}
	query := "SELECT f.* FROM fieldkit.device_firmware AS f JOIN fieldkit.device AS d ON f.device_id = d.source_id WHERE d.key = $1 AND f.module = $2 ORDER BY time DESC LIMIT 1"
	if err := c.options.Database.SelectContext(ctx, &firmwares, query, ctx.DeviceID, ctx.Module); err != nil {
		return err
	}

	if len(firmwares) == 0 {
		return ctx.NotFound()
	}

	fw := firmwares[0]

	log.Infow("Firmware", "time", fw.Time, "url", fw.URL, "etag", fw.ETag, "incoming_etag", incomingETag)

	if incomingETag == fw.ETag {
		return ctx.NotModified()
	}

	if compiled != nil {
		if isOutgoingFirmwareOlderThanIncoming(*compiled, fw) {
			log.Infow("Refusing to apply firmware compiled before device firmware", "incoming", compiled, "outgoing", fw.Time)
			return ctx.NotFound()
		}
	}

	resp, err := http.Get(fw.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.NotFound()
	}

	ctx.ResponseData.Header().Set("ETag", fmt.Sprintf("\"%s\"", fw.ETag))
	ctx.ResponseData.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	ctx.ResponseData.WriteHeader(http.StatusOK)

	n, err := io.Copy(ctx.ResponseData, resp.Body)
	if err != nil {
		return err
	}

	log.Infow("Firmware sent", "bytes", n)

	return nil
}

func (c *FirmwareController) Update(ctx *app.UpdateFirmwareContext) error {
	log := Logger(ctx).Sugar()
	log.Infow("Device", "device_id", ctx.Payload.DeviceID, "firmware_id", ctx.Payload.FirmwareID)

	device, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.Payload.DeviceID))
	if err != nil {
		return err
	}

	firmwares := []*data.Firmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares, "SELECT f.* FROM fieldkit.firmware AS f WHERE f.id = $1", ctx.Payload.FirmwareID); err != nil {
		return err
	}

	if len(firmwares) != 1 {
		return ctx.NotFound()
	}

	firmware := firmwares[0]

	deviceFirmware := data.DeviceFirmware{
		DeviceID: int64(device.ID),
		Time:     time.Now(),
		Module:   firmware.Module,
		Profile:  firmware.Profile,
		URL:      firmware.URL,
		ETag:     firmware.ETag,
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_firmware (device_id, time, module, profile, url, etag)
		   VALUES (:device_id, :time, :module, :profile, :url, :etag)
		   `, deviceFirmware); err != nil {
		return err
	}

	log.Infow("Update firmware", "device_id", ctx.Payload.DeviceID, "firmware_id", ctx.Payload.FirmwareID, "module", firmware.Module, "profile", firmware.Profile, "url", firmware.URL)

	return ctx.OK([]byte("OK"))
}

func (c *FirmwareController) Add(ctx *app.AddFirmwareContext) error {
	log := Logger(ctx).Sugar()

	metaMap := make(map[string]string)
	err := json.Unmarshal([]byte(ctx.Payload.Meta), &metaMap)
	if err != nil {
		return err
	}

	log.Infow("Firmware", "etag", ctx.Payload.Etag, "url", ctx.Payload.URL, "module", ctx.Payload.Module, "profile", ctx.Payload.Profile, "meta", metaMap)

	firmware := data.Firmware{
		Time:    time.Now(),
		Module:  ctx.Payload.Module,
		Profile: ctx.Payload.Profile,
		URL:     ctx.Payload.URL,
		ETag:    ctx.Payload.Etag,
		Meta:    []byte(ctx.Payload.Meta),
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.firmware (time, module, profile, url, etag, meta)
		   VALUES (:time, :module, :profile, :url, :etag, :meta)
		   `, firmware); err != nil {
		return err
	}

	return ctx.OK([]byte("OK"))
}

func FirmwareSummaryType(fw *data.Firmware) *app.FirmwareSummary {
	return &app.FirmwareSummary{
		ID:      int(fw.ID),
		Time:    fw.Time,
		Module:  fw.Module,
		Profile: fw.Profile,
		Etag:    fw.ETag,
		URL:     fw.URL,
	}
}

func FirmwareSummariesType(firmwares []*data.Firmware) []*app.FirmwareSummary {
	summaries := make([]*app.FirmwareSummary, len(firmwares))
	for i, summary := range firmwares {
		summaries[i] = FirmwareSummaryType(summary)
	}
	return summaries
}

func FirmwaresType(firmwares []*data.Firmware) *app.Firmwares {
	return &app.Firmwares{
		Firmwares: FirmwareSummariesType(firmwares),
	}
}

func getBucketAndKey(s3Url string) (bucket, key string, err error) {
	u, err := url.Parse(s3Url)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(u.Host, ".")

	return parts[0], u.Path[1:], nil
}

func signFirmwareURL(svc *s3.S3, f *data.Firmware) error {
	bucket, key, err := getBucketAndKey(f.URL)
	if err != nil {
		return err
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	signed, err := req.Presign(1 * time.Hour)
	if err != nil {
		return err
	}

	f.URL = signed

	return nil
}

func (c *FirmwareController) List(ctx *app.ListFirmwareContext) error {
	firmwares := []*data.Firmware{}

	if err := c.options.Database.SelectContext(ctx, &firmwares,
		`SELECT f.* FROM fieldkit.firmware AS f WHERE (f.module = $1 OR $1 IS NULL) AND (f.profile = $2 OR $2 IS NULL) ORDER BY time DESC LIMIT 100`, ctx.Module, ctx.Profile); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	for _, f := range firmwares {
		err := signFirmwareURL(svc, f)
		if err != nil {
			return err
		}
	}

	return ctx.OK(FirmwaresType(firmwares))
}

func (c *FirmwareController) ListDevice(ctx *app.ListDeviceFirmwareContext) error {
	log := Logger(ctx).Sugar()

	log.Infow("Device", "device_id", ctx.DeviceID)

	firmwares := []*data.Firmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares,
		`SELECT f.id, f.time, f.module, f.profile, f.etag, f.url FROM fieldkit.device_firmware AS f JOIN fieldkit.device AS d ON f.device_id = d.source_id WHERE d.key = $1 ORDER BY time DESC`, ctx.DeviceID); err != nil {
		return err
	}

	return ctx.OK(FirmwaresType(firmwares))
}
