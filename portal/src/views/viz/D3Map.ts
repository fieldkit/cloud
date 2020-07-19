import _ from "lodash";
import Vue from "vue";
import Mapbox from "mapbox-gl-vue";
import { LngLatBounds } from "mapbox-gl";

import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";

import Config from "@/secrets";

export const D3Map = Vue.extend({
    name: "D3Map",
    components: {
        Mapbox,
    },
    data() {
        return {
            coordinates: [-118, 34],
            mapboxToken: Config.MAPBOX_ACCESS_TOKEN,
            map: null,
        };
    },
    props: {
        viz: {
            type: Graph,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    computed: {
        data(): QueriedData {
            if (this.viz.graphing && !this.viz.graphing.empty) {
                return this.viz.graphing;
            }
            return null;
        },
    },
    watch: {
        data(newValue, oldValue) {
            this.viz.log("graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        this.viz.log("mounted");
        this.refresh();
    },
    updated() {
        this.viz.log("updated");
    },
    methods: {
        refresh() {
            if (!this.data) {
                return;
            }
            if (!this.map) {
                return;
            }
            const data = this.data.data;
            const located = data.filter((d) => d.location && d.location.length);
            const vizInfo = this.workspace.vizInfo(this.viz);
            const colors = vizInfo.colorScale;

            this.viz.log("data", located.length);

            if (located.length == 0) {
                return;
            }

            // I wonder if there's a more d3 way to do this.
            this.removePreviousMapped();

            const geojson = {
                type: "Feature",
                properties: {},
                geometry: {
                    type: "LineString",
                    coordinates: located.map((dr) => dr.location),
                },
            };
            this.map.addSource("route", {
                type: "geojson",
                data: geojson,
            });
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
            const geoDots = data.map((dr) => {
                return {
                    type: "Feature",
                    properties: {
                        color: colors(dr.value),
                    },
                    geometry: {
                        type: "Point",
                        coordinates: dr.location,
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
                        stops: [
                            [12, 3],
                            [22, 40],
                        ],
                    },
                    "circle-color": ["get", "color"],
                },
            });
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

            const coordinates = geojson.geometry.coordinates;
            const bounds = coordinates.reduce((bounds, c) => bounds.extend(c, c), new LngLatBounds(coordinates[0], coordinates[0]));

            this.viz.log("bounds", bounds.toArray());

            this.map.fitBounds(bounds.toArray(), {
                padding: 20,
                duration: 0,
            });

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
        mapLoaded(map) {
            this.viz.log("map: ready");
            this.map = map;
            this.map.resize();
            this.refresh();
        },
    },
    template: `
	<div class="viz map">
        <mapbox
            class="viz-map"
            :access-token="mapboxToken"
            :map-options="{
                container: this.viz.id + '-map',
                style: 'mapbox://styles/mapbox/light-v10',
                center: coordinates,
                zoom: 10,
            }"
            :nav-control="{
                show: true,
                position: 'bottom-left',
            }"
            @map-load="mapLoaded"
        />
        <div class="no-data-with-location" v-if="noData">No {{ chart.sensor.label }} data found with location information.</div>
	</div>`,
});
