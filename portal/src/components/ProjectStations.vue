<template>
    <div class="stations-container" v-bind:style="{ width: mapSize.containerWidth }">
        <div class="section-heading stations-heading">FieldKit Stations</div>
        <div class="space"></div>
        <div class="stations-list" v-bind:style="{ width: listSize.width, height: listSize.height }">
            <div v-for="station in projectStations" v-bind:key="station.id">
                <div class="station-box" v-bind:style="{ width: listSize.boxWidth }">
                    <span class="station-name" v-on:click="showStation(station)">
                        {{ station.name }}
                    </span>
                    <div class="last-seen">Last seen {{ getUpdatedDate(station) }}</div>
                </div>
            </div>
        </div>
        <div class="stations-map" v-bind:style="{ width: mapSize.width, height: mapSize.height }">
            <mapbox
                :access-token="mapboxToken"
                :map-options="{
                    style: 'mapbox://styles/mapbox/outdoors-v11',
                    center: coordinates,
                    zoom: 10,
                }"
                :nav-control="{
                    show: true,
                    position: 'bottom-left',
                }"
                @map-init="mapInitialized"
            />
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as utils from "../utilities";
import FKApi from "../api/api";
import Mapbox from "mapbox-gl-vue";
import { MAPBOX_ACCESS_TOKEN } from "../secrets";

export default {
    name: "ProjectStations",
    components: {
        Mapbox,
    },
    data: () => {
        return {
            projectStations: [],
            coordinates: [-118, 34],
            mapboxToken: MAPBOX_ACCESS_TOKEN,
            following: false,
        };
    },
    props: ["project", "mapSize", "listSize"],
    async beforeCreate() {
        this.api = new FKApi();
    },
    mounted() {
        this.fetchStations();
    },
    methods: {
        mapInitialized(map) {
            this.map = map;
        },
        fetchStations() {
            this.api.getStationsByProject(this.project.id).then(result => {
                this.projectStations = result.stations;
                let modules = [];
                if (this.projectStations) {
                    this.projectStations.forEach((s, i) => {
                        if (i == 0 && s.status_json.latitude && this.map) {
                            this.map.setCenter({
                                lat: parseFloat(s.status_json.latitude),
                                lng: parseFloat(s.status_json.longitude),
                            });
                        }
                        if (s.status_json.moduleObjects) {
                            s.status_json.moduleObjects.forEach(m => {
                                modules.push(this.getModuleImg(m));
                            });
                        } else if (s.status_json.statusJson && s.status_json.statusJson.modules) {
                            s.status_json.statusJson.modules.forEach(m => {
                                modules.push(this.getModuleImg(m));
                            });
                        }
                        // could also use readings, if present
                    });
                    modules = _.uniq(modules);
                    this.$emit("loaded", modules);
                }
            });
        },
        getUpdatedDate(station) {
            if (!station.status_json) {
                return "N/A";
            }
            const date = station.status_json.updated;
            const d = new Date(date);
            return d.toLocaleDateString("en-US");
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        getModuleImg(module) {
            let imgPath = require.context("../assets/modules-lg/", false, /\.png$/);
            let img = utils.getModuleImg(module);
            return imgPath("./" + img);
        },
    },
};
</script>

<style scoped>
.section {
    float: left;
}
.section-heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 0 0 35px 0;
}
.stations-heading {
    margin: 25px 0 25px 25px;
}
.stations-container .space {
    margin: 0;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.last-seen {
    font-size: 12px;
    font-weight: 600;
    color: #6a6d71;
}
.space {
    width: 100%;
    float: left;
    margin: 30px 0 0 0;
    border-bottom: solid 1px #d8dce0;
}
.stations-container {
    float: left;
    margin: 22px 0 0 0;
    border: 1px solid #d8dce0;
}
.stations-list {
    float: left;
}
.stations-map {
    float: left;
}
#map {
    width: inherit;
    height: inherit;
}
.station-box {
    height: 38px;
    margin: 20px auto;
    padding: 10px;
    border: 1px solid #d8dce0;
}
</style>
