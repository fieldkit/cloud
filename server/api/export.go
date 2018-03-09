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

func (c *ExportController) ListByInput(ctx *app.ListByInputExportContext) error {
	ctx.ResponseData.Header().Set("Content-Type", "text/csv")
	ctx.ResponseData.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"source-%d.csv\"", ctx.InputID))
	ctx.ResponseData.WriteHeader(200)

	writer := bufio.NewWriter(ctx.ResponseData)

	var token *backend.PagingToken
	var keys []string

	for {
		page, newToken, err := c.options.Backend.ListDocumentsByInput(ctx, ctx.InputID, false, token)
		if err != nil {
			log.Printf("Error querying: %v", err)
			break
		}
		if len(page.Documents) == 0 {
			break
		}

		for _, doc := range page.Documents {
			fields, err := doc.GetParsedFields()
			if err != nil {
				log.Printf("Error parsing fields: %v", err)
			} else {
				if keys == nil {
					for key, _ := range fields {
						keys = append(keys, key)
					}
					sort.Strings(keys)
					row := []string{"ID", "Time", "Location"}
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
					fmt.Sprintf("%d", doc.ID),
					fmt.Sprintf("%v", doc.Timestamp),
					doc.Location.String(),
				}
				row = append(row, fieldValues...)
				writer.WriteString(strings.Join(row, ",") + "\n")
			}
		}

		token = newToken
	}
	return nil
}
