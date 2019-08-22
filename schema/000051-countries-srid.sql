/**
 * Set the SRID on the countries data.
 */
SELECT UpdateGeometrySRID('countries', 'geom', 4326);
