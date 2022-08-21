<template>
    <div class="js-cupertinoPane">
        <div
            class="station-hover-summary js-paneContent"
            ref="paneContent"
            :class="{ 'is-pane': hasCupertinoPane }"
            v-if="viewingSummary && station"
        >
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
                        <span v-else-if="hasData">{{ $t("noRecentData") }}</span>
                        <span v-else class="no-data">{{ $t("noData") }}</span>
                    </template>
                    <template v-if="station.status === StationStatus.down">{{ $t("station.inactive") }}</template>
                    <i v-if="latestPrimaryLevel !== null" :style="{ 'background-color': latestPrimaryColor }">
                        <template v-if="station.status === StationStatus.down">-</template>
                        <template v-else>{{ visibleReadingValue | prettyNum }}</template>
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

import { ModuleSensorMeta, SensorDataQuerier, SensorMeta, RecentlyAggregatedWindows } from "@/views/shared/sensor_data_querier";

import { getBatteryIcon } from "@/utilities";
import { BookmarkFactory, ExploreContext, serializeBookmark } from "@/views/viz/viz";
import { interpolatePartner, isCustomisationEnabled } from "./partners";
import { StationStatus } from "@/api";
import { CupertinoPane } from "cupertino-pane";
import { CupertinoEvents } from "cupertino-pane/dist/types/models";

export enum VisibleReadings {
    Current,
    Last72h,
}

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
            this.updatePaneHeights();
        },
    },
    data(): {
        viewingSummary: boolean;
        sensorMeta: SensorMeta | null;
        readings: RecentlyAggregatedWindows | null;
        hasData: boolean;
        StationStatus: any;
        isMobileView: boolean;
        cupertinoPane: CupertinoPane | null;
    } {
        return {
            viewingSummary: true,
            sensorMeta: null,
            readings: null,
            hasData: false,
            StationStatus: StationStatus,
            isMobileView: window.screen.availWidth < 500,
            cupertinoPane: null,
        };
    },
    async mounted() {
        if (this.hasCupertinoPane && this.isMobileView) {
            this.initCupertinoPane();
        }

        if (this.sensorDataQuerier) {
            const recently = await this.sensorDataQuerier.queryRecently(this.station.id);
            this.sensorMeta = await this.sensorDataQuerier.querySensorMeta();
            this.readings = recently.windows;
            this.hasData = recently.stations[this.station.id].last != null;
        }
    },
    destroyed() {
        this.cupertinoPane.destroy({ animate: false });
    },
    computed: {
        ...mapGetters({ projectsById: "projectsById" }),
        visibleSensor(): ModuleSensorMeta | null {
            if (this.sensorMeta && this.readings && 72 in this.readings && this.readings[72].length > 0) {
                const sensorId = this.readings[72][0].sensorId;
                return this.sensorMeta.findSensorById(sensorId);
            }
            return null;
        },
        visibleReadingValue(): number | null {
            const sensor = this.visibleSensor;
            if (sensor && this.readings && 72 in this.readings && this.readings[72].length > 0) {
                // console.log(this.visibleReadings, this.readings);
                let value: number | undefined;
                if (this.visibleReadings == VisibleReadings.Current) {
                    value = this.readings[72][0].last;
                } else {
                    if (sensor.aggregationFunction == "max") {
                        // TODO Pull into helper
                        value = this.readings[72][0].max;
                    } else {
                        value = this.readings[72][0].avg;
                    }
                }
                return value === undefined ? null : value;
            }

            return null;
        },
        thresholds() {
            const sensor = this.visibleSensor;
            if (sensor && sensor.viz && sensor.viz.length > 0) {
                if (sensor.viz[0].thresholds) {
                    return sensor.viz[0].thresholds;
                }
            }
            return null;
        },
        visibleLevel() {
            const value = this.visibleReadingValue;
            if (value !== null) {
                const thresholds = this.thresholds;
                if (thresholds) {
                    const level = thresholds.levels.find((level) => level.start <= value && level.value > value);
                    return level ?? null;
                }
            }
            return null;
        },
        latestPrimaryLevel(): any {
            if (this.visibleReadingValue === null) {
                return null;
            }
            const level = this.visibleLevel;
            return level?.plainLabel?.enUS || level?.mapKeyLabel?.enUS;
        },
        latestPrimaryColor(): string {
            if (this.visibleReadingValue === null) {
                return "#cccccc";
            }
            const level = this.visibleLevel;
            return level?.color || "#00CCFF";
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
        async initCupertinoPane(): void {
            const paneContentEl = document.querySelector(".js-paneContent") as HTMLElement;
            const generalRowEl = document.querySelector(".js-generalRow") as HTMLElement;
            this.cupertinoPane = new CupertinoPane(".js-cupertinoPane", {
                parentElement: "body",
                breaks: {
                    top: { enabled: true, height: paneContentEl.scrollHeight, bounce: true },
                    // add padding top of container and margin of general row
                    middle: { enabled: true, height: generalRowEl.scrollHeight + 25 + 10, bounce: true },
                    bottom: { enabled: true, height: generalRowEl.scrollHeight + 25 + 10, bounce: true },
                },
                bottomClose: true,
                buttonDestroy: false,
            });

            this.cupertinoPane.present({ animate: true }).then();
        },
        updatePaneHeights(): void {
            this.$nextTick(() => {
                const paneContentEl = document.querySelector(".js-paneContent") as HTMLElement;
                const generalRowEl = document.querySelector(".js-generalRow") as HTMLElement;
                this.cupertinoPane.setBreakpoints({
                    top: { enabled: true, height: paneContentEl.scrollHeight, bounce: true },
                    middle: { enabled: true, height: generalRowEl.scrollHeight + 25 + 10, bounce: true },
                    bottom: { enabled: true, height: generalRowEl.scrollHeight + 25 + 10, bounce: true },
                });
            });
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
        color: #cccccc;
        font-family: $font-family-bold;
    }
}
</style>
