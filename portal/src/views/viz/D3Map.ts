import _ from "lodash";
import Vue from "vue";
import Mapbox from "mapbox-gl-vue";
import { LngLatBounds } from "mapbox-gl";

import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, GeoZoom } from "./viz";
import { MapStore, Map } from "./MapStore";

import Config from "@/secrets";

const mapStore = new MapStore();

export const D3Map = Vue.extend({
    name: "D3Map",
    components: {
        Mapbox,
    },
    data() {
        return {
            mapboxToken: Config.MAPBOX_ACCESS_TOKEN,
            refreshed: false,
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
    destroyed() {
        mapStore.remove(this.viz.id);
    },
    updated() {
        this.viz.log("updated");
    },
    methods: {
        getMap(): Map {
            return mapStore.get(this.viz.id);
        },
        ready() {
            if (!this.data) {
                return false;
            }
            if (!this.getMap()) {
                return false;
            }

            const viewableData = this.data.data.filter((d) => d.location && d.location.length);

            if (viewableData.length == 0) {
                return false;
            }

            return true;
        },
        refresh() {
            if (!this.ready()) {
                return;
            }
            const map = this.getMap();
            const data = this.data.data;
            const located = data.filter((row) => row.location && row.location.length);
            const vizInfo = this.workspace.vizInfo(this.viz);
            const colors = vizInfo.colorScale;
            const enabled = true;

            this.viz.log("map", map);
            this.viz.log("map-refresh: data", located.length);

            const geojson = {
                type: "Feature",
                properties: {},
                geometry: {
                    type: "LineString",
                    coordinates: located.map((row) => row.location),
                },
            };

            // I wonder if there's a more d3 way to do this.
            this.removePreviousMapped(map);

            map.addSource("route", {
                type: "geojson",
                data: geojson,
            });
            map.addLayer({
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
            const geoDots = located.map((row) => {
                return {
                    type: "Feature",
                    properties: {
                        color: colors(row.value),
                    },
                    geometry: {
                        type: "Point",
                        coordinates: row.location,
                    },
                };
            });
            map.addSource("points", {
                type: "geojson",
                data: {
                    type: "FeatureCollection",
                    features: geoDots,
                },
            });
            map.addLayer({
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
            map.addLayer({
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

            if (this.viz.geo) {
                const bounds = this.viz.geo.bounds;

                this.viz.log("map-refresh: bounds(viz)", bounds);

                map.fitBounds(bounds, {
                    padding: 20,
                    duration: 0,
                });
            } else {
                const coordinates = geojson.geometry.coordinates;
                const bounds = coordinates.reduce((bounds, c) => bounds.extend(c, c), new LngLatBounds(coordinates[0], coordinates[0]));

                this.viz.log("map-refresh: bounds(data)", bounds.toArray());

                map.fitBounds(bounds.toArray(), {
                    padding: 20,
                    duration: 0,
                });
            }

            const z = map.getZoom();
            if (z > 19) {
                map.setZoom(19);
            }

            this.viz.log("map-refresh: done");

            this.refreshed = true;
        },
        removePreviousMapped(map) {
            if (map.getLayer("arrow-layer")) {
                map.removeLayer("arrow-layer");
            }
            if (map.getLayer("station-markers")) {
                map.removeLayer("station-markers");
            }
            if (map.getSource("points")) {
                map.removeSource("points");
            }
            if (map.getLayer("route")) {
                map.removeLayer("route");
            }
            if (map.getSource("route")) {
                map.removeSource("route");
            }
        },
        mapLoaded(map) {
            this.viz.log("map: ready");
            mapStore.set(this.viz.id, map).resize();
            this.refresh();
        },
        mapMoveEnd(...args) {
            if (this.ready() && this.refreshed) {
                this.viz.log("map-move-end");
                this.$emit(
                    "viz-geo-zoomed",
                    new GeoZoom(
                        this.getMap()
                            .getBounds()
                            .toArray()
                    )
                );
            } else {
                this.viz.log("map-move-end(ignored)");
            }
        },
        mapZoomEnd(...args) {
            // Our moveEnd handler above is enough.
            /*
            if (this.ready() && this.refreshed) {
                this.viz.log("map-zoom-end");
                this.$emit("viz-geo-zoomed", new GeoZoom(this.getMap().getBounds().toArray()));
            } else {
                this.viz.log("map-zoom-end(ignored)");
            }
			*/
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
                zoom: 10,
            }"
            :nav-control="{
                show: true,
                position: 'bottom-left',
            }"
            @map-load="mapLoaded"
			@map-moveend="mapMoveEnd"
			@map-zoomend="mapZoomEnd"
        />
	</div>`,
});
