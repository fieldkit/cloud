UPDATE fieldkit.station_module SET name = 'fk.weather' WHERE name = 'modules.weather' OR name = 'weather';
UPDATE fieldkit.station_module SET name = 'fk.water.temp' WHERE name = 'modules.water.temp' OR name = 'water.temp';
UPDATE fieldkit.station_module SET name = 'fk.water.ph' WHERE name = 'modules.water.ph' OR name = 'water.ph';
UPDATE fieldkit.station_module SET name = 'fk.water.orp' WHERE name = 'modules.water.orp' OR name = 'water.orp';
UPDATE fieldkit.station_module SET name = 'fk.water.do' WHERE name = 'modules.water.do' OR name = 'water.do';
UPDATE fieldkit.station_module SET name = 'fk.water.do' WHERE name = 'modules.water.dox' OR name = 'water.dox';
UPDATE fieldkit.station_module SET name = 'fk.water.ec' WHERE name = 'modules.water.ec' OR name = 'water.ec';
UPDATE fieldkit.station_module SET name = 'fk.distance' WHERE name = 'modules.distance' OR name = 'distance';
UPDATE fieldkit.station_module SET name = 'fk.diagnostics' WHERE name = 'modules.diagnostics' OR name = 'diagnostics';
UPDATE fieldkit.station_module SET name = 'fk.random' WHERE name = 'modules.random' OR name = 'random' OR name = 'random-module-2';
UPDATE fieldkit.station_module SET name = 'fk.modules.unknown' WHERE name = 'modules.unknown';
UPDATE fieldkit.station_module SET name = 'fk.modules.water.unknown' WHERE name = 'modules.water.unknown';
UPDATE fieldkit.station_module SET name = 'fk.water' WHERE name = 'water';

CREATE FUNCTION fkLowerCaseFirst(word text)
RETURNS TEXT
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
  RETURN LOWER(LEFT(word, 1)) || RIGHT(word, -1);
END;
$$;

CREATE FUNCTION fkCamelCase(snake_case text)
RETURNS TEXT
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
  RETURN
    REPLACE(
      INITCAP(
        REPLACE(snake_case, '_', ' ')
      ),
      ' ', ''
    );
END;
$$;

UPDATE fieldkit.module_sensor SET name = REPLACE(name, '-', '_');
UPDATE fieldkit.module_sensor SET name = fkLowerCaseFirst(fkCamelCase(name)) WHERE name LIKE '%\_%';

DROP FUNCTION fkCamelCase;
DROP FUNCTION fkLowerCaseFirst;