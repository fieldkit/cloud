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
                <!-- TODO This should be handled via partner customization. -->
                <div class="map-legend" v-if="id === 174 && levels.length > 0" :class="{ collapsed: legendCollapsed }">
                    <a class="legend-toggle" @click="legendCollapsed = !legendCollapsed">
                        <i class="icon icon-chevron-right"></i>
                    </a>
                    <div class="legend-content">
                        <div class="legend-content-wrap">
                            <h4>{{ keyTitle }}</h4>
                            <div class="legend-item" v-for="(item, idx) in levels" :key="idx">
                                <span class="legend-dot" :style="{ color: item.color }">&#x25CF;</span>
                                <span>{{ item.mapKeyLabel ? item.mapKeyLabel["enUS"] : item.label["enUS"] }}</span>
                            </div>
                            <div class="legend-item">
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
                        :visibleReadings="visibleReadings"
                        :mapBounds="mapBounds"
                    />

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
                </div>
            </template>
        </div>

        <MapViewTypeToggle
            :routes="[
                { name: 'viewProjectBigMap', label: 'map.toggle.map', viewType: 'map', params: { id: id } },
                { name: 'viewProjectBigMapList', label: 'map.toggle.list', viewType: 'list', params: { id: id } },
            ]"
        ></MapViewTypeToggle>
    </StandardLayout>
</template>

<script lang="ts">
import _ from "lodash";
import * as utils from "../../utilities";

import { mapGetters, mapState } from "vuex";
import {
    ActionTypes,
    BoundingRectangle,
    DisplayStation,
    GlobalState,
    MappedStations,
    Project,
    ProjectModule,
    VisibleReadings,
} from "@/store";
import { SensorDataQuerier } from "@/views/shared/sensor_data_querier";

import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import StationsMap from "../shared/StationsMap.vue";
import StationHoverSummary from "@/views/shared/StationHoverSummary.vue";
import TinyChart from "@/views/viz/TinyChart.vue";
import CommonComponents from "@/views/shared";
import ProjectDetailCard from "@/views/projects/ProjectDetailCard.vue";
import { ExploreContext } from "@/views/viz/common";

import { getPartnerCustomizationWithDefault, isCustomisationEnabled } from "@/views/shared/partners";
import MapViewTypeToggle from "@/views/shared/MapViewTypeToggle.vue";
import { MapViewType } from "@/api/api";

export default Vue.extend({
    name: "ProjectBigMap",
    components: {
        ...CommonComponents,
        StationsMap,
        StationHoverSummary,
        StandardLayout,
        ProjectDetailCard,
        TinyChart,
        MapViewTypeToggle,
    },
    data(): {
        layoutChanges: number;
        activeStationId: number | null;
        recentMapMode: boolean;
        legendCollapsed: boolean;
        isMobileView: boolean;
    } {
        return {
            layoutChanges: 0,
            activeStationId: null,
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
            return new SensorDataQuerier(this.$services.api);
        },
        projectStations(): DisplayStation[] {
            const stations = this.$getters.projectsById[this.id].stations;
            const sortFactors = _.fromPairs(stations.map((station) => [station.id, station.getSortOrder(this.visibleReadings)]));
            return _.orderBy(
                this.$getters.projectsById[this.id].stations.slice(),
                [(station) => sortFactors[station.id][0], (station) => sortFactors[station.id][1], (station) => sortFactors[station.id][2]],
                ["asc", "desc", "asc"]
            );
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
            return this.displayProject.stations.filter((station) => station.hasData);
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
        viewType(): MapViewType {
            if (this.$route.meta?.viewType) {
                return this.$route.meta.viewType;
            }
            return MapViewType.map;
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
    methods: {
        switchView(): void {
            this.activeStationId = null;
            this.layoutChanges++;
        },
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        showSummary(station: DisplayStation): void {
            console.log("map: show-summay", station);
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
}

.stations-list {
    @include flex();
    flex-wrap: wrap;
    padding: 160px 40px;
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

    ::v-deep .station-hover-summary {
        position: unset;
    }

    ::v-deep .summary-container {
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

::v-deep .mapboxgl-ctrl-geocoder {
    margin: 30px 0 0 30px;

    @include bp-down($sm) {
        margin: 61px 0 0 10px;
    }
}

::v-deep .mapboxgl-ctrl-bottom-left {
    margin-left: 25px;
}
</style>
