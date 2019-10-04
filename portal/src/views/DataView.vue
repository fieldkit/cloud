<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="data" />
        <div id="data-view-background">
            <div class="main-panel">
                <div id="data-container">
                    <router-link :to="{ name: 'stations' }">
                        <div class="map-link"><span class="small-arrow">&lt;</span> Stations Map</div>
                    </router-link>
                    <div id="station-name">{{ this.station ? this.station.name : "Data" }}</div>
                    <DataChart
                        :summary="summary"
                        :stationData="stationData"
                        :station="station"
                        :selectedSensor="selectedSensor"
                        @switchedSensor="onSensorSwitch"
                    />
                    <NotesList :station="station" />
                    <SensorSummary :sensor="selectedSensor" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
// import _ from "lodash";
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import DataChart from "../components/DataChart";
import NotesList from "../components/NotesList";
import SensorSummary from "../components/SensorSummary";

export default {
    name: "DataView",
    components: {
        HeaderBar,
        SidebarNav,
        DataChart,
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
            isAuthenticated: false
        };
    },
    async mounted() {
        if (this.stationParam) {
            this.station = this.stationParam;
            this.summary = await this.api.getStationDataSummaryByDeviceId(this.station.device_id);
            this.stationData = await this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
            this.selectedSensor = this.station.status_json.moduleObjects[0].sensorObjects[0];
        } else if (this.id) {
            this.api.getStation(this.id).then(station => {
                this.station = station;
                this.selectedSensor = this.station.status_json.moduleObjects[0].sensorObjects[0];
                this.api.getStationDataSummaryByDeviceId(this.station.device_id).then(summary => {
                    this.summary = summary;
                });
                this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000).then(data => {
                    this.stationData = data;
                });
            });
        }
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api.getCurrentUser().then(user => {
            this.user = user;
            this.isAuthenticated = true;
        });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        onSensorSwitch(sensor) {
            this.selectedSensor = sensor;
        }
    }
};
</script>

<style scoped>
#data-view-background {
    float: left;
    background-color: rgb(252, 252, 252);
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
