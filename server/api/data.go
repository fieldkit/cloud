package api

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type DataControllerOptions struct {
	Config   *ApiConfiguration
	Session  *session.Session
	Database *sqlxcache.DB
}

type DataController struct {
	options DataControllerOptions
	*goa.Controller
}

func NewDataController(ctx context.Context, service *goa.Service, options DataControllerOptions) *DataController {
	return &DataController{
		options:    options,
		Controller: service.NewController("DataController"),
	}
}

func (c *DataController) Process(ctx *app.ProcessDataContext) error {
	log := Logger(ctx).Sugar()

	ir, err := NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	pending, err := ir.QueryPending(ctx)
	if err != nil {
		return err
	}

	recordAdder := backend.NewRecordAdder(c.options.Session, c.options.Database)

	for _, i := range pending {
		log.Infow("pending", "ingestion", i)

		err := recordAdder.WriteRecords(ctx, i)
		if err != nil {
			log.Errorw("error", "error", err)
			err := ir.MarkProcessedHasErrors(ctx, i.ID)
			if err != nil {
				return err
			}
		} else {
			err := ir.MarkProcessedDone(ctx, i.ID)
			if err != nil {
				return err
			}
		}
	}

	log.Infow("done", "elapsed", 0)

	return nil
}

func (c *DataController) Delete(ctx *app.DeleteDataContext) error {
	log := Logger(ctx).Sugar()

	ir, err := NewIngestionRepository(c.options.Database)
	if err != nil {
		return err
	}

	log.Infow("deleting", "ingestion_id", ctx.IngestionID)

	i, err := ir.QueryById(ctx, int64(ctx.IngestionID))
	if err != nil {
		return err
	}
	if i == nil {
		return ctx.NotFound()
	}

	if err := ir.Delete(ctx, int64(ctx.IngestionID)); err != nil {
		return err
	}

	svc := s3.New(c.options.Session)

	object, err := backend.GetBucketAndKey(i.URL)
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

	return ctx.OK([]byte("deleted"))
}

func (c *DataController) Device(ctx *app.DeviceDataContext) error {
	log := Logger(ctx).Sugar()

	if false {
		ir, err := NewIngestionRepository(c.options.Database)
		if err != nil {
			return err
		}

		all, err := ir.QueryPending(ctx)
		if err != nil {
			return err
		}

		for _, i := range all {
			log.Infow("ingestion", "device_id", i.DeviceID)
		}
	}

	rr, err := NewRecordRepository(c.options.Database)
	if err != nil {
		return err
	}

	data, err := rr.QueryDevice(ctx, ctx.DeviceID)
	if err != nil {
		return err
	}

	log.Infow("data", "records", len(data))

	dataVms := make([]*app.DeviceDataRecord, 0)
	for _, r := range data {
		data, err := r.GetData()
		if err != nil {
			return err
		}
		dataVms = append(dataVms, &app.DeviceDataRecord{
			Time:   r.Time,
			Record: int(r.Number),
			Data:   data,
		})
	}

	return ctx.OK(&app.DeviceDataRecordsResponse{
		Summary: &app.DeviceDataStreamsSummary{},
		Meta:    []*app.DeviceMetaRecord{},
		Data:    dataVms,
	})
}

type IngestionRepository struct {
	Database *sqlxcache.DB
}

func NewIngestionRepository(database *sqlxcache.DB) (ir *IngestionRepository, err error) {
	return &IngestionRepository{Database: database}, nil
}

func (r *IngestionRepository) QueryById(ctx context.Context, id int64) (i *data.Ingestion, err error) {
	found := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &found, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.id = $1)`, id); err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (r *IngestionRepository) Delete(ctx context.Context, id int64) (err error) {
	if _, err := r.Database.ExecContext(ctx, `DELETE FROM fieldkit.ingestion WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) QueryPending(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.errors IS NULL) ORDER BY i.size ASC, i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) QueryAll(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i ORDER BY i.time DESC`); err != nil {
		return nil, err
	}
	return pending, nil
}

func (r *IngestionRepository) MarkProcessedHasErrors(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = true, attempted = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

func (r *IngestionRepository) MarkProcessedDone(ctx context.Context, id int64) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = false, completed = NOW() WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}

type RecordRepository struct {
	Database *sqlxcache.DB
}

func NewRecordRepository(database *sqlxcache.DB) (rr *RecordRepository, err error) {
	return &RecordRepository{Database: database}, nil
}

func (r *RecordRepository) QueryDevice(ctx context.Context, deviceId string) (all []*data.DataRecord, err error) {
	log := Logger(ctx).Sugar()

	pageSize := 200
	page := 0

	deviceIdBytes, err := base64.StdEncoding.DecodeString(deviceId)
	if err != nil {
		return nil, err
	}

	log.Infow("querying", "device_id", deviceIdBytes)

	data := []*data.DataRecord{}
	if err := r.Database.SelectContext(ctx, &data, `
	    SELECT r.*
	    FROM fieldkit.data_record AS r JOIN fieldkit.ingestion AS i ON (r.ingestion_id = i.id)
	    WHERE (i.device_id = $1) AND ((i.errors != true) OR (i.errors IS NULL))
	    ORDER BY r.time DESC LIMIT $2 OFFSET $3
	`, deviceIdBytes, pageSize, pageSize*page); err != nil {
		return nil, err
	}

	return data, nil
}
