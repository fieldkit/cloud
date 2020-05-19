<template>
    <div
        id="station-summary-container"
        :style="{ width: summarySize.width, top: summarySize.top, left: summarySize.left }"
        v-if="viewingSummary"
    >
        <div v-if="station" id="close-form-btn" v-on:click="closeSummary">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <div class="station-container" v-if="station">
            <div class="left image-container">
                <img alt="Station image" :src="stationSmallPhoto" class="station-element" />
            </div>
            <div class="left">
                <div id="station-name" class="station-element">{{ station.name }}</div>
                <div class="station-element">
                    Last Synced
                    <span class="small-light">{{ getSyncedDate() }}</span>
                </div>
                <div class="station-element">
                    <img id="battery" alt="Battery level" :src="getBatteryImg()" />
                    <span class="small-light">{{ station.status_json.batteryLevel }}%</span>
                </div>
                <div>
                    <img
                        v-for="module in station.modules"
                        v-bind:key="module.id"
                        alt="Module icon"
                        class="small-space"
                        :src="getModuleImg(module)"
                    />
                </div>
            </div>
            <div class="spacer"></div>
            <div id="location-container" class="section">
                <div class="location-item" v-if="placeName">
                    <img alt="location-icon" src="../assets/icon-location.png" />
                    {{ station.status_json.locationName ? station.status_json.locationName : placeName }}
                </div>
                <div class="location-item" v-if="nativeLand.length > 0">
                    <img alt="location-icon" src="../assets/icon-location.png" />
                    Native Land:
                    <span v-for="(n, i) in nativeLand" v-bind:key="n.url" class="note-container">
                        <a :href="n.url" class="native-land-link" target="_blank">{{ n.name }}</a>
                        <span>{{ i == nativeLand.length - 1 ? "" : ", " }}</span>
                    </span>
                </div>
                <div class="left">
                    {{ getLat() || "--" }}
                    <br />
                    Latitude
                </div>
                <div class="left">
                    {{ getLong() || "--" }}
                    <br />
                    Longitude
                </div>
            </div>
            <template v-if="!compact">
                <div class="spacer"></div>
                <div id="readings-container" class="section">
                    <div id="readings-label">Latest Reading</div>
                    <div v-for="(module, moduleIndex) in station.modules" v-bind:key="module.id">
                        <template v-if="module.position < 5">
                            <div
                                v-for="(sensor, sensorIndex) in module.sensors"
                                v-bind:key="sensor.id"
                                :class="getCounter(moduleIndex, sensorIndex) % 2 == 1 ? 'left-reading' : 'right-reading'"
                            >
                                <div class="left sensor-name">{{ getSensorName(module, sensor) }}</div>
                                <div class="right sensor-unit">
                                    {{ sensor.unit_of_measure }}
                                </div>
                                <div class="right sensor-reading">
                                    {{ getReading(sensor) }}
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
                <router-link :to="{ name: 'viewData', query: { stationId: station.id } }">
                    <div id="view-data-btn" class="section">
                        Explore Data
                    </div>
                </router-link>
            </template>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as utils from "../utilities";
import FKApi from "../api/api";
import { makeAuthenticatedApiUrl } from "@/api/api";

export default {
    name: "StationSummary",
    data: () => {
        return {
            moduleSensorCounter: 0,
            modulesSensors: {},
            viewingSummary: false,
            nativeLand: [],
            placeName: "",
        };
    },
    props: ["station", "summarySize", "compact"],
    computed: {
        stationSmallPhoto: function() {
            return makeAuthenticatedApiUrl(this.station.photos.small);
        },
    },
    watch: {
        station() {
            if (this.station) {
                this.getPlaceName()
                    .then(this.getNativeLand)
                    .then(result => {
                        this.nativeLand = _.map(_.flatten(_.map(result, "properties")), r => {
                            return { name: r.Name, url: r.description };
                        });
                    });
            }
        },
    },
    async beforeCreate() {
        this.api = new FKApi();
    },
    methods: {
        viewSummary() {
            this.viewingSummary = true;
        },

        getPlaceName() {
            const longLatMapbox = this.station.status_json.longitude + "," + this.station.status_json.latitude;
            return this.api.getPlaceName(longLatMapbox).then(result => {
                this.placeName = result.features[0] ? result.features[0].place_name : "Unknown area";
            });
        },

        getNativeLand() {
            const latLongNative = this.station.status_json.latitude + "," + this.station.status_json.longitude;
            return this.api.getNativeLand(latLongNative);
        },

        getReading(sensor) {
            if (!sensor.reading) {
                return "--";
            }
            return sensor.reading.last || sensor.reading.last == 0 ? sensor.reading.last.toFixed(1) : "--";
        },

        getCounter(moduleIndex, sensorIndex) {
            if (this.modulesSensors[moduleIndex]) {
                if (!this.modulesSensors[moduleIndex][sensorIndex]) {
                    this.moduleSensorCounter += 1;
                    this.modulesSensors[moduleIndex][sensorIndex] = this.moduleSensorCounter;
                }
            } else {
                this.moduleSensorCounter += 1;
                this.modulesSensors[moduleIndex] = {};
                this.modulesSensors[moduleIndex][sensorIndex] = this.moduleSensorCounter;
            }
            return this.modulesSensors[moduleIndex][sensorIndex];
        },

        getSyncedDate() {
            return utils.getUpdatedDate(this.station);
        },

        getBatteryImg() {
            const imgPath = require.context("../assets/battery/", false, /\.png$/);
            const battery = this.station.status_json.batteryLevel;
            let img = "";
            if (battery == 0) {
                img = "0.png";
            } else if (battery <= 20) {
                img = "20.png";
            } else if (battery <= 40) {
                img = "40.png";
            } else if (battery <= 60) {
                img = "60.png";
            } else if (battery <= 80) {
                img = "80.png";
            } else {
                img = "100.png";
            }
            return imgPath("./" + img);
        },

        getSensorName(module, sensor) {
            return this.$t(module.name + "." + sensor.name);
        },

        getModuleImg(module) {
            let imgPath = require.context("../assets/", false, /\.png$/);
            let img = utils.getModuleImg(module);
            return imgPath("./" + img);
        },

        getLat() {
            if (!this.station.status_json.latitude) {
                return false;
            }
            let lat = parseFloat(this.station.status_json.latitude);
            return lat.toFixed(5);
        },

        getLong() {
            if (!this.station.status_json.longitude) {
                return false;
            }
            let long = parseFloat(this.station.status_json.longitude);
            return long.toFixed(5);
        },

        closeSummary() {
            this.viewingSummary = false;
            if (!this.compact) {
                if (this.$route.name != "stations") {
                    this.$router.push({ name: "stations" });
                }
            }
        },
    },
};
</script>

<style scoped>
#station-summary-container {
    position: relative;
    background-color: #ffffff;
    padding: 0 20px 20px 10px;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
}
.station-container {
    margin: 20px 0 0 0;
    font-size: 14px;
    font-weight: lighter;
    overflow: hidden;
}
.image-container {
    width: 124px;
    height: 100px;
    text-align: center;
}
.image-container img {
    max-width: 124px;
    max-height: 100px;
}
.section {
    width: 100%;
    float: left;
}
.spacer {
    float: left;
    width: 100%;
    margin: 5px 0 10px 10px;
    border-bottom: 1px solid #f1eeee;
    height: 1px;
}
.small-light {
    font-size: 12px;
    color: #6a6d71;
}
.location-item {
    margin: 10px;
    width: 100%;
}
.location-item img {
    float: left;
    margin: 2px 5px 0 0;
}
.native-land-link {
    font-weight: bold;
    text-decoration: underline;
}
.left {
    float: left;
}
.right {
    float: right;
}
.station-element {
    margin: 5px 5px 0 5px;
}
#station-name {
    font-size: 18px;
    font-weight: bold;
}
#battery {
    margin-right: 5px;
    width: 20px;
    height: 11px;
}
.small-space {
    margin: 3px;
}
#location-container .left {
    margin-left: 10px;
}
#readings-label {
    margin: 0 0 10px 10px;
}
.left-reading,
.right-reading {
    width: 42%;
    margin: 8px 0;
    padding: 5px 10px;
    background-color: #f4f5f7;
    border-radius: 2px;
}
.left-reading {
    float: left;
    margin-left: 10px;
}
.right-reading {
    float: right;
    margin-right: 3px;
}
.sensor-name {
    font-size: 11px;
    line-height: 20px;
}
.sensor-reading {
    font-size: 16px;
}
.sensor-unit {
    font-size: 10px;
    margin: 6px 0 0 4px;
}
#view-data-btn {
    width: 90%;
    font-size: 18px;
    font-weight: bold;
    color: #ffffff;
    text-align: center;
    padding: 10px;
    margin: 30px 0 0 10px;
    background-color: #ce596b;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
#close-form-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    cursor: pointer;
}
</style>
