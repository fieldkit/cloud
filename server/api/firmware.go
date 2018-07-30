package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type FirmwareControllerOptions struct {
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

func (c *FirmwareController) Check(ctx *app.CheckFirmwareContext) error {
	log := Logger(ctx).Sugar()

	incomingETag := stripQuotes(ctx.IfNoneMatch)

	log.Infow("Device", "device_id", ctx.DeviceID, "incoming_etag", incomingETag)

	firmwares := []*data.DeviceFirmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares, "SELECT f.* FROM fieldkit.device_firmware AS f JOIN fieldkit.device AS d ON f.device_id = d.source_id WHERE d.key = $1 ORDER BY time DESC LIMIT 1", ctx.DeviceID); err != nil {
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

	resp, err := http.Get(fw.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

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
	log.Infow("Device", "device_id", ctx.Payload.DeviceID, "etag", ctx.Payload.Etag, "url", ctx.Payload.URL)

	device, err := c.options.Backend.GetDeviceSourceByID(ctx, int32(ctx.Payload.DeviceID))
	if err != nil {
		return err
	}

	firmware := data.DeviceFirmware{
		DeviceID: int64(device.ID),
		URL:      ctx.Payload.URL,
		ETag:     ctx.Payload.Etag,
		Time:     time.Now(),
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.device_firmware (device_id, time, url, etag)
		   VALUES (:device_id, :time, :url, :etag)
		   `, firmware); err != nil {
		return err
	}

	return ctx.OK([]byte("OK"))
}

func (c *FirmwareController) Add(ctx *app.AddFirmwareContext) error {
	log := Logger(ctx).Sugar()
	log.Infow("Device", "etag", ctx.Payload.Etag, "url", ctx.Payload.URL)

	firmware := data.Firmware{
		URL:  ctx.Payload.URL,
		ETag: ctx.Payload.Etag,
		Time: time.Now(),
	}

	if _, err := c.options.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.firmware (time, url, etag)
		   VALUES (:time, :url, :etag)
		   `, firmware); err != nil {
		return err
	}

	return ctx.OK([]byte("OK"))
}
