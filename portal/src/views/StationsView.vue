<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="stations" />
        <div class="main-panel">
            <div id="stations-container">
                <div class="container">
                    <h1>Stations</h1>
                    <StationsList :stations="stations" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import StationsList from "../components/StationsList";

export default {
    name: "StationsView",
    components: {
        HeaderBar,
        SidebarNav,
        StationsList
    },
    data: () => {
        return {
            user: {},
            stations: {},
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
            }
            this.isAuthenticated = true;
            console.log("this is the user info", this.user);
            console.log("this is the station info", this.stations);
        }
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        }
    }
};
</script>

<style scoped>
#stations-container .container {
    padding: 0;
}
#stations-container {
    margin-bottom: 60px;
}
</style>
