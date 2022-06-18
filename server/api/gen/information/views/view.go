// Code generated by goa v3.2.4, DO NOT EDIT.
//
// information views
//
// Command:
// $ goa gen github.com/fieldkit/cloud/server/api/design

package views

import (
	goa "goa.design/goa/v3/pkg"
)

// DeviceLayoutResponse is the viewed result type that is projected based on a
// view.
type DeviceLayoutResponse struct {
	// Type to project
	Projected *DeviceLayoutResponseView
	// View to render
	View string
}

// DeviceLayoutResponseView is a type that runs validations on a projected type.
type DeviceLayoutResponseView struct {
	Configurations []*StationConfigurationView
	Sensors        map[string][]*StationSensorView
}

// StationConfigurationView is a type that runs validations on a projected type.
type StationConfigurationView struct {
	ID           *int64
	Time         *int64
	ProvisionID  *int64
	MetaRecordID *int64
	SourceID     *int32
	Modules      []*StationModuleView
}

// StationModuleView is a type that runs validations on a projected type.
type StationModuleView struct {
	ID               *int64
	HardwareID       *string
	HardwareIDBase64 *string
	MetaRecordID     *int64
	Name             *string
	Position         *int32
	Flags            *int32
	Internal         *bool
	FullKey          *string
	Sensors          []*StationSensorView
	Meta             map[string]interface{}
}

// StationSensorView is a type that runs validations on a projected type.
type StationSensorView struct {
	Name          *string
	UnitOfMeasure *string
	Reading       *SensorReadingView
	Key           *string
	FullKey       *string
	Ranges        []*SensorRangeView
	Meta          map[string]interface{}
}

// SensorReadingView is a type that runs validations on a projected type.
type SensorReadingView struct {
	Last *float32
	Time *int64
}

// SensorRangeView is a type that runs validations on a projected type.
type SensorRangeView struct {
	Minimum *float32
	Maximum *float32
}

var (
	// DeviceLayoutResponseMap is a map of attribute names in result type
	// DeviceLayoutResponse indexed by view name.
	DeviceLayoutResponseMap = map[string][]string{
		"default": []string{
			"configurations",
			"sensors",
		},
	}
)

// ValidateDeviceLayoutResponse runs the validations defined on the viewed
// result type DeviceLayoutResponse.
func ValidateDeviceLayoutResponse(result *DeviceLayoutResponse) (err error) {
	switch result.View {
	case "default", "":
		err = ValidateDeviceLayoutResponseView(result.Projected)
	default:
		err = goa.InvalidEnumValueError("view", result.View, []interface{}{"default"})
	}
	return
}

// ValidateDeviceLayoutResponseView runs the validations defined on
// DeviceLayoutResponseView using the "default" view.
func ValidateDeviceLayoutResponseView(result *DeviceLayoutResponseView) (err error) {
	if result.Configurations == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("configurations", "result"))
	}
	if result.Sensors == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("sensors", "result"))
	}
	for _, e := range result.Configurations {
		if e != nil {
			if err2 := ValidateStationConfigurationView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationConfigurationView runs the validations defined on
// StationConfigurationView.
func ValidateStationConfigurationView(result *StationConfigurationView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.ProvisionID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("provisionId", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	if result.Modules == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("modules", "result"))
	}
	for _, e := range result.Modules {
		if e != nil {
			if err2 := ValidateStationModuleView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationModuleView runs the validations defined on StationModuleView.
func ValidateStationModuleView(result *StationModuleView) (err error) {
	if result.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "result"))
	}
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.Position == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("position", "result"))
	}
	if result.Flags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("flags", "result"))
	}
	if result.Internal == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("internal", "result"))
	}
	if result.FullKey == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fullKey", "result"))
	}
	if result.Sensors == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("sensors", "result"))
	}
	for _, e := range result.Sensors {
		if e != nil {
			if err2 := ValidateStationSensorView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateStationSensorView runs the validations defined on StationSensorView.
func ValidateStationSensorView(result *StationSensorView) (err error) {
	if result.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "result"))
	}
	if result.UnitOfMeasure == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("unitOfMeasure", "result"))
	}
	if result.Key == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("key", "result"))
	}
	if result.FullKey == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fullKey", "result"))
	}
	if result.Ranges == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("ranges", "result"))
	}
	if result.Reading != nil {
		if err2 := ValidateSensorReadingView(result.Reading); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range result.Ranges {
		if e != nil {
			if err2 := ValidateSensorRangeView(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateSensorReadingView runs the validations defined on SensorReadingView.
func ValidateSensorReadingView(result *SensorReadingView) (err error) {
	if result.Last == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("last", "result"))
	}
	if result.Time == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("time", "result"))
	}
	return
}

// ValidateSensorRangeView runs the validations defined on SensorRangeView.
func ValidateSensorRangeView(result *SensorRangeView) (err error) {
	if result.Minimum == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minimum", "result"))
	}
	if result.Maximum == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("maximum", "result"))
	}
	return
}
