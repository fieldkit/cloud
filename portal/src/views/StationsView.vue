<template>
    <div>
        <SidebarNav viewing="stations" :stations="stations" :projects="projects" @showStation="showSummary" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" @sidebarToggled="onSidebarToggle" />
        <div id="stations-view-panel" class="main-panel">
            <mapbox
                :access-token="mapboxToken"
                :map-options="{
                    style: 'mapbox://styles/mapbox/light-v10',
                    center: coordinates,
                    zoom: 14
                }"
                :nav-control="{
                    show: true,
                    position: 'bottom-left'
                }"
                @map-init="mapInitialized"
            />
            <StationSummary
                :isAuthenticated="isAuthenticated"
                :station="activeStation"
                ref="stationSummary"
            />
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
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
        StationSummary
    },
    props: ["id"],
    data: () => {
        return {
            user: {},
            projects: [],
            stations: {
                stations: []
            },
            activeStation: null,
            isAuthenticated: false,
            coordinates: [-96, 37.8],
            mapboxToken: MAPBOX_ACCESS_TOKEN
        };
    },
    watch: {
        $route(to) {
            if (to.name == "stations") {
                this.$refs.stationSummary.closeSummary();
            }
        }
    },
    async beforeCreate() {
        const api = new FKApi();
        api.getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                api.getStations().then(stations => {
                    this.stations = stations && stations.stations ? stations.stations : [];
                    if (this.map) {
                        this.initStations();
                    } else {
                        this.waitingForMap = true;
                    }
                    if (this.id) {
                        let station = this.stations.find(s => {
                            return s.id == this.id;
                        });
                        this.activeStation = station;
                        this.$refs.stationSummary.viewSummary();
                    }
                });
                api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
            })
            .catch(() => {
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
                this.$refs.stationSummary.viewSummary();
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
                    s.status_json.latitude &&
                    s.status_json.longitude &&
                    s.status_json.latitude != 1000 &&
                    s.status_json.longitude != 1000
                );
            });
            mappable.forEach(s => {
                let coordinates = [s.status_json.latitude, s.status_json.longitude];
                if (mappable.length == 1) {
                    this.map.setCenter({
                        lat: coordinates[0],
                        lng: coordinates[1]
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
                        coordinates: [coordinates[1], coordinates[0]]
                    },
                    properties: {
                        title: s.name,
                        icon: "marker"
                    }
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
                        features: stationFeatures
                    }
                },
                layout: {
                    "icon-image": "{icon}-15",
                    "icon-size": 2,
                    "text-field": "{title}",
                    "text-font": ["Open Sans Semibold", "Arial Unicode MS Bold"],
                    "text-offset": [0, 0.75],
                    "text-anchor": "top"
                }
                // paint: {
                //     "text-color": "#0000ff"
                // }
            });

            if (mappable.length > 1) {
                let bounds = [
                    {
                        lat: latMin,
                        lng: longMin
                    },
                    {
                        lat: latMax,
                        lng: longMax
                    }
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
                stationsView.activeStation = station;
                stationsView.$refs.stationSummary.viewSummary();
                stationsView.updateStationRoute(station);
            });
        },

        showSummary(station) {
            this.activeStation = station;
            this.$refs.stationSummary.viewSummary();
            this.updateStationRoute(station);
        },

        updateStationRoute(station) {
            if (this.$route.name != "viewStation" || this.$route.params.id != station.id) {
                this.$router.push({ name: "viewStation", params: { id: station.id } });
            }
        }
    }
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
</style>
