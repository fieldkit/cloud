<template>
    <div>
        <SidebarNav viewing="data" :stations="stations" :projects="projects" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <div id="data-view-background" class="main-panel" v-show="isAuthenticated">
            <div id="data-container">
                <router-link :to="{ name: 'stations' }">
                    <div class="map-link"><span class="small-arrow">&lt;</span> Stations Map</div>
                </router-link>
                <div>
                    <div id="station-name">
                        {{ this.station ? this.station.name : "" }}
                    </div>
                    <div class="small-label">{{ this.station ? getSyncedDate() : "" }}</div>
                    <div class="block-label">Data visualization</div>
                </div>
                <DataChartControl
                    ref="dataChartControl"
                    :combinedStationInfo="combinedStationInfo"
                    :station="station"
                    :noStation="noStation"
                    :labels="labels"
                    :totalTime="stationTimeRange"
                    @switchedSensor="onSensorSwitch"
                    @timeChanged="onTimeChange"
                    @chartTimeChanged="onChartTimeChange"
                />
                <div id="lower-container">
                    <NotesList
                        :station="station"
                        :selectedSensor="selectedSensor"
                        :isAuthenticated="isAuthenticated"
                    />
                    <SensorSummary
                        ref="sensorSummary"
                        :sensors="sensors"
                        :selectedSensor="selectedSensor"
                        :stationData="stationData"
                        :timeRange="timeRange"
                        :labels="labels"
                    />
                </div>
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
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import DataChartControl from "../components/DataChartControl";
import NotesList from "../components/NotesList";
import SensorSummary from "../components/SensorSummary";
import * as tempStations from "../assets/ancientGoose.json";
import * as utils from "../utilities";

export default {
    name: "DataView",
    components: {
        HeaderBar,
        SidebarNav,
        DataChartControl,
        NotesList,
        SensorSummary
    },
    props: ["stationParam", "id"],
    data: () => {
        return {
            user: {},
            noStation: false,
            station: null,
            stationId: null,
            stationTimeRange: [],
            stationData: [],
            stations: [],
            projects: [],
            sensors: [],
            combinedStationInfo: {},
            selectedSensor: null,
            isAuthenticated: false,
            failedAuth: false,
            timeRange: null,
            // temporary label system
            labels: {
                ph: "pH",
                do: "Dissolved Oxygen",
                ec: "Electrical Conductivity",
                tds: "Total Dissolved Solids",
                salinity: "Salinity",
                temp: "Temperature",
                humidity: "Humidity",
                temperature1: "Temperature",
                pressure: "Pressure",
                temperature2: "Temperature 2",
                rain: "Rain",
                windSpeed: "Wind Speed",
                windDir: "Wind Direction",
                windDirMv: "Wind Direction Raw ADC",
                windHrMaxSpeed: "Wind Max Speed (1 hour)",
                windHrMaxDir: "Wind Max Direction (1 hour)",
                wind10mMaxSpeed: "Wind Max Speed (10 min)",
                wind10mMaxDir: "Wind Max Direction (10 min)",
                wind2mAvgSpeed: "Wind Average Speed (2 min)",
                wind2mAvgDir: "Wind Average Direction (2 min)",
                rainThisHour: "Rain This Hour",
                rainPrevHour: "Rain Previous Hour",
                distance0: "Distance 0",
                distance1: "Distance 1",
                distance2: "Distance 2",
                calibration: "Calibration"
            }
        };
    },
    async beforeCreate() {
        const dataView = this;
        window.onpopstate = function(event) {
            // Note: event.state.key changes
            dataView.componentKey = event.state ? event.state.key : 0;
            if (dataView.id != dataView.$route.params.id) {
                dataView.getStationFromRoute();
            } else {
                dataView.$refs.dataChartControl.refresh(dataView.stationData);
            }
        };
        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                if (this.stationParam) {
                    this.station = this.stationParam;
                }
                if (this.id) {
                    this.stationId = this.id;
                }
                this.user = user;
                this.isAuthenticated = true;
                this.api.getStations().then(s => {
                    this.stations = s.stations;
                });
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
                this.api
                    .getModulesMeta()
                    .then(this.processModulesMeta)
                    .then(this.initiateDataRetrieval);
            })
            .catch(() => {
                this.failedAuth = true;
            });
    },
    methods: {
        async fetchSummary(start, end) {
            if (this.station) {
                return this.api.getStationDataSummaryByDeviceId(this.station.device_id, start, end);
            } else if (this.stationId) {
                // temporarily show Ancient Goose 81 (station 159) to anyone who views /dashboard/data/0
                if (this.stationId == 0) {
                    this.station = tempStations.stations[0];
                    return this.api.getStationDataSummaryByDeviceId(this.station.device_id, start, end);
                } else {
                    return this.api.getStation(this.stationId).then(station => {
                        this.station = station;
                        return this.api.getStationDataSummaryByDeviceId(this.station.device_id, start, end);
                    });
                }
            } else {
                this.noStation = true;
            }
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
                    this.sensors.push({
                        colorScale: colors,
                        unit: s.unit_of_measure,
                        key: s.key,
                        name: null
                    });
                });
            });
            return;
        },

        initiateDataRetrieval() {
            this.fetchSummary().then(result => {
                this.handleInitialDataSummary(result);
            });
        },

        handleInitialDataSummary(result) {
            this.stationData = [];
            if (!result) {
                this.combinedStationInfo = {
                    sensors: [],
                    stationData: this.stationData
                };
                return;
            }
            const processedData = this.processData(result);
            let sensors = processedData.sensors;
            let processed = processedData.data;
            //sort data by date
            processed.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });
            this.stationTimeRange[0] = processed[0].date;
            this.stationTimeRange[1] = processed[processed.length - 1].date;
            this.timeRange = { start: this.stationTimeRange[0], end: this.stationTimeRange[1] };
            this.stationData = processed;
            this.combinedStationInfo = { sensors: sensors, stationData: processed };
            if (processed.length > 0) {
                // resize div to fit sections
                document.getElementById("lower-container").style["min-width"] = "1100px";
            }
        },

        processData(result) {
            let processed = [];
            const resultSensors = _.flatten(_.map(result.modules, "sensors"));
            const sensors = _.intersectionBy(this.sensors, resultSensors, s => {
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

        handleSummary(result) {
            const processedData = this.processData(result);
            this.stationData = processedData.data;
            this.$refs.dataChartControl.updateAll(this.stationData);
        },

        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },

        onSensorSwitch(sensor) {
            this.selectedSensor = sensor;
        },

        onTimeChange(range) {
            this.timeRange = range;
            const start = this.timeRange.start.getTime();
            const end = this.timeRange.end.getTime();
            this.fetchSummary(start, end).then(result => {
                this.handleSummary(result);
            });
        },

        onChartTimeChange(range, chartId) {
            const start = range.start.getTime();
            const end = range.end.getTime();
            // TODO: perhaps don't rely on parent chart having this id:
            if (chartId == "chart-1") {
                // change time range for SensorSummary
                this.timeRange = range;
            }
            this.fetchSummary(start, end).then(result => {
                const processedData = this.processData(result);
                this.$refs.dataChartControl.updateChartData(processedData.data, chartId);
            });
        },

        getSyncedDate() {
            return "Last synced " + utils.getUpdatedDate(this.station);
        },

        getStationFromRoute() {
            this.stationId = this.$route.params.id;
            this.station = null;
            this.showStation();
        },

        showStation(station) {
            if (station) {
                this.$router.push({ name: "viewData", params: { id: station.id } });
                this.station = station;
            }
            this.$refs.dataChartControl.prepareNewStation();
            this.selectedSensor = null;
            this.noStation = false;
            document.getElementById("lower-container").style["min-width"] = "700px";
            this.initiateDataRetrieval();
        }
    }
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
#station-name {
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
    font-size: 14px;
    display: block;
}
</style>
