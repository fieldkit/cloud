package api

import (
	"context"
	"encoding/json"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
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
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	data_records := make([]*data.DataRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &data_records, `SELECT * FROM fieldkit.data_record WHERE (id = $1)`, ctx.RecordID); err != nil {
		return err
	}

	if len(data_records) == 0 {
		return ctx.NotFound()
	}

	meta_records := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &meta_records, `SELECT * FROM fieldkit.meta_record WHERE (id = $1)`, ctx.RecordID); err != nil {
		return err
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
	p, err := NewPermissions(ctx)
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
