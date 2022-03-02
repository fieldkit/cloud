<template>
    <div>
        {{ label }}

        <slot />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { TranslateResult } from "vue-i18n";
import { isCustomisationEnabled } from "../PartnerCustomisationHelper";

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
            if (isCustomisationEnabled()) {
                return this.$t(this.sensorsKey);
            }
            return this.$t(this.stationsKey);
        },
    },
});
</script>

<style scoped></style>
