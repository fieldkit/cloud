UPDATE fieldkit.sensor_meta SET viz = '[{"name":"D3TimeSeriesGraph","disabled":false,"minimumGap":600,"thresholds":{"label":{"en-US":"Flood Depth"},"levels":[{"label":{"en-US":"None or Minimal"},"mapKeyLabel":{"en-US":"None or Minimal"},"plainLabel":{"en-US":"None or Minimal"},"value":4,"color":"#00CCFF","hidden":true,"start":0},{"label":{"en-US":"Minor (4\")"},"mapKeyLabel":{"en-US":"Minor (4\"-12\")"},"plainLabel":{"en-US":"Minor"},"value":12,"color":"#fdbf4b","start":4},{"label":{"en-US":"Moderate (12\")"},"mapKeyLabel":{"en-US":"Moderate (12\"-24\")"},"plainLabel":{"en-US":"Moderate"},"value":24,"color":"#fe4d4c","start":12},{"label":{"en-US":"Major (24\")"},"mapKeyLabel":{"en-US":"Major (> 24\")"},"plainLabel":{"en-US":"Major"},"value":100,"color":"#d74dfe","start":24}]}},{"name":"D3Map","disabled":false,"thresholds":null}]' WHERE id = 70;
UPDATE fieldkit.sensor_meta SET viz = '[{"name":"D3TimeSeriesGraph","disabled":false,"minimumGap":600,"thresholds":null},{"name":"D3Map","disabled":false,"thresholds":null}]' WHERE id = 74;
