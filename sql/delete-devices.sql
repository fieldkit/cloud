delete from fieldkit.record where source_id in (select id from fieldkit.source where name = 'weather-proxy');
delete from fieldkit.device_location where device_id in (select id from fieldkit.source where name = 'weather-proxy');
delete from fieldkit.device where source_id in (select id from fieldkit.source where name = 'weather-proxy');
delete from fieldkit.source where name = 'weather-proxy';
