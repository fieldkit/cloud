delete from fieldkit.document where input_id in (select id from fieldkit.input where name = 'weather-proxy');
delete from fieldkit.device_location where device_id in (select id from fieldkit.input where name = 'weather-proxy');
delete from fieldkit.device where input_id in (select id from fieldkit.input where name = 'weather-proxy');
delete from fieldkit.input where name = 'weather-proxy';
