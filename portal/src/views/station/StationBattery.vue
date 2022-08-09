<template>
    <div class="station-battery-container">
        <div class="station-seen" v-if="station.lastReadingAt">
            {{ $t("station.lastSeen") }}
            <span class="small-light">{{ station.lastReadingAt | prettyDateTime }}</span>
        </div>
        <span v-if="station.status === StationStatus.down" class="small-light inactive">({{ $t("station.inactive") }})</span>
        <div class="station-battery" v-if="station.battery">
            <img class="battery" :alt="$t('station.batteryLevel')" :src="getBatteryIcon()" />
            <span class="small-light">{{ station.battery | integer }}%</span>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { DisplayStation, StationStatus } from "@/store";
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
    data: () => {
        return {
            StationStatus: StationStatus,
        };
    },
    methods: {
        getBatteryIcon() {
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/variables";

.station-seen {
    font-size: 14px;
    font-family: $font-family-bold;
    align-self: flex-start;
    margin-right: 5px;
}

.station-battery-container {
    display: flex;
}

.battery {
    width: 20px;
    height: 11px;
    padding-right: 3px;
    margin-left: 10px;
}

.inactive {
    font-size: 14px;
}
</style>
