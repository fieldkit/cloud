package bulk

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"time"
)

const (
	AggregatingBatchSize = 100
)

type CsvModel struct {
	DeviceIDColumn   uint   `json:"device_id_column"`
	DeviceNameColumn uint   `json:"device_name_column"`
	TimeColumn       uint   `json:"time_column"`
	TimeLayout       string `json:"time_layout"`
	BatteryColumn    int    `json:"battery_column"`
	SensorColumns    []uint `json:"sensor_columns"`
}

type CsvBulkSource struct {
	source  BulkSource
	path    string
	file    io.Reader
	reader  *csv.Reader
	model   *CsvModel
	columns []string
	verbose bool
}

func NewCsvBulkSource(path string, model *CsvModel) *CsvBulkSource {
	return &CsvBulkSource{
		path:  path,
		model: model,
	}
}

func (s *CsvBulkSource) NextBatch(ctx context.Context, batch *BulkMessageBatch) error {
	log := Logger(ctx).Sugar()

	if s.reader == nil {
		log.Infow("opening", "file", s.path)

		file, err := os.Open(s.path)
		if err != nil {
			return fmt.Errorf("opening %v (%v)", s.path, err)
		}

		s.reader = csv.NewReader(file)
	}

	if batch.Messages == nil {
		batch.Messages = make([]*BulkMessage, 0)
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

		deviceIDColumn := row[s.model.DeviceIDColumn]
		if deviceIDColumn == "" {
			if s.verbose {
				log.Infow("missing-device-id", "row", row)
			}
			continue
		}

		deviceNameColumn := row[s.model.DeviceNameColumn]
		if deviceNameColumn == "" {
			log.Infow("missing-device-name", "row", row)
			continue
		}

		timeColumn := row[s.model.TimeColumn]
		if timeColumn == "" {
			log.Infow("missing-time", "row", row)
			continue
		}

		stamp, err := time.Parse(s.model.TimeLayout, timeColumn)
		if err != nil {
			log.Infow("malformed-time", "column", timeColumn, "error", err)
			continue
		}

		message := &BulkMessage{
			Time:       stamp,
			DeviceID:   []byte(deviceIDColumn),
			DeviceName: deviceNameColumn,
			Readings:   make([]BulkReading, 0),
		}

		for _, sensorIndex := range s.model.SensorColumns {
			key := s.columns[sensorIndex]
			valueColumn := row[sensorIndex]
			value, ok := toFloat(valueColumn)
			if ok {
				message.Readings = append(message.Readings, BulkReading{
					Key:        key,
					Calibrated: value,
					Battery:    s.model.BatteryColumn >= 0 && int(sensorIndex) == s.model.BatteryColumn,
				})
			}
		}

		batch.Messages = append(batch.Messages, message)

		if len(batch.Messages) == AggregatingBatchSize {
			break
		}

		if false {
			log.Infow("csv", "message", message)
		}
	}

	if len(batch.Messages) == 0 {
		return io.EOF
	}

	return nil
}

func toFloat(x interface{}) (float64, bool) {
	switch x := x.(type) {
	case int:
		return float64(x), true
	case float64:
		return x, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	case *big.Int:
		f, err := strconv.ParseFloat(x.String(), 64)
		return f, err == nil
	default:
		return 0.0, false
	}
}
