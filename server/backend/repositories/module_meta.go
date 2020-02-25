package repositories

import (
	"fmt"

	pb "github.com/fieldkit/data-protocol"
)

const (
	ManufacturerConservify = 0x01
	ConservifyWeather      = 0x01
	ConservifyWater        = 0x02
	ConservifyDistance     = 0x03
	ConservifyRandom       = 0xa0
	ConservifyDiagnostics  = 0xa1
)

type ModuleMetaRepository struct {
}

func NewModuleMetaRepository() *ModuleMetaRepository {
	return &ModuleMetaRepository{}
}

func (r *ModuleMetaRepository) FindModuleMeta(m *pb.ModuleHeader) (mm *ModuleMeta, err error) {
	all, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, err
	}
	for _, module := range all {
		if module.Header.Manufacturer == m.Manufacturer && module.Header.Kind == m.Kind {
			return module, nil
		}
	}
	return nil, fmt.Errorf("missing module meta")
}

func (r *ModuleMetaRepository) FindSensor(m *pb.ModuleHeader, sensor string) (mm *SensorMeta, err error) {
	all, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, err
	}
	for _, module := range all {
		if module.Header.Manufacturer == m.Manufacturer && module.Header.Kind == m.Kind {
			for _, s := range module.Sensors {
				if s.Key == sensor {
					return s, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("missing sensor meta")
}

func (r *ModuleMetaRepository) FindAllModulesMeta() (mm []*ModuleMeta, err error) {
	mm = []*ModuleMeta{
		&ModuleMeta{
			Key: "modules.water.ph",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "ph",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 14.0,
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.water.ec",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "ec",
					UnitOfMeasure: "μS/cm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 500000.0,
						},
					},
				},
				&SensorMeta{
					Key:           "tds",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "salinity",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.water.do",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "do",
					UnitOfMeasure: "mg/L",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0,
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.water.orp",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "orp",
					UnitOfMeasure: "mV",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -1019.9,
							Maximum: 1019.9,
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.water.temp",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "temp",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -126.0,
							Maximum: 1254.0,
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.weather",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWeather,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "humidity",
					UnitOfMeasure: "%",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature_1",
					UnitOfMeasure: "C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "pressure",
					UnitOfMeasure: "kPa",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 200000.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature_2",
					UnitOfMeasure: "C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "rain",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0, // TODO Arbitrary
						},
					},
				},
				&SensorMeta{
					Key:           "wind_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_dir",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_dir_mv",
					UnitOfMeasure: "mV",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_hr_max_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_hr_max_dir",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_10m_max_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_10m_max_dir",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_2m_avg_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 0.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind_2m_avg_dir",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "rain_this_hour",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0, // TODO Arbitrary
						},
					},
				},
				&SensorMeta{
					Key:           "rain_prev_hour",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0, // TODO Arbitrary
						},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.distance",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyDistance,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "distance",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance_0",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance_1",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance_2",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "calibration",
					UnitOfMeasure: "mm",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.diagnostics",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyDiagnostics,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "battery_charge",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "battery_voltage",
					UnitOfMeasure: "V",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "battery_vbus",
					UnitOfMeasure: "V",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "battery_vs",
					UnitOfMeasure: "mv",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "battery_ma",
					UnitOfMeasure: "ma",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "battery_power",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "free_memory",
					UnitOfMeasure: "bytes",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "uptime",
					UnitOfMeasure: "ms",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "temperature",
					UnitOfMeasure: "C",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "modules.random",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyRandom,
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "random_0",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_1",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_2",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_3",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_4",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_5",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_6",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_7",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_8",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random_9",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
			},
		},
	}

	return
}
