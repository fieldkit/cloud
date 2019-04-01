SET CLIENT_ENCODING TO UTF8;
SET STANDARD_CONFORMING_STRINGS TO ON;
BEGIN;
CREATE TABLE "fieldkit"."countries" (gid serial,
"fips" varchar(2),
"iso2" varchar(2),
"iso3" varchar(3),
"un" int2,
"name" varchar(50),
"area" int4,
"pop2005" int8,
"region" int2,
"subregion" int2,
"lon" float8,
"lat" float8);
ALTER TABLE "fieldkit"."countries" ADD PRIMARY KEY (gid);
SELECT AddGeometryColumn('fieldkit','countries','geom','0','MULTIPOLYGON',2);
COMMIT;
ANALYZE "fieldkit"."countries";
