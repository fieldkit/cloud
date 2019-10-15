<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="data" />
        <div id="data-view-background" v-show="isAuthenticated">
            <div class="main-panel">
                <div id="data-container">
                    <router-link :to="{ name: 'stations' }">
                        <div class="map-link"><span class="small-arrow">&lt;</span> Stations Map</div>
                    </router-link>
                    <div id="station-name">{{ this.station ? this.station.name : "Data" }}</div>
                    <DataChartControl
                        ref="dataChartControl"
                        :stationData="stationData"
                        :station="station"
                        :sensors="sensors"
                        :selectedSensor="selectedSensor"
                        :labels="labels"
                        @switchedSensor="onSensorSwitch"
                        @timeChanged="onTimeChange"
                    />
                    <NotesList :station="station" />
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
// import _ from "lodash";
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
            user: {},
            station: null,
            stationData: [],
            sensors: [],
            selectedSensor: null,
            summary: [],
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
            dataView.setSensor();
            dataView.setTimeWindow();
            dataView.$refs.dataChartControl.refresh();
        };

        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                this.setTimeWindow();
                this.fetchData().then(result => {
                    this.processData(result);
                });
            })
            .catch(() => {
                this.failedAuth = true;
            });
    },
    methods: {
        async fetchData() {
            if (this.stationParam) {
                this.station = this.stationParam;
                return this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
            } else if (this.id) {
                // temporarily show Ancient Goose 81 to anyone who views /dashboard/data/0
                if (this.id == 0) {
                    this.station = tempStations.stations[0];
                    return this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
                } else {
                    this.api.getStation(this.id).then(station => {
                        this.station = station;
                        return this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
                    });
                }
            }
        },
        processData(result) {
            let processed = [];
            let sensors = [];
            result.versions.forEach(v => {
                v.meta.station.modules.forEach(m => {
                    m.sensors.forEach(s => {
                        sensors.push(s);
                    });
                });
                v.data.forEach(d => {
                    // HACK: for now only including ones with
                    // sensor readings
                    if (Object.keys(d.d).length > 0) {
                        d.d.date = new Date(d.time * 1000);
                        processed.push(d.d);
                    }
                });
            });
            //sort data by date
            processed.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });
            // get most recent reading for each sensor
            let recent = processed[processed.length - 1];
            sensors.forEach(s => {
                s.currentReading = recent[s.key];
            });
            this.sensors = sensors;
            this.stationData = processed;
            this.setSensor();
        },
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        setSensor() {
            // get selected sensor from url
            if (this.$route.query.sensor) {
                this.sensors.forEach(s => {
                    if (s.key == this.$route.query.sensor) {
                        this.selectedSensor = s;
                    }
                });
            }
            // or set the first sensor to be selected sensor
            if (!this.$route.query.sensor || !this.selectedSensor) {
                this.selectedSensor = this.sensors[0];
            }
        },
        setTimeWindow() {
            if (this.$route.query.start && this.$route.query.end) {
                let newStart = new Date(parseInt(this.$route.query.start));
                let newEnd = new Date(parseInt(this.$route.query.end));
                this.timeRange = { start: newStart, end: newEnd };
            } else {
                this.timeRange = null;
            }
        },
        onSensorSwitch(sensor) {
            this.selectedSensor = sensor;
        },
        onTimeChange(range) {
            this.timeRange = range;
        }
    }
};
</script>

<style scoped>
#data-view-background {
    float: left;
    background-color: rgb(252, 252, 252);
}
.no-auth-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 300px;
}
.show-link {
    text-decoration: underline;
}
.main-panel {
    margin-left: 280px;
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
