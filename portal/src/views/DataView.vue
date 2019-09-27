<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="data" />
        <div class="main-panel">
            <div id="data-container">
                <div class="container">
                    <h1>Data</h1>
                    <DataExample :summary="summary" :stationData="stationData" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import DataExample from "../components/DataExample";

export default {
    name: "DataView",
    components: {
        HeaderBar,
        SidebarNav,
        DataExample
    },
    data: () => {
        return {
            user: {},
            stations: {},
            stationData: [],
            summary: [],
            isAuthenticated: false
        };
    },
    async beforeCreate() {
        const api = new FKApi();
        if (api.authenticated()) {
            this.user = await api.getCurrentUser();
            this.stations = await api.getStations();
            if (this.stations) {
                this.stations = this.stations.stations;
                if (this.stations.length > 0) {
                    const id = _(this.stations).first().device_id;
                    this.summary = await api.getStationDataSummaryByDeviceId(id);
                    this.stationData = await api.getJSONDataByDeviceId(id, 0, 1000);
                }
            }
            this.isAuthenticated = true;
        }
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        }
    }
};
</script>

<style scoped></style>
