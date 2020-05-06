<template>
    <div id="project-summary-container" v-if="viewingSummary">
        <div class="project-container" v-if="project">
            <div class="project-profile-container">
                <div class="section left-section">
                    <img alt="Fieldkit Project" v-if="project.media_url" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                </div>
                <div class="section right-section">
                    <div id="project-name">{{ project.name }}</div>
                    <div class="location" v-if="project.location">
                        <img alt="Location" src="../assets/icon-location.png" class="icon" />
                        {{ project.location }}
                    </div>
                    <div class="time-container">
                        <div class="time" v-if="displayStartDate">Started {{ displayStartDate }}</div>
                        <span v-if="displayStartDate && displayRunTime">&nbsp;|&nbsp;</span>
                        <div class="time" v-if="displayRunTime">{{ displayRunTime }}</div>
                    </div>
                    <div class="project-detail">{{ project.description }}</div>
                    <div class="module-icons">
                        <img v-for="module in modules" v-bind:key="module" alt="Module icon" class="module-icon" :src="module" />
                    </div>
                </div>

                <div class="follow-btn">
                    <span v-if="following" v-on:click="unfollowProject">
                        Following
                    </span>
                    <span v-else v-on:click="followProject">
                        <img alt="Follow" src="../assets/heart_gray.png" class="icon" />
                        Follow
                    </span>
                </div>
            </div>

            <div class="stations-container">
                <div class="section-heading stations-heading">FieldKit Stations</div>
                <div class="space"></div>
                <div class="stations-list">
                    <div v-for="station in projectStations" v-bind:key="station.id">
                        <div class="station-box">
                            <span class="station-name" v-on:click="showStation(station)">
                                {{ station.name }}
                            </span>
                            <div class="last-seen">Last seen {{ getUpdatedDate(station) }}</div>
                        </div>
                    </div>
                </div>
                <div class="stations-map">
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
            <div class="team-container">
                <div class="section-heading">{{ getTeamHeading() }}</div>
                <div v-for="user in projectUsers" v-bind:key="user.user.id" class="team-member">
                    <img v-if="user.user.media_url" alt="User image" :src="user.userImage" class="user-icon" />
                    <span class="user-name">{{ user.user.name }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as utils from "../utilities";
import FKApi from "../api/api";
import { API_HOST } from "../secrets";
import Mapbox from "mapbox-gl-vue";
import { MAPBOX_ACCESS_TOKEN } from "../secrets";

export default {
    name: "ProjectProfile",
    components: {
        Mapbox,
    },
    data: () => {
        return {
            baseUrl: API_HOST,
            projectStations: [],
            viewingSummary: false,
            displayStartDate: "",
            displayRunTime: "",
            projectUsers: [],
            timeUnits: ["seconds", "minutes", "hours", "days", "weeks", "months", "years"],
            modules: [],
            coordinates: [-118, 34],
            mapboxToken: MAPBOX_ACCESS_TOKEN,
            following: false,
        };
    },
    props: ["user", "project", "userStations", "users"],
    watch: {
        users() {
            if (this.users) {
                this.projectUsers = this.users.map(u => {
                    u.userImage = this.baseUrl + "/user/" + u.user.id + "/media";
                    return u;
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
            this.reset();
        },
        followProject() {
            this.api.followProject(this.project.id).then(() => {
                this.following = true;
            });
        },
        unfollowProject() {
            this.api.unfollowProject(this.project.id).then(() => {
                this.following = false;
            });
        },
        reset() {
            this.projectStations = [];
            this.fetchStations();
            this.fetchFollowers();
            this.updateDisplayDates();
        },
        mapInitialized(map) {
            this.map = map;
        },
        fetchFollowers() {
            this.api.getProjectFollows(this.project.id).then(result => {
                result.followers.forEach(f => {
                    if (f.id == this.user.id) {
                        this.following = true;
                    }
                });
            });
        },
        fetchStations() {
            this.api.getStationsByProject(this.project.id).then(result => {
                this.projectStations = result.stations;
                this.modules = [];
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
                                this.modules.push(this.getModuleImg(m));
                            });
                        } else if (s.status_json.statusJson && s.status_json.statusJson.modules) {
                            s.status_json.statusJson.modules.forEach(m => {
                                this.modules.push(this.getModuleImg(m));
                            });
                        }
                        // could also use readings, if present
                    });
                    this.modules = _.uniq(this.modules);
                }
            });
        },
        getImageUrl(project) {
            return this.baseUrl + "/projects/" + project.id + "/media";
        },
        updateDisplayDates() {
            this.displayRunTime = "";
            this.displayStartDate = "";
            if (this.project.start_time) {
                let d = new Date(this.project.start_time);
                this.displayStartDate = d.toLocaleDateString("en-US");
                this.getRunTime();
            }
        },
        getUpdatedDate(station) {
            if (!station.status_json) {
                return "N/A";
            }
            const date = station.status_json.updated;
            const d = new Date(date);
            return d.toLocaleDateString("en-US");
        },
        closeSummary() {
            this.viewingSummary = false;
        },
        showStation(station) {
            this.$emit("showStation", station);
        },
        getRunTime() {
            let start = new Date(this.project.start_time);
            let end, runTense;
            if (this.project.end_time) {
                end = new Date(this.project.end_time);
                runTense = "Ran for ";
            } else {
                // assume it's still running?
                end = new Date();
                runTense = "Running for ";
            }
            // get interval and convert to seconds
            const interval = (end.getTime() - start.getTime()) / 1000;
            let displayValue = interval;
            let unit = 0;
            // unit is an index into this.timeUnits
            if (interval < 60) {
                // already set to seconds
            } else if (interval < 3600) {
                // minutes
                unit = 1;
                displayValue /= 60;
                displayValue = Math.round(displayValue);
            } else if (interval < 86400) {
                // hours
                unit = 2;
                displayValue /= 3600;
                displayValue = Math.round(displayValue);
            } else if (interval < 604800) {
                // days
                unit = 3;
                displayValue /= 86400;
                displayValue = Math.round(displayValue);
            } else if (interval < 2628000) {
                // weeks
                unit = 4;
                displayValue /= 604800;
                displayValue = Math.round(displayValue);
            } else if (interval < 31535965) {
                // months
                unit = 5;
                displayValue /= 2628000;
                displayValue = Math.round(displayValue);
            } else {
                // years
                unit = 6;
                displayValue /= 31535965;
                displayValue = Math.round(displayValue);
            }
            this.displayRunTime = runTense + displayValue + " " + this.timeUnits[unit];
        },
        getModuleImg(module) {
            let imgPath = require.context("../assets/modules-lg/", false, /\.png$/);
            let img = utils.getModuleImg(module);
            return imgPath("./" + img);
        },
        getModuleName(module) {
            const newName = utils.convertOldFirmwareResponse(module);
            return this.$t(newName + ".name");
        },
        getTeamHeading() {
            const members = this.projectUsers.length == 1 ? "member" : "members";
            return "Project Team (" + this.projectUsers.length + " " + members + ")";
        },
    },
};
</script>

<style scoped>
#project-summary-container {
    width: 1080px;
    margin: 0 0 0 30px;
    background-color: #ffffff;
    z-index: 2;
}
.project-profile-container {
    float: left;
    width: 820px;
    padding: 20px;
    border: 1px solid #d8dce0;
}
#project-name {
    font-size: 24px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
}
.show-link {
    text-decoration: underline;
}
.project-container {
    font-size: 16px;
    font-weight: lighter;
    overflow: hidden;
}
.project-image {
    max-width: 288px;
    max-height: 139px;
}
.section {
    float: left;
}
.left-section {
    width: 380px;
    text-align: center;
}
.right-section {
    margin-left: 30px;
}
.section-heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 0 0 35px 0;
}
.time {
    font-size: 14px;
    display: inline-block;
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
}
.last-seen {
    font-size: 12px;
    font-weight: 600;
    color: #6a6d71;
}
.project-detail {
    font-size: 16px;
    line-height: 24px;
}
.follow-btn {
    float: right;
    clear: both;
    margin: -20px 0 0 0;
    border: 1px solid #cccdcf;
    border-radius: 3px;
    width: 80px;
    height: 23px;
    font-size: 14px;
    font-weight: bold;
    padding: 5px 5px 0 5px;
    text-align: center;
    cursor: pointer;
}
.follow-btn img {
    vertical-align: middle;
    margin: 0 3px 4px 0;
}
.space {
    width: 100%;
    float: left;
    margin: 30px 0 0 0;
    border-bottom: solid 1px #d8dce0;
}
.team-icons,
.module-icons {
    width: 225px;
    margin: 10px 0;
    float: left;
}
.user-icon,
.module-icon {
    width: 35px;
    margin: 0 5px;
    float: left;
}

.stations-container {
    width: 860px;
    float: left;
    margin: 22px 0 0 0;
    border: 1px solid #d8dce0;
}
.stations-list {
    width: 320px;
    height: 332px;
    float: left;
}
.stations-map {
    width: 540px;
    height: 332px;
    float: left;
}
#map {
    width: 540px;
    height: 332px;
}
.station-box {
    width: 250px;
    height: 38px;
    margin: 20px auto;
    padding: 10px;
    border: 1px solid #d8dce0;
}

.team-container {
    width: 310px;
    margin: 22px 0 0 0;
    padding: 20px;
    border: 1px solid #d8dce0;
    float: left;
    clear: both;
}
.team-member {
    float: left;
    clear: both;
}
.user-name {
}

#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
