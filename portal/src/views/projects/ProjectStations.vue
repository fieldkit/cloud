<template>
    <div class="stations-container">
        <div class="section-heading stations-heading">
            FieldKit Stations
            <div class="add-station" v-on:click="showStationSelect" v-if="admin && !addingStation">
                <img src="@/assets/add.png" class="add-station-btn" />
                Add Station
            </div>
            <div class="station-dropdown" v-if="addingStation">
                Add a station:
                <select v-model="selectedStationId" v-on:change="stationSelected">
                    <option v-for="station in userStations" v-bind:value="station.id" v-bind:key="station.id">
                        {{ station.name }}
                    </option>
                </select>
            </div>
        </div>
        <div class="section-body">
            <div id="stations-list">
                <div class="toggle-icon-container" v-on:click="toggleStations" v-if="false">
                    <img v-if="showStationsList" alt="Collapse List" src="@/assets/tab-collapse.png" class="toggle-icon" />
                    <img v-if="!showStationsList" alt="Expand List" src="@/assets/tab-expand.png" class="toggle-icon" />
                </div>
                <div v-if="projectStations.length == 0" class="project-stations-no-stations">
                    <h3>No Stations</h3>
                    <p>
                        Add a station to this project to include its recent data and activities.
                    </p>
                </div>
                <div v-if="projectStations.length > 0">
                    <div v-for="station in projectStations" v-bind:key="station.id">
                        <div class="station-box" :style="{ width: listSize.boxWidth }">
                            <div class="delete-link">
                                <img alt="Delete" src="@/assets/Delete.png" :data-id="station.id" v-on:click="deleteStation(station)" />
                            </div>
                            <span class="station-name" v-on:click="showStation(station)">
                                {{ station.name }}
                            </span>
                            <div class="last-seen">Last seen {{ station.updated | prettyDate }}</div>
                        </div>
                    </div>
                </div>
            </div>
            <div id="stations-map-container">
                <StationsMap @mapReady="onMapReady" @showSummary="showSummary" ref="stationsMap" :mapped="mappedProject" />
                <StationSummary
                    v-show="activeStation"
                    :station="activeStation"
                    :compact="true"
                    :summarySize="summarySize"
                    ref="stationSummary"
                />
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as utils from "../../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";
import StationSummary from "@/components/StationSummary";
import StationsMap from "@/components/StationsMap";

export default {
    name: "ProjectStations",
    components: {
        StationSummary,
        StationsMap,
    },
    data: () => {
        return {
            activeStation: null,
            following: false,
            showStationsList: true,
            addingStation: false,
            selectedStationId: null,
            summarySize: {
                width: "359px",
                top: "-300px",
                left: "122px",
                constrainTop: "230px",
            },
        };
    },
    props: {
        project: { required: true },
        admin: { required: true },
        mapContainerSize: { required: true },
        listSize: { required: true },
        userStations: { required: true },
    },
    computed: {
        projectStations() {
            return this.$getters.projectsById[this.project.id].stations;
        },
        mappedProject() {
            return this.$getters.projectsById[this.project.id].mapped;
        },
    },
    methods: {
        onMapReady(map) {
            console.log("map ready resize");
            this.map = map;
            this.map.resize();
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        showStationSelect() {
            this.addingStation = true;
        },
        async stationSelected() {
            const payload = {
                projectId: this.project.id,
                stationId: this.selectedStationId,
            };
            await this.$store.dispatch(ActionTypes.STATION_PROJECT_ADD, payload);
        },
        async deleteStation(station) {
            if (window.confirm("Are you sure you want to remove this station?")) {
                const payload = {
                    projectId: this.project.id,
                    stationId: station.Id,
                };
                await this.$store.dispatch(ActionTypes.STATION_PROJECT_REMOVE, payload);
            }
        },
        toggleStations() {
            const stationsMap = document.getElementById("stations-map-container");
            this.showStationsList = !this.showStationsList;
            if (this.showStationsList) {
                document.getElementById("stations-list").style.width = this.listSize.width;
                stationsMap.style.width = this.mapContainerSize.width;
                stationsMap.style["margin-left"] = "0";
                this.map.resize();
                document.getElementById("stations-map-container").style.transition = "width 0.5s";
                const boxes = document.getElementsByClassName("station-box");
                Array.from(boxes).forEach((b) => {
                    b.style.opacity = 1;
                });
            } else {
                const boxes = document.getElementsByClassName("station-box");
                Array.from(boxes).forEach((b) => {
                    b.style.opacity = 0;
                });
                document.getElementById("stations-list").style.width = "1px";
                stationsMap.addEventListener("transitionend", this.postExpandMap, true);
                stationsMap.style.width = this.mapContainerSize.outerWidth;
                stationsMap.style["margin-left"] = "-1px";
            }
        },
        postExpandMap() {
            this.map.resize();
            document.getElementById("stations-map-container").style.transition = "width 0s";
            document.getElementById("stations-map-container").removeEventListener("transitionend", this.postExpandMap, true);
        },
        showSummary(station) {
            console.log("this.activeStationId", station.id);
            this.activeStationId = station.id;
        },
    },
};
</script>

<style scoped>
.toggle-icon-container {
    float: right;
    margin: 16px -34px 0 0;
    position: relative;
    z-index: 2;
    cursor: pointer;
}
.section {
    float: left;
}
.section-heading {
    font-size: 20px;
    font-weight: 600;
    padding-top: 15px;
    padding-bottom: 15px;
    padding-left: 35px;
    border-bottom: 1px solid lightgray;
}
.stations-heading {
    display: flex;
    flex-direction: row;
}
.section-body {
    display: flex;
    flex-direction: row;
    height: 100%;
}
.stations-container {
    height: 400px;
    display: flex;
    flex-direction: column;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.station-dropdown,
.add-station {
    margin-left: auto;
    font-size: 14px;
    margin-right: 50px;
    cursor: pointer;
}
.add-station-btn {
    width: 18px;
    vertical-align: text-top;
}
.station-dropdown select {
    font-size: 16px;
    border: 1px solid lightgray;
    border-radius: 4px;
    padding: 2px 4px;
}
.last-seen {
    font-size: 12px;
    font-weight: 600;
    color: #6a6d71;
}
.space {
    width: 100%;
    float: left;
    margin: 30px 0 0 0;
    border-bottom: solid 1px #d8dce0;
}
.stations-container {
    margin: 22px 0 0 0;
    border: 2px solid #d8dce0;
    background-color: #ffffff;
}
#stations-list {
    transition: width 0.5s;
    flex: 1;
}
#stations-map-container {
    transition: width 0.5s;
    position: relative;
    flex: 2;
}
.station-box {
    height: 38px;
    margin: 20px auto;
    padding: 10px;
    border: 1px solid #d8dce0;
    transition: opacity 0.25s;
}
.delete-link {
    float: right;
    opacity: 0;
}
.delete-link:hover {
    opacity: 1;
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
</style>
