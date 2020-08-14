<template>
    <div v-if="mapped.valid" class="ignored" style="height: 100%; width: 100%;">
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
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import Mapbox from "mapbox-gl-vue";
import Config from "@/secrets";

export default Vue.extend({
    name: "StationsMap",
    components: {
        Mapbox,
    },
    data: () => {
        return {
            mapboxToken: Config.MAPBOX_ACCESS_TOKEN,
        };
    },
    props: {
        mapped: {
            type: Object,
            required: true,
        },
        layoutChanges: {
            type: Number,
            default: 0,
        },
    },
    watch: {
        layoutChanges(this: any) {
            if (this.map) {
                this.map.resize();
            }
        },
    },
    methods: {
        onMapInitialized(this: any, map) {
            console.log("map: initialized");
            this.map = map;
        },
        onMapLoaded(this: any, map) {
            console.log("map: loaded");
            this.map = map;

            if (!this.map.hasImage("dot")) {
                const compass = this.$loadAsset("Icon_Map_Dot.png");
                this.map.loadImage(compass, (error, image) => {
                    if (error) throw error;
                    if (!this.map.hasImage("dot")) {
                        this.map.addImage("dot", image);
                    }
                });
            }

            this.map.resize();

            this.updateMap();
        },
        updateMap(this: any) {
            if (!this.map) {
                return;
            }

            if (!this.mapped || !this.mapped.valid) {
                return;
            }

            if (!this.map.getLayer("station-markers")) {
                console.log("map: updating", this.mapped);

                this.map.addSource("stations", {
                    type: "geojson",
                    data: {
                        type: "FeatureCollection",
                        features: this.mapped.features,
                    },
                });

                this.map.addLayer({
                    id: "regions",
                    type: "fill",
                    source: "stations",
                    paint: {
                        "fill-color": "#aaaaaa",
                        "fill-opacity": 0.2,
                    },
                    filter: ["==", "$type", "Polygon"],
                });

                this.map.addLayer({
                    id: "station-markers",
                    type: "symbol",
                    source: "stations",
                    filter: ["==", "$type", "Point"],
                    layout: {
                        "icon-image": "dot",
                        "text-field": "{title}",
                        "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                        "text-offset": [0, 0.75],
                        "text-anchor": "top",
                    },
                });

                this.map.on("click", "station-markers", (e) => {
                    const id = e.features[0].properties.id;
                    console.log("map: click", id);
                    this.$emit("show-summary", { id: id });
                });
            } else {
                console.log("map: ignored", this.mapped);
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
