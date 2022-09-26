package webhook

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type CsvMessageSource struct {
	path     string
	schemaID int32
	reader   *csv.Reader
	columns  []string
	verbose  bool
}

func NewCsvMessageSource(path string, schemaID int32, verbose bool) *CsvMessageSource {
	return &CsvMessageSource{
		path:     path,
		schemaID: schemaID,
		verbose:  verbose,
	}
}

func (s *CsvMessageSource) NextBatch(ctx context.Context, batch *MessageBatch) error {
	log := Logger(ctx).Sugar()

	if s.reader == nil {
		log.Infow("opening", "file", s.path)

		file, err := os.Open(s.path)
		if err != nil {
			return fmt.Errorf("opening %v (%w)", s.path, err)
		}

		s.reader = csv.NewReader(file)
	}

	if batch.Messages == nil {
		batch.Schemas = make(map[int32]*MessageSchemaRegistration)
		batch.Messages = make([]*WebHookMessage, 0)
	} else {
		batch.Messages = batch.Messages[:0]
	}

	for {
		row, err := s.reader.Read()
		if err == io.EOF {
			break
		}

		if s.verbose {
			log.Infow("row", "row", row)
		}

		if s.columns == nil {
			s.columns = row
			continue
		}

		jsonMap := make(map[string]interface{})
		for i, column := range s.columns {
			jsonMap[column] = row[i]
		}

		if s.verbose {
			log.Infow("json", "json", jsonMap)
		}

		body, err := json.Marshal(jsonMap)
		if err != nil {
			return fmt.Errorf("marshal json (%w)", err)
		}

		message := &WebHookMessage{
			ID:        int64(0),
			CreatedAt: time.Time{},
			SchemaID:  &s.schemaID,
			Headers:   nil,
			Body:      body,
		}

		if s.verbose {
			log.Infow("csv", "message", message)
		}

		batch.Messages = append(batch.Messages, message)

		if len(batch.Messages) == AggregatingBatchSize {
			break
		}
	}

	if len(batch.Messages) == 0 {
		return io.EOF
	}

	return nil
}
