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

            <template v-if="isCustomisationEnabled()">
                <div class="row where-row">
                    <div v-if="neighborhood || borough" class="location-container">
                        <i class="icon icon-location" />
                        <template v-if="neighborhood">{{ neighborhood }}</template>
                        <template v-if="neighborhood && borough">{{ ", " }}</template>
                        <template v-if="borough">{{ borough }}</template>
                    </div>
                    <div v-if="deploymentDate || deployedBy" class="location-container">
                        <template v-if="deploymentDate">{{ $t("station.deployedOn") }} {{ deploymentDate }}</template>
                        <template v-if="deployedBy">{{ $t("station.by") }} {{ deployedBy }}</template>
                    </div>
                </div>
            </template>

            <template v-else>
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
                            <i class="icon icon-location" />
                            <span>
                                Native Lands:
                                <span class="bold">{{ station.placeNameNative }}</span>
                            </span>
                        </div>
                    </div>
                </div>
            </template>

            <div v-if="!isCustomisationEnabled()" class="station-modules">
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
import { getPartnerCustomizationWithDefault, isCustomisationEnabled, PartnerCustomization } from "@/views/shared/partners";

export default Vue.extend({
    name: "StationSummaryContent",
    components: {
        ...CommonComponents,
    },
    data: () => {
        return {
            neighborhood: null,
            borough: null,
            deployedBy: null,
            deploymentDate: null,
        };
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
    mounted(): void {
        this.setPartnerAttributes();
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
        isCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
        getAttributeValue(attrName: string): any {
            if (this.station) {
                const value = this.station.attributes.find((attr) => attr.name === attrName)?.stringValue;
                return value ? value : null;
            }
        },
        setPartnerAttributes(): void {
            this.neighborhood = this.getAttributeValue("Neighborhood");
            this.borough = this.getAttributeValue("Borough");
            this.deploymentDate = this.getAttributeValue("Deployment Date");
            this.deployedBy = this.getAttributeValue("Deployed By");
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
    display: flex;
    margin-bottom: 2px;
    font-size: 14px;

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
    transform: translateY(1px);
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
    margin-top: 5px;
    color: var(--color-dark);
}
</style>
