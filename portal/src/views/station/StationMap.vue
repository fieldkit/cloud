<template v-if="ready">
    <mapbox
        class="station-map"
        :access-token="mapbox.token"
        :map-options="{ style: mapbox.style, bounds: bounds, zoom: 10 }"
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
import Vue from "vue";
import { BoundingRectangle, DisplayStation, LngLat, MappedStations, ProjectModule } from "@/store";
import * as utils from "@/utilities";
import Config from "@/secrets";
import Mapbox from "mapbox-gl-vue";
import { ProtectedData } from "@/views/shared/StationsMap.vue";
import ValueMarker from "@/views/shared/ValueMarker.vue";
import mapboxgl from "mapbox-gl";
import * as d3 from "d3";

export default Vue.extend({
    name: "StationView",
    components: {
        Mapbox,
    },
    data(): {
        mapbox: { token: string; style: string };
        ready: boolean;
        thresholds: object;
    } {
        return {
            mapbox: Config.mapbox,
            ready: false,
            thresholds: {
                label: { "en-US": "Severity of Flooding" },
                levels: [
                    { label: { "en-US": "None" }, value: 8, color: "#00CCFF" },
                    { label: { "en-US": "Almost" }, value: 8.5, color: "#0099FF" },
                    { label: { "en-US": "Some" }, value: 9, color: "#0066FF" },
                    { label: { "en-US": "Flooding" }, value: 9.5, color: "#0033FF" },
                    { label: { "en-US": "Severe" }, value: 10, color: "#0000FF" },
                ],
            },
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
    },
    computed: {
        headerSubtitle() {
            let subtitle;
            if (this.station && this.station.deployedAt && this.$options.filters?.prettyDate) {
                if (this.station.deployedAt) {
                    subtitle = this.$tc("station.deployed") + " " + this.$options.filters.prettyDate(this.station.deployedAt);
                } else {
                    subtitle = this.$tc("station.readyToDeploy");
                }
            }
            return subtitle;
        },
        // Mapbox maps absolutely hate being mangled by Vue
        protectedData(): ProtectedData {
            return (this as unknown) as ProtectedData;
        },
        bounds(): LngLat[] | null {
            if (this.value) {
                //    return this.value.lngLat();
            }

            // return this.mapBounds.lngLat();
            return null;
        },
    },
    methods: {
        layoutChange() {
            this.$emit("layoutChange");
        },
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





            if (!map.getLayer("station-markers")) {
                console.log("map: updating", this.mapped);


               /* map.addLayer({
                    id: "regions",
                    type: "fill",
                    source: "stations",
                    paint: {
                        "fill-color": "#aaaaaa",
                        "fill-opacity": 1,
                    },
                    filter: ["==", "$type", "Polygon"],
                });
*/
               /* if (!this.mapped.isSingleType) {
                   /!* map.addLayer({
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
                        },*!/
                    });
                }*/

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

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/layout";

.station-map {
    height: 400px;
}
</style>
