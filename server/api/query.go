package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goadesign/goa"

	"github.com/Conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/data"
)

type QueryControllerOptions struct {
	Database *sqlxcache.DB
	Backend  *backend.Backend
}

type QueryController struct {
	*goa.Controller
	options QueryControllerOptions
}

func NewQueryController(service *goa.Service, options QueryControllerOptions) *QueryController {
	return &QueryController{
		Controller: service.NewController("QueryController"),
		options:    options,
	}
}

type QueryCriteria struct {
	Keys      []string
	SourceID  int
	StartTime time.Time
	EndTime   time.Time
}

func ExtractQueryCriteria(ctx *app.ListBySourceQueryContext) (criteria *QueryCriteria, err error) {
	startTime, err := strconv.ParseInt(ctx.RequestData.Params.Get("startTime"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing startTime: %v", err)
	}

	endTime, err := strconv.ParseInt(ctx.RequestData.Params.Get("endTime"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing endTime: %v", err)
	}

	keys := strings.Split(ctx.RequestData.Params.Get("keys"), ",")

	criteria = &QueryCriteria{
		SourceID:  ctx.SourceID,
		StartTime: time.Unix(startTime, 0),
		EndTime:   time.Unix(endTime, 0),
		Keys:      keys,
	}

	return
}

func (c *QueryController) ListBySource(ctx *app.ListBySourceQueryContext) error {
	criteria, err := ExtractQueryCriteria(ctx)
	if err != nil {
		return err
	}

	pageSize := 200
	page := 0
	series := make([]*app.SeriesData, 0)

	for _, key := range criteria.Keys {
		series = append(series, &app.SeriesData{
			Name: key,
			Rows: make([]interface{}, 0),
		})
	}

	for {
		records := []*data.Record{}

		if err := c.options.Database.SelectContext(ctx, &records,
			`SELECT r.id, r.timestamp, r.data FROM fieldkit.record_visible AS r WHERE r.source_id = $3 AND r.timestamp BETWEEN $4 AND $5 ORDER BY r.timestamp LIMIT $1 OFFSET $2`,
			pageSize, page*pageSize, criteria.SourceID, criteria.StartTime, criteria.EndTime); err != nil {
			return err
		}

		if len(records) == 0 {
			break
		}

		for _, record := range records {
			timestamp := record.Timestamp
			parsed, err := record.GetParsedFields()
			if err != nil {
				return err
			}
			for seriesIndex, key := range criteria.Keys {
				if parsed[key] != nil {
					sample := []interface{}{
						timestamp.Unix(),
						parsed[key],
					}
					series[seriesIndex].Rows = append(series[seriesIndex].Rows, sample)
				}
			}
		}

		page += 1
	}

	qd := &app.QueryData{
		Series: series,
	}

	return ctx.OK(qd)
}
