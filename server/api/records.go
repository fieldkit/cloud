package api

import (
	"context"
	"errors"

	"goa.design/goa/v3/security"

	records "github.com/fieldkit/cloud/server/api/gen/records"

	"github.com/fieldkit/cloud/server/backend/repositories"
	serrors "github.com/fieldkit/cloud/server/common/errors"
	"github.com/fieldkit/cloud/server/data"
)

type RecordsService struct {
	options *ControllerOptions
}

func NewRecordsService(ctx context.Context, options *ControllerOptions) *RecordsService {
	return &RecordsService{
		options: options,
	}
}

func (c *RecordsService) Data(ctx context.Context, payload *records.DataPayload) (*records.DataResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	dataRecords := make([]*data.DataRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &dataRecords, `
		SELECT id, provision_id, time, number, meta_record_id, ST_AsBinary(location) AS location, raw, pb FROM fieldkit.data_record WHERE (id = $1)
		`, payload.RecordID); err != nil {
		return nil, serrors.Structured(err, "data_record_id", payload.RecordID)
	}

	if len(dataRecords) == 0 {
		return nil, records.MakeNotFound(errors.New("not found"))
	}

	metaRecords := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &metaRecords, `
		SELECT * FROM fieldkit.meta_record WHERE (id = $1)
		`, dataRecords[0].MetaRecordID); err != nil {
		return nil, serrors.Structured(err, "meta_record_id", dataRecords[0].MetaRecordID)
	}

	if len(metaRecords) == 0 {
		obj := struct {
			Data *data.DataRecord `json:"data"`
		}{
			dataRecords[0],
		}

		return &records.DataResult{
			Object: obj,
		}, nil
	}

	_ = p

	obj := struct {
		Data *data.DataRecord `json:"data"`
		Meta *data.MetaRecord `json:"meta"`
	}{
		dataRecords[0],
		metaRecords[0],
	}

	return &records.DataResult{
		Object: obj,
	}, nil
}

func (c *RecordsService) Meta(ctx context.Context, payload *records.MetaPayload) (*records.MetaResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	metaRecords := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &metaRecords, `SELECT * FROM fieldkit.meta_record WHERE (id = $1)`, payload.RecordID); err != nil {
		return nil, err
	}

	if len(metaRecords) == 0 {
		return nil, records.MakeNotFound(errors.New("not found"))
	}

	_ = p

	obj := struct {
		Meta *data.MetaRecord `json:"meta"`
	}{
		metaRecords[0],
	}

	return &records.MetaResult{
		Object: obj,
	}, nil
}

func (c *RecordsService) Resolved(ctx context.Context, payload *records.ResolvedPayload) (*records.ResolvedResult, error) {
	p, err := NewPermissions(ctx, c.options).Unwrap()
	if err != nil {
		return nil, err
	}

	dbDatas := make([]*data.DataRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &dbDatas, `
		SELECT id, provision_id, time, number, meta_record_id, ST_AsBinary(location) AS location, raw, pb FROM fieldkit.data_record WHERE (id = $1)
		`, payload.RecordID); err != nil {
		return nil, err
	}

	if len(dbDatas) == 0 {
		return nil, records.MakeNotFound(errors.New("not found"))
	}

	dbMetas := make([]*data.MetaRecord, 0)
	if err := c.options.Database.SelectContext(ctx, &dbMetas, `
		SELECT * FROM fieldkit.meta_record WHERE (id = $1)
		`, dbDatas[0].MetaRecordID); err != nil {
		return nil, err
	}

	metaFactory := repositories.NewMetaFactory()
	for _, dbMeta := range dbMetas {
		_, err := metaFactory.Add(ctx, dbMeta, false)
		if err != nil {
			return nil, err
		}
	}

	filtered, err := metaFactory.Resolve(ctx, dbDatas[0], true, false)
	if err != nil {
		return nil, err
	}

	_ = p

	obj := struct {
		Data     *data.DataRecord             `json:"data"`
		Meta     *data.MetaRecord             `json:"meta"`
		Filtered *repositories.FilteredRecord `json:"filtered"`
	}{
		dbDatas[0],
		dbMetas[0],
		filtered,
	}

	return &records.ResolvedResult{
		Object: obj,
	}, nil
}

func (s *RecordsService) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return Authenticate(ctx, AuthAttempt{
		Token:        token,
		Scheme:       scheme,
		Key:          s.options.JWTHMACKey,
		NotFound:     func(m string) error { return records.MakeNotFound(errors.New(m)) },
		Unauthorized: func(m string) error { return records.MakeUnauthorized(errors.New(m)) },
		Forbidden:    func(m string) error { return records.MakeForbidden(errors.New(m)) },
	})
}
