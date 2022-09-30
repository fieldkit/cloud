<template>
    <div class="js-cupertinoPane">
        <div class="station-hover-summary" ref="paneContent" :class="{ 'is-pane': hasCupertinoPane }" v-if="viewingSummary && station">
            <StationSummaryContent ref="summaryContent" :station="station">
                <template #top-right-actions>
                    <img alt="Close" src="@/assets/icon-close.svg" class="close-button" v-on:click="wantCloseSummary" />
                    <img
                        :alt="$tc('station.navigateToStation')"
                        class="navigate-button"
                        :src="$loadAsset(interpolatePartner('tooltip-') + '.svg')"
                        @click="openStationPageTab"
                    />
                </template>
            </StationSummaryContent>

            <template v-if="isPartnerCustomisationEnabled()">
                <div class="latest-primary" :style="{ color: latestPrimaryColor }">
                    <template v-if="station.status === StationStatus.up">
                        <template v-if="latestPrimaryLevel !== null">{{ latestPrimaryLevel }}</template>
                        <span v-else-if="hasData" class="no-data">{{ $t("noRecentData") }}</span>
                        <span v-else class="no-data">{{ $t("noData") }}</span>
                    </template>
                    <template v-if="station.status === StationStatus.down">{{ $t("station.inactive") }}</template>
                    <i v-if="latestPrimaryLevel !== null" :style="{ 'background-color': latestPrimaryColor }">
                        <template v-if="station.status === StationStatus.down">-</template>
                        <template v-else>{{ visibleReadingValue | prettyReadingNarrowSpace }}</template>
                    </i>
                    <i v-else :style="{ 'background-color': latestPrimaryColor }">
                        –
                    </i>
                </div>
            </template>

            <slot :station="station" :sensorDataQuerier="sensorDataQuerier"></slot>

            <div class="explore-button" v-if="explore" v-on:click="onClickExplore">Explore Data</div>

            <StationBattery :station="station" />
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";

import Vue, { PropType } from "vue";
import { mapGetters } from "vuex";

import CommonComponents from "@/views/shared";
import StationBattery from "@/views/station/StationBattery.vue";
import StationSummaryContent from "./StationSummaryContent.vue";

import { ModuleSensorMeta, SensorDataQuerier, SensorMeta } from "@/views/shared/sensor_data_querier";
import { VisibleReadings, DecoratedReading } from "@/store";

import { getBatteryIcon } from "@/utilities";
import { BookmarkFactory, ExploreContext, serializeBookmark } from "@/views/viz/viz";
import { interpolatePartner, isCustomisationEnabled } from "./partners";
import { StationStatus } from "@/api";
import { CupertinoPane } from "cupertino-pane";

export default Vue.extend({
    name: "StationHoverSummary",
    components: {
        ...CommonComponents,
        StationSummaryContent,
        StationBattery,
    },
    props: {
        station: {
            type: Object,
            required: true,
        },
        explore: {
            type: Boolean,
            default: true,
        },
        exploreContext: {
            type: Object as PropType<ExploreContext>,
            default: () => {
                return new ExploreContext();
            },
        },
        sensorDataQuerier: {
            type: Object as PropType<SensorDataQuerier>,
            required: false,
        },
        visibleReadings: {
            type: Number as PropType<VisibleReadings>,
            default: VisibleReadings.Current,
        },
        hasCupertinoPane: {
            type: Boolean,
            default: false,
        },
    },
    filters: {
        integer: (value) => {
            if (!value) return "";
            return Math.round(value);
        },
    },
    watch: {
        station(this: any) {
            if (this.cupertinoPane) {
                this.cupertinoPane.present({ animate: true });
            }
        },
    },
    data(): {
        viewingSummary: boolean;
        sensorMeta: SensorMeta | null;
        StationStatus: any;
        isMobileView: boolean;
        cupertinoPane: CupertinoPane | null;
    } {
        return {
            viewingSummary: true,
            sensorMeta: null,
            StationStatus: StationStatus,
            isMobileView: window.screen.availWidth < 500,
            cupertinoPane: null,
        };
    },
    async mounted() {
        if (this.hasCupertinoPane && this.isMobileView) {
            this.initCupertinoPane();
        }
        this.sensorMeta = await this.sensorDataQuerier.querySensorMeta();
    },
    destroyed() {
        this.destroyCupertinoPane();
    },
    computed: {
        ...mapGetters({ projectsById: "projectsById" }),
        visibleSensor(): ModuleSensorMeta | null {
            const primarySensor = this.station.primarySensor;
            if (this.sensorMeta && primarySensor) {
                return this.sensorMeta.findSensorByKey(primarySensor.fullKey);
            }
            return null;
        },
        hasData(): boolean {
            return this.station.hasData;
        },
        decoratedReading(): DecoratedReading | null {
            const readings = this.station.getDecoratedReadings(this.visibleReadings);
            if (readings && readings.length > 0) {
                return readings[0];
            }
            return null;
        },
        visibleReadingValue(): number | null {
            const reading: DecoratedReading | null = this.decoratedReading;
            if (reading) {
                return reading.value;
            }
            return null;
        },
        latestPrimaryLevel(): any {
            const reading: DecoratedReading | null = this.decoratedReading;
            if (reading) {
                return reading?.thresholdLabel;
            }
            return null;
        },
        latestPrimaryColor(): string {
            const reading: DecoratedReading | null = this.decoratedReading;
            if (reading === null) {
                return "#777a80";
            }
            if (reading) {
                return reading?.color;
            }
            return "#00CCFF";
        },
    },
    methods: {
        viewSummary() {
            this.viewingSummary = true;
        },
        onClickExplore() {
            const bm = BookmarkFactory.forStation(this.station.id, this.exploreContext);
            return this.$router.push({
                name: "exploreBookmark",
                query: { bookmark: serializeBookmark(bm) },
            });
        },
        getBatteryIcon() {
            return this.$loadAsset(getBatteryIcon(this.station.battery));
        },
        wantCloseSummary() {
            this.$emit("close");
        },
        layoutChange() {
            this.$emit("layoutChange");
        },
        openStationPageTab() {
            const routeData = this.$router.resolve({ name: "viewStationFromMap", params: { stationId: this.station.id } });
            window.open(routeData.href, "_blank");
        },
        interpolatePartner(baseString) {
            return interpolatePartner(baseString);
        },
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
        async initCupertinoPane(): Promise<void> {
            const paneContentEl = this.$refs["paneContent"] as HTMLDivElement;
            const generalRowEl = (this.$refs["summaryContent"] as Vue).$refs["summaryGeneralRow"] as HTMLDivElement;
            this.cupertinoPane = new CupertinoPane(".js-cupertinoPane", {
                parentElement: "body",
                breaks: {
                    top: { enabled: true, height: paneContentEl.scrollHeight, bounce: true },
                    // add padding top of container and margin of general row
                    middle: { enabled: true, height: generalRowEl.scrollHeight + 25 + 10, bounce: true },
                    bottom: { enabled: true, height: 0 },
                },
                bottomClose: true,
                buttonDestroy: false,
            });
            this.cupertinoPane.present({ animate: true });
        },
        destroyCupertinoPane(): void {
            if (this.cupertinoPane) {
                this.cupertinoPane.destroy();
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.station-hover-summary {
    position: absolute;
    background-color: #ffffff;
    border: solid 1px #d8dce0;
    z-index: 2;
    display: flex;
    flex-direction: column;
    padding: 27px 20px 17px;
    width: 389px;
    box-sizing: border-box;

    ::v-deep .station-name {
        font-size: 16px;
    }

    * {
        font-family: $font-family-light;
    }
}

.location-coordinates {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    flex-basis: 50%;
}

.location-coordinates .coordinate {
    display: flex;
    flex-direction: column;

    > div {
        margin-bottom: 2px;

        &:nth-of-type(2) {
            font-size: 12px;
        }
    }
}

.close-button {
    cursor: pointer;
    @include position(absolute, -7px -5px null null);
}

.navigate-button {
    cursor: pointer;
    @include position(absolute, -10px 20px null null);
}

.readings-container {
    margin-top: 7px;
    padding-top: 9px;
    border-top: 1px solid #f1eeee;
    text-align: left;
    font-size: 14px;
    color: #2c3e50;
}

.readings-container div.title {
    padding-bottom: 13px;
    font-family: var(--font-family-medium);

    body.floodnet & {
        font-family: var(--font-family-bold);
    }
}

.explore-button {
    font-size: 18px;
    font-family: var(--font-family-bold);
    color: #ffffff;
    text-align: center;
    padding: 10px;
    margin: 24px 0 14px 0px;
    background-color: var(--color-secondary);
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}

.icon {
    padding-right: 7px;
}

::v-deep .reading {
    height: 35px;

    .name {
        font-size: 11px;
    }
}

.latest-primary {
    font-size: 12px;
    font-family: $font-family-bold;
    @include flex(center, flex-end);

    @include bp-down($xs) {
        margin-top: 10px;
    }

    i {
        font-style: normal;
        font-size: 10px;
        font-weight: 900;
        border-radius: 50%;
        padding: 5px;
        color: #fff;
        margin-left: 5px;
        min-width: 1.2em;
        text-align: center;
    }

    .no-data {
        color: #777a80;
        font-family: $font-family-bold;
    }
}
</style>
