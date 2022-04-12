<template>
    <StandardLayout @sidebar-toggle="layoutChanges++" :sidebarNarrow="true">
        <div class="project-public project-container" v-if="displayProject">
            <div class="project-detail-card">
                <div class="photo-container">
                    <ProjectPhoto :project="project" :image-size="150" />
                </div>
                <div class="detail-container">
                    <p class="detail-title">{{ project.name }}</p>
                    <div class="detail-description">{{ project.description }}</div>
                    <router-link :to="{ name: 'viewProject', params: { id: id } }" class="link">Project Dashboard ></router-link>
                </div>
            </div>
            <!-- fixme: currently restricted to floodnet project -->
            <div class="map-legend" v-if="id === 174 && levels.length > 0">
                <h4>{{ keyTitle }}</h4>
                <div class="legend-item" v-for="(item, idx) in levels" :key="idx">
                    <span class="legend-dot" :style="{ color: item.color }">&#x25CF;</span>
                    <span>{{ item.label["enUS"] }}</span>
                </div>
                <div class="legend-item" v-if="hasStationsWithoutData">
                    <span class="legend-dot" style="color: #ccc">&#x25CF;</span>
                    <span>No Data</span>
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
                :readings="false"
                :exploreContext="exploreContext"
                @close="onCloseSummary"
                v-bind:key="activeStation.id"
            />
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import { GlobalState } from "@/store/modules/global";
import * as ActionTypes from "@/store/actions";
import * as utils from "../../utilities";
import { ProjectModule, DisplayStation, Project, MappedStations, BoundingRectangle } from "@/store";
import StationsMap from "../shared/StationsMap.vue";
import StationHoverSummary from "@/views/shared/StationHoverSummary.vue";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";

import { ExploreContext } from "@/views/viz/common";

export default Vue.extend({
    name: "ProjectBigMap",
    components: {
        ...CommonComponents,
        StationsMap,
        StationHoverSummary,
        StandardLayout,
    },
    data(): {
        layoutChanges: number;
        activeStationId: number | null;
    } {
        return {
            layoutChanges: 0,
            activeStationId: null,
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
        project(): Project {
            return this.displayProject.project;
        },
        projectStations(): DisplayStation[] {
            return this.$getters.projectsById[this.id].stations;
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
            return new ExploreContext(this.project.id);
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
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/project";
@import "../../scss/global";

.project-detail-card {
    display: flex;
    border: 1px solid var(--color-border);
    padding: 1px;
    border-radius: 3px;
    position: relative;
    z-index: $z-index-top;
    width: 349px;
    position: absolute;
    top: 95px;
    right: 28px;
    box-sizing: border-box;
    background-color: #ffffff;
    text-align: left;

    @include bp-down($sm) {
        width: 100%;
        position: fixed;
        border-top-right-radius: 10px;
        border-top-left-radius: 10px;
        bottom: 0px;
        text-align: center;
        padding-bottom: 10px;
        top: auto;
        right: auto;
        align-items: center;
        justify-content: center;
    }

    .link {
        color: $color-fieldkit-primary;
        font-size: 12px;
        padding-bottom: 15px;
        margin-bottom: 10pxs;
    }
}
.detail-title {
    font-family: $font-family-bold;
    font-size: 18px;
    margin-top: 15px;
    margin-bottom: 5px;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    max-height: 1em;
    width: 240px;
}
.detail-container {
    width: 75%;
    margin-bottom: 10px;
}
.detail-description {
    font-family: $font-family-light;
    font-size: 14px;
    max-height: 35px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    max-height: 35px;
    overflow: hidden;
    margin-bottom: 10px;
    margin-right: 10px;
}

.details {
    display: flex;
    border-radius: 2px;
    border: solid 1px var(--color-border);
    background-color: white;
}

.photo-container {
    width: 75px;
    height: 75px;
    margin: 10px;
    float: left;

    @include bp-down($sm) {
        display: none;
    }
}
.container-map {
    width: 100%;
    height: calc(100% - 33px);
    margin-top: 0;
    @include position(absolute, 66px null null 0);

    @include bp-down($sm) {
        top: 54px;
        height: calc(100% - 54px);
    }
}
.map-legend {
    display: flex;
    flex-direction: column;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    position: relative;
    z-index: $z-index-top;
    position: absolute;
    bottom: 37px;
    left: 65px;
    box-sizing: border-box;
    background-color: #ffffff;
    text-align: left;
    padding: 15px;
    padding-right: 30px;
    font-family: $font-family-medium;

    h4 {
        margin: 0 0 0.5em 0;
        font-family: $font-family-bold;
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

    @include bp-down($sm) {
        display: none;
    }
}

::v-deep .station-hover-summary {
    width: 359px;
    top: 122px;
    left: 300px;
}
</style>
