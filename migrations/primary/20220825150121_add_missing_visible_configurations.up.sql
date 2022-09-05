INSERT INTO fieldkit.visible_configuration (configuration_id, station_id)		
SELECT
    sc.id AS configuration_id,
    (SELECT id FROM fieldkit.station WHERE device_id = p.device_id) AS station_id
FROM fieldkit.station_configuration AS sc JOIN fieldkit.provision AS p ON (sc.provision_id = p.id)
WHERE p.id IN (
    SELECT id FROM fieldkit.provision WHERE device_id IN (SELECT device_id FROM fieldkit.station WHERE id NOT IN (
        SELECT station_id FROM fieldkit.visible_configuration)
    )
);