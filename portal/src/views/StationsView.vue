<template>
    <div class="full-height">
        <SidebarNav
            class="full-height left"
            viewingStations="true"
            :isAuthenticated="isAuthenticated"
            :stations="stations"
            :projects="userProjects"
            @showStation="showSummary"
        />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" @sidebarToggled="onSidebarToggle" />
        <div id="stations-view-panel" class="main-panel full-height">
            <div id="summary-and-map">
                <StationsMap :stations="stations" @mapReady="onMapReady" @showSummary="showSummary" ref="stationsMap" />
                <StationSummary
                    v-if="activeStationId"
                    @closeSummary="closeSummary"
                    class="summary-container"
                    :station="activeStation"
                    :summarySize="summarySize"
                    ref="stationSummary"
                />
            </div>
            <div v-if="isAuthenticated && showNoStationsMessage && hasNoStations" id="no-stations">
                <div id="close-notice-btn" v-on:click="closeNotice">
                    <img alt="Close" src="../assets/close.png" />
                </div>
                <p class="heading">Add a New Station</p>
                <p class="text">
                    You currently don't have any stations on your account. Download the FieldKit app and connect your station to add them to
                    your account.
                </p>
                <a href="https://apps.apple.com/us/app/fieldkit-org/id1463631293?ls=1" target="_blank">
                    <img alt="App store" src="../assets/appstore.png" class="app-btn" />
                </a>
                <a href="https://play.google.com/store/apps/details?id=com.fieldkit&hl=en_US" target="_blank">
                    <img alt="Google Play" src="../assets/googleplay.png" class="app-btn" />
                </a>
            </div>
            <div v-if="!isAuthenticated" id="no-user">
                <p>
                    Please
                    <router-link v-if="$route" :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="show-link">
                        log in
                    </router-link>
                    <router-link v-if="!$route" :to="{ name: 'login' }" class="show-link">log in</router-link>
                    to view stations.
                </p>
            </div>
        </div>
    </div>
</template>

<script>
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import StationSummary from "../components/StationSummary";
import StationsMap from "../components/StationsMap";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
    name: "StationsView",
    components: {
        HeaderBar,
        SidebarNav,
        StationsMap,
        StationSummary,
    },
    props: {
        id: { type: Number },
    },
    data: () => {
        return {
            activeStationId: null,
            showNoStationsMessage: true,
            summarySize: {
                width: "415px",
                top: "120px",
                left: "120px",
                constrainTop: "285px",
            },
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: s => s.user.user,
            hasNoStations: s => s.stations.hasNoStations,
            stations: s => s.stations.stations.user,
            userProjects: s => s.stations.projects.user,
            anyStations: s => s.stations.stations.user.length > 0,
        }),
        activeStation() {
            return this.$store.state.stations.stations.all[this.activeStationId];
        },
    },
    beforeMount() {
        this.$store.dispatch(ActionTypes.NEED_PROJECTS);
        this.$store.dispatch(ActionTypes.NEED_STATIONS);
        if (this.id) {
            this.activeStationId = this.id;
            return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.id });
        }
    },
    methods: {
        goBack() {
            if (window.history.length) {
                return this.$router.go(-1);
            } else {
                return this.$router.push("/");
            }
        },
        showSummary(station, preserveRoute) {
            this.activeStationId = station.id;
            this.updateStationRoute(station);
        },
        closeSummary() {
            this.activeStationId = null;
        },
        updateStationRoute(station) {
            if (this.$route.name != "viewStation" || this.$route.params.id != station.id) {
                this.$router.push({ name: "viewStation", params: { id: station.id } });
            }
        },
        onMapReady(map) {
            this.map = map;
        },
        onSidebarToggle() {
            if (this.map) {
                this.map.resize();
            }
        },
    },
};
</script>

<style scoped>
.full-height {
    height: 100%;
}
.left {
    float: left;
}
#stations-view-panel {
    margin: 0;
    overflow: scroll;
}
#no-user {
    font-size: 20px;
    background-color: #ffffff;
    width: 400px;
    position: absolute;
    top: 40px;
    left: 260px;
    padding: 0 15px 15px 15px;
    margin: 60px;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
}
#no-stations {
    background-color: #ffffff;
    width: 360px;
    position: absolute;
    top: 70px;
    left: 50%;
    padding: 75px 60px 75px 60px;
    margin: 135px 0 0 -120px;
    text-align: center;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
}
#no-stations .heading {
    font-size: 18px;
    font-weight: bold;
}
#no-stations .text {
    font-size: 14px;
}
#no-stations .app-btn {
    margin: 20px 14px;
}
#close-notice-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    cursor: pointer;
}
.show-link {
    text-decoration: underline;
}
</style>
