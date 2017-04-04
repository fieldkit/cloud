package api

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/goadesign/goa"

	"github.com/O-C-R/fieldkit/server/api/app"
	"github.com/O-C-R/fieldkit/server/backend"
	"github.com/O-C-R/fieldkit/server/data"
	"github.com/O-C-R/fieldkit/server/data/jsondocument"
)

const (
	VarintField  = "varint"
	UvarintField = "uvarint"
	Float32Field = "float32"
	Float64Field = "float64"
)

func FieldkitInputType(fieldkitInput *data.FieldkitInput) *app.FieldkitInput {
	fieldkitInputType := &app.FieldkitInput{
		ID:           int(fieldkitInput.ID),
		ExpeditionID: int(fieldkitInput.ExpeditionID),
		Name:         fieldkitInput.Name,
	}

	if fieldkitInput.TeamID != nil {
		teamID := int(*fieldkitInput.TeamID)
		fieldkitInputType.TeamID = &teamID
	}

	if fieldkitInput.UserID != nil {
		userID := int(*fieldkitInput.UserID)
		fieldkitInputType.UserID = &userID
	}

	return fieldkitInputType
}

func FieldkitInputsType(fieldkitInputs []*data.FieldkitInput) *app.FieldkitInputs {
	fieldkitInputsCollection := make([]*app.FieldkitInput, len(fieldkitInputs))
	for i, fieldkitInput := range fieldkitInputs {
		fieldkitInputsCollection[i] = FieldkitInputType(fieldkitInput)
	}

	return &app.FieldkitInputs{
		FieldkitInputs: fieldkitInputsCollection,
	}
}

func FieldkitBinaryType(fieldkitBinary *data.FieldkitBinary) *app.FieldkitBinary {
	return &app.FieldkitBinary{
		ID:        int(fieldkitBinary.ID),
		InputID:   int(fieldkitBinary.InputID),
		SchemaID:  int(fieldkitBinary.SchemaID),
		Fields:    fieldkitBinary.Fields,
		Mapper:    fieldkitBinary.Mapper.Pointers(),
		Longitude: fieldkitBinary.Longitude,
		Latitude:  fieldkitBinary.Latitude,
	}
}

type FieldkitControllerOptions struct {
	Backend *backend.Backend
}

// FieldkitController implements the twitter resource.
type FieldkitController struct {
	*goa.Controller
	options FieldkitControllerOptions
}

func NewFieldkitController(service *goa.Service, options FieldkitControllerOptions) *FieldkitController {
	return &FieldkitController{
		Controller: service.NewController("FieldkitController"),
		options:    options,
	}
}

func (c *FieldkitController) Add(ctx *app.AddFieldkitContext) error {
	fieldkitInput := &data.FieldkitInput{}
	fieldkitInput.ExpeditionID = int32(ctx.ExpeditionID)
	fieldkitInput.Name = ctx.Payload.Name
	if err := c.options.Backend.AddInput(ctx, &fieldkitInput.Input); err != nil {
		return err
	}

	if err := c.options.Backend.AddFieldkitInput(ctx, fieldkitInput); err != nil {
		return err
	}

	return ctx.OK(FieldkitInputType(fieldkitInput))
}

func (c *FieldkitController) SetBinary(ctx *app.SetBinaryFieldkitContext) error {
	fieldkitBinary := &data.FieldkitBinary{
		ID:       uint16(ctx.BinaryID),
		InputID:  int32(ctx.InputID),
		SchemaID: int32(ctx.Payload.SchemaID),
		Fields:   ctx.Payload.Fields,
		Mapper:   data.NewMapper(ctx.Payload.Mapper),
	}

	if err := c.options.Backend.SetFieldkitInputBinary(ctx, fieldkitBinary); err != nil {
		return err
	}

	return ctx.OK(FieldkitBinaryType(fieldkitBinary))
}

func (c *FieldkitController) addMessage(ctx context.Context, input *data.Input, fieldkitBinary *data.FieldkitBinary, timestamp int64, message interface{}) error {
	messageJSON, err := jsondocument.UnmarshalGo(message)
	if err != nil {
		return err
	}

	documentDataJSON, err := fieldkitBinary.Mapper.Map(messageJSON)
	if err != nil {
		return err
	}

	longitude, latitude := 0., 0.
	if fieldkitBinary.Longitude != nil {
		longitudeJSON, err := messageJSON.Document(*fieldkitBinary.Longitude)
		if err != nil {
			return err
		}

		longitude, err = longitudeJSON.Number()
		if err != nil {
			return err
		}
	}

	if fieldkitBinary.Latitude != nil {
		latitudeJSON, err := messageJSON.Document(*fieldkitBinary.Latitude)
		if err != nil {
			return err
		}

		latitude, err = latitudeJSON.Number()
		if err != nil {
			return err
		}
	}

	document := &data.Document{
		InputID:   input.ID,
		SchemaID:  fieldkitBinary.SchemaID,
		TeamID:    input.TeamID,
		UserID:    input.UserID,
		Location:  data.NewLocation(longitude, latitude),
		Timestamp: time.Unix(timestamp, 0),
	}

	if err := document.SetData(documentDataJSON.Interface()); err != nil {
		return err
	}

	return c.options.Backend.AddDocument(ctx, document)
}

func (c *FieldkitController) SendCsv(ctx *app.SendCsvFieldkitContext) error {

	// Get the input.
	input, err := c.options.Backend.Input(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	// Iterate through the CSV.
	reader := csv.NewReader(ctx.RequestData.Body)
	for {

		// Read a next line, breaking on EOF.
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(line) == 0 {
			continue
		}

		// Parse the first element of line.
		id, err := strconv.ParseUint(line[0], 10, 16)
		if err != nil {
			return err
		}

		// Get the fieldkit binary format.
		fieldkitBinary, err := c.options.Backend.FieldkitInputBinary(ctx, input.ID, uint16(id))
		if err != nil {
			return err
		}

		// Parse the first element of line.
		timestamp, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			return err
		}

		line = line[2:]
		message := make([]interface{}, len(fieldkitBinary.Fields))
		for i, field := range fieldkitBinary.Fields {
			switch field {
			case VarintField:
				data, err := strconv.ParseInt(line[i], 10, 64)
				if err != nil {
					return err
				}

				message[i] = float64(data)

			case UvarintField:
				data, err := strconv.ParseUint(line[i], 10, 64)
				if err != nil {
					return err
				}

				message[i] = float64(data)

			case Float32Field:
				data, err := strconv.ParseFloat(line[i], 32)
				if err != nil {
					return err
				}

				message[i] = data

			case Float64Field:
				data, err := strconv.ParseFloat(line[i], 64)
				if err != nil {
					return err
				}

				message[i] = data

			default:
				return fmt.Errorf("unknown field type, %s", field)
			}
		}

		if err := c.addMessage(ctx, input, fieldkitBinary, timestamp, message); err != nil {
			return err
		}
	}

	return ctx.NoContent()
}

func (c *FieldkitController) SendBinary(ctx *app.SendBinaryFieldkitContext) error {

	// Get the input.
	input, err := c.options.Backend.Input(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	reader := data.NewByteReader(ctx.RequestData.Body)
	for {
		id, err := binary.ReadUvarint(reader)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		timestamp, err := binary.ReadVarint(reader)
		if err != nil {
			return err
		}

		fieldkitBinary, err := c.options.Backend.FieldkitInputBinary(ctx, input.ID, uint16(id))
		if err != nil {
			return err
		}

		message := make([]interface{}, len(fieldkitBinary.Fields))
		for i, field := range fieldkitBinary.Fields {
			switch field {
			case VarintField:
				data, err := binary.ReadVarint(reader)
				if err != nil {
					return err
				}

				message[i] = float64(data)

			case UvarintField:
				data, err := binary.ReadUvarint(reader)
				if err != nil {
					return err
				}

				message[i] = float64(data)

			case Float32Field:
				data := make([]byte, 4)
				if _, err := io.ReadFull(reader, data); err != nil {
					return err
				}

				number := float64(math.Float32frombits(binary.LittleEndian.Uint32(data)))
				message[i] = float64(number)

			case Float64Field:
				data := make([]byte, 8)
				if _, err := io.ReadFull(reader, data); err != nil {
					return err
				}

				number := math.Float64frombits(binary.LittleEndian.Uint64(data))
				message[i] = float64(number)

			default:
				return fmt.Errorf("unknown field type, %s", field)
			}
		}

		if err := c.addMessage(ctx, input, fieldkitBinary, timestamp, message); err != nil {
			return err
		}
	}

	return ctx.NoContent()
}

func (c *FieldkitController) GetID(ctx *app.GetIDFieldkitContext) error {
	fieldkitInput, err := c.options.Backend.FieldkitInput(ctx, int32(ctx.InputID))
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputType(fieldkitInput))
}

func (c *FieldkitController) ListID(ctx *app.ListIDFieldkitContext) error {
	fieldkitInputs, err := c.options.Backend.ListFieldkitInputsByID(ctx, int32(ctx.ExpeditionID))
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputsType(fieldkitInputs))
}

func (c *FieldkitController) List(ctx *app.ListFieldkitContext) error {
	fieldkitInputs, err := c.options.Backend.ListFieldkitInputs(ctx, ctx.Project, ctx.Expedition)
	if err != nil {
		return err
	}

	return ctx.OK(FieldkitInputsType(fieldkitInputs))
}
