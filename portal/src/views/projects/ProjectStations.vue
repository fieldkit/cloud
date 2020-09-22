<template>
    <div class="stations-container">
        <StationPickerModal :stations="userStations" @add="onAddStation" @close="onCloseStationPicker" v-if="addingStation" />
        <div class="section-heading stations-heading">
            FieldKit Stations
            <div class="add-station" v-on:click="showStationPicker" v-if="admin">
                <img src="@/assets/icon-plus-round.svg" class="add-station-btn" />
                Add Station
            </div>
        </div>
        <div class="section-body">
            <div class="stations-panel" v-show="showStationsPanel">
                <div v-if="projectStations.length == 0" class="project-stations-no-stations">
                    <h3>No Stations</h3>
                    <p>
                        Add a station to this project to include its recent data and activities.
                    </p>
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
                            <img class="notes" v-on:click="openNotes(station)" src="@/assets/icon-field-notes.svg" />
                            <div class="remove" v-on:click="removeStation(station)">Delete</div>
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
                <StationsMap @show-summary="showSummary" :mapped="mappedProject" :layoutChanges="layoutChanges" />
                <StationSummary
                    v-if="activeStation"
                    :station="activeStation"
                    :readings="false"
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
import * as utils from "../../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";
import StationSummary from "@/views/shared/StationSummary.vue";
import StationsMap from "@/views/shared/StationsMap.vue";
import StationPickerModal from "@/views/shared/StationPickerModal.vue";
import TinyStation from "@/views/shared/TinyStation.vue";
import PaginationControls from "@/views/shared/PaginationControls.vue";

export default Vue.extend({
    name: "ProjectStations",
    components: {
        StationSummary,
        StationPickerModal,
        StationsMap,
        PaginationControls,
        TinyStation,
    },
    data: () => {
        return {
            activeStationId: null,
            layoutChanges: 0,
            showStationsPanel: true,
            addingStation: false,
            page: 0,
            pageSize: 4,
        };
    },
    props: {
        project: { required: true },
        admin: { required: true },
        userStations: { required: true },
    },
    watch: {
        project() {
            this.activeStationId = null;
        },
    },
    computed: {
        projectStations(this: any) {
            return this.$getters.projectsById[this.project.id].stations;
        },
        mappedProject(this: any) {
            return this.$getters.projectsById[this.project.id].mapped;
        },
        activeStation(this: any) {
            if (this.activeStationId) {
                return this.$getters.stationsById[this.activeStationId];
            }
            return null;
        },
        visibleStations(this: any) {
            if (!this.projectStations) {
                return [];
            }
            const start = this.page * this.pageSize;
            const end = start + this.pageSize;
            return this.projectStations.slice(start, end);
        },
        totalPages() {
            if (this.projectStations) {
                return Math.ceil(this.projectStations.length / this.pageSize);
            }
            return 0;
        },
    },
    methods: {
        showStation(this: any, station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        showStationPicker(this: any) {
            this.addingStation = true;
        },
        onAddStation(this: any, stationId) {
            this.addingStation = false;
            const payload = {
                projectId: this.project.id,
                stationId: stationId,
            };
            return this.$store.dispatch(ActionTypes.STATION_PROJECT_ADD, payload);
        },
        onCloseStationPicker() {
            this.addingStation = false;
        },
        showSummary(station) {
            console.log("showSummay", station);
            this.activeStationId = station.id;
        },
        removeStation(this: any, station) {
            console.log("remove", station);
            if (!window.confirm("Are you sure you want to remove this station?")) {
                return;
            }
            const payload = {
                projectId: this.project.id,
                stationId: station.id,
            };
            return this.$store.dispatch(ActionTypes.STATION_PROJECT_REMOVE, payload);
        },
        openNotes(this: any, station) {
            return this.$router.push({
                name: "viewProjectStationNotes",
                params: {
                    projectId: this.project.id,
                    stationId: station.id,
                },
            });
        },
        onCloseSummary() {
            this.activeStationId = null;
        },
        toggleStationsPanel() {
            this.layoutChanges++;
            this.showStationsPanel = !this.showStationsPanel;
            console.log("toggle", this.showStationsPanel, this.layoutChanges);
        },
        onNewPage(this: any, page: number) {
            this.page = page;
        },
    },
});
</script>


<style scoped lang="scss">
@import '../../scss/mixins';

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
    border-bottom: 1px solid #d8dce0;
}
.stations-heading {
    display: flex;
    align-items: center;
    flex-direction: row;
}
.section-body {
    display: flex;
    flex-direction: row;
    height: 352px;
}
.stations-container {
    display: flex;
    flex-direction: column;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.add-station {
    margin-left: auto;
    font-size: 14px;
    margin-right: 1em;
    cursor: pointer;
    display: flex;
}
.add-station-btn {
    width: 18px;
    vertical-align: text-top;
    margin-right: 7px;
}
.last-seen {
    font-size: 12px;
    font-weight: 600;
    color: #6a6d71;
}
.stations-container {
    margin: 25px 0 0 0;
    border-radius: 1px;
    border: solid 1px #d8dce0;
    background-color: #ffffff;
}
.stations-panel {
    transition: width 0.5s;
    flex: 1;
    display: flex;
    flex-direction: column;
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
}
.station-box {
    height: 38px;
    margin: 20px auto;
    padding: 1em;
    border: 1px solid #d8dce0;
    transition: opacity 0.25s;
}
.project-stations-no-stations {
    width: 80%;
    text-align: center;
    margin: auto;
    padding-top: 20px;
}
.project-stations-no-stations h1 {
}
.project-stations-no-stations p {
}

.station-links {
    margin-left: auto;
    flex: 0 0 61px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    border-left: 1px solid #d8dce0;
    text-align: center;
}

.station-links .remove {
    cursor: pointer;
    font-size: 12px;
}

.station-links .notes {
    cursor: pointer;
    font-size: 14px;
}

.pagination {
    margin-top: auto;
    padding-bottom: 1em;
}

::v-deep .station-hover-summary {
    width: 359px;
    top: 30px;
    left: 122px;
}
</style>
