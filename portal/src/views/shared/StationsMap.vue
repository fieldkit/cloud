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
            show: !isMobileView,
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

import _ from "lodash";
import Config from "@/secrets";
import { MappedStations, LngLat, BoundingRectangle, VisibleReadings, DecoratedReading } from "@/store";

import mapboxgl from "mapbox-gl";
import MapboxGeocoder from "@mapbox/mapbox-gl-geocoder";
import "@mapbox/mapbox-gl-geocoder/dist/mapbox-gl-geocoder.css";

import Vue, { PropType } from "vue";
import ValueMarker from "./ValueMarker.vue";
import Mapbox from "mapbox-gl-vue";

export interface ProtectedData {
    map: any;
    markers: { [index: number]: any };
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
        sensorMeta: Map<string, any> | null;
        hasGeocoder: boolean;
        isMobileView: boolean;
    } {
        return {
            mapbox: Config.mapbox,
            ready: false,
            sensorMeta: null,
            hasGeocoder: false,
            isMobileView: window.screen.availWidth <= 768,
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
        visibleReadings: {
            type: Number as PropType<VisibleReadings>,
            default: VisibleReadings.Current,
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
        visibleReadings(): void {
            console.log("map: visible-readings");
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

            if (!this.hasGeocoder) {
                map.addControl(
                    new MapboxGeocoder({
                        accessToken: this.mapbox.token,
                        mapboxgl: mapboxgl,
                        collapsed: true,
                        marker: false,
                    }),
                    "top-left"
                );
                this.hasGeocoder = true;
            }

            if (!map.getLayer("station-markers") && this.showStations) {
                const stationsSource = map.getSource("stations");
                if (!stationsSource) {
                    console.log("map: updating", this.mapped);

                    map.addSource("stations", {
                        type: "geojson",
                        data: {
                            type: "FeatureCollection",
                            features: this.mapped.features,
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
            } else {
                console.log("map: keeping", this.mapped);
            }

            if (this.bounds) {
                map.fitBounds(this.bounds, { duration: 0 });
            }

            // Regenerate custom map markers
            if (this.protectedData.markers) {
                for (const entry of this.protectedData.markers) {
                    const { marker } = entry;
                    marker.remove();
                }
                this.protectedData.markers = {};
            }

            const ValueMarkerCtor = Vue.extend(ValueMarker);
            const markers = [];
            const sortFactors = _.fromPairs(
                this.mapped.features.map((feature) => [feature.properties?.id, feature.station.getSortOrder(this.visibleReadings)])
            );
            const sorted = _.reverse(
                _.orderBy(
                    _.cloneDeep(this.mapped.features),
                    [
                        (feature) => sortFactors[feature.properties.id][0],
                        (feature) => sortFactors[feature.properties.id][1],
                        (feature) => sortFactors[feature.properties.id][2],
                    ],
                    ["asc", "desc", "asc"]
                )
            );
            for (const feature of sorted) {
                const readings = feature.station.inactive ? null : feature.station.getDecoratedReadings(this.visibleReadings);
                const instance = new ValueMarkerCtor({
                    propsData: {
                        ...(readings && readings.length > 0 && { color: readings[0].color }),
                        ...{ value: readings && readings.length > 0 ? readings[0].value : null },
                        ...{ id: feature.properties.id },
                    },
                });
                instance.$mount();
                instance.$on("marker-click", (evt) => {
                    this.$emit("show-summary", { id: evt.id });
                });

                const marker = new mapboxgl.Marker(instance.$el).setLngLat(feature.geometry.coordinates).addTo(map);
                markers.push({ marker: marker, instance: instance });
            }
            this.protectedData.markers = markers;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/global";

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

::v-deep .mapboxgl-ctrl-geocoder {
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
    border: solid 1px #f4f5f7;
    border-radius: 0;
    height: 40px;

    @include bp-down($sm) {
        width: 40px;
    }

    &.mapboxgl-ctrl-geocoder--collapsed {
        min-width: 40px;
    }

    &:not(.mapboxgl-ctrl-geocoder--collapsed) {
        @include bp-down($xs) {
            min-width: calc(100vw - 20px) !important;
        }
    }

    input {
        outline: none;
        height: 37px;
        padding-left: 38px;
        font-size: 16px;
    }
}

::v-deep .mapboxgl-ctrl-geocoder--icon-search {
    top: 9px;
    left: 8px;
}

::v-deep .mapboxgl-ctrl-geocoder--icon-close {
    margin-top: 4px;

    @include bp-down($sm) {
        margin-top: 3px;
    }
}
</style>
