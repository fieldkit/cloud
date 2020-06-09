<template>
    <div class="stations-container" :style="{ width: mapContainerSize.outerWidth }">
        <div class="section-heading stations-heading">
            FieldKit Stations
            <div class="add-station" v-on:click="showStationSelect" v-if="admin && !addingStation">
                <img src="../assets/add.png" class="add-station-btn" />
                Add Station
            </div>
            <div class="station-dropdown" v-if="addingStation">
                Add a station:
                <select v-model="stationOption" v-on:change="stationSelected">
                    <option v-for="station in userStations" v-bind:value="station.id" v-bind:key="station.id">
                        {{ station.name }}
                    </option>
                </select>
            </div>
        </div>
        <div class="space"></div>
        <div id="stations-list" :style="{ width: listSize.width, height: listSize.height }">
            <div class="toggle-icon-container" v-on:click="toggleStations">
                <img v-if="showStationsList" alt="Collapse List" src="../assets/tab-collapse.png" class="toggle-icon" />
                <img v-if="!showStationsList" alt="Expand List" src="../assets/tab-expand.png" class="toggle-icon" />
            </div>
            <div v-for="station in projectStations" v-bind:key="station.id">
                <div class="station-box" :style="{ width: listSize.boxWidth }">
                    <div class="delete-link">
                        <img alt="Delete" src="../assets/Delete.png" :data-id="station.id" v-on:click="deleteStation" />
                    </div>
                    <span class="station-name" v-on:click="showStation(station)">
                        {{ station.name }}
                    </span>
                    <div class="last-seen">Last seen {{ getUpdatedDate(station) }}</div>
                </div>
            </div>
        </div>
        <div id="stations-map-container" :style="{ width: mapContainerSize.width, height: mapContainerSize.height }">
            <StationsMap
                :mapSize="mapSize"
                :stations="projectStations"
                @mapReady="onMapReady"
                @showSummary="showSummary"
                ref="stationsMap"
            />
            <StationSummary
                v-show="activeStation"
                :station="activeStation"
                :compact="true"
                :summarySize="summarySize"
                ref="stationSummary"
            />
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as utils from "../utilities";
import FKApi from "../api/api";
import StationSummary from "../components/StationSummary";
import StationsMap from "../components/StationsMap";

export default {
    name: "ProjectStations",
    components: {
        StationSummary,
        StationsMap,
    },
    data: () => {
        return {
            projectStations: [],
            activeStation: null,
            following: false,
            showStationsList: true,
            addingStation: false,
            stationOption: "",
            summarySize: {
                width: "359px",
                top: "-300px",
                left: "122px",
                constrainTop: "230px",
            },
            mapSize: {
                width: "inherit",
                height: "inherit",
                position: "relative",
            },
        };
    },
    props: ["project", "admin", "mapContainerSize", "listSize", "userStations"],
    async beforeCreate() {
        this.api = new FKApi();
    },
    mounted() {
        this.fetchStations();
    },
    methods: {
        onMapReady(map) {
            this.map = map;
        },
        fetchStations() {
            this.api.getStationsByProject(this.project.id).then(result => {
                this.projectStations = result.stations;
                let modules = [];
                if (this.projectStations) {
                    this.projectStations.forEach((s, i) => {
                        if (i == 0 && s.location && s.location.latitude && this.map) {
                            this.map.setCenter({
                                lat: parseFloat(s.location.latitude),
                                lng: parseFloat(s.location.longitude),
                            });
                        }
                        s.modules.forEach(m => {
                            if (!m.internal) {
                                modules.push(this.getModuleImg(m));
                            }
                        });
                    });
                    modules = _.uniq(modules);
                    this.$emit("loaded", { modules: modules, projectStations: this.projectStations });
                }
            });
        },
        getUpdatedDate(station) {
            return utils.getUpdatedDate(station);
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        showStationSelect() {
            this.addingStation = true;
        },
        stationSelected() {
            const params = {
                projectId: this.project.id,
                stationId: this.stationOption,
            };
            this.api.addStationToProject(params).then(() => {
                this.fetchStations();
            });
        },
        deleteStation(event) {
            const stationId = event.target.getAttribute("data-id");
            if (window.confirm("Are you sure you want to remove this station?")) {
                const params = {
                    projectId: this.project.id,
                    stationId: stationId,
                };
                this.api.removeStationFromProject(params).then(() => {
                    this.fetchStations();
                });
            }
        },
        getModuleImg(module) {
            let imgPath = require.context("../assets/modules-lg/", false, /\.png$/);
            let img = utils.getModuleImg(module);
            return imgPath("./" + img);
        },
        toggleStations() {
            let stationsMap = document.getElementById("stations-map-container");
            this.showStationsList = !this.showStationsList;
            if (this.showStationsList) {
                document.getElementById("stations-list").style.width = this.listSize.width;
                stationsMap.style.width = this.mapContainerSize.width;
                stationsMap.style["margin-left"] = "0";
                this.map.resize();
                document.getElementById("stations-map-container").style.transition = "width 0.5s";
                let boxes = document.getElementsByClassName("station-box");
                Array.from(boxes).forEach(b => {
                    b.style.opacity = 1;
                });
            } else {
                let boxes = document.getElementsByClassName("station-box");
                Array.from(boxes).forEach(b => {
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
            this.activeStation = station;
            this.$refs.stationSummary.viewSummary();
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
    float: left;
    margin: 0 0 35px 0;
}
.stations-heading {
    width: 100%;
    margin: 25px 0 25px 25px;
}
.stations-container .space {
    margin: 0;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.station-dropdown,
.add-station {
    float: right;
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
    float: left;
    margin: 22px 0 0 0;
    border: 1px solid #d8dce0;
}
#stations-list {
    float: left;
    transition: width 0.5s;
}
#stations-map-container {
    float: left;
    transition: width 0.5s;
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
</style>
