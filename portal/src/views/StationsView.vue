<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="stations" :stations="stations" @showStation="showSummary" />
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
    data: () => {
        return {
            user: {},
            stations: {
                stations: []
            },
            activeStation: null,
            isAuthenticated: false,
            coordinates: [-96, 37.8],
            mapboxToken: MAPBOX_ACCESS_TOKEN
        };
    },
    mounted() {
        // HACK
        const newWidth = window.innerWidth - 240; // the sidebar is 240px
        document.getElementById("stations-view-panel").style.width = newWidth + "px";
        this.map.resize();
    },
    async beforeCreate() {
        const api = new FKApi();
        api.getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                api.getStations().then(stations => {
                    this.stations = stations;
                    if (this.map) {
                        this.initStations();
                    } else {
                        this.waitingForMap = true;
                    }
                    console.log("this is the user info", this.user);
                    console.log("this is the station info", this.stations);
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

        mapInitialized(map) {
            this.map = map;
            if (this.waitingForMap) {
                this.initStations();
            }
        },

        initStations() {
            if (this.stations && this.stations.stations.length > 0) {
                this.stations = this.stations.stations;
                this.updateMap();
            } else {
                this.$refs.stationSummary.viewSummary();
                this.map.setZoom(6);
            }
        },

        updateMap() {
            let bounds = [];
            let stationFeatures = [];
            let mappable = this.stations.filter(s => {
                return s.status_json.latitude && s.status_json.longitude;
            });
            mappable.forEach(s => {
                let coordinates = [s.status_json.latitude, s.status_json.longitude];
                if (mappable.length == 1) {
                    this.map.setCenter({
                        lat: coordinates[0],
                        lng: coordinates[1]
                    });
                } else {
                    bounds.push({
                        lat: coordinates[0],
                        lng: coordinates[1]
                    });
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
            });
        },

        showSummary(station) {
            this.activeStation = station;
            this.$refs.stationSummary.viewSummary();
        }
    }
};
</script>

<style scoped>
.main-panel {
    width: 80%;
}
#stations-container .container {
    padding: 0;
}
#stations-container {
    margin-bottom: 60px;
}
#map {
    width: 100%;
    height: 100%;
    min-height: 500px;
}
</style>
