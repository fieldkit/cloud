<template>
    <mapbox
        class="stations-map"
        :access-token="mapboxToken"
        :map-options="{
            style: 'mapbox://styles/mapbox/outdoors-v11',
            center: coordinates,
            zoom: 10,
        }"
        :nav-control="{
            show: true,
            position: 'bottom-left',
        }"
        @map-init="mapInitialized"
    />
</template>

<script>
import Mapbox from "mapbox-gl-vue";
import { MAPBOX_ACCESS_TOKEN } from "../secrets";

export default {
    name: "StationsMap",
    components: {
        Mapbox,
    },
    data: () => {
        return {
            coordinates: [-118, 34],
            mapboxToken: MAPBOX_ACCESS_TOKEN,
        };
    },
    props: ["mapSize", "stations"],
    watch: {
        stations() {
            if (this.stations) {
                this.updateMap();
            }
        },
    },
    mounted() {
        let mapDiv = document.getElementById("map");
        mapDiv.style.width = this.mapSize.width;
        mapDiv.style.height = this.mapSize.height;
        mapDiv.style.position = this.mapSize.position;
        this.map.resize();
    },
    methods: {
        mapInitialized(map) {
            this.map = map;
            const view = this;
            let imgData = require.context("../assets/", false, /\.png$/);
            imgData = imgData("./" + "Icon_Map_Dot.png");
            this.map.loadImage(imgData, function(error, image) {
                if (error) throw error;
                if (!view.map.hasImage("dot")) view.map.addImage("dot", image);
            });
            this.$emit("mapReady", map);
        },
        updateMap() {
            let longMax = -180;
            let longMin = 180;
            let latMin = 90;
            let latMax = -90;
            let stationFeatures = [];
            let mappable = this.stations.filter(s => {
                return (
                    s.status_json.latitude && s.status_json.longitude && s.status_json.latitude != 1000 && s.status_json.longitude != 1000
                );
            });
            mappable.forEach(s => {
                let coordinates = [s.status_json.latitude, s.status_json.longitude];
                if (mappable.length == 1) {
                    this.map.setCenter({
                        lat: coordinates[0],
                        lng: coordinates[1],
                    });
                } else {
                    if (s.status_json.latitude > latMax) {
                        latMax = s.status_json.latitude;
                    }
                    if (s.status_json.latitude < latMin) {
                        latMin = s.status_json.latitude;
                    }
                    if (s.status_json.longitude > longMax) {
                        longMax = s.status_json.longitude;
                    }
                    if (s.status_json.longitude < longMin) {
                        longMin = s.status_json.longitude;
                    }
                }
                stationFeatures.push({
                    type: "Feature",
                    geometry: {
                        type: "Point",
                        coordinates: [coordinates[1], coordinates[0]],
                    },
                    properties: {
                        title: s.name,
                        icon: "marker",
                    },
                });
            });
            // sort by latitude
            stationFeatures.sort(function(a, b) {
                var latA = a.geometry.coordinates[0];
                var latB = b.geometry.coordinates[0];
                if (latA < latB) {
                    return -1;
                }
                if (latA > latB) {
                    return 1;
                }
                // lats must be equal
                return 0;
            });
            // could sort by and compare longitude, too

            // TODO: handle this better - merge into aggregate markers
            // currently if stations overlap exactly, only one shows
            stationFeatures.forEach((f, i) => {
                if (i > 0) {
                    if (f.geometry.coordinates[0] == stationFeatures[i - 1].geometry.coordinates[0]) {
                        // TODO: what is an acceptable amount to offset labels by?
                        // less than this, and still only one shows
                        f.geometry.coordinates[0] += 0.00003;
                    }
                }
            });

            this.map.addLayer({
                id: "station-marker",
                type: "symbol",
                source: {
                    type: "geojson",
                    data: {
                        type: "FeatureCollection",
                        features: stationFeatures,
                    },
                },
                layout: {
                    "icon-image": "dot",
                    "text-field": "{title}",
                    "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                    "text-offset": [0, 0.75],
                    "text-anchor": "top",
                },
                // paint: {
                //     "text-color": "#0000ff"
                // }
            });

            if (mappable.length > 1) {
                let bounds = [
                    {
                        lat: latMin,
                        lng: longMin,
                    },
                    {
                        lat: latMax,
                        lng: longMax,
                    },
                ];
                this.map.fitBounds(bounds, { duration: 0 });
                // fitBounds is cropped too close: zoom out
                const z = this.map.getZoom();
                this.map.setZoom(z - 1.5);
            } else if (mappable.length == 0) {
                this.map.setZoom(6);
            }

            const view = this;
            this.map.on("click", "station-marker", function(e) {
                const name = e.features[0].properties.title;
                const station = view.stations.find(s => {
                    return s.name == name;
                });
                // view.showSummary(station);
                view.$emit("showSummary", station);
            });
        },
    },
};
</script>

<style scoped></style>
