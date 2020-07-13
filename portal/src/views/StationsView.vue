<template>
    <StandardLayout @sidebar-toggle="onSidebarToggle" :viewingStations="true">
        <div class="container-map">
            <StationsMap @mapReady="onMapReady" @showSummary="showSummary" :mapped="mapped" />
            <StationSummary
                v-if="activeStation"
                @closeSummary="closeSummary"
                class="summary-container"
                :station="activeStation"
                :summarySize="summarySize"
                ref="stationSummary"
            />
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
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "./StandardLayout.vue";
import StationSummary from "../components/StationSummary.vue";
import StationsMap from "../components/StationsMap.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "StationsView",
    components: {
        StandardLayout,
        StationsMap,
        StationSummary,
    },
    props: {
        id: { type: Number },
    },
    data: () => {
        return {
            map: null,
            showNoStationsMessage: true,
            summarySize: {
                width: "415px",
                top: "120px",
                left: "360px",
                constrainTop: "285px",
            },
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy", mapped: "mapped" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            hasNoStations: (s: GlobalState) => s.stations.hasNoStations,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
            anyStations: (s: GlobalState) => s.stations.user.stations.length > 0,
        }),
        activeStation() {
            console.log("station", this.$store.state.stations.stations[this.id]);
            return this.$store.state.stations.stations[this.id];
        },
    },
    beforeMount() {
        if (this.id) {
            return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.id });
        }
    },
    watch: {
        id() {
            if (this.id) {
                return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.id });
            }
        },
    },
    methods: {
        goBack() {
            if (window.history.length) {
                return this.$router.go(-1);
            } else {
                return this.$router.push("/");
            }
        },
        showSummary(params, preserveRoute) {
            return this.$router.push({ name: "viewStation", params: params });
        },
        closeSummary() {
            return this.$router.push({ name: "stations" });
        },
        onMapReady(map) {
            this.map = map;
        },
        onSidebarToggle() {
            console.log("stations toggle");
            if (this.map) {
                this.map.resize();
            }
        },
        closeNotice() {
            this.showNoStationsMessage = false;
        },
    },
});
</script>

<style scoped>
.container-ignored {
    height: 100%;
}
.container-top {
    display: flex;
    flex-direction: row;
    height: 100%;
}
.container-side {
    width: 240px;
    height: 100%;
}
.container-main {
    flex-grow: 1;
    flex-direction: column;
    display: flex;
}
.container-header {
    height: 70px;
}
.container-map {
    flex-grow: 1;
}
#stations-view-panel {
    margin: 0;
}
.closeSummary {
    position: absolute;
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
