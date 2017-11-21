package ingestion

func AddLegacyRockBlockSchemas(sr *InMemorySchemas, rb DeviceId) {
	sr.DefineSchema(NewSchemaId(rb, "AT"), &JsonMessageSchema{
		HasTime: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Atlas Battery", Type: "float32"},
			JsonSchemaField{Name: "Orp", Type: "float32"},
			JsonSchemaField{Name: "DO", Type: "float32"},
			JsonSchemaField{Name: "pH", Type: "float32"},
			JsonSchemaField{Name: "EC 1", Type: "float32"},
			JsonSchemaField{Name: "EC 2", Type: "float32"},
			JsonSchemaField{Name: "EC 3", Type: "float32"},
			JsonSchemaField{Name: "EC 4", Type: "float32"},
			JsonSchemaField{Name: "Water Temp", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
			JsonSchemaField{Name: "Ignore", Type: "float32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "WE"), &JsonMessageSchema{
		HasTime: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},
			JsonSchemaField{Name: "Old Voltage", Type: "float32"},

			JsonSchemaField{Name: "Temp", Type: "float32"},
			JsonSchemaField{Name: "Humidity", Type: "float32"},
			JsonSchemaField{Name: "Pressure", Type: "float32"},
			JsonSchemaField{Name: "WindSpeed2", Type: "float32"},
			JsonSchemaField{Name: "WindDir2", Type: "float32"},
			JsonSchemaField{Name: "WindGust10", Type: "float32"},
			JsonSchemaField{Name: "WindGustDir10", Type: "float32"},
			JsonSchemaField{Name: "DailyRain", Type: "float32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "SO"), &JsonMessageSchema{
		HasTime: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Sonar Battery", Type: "float32"},
			JsonSchemaField{Name: "Depth 1", Type: "float32"},
			JsonSchemaField{Name: "Depth 2", Type: "float32"},
			JsonSchemaField{Name: "Depth 3", Type: "float32"},
			JsonSchemaField{Name: "Depth 4", Type: "float32"},
			JsonSchemaField{Name: "Depth 5", Type: "float32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "LO"), &JsonMessageSchema{
		HasTime:     true,
		HasLocation: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "Latitude", Type: "float32"},
			JsonSchemaField{Name: "Longitude", Type: "float32"},
			JsonSchemaField{Name: "Altitude", Type: "float32"},
			JsonSchemaField{Name: "Uptime", Type: "float32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "ST"), &JsonMessageSchema{
		HasTime: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "SD Available", Type: "boolean"},
			JsonSchemaField{Name: "GPS Fix", Type: "boolean"},
			JsonSchemaField{Name: "Battery Sleep Time", Type: "uint32"},
			JsonSchemaField{Name: "Tx Failures", Type: "uint32"},
			JsonSchemaField{Name: "Tx Skips", Type: "uint32"},
			JsonSchemaField{Name: "Weather Readings Received", Type: "uint32"},
			JsonSchemaField{Name: "Atlas Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Sonar Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Dead For", Type: "uint32"},

			JsonSchemaField{Name: "Avg Tx Time", Type: "float32"},

			JsonSchemaField{Name: "Idle Interval", Type: "uint32"},
			JsonSchemaField{Name: "Airwaves Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Start", Type: "uint32"},
			JsonSchemaField{Name: "Weather Ignore", Type: "uint32"},
			JsonSchemaField{Name: "Weather Off", Type: "uint32"},
			JsonSchemaField{Name: "Weather Reading", Type: "uint32"},

			JsonSchemaField{Name: "Tx A Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx A Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Interval", Type: "uint32"},

			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "ST"), &JsonMessageSchema{
		HasTime: true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Time", Type: "uint32"},
			JsonSchemaField{Name: "Station", Type: "string"},
			JsonSchemaField{Name: "Type", Type: "string"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},

			JsonSchemaField{Name: "SD Available", Type: "boolean"},
			JsonSchemaField{Name: "GPS Fix", Type: "boolean"},
			JsonSchemaField{Name: "Battery Sleep Time", Type: "uint32"},
			JsonSchemaField{Name: "Deep Sleep Time", Type: "uint32", Optional: true},
			JsonSchemaField{Name: "Tx Failures", Type: "uint32"},
			JsonSchemaField{Name: "Tx Skips", Type: "uint32"},
			JsonSchemaField{Name: "Weather Readings Received", Type: "uint32"},
			JsonSchemaField{Name: "Atlas Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Sonar Packets Rx", Type: "uint32"},
			JsonSchemaField{Name: "Dead For", Type: "uint32"},

			JsonSchemaField{Name: "Avg Tx Time", Type: "float32"},

			JsonSchemaField{Name: "Idle Interval", Type: "uint32"},
			JsonSchemaField{Name: "Airwaves Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Interval", Type: "uint32"},
			JsonSchemaField{Name: "Weather Start", Type: "uint32"},
			JsonSchemaField{Name: "Weather Ignore", Type: "uint32"},
			JsonSchemaField{Name: "Weather Off", Type: "uint32"},
			JsonSchemaField{Name: "Weather Reading", Type: "uint32"},

			JsonSchemaField{Name: "Tx A Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx A Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx B Interval", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Offset", Type: "uint32"},
			JsonSchemaField{Name: "Tx C Interval", Type: "uint32"},

			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})

	sr.DefineSchema(NewSchemaId(rb, "1"), &JsonMessageSchema{
		UseProviderTime: true,
		HasLocation:     true,
		Fields: []JsonSchemaField{
			JsonSchemaField{Name: "Latitude", Type: "float32"},
			JsonSchemaField{Name: "Longitude", Type: "float32"},
			JsonSchemaField{Name: "Altitude", Type: "float32"},
			JsonSchemaField{Name: "Temperature", Type: "float32"},
			JsonSchemaField{Name: "Humidity", Type: "float32"},
			JsonSchemaField{Name: "Battery Percentage", Type: "float32"},
			JsonSchemaField{Name: "Battery Voltage", Type: "float32"},
			JsonSchemaField{Name: "Uptime", Type: "uint32"},
		},
	})
}
