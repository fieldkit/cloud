ALTER TABLE fieldkit.sensor_meta ADD COLUMN bucket_widths INTEGER[];

UPDATE fieldkit.sensor_meta SET bucket_widths = ARRAY[ 6*3600 + 12*60 + 30 ] WHERE full_key = 'wh.floodnet.tideFeet';