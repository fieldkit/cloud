<template>
    <div>
        <SidebarNav viewing="data" :stations="stations" :projects="projects" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <div id="data-view-background" class="main-panel" v-show="isAuthenticated">
            <div id="data-container">
                <router-link :to="{ name: 'stations' }">
                    <div class="map-link"><span class="small-arrow">&lt;</span> Stations Map</div>
                </router-link>
                <div id="station-name">{{ this.station ? this.station.name : "" }}</div>
                <DataChartControl
                    ref="dataChartControl"
                    :combinedStationInfo="combinedStationInfo"
                    :station="station"
                    :noStation="noStation"
                    :labels="labels"
                    @switchedSensor="onSensorSwitch"
                    @timeChanged="onTimeChange"
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
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import DataChartControl from "../components/DataChartControl";
import NotesList from "../components/NotesList";
import SensorSummary from "../components/SensorSummary";
import * as tempStations from "../assets/ancientGoose.json";

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
            requestPage: 0,
            user: {},
            noStation: false,
            station: null,
            stationId: null,
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
                rainPrevHour: "Rain Previous Hour"
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
                this.fetchSummary().then(result => {
                    this.handleInitialDataSummary(result);
                });
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

        handleInitialDataSummary(result) {
            if (!result) {
                this.sensors = [];
                this.stationData = [];
                this.combinedStationInfo = { sensors: this.sensors, stationData: this.stationData };
                return;
            }
            const processedData = this.processData(result);
            let sensors = processedData.sensors;
            let processed = processedData.data;
            //sort data by date
            processed.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });
            // get most recent reading for each sensor

            // TODO: when the window loads with a specific time range,
            // these might not be most recent

            sensors = _.uniqBy(sensors, "key");
            const numSensors = sensors.length;
            let numReadings = 0;
            if (processed.length > 0) {
                for (let i = processed.length - 1; i > 0; i--) {
                    let recent = processed[i];
                    sensors.forEach(s => {
                        if (!s.currentReading && recent[s.key]) {
                            numReadings += 1;
                            s.currentReading = recent[s.key];
                            if (numReadings == numSensors) {
                                i == 0;
                            }
                        }
                    });
                }

                // resize div to fit sections
                document.getElementById("lower-container").style["min-width"] = "1100px";
            }
            this.sensors = sensors;
            this.stationData = processed;
            this.combinedStationInfo = { sensors: sensors, stationData: processed };
        },

        processData(result) {
            let processed = [];
            let sensors = [];
            // TODO: handle possible changes to sensors
            result.modules.forEach(m => {
                if (m.sensors) {
                    m.sensors.forEach(s => {
                        sensors.push(s);
                    });
                }
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
            this.$refs.dataChartControl.updateData(this.stationData);
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
            this.fetchSummary().then(result => {
                this.handleInitialDataSummary(result);
            });
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
    font-weight: bold;
    margin: 20px 0;
}
</style>
