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
        <div id="stations-view-panel" class="main-panel">
            <mapbox
                class="stations-map"
                :access-token="mapboxToken"
                :map-options="{
                    style: 'mapbox://styles/mapbox/outdoors-v11',
                    center: coordinates,
                    zoom: 14,
                }"
                :nav-control="{
                    show: true,
                    position: 'bottom-left',
                }"
                @map-init="mapInitialized"
            />
            <StationSummary
                v-if="activeStation"
                :isAuthenticated="isAuthenticated"
                :station="activeStation"
                :placeName="placeName"
                :nativeLand="nativeLand"
                ref="stationSummary"
            />
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
            <div v-if="failedAuth" id="no-auth">
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
import Mapbox from "mapbox-gl-vue";
import { MAPBOX_ACCESS_TOKEN } from "../secrets";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import StationSummary from "../components/StationSummary";

export default {
    name: "StationsView",
    components: {
        Mapbox,
        HeaderBar,
        SidebarNav,
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
            failedAuth: false,
            coordinates: [-96, 37.8],
            mapboxToken: MAPBOX_ACCESS_TOKEN,
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
                        if (this.map) {
                            this.initStations();
                        } else {
                            this.waitingForMap = true;
                        }
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
                // NOTE: This isn't necessarily true, we should be
                // checking for actual authentication error status codes
                // in the Api code and then raising some kind of event to
                // indicate they aren't authenticated. Added this message
                // because I saw this happen even though I was logged in
                // and was curious why but couldn't see.
                console.log("error", e);
                this.failedAuth = true;
                // handle non-logged in state
                if (this.map) {
                    this.initStations();
                } else {
                    this.waitingForMap = true;
                }
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },

        onSidebarToggle() {
            this.map.resize();
        },

        mapInitialized(map) {
            this.map = map;
            if (this.waitingForMap) {
                this.initStations();
            }
        },

        initStations() {
            if (this.stations && this.stations.length > 0) {
                this.updateMap();
            } else {
                this.map.setZoom(6);
            }
        },

        updateMap() {
            let longMax = -180;
            let longMin = 180;
            let latMin = 90;
            let latMax = -90;
            let stationFeatures = [];
            let mappable = this.stations.filter(s => {
                return (
                    s.status_json.latitude && s.status_json.longitude && s.status_json.latitude != 1000 && s.status_json.longitude != 1000
                );
            });
            mappable.forEach(s => {
                let coordinates = [s.status_json.latitude, s.status_json.longitude];
                if (mappable.length == 1) {
                    this.map.setCenter({
                        lat: coordinates[0],
                        lng: coordinates[1],
                    });
                } else {
                    if (s.status_json.latitude > latMax) {
                        latMax = s.status_json.latitude;
                    }
                    if (s.status_json.latitude < latMin) {
                        latMin = s.status_json.latitude;
                    }
                    if (s.status_json.longitude > longMax) {
                        longMax = s.status_json.longitude;
                    }
                    if (s.status_json.longitude < longMin) {
                        longMin = s.status_json.longitude;
                    }
                }
                stationFeatures.push({
                    type: "Feature",
                    geometry: {
                        type: "Point",
                        coordinates: [coordinates[1], coordinates[0]],
                    },
                    properties: {
                        title: s.name,
                        icon: "marker",
                    },
                });
            });
            // sort by latitude
            stationFeatures.sort(function(a, b) {
                var latA = a.geometry.coordinates[0];
                var latB = b.geometry.coordinates[0];
                if (latA < latB) {
                    return -1;
                }
                if (latA > latB) {
                    return 1;
                }
                // lats must be equal
                return 0;
            });
            // could sort by and compare longitude, too

            // TODO: handle this better - merge into aggregate markers
            // currently if stations overlap exactly, only one shows
            stationFeatures.forEach((f, i) => {
                if (i > 0) {
                    if (f.geometry.coordinates[0] == stationFeatures[i - 1].geometry.coordinates[0]) {
                        // TODO: what is an acceptable amount to offset labels by?
                        // less than this, and still only one shows
                        f.geometry.coordinates[0] += 0.00003;
                    }
                }
            });

            this.map.addLayer({
                id: "station-marker",
                type: "symbol",
                source: {
                    type: "geojson",
                    data: {
                        type: "FeatureCollection",
                        features: stationFeatures,
                    },
                },
                layout: {
                    "icon-image": "{icon}-15",
                    "icon-size": 2,
                    "text-field": "{title}",
                    "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                    "text-offset": [0, 0.75],
                    "text-anchor": "top",
                },
                // paint: {
                //     "text-color": "#0000ff"
                // }
            });

            if (mappable.length > 1) {
                let bounds = [
                    {
                        lat: latMin,
                        lng: longMin,
                    },
                    {
                        lat: latMax,
                        lng: longMax,
                    },
                ];
                this.map.fitBounds(bounds, { duration: 0 });
                // fitBounds is cropped too close: zoom out
                const z = this.map.getZoom();
                this.map.setZoom(z - 1.5);
            } else if (mappable.length == 0) {
                this.map.setZoom(6);
            }

            const stationsView = this;
            this.map.on("click", "station-marker", function(e) {
                const name = e.features[0].properties.title;
                const station = stationsView.stations.find(s => {
                    return s.name == name;
                });
                stationsView.showSummary(station);
            });
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
#no-auth {
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
