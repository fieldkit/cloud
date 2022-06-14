<template>
    <div class="row general-row" v-if="station">
        <div class="image-container">
            <StationPhoto :station="station" />
        </div>
        <div class="station-details">
            <div class="station-name">{{ station.name }}</div>

            <div class="row where-row" v-if="station.placeNameNative || station.placeNameOther || station.placeNameNative">
                <div class="location-container" v-if="station.locationName ? station.locationName : station.placeNameOther">
                    <img alt="Location" src="@/assets/icon-location.svg" class="icon" />
                    <template>{{ station.locationName ? station.locationName : station.placeNameOther }}</template>
                </div>
                <div class="location-container" v-if="station.placeNameNative">
                    <img alt="Location" src="@/assets/icon-location.svg" class="icon" />
                    <span>
                        Native Lands:
                        <span class="bold">{{ station.placeNameNative }}</span>
                    </span>
                </div>

                <div class="flex">
                    <div v-for="module in station.modules" v-bind:key="module.id" class="module-icon-container">
                        <img v-if="!module.internal" alt="Module Icon" class="small-space" :src="getModuleIcon(module)" />
                    </div>
                </div>

                <div v-if="station.location" class="coordinates-row">
                    <div class="coordinate latitude">
                        <div>{{ station.location.latitude | prettyCoordinate }}</div>
                        <div>Latitude</div>
                    </div>
                    <div class="coordinate longitude">
                        <div>{{ station.location.longitude | prettyCoordinate }}</div>
                        <div>Longitude</div>
                    </div>
                </div>
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
    width: 93px;
    text-align: center;
    padding-right: 11px;
}

.image-container img {
    width: 100%;
}

.station-name {
    font-size: 18px;
    font-family: var(--font-family-bold);
    margin-bottom: 5px;
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
    margin-right: 6px;
    border-radius: 50%;
    display: flex;

    img {
        width: 24px;
        height: 24px;
    }
}

.general-row {
    display: flex;
    flex-direction: row;
    position: relative;
}

.station-details {
    text-align: left;
    padding-right: 15px;
}

.icon {
    padding-right: 6px;
}

.small-light {
    font-size: 12px;
    color: #6a6d71;
}

.station-seen {
    font-size: 12px;
}

.coordinates-row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-top: 7px;
}

.where-row {
    display: flex;
    flex-direction: column;
    text-align: left;
    font-size: 14px;
    color: #2c3e50;
}

.where-row > div {
    padding-bottom: 5px;
    @include flex(flex-start);
}
</style>
