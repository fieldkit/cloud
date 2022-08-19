<template>
    <StandardLayout @sidebar-toggle="sidebarToggle()" :sidebarNarrow="true" :clipStations="true">
        <div class="project-public project-container" v-if="displayProject">
            <ProjectDetailCard :project="project"></ProjectDetailCard>

            <template v-if="viewType === 'list'">
                <div class="stations-list" v-if="projectStations && projectStations.length > 0">
                    <div v-for="station in projectStations" v-bind:key="station.id">
                        <StationHoverSummary
                            class="summary-container"
                            :station="station"
                            :sensorDataQuerier="sensorDataQuerier"
                            :exploreContext="exploreContext"
                            :visibleReadings="visibleReadings"
                            v-slot="{ sensorDataQuerier }"
                        >
                            <TinyChart :station-id="station.id" :station="station" :querier="sensorDataQuerier" />
                        </StationHoverSummary>
                    </div>
                </div>
            </template>

            <template v-if="viewType === 'map'">
                <!-- fixme: currently restricted to floodnet project -->
                <div class="map-legend" v-if="id === 174 && levels.length > 0" :class="{ collapsed: legendCollapsed }">
                    <a class="legend-toggle" @click="legendCollapsed = !legendCollapsed">
                        <i class="icon icon-chevron-right"></i>
                    </a>
                    <div class="legend-content">
                        <div class="legend-content-wrap">
                            <h4>Current {{ keyTitle }}</h4>
                            <div class="legend-item" v-for="(item, idx) in levels" :key="idx">
                                <span class="legend-dot" :style="{ color: item.color }">&#x25CF;</span>
                                <span>{{ item.mapKeyLabel ? item.mapKeyLabel["enUS"] : item.label["enUS"] }}</span>
                            </div>
                            <div class="legend-item" v-if="hasStationsWithoutData">
                                <span class="legend-dot" style="color: #ccc">&#x25CF;</span>
                                <span>{{ $t("map.legend.noData") }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="container-map">
                    <StationsMap
                        @show-summary="showSummary"
                        :mapped="mappedProject"
                        :layoutChanges="layoutChanges"
                        :showStations="project.showStations"
                        :mapBounds="mapBounds"
                    />
                </div>

                <StationHoverSummary
                    v-if="activeStation"
                    :station="activeStation"
                    :sensorDataQuerier="sensorDataQuerier"
                    :exploreContext="exploreContext"
                    :visibleReadings="visibleReadings"
                    :hasCupertinoPane="true"
                    @close="onCloseSummary"
                    v-slot="{ sensorDataQuerier }"
                >
                    <TinyChart :station-id="activeStation.id" :station="activeStation" :querier="sensorDataQuerier" />
                </StationHoverSummary>
            </template>
        </div>
        <div class="view-type-container">
            <label class="toggle-btn">
                <input type="checkbox" v-model="recentMapMode" />
                <span :class="{ active: !recentMapMode }">{{ $t("map.toggle.current") }}</span>
                <i></i>
                <span :class="{ active: recentMapMode }">{{ $t("map.toggle.recent") }}</span>
            </label>
            <div class="view-type">
                <div class="view-type-map" v-bind:class="{ active: viewType === 'map' }" v-on:click="switchView('map')">
                    {{ $t("map.toggle.map") }}
                </div>
                <div class="view-type-list" v-bind:class="{ active: viewType === 'list' }" v-on:click="switchView('list')">
                    {{ $t("map.toggle.list") }}
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import * as utils from "../../utilities";

import { mapState, mapGetters } from "vuex";
import { ActionTypes, GlobalState, ProjectModule, DisplayStation, Project, MappedStations, BoundingRectangle } from "@/store";
import { SensorDataQuerier } from "@/views/shared/sensor_data_querier";

import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import StationsMap from "../shared/StationsMap.vue";
import StationHoverSummary, { VisibleReadings } from "@/views/shared/StationHoverSummary.vue";
import TinyChart from "@/views/viz/TinyChart.vue";
import CommonComponents from "@/views/shared";
import ProjectDetailCard from "@/views/projects/ProjectDetailCard.vue";
import { ExploreContext } from "@/views/viz/common";

import { getPartnerCustomizationWithDefault, isCustomisationEnabled } from "@/views/shared/partners";

export default Vue.extend({
    name: "ProjectBigMap",
    components: {
        ...CommonComponents,
        StationsMap,
        StationHoverSummary,
        StandardLayout,
        ProjectDetailCard,
        TinyChart,
    },
    data(): {
        layoutChanges: number;
        activeStationId: number | null;
        viewType: string;
        recentMapMode: boolean;
        legendCollapsed: boolean;
        isMobileView: boolean;
    } {
        return {
            layoutChanges: 0,
            activeStationId: null,
            viewType: "map",
            recentMapMode: false,
            legendCollapsed: false,
            isMobileView: window.screen.availWidth <= 768,
        };
    },
    props: {
        id: {
            required: true,
            type: Number,
        },
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            userStations: (s: GlobalState) => Object.values(s.stations.user.stations),
            displayProject() {
                return this.$getters.projectsById[this.id];
            },
        }),
        visibleReadings(): VisibleReadings {
            return this.recentMapMode ? VisibleReadings.Last72h : VisibleReadings.Current;
        },
        project(): Project {
            return this.displayProject.project;
        },
        sensorDataQuerier(): SensorDataQuerier {
            return new SensorDataQuerier(
                this.$services.api,
                this.projectStations.map((s: DisplayStation) => s.id)
            );
        },
        projectStations(): DisplayStation[] {
            return this.$getters.projectsById[this.id].stations.slice().sort((a, b) => b.latestPrimary - a.latestPrimary);
        },
        mappedProject(): MappedStations | null {
            return this.$getters.projectsById[this.id].mapped;
        },
        activeStation(): DisplayStation | null {
            if (this.activeStationId) {
                return this.$getters.stationsById[this.activeStationId];
            }
            return null;
        },
        projectModules(): { name: string; url: string }[] {
            return this.$getters.projectsById[this.displayProject.id].modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
        mapBounds(): BoundingRectangle {
            if (this.project.bounds?.min && this.project.bounds?.max) {
                return new BoundingRectangle(this.project.bounds?.min, this.project.bounds?.max);
            }

            return MappedStations.defaultBounds();
        },
        exploreContext(): ExploreContext {
            return new ExploreContext(this.project.id, true);
        },
        stationsWithData(): DisplayStation[] {
            return this.displayProject.stations.filter((station) => station.latestPrimary != null);
        },
        hasStationsWithoutData(): boolean {
            return this.stationsWithData.length < this.projectStations.length;
        },
        keyTitle(): string {
            if (this.stationsWithData.length > 0) {
                return this.getThresholds(this.stationsWithData).label["enUS"];
            } else {
                return "";
            }
        },
        levels(): object[] {
            if (this.stationsWithData.length > 0) {
                return this.getThresholds(this.stationsWithData).levels.filter((d) => d["label"]);
            } else {
                return [];
            }
        },
        partnerCustomization() {
            return getPartnerCustomizationWithDefault();
        },
    },
    watch: {
        id(): Promise<any> {
            if (this.id) {
                return this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id });
            }
            return Promise.resolve();
        },
    },
    beforeMount(): Promise<any> {
        if (this.id) {
            return this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id });
        }
        return Promise.resolve();
    },
    mounted() {
        console.log("dada");
    },
    methods: {
        switchView(type: string): void {
            this.viewType = type;
            this.layoutChanges++;
        },
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        showSummary(station: DisplayStation): void {
            console.log("showSummay", station);
            this.activeStationId = station.id;
        },
        onCloseSummary(): void {
            this.activeStationId = null;
        },
        getThresholds(stations: DisplayStation[]): object[] {
            try {
                return stations[0].configurations.all[0].modules[0].sensors[0].meta.viz[0].thresholds;
            } catch (error) {
                return [];
            }
        },
        sidebarToggle() {
            setTimeout(() => {
                this.layoutChanges++;
            }, 250);
        },
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/project";
@import "../../scss/global";

.container-map {
    width: 100%;
    height: calc(100% - 157px);
    margin-top: 0;
    @include position(absolute, 157px null null 0);

    @include bp-down($sm) {
        top: 54px;
        height: calc(100% - 54px);
    }

    ::v-deep .stations-map {
        height: 100% !important;
    }
}
.map-legend {
    display: flex;
    flex-direction: column;
    border: 1px solid #f4f5f7;
    border-radius: 3px;
    z-index: $z-index-top;
    position: absolute;
    bottom: 80px;
    right: 0;
    box-sizing: border-box;
    background-color: #fcfcfc;
    text-align: left;
    font-family: $font-family-medium;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);

    h4 {
        margin: 0 0 0.5em 0;
        font-family: $font-family-bold;
    }

    &.collapsed {
        .legend-toggle {
            left: -25px;
        }

        .legend-content {
            max-width: 0;
        }

        .icon-chevron-right {
            transform: rotate(180deg) translateX(0);
        }
    }

    .icon-chevron-right {
        transform: translateX(2px);
    }

    .legend-content {
        overflow: hidden;
        white-space: nowrap;
        max-width: 230px;
        transition: all 0.25s ease-in-out;
    }

    .legend-content-wrap {
        padding: 15px 20px 15px 15px;
    }

    .legend-toggle {
        @include position(absolute, 11px null null -24px);
        @include flex(center, center);
        height: 35px;
        width: 24px;
        font-size: 22px;
        color: #979797;
        background-color: #fcfcfc;
        box-shadow: -1px 2px 2px 0 rgb(0 0 0 / 13%);
        cursor: pointer;
    }

    .legend-item {
        margin-bottom: 0.2em;
    }
    .legend-item:last-child {
        margin-bottom: 0px;
    }
    .legend-dot {
        margin-right: 10px;
        font-size: 1.5em;
    }
}

::v-deep .station-hover-summary {
    width: 359px;
    top: calc(50% - 100px);
    left: calc(50% - 180px);

    @include bp-down($xs) {
        width: 100%;
        left: 0;
        border-radius: 10px;
        padding: 25px 10px 12px 10px;

        &.is-pane {
            top: 0;
        }

        .close-button {
            display: none;
        }

        .navigate-button {
            width: 14px;
            height: 14px;
            right: -3px;
            top: -17px;
        }

        .image-container {
            flex-basis: 62px;
            margin-right: 10px;
        }

        .station-name {
            font-size: 14px;
        }

        .explore-button {
            margin-top: 15px;
            margin-bottom: 10px;
        }
    }
}

::v-deep .stations-list {
    @include flex();
    flex-wrap: wrap;
    padding: 160px 70px;
    margin: -20px;
    width: 100%;
    box-sizing: border-box;

    @include bp-down($md) {
        padding: 100px 20px;
        margin: 30px -20px -20px;
    }

    @include bp-down($sm) {
        justify-content: center;
    }

    @include bp-down($xs) {
        padding: 80px 0px;
        margin: 55px 0 -5px 0;
        transform: translateX(10px);
        width: calc(100% - 20px);
    }

    .summary-container {
        z-index: 0;
        position: unset;
        margin: 20px;
        flex-basis: 389px;
        box-sizing: border-box;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);

        @include bp-down($md) {
            padding: 19px 11px;
            flex-basis: calc(50% - 40px);
        }

        @include bp-down($sm) {
            justify-self: center;
            flex: 1 1 389px;
            max-width: 389px;
            margin: 10px 0;
        }

        @include bp-down($xs) {
            margin: 5px 0;
            width: auto;
        }

        .close-button {
            display: none;
        }

        .navigate-button {
            right: 0;
        }
    }
}

.view-type {
    width: 115px;
    height: 38px;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
    border: solid 1px #f4f5f7;
    background-color: #ffffff;
    cursor: pointer;
    margin-left: 30px;
    @include flex(center, center);

    @include bp-down($sm) {
        width: 98px;
        margin-left: 5px;
    }

    @media screen and (max-width: 350px) {
        width: 88px;
    }

    &-container {
        z-index: $z-index-top;
        margin: 0;
        @include flex(center, center);
        @include position(absolute, 90px 25px null null);

        @include bp-down($sm) {
            @include position(absolute, 115px 10px null null);
        }
    }

    > div {
        flex-basis: 50%;
        height: 100%;
        @include flex(center, center);

        @include bp-down($sm) {
            font-size: 14px;
        }

        &.active {
            font-family: $font-family-bold;
        }
    }

    &-map {
        flex-basis: 50%;
        border-right: solid 1px $color-border;
    }

    .icon {
        font-size: 18px;
    }
}

.toggle-btn {
    cursor: pointer;
    z-index: $z-index-top;
    position: relative;
    font-size: 14px;
    -webkit-tap-highlight-color: transparent;
    font-family: $font-family-medium !important;
    height: 38px;
    display: flex;
    align-items: center;

    @include bp-down($sm) {
        background-color: #fff;
        padding: 0 10px;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
        border: solid 1px #f4f5f7;
    }

    @media screen and (max-width: 350px) {
        padding: 0 8px;
    }

    * {
        font-family: $font-family-medium !important;
    }

    span {
        opacity: 0.5;

        &.active {
            opacity: 1;
            color: $color-fieldkit-primary;

            body.floodnet & {
                color: $color-dark;
            }
        }
    }
}
.toggle-btn i {
    position: relative;
    display: inline-block;
    margin: 0 10px;
    width: 27px;
    height: 16px;
    background-color: #e6e6e6;
    border-radius: 20px;
    vertical-align: text-bottom;
    transition: all 0.3s linear;
    user-select: none;
}
.toggle-btn i::before {
    content: "";
    position: absolute;
    left: 0;
    width: 27px;
    background-color: #fff;
    border-radius: 11px;
    transform: translate3d(2px, 2px, 0) scale3d(1, 1, 1);
    transition: all 0.25s linear;
}
.toggle-btn i::after {
    content: "";
    position: absolute;
    left: 0;
    width: 12px;
    height: 12px;
    background-color: #fff;
    border-radius: 50%;
    transform: translate3d(2px, 2px, 0);
    transition: all 0.2s ease-in-out;
}
.toggle-btn:active i::after {
    width: 28px;
    transform: translate3d(2px, 2px, 0);
}
.toggle-btn:active input:checked ~ i::after {
    transform: translate3d(16px, 2px, 0);
}
.toggle-btn input {
    display: none;
}
.toggle-btn input ~ i {
    background-color: $color-primary;
}
.toggle-btn input:checked ~ i::before {
    transform: translate3d(18px, 2px, 0) scale3d(0, 0, 0);
}
.toggle-btn input:checked ~ i::after {
    transform: translate3d(13px, 2px, 0);
}

::v-deep .mapboxgl-ctrl-geocoder {
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
    border: solid 1px #f4f5f7;
    border-radius: 0;
    height: 40px;
    margin: 30px 0 0 30px;

    @include bp-down($sm) {
        margin: 61px 0 0 10px;
        width: 40px;
    }

    &.mapboxgl-ctrl-geocoder--collapsed {
        min-width: 40px;
    }

    &:not(.mapboxgl-ctrl-geocoder--collapsed) {
        @include bp-down($xs) {
            min-width: calc(100vw - 20px) !important;
        }
    }

    input {
        outline: none;
        height: 37px;
        padding-left: 38px;
        font-size: 16px;
    }
}

::v-deep .mapboxgl-ctrl-geocoder--icon-search {
    top: 9px;
    left: 8px;
}

::v-deep .mapboxgl-ctrl-bottom-left {
    margin-left: 30px;
}
</style>
