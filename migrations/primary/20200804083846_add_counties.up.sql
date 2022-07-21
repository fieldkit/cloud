CREATE TABLE "fieldkit"."counties" (
	"gid" serial,
	"statefp" varchar(2),
	"countyfp" varchar(3),
	"countyns" varchar(8),
	"geoid" varchar(5),
	"name" varchar(100),
	"namelsad" varchar(100),
	"lsad" varchar(2),
	"classfp" varchar(2),
	"mtfcc" varchar(5),
	"csafp" varchar(3),
	"cbsafp" varchar(5),
	"metdivfp" varchar(5),
	"funcstat" varchar(1),
	"aland" float8,
	"awater" float8,
	"intptlat" varchar(11),
	"intptlon" varchar(12)
);
ALTER TABLE "fieldkit"."counties" ADD PRIMARY KEY (gid);
SELECT AddGeometryColumn('fieldkit', 'counties', 'geom', '0', 'MULTIPOLYGON', 2);
