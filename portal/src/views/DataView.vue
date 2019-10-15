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
                        :summary="summary"
                        :stationData="stationData"
                        :station="station"
                        :selectedSensor="selectedSensor"
                        @switchedSensor="onSensorSwitch"
                        @timeChanged="onTimeChange"
                    />
                    <NotesList :station="station" />
                    <SensorSummary
                        ref="sensorSummary"
                        :selectedSensor="selectedSensor"
                        :stationData="stationData"
                        :timeRange="timeRange"
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
            selectedSensor: null,
            summary: [],
            isAuthenticated: false,
            failedAuth: false,
            timeRange: null
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
                this.fetchData();
            })
            .catch(() => {
                this.failedAuth = true;
            });
    },
    methods: {
        async fetchData() {
            if (this.stationParam) {
                this.station = this.stationParam;
                this.summary = await this.api.getStationDataSummaryByDeviceId(this.station.device_id);
                this.stationData = await this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
                this.setSensor();
            } else if (this.id) {
                // temporarily show Ancient Goose 81 to anyone who views /dashboard/data/0
                if (this.id == 0) {
                    this.station = tempStations.stations[0];
                    this.setSensor();
                    this.api.getStationDataSummaryByDeviceId(this.station.device_id).then(summary => {
                        this.summary = summary;
                    });
                    this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000).then(data => {
                        this.stationData = data;
                    });
                } else {
                    this.api.getStation(this.id).then(station => {
                        this.station = station;
                        this.setSensor();
                        this.api.getStationDataSummaryByDeviceId(this.station.device_id).then(summary => {
                            this.summary = summary;
                        });
                        this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000).then(data => {
                            this.stationData = data;
                        });
                    });
                }
            }
        },
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        setSensor() {
            const modules = this.station.status_json.moduleObjects;
            // get selected sensor from url
            if (this.$route.query.sensor) {
                modules.forEach(m => {
                    m.sensorObjects.forEach(s => {
                        if (s.name == this.$route.query.sensor) {
                            this.selectedSensor = s;
                        }
                    });
                });
            }
            // or set the first sensor to be selected sensor
            if (!this.$route.query.sensor || !this.selectedSensor) {
                if (modules.length > 0 && modules[0].sensorObjects.length > 0) {
                    this.selectedSensor = modules[0].sensorObjects[0];
                }
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
