<template>
    <div class="stations-container">
        <StationPickerModal
            :stations="userStations"
            :filter="pickFilter"
            :actionType="StationPickerActionType.add"
            :title="$t('project.stations.add.title')"
            :actionText="$t('project.stations.add.cta')"
            @add="onStationAdd"
            @close="onCloseAddStationModal"
            v-if="addingStation"
        />
        <StationPickerModal
            :stations="projectStations"
            :actionType="StationPickerActionType.remove"
            :title="$t('project.stations.edit.title')"
            :actionText="$t('project.stations.edit.cta')"
            @remove="onStationRemove"
            @close="onCloseEditStationModal"
            v-if="editingStation"
        />
        <StationOrSensor class="section-heading stations-heading">
            <div class="stations-cta-container" v-if="admin">
                <div class="stations-cta" v-on:click="showAddStationPicker">
                    <i class="icon icon-plus-round"></i>
                    {{ $t("project.stations.add.trigger") }}
                </div>
                <div class="stations-cta" v-on:click="showEditStationPicker">
                    <i class="icon icon-minus-round"></i>
                    {{ $t("project.stations.edit.trigger") }}
                </div>
            </div>
        </StationOrSensor>
        <div class="section-body">
            <div class="stations-panel" v-show="showStationsPanel">
                <div v-if="projectStations.length == 0" class="project-stations-no-stations">
                    <h3>{{ $t("project.stations.none.title") }}</h3>
                    <p>{{ $t("project.stations.none.text") }}</p>
                </div>
                <div v-if="projectStations.length > 0" class="stations">
                    <TinyStation
                        v-for="station in visibleStations"
                        v-bind:key="station.id"
                        :station="station"
                        :narrow="true"
                        @selected="showSummary(station)"
                    >
                        <div class="station-links">
                            <i v-on:click="openNotes(station)" class="icon icon-field-notes"></i>
                        </div>
                    </TinyStation>
                </div>
                <PaginationControls :page="page" :totalPages="totalPages" @new-page="onNewPage" />
            </div>
            <div class="toggle-icon-container" v-on:click="toggleStationsPanel">
                <img v-if="showStationsPanel" alt="Collapse List" src="@/assets/icon-tab-collapse.svg" class="toggle-expand" />
                <img v-if="!showStationsPanel" alt="Expand List" src="@/assets/icon-tab-expand.svg" class="toggle-icon" />
            </div>
            <div class="project-stations-map-container">
                <router-link :to="{ name: 'viewProjectBigMap' }" class="link">
                    <div class="map-expand">
                        <img alt="Location" src="@/assets/icon-expand-map.svg" class="icon" />
                    </div>
                </router-link>

                <StationsMap
                    @show-summary="showSummary"
                    :mapped="mappedProject"
                    :layoutChanges="layoutChanges"
                    :showStations="project.showStations"
                    :mapBounds="mapBounds"
                />
                <StationSummary
                    v-if="activeStation"
                    :station="activeStation"
                    :readings="false"
                    :exploreContext="exploreContext"
                    @close="onCloseSummary"
                    v-bind:key="activeStation.id"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";

import { ExploreContext } from "@/views/viz/common";

import { BoundingRectangle, DisplayProject, DisplayStation, MappedStations } from "@/store";
import * as ActionTypes from "@/store/actions";

import StationSummary from "@/views/shared/StationSummary.vue";
import StationsMap from "@/views/shared/StationsMap.vue";
import StationPickerModal from "@/views/shared/StationPickerModal.vue";
import TinyStation from "@/views/shared/TinyStation.vue";
import PaginationControls from "@/views/shared/PaginationControls.vue";
import { StationPickerActionType } from "@/views/shared/StationPicker.vue";
import StationOrSensor from "@/views/shared/partners/StationOrSensor.vue";

export default Vue.extend({
    name: "ProjectStations",
    components: {
        StationSummary,
        StationPickerModal,
        StationsMap,
        PaginationControls,
        TinyStation,
        StationOrSensor,
    },
    data(): {
        activeStationId: number | null;
        layoutChanges: number;
        showStationsPanel: boolean;
        addingStation: boolean;
        editingStation: boolean;
        page: number;
        pageSize: number;
        StationPickerActionType: any;
    } {
        return {
            activeStationId: null,
            layoutChanges: 0,
            showStationsPanel: true,
            addingStation: false,
            editingStation: false,
            page: 0,
            pageSize: 4,
            StationPickerActionType: StationPickerActionType,
        };
    },
    props: {
        project: {
            type: Object as () => DisplayProject,
            required: true,
        },
        admin: {
            type: Boolean,
            required: true,
        },
        userStations: {
            type: Array as () => DisplayStation[],
            required: true,
        },
    },
    watch: {
        project(): void {
            this.activeStationId = null;
        },
    },
    computed: {
        projectStations(): DisplayStation[] {
            return this.$getters.projectsById[this.project.id].stations;
        },
        mappedProject(): MappedStations | null {
            console.log("radoi project mapped", this.$getters.projectsById[this.project.id].mapped);
            return this.$getters.projectsById[this.project.id].mapped;
        },
        activeStation(): DisplayStation | null {
            if (this.activeStationId) {
                return this.$getters.stationsById[this.activeStationId];
            }
            return null;
        },
        visibleStations(): DisplayStation[] {
            if (!this.projectStations) {
                return [];
            }
            const start = this.page * this.pageSize;
            const end = start + this.pageSize;
            return this.projectStations.slice(start, end);
        },
        totalPages(): number {
            if (this.projectStations) {
                return Math.ceil(this.projectStations.length / this.pageSize);
            }
            return 0;
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
    },
    methods: {
        pickFilter(station: DisplayStation): boolean {
            const excluding = _.map(this.projectStations, (s) => s.id);
            return excluding.indexOf(station.id) < 0;
        },
        showStation(station: DisplayStation): Promise<any> {
            // All parameters are strings.
            return this.$router.push({ name: "mapStation", params: { id: String(station.id) } });
        },
        showAddStationPicker(): void {
            this.addingStation = true;
        },
        showEditStationPicker(): void {
            this.editingStation = true;
        },
        onStationAdd(stationIds): void {
            this.addingStation = false;
            stationIds.forEach((stationId) => {
                const payload = {
                    projectId: this.project.id,
                    stationId: stationId,
                };
                this.$store.dispatch(ActionTypes.STATION_PROJECT_ADD, payload);
            });
        },
        onStationRemove(stationIds): void {
            this.editingStation = false;
            stationIds.forEach((stationId) => {
                const payload = {
                    projectId: this.project.id,
                    stationId: stationId,
                };
                this.$store.dispatch(ActionTypes.STATION_PROJECT_REMOVE, payload);
            });
        },
        onCloseAddStationModal(): void {
            this.addingStation = false;
        },
        onCloseEditStationModal(): void {
            this.editingStation = false;
        },
        showSummary(station: DisplayStation): void {
            console.log("showSummay", station);
            this.activeStationId = station.id;
        },
        removeStation(this: any, station: DisplayStation): Promise<any> {
            console.log("remove", station);
            if (!window.confirm("Are you sure you want to remove this station?")) {
                return Promise.resolve();
            }
            const payload = {
                projectId: this.project.id,
                stationId: station.id,
            };
            return this.$store.dispatch(ActionTypes.STATION_PROJECT_REMOVE, payload);
        },
        openNotes(this: any, station: DisplayStation): Promise<any> {
            return this.$router.push({
                name: "viewProjectStationNotes",
                params: {
                    projectId: this.project.id,
                    stationId: station.id,
                },
            });
        },
        onCloseSummary(): void {
            this.activeStationId = null;
        },
        toggleStationsPanel(): void {
            this.layoutChanges++;
            this.showStationsPanel = !this.showStationsPanel;
            console.log("toggle", this.showStationsPanel, this.layoutChanges);
        },
        onNewPage(page: number): void {
            this.page = page;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.toggle-icon-container {
    float: right;
    margin: 16px -38px 0 0;
    position: relative;
    z-index: 2;
    cursor: pointer;
}
.toggle-expand {
    transform: translateX(-4px);
}
.section {
    float: left;
}
.section-heading {
    font-size: 20px;
    font-weight: 500;
    padding-top: 1em;
    padding-bottom: 1em;
    padding-left: 1em;
    border-bottom: 1px solid var(--color-border);

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
}
.stations-heading {
    display: flex;
    align-items: center;
    flex-direction: row;
}
.section-body {
    display: flex;
    flex-direction: row;
    height: 420px;
}
.stations-container {
    display: flex;
    flex-direction: column;
}
.station-name {
    font-size: 14px;
    cursor: pointer;
}
.stations-cta-container {
    margin-left: auto;
    font-size: 14px;
    margin-right: 1em;
    @include flex(center);
}
.stations-cta {
    cursor: pointer;
    @include flex(center);

    &:not(:last-of-type) {
        margin-right: 35px;

        @include bp-down($xs) {
            margin-right: 15px;
        }
    }

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }

    .icon {
        margin-right: 7px;
        margin-top: -3px;
    }
}
.last-seen {
    font-size: 12px;
    font-family: var(--font-family-bold);
    color: #6a6d71;
}
.stations-container {
    margin: 25px 0 0 0;
    border-radius: 1px;
    border: solid 1px var(--color-border);
    background-color: #ffffff;
    overflow: hidden;
}
.stations-panel {
    transition: width 0.5s;
    flex: 1;
    display: flex;
    flex-direction: column;

    @include bp-down($xs) {
        flex-basis: 85%;
    }
}
.stations-panel .stations {
    padding: 20px 25px 0;

    @include bp-down($xs) {
        padding: 20px 10px;
    }
}
.stations-panel .stations > * {
    margin-bottom: 1em;
}
.project-stations-map-container {
    transition: width 0.5s;
    position: relative;
    flex: 2;
    height: 100%;
}
.station-box {
    height: 38px;
    margin: 20px auto;
    padding: 1em;
    border: 1px solid var(--color-border);
    transition: opacity 0.25s;
}
.project-stations-no-stations {
    width: 80%;
    text-align: center;
    margin: auto;
    padding-top: 20px;
}

.project-stations-no-stations p {
    font-family: $font-family-light;
}

.station-links {
    margin-left: auto;
    flex: 0 0 61px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    border-left: 1px solid var(--color-border);
    text-align: center;
}

.station-links .remove {
    cursor: pointer;
    font-size: 12px;
}

.pagination {
    margin-top: auto;
    padding-bottom: 1em;
}

::v-deep .station-hover-summary {
    width: 359px;
    top: 20px;
    left: 122px;
}

.icon-field-notes {
    font-size: 20px;
    cursor: pointer;
    margin: auto;

    &:before {
        color: var(--color-primary);
    }

    body.floodnet & {
        &:before {
            color: var(--color-dark);
        }
    }
}
.map-expand {
    background-color: #ffffff;
    z-index: 50;
    position: absolute;
    right: 10px;
    top: 10px;
    display: inline-block;
    height: 38px;
    border: 1px solid var(--color-border);
    padding: 1px;
    border-radius: 3px;
    box-sizing: border-box;

    img {
        margin: 10px;
    }
}
</style>
