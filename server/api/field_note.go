package api

import (
	"bufio"
	"io"
	"time"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/data"
	"github.com/fieldkit/cloud/server/backend/repositories"
)

type FieldNoteControllerOptions struct {
	Session *session.Session
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

func FieldNotesType(fieldNotes []*data.FieldNoteQueryResult) *app.FieldNotes {
	fieldNotesCollection := make([]*app.FieldNoteQueryResult, len(fieldNotes))
	for i, fieldNote := range fieldNotes {
		fieldNotesCollection[i] = FieldNoteQueryResultType(fieldNote)
	}

	return &app.FieldNotes{
		Notes: fieldNotesCollection,
	}
}

func FieldNoteMediaType(fieldNoteMedia *data.FieldNoteMedia) *app.FieldNoteMedia {
	return &app.FieldNoteMedia {
		ID:          int(fieldNoteMedia.ID),
		UserID:      int(fieldNoteMedia.UserID),
		ContentType: fieldNoteMedia.ContentType,
		Created:     fieldNoteMedia.Created,
		URL:         fieldNoteMedia.URL,
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

func (c *FieldNoteController) SaveMedia(ctx *app.SaveMediaFieldNoteContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(c.options.Session)
	saved, err := mr.Save(ctx, ctx.RequestData)
	if err != nil {
		return err
	}

	fieldNoteMedia := &data.FieldNoteMedia{
		Created:      time.Now(),
		UserID:       p.UserID,
		ContentType:  saved.MimeType,
		URL:          saved.URL,
	}

	if err := c.options.Database.NamedGetContext(ctx, fieldNoteMedia, "INSERT INTO fieldkit.field_note_media (user_id, content_type, created, url) VALUES (:user_id, :content_type, :created, :url) RETURNING *", fieldNoteMedia); err != nil {
		return err
	}

	return ctx.OK(FieldNoteMediaType(fieldNoteMedia))
}

func (c *FieldNoteController) GetMedia(ctx *app.GetMediaFieldNoteContext) error {
	_, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	fieldNoteMedia := &data.FieldNoteMedia{}
	if err := c.options.Database.GetContext(ctx, fieldNoteMedia, "SELECT * FROM fieldkit.field_note_media WHERE id = $1", ctx.MediaID); err != nil {
		return err
	}

	mr := repositories.NewMediaRepository(c.options.Session)

	lm, err := mr.LoadByURL(ctx, fieldNoteMedia.URL)
	if err != nil {
		return err
	}

	if lm != nil {
		writer := bufio.NewWriter(ctx.ResponseData)

		_, err = io.Copy(writer, lm.Reader)
		if err != nil {
			return err
		}
		return nil
	}

	return ctx.OK(nil)
}


func (c *FieldNoteController) Add(ctx *app.AddFieldNoteContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	// default category ID
	categoryId := 1
	if ctx.Payload.CategoryID != nil {
		categoryId = *ctx.Payload.CategoryID
	}

	fieldNote := &data.FieldNote{
		StationID:   int32(ctx.StationID),
		Created:     ctx.Payload.Created,
		UserID:      p.UserID,
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

	fieldNoteQueryResult := &data.FieldNoteQueryResult{}
	if err := c.options.Database.GetContext(ctx, fieldNoteQueryResult, "SELECT fn.id AS ID, fn.created, fn.user_id, fn.note, c.key AS CategoryKey, u.username, m.url AS MediaURL, m.content_type AS MediaContentType FROM fieldkit.field_note AS fn JOIN fieldkit.user u ON (u.id = fn.user_id) JOIN fieldkit.field_note_category AS c ON (c.id = fn.category_id) LEFT OUTER JOIN fieldkit.field_note_media AS m ON (m.id = fn.media_id) WHERE fn.id = $1", fieldNote.ID); err != nil {
		return err
	}

	return ctx.OK(FieldNoteQueryResultType(fieldNoteQueryResult))
}

func (c *FieldNoteController) Update(ctx *app.UpdateFieldNoteContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	// default category ID
	categoryId := 1
	if ctx.Payload.CategoryID != nil {
		categoryId = *ctx.Payload.CategoryID
	}

	fieldNote := &data.FieldNote{
		ID:          int32(ctx.FieldNoteID),
		StationID:   int32(ctx.StationID),
		Created:     ctx.Payload.Created,
		UserID:      p.UserID,
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

	fieldNoteQueryResult := &data.FieldNoteQueryResult{}
	if err := c.options.Database.GetContext(ctx, fieldNoteQueryResult, "SELECT fn.id AS ID, fn.created, fn.user_id, fn.note, c.key AS CategoryKey, u.username, m.url AS MediaURL, m.content_type AS MediaContentType FROM fieldkit.field_note AS fn JOIN fieldkit.user u ON (u.id = fn.user_id) JOIN fieldkit.field_note_category AS c ON (c.id = fn.category_id) LEFT OUTER JOIN fieldkit.field_note_media AS m ON (m.id = fn.media_id) WHERE fn.id = $1", fieldNote.ID); err != nil {
		return err
	}

	return ctx.OK(FieldNoteQueryResultType(fieldNoteQueryResult))
}

func (c *FieldNoteController) Get(ctx *app.GetFieldNoteContext) error {
	fieldNotes := []*data.FieldNoteQueryResult{}

	if err := c.options.Database.SelectContext(ctx, &fieldNotes, "SELECT fn.id AS ID, fn.created, fn.user_id, fn.note, c.key AS CategoryKey, u.username, m.url AS MediaURL, m.content_type AS MediaContentType FROM fieldkit.field_note AS fn JOIN fieldkit.user u ON (u.id = fn.user_id) JOIN fieldkit.field_note_category AS c ON (c.id = fn.category_id) LEFT OUTER JOIN fieldkit.field_note_media AS m ON (m.id = fn.media_id) WHERE fn.station_id = $1 ORDER BY fn.created DESC", ctx.StationID); err != nil {
		return err
	}

	return ctx.OK(FieldNotesType(fieldNotes))
}

func (c *FieldNoteController) Delete(ctx *app.DeleteFieldNoteContext) error {
	p, err := NewPermissions(ctx)
	if err != nil {
		return err
	}

	err = p.CanModifyStationByStationID(int32(ctx.StationID))
	if err != nil {
		return err
	}

	if _, err := c.options.Database.ExecContext(ctx, "DELETE FROM fieldkit.field_note WHERE id = $1", ctx.FieldNoteID); err != nil {
		return err
	}

	return ctx.OK()
}
