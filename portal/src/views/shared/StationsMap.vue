<template v-if="mapped.valid && ready">
    <mapbox
        class="stations-map"
        :access-token="mapbox.token"
        :map-options="{
            style: mapbox.style,
            bounds: bounds,
            zoom: 10,
        }"
        :nav-control="{
            show: true,
            position: 'bottom-left',
        }"
        @map-init="onMapInitialized"
        @map-load="onMapLoaded"
        @zoomend="newBounds"
        @dragend="newBounds"
    />
</template>

<script lang="ts">
/* eslint-disable vue/no-unused-components */

import Vue from "vue";
import Mapbox from "mapbox-gl-vue";
import Config from "@/secrets";
import { MappedStations, LngLat, BoundingRectangle } from "@/store";
import ValueMarker from "./ValueMarker.vue";

import * as d3 from "d3";
import mapboxgl from "mapbox-gl";

interface ProtectedData {
    map: any;
}

export default Vue.extend({
    name: "StationsMap",
    components: {
        Mapbox,
        ValueMarker,
    },
    data(): {
        mapbox: { token: string; style: string };
        ready: boolean;
        sensorMeta: Map<string, any>;
    } {
        return {
            mapbox: Config.mapbox,
            ready: false,
            sensorMeta: null,
        };
    },
    props: {
        mapped: {
            type: MappedStations,
        },
        value: {
            type: BoundingRectangle,
        },
        mapBounds: {
            type: BoundingRectangle,
        },
        showStations: {
            type: Boolean,
            default: false,
        },
        layoutChanges: {
            type: Number,
            default: 0,
        },
    },
    computed: {
        // Mapbox maps absolutely hate being mangled by Vue
        protectedData(): ProtectedData {
            return (this as unknown) as ProtectedData;
        },
        bounds(): LngLat[] | null {
            if (this.value) {
                return this.value.lngLat();
            }

            return this.mapBounds ? this.mapBounds.lngLat() : this.mapped.boundsLngLat();
        },
    },
    watch: {
        layoutChanges(): void {
            console.log("map: layout changed");
            if (this.protectedData.map) {
                // TODO Not a fan of this.
                this.$nextTick(() => {
                    this.protectedData.map.resize();
                });
            }
        },
        mapped(): void {
            console.log("map: mapped changed", this.mapped);
            this.updateMap();
        },
        showStations(): void {
            this.updateMap();
        },
    },
    methods: {
        onMapInitialized(map: any): void {
            console.log("map: initialized");
            this.protectedData.map = map;
        },
        onMapLoaded(map: any): void {
            console.log("map: loaded");
            this.protectedData.map = map;

            if (!map.hasImage("dot")) {
                const compass = this.$loadAsset("Icon_Map_Dot.png");
                map.loadImage(compass, (error, image) => {
                    if (error) throw error;
                    if (!map.hasImage("dot")) {
                        map.addImage("dot", image);
                    }
                });
            }

            setTimeout(() => {
                map.resize();

                this.ready = true;
                this.updateMap();

                // Force model to update.
                this.newBounds();
            }, 100);
        },
        newBounds() {
            const map = this.protectedData.map;
            const bounds = map.getBounds();
            this.$emit("input", new BoundingRectangle([bounds._sw.lng, bounds._sw.lat], [bounds._ne.lng, bounds._ne.lat]));
        },
        updateMap(): void {
            if (!this.protectedData.map) {
                console.log("map: update-skip.1");
                return;
            }

            if (!this.mapped || !this.mapped.valid || !this.ready) {
                console.log("map: update-skip.2", this.mapped?.valid, this.ready);
                return;
            }

            const map = this.protectedData.map;

            // Marker color scale
            const appendColor = (features) => {
                return features.map((d) => {
                    if (d.properties.thresholds) {
                        const markerScale = d3
                            .scaleThreshold()
                            .domain(d.properties.thresholds.levels.map((d) => d.value))
                            .range(d.properties.thresholds.levels.map((d) => d.color));

                        d.properties.color = markerScale(d.properties.value);
                    } else {
                        //default color
                        d.properties.color = "#00CCFF";
                    }

                    return d;
                });
            };

            if (!map.getLayer("station-markers") && this.showStations) {
                console.log("map: updating", this.mapped);

                map.addSource("stations", {
                    type: "geojson",
                    data: {
                        type: "FeatureCollection",
                        features: appendColor(this.mapped.features),
                    },
                });

                map.addLayer({
                    id: "regions",
                    type: "fill",
                    source: "stations",
                    paint: {
                        "fill-color": "#aaaaaa",
                        "fill-opacity": 0.2,
                    },
                    filter: ["==", "$type", "Polygon"],
                });

                if (!this.mapped.isSingleType) {
                    map.addLayer({
                        id: "station-markers",
                        type: "symbol",
                        source: "stations",
                        filter: ["==", "$type", "Point"],
                        layout: {
                            "icon-image": "dot",
                            "text-field": "{title}",
                            "icon-ignore-placement": true,
                            "icon-allow-overlap": true,
                            "text-allow-overlap": true,
                            "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                            "text-offset": [0, 0.75],
                            "text-variable-anchor": ["top", "right", "bottom", "left"],
                        },
                    });
                }

                map.on("click", "station-markers", (e) => {
                    const id = e.features[0].properties.id;
                    console.log("map: click", id);
                    this.$emit("show-summary", { id: id });
                });
            } else {
                console.log("map: keeping", this.mapped);
            }

            if (this.bounds) {
                map.fitBounds(this.bounds, { duration: 0 });
            }

            //Generate custom map markers
            const valueMarker = Vue.extend(ValueMarker);

            for (const feature of this.mapped.features) {
                const instance = new valueMarker({
                    propsData: { color: feature.properties.color, value: feature.properties.value, id: feature.properties.id },
                });
                instance.$mount();
                instance.$on("marker-click", (evt) => {
                    this.$emit("show-summary", { id: evt.id });
                });

                new mapboxgl.Marker(instance.$el).setLngLat(feature.geometry.coordinates).addTo(map);
            }
        },
    },
});
</script>

<style scoped>
.map-view #map {
    height: 100%;
    position: relative;
    width: inherit;
}
.project-container #map {
    height: inherit;
    position: inherit;
    width: inherit;
}
.marker {
    height: 10px;
    width: 10px;
}
</style>
