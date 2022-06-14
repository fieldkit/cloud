<template>
    <div class="station-hover-summary" v-if="viewingSummary && station">
        <StationSummaryContent :station="station">
            <img alt="Close" src="@/assets/icon-close.svg" class="close-button" v-on:click="wantCloseSummary" />
            <img :alt="$tc('station.navigateToStation')" class="navigate-button" src="@/assets/tooltip-blue.svg" @click="openStationPageTab" />
        </StationSummaryContent>

        <div class="readings-container" v-if="readings">
            <div class="title">Latest Readings</div>
            <LatestStationReadings :id="station.id" @layoutChange="layoutChange" />
        </div>

        <div class="explore-button" v-if="explore" v-on:click="onClickExplore">Explore Data</div>

        <div class="flex">
            <div class="station-seen" v-if="station.updatedAt">
                Last Seen
                <span class="small-light">{{ station.updatedAt | prettyDateTime }}</span>
            </div>

            <div class="station-battery" v-if="station.battery">
                <img class="battery" alt="Battery Level" :src="getBatteryIcon()" />
                <span class="small-light">{{ station.battery | integer }}%</span>
            </div>
        </div>
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
        openStationPageTab() {
            const routeData = this.$router.resolve({ name: "viewStationFromMap", params: { stationId: this.station.id } });
            window.open(routeData.href, "_blank");
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
    padding: 27px 20px 17px;
    width: 389px;
    box-sizing: border-box;

    ::v-deep .station-name {
        font-size: 16px;
    }

    ::v-deep .image-container {
        border-radius: 5px;
        padding: 0;
        margin-right: 14px;
        overflow: hidden;

        img {
            padding: 0;
        }
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

.small-light {
    font-size: 12px;
    color: #6a6d71;
}

.station-seen {
    font-size: 14px;
    font-family: var(--font-family-bold);
    align-self: flex-start;
    margin-right: 15px;
}

.station-battery {
    display: flex;
    align-items: center;
    line-height: 13px;

    .battery {
        width: 22px;
        margin-top: -2px;
        margin-right: 3px;
    }
}
</style>
