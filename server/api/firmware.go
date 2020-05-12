package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/goadesign/goa"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"
)

type FirmwareController struct {
	*goa.Controller
	options *ControllerOptions
}

func NewFirmwareController(service *goa.Service, options *ControllerOptions) *FirmwareController {
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

func FirmwareSummaryType(fw *data.Firmware) *app.FirmwareSummary {
	buildNumber := 0
	buildTime := 0
	metaFields, _ := fw.GetMeta()
	if metaFields != nil {
		if raw, ok := metaFields["Build-Number"]; ok {
			if v, err := strconv.ParseInt(raw.(string), 10, 32); err == nil {
				buildNumber = int(v)
			}
		}
		if raw, ok := metaFields["Build-Time"]; ok {
			if p, err := time.Parse("20060102_150405", raw.(string)); err == nil {
				buildTime = int(p.Unix())
			}
		}
	}

	return &app.FirmwareSummary{
		ID:          int(fw.ID),
		Time:        fw.Time,
		Module:      fw.Module,
		Profile:     fw.Profile,
		Etag:        fw.ETag,
		URL:         fw.URL,
		Meta:        metaFields,
		BuildNumber: buildNumber,
		BuildTime:   buildTime,
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

func (c *FirmwareController) Download(ctx *app.DownloadFirmwareContext) error {
	log := Logger(ctx).Sugar()

	firmwares := []*data.Firmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares, `SELECT * FROM fieldkit.firmware WHERE id = $1`, ctx.FirmwareID); err != nil {
		return err
	}

	if len(firmwares) == 0 {
		log.Errorw("firmware missing", "firmware_id", ctx.FirmwareID)
		return ctx.NotFound()
	}

	fw := firmwares[0]

	bucketAndKey, err := common.GetBucketAndKey(fw.URL)
	if err != nil {
		return err
	}

	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucketAndKey.Bucket),
		Key:    aws.String(bucketAndKey.Key),
	}

	svc := s3.New(c.options.Session)

	obj, err := svc.GetObject(goi)
	if err != nil {
		return fmt.Errorf("error reading object %v: %v", bucketAndKey.Key, err)
	}

	contentLength := 0

	if obj.ContentLength != nil {
		contentLength = int(*obj.ContentLength)
	}

	ctx.ResponseData.Header().Set("ETag", fmt.Sprintf("\"%s\"", fw.ETag))
	ctx.ResponseData.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	ctx.ResponseData.WriteHeader(http.StatusOK)

	n, err := io.Copy(ctx.ResponseData, obj.Body)
	if err != nil {
		return err
	}

	log.Infow("firmware sent", "bytes", n)

	return nil
}

func (c *FirmwareController) Add(ctx *app.AddFirmwareContext) error {
	log := Logger(ctx).Sugar()

	metaMap := make(map[string]string)
	err := json.Unmarshal([]byte(ctx.Payload.Meta), &metaMap)
	if err != nil {
		return err
	}

	log.Infow("add firmware", "etag", ctx.Payload.Etag, "url", ctx.Payload.URL, "module", ctx.Payload.Module, "profile", ctx.Payload.Profile, "meta", metaMap)

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

func (c *FirmwareController) List(ctx *app.ListFirmwareContext) error {
	log := Logger(ctx).Sugar()

	firmwareTester := false

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err == nil {
		user := &data.User{}
		if err := c.options.Database.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE id = $1`, p.UserID()); err != nil {
			return err

		}

		log.Infow("firmware", "user", user.ID, "firmware_tester", user.FirmwareTester)

		firmwareTester = user.FirmwareTester
	} else {
		log.Infow("firmware", "user", "anonymous")
	}

	page := 0
	if ctx.Page != nil {
		page = *ctx.Page
	}

	pageSize := 10
	if ctx.PageSize != nil {
		pageSize = *ctx.PageSize
	}

	firmwares := []*data.Firmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares, `
		SELECT f.*
		FROM fieldkit.firmware AS f
		WHERE (f.module = $1 OR $1 IS NULL) AND
			  (f.profile = $2 OR $2 IS NULL) AND
			  (f.available OR $5)
		ORDER BY time DESC LIMIT $3 OFFSET $4
		`, ctx.Module, ctx.Profile, pageSize, page*pageSize, firmwareTester); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	if false {
		for _, f := range firmwares {
			signed, err := backend.SignS3URL(svc, f.URL)
			if err != nil {
				return err
			}
			f.URL = signed
		}
	} else {
		for _, f := range firmwares {
			f.URL = fmt.Sprintf("/firmware/%d/download", f.ID)
		}
	}

	return ctx.OK(FirmwaresType(firmwares))
}

func (c *FirmwareController) Delete(ctx *app.DeleteFirmwareContext) error {
	log := Logger(ctx).Sugar()

	log.Infow("deleting", "firmware_id", ctx.FirmwareID)

	firmwares := []*data.Firmware{}
	if err := c.options.Database.SelectContext(ctx, &firmwares, `SELECT * FROM fieldkit.firmware WHERE id = $1`, ctx.FirmwareID); err != nil {
		return err
	}

	if len(firmwares) == 0 {
		return ctx.NotFound()
	}

	svc := s3.New(c.options.Session)

	for _, fw := range firmwares {
		object, err := common.GetBucketAndKey(fw.URL)
		if err != nil {
			return fmt.Errorf("error parsing URL: %v", err)
		}

		log.Infow("deleting", "url", fw.URL)

		_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(object.Bucket), Key: aws.String(object.Key)})
		if err != nil {
			return fmt.Errorf("unable to delete object %q from bucket %q, %v", object.Key, object.Bucket, err)
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(object.Bucket),
			Key:    aws.String(object.Key),
		})
		if err != nil {
			return err
		}

		if _, err := c.options.Database.ExecContext(ctx, `DELETE FROM fieldkit.firmware WHERE id = $1`, ctx.FirmwareID); err != nil {
			return err
		}
	}

	return ctx.OK([]byte("{}"))
}
