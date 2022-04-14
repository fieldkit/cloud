<template>
    <div class="station-hover-summary" v-if="viewingSummary && station">
        <StationSummaryContent :station="station">
            <img alt="Close" src="@/assets/icon-close.svg" class="close-button" v-on:click="wantCloseSummary" />
        </StationSummaryContent>

        <div class="row where-row" v-if="station.placeNameNative || station.placeNameOther || station.placeNameNative">
            <div class="location-container" v-if="station.locationName ? station.locationName : station.placeNameOther">
                <img alt="Location" src="@/assets/icon-location.svg" class="icon" />
                <template>{{ station.locationName ? station.locationName : station.placeNameOther }}</template>
            </div>
            <div class="location-container" v-if="station.placeNameNative">
                <img alt="Location" src="@/assets/icon-location.svg" class="icon" />
                <template>Native Lands: {{ station.placeNameNative }}</template>
            </div>

            <div class="coordinates-row">
                <div v-if="station.location" class="location-coordinates valid">
                    <div class="coordinate latitude">
                        <div>{{ station.location.latitude | prettyCoordinate }}</div>
                        <div>Latitude</div>
                    </div>
                    <div class="coordinate longitude">
                        <div>{{ station.location.longitude | prettyCoordinate }}</div>
                        <div>Longitude</div>
                    </div>
                </div>
                <!--
                <div v-else class="location-coordinates missing">
                    <div class="coordinate latitude">
                        --
                        <br />
                        Latitude
                    </div>
                    <div class="coordinate longitude">
                        --
                        <br />
                        Longitude
                    </div>
                </div>
                -->
                <div class="empty"></div>
            </div>
        </div>

        <div class="readings-container" v-if="readings">
            <div class="title">Latest Readings</div>
            <LatestStationReadings :id="station.id" @layoutChange="layoutChange" />
        </div>

        <div class="explore-button" v-if="explore" v-on:click="onClickExplore">Explore Data</div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import StationSummaryContent from "./StationSummaryContent.vue";
import { BookmarkFactory, serializeBookmark, ExploreContext } from "@/views/viz/viz";
import * as utils from "@/utilities";

export default Vue.extend({
    name: "StationHoverSummary",
    components: {
        ...CommonComponents,
        StationSummaryContent,
    },
    data: () => {
        return {
            viewingSummary: true,
        };
    },
    props: {
        station: {
            type: Object,
            required: true,
        },
        readings: {
            type: Boolean,
            default: true,
        },
        explore: {
            type: Boolean,
            default: true,
        },
        exploreContext: {
            type: Object,
            default: () => {
                return new ExploreContext();
            },
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
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
        getModuleIcon(module) {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        wantCloseSummary() {
            this.$emit("close");
        },
        layoutChange() {
            this.$emit("layoutChange");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.station-hover-summary {
    position: absolute;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
    display: flex;
    flex-direction: column;
    padding: 23px 21px;
    width: 389px;
    box-sizing: border-box;

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
.coordinates-row {
    display: flex;
    flex-direction: row;
    margin-top: 7px;
}
.where-row > div {
    padding-bottom: 8px;
    @include flex(center);
}
.where-row {
    border-top: 1px solid #f1eeee;
    margin-top: 18px;
    padding-top: 12px;
    display: flex;
    flex-direction: column;
}
.close-button {
    cursor: pointer;
    @include position(absolute, -7px -5px null null);
}
.where-row {
    text-align: left;
    font-size: 14px;
    color: #2c3e50;
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
    margin: 24px 0 0 0px;
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
</style>
