<template>
    <div id="project-summary-container">
        <div class="project-container" v-if="project">
            <div class="project-profile-container">
                <div class="section left-section">
                    <img alt="Fieldkit Project" v-if="project.mediaUrl" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                </div>
                <div class="section right-section">
                    <div id="project-name">{{ project.name }}</div>
                    <div class="location" v-if="project.location">
                        <img alt="Location" src="../assets/icon-location.png" class="icon" />
                        {{ project.location }}
                    </div>
                    <div class="location" v-if="displayProject.places.native">
                        <img alt="Location" src="../assets/icon-location.png" class="icon" />
                        Native Lands: {{ displayProject.places.native }}
                    </div>
                    <div class="time-container">
                        <div class="time" v-if="project.startTime">Started: {{ project.startTime | prettyDate }}</div>
                        <span v-if="project.startTime && displayProject.duration">&nbsp;|&nbsp;</span>
                        <div class="time" v-if="displayProject.duration">{{ displayProject.duration | prettyDuration }}</div>
                    </div>
                    <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                    <div class="project-detail">{{ project.description }}</div>
                    <div class="module-icons">
                        <img
                            v-for="module in projectModules"
                            v-bind:key="module.name"
                            alt="Module icon"
                            class="module-icon"
                            :src="module.url"
                        />
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

            <div id="recent-update" v-if="mostRecentUpdate">
                <div class="activity-icon">
                    <img :src="mostRecentUpdate.icon" />
                </div>
                <div class="activity-heading">Project Update</div>
                <div class="activity-byline">by {{ mostRecentUpdate.name }} | {{ mostRecentUpdate.time.toLocaleDateString() }}</div>
                <div class="activity-text">
                    {{ mostRecentUpdate.text }}
                </div>
            </div>

            <ProjectStations
                :project="project"
                :admin="false"
                :mapContainerSize="mapContainerSize"
                :listSize="listSize"
                :userStations="userStations"
            />

            <div class="team-container">
                <div class="section-heading">{{ getTeamHeading() }}</div>
                <div v-for="user in displayProject.users" v-bind:key="user.user.id" class="team-member">
                    <img v-if="user.user.mediaUrl" alt="User image" :src="getUserImage(user)" class="user-icon" />
                    <span class="user-name">{{ user.user.name }}</span>
                </div>
            </div>
            <div id="public-activity-feed-container">
                <div class="heading">Recent Activity</div>
                <ProjectActivity :displayProject="displayProject" :viewing="true" />
            </div>
        </div>
    </div>
</template>

<script>
import * as utils from "../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "../api/api";
import ProjectStations from "../components/ProjectStations";
import ProjectActivity from "../components/ProjectActivity";

export default {
    name: "ProjectPublic",
    components: {
        ProjectStations,
        ProjectActivity,
    },
    data: () => {
        return {
            displayStartDate: "",
            displayRunTime: "",
            mostRecentUpdate: null,
            following: false,
            mapContainerSize: {
                width: "540px",
                height: "332px",
                outerWidth: "860px",
            },
            listSize: {
                width: "320px",
                height: "332px",
                boxWidth: "250px",
            },
        };
    },
    props: { user: {}, userStations: {}, displayProject: {} },
    computed: {
        project() {
            return this.displayProject.project;
        },
        projectStations() {
            return this.$store.getters.projectsById[this.displayProject.id].stations;
        },
        projectModules() {
            return this.$store.getters.projectsById[this.displayProject.id].modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
    },
    methods: {
        async followProject() {
            await this.$store.dispatch(ActionTypes.PROJECT_FOLLOW, { projectId: this.project.id });
            this.following = true;
        },
        async unfollowProject() {
            await this.$store.dispatch(ActionTypes.PROJECT_UNFOLLOW, { projectId: this.project.id });
            this.following = false;
        },
        getImageUrl(project) {
            return this.$config.baseUrl + "/projects/" + project.id + "/media";
        },
        getUserImage(projectUser) {
            return this.$config.baseUrl + "/user/" + projectUser.user.id + "/media";
        },
        getModuleImg(module) {
            return this.$loadAsset("modules-lg/" + utils.getModuleImg(module));
        },
        getTeamHeading() {
            const members = this.displayProject.users.length == 1 ? "member" : "members";
            return "Project Team (" + this.displayProject.users.length + " " + members + ")";
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
.project-detail {
    overflow-wrap: break-word;
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
#recent-update {
    width: 820px;
    float: left;
    margin: 22px 0 0 0;
    padding: 20px;
    border-radius: 2px;
    border: 1px solid #d8dce0;
}
.activity-icon {
    float: left;
    margin-left: 40px;
}
.activity-icon img {
    width: 35px;
}
.activity-heading {
    float: left;
    margin-left: 14px;
    font-size: 20px;
    font-weight: 600;
}
.activity-byline {
    float: left;
    clear: both;
    margin: -14px 0 0 90px;
    font-size: 16px;
}
.activity-text {
    float: left;
    clear: both;
    margin-left: 90px;
    font-size: 16px;
    line-height: 24px;
    margin-top: 9px;
}
.team-icons,
.module-icons {
    width: 225px;
    margin: 10px 0;
    float: left;
}
.user-icon,
.module-icon {
    width: 32px;
    margin: 2px 5px;
    float: left;
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
#public-activity-feed-container {
    float: left;
    width: 480px;
    height: 312px;
    margin: 22px 0 0 32px;
    border: 1px solid #d8dce0;
    background: white;
}
#public-activity-feed-container .heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 20px;
}
</style>
