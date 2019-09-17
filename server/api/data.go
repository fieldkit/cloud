package api

import (
	"context"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"

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

	pending, err := ir.GetPending(ctx)
	if err != nil {
		return err
	}

	recordAdder := backend.NewRecordAdder(c.options.Session, c.options.Database)

	for _, i := range pending {
		log.Infow("pending", "ingestion", i)

		err := recordAdder.WriteRecords(ctx, i)
		if err != nil {
			log.Errorw("error", "error", err)
			err := ir.MarkHasErrors(ctx, i.ID, true)
			if err != nil {
				return err
			}
		}
	}

	log.Infow("done", "elapsed", 0)

	return nil
}

type IngestionRepository struct {
	Database *sqlxcache.DB
}

func NewIngestionRepository(database *sqlxcache.DB) (ir *IngestionRepository, err error) {
	return &IngestionRepository{Database: database}, nil
}

func (r *IngestionRepository) GetPending(ctx context.Context) (all []*data.Ingestion, err error) {
	pending := []*data.Ingestion{}
	if err := r.Database.SelectContext(ctx, &pending, `SELECT i.* FROM fieldkit.ingestion AS i WHERE (i.errors IS NULL) ORDER BY i.time DESC`); err != nil {
		return nil, err
	}

	all = pending

	return
}

func (r *IngestionRepository) MarkHasErrors(ctx context.Context, id int64, errors bool) error {
	if _, err := r.Database.ExecContext(ctx, `UPDATE fieldkit.ingestion SET errors = $2 WHERE id = $1`, id, errors); err != nil {
		return err
	}

	return nil
}
