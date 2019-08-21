package api

import (
	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
)

type StationLogControllerOptions struct {
	Database *sqlxcache.DB
}

func StationLogType(stationLog *data.StationLog) *app.StationLog {
	return &app.StationLog{
		ID:        int(stationLog.ID),
		Body:      stationLog.Body,
		Timestamp: stationLog.Timestamp,
		StationID: int(stationLog.StationID),
	}
}

func StationLogsType(stationLogs []*data.StationLog) *app.StationLogs {
	stationLogsCollection := make([]*app.StationLog, len(stationLogs))
	for i, stationLog := range stationLogs {
		stationLogsCollection[i] = StationLogType(stationLog)
	}
	return &app.StationLogs{
		StationLogs: stationLogsCollection,
	}
}

type StationLogController struct {
	*goa.Controller
	options StationLogControllerOptions
}

func NewStationLogController(service *goa.Service, options StationLogControllerOptions) *StationLogController {
	return &StationLogController{
		Controller: service.NewController("StationLogController"),
		options:    options,
	}
}

func (c *StationLogController) Add(ctx *app.AddStationLogContext) error {
	stationLog := &data.StationLog{
		ID:        int32(ctx.Payload.ID),
		StationID: int32(ctx.Payload.StationID),
		Body:      ctx.Payload.Body,
		Timestamp: ctx.Payload.Timestamp,
	}

	if err := c.options.Database.NamedGetContext(ctx, stationLog, "INSERT INTO fieldkit.station_log (name, body, timestamp, station_id) VALUES (:name, :body, :timestamp, :station_id) RETURNING *", stationLog); err != nil {
		return err
	}

	return ctx.OK(StationLogType(stationLog))
}

func (c *StationLogController) AddMultiple(ctx *app.AddMultipleStationLogContext) error {
	for _, addLog := range ctx.Payload.StationLogs {
		stationLog := &data.StationLog{
			StationID: int32(addLog.StationID),
			Body:      addLog.Body,
			Timestamp: addLog.Timestamp,
		}

		if err := c.options.Database.NamedGetContext(ctx, stationLog, "INSERT INTO fieldkit.station_log (station_id, body, timestamp) VALUES (:id, :station_id, :body, :timestamp) RETURNING *", stationLog); err != nil {
			return err
		}
	}

	return nil
}

func (c *StationLogController) Update(ctx *app.UpdateStationLogContext) error {
	stationLog := &data.StationLog{
		ID:        int32(ctx.StationLogID),
		StationID: int32(ctx.Payload.StationID),
		Body:      ctx.Payload.Body,
		Timestamp: ctx.Payload.Timestamp,
	}
	if err := c.options.Database.NamedGetContext(ctx, stationLog, "UPDATE fieldkit.station_log SET name = :name, station_id = :station_id, body = :body, timestamp = :timestamp, WHERE id = :id RETURNING *", stationLog); err != nil {
		return err
	}

	return ctx.OK(StationLogType(stationLog))
}

func (c *StationLogController) Get(ctx *app.GetStationLogContext) error {
	stationLog := &data.StationLog{}
	if err := c.options.Database.GetContext(ctx, stationLog, "SELECT * FROM fieldkit.station_log WHERE name = $1", ctx.StationLog); err != nil {
		return err
	}

	return ctx.OK(StationLogType(stationLog))
}
