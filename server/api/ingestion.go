package api

import (
	"context"
	"fmt"

	"goa.design/goa/v3/security"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	ingestion "github.com/fieldkit/cloud/server/api/gen/ingestion"

	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/common"
	"github.com/fieldkit/cloud/server/messages"
)

type IngestionService struct {
	options *ControllerOptions
}

func NewIngestionService(ctx context.Context, options *ControllerOptions) *IngestionService {
	return &IngestionService{
		options: options,
	}
}

func (c *IngestionService) ProcessPending(ctx context.Context, payload *ingestion.ProcessPendingPayload) (err error) {
	log := Logger(ctx).Sugar()

	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	ingestions, err := ir.QueryPending(ctx)
	if err != nil {
		return err
	}

	log.Infow("queueing", "ingestions", len(ingestions))

	for _, i := range ingestions {
		c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
			Time:   i.Time,
			ID:     i.ID,
			URL:    i.URL,
			UserID: p.UserID(),
		})
	}

	return nil
}

func (c *IngestionService) ProcessStation(ctx context.Context, payload *ingestion.ProcessStationPayload) (err error) {
	log := Logger(ctx).Sugar()

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("processing", "station_id", payload.StationID)

	p, err := NewPermissions(ctx, c.options).ForStationByID(int(payload.StationID))
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	ingestions, err := ir.QueryByStationID(ctx, int64(payload.StationID))
	if err != nil {
		return err
	}

	log.Infow("queueing", "ingestions", len(ingestions))

	for _, i := range ingestions {
		c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
			Time:   i.Time,
			ID:     i.ID,
			URL:    i.URL,
			UserID: p.UserID(),
		})
	}

	return nil
}

func (c *IngestionService) ProcessIngestion(ctx context.Context, payload *ingestion.ProcessIngestionPayload) (err error) {
	log := Logger(ctx).Sugar()

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("processing", "ingestion_id", payload.IngestionID)

	i, err := ir.QueryByID(ctx, int64(payload.IngestionID))
	if err != nil {
		return err
	}
	if i == nil {
		return ingestion.NotFound("not found")
	}

	p, err := NewPermissions(ctx, c.options).ForStationByDeviceID(i.DeviceID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	c.options.Publisher.Publish(ctx, &messages.IngestionReceived{
		Time:    i.Time,
		ID:      i.ID,
		URL:     i.URL,
		UserID:  p.UserID(),
		Verbose: true,
	})

	return nil
}

func (c *IngestionService) Delete(ctx context.Context, payload *ingestion.DeletePayload) (err error) {
	log := Logger(ctx).Sugar()

	ir, err := repositories.NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("deleting", "ingestion_id", payload.IngestionID)

	i, err := ir.QueryByID(ctx, int64(payload.IngestionID))
	if err != nil {
		return err
	}
	if i == nil {
		return ingestion.NotFound("not found")
	}

	p, err := NewPermissions(ctx, c.options).ForStationByDeviceID(i.DeviceID)
	if err != nil {
		return err
	}

	if err := p.CanModify(); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	object, err := common.GetBucketAndKey(i.URL)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %v", err)
	}

	log.Infow("deleting", "url", i.URL)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(object.Bucket), Key: aws.String(object.Key)})
	if err != nil {
		return fmt.Errorf("Unable to delete object %q from bucket %q, %v", object.Key, object.Bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(object.Bucket),
		Key:    aws.String(object.Key),
	})
	if err != nil {
		return err
	}

	if err := ir.Delete(ctx, int64(payload.IngestionID)); err != nil {
		return err
	}

	return nil
}

func (s *IngestionService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return ingestion.NotFound(m) },
		Unauthorized: func(m string) error { return ingestion.Unauthorized(m) },
		Forbidden:    func(m string) error { return ingestion.Forbidden(m) },
	})
}
