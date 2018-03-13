package api

import (
	"bufio"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/goadesign/goa"

	"github.com/fieldkit/cloud/server/api/app"
	"github.com/fieldkit/cloud/server/backend"
)

type ExportControllerOptions struct {
	Backend *backend.Backend
}

type ExportController struct {
	*goa.Controller
	options ExportControllerOptions
}

func NewExportController(service *goa.Service, options ExportControllerOptions) *ExportController {
	return &ExportController{
		Controller: service.NewController("ExportController"),
		options:    options,
	}
}

func BoolTo0or1(flag bool) int {
	if flag {
		return 1
	}
	return 0
}

func (c *ExportController) ListBySource(ctx *app.ListBySourceExportContext) error {
	source, err := c.options.Backend.GetSourceByID(ctx, int32(ctx.SourceID))
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("source-%d-%s.csv", source.ID, source.Name)

	ctx.ResponseData.Header().Set("Content-Type", "text/csv")
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	ctx.ResponseData.WriteHeader(200)

	writer := bufio.NewWriter(ctx.ResponseData)

	var token *backend.PagingToken
	var keys []string

	for {
		page, newToken, err := c.options.Backend.ListRecordsBySource(ctx, ctx.SourceID, false, true, token)
		if err != nil {
			return err
		}
		if len(page.Records) == 0 {
			break
		}

		for _, record := range page.Records {
			fields, err := record.GetParsedFields()
			if err != nil {
				log.Printf("Error parsing fields: %v", err)
			} else {
				if keys == nil {
					for key, _ := range fields {
						keys = append(keys, key)
					}
					sort.Strings(keys)
					row := []string{"id", "time", "longitude", "latitude", "visible", "fixed", "outlier", "manually_excluded"}
					row = append(row, keys...)
					writer.WriteString(strings.Join(row, ",") + "\n")
				}

				fieldValues := make([]string, 0)
				for _, key := range keys {
					value := fields[key]
					if value != nil {
						fieldValues = append(fieldValues, fmt.Sprintf("%v", value))
					} else {
						fieldValues = append(fieldValues, "")
					}
				}
				row := []string{
					fmt.Sprintf("%d", record.ID),
					fmt.Sprintf("%v", record.Timestamp),
					fmt.Sprintf("%f", record.Location.Coordinates()[0]),
					fmt.Sprintf("%f", record.Location.Coordinates()[1]),
					fmt.Sprintf("%d", BoolTo0or1(record.Visible)),
					fmt.Sprintf("%d", BoolTo0or1(record.Fixed)),
					fmt.Sprintf("%d", BoolTo0or1(record.Outlier)),
					fmt.Sprintf("%d", BoolTo0or1(record.ManuallyExcluded)),
				}
				row = append(row, fieldValues...)
				writer.WriteString(strings.Join(row, ",") + "\n")
			}
		}

		token = newToken
	}
	return nil
}
