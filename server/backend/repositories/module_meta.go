package repositories

import (
	"fmt"
	"strings"

	"github.com/fieldkit/cloud/server/common/errors"
)

const (
	ManufacturerConservify = 0x01
	ConservifyWeather      = 0x01
	ConservifyWater        = 0x02
	ConservifyDistance     = 0x03
	ConservifyWaterPh      = 0x04
	ConservifyWaterEc      = 0x05
	ConservifyWaterDo      = 0x06
	ConservifyWaterTemp    = 0x07
	ConservifyWaterOrp     = 0x08
	ConservifyRandom       = 0xa0
	ConservifyDiagnostics  = 0xa1

	MaximumWindSpeed = 500    // km/hr world record is 371
	MaximumRain      = 1000.0 // mm
)

type HeaderFields struct {
	Manufacturer uint32
	Kind         uint32
}

type ModuleMetaRepository struct {
}

func NewModuleMetaRepository() *ModuleMetaRepository {
	return &ModuleMetaRepository{}
}

func (r *ModuleMetaRepository) FindByFullKey(fullKey string) (mm *SensorAndModuleMeta, err error) {
	all, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, err
	}

	for _, module := range all {
		for _, sensor := range module.Sensors {
			if sensor.FullKey == fullKey {
				mm = &SensorAndModuleMeta{
					Module: module,
					Sensor: sensor,
				}
				return
			}
		}
	}

	return nil, fmt.Errorf("unknown sensor: %s", fullKey)
}

func (r *ModuleMetaRepository) FindModuleMeta(m *HeaderFields) (mm *ModuleMeta, err error) {
	all, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, err
	}
	for _, module := range all {
		if module.Header.Manufacturer == m.Manufacturer && module.Header.Kind == m.Kind {
			return module, nil
		}
	}

	message := fmt.Sprintf("missing sensor meta (%v, %v)", m.Manufacturer, m.Kind)
	return nil, errors.Structured(message, "manufacturer", m.Manufacturer, "kind", m.Kind)
}

func (r *ModuleMetaRepository) FindSensorMeta(m *HeaderFields, sensor string) (mm *ModuleMeta, sm *SensorMeta, err error) {
	all, err := r.FindAllModulesMeta()
	if err != nil {
		return nil, nil, err
	}

	// Very old firmware keys. We should sanitize these earlier in the process.
	weNeedToCleanThisUp := strings.ReplaceAll(strings.ReplaceAll(sensor, " ", "_"), "-", "_")

	for _, module := range all {
		sameKind := module.Header.Kind == m.Kind
		if !sameKind {
			for _, k := range module.Header.AllKinds {
				if m.Kind == k {
					sameKind = true
					break
				}
			}
		}

		if module.Header.Manufacturer == m.Manufacturer && sameKind {
			for _, s := range module.Sensors {
				if s.Key == sensor || s.FirmwareKey == sensor {
					return module, s, nil
				}
				if s.Key == weNeedToCleanThisUp || s.FirmwareKey == weNeedToCleanThisUp {
					return module, s, nil
				}
			}
		}
	}

	message := fmt.Sprintf("missing sensor meta (%v, %v, %v)", m.Manufacturer, m.Kind, sensor)
	return nil, nil, errors.Structured(message, "manufacturer", m.Manufacturer, "kind", m.Kind, "sensor", sensor)
}

func (r *ModuleMetaRepository) FindAllModulesMeta() (mm []*ModuleMeta, err error) {
	mm = []*ModuleMeta{
		&ModuleMeta{
			Key: "fk.water.ph",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterPh,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "ph",
					FirmwareKey:   "ph",
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
			Key: "fk.water.ec",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterEc,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "ec",
					FirmwareKey:   "ec",
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
					FirmwareKey:   "tds",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{

							SensorRanges{
								Minimum: 0.0,
								Maximum: 0.0,
							},

					},
				},
				&SensorMeta{
					Key:           "salinity",
					FirmwareKey:   "salinity",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{
							SensorRanges{
								Minimum: 0.0,
								Maximum: 0.0,
							},
					},
				},
			},
		},
		&ModuleMeta{
			Key: "fk.water.dox",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterDo,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "dox",
					FirmwareKey:   "dox",
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
			Key: "fk.water.do",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterDo,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "do",
					FirmwareKey:   "do",
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
			Key: "fk.water.orp",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterOrp,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "orp",
					FirmwareKey:   "orp",
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
			Key: "fk.water.temp",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWaterTemp,
				AllKinds:     []uint32{ConservifyWater},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "temp",
					FirmwareKey:   "temp",
					UnitOfMeasure: "°C",
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
			Key: "fk.weather",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWeather,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "humidity",
					FirmwareKey:   "humidity",
					UnitOfMeasure: "%",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature1",
					FirmwareKey:   "temperature_1",
					UnitOfMeasure: "°C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "pressure",
					FirmwareKey:   "pressure",
					UnitOfMeasure: "kPa",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 1.0,
							Maximum: 200000.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature2",
					FirmwareKey:   "temperature_2",
					UnitOfMeasure: "°C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "rain",
					FirmwareKey:   "rain",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumRain,
						},
					},
				},
				&SensorMeta{
					Key:           "windSpeed",
					FirmwareKey:   "wind_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumWindSpeed,
						},
					},
				},
				&SensorMeta{
					Key:           "windDir",
					FirmwareKey:   "wind_dir",
					UnitOfMeasure: "°",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "windDirMv",
					FirmwareKey:   "wind_dir_mv",
					UnitOfMeasure: "mV",
					Ranges:        []SensorRanges{
						/*
							SensorRanges{
								Minimum: 0.0,
								Maximum: 0.0,
							},
						*/
					},
				},
				&SensorMeta{
					Key:           "windHrMaxSpeed",
					FirmwareKey:   "wind_hr_max_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumWindSpeed,
						},
					},
				},
				&SensorMeta{
					Key:           "windHrMaxDir",
					FirmwareKey:   "wind_hr_max_dir",
					UnitOfMeasure: "°",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind10mMaxSpeed",
					FirmwareKey:   "wind_10m_max_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumWindSpeed,
						},
					},
				},
				&SensorMeta{
					Key:           "wind10mMaxDir",
					FirmwareKey:   "wind_10m_max_dir",
					UnitOfMeasure: "°",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "wind2mAvgSpeed",
					FirmwareKey:   "wind_2m_avg_speed",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumWindSpeed,
						},
					},
				},
				&SensorMeta{
					Key:           "wind2mAvgDir",
					FirmwareKey:   "wind_2m_avg_dir",
					UnitOfMeasure: "°",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "rainThisHour",
					FirmwareKey:   "rain_this_hour",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumRain,
						},
					},
				},
				&SensorMeta{
					Key:           "rainPrevHour",
					FirmwareKey:   "rain_prev_hour",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumRain,
						},
					},
				},
			},
		},
		// This is from a very old version of the Weather firmware, we need to find a way to phase this out.
		&ModuleMeta{
			Key: "fk.weather",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWeather,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "humidity",
					FirmwareKey:   "humidity",
					UnitOfMeasure: "%",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 100.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature1",
					FirmwareKey:   "temperature_1",
					UnitOfMeasure: "°C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "pressure",
					FirmwareKey:   "pressure",
					UnitOfMeasure: "kPa",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 1.0,
							Maximum: 200000.0,
						},
					},
				},
				&SensorMeta{
					Key:           "temperature2",
					FirmwareKey:   "temperature_2",
					UnitOfMeasure: "°C",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -100.0,
							Maximum: 200.0,
						},
					},
				},
				&SensorMeta{
					Key:           "rain",
					FirmwareKey:   "rain",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumRain,
						},
					},
				},
				&SensorMeta{
					Key:           "wind",
					FirmwareKey:   "wind",
					UnitOfMeasure: "km/hr",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: MaximumWindSpeed,
						},
					},
				},
				&SensorMeta{
					Key:           "windDir",
					FirmwareKey:   "wind_dir",
					UnitOfMeasure: "°",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 360.0,
						},
					},
				},
				&SensorMeta{
					Key:           "windMv",
					FirmwareKey:   "wind_mv",
					UnitOfMeasure: "mV",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "fk.distance",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyDistance,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "distance",
					FirmwareKey:   "distance",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance0",
					FirmwareKey:   "distance_0",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance1",
					FirmwareKey:   "distance_1",
					UnitOfMeasure: "mm",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0.0,
							Maximum: 4999.0,
						},
					},
				},
				&SensorMeta{
					Key:           "distance2",
					FirmwareKey:   "distance_2",
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
					FirmwareKey:   "calibration",
					UnitOfMeasure: "mm",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "fk.diagnostics",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyDiagnostics,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Internal: true,
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "batteryCharge",
					FirmwareKey:   "battery_charge",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "batteryVoltage",
					FirmwareKey:   "battery_voltage",
					UnitOfMeasure: "V",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "batteryVbus",
					FirmwareKey:   "battery_vbus",
					UnitOfMeasure: "V",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "batteryVs",
					FirmwareKey:   "battery_vs",
					UnitOfMeasure: "mv",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "batteryMa",
					FirmwareKey:   "battery_ma",
					UnitOfMeasure: "ma",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "batteryPower",
					FirmwareKey:   "battery_power",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "solarVbus",
					FirmwareKey:   "solar_vbus",
					UnitOfMeasure: "V",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "solarVs",
					FirmwareKey:   "solar_vs",
					UnitOfMeasure: "mv",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "solarMa",
					FirmwareKey:   "solar_ma",
					UnitOfMeasure: "ma",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "solarPower",
					FirmwareKey:   "solar_power",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "freeMemory",
					FirmwareKey:   "free_memory",
					UnitOfMeasure: "bytes",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "uptime",
					FirmwareKey:   "uptime",
					UnitOfMeasure: "ms",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "temperature",
					FirmwareKey:   "temperature",
					UnitOfMeasure: "°C",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "fk.random",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyRandom,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Internal: true,
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "random0",
					FirmwareKey:   "random_0",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random1",
					FirmwareKey:   "random_1",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random2",
					FirmwareKey:   "random_2",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random3",
					FirmwareKey:   "random_3",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random4",
					FirmwareKey:   "random_4",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random5",
					FirmwareKey:   "random_5",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random6",
					FirmwareKey:   "random_6",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random7",
					FirmwareKey:   "random_7",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random8",
					FirmwareKey:   "random_8",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
				&SensorMeta{
					Key:           "random9",
					FirmwareKey:   "random_9",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{},
				},
			},
		},
		&ModuleMeta{
			Key: "fk.water",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyWater,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{},
		},
		&ModuleMeta{
			Key: "fk.testing",
			Header: ModuleHeader{
				Manufacturer: ManufacturerConservify,
				Kind:         ConservifyDiagnostics,
				AllKinds:     []uint32{},
				Version:      0x1,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "sin",
					FirmwareKey:   "sin",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: -1,
							Maximum: 1.0,
						},
					},
				},
				&SensorMeta{
					Key:           "saw.weekly",
					FirmwareKey:   "saw.weekly",
					UnitOfMeasure: "",
					Ranges: []SensorRanges{
						SensorRanges{
							Minimum: 0,
							Maximum: 1000,
						},
					},
				},
			},
		},
		// TODO This should move to the schema definition.
		&ModuleMeta{
			Key: "wh.floodnet",
			Header: ModuleHeader{
				Manufacturer: 0,
				Kind:         0,
				AllKinds:     []uint32{},
				Version:      0,
			},
			Sensors: []*SensorMeta{
				&SensorMeta{
					Key:           "humidity",
					FirmwareKey:   "humidity",
					UnitOfMeasure: "%",
					Ranges:        []SensorRanges{
					    SensorRanges{
                            Minimum: 0.0,
                            Maximum: 100.0,
                        },
					},
				},
				&SensorMeta{
					Key:           "pressure",
					FirmwareKey:   "pressure",
					UnitOfMeasure: "kPa",
					Ranges:        []SensorRanges{
					    SensorRanges{
                            Minimum: 1.0,
                            Maximum: 200000.0,
                        },
					},
				},
				&SensorMeta{
					Key:           "altitude",
					FirmwareKey:   "altitude",
					UnitOfMeasure: "m",
					Ranges:        []SensorRanges{
					SensorRanges{
                            Minimum: 3000,
                            Maximum: -500,
                        },
                    },
				},
				&SensorMeta{
					Key:           "battery",
					FirmwareKey:   "battery",
					UnitOfMeasure: "%",
					Ranges:        []SensorRanges{
					    SensorRanges{
                            Minimum: 0,
                            Maximum: 100,
                        },
                    },
				},
				&SensorMeta{
					Key:           "distance",
					FirmwareKey:   "distance",
					UnitOfMeasure: "mm",
					Ranges:        []SensorRanges{
					    SensorRanges{
                            Minimum: 0.0,
                            Maximum: 4999.0,
                        },
					},
				},
				&SensorMeta{
					Key:           "temperature",
					FirmwareKey:   "temperature",
					UnitOfMeasure: "°C",
					Ranges:        []SensorRanges{
                        SensorRanges{
                                Minimum: -100.0,
                                Maximum: 200.0,
                        },
                    },
				},
				&SensorMeta{
					Key:           "sdError",
					FirmwareKey:   "sdError",
					UnitOfMeasure: "",
					Ranges:        []SensorRanges{
					    SensorRanges{
                                Minimum: 0,
                                Maximum: 200.0,
                        },
					},
				},
			},
		},
	}

	for _, m := range mm {
		for _, s := range m.Sensors {
			s.FullKey = m.Key + "." + s.Key
		}
	}

	return
}
