<template>
    <div id="station-summary-container" v-if="viewingSummary">
        <div v-if="station" id="close-form-btn" v-on:click="closeSummary">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <div v-if="!isAuthenticated" class="no-stations-message">
            <p>
                Please
                <router-link :to="{ name: 'login' }" class="show-link">
                    log in
                </router-link>
                to view stations.
            </p>
        </div>
        <div v-if="isAuthenticated && !station" class="no-stations-message">
            <p>No stations added yet.</p>
        </div>
        <div class="station-container" v-if="station">
            <div class="left">
                <img style="width: 124px; height: 100px;" alt="Station image" :src="stationSmallPhoto" class="station-element" />
            </div>
            <div class="left">
                <div id="station-name" class="station-element">{{ station.name }}</div>
                <div class="station-element">Last Synced {{ getSyncedDate() }}</div>
                <div class="station-element">
                    <img id="battery" alt="Battery level" :src="getBatteryImg()" />
                    <span>{{ station.status_json.batteryLevel }}%</span>
                </div>
                <div>
                    <img
                        v-for="module in station.status_json.moduleObjects"
                        v-bind:key="module.id"
                        alt="Module icon"
                        class="small-space"
                        :src="getModuleImg(module)"
                    />
                </div>
            </div>
            <div class="spacer"></div>
            <div id="location-container" class="section">
                <div>{{ this.station.status_json.locationName }}</div>
                <div class="left">
                    {{ getLat() || "--" }}
                    <br />
                    Latitude
                </div>
                <div class="right">
                    {{ getLong() || "--" }}
                    <br />
                    Longitude
                </div>
            </div>
            <div class="spacer"></div>
            <div id="readings-container" class="section">
                <div id="readings-label">Latest Reading</div>
                <div v-for="module in station.status_json.moduleObjects" v-bind:key="module.id">
                    <div v-for="sensor in module.sensorObjects" v-bind:key="sensor.id" class="reading">
                        <div class="left">{{ sensor.name }}</div>
                        <div class="right">
                            {{ sensor.currentReading ? sensor.currentReading.toFixed(1) : "ncR" }}
                        </div>
                    </div>
                </div>
            </div>
            <router-link :to="{ name: 'viewData', params: { id: station.id } }">
                <div id="view-data-btn" class="section">
                    View Data
                </div>
            </router-link>
        </div>
    </div>
</template>

<script>
import * as utils from "../utilities";
import { API_HOST } from "../secrets";

export default {
    name: "StationSummary",
    data: () => {
        return {
            viewingSummary: false,
        };
    },
    props: ["station", "isAuthenticated"],
    computed: {
        stationSmallPhoto: function() {
            return API_HOST + this.station.photos.small;
        },
    },
    methods: {
        viewSummary() {
            this.viewingSummary = true;
        },

        getSyncedDate() {
            return utils.getUpdatedDate(this.station);
        },

        getBatteryImg() {
            const images = require.context("../assets/battery/", false, /\.png$/);
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
            return images("./" + img);
        },

        getModuleImg(module) {
            let images = require.context("../assets/", false, /\.png$/);
            let img = "placeholder.png";
            // Note: this is not a trustworthy way of figuring out what icons to show,
            // as the user could rename their module anything
            if (module.name.indexOf("Water") > -1) {
                img = "water.png";
            }
            if (module.name.indexOf("Weather") > -1) {
                img = "weather.png";
            }
            if (module.name.indexOf("Ocean") > -1) {
                img = "ocean.png";
            }
            return images("./" + img);
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
            if (this.$route.name != "stations") {
                this.$router.push({ name: "stations" });
            }
        },
    },
};
</script>

<style scoped>
#station-summary-container {
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
.no-stations-message {
    font-size: 20px;
}
.show-link {
    text-decoration: underline;
}
.station-container {
    padding: 10px;
    margin: 20px 0;
    font-size: 14px;
    font-weight: lighter;
    overflow: hidden;
}
.section {
    width: 100%;
    float: left;
}
.spacer {
    float: left;
    width: 100%;
    margin: 20px 0;
    border-bottom: 1px solid rgb(215, 220, 225);
    height: 1px;
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
}
.small-space {
    margin: 3px;
}
#location-container {
    width: 50%;
}
#location-container .left {
    margin-left: 10px;
}
#readings-label {
    margin-bottom: 10px;
}
.reading {
    width: 40%;
    float: left;
    margin: 5px 10px 5px 0;
    padding: 5px 10px;
    background-color: rgb(233, 233, 233);
    border: 1px solid rgb(200, 200, 200);
}
#view-data-btn {
    width: 90%;
    font-size: 20px;
    text-align: center;
    padding: 10px;
    margin: 20px 0;
    background-color: rgb(243, 243, 243);
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
