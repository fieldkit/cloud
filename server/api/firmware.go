package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"goa.design/goa/v3/security"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/backend"
	_ "github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/data"

	firmware "github.com/fieldkit/cloud/server/api/gen/firmware"
)

type FirmwareService struct {
	options *ControllerOptions
}

func NewFirmwareService(ctx context.Context, options *ControllerOptions) *FirmwareService {
	return &FirmwareService{
		options: options,
	}
}

func FirmwareSummaryType(fw *data.Firmware) *firmware.FirmwareSummary {
	buildNumber := int32(0)
	buildTime := int64(0)
	metaFields, _ := fw.GetMeta()
	if metaFields != nil {
		if raw, ok := metaFields["Build-Number"]; ok {
			if v, err := strconv.ParseInt(raw.(string), 10, 32); err == nil {
				buildNumber = int32(v)
			}
		}
		if raw, ok := metaFields["Build-Time"]; ok {
			if p, err := time.Parse("20060102_150405", raw.(string)); err == nil {
				buildTime = p.Unix()
			}
		}
	}

	return &firmware.FirmwareSummary{
		ID:          fw.ID,
		Time:        fw.Time.String(),
		Module:      fw.Module,
		Profile:     fw.Profile,
		Etag:        fw.ETag,
		URL:         fw.URL,
		Meta:        metaFields,
		BuildNumber: buildNumber,
		BuildTime:   buildTime,
	}
}

func FirmwareSummariesType(firmwares []*data.Firmware) []*firmware.FirmwareSummary {
	summaries := make([]*firmware.FirmwareSummary, len(firmwares))
	for i, summary := range firmwares {
		summaries[i] = FirmwareSummaryType(summary)
	}
	return summaries
}

func FirmwaresType(firmwares []*data.Firmware) *firmware.Firmwares {
	return &firmware.Firmwares{
		Firmwares: FirmwareSummariesType(firmwares),
	}
}

func (s *FirmwareService) Download(ctx context.Context, payload *firmware.DownloadPayload) (*firmware.DownloadResult, io.ReadCloser, error) {
	log := Logger(ctx).Sugar()

	firmwares := []*data.Firmware{}
	if err := s.options.Database.SelectContext(ctx, &firmwares, `
		SELECT * FROM fieldkit.firmware WHERE id = $1
		`, payload.FirmwareID); err != nil {
		return nil, nil, err
	}

	if len(firmwares) == 0 {
		log.Errorw("firmware missing", "firmware_id", payload.FirmwareID)
		return nil, nil, firmware.MakeNotFound(errors.New("not found"))
	}

	fw := firmwares[0]

	bucketAndKey, err := common.GetBucketAndKey(fw.URL)
	if err != nil {
		return nil, nil, err
	}

	goi := &s3.GetObjectInput{
		Bucket: aws.String(bucketAndKey.Bucket),
		Key:    aws.String(bucketAndKey.Key),
	}

	svc := s3.New(s.options.Session)

	obj, err := svc.GetObject(goi)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading object %v: %v", bucketAndKey.Key, err)
	}

	contentLength := int64(0)

	if obj.ContentLength != nil {
		contentLength = int64(*obj.ContentLength)
	}

	log.Infow("firmware sent", "bytes", contentLength)

	return &firmware.DownloadResult{
		Length:      contentLength,
		ContentType: "application",
	}, obj.Body, nil
}

func (s *FirmwareService) Add(ctx context.Context, payload *firmware.AddPayload) error {
	log := Logger(ctx).Sugar()

	metaMap := make(map[string]string)
	err := json.Unmarshal([]byte(payload.Firmware.Meta), &metaMap)
	if err != nil {
		return err
	}

	log.Infow("add firmware", "etag", payload.Firmware.Etag, "url", payload.Firmware.URL, "module", payload.Firmware.Module, "profile", payload.Firmware.Profile, "meta", metaMap)

	fw := data.Firmware{
		Time:    time.Now(),
		Module:  payload.Firmware.Module,
		Profile: payload.Firmware.Profile,
		URL:     payload.Firmware.URL,
		ETag:    payload.Firmware.Etag,
		Meta:    []byte(payload.Firmware.Meta),
	}

	if _, err := s.options.Database.NamedExecContext(ctx, `
		   INSERT INTO fieldkit.firmware (time, module, profile, url, etag, meta)
		   VALUES (:time, :module, :profile, :url, :etag, :meta)
		   `, fw); err != nil {
		return err
	}

	return nil
}

func (s *FirmwareService) List(ctx context.Context, payload *firmware.ListPayload) (*firmware.Firmwares, error) {
	log := Logger(ctx).Sugar()

	firmwareTester := false

	p, err := NewPermissions(ctx, s.options).Unwrap()
	if err == nil {
		user := &data.User{}
		if err := s.options.Database.GetContext(ctx, user, `SELECT * FROM fieldkit.user WHERE id = $1`, p.UserID()); err != nil {
			return nil, err
		}

		log.Infow("firmware", "user", user.ID, "firmware_tester", user.FirmwareTester)

		firmwareTester = user.FirmwareTester
	} else {
		log.Infow("firmware", "user", "anonymous")
	}

	page := int32(0)
	if payload.Page != nil {
		page = *payload.Page
	}

	pageSize := int32(10)
	if payload.PageSize != nil {
		pageSize = *payload.PageSize
	}

	firmwares := []*data.Firmware{}
	if err := s.options.Database.SelectContext(ctx, &firmwares, `
		SELECT f.*
		FROM fieldkit.firmware AS f
		WHERE (f.module = $1 OR $1 IS NULL) AND
			  (f.profile = $2 OR $2 IS NULL) AND
			  (f.available OR $5)
		ORDER BY time DESC LIMIT $3 OFFSET $4
		`, payload.Module, payload.Profile, pageSize, page*pageSize, firmwareTester); err != nil {
		return nil, err
	}

	svc := s3.New(s.options.Session)

	if false {
		for _, f := range firmwares {
			signed, err := backend.SignS3URL(svc, f.URL)
			if err != nil {
				return nil, err
			}
			f.URL = signed
		}
	} else {
		for _, f := range firmwares {
			f.URL = fmt.Sprintf("/firmware/%d/download", f.ID)
		}
	}

	return FirmwaresType(firmwares), nil
}

func (s *FirmwareService) Delete(ctx context.Context, payload *firmware.DeletePayload) error {
	log := Logger(ctx).Sugar()

	log.Infow("deleting", "firmware_id", payload.FirmwareID)

	firmwares := []*data.Firmware{}
	if err := s.options.Database.SelectContext(ctx, &firmwares, `
		SELECT * FROM fieldkit.firmware WHERE id = $1
		`, payload.FirmwareID); err != nil {
		return err
	}

	if len(firmwares) == 0 {
		return firmware.MakeNotFound(errors.New("not found"))
	}

	svc := s3.New(s.options.Session)

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

		if _, err := s.options.Database.ExecContext(ctx, `
			DELETE FROM fieldkit.firmware WHERE id = $1
			`, payload.FirmwareID); err != nil {
			return err
		}
	}

	return nil
}

func (s *FirmwareService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return firmware.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return firmware.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return firmware.MakeForbidden(errors.New(m)) },
	})
}
