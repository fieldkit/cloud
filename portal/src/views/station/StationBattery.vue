<template>
    <div class="station-battery-container">
        <div class="station-seen" v-if="station.updatedAt">
            {{ $t("station.lastSeen") }}
            <span class="small-light">{{ station.updatedAt | prettyDateTime }}</span>
        </div>
        <div class="station-battery" v-if="station.battery">
            <img class="battery" :alt="$t('station.batteryLevel')" :src="getBatteryIcon()" />
            <span class="small-light">{{ station.battery | integer }}%</span>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { DisplayStation } from "@/store";
import * as utils from "@/utilities";

export default Vue.extend({
    name: "StationBattery.vue",
    props: {
        station: {
            type: Object as PropType<DisplayStation>,
            default: null,
        },
    },
    filters: {
        integer: (value) => {
            if (!value) return "";
            return Math.round(value);
        },
    },
    methods: {
        getBatteryIcon() {
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
    },
});
</script>

<style scoped lang="scss">
.station-seen {
    font-size: 14px;
    font-family: var(--font-family-bold);
    align-self: flex-start;
    margin-right: 15px;
}

.small-light {
    font-size: 12px;
    color: #6a6d71;
    font-family: var(--font-family-light);
}

.station-battery-container {
    display: flex;
}

.battery {
    width: 20px;
    height: 11px;
    padding-right: 3px;
}
</style>
