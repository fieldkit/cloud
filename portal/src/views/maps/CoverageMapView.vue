<template>
    <StandardLayout @sidebar-toggle="layoutChanges++">
        <div class="container-map">
            <template v-if="ready">
                <mapbox
                    class="stations-map"
                    :access-token="mapbox.token"
                    :map-options="{
                        style: mapbox.style,
                        bounds: bounds,
                        zoom: 10,
                    }"
                    :nav-control="{
                        show: false,
                        position: 'bottom-left',
                    }"
                    @map-init="onMapInitialized"
                    @map-load="onMapLoaded"
                    @zoomend="newBounds"
                    @dragend="newBounds"
                />
            </template>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";
import { LngLat, BoundingRectangle } from "@/store";

import Vue, { PropType } from "vue";
import StandardLayout from "../StandardLayout.vue";

import Config from "@/secrets";
import mapboxgl from "mapbox-gl";
import Mapbox from "mapbox-gl-vue";

export interface ProtectedData {
    map: any;
}

export default Vue.extend({
    name: "CoverageMapView",
    components: {
        StandardLayout,
        Mapbox,
    },
    props: {
        bounds: {
            type: (Array as unknown) as PropType<[[number, number], [number, number]]>,
            required: false,
        },
    },
    data(): {
        mapbox: { token: string; style: string };
        ready: boolean;
    } {
        return {
            mapbox: Config.mapbox,
            ready: false,
        };
    },
    computed: {
        // Mapbox maps absolutely hate being mangled by Vue
        protectedData(): ProtectedData {
            return (this as unknown) as ProtectedData;
        },
    },
    async beforeMount(): Promise<any> {
        console.log("before-mount", this.bounds);
        const map = await this.$services.api.coverageMap();
        this.ready = true;
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
    },
    methods: {
        onMapInitialized(map: any): void {
            console.log("map: initialized");
            this.protectedData.map = map;
        },
        onMapLoaded(map: any): void {
            console.log("map: loaded");
            this.protectedData.map = map;
        },
        async newBounds() {
            const map = this.protectedData.map;
            const bounds = map.getBounds();
            const br = new BoundingRectangle([bounds._sw.lng, bounds._sw.lat], [bounds._ne.lng, bounds._ne.lat]);
            const param = JSON.stringify([br.min, br.max]);
            console.log("new-bounds", br, param);
            await this.$router.push({ name: "mapsCoverageWithBounds", params: { bounds: param } });
        },
        updateMap(): void {
            if (!this.protectedData.map) {
                console.log("map: update-skip.1");
                return;
            }

            const map = this.protectedData.map;

            if (this.bounds) {
                map.fitBounds(this.bounds, { duration: 0 });
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins.scss";

.container-map {
    width: 100%;
    height: calc(100% - 66px);
    margin-top: 0;
    @include position(absolute, 66px null null 0);

    @include bp-down($sm) {
        top: 54px;
        height: calc(100% - 54px);
    }

    ::v-deep .stations-map {
        height: 100% !important;
    }
}

</style>

