package api

import (
	"context"
	"errors"
	_ "fmt"
	_ "strconv"
	_ "strings"
	_ "time"

	_ "github.com/lib/pq"

	_ "github.com/jmoiron/sqlx"

	"github.com/conservify/sqlxcache"

	"goa.design/goa/v3/security"

	csvService "github.com/fieldkit/cloud/server/api/gen/csv"

	_ "github.com/fieldkit/cloud/server/backend/handlers"
	_ "github.com/fieldkit/cloud/server/backend/repositories"
	_ "github.com/fieldkit/cloud/server/data"
)

func NewRawQueryParamsFromCsvExport(payload *csvService.ExportPayload) (*RawQueryParams, error) {
	return &RawQueryParams{
		Start:      payload.Start,
		End:        payload.End,
		Resolution: payload.Resolution,
		Stations:   payload.Stations,
		Sensors:    payload.Sensors,
		Aggregate:  payload.Aggregate,
		Tail:       payload.Tail,
		Complete:   payload.Complete,
	}, nil
}

type CsvService struct {
	options *ControllerOptions
	db      *sqlxcache.DB
}

func NewCsvService(ctx context.Context, options *ControllerOptions) *CsvService {
	return &CsvService{
		options: options,
		db:      options.Database,
	}
}

func (c *CsvService) Export(ctx context.Context, payload *csvService.ExportPayload) (*csvService.ExportResult, error) {
	log := Logger(ctx).Sugar()

	rawParams, err := NewRawQueryParamsFromCsvExport(payload)
	if err != nil {
		return nil, err
	}

	qp, err := rawParams.BuildQueryParams()
	if err != nil {
		return nil, err
	}

	if len(qp.Sensors) == 0 {
		return nil, errors.New("sensors is required")
	}

	log.Infow("query_parameters", "start", qp.Start, "end", qp.End, "sensors", qp.Sensors, "stations", qp.Stations, "resolution", qp.Resolution, "aggregate", qp.Aggregate)

	aggregateName := "1m"
	aqp, err := NewAggregateQueryParams(qp, aggregateName, nil)
	if err != nil {
		return nil, err
	}

	dq := NewDataQuerier(c.db)

	queried, err := dq.QueryAggregate(ctx, aqp)
	if err != nil {
		return nil, err
	}

	defer queried.Close()

	for queried.Next() {
		row := &DataRow{}
		if err = queried.StructScan(row); err != nil {
			return nil, err
		}
	}

	return &csvService.ExportResult{
		Object: nil,
	}, nil
}

func (s *CsvService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     nil,
		Unauthorized: func(m string) error { return csvService.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return csvService.MakeForbidden(errors.New(m)) },
	})
}
