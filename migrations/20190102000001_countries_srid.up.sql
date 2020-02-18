/**
 * Set the SRID on the countries data.
 */
SELECT UpdateGeometrySRID('fieldkit', 'countries', 'geom', 4326);
