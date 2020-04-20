<template>
    <div>
        <SidebarNav
            :isAuthenticated="isAuthenticated"
            viewing="data"
            :stations="stations"
            :projects="projects"
            @showStation="showStation"
        />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <div id="data-view-background" class="main-panel" v-show="isAuthenticated">
            <div id="data-container">
                <router-link :to="{ name: 'stations' }">
                    <div class="map-link">
                        <span class="small-arrow">&lt;</span>
                        Stations Map
                    </div>
                </router-link>
                <div>
                    <div id="project-name">
                        {{ projects[0] ? projects[0].name : "" }}
                    </div>
                    <div class="block-label">Data visualization</div>
                </div>
                <DataChartControl
                    ref="dataChartControl"
                    :treeSelectOptions="treeSelectOptions"
                    @stationAdded="getInitialStationData"
                    @stationChanged="onStationChange"
                    @switchedSensor="onSensorSwitch"
                    @timeChanged="onTimeChange"
                />
                <!-- <div id="lower-container">
                    <NotesList :station="station" :selectedSensor="selectedSensor" :isAuthenticated="isAuthenticated" />
                    <SensorSummary
                        ref="sensorSummary"
                        :sensors="sensors"
                        :selectedSensor="selectedSensor"
                        :stationData="stationData"
                        :timeRange="timeRange"
                        :labels="labels"
                    />
                </div> -->
            </div>
        </div>
        <div v-if="failedAuth" class="no-auth-message">
            <p>
                Please
                <router-link :to="{ name: 'login' }" class="show-link">
                    log in
                </router-link>
                to view data.
            </p>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as d3 from "d3";
import FKApi from "@/api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import DataChartControl from "../components/DataChartControl";
// import NotesList from "../components/NotesList";
// import SensorSummary from "../components/SensorSummary";
// import * as tempStations from "../assets/ancientGoose.json";
import * as utils from "../utilities";

export default {
    name: "DataView",
    components: {
        HeaderBar,
        SidebarNav,
        DataChartControl,
        // NotesList,
        // SensorSummary,
    },
    props: [],
    data: () => {
        return {
            user: {},
            stationId: null,
            stationData: {},
            stations: [],
            projects: [],
            allSensors: [],
            selectedSensor: null,
            isAuthenticated: false,
            failedAuth: false,
            timeRange: null,
            treeSelectOptions: [],
        };
    },
    async beforeCreate() {
        const dataView = this;
        window.onpopstate = function(event) {
            // Note: event.state.key changes
            dataView.componentKey = event.state ? event.state.key : 0;
            dataView.$refs.dataChartControl.refresh(dataView.stationData);
        };
        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                this.api.getStations().then(s => {
                    this.stations = s.stations;
                    this.api
                        .getModulesMeta()
                        .then(this.processModulesMeta)
                        .then(this.initDataChartControl)
                        .then(this.initTreeSelect);
                });
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
            })
            .catch(() => {
                this.failedAuth = true;
            });
    },
    methods: {
        initDataChartControl() {
            if (this.$refs.dataChartControl) {
                this.$refs.dataChartControl.initialize();
            }
        },

        getInitialStationData(id, chartId) {
            this.api
                .getStation(id)
                .then(station => {
                    const deviceId = station.device_id;
                    this.stationData[deviceId] = {};
                    this.stationData[deviceId].station = station;
                    this.api
                        .getStationDataSummaryByDeviceId(deviceId)
                        .then(result => {
                            this.handleInitialDataSummary(result, deviceId, chartId);
                        })
                        .catch(() => {
                            this.$refs.dataChartControl.noInitialDataFound(station, chartId);
                        });
                })
                .catch(() => {
                    this.$refs.dataChartControl.noStationFound();
                });
        },

        handleInitialDataSummary(result, deviceId, chartId, reset) {
            if (!result) {
                this.stationData[deviceId] = {
                    sensors: [],
                    data: [],
                };
                this.$refs.dataChartControl.initializeChart(this.stationData, deviceId, chartId);
                return;
            }
            const processedData = this.processData(result);
            let processed = processedData.data;
            //sort data by date
            processed.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });
            const origEnd = processed[processed.length - 1].date;
            // convert to midnight of the given day
            const endDate = new Date(origEnd.getFullYear(), origEnd.getMonth(), origEnd.getDate(), 23, 59, 59);

            this.stationData[deviceId].timeRange = [processed[0].date, endDate];
            this.stationData[deviceId].data = processed;
            this.stationData[deviceId].sensors = processedData.sensors;
            if (processed.length > 0) {
                // resize div to fit sections
                // document.getElementById("lower-container").style["min-width"] = "1100px";
            }
            if (reset) {
                this.$refs.dataChartControl.resetChart(this.stationData, deviceId, chartId);
            } else {
                this.$refs.dataChartControl.initializeChart(this.stationData, deviceId, chartId);
            }
        },

        async fetchSummary(deviceId, start, end) {
            return this.api.getStationDataSummaryByDeviceId(deviceId, start, end);
        },

        processModulesMeta(meta) {
            // temp color function for sensors without ranges
            const black = function() {
                return "#000000";
            };

            meta.forEach(m => {
                // m.key
                m.sensors.forEach(s => {
                    let colors;
                    if (s.ranges.length > 0) {
                        colors = d3
                            .scaleSequential()
                            .domain([s.ranges[0].minimum, s.ranges[0].maximum])
                            .interpolator(d3.interpolatePlasma);
                    } else {
                        colors = d3
                            .scaleSequential()
                            .domain([0, 1])
                            .interpolator(black);
                    }
                    this.allSensors.push({
                        colorScale: colors,
                        unit: s.unit_of_measure,
                        key: s.key,
                        firmwareKey: s.firmware_key,
                        name: null,
                    });
                });
            });
            return;
        },

        processData(result) {
            let processed = [];
            const resultSensors = _.flatten(_.map(result.modules, "sensors"));
            const sensors = _.intersectionBy(this.allSensors, resultSensors, s => {
                return s.key;
            });
            result.data.forEach(d => {
                // Only including ones with sensor readings
                if (Object.keys(d.d).length > 0) {
                    d.d.date = new Date(d.time);
                    // filter out ones with dates before 2018
                    if (d.d.date.getFullYear() > 2018) {
                        processed.push(d.d);
                    }
                }
            });
            return { data: processed, sensors: sensors };
        },

        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },

        initTreeSelect() {
            this.treeSelectOptions = [];
            let counterId = 0;
            console.log("stations", this.stations.length);
            this.stations.forEach(s => {
                let modules = [];
                if (s.status_json.moduleObjects) {
                    s.status_json.moduleObjects.forEach(m => {
                        console.log("station moduleObject", m.name);
                        let sensors = [];
                        let addModule = true;
                        m.sensorObjects.forEach(sensor => {
                            console.log("station sensorObject", sensor.name);
                            let dataViewSensor = this.allSensors.find(sr => {
                                return sr.firmwareKey == sensor.name;
                            });
                            if (dataViewSensor) {
                                console.log("found dataViewSensor", dataViewSensor);
                                const sensorLabel = this.getSensorName(m, sensor);
                                console.log("and label", sensorLabel);
                                sensors.push({
                                    id: s.name + dataViewSensor.key,
                                    label: sensorLabel,
                                    customLabel: s.name + " : " + sensorLabel,
                                    key: dataViewSensor.key,
                                    stationId: s.id,
                                });
                            } else {
                                // don't add module if sensor wasn't found
                                addModule = false;
                            }
                        });
                        if (addModule) {
                            counterId += 1;
                            modules.push({
                                id: counterId,
                                label: this.getModuleName(m),
                                customLabel: s.name + " : " + this.getModuleName(m),
                                stationId: s.id,
                                children: sensors,
                            });
                        }
                    });
                    console.log("modules added", modules, "for", s.name);
                    counterId += 1;
                    this.treeSelectOptions.push({
                        id: counterId,
                        label: s.name,
                        customLabel: s.name,
                        stationId: s.id,
                        children: modules,
                    });
                }
            });
        },

        onSensorSwitch(sensor) {
            this.selectedSensor = sensor;
        },

        onStationChange(stationId, chart) {
            this.api.getStation(stationId).then(station => {
                const deviceId = station.device_id;
                if (this.stationData[deviceId]) {
                    // use this station's data if we already have it
                    this.$refs.dataChartControl.resetChart(this.stationData, deviceId, chart.id);
                    this.fetchSummary(deviceId, chart.start.getTime(), chart.end.getTime()).then(result => {
                        const processedData = this.processData(result);
                        this.$refs.dataChartControl.updateChartData(processedData.data, chart.id, station);
                    });
                } else {
                    this.stationData[deviceId] = {};
                    this.stationData[deviceId].station = station;
                    // otherwise, need to fetch data for initializing ranges, etc
                    this.api
                        .getStationDataSummaryByDeviceId(deviceId)
                        .then(result => {
                            const reset = true;
                            this.handleInitialDataSummary(result, deviceId, chart.id, reset);
                            // and then get the correct time window summary
                            this.fetchSummary(deviceId, chart.start.getTime(), chart.end.getTime()).then(result => {
                                const processedData = this.processData(result);
                                this.$refs.dataChartControl.updateChartData(processedData.data, chart.id, station);
                            });
                        })
                        .catch(() => {
                            this.$refs.dataChartControl.noDataFound(station, chart.id);
                        });
                }
            });
        },

        onTimeChange(range, chart) {
            const start = range.start.getTime();
            const end = range.end.getTime();
            // TODO: perhaps don't rely on parent chart having this id:
            // if (chart.id == "chart-1") {
            //     // change time range for SensorSummary
            //     this.timeRange = range;
            // }
            this.fetchSummary(chart.station.device_id, start, end).then(result => {
                const processedData = this.processData(result);
                this.$refs.dataChartControl.updateChartData(processedData.data, chart.id);
            });
        },

        getModuleName(module) {
            const newName = utils.convertOldFirmwareResponse(module);
            return this.$t(newName + ".name");
        },

        getSensorName(module, sensor) {
            const newName = utils.convertOldFirmwareResponse(module);
            return this.$t(newName + ".sensors." + sensor.name);
        },

        showStation(station) {
            if (station) {
                this.$router.push({ name: "viewData", query: { stationId: station.id } });
                this.$refs.dataChartControl.refresh();
            }
            // this.selectedSensor = null;
            // document.getElementById("lower-container").style["min-width"] = "700px";
        },
    },
};
</script>

<style scoped>
#data-view-background {
    color: rgb(41, 61, 81);
    background-color: rgb(252, 252, 252);
}
#data-container {
    margin: 20px;
    overflow: scroll;
}
#lower-container {
    clear: both;
}
.no-auth-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 40px;
}
.show-link {
    text-decoration: underline;
}
.small-arrow {
    font-size: 9px;
    vertical-align: middle;
}
.map-link {
    margin: 10px 0;
    font-size: 13px;
}
#project-name {
    font-size: 24px;
    line-height: 24px;
    font-weight: bold;
    margin: 20px 10px 0 0;
    display: inline-block;
}
.small-label {
    font-weight: normal;
    font-size: 12px;
    display: inline-block;
}
.block-label {
    margin-top: 5px;
    font-weight: normal;
    font-size: 18px;
    display: block;
}
</style>
