<template v-if="mapped.valid">
    <mapbox
        class="stations-map"
        :access-token="mapboxToken"
        :map-options="{
            style: 'mapbox://styles/mapbox/outdoors-v11',
            bounds: mapped.boundsLngLat(),
            zoom: 10,
        }"
        :nav-control="{
            show: true,
            position: 'bottom-left',
        }"
        @map-init="onMapInitialized"
        @map-load="onMapLoaded"
    />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import Mapbox from "mapbox-gl-vue";
import Config from "@/secrets";
import { promiseAfter } from "@/utilities";
import { MappedStations } from "@/store";

interface ProtectedData {
    map: any;
}

export default Vue.extend({
    name: "StationsMap",
    components: {
        Mapbox,
    },
    data(): { mapboxToken: string } {
        return {
            mapboxToken: Config.MAPBOX_ACCESS_TOKEN,
        };
    },
    props: {
        mapped: {
            type: MappedStations,
            required: true,
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
    },
    watch: {
        layoutChanges(): void {
            console.log("map: layout changed");
            if (this.protectedData.map) {
                // TODO Not a fan of this.
                promiseAfter(250, {}).then(() => {
                    this.protectedData.map.resize();
                });
            }
        },
        mapped(): void {
            console.log("map: mapped changed");
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

            map.resize();

            this.updateMap();
        },
        updateMap(): void {
            if (!this.protectedData.map) {
                return;
            }

            if (!this.mapped || !this.mapped.valid) {
                return;
            }

            const map = this.protectedData.map;

            if (!map.getLayer("station-markers")) {
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
                        "text-allow-overlap": false,
                        "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                        "text-offset": [0, 0.75],
                        "text-anchor": "top",
                    },
                });

                map.on("click", "station-markers", (e) => {
                    const id = e.features[0].properties.id;
                    console.log("map: click", id);
                    this.$emit("show-summary", { id: id });
                });
            } else {
                console.log("map: keeping features", this.mapped);
                map.fitBounds(this.mapped.boundsLngLat(), { duration: 0 });
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
</style>
