package api

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/conservify/sqlxcache"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
)

type FieldNoteControllerOptions struct {
	Database *sqlxcache.DB
}

func FieldNoteQueryResultType(fieldNote *data.FieldNoteQueryResult) *app.FieldNoteQueryResult {
	fieldNoteQueryResultType := &app.FieldNoteQueryResult{
		ID:          int(fieldNote.ID),
		Created:     fieldNote.Created,
		UserID:      int(fieldNote.UserID),
		CategoryKey: fieldNote.CategoryKey,
		Username:    fieldNote.Username,
	}

	if fieldNote.Note != nil {
		fieldNoteQueryResultType.Note = fieldNote.Note
	}

	if fieldNote.MediaURL != nil {
		fieldNoteQueryResultType.MediaURL = fieldNote.MediaURL
	}

	if fieldNote.MediaContentType != nil {
		fieldNoteQueryResultType.MediaContentType = fieldNote.MediaContentType
	}

	return fieldNoteQueryResultType
}

func FieldNoteType(fieldNote *data.FieldNote) *app.FieldNote {
	fieldNoteType := &app.FieldNote{
		ID:          int(fieldNote.ID),
		StationID:   int(fieldNote.StationID),
		Created:     fieldNote.Created,
		UserID:      int(fieldNote.UserID),
		CategoryID:  int(fieldNote.CategoryID),
	}

	if fieldNote.Note != nil {
		fieldNoteType.Note = *fieldNote.Note
	}

	if fieldNote.MediaID != nil {
		mediaID := int(*fieldNote.MediaID)
		fieldNoteType.MediaID = &mediaID
	}

	return fieldNoteType
}

func FieldNotesType(fieldNotes []*data.FieldNoteQueryResult) *app.FieldNotes {
	fieldNotesCollection := make([]*app.FieldNoteQueryResult, len(fieldNotes))
	for i, fieldNote := range fieldNotes {
		fieldNotesCollection[i] = FieldNoteQueryResultType(fieldNote)
	}

	return &app.FieldNotes{
		Notes: fieldNotesCollection,
	}
}

type FieldNoteController struct {
	*goa.Controller
	options FieldNoteControllerOptions
}

func NewFieldNoteController(service *goa.Service, options FieldNoteControllerOptions) *FieldNoteController {
	return &FieldNoteController{
		Controller: service.NewController("FieldNoteController"),
		options:    options,
	}
}

func (c *FieldNoteController) Add(ctx *app.AddFieldNoteContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	// TODO: use true default category ID
	categoryId := 0
	if ctx.Payload.CategoryID != nil {
		categoryId = *ctx.Payload.CategoryID
	}

	fieldNote := &data.FieldNote{
		StationID:   int32(ctx.Payload.StationID),
		Created:     ctx.Payload.Created,
		UserID:      int32(ctx.Payload.UserID),
		CategoryID:  int32(categoryId),
	}

	if ctx.Payload.Note != nil {
		note := string(*ctx.Payload.Note)
		fieldNote.Note = &note
	} else {
		fieldNote.Note = nil
	}

	if ctx.Payload.MediaID != nil {
		mediaID := int32(*ctx.Payload.MediaID)
		fieldNote.MediaID = &mediaID
	} else {
		fieldNote.MediaID = nil
	}

	if err := c.options.Database.NamedGetContext(ctx, fieldNote, "INSERT INTO fieldkit.field_note (station_id, created, user_id, category_id, note, media_id) VALUES (:station_id, :created, :user_id, :category_id, :note, :media_id) RETURNING *", fieldNote); err != nil {
		return err
	}

	return ctx.OK(FieldNoteType(fieldNote))
}

func (c *FieldNoteController) Update(ctx *app.UpdateFieldNoteContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	// TODO: use true default category ID
	categoryId := 0
	if ctx.Payload.CategoryID != nil {
		categoryId = *ctx.Payload.CategoryID
	}

	fieldNote := &data.FieldNote{
		ID:          int32(ctx.FieldNoteID),
		StationID:   int32(ctx.Payload.StationID),
		Created:     ctx.Payload.Created,
		UserID:      int32(ctx.Payload.UserID),
		CategoryID:  int32(categoryId),
	}

	if ctx.Payload.Note != nil {
		note := string(*ctx.Payload.Note)
		fieldNote.Note = &note
	} else {
		fieldNote.Note = nil
	}

	if ctx.Payload.MediaID != nil {
		mediaID := int32(*ctx.Payload.MediaID)
		fieldNote.MediaID = &mediaID
	} else {
		fieldNote.MediaID = nil
	}

	if err := c.options.Database.NamedGetContext(ctx, fieldNote, "UPDATE fieldkit.field_note SET station_id = :station_id, created = :created, user_id = :user_id, category_id = :category_id, note = :note, media_id = :media_id WHERE id = :id RETURNING *", fieldNote); err != nil {
		return err
	}

	return ctx.OK(FieldNoteType(fieldNote))
}

func (c *FieldNoteController) Get(ctx *app.GetFieldNoteContext) error {
	fieldNotes := []*data.FieldNoteQueryResult{}

	if err := c.options.Database.SelectContext(ctx, &fieldNotes, "SELECT fn.id AS ID, fn.created, fn.user_id, fn.note, c.key AS CategoryKey, u.username, m.url AS MediaURL, m.content_type AS MediaContentType FROM fieldkit.field_note AS fn JOIN fieldkit.user u ON (u.id = fn.user_id) JOIN fieldkit.field_note_category AS c ON (c.id = fn.category_id) LEFT OUTER JOIN fieldkit.field_note_media AS m ON (m.id = fn.media_id) WHERE fn.station_id = $1", ctx.StationID); err != nil {
		return err
	}

	return ctx.OK(FieldNotesType(fieldNotes))
}

func (c *FieldNoteController) Delete(ctx *app.DeleteFieldNoteContext) error {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.field_note WHERE id = $1", ctx.FieldNoteID); err != nil {
		return err
	}

	return ctx.OK()
}
