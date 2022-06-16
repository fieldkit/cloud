<template>
    <div class="row general-row" v-if="station">
        <div class="image-container">
            <StationPhoto :station="station" />
        </div>
        <div class="station-details">
            <div class="station-name">
                {{ station.name }}
                <slot name="top-right-actions"></slot>
            </div>

            <div
                class="row where-row"
                v-if="stationLocationName || station.placeNameNative || station.placeNameOther || station.placeNameNative"
            >
                <div class="location-container">
                    <div v-if="stationLocationName || station.placeNameOther">
                        <i class="icon icon-location" />
                        <template>
                            {{ stationLocationName ? stationLocationName : station.placeNameOther }}
                        </template>
                    </div>
                    <div v-if="station.placeNameNative">
                        <img alt="Location" src="@/assets/icon-location.svg" class="icon" />
                        <span>
                            Native Lands:
                            <span class="bold">{{ station.placeNameNative }}</span>
                        </span>
                    </div>
                </div>
            </div>

            <div class="station-modules">
                <div v-for="(module, index) in station.modules" v-bind:key="index" class="module-icon-container">
                    <img alt="Module Icon" class="small-space" :src="getModuleIcon(module)" />
                </div>
            </div>

            <slot name="extra-detail"></slot>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import * as utils from "@/utilities";
import { DisplayStation } from "@/store";
import { getPartnerCustomizationWithDefault, PartnerCustomization } from "@/views/shared/partners";

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
            type: Object as PropType<DisplayStation>,
            default: null,
        },
    },
    computed: {
        stationLocationName(): string {
            return this.partnerCustomization().stationLocationName(this.station);
        },
    },
    methods: {
        getBatteryIcon() {
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
        getModuleIcon(module) {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        partnerCustomization(): PartnerCustomization {
            return getPartnerCustomizationWithDefault();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.image-container {
    flex: 0 0 93px;
    text-align: center;
    margin-right: 11px;
}

.image-container img {
    width: 100%;
    border-radius: 3px;
}

.station-name {
    font-size: 18px;
    font-family: var(--font-family-bold);
    margin-bottom: 3px;
    padding-right: 45px;
}

.station-synced {
    font-size: 14px;
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
}

.location-container {
    flex-direction: column;
    display: flex;
    margin-top: 2px;

    .summary-content & {
        flex-direction: row;
    }

    > div {
        @include flex(flex-start);
        margin-bottom: 5px;

        &:first-of-type {
            .summary-content & {
                margin-right: 10px;
            }
        }
    }
}

.icon {
    padding-right: 5px;
}

.station-modules {
    @include flex();
    margin-left: -2px;
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
    padding-bottom: 5px;
    color: var(--color-dark);
}
</style>
