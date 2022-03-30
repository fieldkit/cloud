<template>
    <div>
        {{ label }}

        <slot />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { TranslateResult } from "vue-i18n";
import { isCustomisationEnabled } from "./partners";

// Keeping this logic in this file.
export function callStationsStations() {
    return !isCustomisationEnabled();
}

export default Vue.extend({
    name: "StationOrSensor",
    props: {
        stationsKey: {
            type: String,
            default: "layout.side.stations.title",
        },
        sensorsKey: {
            type: String,
            default: "layout.side.sensors.title",
        },
    },
    computed: {
        label(): TranslateResult {
            if (callStationsStations()) {
                return this.$t(this.stationsKey);
            }
            return this.$t(this.sensorsKey);
        },
    },
});
</script>

<style scoped></style>
