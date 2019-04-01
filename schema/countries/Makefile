LOCAL := $(shell pwd)
LOCAL_DATA := $(LOCAL)/data

default: image run

image:
	docker build --rm -t conservify/postgis-tools .

run: data/TM_WORLD_BORDERS_SIMPL-0.3.shp data/TM_WORLD_BORDERS-0.3.shp
	sudo whoami
	docker run --rm -v $(LOCAL_DATA):/data conservify/postgis-tools /bin/sh -c '/usr/local/bin/shp2pgsql -p /data/TM_WORLD_BORDERS_SIMPL-0.3.shp fieldkit.countries > /data/schema.sql'
	docker run --rm -v $(LOCAL_DATA):/data conservify/postgis-tools /bin/sh -c '/usr/local/bin/shp2pgsql -a /data/TM_WORLD_BORDERS_SIMPL-0.3.shp fieldkit.countries > /data/TM_WORLD_BORDERS_SIMPL-0.3.sql'
	docker run --rm -v $(LOCAL_DATA):/data conservify/postgis-tools /bin/sh -c '/usr/local/bin/shp2pgsql -a /data/TM_WORLD_BORDERS-0.3.shp fieldkit.countries > /data/TM_WORLD_BORDERS-0.3.sql'
	sudo chown -R $(USER). $(LOCAL_DATA)

data/TM_WORLD_BORDERS_SIMPL-0.3.shp:
	cd data && unzip -n TM_WORLD_BORDERS_SIMPL-0.3.zip

data/TM_WORLD_BORDERS-0.3.shp:
	cd data && unzip -n TM_WORLD_BORDERS-0.3.zip
