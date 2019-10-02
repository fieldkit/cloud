<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="data" />
        <div class="main-panel">
            <div id="data-container">
                <router-link :to="{ name: 'stations' }">
                    <div class="map-link">&lt; Stations Map</div>
                </router-link>
                <div id="station-name">{{ this.station ? this.station.name : "Data" }}</div>
                <div class="spacer"></div>

                <div id="readings-container" class="section" v-if="this.station">
                    <div id="readings-label">
                        Latest Reading <span class="synced">Last synced {{ getSyncedDate() }}</span>
                    </div>
                    <div v-for="module in this.station.status_json.moduleObjects" v-bind:key="module.id">
                        <div v-for="sensor in module.sensorObjects" v-bind:key="sensor.id" class="reading">
                            <div class="left">{{ sensor.name }}</div>
                            <div class="right">
                                {{ sensor.current_reading + " " + (sensor.unit ? sensor.unit : "") }}
                            </div>
                        </div>
                    </div>
                </div>
                <DataChart :summary="summary" :stationData="stationData" :station="station" />
                <NotesList :station="station" />
                <SensorSummary :sensor="selectedSensor" />
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
    props: ["station"],
    data: () => {
        return {
            user: {},
            stationData: [],
            selectedSensor: null,
            summary: [],
            isAuthenticated: false
        };
    },
    async mounted() {
        if (this.station) {
            this.summary = await this.api.getStationDataSummaryByDeviceId(this.station.device_id);
            this.stationData = await this.api.getJSONDataByDeviceId(this.station.device_id, 0, 1000);
            this.selectedSensor = this.station.status_json.moduleObjects[0].sensorObjects[0];
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

        getSyncedDate() {
            const date = this.station.status_json.updated;
            const d = new Date(date);
            return d.toLocaleDateString("en-US");
        }
    }
};
</script>

<style scoped>
.main-panel {
    width: 80%;
    margin-left: 20px;
}
.map-link {
    margin: 10px 0;
    font-size: 13px;
}
#station-name {
    font-size: 20px;
    font-weight: bold;
}
.spacer {
    width: 100%;
    margin: 10px 0 20px 0;
    border-top: 2px solid rgba(230, 230, 230);
}
.synced {
    margin-left: 10px;
    font-size: 11px;
}
#readings-label {
    margin-bottom: 10px;
}
.reading {
    font-size: 11px;
    width: 145px;
    float: left;
    margin: 5px 10px 5px 0;
    padding: 5px 10px;
    background-color: rgb(233, 233, 233);
    border: 1px solid rgb(200, 200, 200);
}
.left {
    float: left;
}
.right {
    float: right;
}
</style>
