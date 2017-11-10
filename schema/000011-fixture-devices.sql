INSERT INTO fieldkit.input (expedition_id, name, team_id, user_id, active) VALUES (
  (SELECT id FROM fieldkit.expedition WHERE slug = 'www-expedition' AND project_id = (SELECT id FROM fieldkit.project WHERE slug = 'www')),
  'FK Device #1',
  NULL,
  NULL,
  TRUE
) ON CONFLICT DO NOTHING;

INSERT INTO fieldkit.device (input_id, token) VALUES (
  (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
  '05a463d9-1c9c-4810-b3a7-8cc1c9452e89'
) ON CONFLICT DO NOTHING;

WITH the_schema AS
  (
    INSERT INTO fieldkit.schema (project_id, json_schema) VALUES (
      (SELECT id FROM fieldkit.project WHERE name = 'www'),
      '{}'
    ) RETURNING id
  )

  INSERT INTO fieldkit.stream (name, device_id, schema_id) VALUES (
    'FK Device #1 Status',
    (SELECT id FROM fieldkit.input WHERE name = 'FK Device #1'),
    (SELECT id FROM the_schema)
  );

