package webhook

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
)

type Extracted struct {
}

type Extractor interface {
	Extract(ctx context.Context, complete interface{}, source interface{}) (interface{}, error)
}

type FkStatusExtractor struct {
}

type FkLoRaStatus struct {
	BatteryBusVoltage float32 `json:"battery_bus_v"`
	BatteryMinV       float32 `json:"battery_min_v"`
	BatteryMaxV       float32 `json:"battery_max_v"`
	SolarMinV         float32 `json:"solar_min_v"`
	SolarMaxV         float32 `json:"solar_max_v"`
}

func (e *FkStatusExtractor) Extract(ctx context.Context, complete interface{}, source interface{}) (interface{}, error) {
	asString, ok := source.(string)
	if !ok {
		return nil, fmt.Errorf("fk:status expected string")
	}

	log := Logger(ctx).Sugar()

	raw, err := base64.StdEncoding.DecodeString(asString)
	if err != nil {
		return nil, fmt.Errorf("base64-error: %w", err)
	}

	status := FkLoRaStatus{}

	buf := bytes.NewReader(raw)
	if err := binary.Read(buf, binary.LittleEndian, &status.BatteryBusVoltage); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &status.BatteryMinV); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &status.BatteryMaxV); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &status.SolarMinV); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &status.SolarMaxV); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	log.Infow("extractor:fk:status", "status", status)

	return status, nil
}

type FkReadingsExtractor struct {
}

type FkLoRaReadings struct {
	Packet   byte      `json:"packet"`
	Age      uint64    `json:"age"`
	Block    uint64    `json:"block"`
	Readings []float32 `json:"readings"`
}

func (e *FkReadingsExtractor) Extract(ctx context.Context, complete interface{}, source interface{}) (interface{}, error) {
	asString, ok := source.(string)
	if !ok {
		return nil, fmt.Errorf("fk:readings expected string")
	}

	log := Logger(ctx).Sugar()

	raw, err := base64.StdEncoding.DecodeString(asString)
	if err != nil {
		return nil, fmt.Errorf("base64-error: %w", err)
	}

	var packet byte = 0
	buf := bytes.NewReader(raw)
	if err := binary.Read(buf, binary.LittleEndian, &packet); err != nil {
		return nil, fmt.Errorf("read-error(packet): %w", err)
	}

	readings := FkLoRaReadings{
		Packet:   packet,
		Age:      0,
		Block:    0,
		Readings: make([]float32, 0),
	}

	if packet == 0 {
		age, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, fmt.Errorf("read-error(age): %w", err)
		}

		block, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, fmt.Errorf("read-error(reading): %w", err)
		}

		readings.Age = age
		readings.Block = block
	}

	for buf.Len() > 0 {
		var value float32
		if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("read-error(reading): %w", err)
			} else {
				break
			}
		}

		readings.Readings = append(readings.Readings, value)
	}

	log.Infow("extractor:fk:readings", "readings", readings)

	return readings, nil
}

type FkLocationExtractor struct {
}

type FkLoRaLocation struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Altitude  float32 `json:"altitude"`
}

func (e *FkLocationExtractor) Extract(ctx context.Context, complete interface{}, source interface{}) (interface{}, error) {
	asString, ok := source.(string)
	if !ok {
		return nil, fmt.Errorf("fk:location expected string")
	}

	log := Logger(ctx).Sugar()

	raw, err := base64.StdEncoding.DecodeString(asString)
	if err != nil {
		return nil, fmt.Errorf("base64-error: %w", err)
	}

	location := FkLoRaLocation{}

	buf := bytes.NewReader(raw)
	if err := binary.Read(buf, binary.LittleEndian, &location.Longitude); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &location.Latitude); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &location.Altitude); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	log.Infow("extractor:fk:location", "location", location)

	return location, nil
}

func FindExtractor(ctxt context.Context, name string) (Extractor, error) {
	switch name {
	case "fk:status":
		return &FkStatusExtractor{}, nil
	case "fk:location":
		return &FkLocationExtractor{}, nil
	case "fk:readings":
		return &FkReadingsExtractor{}, nil
	}
	return nil, nil
}
