<template>
    <div>
        <SidebarNav
            viewing="stations"
            :isAuthenticated="isAuthenticated"
            :stations="stations"
            :projects="projects"
            @showStation="showSummary"
        />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" @sidebarToggled="onSidebarToggle" />
            <StationSummary
                v-if="activeStation"
                :isAuthenticated="isAuthenticated"
                :station="activeStation"
                :placeName="placeName"
                :nativeLand="nativeLand"
                ref="stationSummary"
            />
        <div id="stations-view-panel" class="main-panel full-height">
            <div id="summary-and-map">
                <StationsMap :mapSize="mapSize" :stations="stations" @mapReady="onMapReady" @showSummary="showSummary" ref="stationsMap" />
            <div v-if="isAuthenticated && showNotice" id="no-stations">
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
            <div v-if="noCurrentUser" id="no-user">
                <p>
                    Please
                    <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="show-link">
                        log in
                    </router-link>
                    to view stations.
                </p>
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import FKApi from "@/api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import StationSummary from "../components/StationSummary";
import StationsMap from "../components/StationsMap";

export default {
    name: "StationsView",
    components: {
        HeaderBar,
        SidebarNav,
        StationsMap,
        StationSummary,
    },
    props: ["id"],
    data: () => {
        return {
            user: {},
            projects: [],
            stations: [],
            activeStation: null,
            showNotice: false,
            placeName: "",
            nativeLand: [],
            isAuthenticated: false,
            noCurrentUser: false,
            mapSize: {
                width: "100%",
                height: "100vh",
                position: "absolute",
            },
        };
    },
    watch: {
        $route(to) {
            if (to.name == "stations") {
                this.$refs.stationSummary.closeSummary();
            }
        },
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;

                return Promise.all([
                    Promise.resolve().then(() => {
                        if (this.id) {
                            return this.api.getStation(this.id).then(station => {
                                return this.showSummary(station, true);
                            });
                        }
                    }),
                    this.api.getStations().then(s => {
                        this.stations = s.stations;
                    }),
                    this.api.getUserProjects().then(projects => {
                        if (projects && projects.projects.length > 0) {
                            this.projects = projects.projects;
                        }
                    }),
                ]);
            })
            .then(() => {
                if (this.stations.length == 0 && !this.id) {
                    this.showNotice = true;
                }
            })
            .catch(e => {
                console.log("error", e);
                this.noCurrentUser = true;
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },

        onMapReady(map) {
            this.map = map;
        },

        getPlaceName() {
            const longLatMapbox = this.activeStation.status_json.longitude + "," + this.activeStation.status_json.latitude;
            return this.api.getPlaceName(longLatMapbox).then(result => {
                this.placeName = result.features[0] ? result.features[0].place_name : "Unknown area";
            });
        },

        getNativeLand() {
            const latLongNative = this.activeStation.status_json.latitude + "," + this.activeStation.status_json.longitude;
            return this.api.getNativeLand(latLongNative);
        },

        showSummary(station, preserveRoute) {
            this.activeStation = station;
            this.getPlaceName()
                .then(this.getNativeLand)
                .then(result => {
                    this.nativeLand = _.map(_.flatten(_.map(result, "properties")), r => {
                        return { name: r.Name, url: r.description };
                    });
                    this.$refs.stationSummary.viewSummary();
                    if (!preserveRoute) {
                        this.updateStationRoute(station);
                    }
                });
        },

        updateStationRoute(station) {
            if (this.$route.name != "viewStation" || this.$route.params.id != station.id) {
                this.$router.push({ name: "viewStation", params: { id: station.id } });
            }
        },

        closeNotice() {
            this.showNotice = false;
        },
    },
};
</script>

<style scoped>
#stations-view-panel {
    margin: 0;
}
#map {
    width: 100%;
    height: 100vh;
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
