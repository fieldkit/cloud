INSERT INTO fieldkit.input (expedition_id, name, team_id, user_id, active) VALUES (
  (SELECT id FROM fieldkit.expedition WHERE slug = 'www-expedition' AND project_id = (SELECT id FROM fieldkit.project WHERE slug = 'www')),
  'FK Device #1',
  NULL,
  NULL,
  TRUE
) ON CONFLICT DO NOTHING;

INSERT INTO fieldkit.device (input_id, token, key) VALUES (
  (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
  '05a463d9-1c9c-4810-b3a7-8cc1c9452e89',
  'ROCKBLOCK-11380'
) ON CONFLICT DO NOTHING;

WITH the_schema AS
  (
    INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (
      (SELECT id FROM fieldkit.project WHERE slug = 'www'), '{
  "UseProviderTime": false,
  "UseProviderLocation": false,
  "HasTime": true,
  "HasLocation": false,
  "Fields": [
    {
      "Name": "Time",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Station",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Type",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Battery Voltage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Battery Percentage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "SD Available",
      "Type": "boolean",
      "Optional": false
    },
    {
      "Name": "GPS Fix",
      "Type": "boolean",
      "Optional": false
    },
    {
      "Name": "Battery Sleep Time",
      "Type": "uint32",
      "Optional": false
    },
    {
    "Name": "Deep Sleep Time",
    "Type": "uint32",
    "Optional": true
    },
    {
      "Name": "Tx Failures",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx Skips",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Readings Received",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Atlas Packets Rx",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Sonar Packets Rx",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Dead For",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Avg Tx Time",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Idle Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Airwaves Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Start",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Ignore",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Off",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Reading",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx A Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx A Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx B Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx B Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx C Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx C Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Uptime",
      "Type": "uint32",
      "Optional": false
    }
  ]
}'
    ) RETURNING id
  )

  INSERT INTO fieldkit.device_schema (device_id, schema_id, key) VALUES (
    (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
    (SELECT id FROM the_schema),
    'ST'
  );

WITH the_schema AS
  (
    INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (
      (SELECT id FROM fieldkit.project WHERE slug = 'www'), '{
  "UseProviderTime": false,
  "UseProviderLocation": false,
  "HasTime": true,
  "HasLocation": false,
  "Fields": [
    {
      "Name": "Time",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Station",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Type",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Battery Voltage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Battery Percentage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "SD Available",
      "Type": "boolean",
      "Optional": false
    },
    {
      "Name": "GPS Fix",
      "Type": "boolean",
      "Optional": false
    },
    {
      "Name": "Battery Sleep Time",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx Failures",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx Skips",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Readings Received",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Atlas Packets Rx",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Sonar Packets Rx",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Dead For",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Avg Tx Time",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Idle Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Airwaves Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Start",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Ignore",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Off",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Weather Reading",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx A Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx A Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx B Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx B Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx C Offset",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Tx C Interval",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Uptime",
      "Type": "uint32",
      "Optional": false
    }
  ]
}'
    ) RETURNING id
  )

  INSERT INTO fieldkit.device_schema (device_id, schema_id, key) VALUES (
    (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
    (SELECT id FROM the_schema),
    'ST'
  );

WITH the_schema AS
  (
    INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (
      (SELECT id FROM fieldkit.project WHERE slug = 'www'), '{
  "UseProviderTime": false,
  "UseProviderLocation": false,
  "HasTime": true,
  "HasLocation": false,
  "Fields": [
    {
      "Name": "Time",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Station",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Type",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Battery Voltage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Battery Percentage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Old Voltage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Temp",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Humidity",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Pressure",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "WindSpeed2",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "WindDir2",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "WindGust10",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "WindGustDir10",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "DailyRain",
      "Type": "float32",
      "Optional": false
    }
  ]
}'
    ) RETURNING id
  )

  INSERT INTO fieldkit.device_schema (device_id, schema_id, key) VALUES (
    (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
    (SELECT id FROM the_schema),
    'WE'
  );

WITH the_schema AS
  (
    INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (
      (SELECT id FROM fieldkit.project WHERE slug = 'www'), '{
  "UseProviderTime": false,
  "UseProviderLocation": false,
  "HasTime": true,
  "HasLocation": true,
  "Fields": [
    {
      "Name": "Time",
      "Type": "uint32",
      "Optional": false
    },
    {
      "Name": "Station",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Type",
      "Type": "string",
      "Optional": false
    },
    {
      "Name": "Battery Voltage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Battery Percentage",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Latitude",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Longitude",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Altitude",
      "Type": "float32",
      "Optional": false
    },
    {
      "Name": "Uptime",
      "Type": "float32",
      "Optional": false
    }
  ]
}'
    ) RETURNING id
  )

  INSERT INTO fieldkit.device_schema (device_id, schema_id, key) VALUES (
    (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
    (SELECT id FROM the_schema),
    'LO'
  );
