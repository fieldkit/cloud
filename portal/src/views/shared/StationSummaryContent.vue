<template>
    <div class="row general-row" v-if="station">
        <div class="image-container">
            <StationPhoto :station="station" />
        </div>
        <div class="station-details">
            <div class="station-name">{{ station.name }}</div>
            <div class="station-seen" v-if="station.updatedAt">
                Last Seen
                <span class="small-light">{{ station.updatedAt | prettyDate }}</span>
            </div>
            <div class="station-battery" v-if="station.battery">
                <img class="battery" alt="Battery Level" :src="getBatteryIcon()" />
                <span class="small-light">{{ station.battery | integer }}%</span>
            </div>
            <div v-for="module in station.modules" v-bind:key="module.id" class="module-icon-container">
                <img v-if="!module.internal" alt="Module Icon" class="small-space" :src="getModuleIcon(module)" />
            </div>
        </div>
        <slot></slot>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import * as utils from "@/utilities";

export default Vue.extend({
    name: "StationSummaryContent",
    components: {
        ...CommonComponents,
    },
    data: () => {
        return {};
    },
    props: {
        station: {
            type: Object,
            default: null,
        },
    },
    methods: {
        getBatteryIcon() {
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
        getModuleIcon(module) {
            return this.$loadAsset(utils.getModuleImg(module));
        },
    },
    filters: {
        integer: (value) => {
            if (!value) return "";
            return Math.round(value);
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.image-container {
    width: 124px;
    text-align: center;
    padding-right: 11px;
    height: 100px;
}
.image-container img {
    max-width: 124px;
    max-height: 100px;
    padding: 5px;
}
.station-name {
    font-size: 18px;
    font-family: var(--font-family-bold);
    margin-bottom: 2px;
}
.station-synced {
    font-size: 14px;
}
.station-battery {
    margin: 5px 0 0;
    display: flex;
    line-height: 13px;
}
.battery {
    width: 20px;
    height: 11px;
    padding-right: 5px;
}
.module-icon-container {
    float: left;
    margin-right: 10px;
    margin-top: 5px;
    box-shadow: 0 2px 4px 0 var(--color-border);
    border-radius: 50%;
    display: flex;

    img {
        width: 22px;
        height: 22px;
    }
}
.general-row {
    display: flex;
    flex-direction: row;
    position: relative;
    padding-right: 23px;
}
.station-details {
    text-align: left;
}
.icon {
    padding-right: 7px;
}
.small-light {
    font-size: 12px;
    color: #6a6d71;
}
</style>
