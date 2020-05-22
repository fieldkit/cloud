package api

import (
	"context"
	"encoding/json"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend/repositories"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/errors"
)

type RecordsController struct {
	options *ControllerOptions
	*goa.Controller
}

func NewRecordsController(ctx context.Context, service *goa.Service, options *ControllerOptions) *RecordsController {
	return &RecordsController{
		options:    options,
		Controller: service.NewController("RecordsController"),
	}
}

func writeJSON(responseData *goa.ResponseData, obj interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	responseData.Header().Set("Content-Type", "application/json")
	responseData.WriteHeader(200)
	responseData.Write(bytes)
	return nil
}

func (c *RecordsController) Data(ctx *app.DataRecordsContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	data_records := make([]*data.DataRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &data_records, `SELECT id, provision_id, time, number, meta_record_id, ST_AsBinary(location) AS location, raw FROM fieldkit.data_record WHERE (id = $1)`, ctx.RecordID); err != nil {
		return errors.Structured(err, "data_record_id", ctx.RecordID)
	}

	if len(data_records) == 0 {
		return ctx.NotFound()
	}

	meta_records := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &meta_records, `SELECT * FROM fieldkit.meta_record WHERE (id = $1)`, data_records[0].MetaRecordID); err != nil {
		return errors.Structured(err, "meta_record_id", data_records[0].MetaRecordID)
	}

	if len(meta_records) == 0 {
		return writeJSON(ctx.ResponseData, struct {
			Data *data.DataRecord `json:"data"`
		}{
			data_records[0],
		})
	}

	_ = p

	return writeJSON(ctx.ResponseData, struct {
		Data *data.DataRecord `json:"data"`
		Meta *data.MetaRecord `json:"meta"`
	}{
		data_records[0],
		meta_records[0],
	})
}

func (c *RecordsController) Meta(ctx *app.MetaRecordsContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	records := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &records, `SELECT * FROM fieldkit.meta_record WHERE (id = $1)`, ctx.RecordID); err != nil {
		return err
	}

	if len(records) == 0 {
		return ctx.NotFound()
	}

	_ = p

	return writeJSON(ctx.ResponseData, struct {
		Meta *data.MetaRecord `json:"meta"`
	}{
		records[0],
	})
}

func (c *RecordsController) Resolved(ctx *app.ResolvedRecordsContext) error {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return err
	}

	dbDatas := make([]*data.DataRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &dbDatas, `SELECT id, provision_id, time, number, meta_record_id, ST_AsBinary(location) AS location, raw FROM fieldkit.data_record WHERE (id = $1)`, ctx.RecordID); err != nil {
		return err
	}

	if len(dbDatas) == 0 {
		return ctx.NotFound()
	}

	dbMetas := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &dbMetas, `SELECT * FROM fieldkit.meta_record WHERE (id = $1)`, dbDatas[0].MetaRecordID); err != nil {
		return err
	}

	metaFactory := repositories.NewMetaFactory()
	for _, dbMeta := range dbMetas {
		_, err := metaFactory.Add(ctx, dbMeta)
		if err != nil {
			return err
		}
	}

	filtered, err := metaFactory.Resolve(ctx, dbDatas[0], true)
	if err != nil {
		return err
	}

	_ = p

	return writeJSON(ctx.ResponseData, struct {
		Data     *data.DataRecord             `json:"data"`
		Meta     *data.MetaRecord             `json:"meta"`
		Filtered *repositories.FilteredRecord `json:"filtered"`
	}{
		dbDatas[0],
		dbMetas[0],
		filtered,
	})
}

func (c *RecordsController) Filtered(ctx *app.FilteredRecordsContext) error {
	return writeJSON(ctx.ResponseData, struct{}{})
}
