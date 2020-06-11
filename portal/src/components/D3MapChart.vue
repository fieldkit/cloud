<template>
    <div v-show="activeMode" class="cover-main-chart">
        <mapbox
            class="stations-map"
            :access-token="mapboxToken"
            :map-options="{
                container: chart.id + '-map',
                style: 'mapbox://styles/mapbox/light-v10',
                center: coordinates,
                zoom: 10,
            }"
            :nav-control="{
                show: true,
                position: 'bottom-left',
            }"
            @map-init="mapInitialized"
        />
        <div class="no-data-with-location" v-if="noData">No {{ chart.sensor.label }} data found with location information.</div>
    </div>
</template>

<script>
import Mapbox from "mapbox-gl-vue";
import { MAPBOX_ACCESS_TOKEN } from "../secrets";

export default {
    // keeping this name even though it doesn't currently use d3,
    // as it's still a child component of D3Chart
    name: "D3MapChart",
    props: ["chart", "layout"],
    data: () => {
        return {
            activeMode: false,
            coordinates: [-118, 34],
            mapboxToken: MAPBOX_ACCESS_TOKEN,
            noData: false,
        };
    },
    components: {
        Mapbox,
    },
    watch: {
        chart: function() {
            if (this.activeMode) {
                this.makeMap();
            }
        },
    },
    mounted() {
        let mapDiv = document.getElementById(this.chart.id + "-map");
        if (mapDiv) {
            mapDiv.style.width = this.layout.width + "px";
            mapDiv.style.height = this.layout.height + "px";
        }
    },
    methods: {
        setStatus(status) {
            this.activeMode = status;
        },
        dataChanged() {
            if (this.activeMode) {
                this.makeMap();
            }
        },
        mapInitialized(map) {
            this.map = map;
            this.map.scrollZoom.disable();
            const view = this;
            let imgData = require.context("../assets/", false, /\.png$/);
            imgData = imgData("./" + "map_arrow.png");
            this.map.loadImage(imgData, function(error, image) {
                if (error) throw error;
                if (!view.map.hasImage("arrow")) view.map.addImage("arrow", image);
            });
        },
        makeMap() {
            this.noData = false;
            this.hideLoading();
            setTimeout(() => {
                this.map.resize();
                this.displayData();
            }, 500);
        },
        displayData() {
            this.removePreviousMapped();
            this.chartData = this.chart.data.filter(d => {
                return d.location && d.location[0] != 0 && d.location[1] != 0;
            });
            if (this.chartData.length == 0) {
                this.noData = true;
                return;
            }
            // need long, lat, and data value for color:
            const coords = this.chartData.map(d => {
                return [d.location[0], d.location[1], d[this.chart.sensor.key]];
            });
            const geojson = {
                type: "Feature",
                properties: {},
                geometry: {
                    type: "LineString",
                    coordinates: coords,
                },
            };
            this.map.addSource("route", {
                type: "geojson",
                data: geojson,
            });
            // add line
            this.map.addLayer({
                id: "route",
                type: "line",
                source: "route",
                layout: {
                    "line-join": "round",
                    "line-cap": "round",
                },
                paint: {
                    "line-color": "#888",
                    "line-width": 3,
                },
            });
            // add dots
            const geoDots = geojson.geometry.coordinates.map(c => {
                return {
                    type: "Feature",
                    properties: {
                        color: this.chart.colors(c[2]),
                    },
                    geometry: {
                        type: "Point",
                        coordinates: c,
                    },
                };
            });
            this.map.addSource("points", {
                type: "geojson",
                data: {
                    type: "FeatureCollection",
                    features: geoDots,
                },
            });
            this.map.addLayer({
                id: "station-markers",
                type: "circle",
                source: "points",
                paint: {
                    // make circles larger as the user zooms from z12 to z22
                    "circle-radius": {
                        base: 1.75,
                        stops: [[12, 3], [22, 40]],
                    },
                    "circle-color": ["get", "color"],
                },
            });

            // add arrows
            this.map.addLayer({
                id: "arrow-layer",
                type: "symbol",
                source: "route",
                layout: {
                    "symbol-placement": "line",
                    "symbol-spacing": 100,
                    "icon-allow-overlap": true,
                    "icon-image": "arrow",
                    "icon-size": 0.75,
                    visibility: "visible",
                    // to reverse them:
                    // "icon-rotate": 180,
                },
            });

            let coordinates = geojson.geometry.coordinates;
            // Pass the first coordinates in the LineString to `lngLatBounds` &
            // wrap each coordinate pair in `extend` to include them in the bounds
            // result. A variation of this technique could be applied to zooming
            // to the bounds of multiple Points or Polygon geomteries - it just
            // requires wrapping all the coordinates with the extend method.
            let bounds = coordinates.reduce((bounds, coord) => {
                return bounds.extend(coord);
            }, new mapboxgl.LngLatBounds(coordinates[0], coordinates[0]));

            this.map.fitBounds(bounds, {
                padding: 20,
                duration: 0,
            });
            // bounds can be very close: zoom out if needed
            const z = this.map.getZoom();
            if (z > 19) {
                this.map.setZoom(19);
            }
        },
        removePreviousMapped() {
            if (this.map.getLayer("arrow-layer")) {
                this.map.removeLayer("arrow-layer");
            }
            if (this.map.getLayer("station-markers")) {
                this.map.removeLayer("station-markers");
            }
            if (this.map.getSource("points")) {
                this.map.removeSource("points");
            }
            if (this.map.getLayer("route")) {
                this.map.removeLayer("route");
            }
            if (this.map.getSource("route")) {
                this.map.removeSource("route");
            }
        },
        hideLoading() {
            if (document.getElementById("main-loading")) {
                document.getElementById("main-loading").style.display = "none";
            }
            document.getElementById(this.chart.id + "-loading").style.display = "none";
        },
    },
};
</script>

<style scoped>
.cover-main-chart {
    margin-top: -406px;
}
.no-data-with-location {
    float: left;
    position: relative;
    margin: -300px 0 0 60px;
    font-size: 20px;
    padding: 20px;
    background: rgba(255, 255, 255, 0.9);
    border-radius: 3px;
}
</style>
