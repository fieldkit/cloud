SELECT
s.id AS station_id, fn.id AS note_id, fnm.id AS media_id, fnm.content_type, fnm.url
FROM
fieldkit.station AS s LEFT JOIN fieldkit.field_note AS fn ON (fn.station_id = s.id)
LEFT JOIN fieldkit.field_note_media AS fnm ON (fn.media_id = fnm.id)
WHERE s.owner_id = 5;
