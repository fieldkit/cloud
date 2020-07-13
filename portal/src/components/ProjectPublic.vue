<template>
    <div class="project-public project-container" v-if="project">
        <div class="header">
            <div class="left">
                <div class="project-name">{{ project.name }}</div>
                <div class="project-dashboard">Project Dashboard</div>
            </div>
            <div class="right activity">Activity</div>
        </div>
        <div class="details">
            <div class="left">
                <div class="photo">
                    <img alt="FieldKit Project" v-if="project.photo" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default FieldKit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
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

            <div class="right">
                <div class="details-heading">
                    Project Details
                </div>
                <div class="details-top">
                    <div class="details-left">
                        <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="details-right">
                        <div class="time-container" v-if="project.startTime">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>Started: {{ project.startTime | prettyDate }}</template>
                        </div>
                        <div class="duration-container" v-if="displayProject.duration">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>{{ displayProject.duration | prettyDuration }}</template>
                        </div>
                        <div class="location-container" v-if="project.location">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>{{ project.location }}</template>
                        </div>
                        <div class="location-container" v-if="displayProject.places.native">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>Native Lands: {{ displayProject.places.native }}</template>
                        </div>
                    </div>
                </div>

                <div class="details-bottom">
                    <div class="details-team">
                        <div class="title">Team</div>
                        <UserPhoto
                            v-for="projectUser in displayProject.users"
                            v-bind:key="projectUser.user.email"
                            :user="projectUser.user"
                        />
                    </div>
                    <div class="details-modules">
                        <div class="title">Modules</div>

                        <img
                            v-for="module in projectModules"
                            v-bind:key="module.name"
                            alt="Module icon"
                            class="module-icon"
                            :src="module.url"
                        />
                    </div>
                </div>
            </div>
        </div>

        <div class="project-stations">
            <ProjectStations
                :project="project"
                :admin="false"
                :mapContainerSize="mapContainerSize"
                :listSize="listSize"
                :userStations="userStations"
            />
        </div>

        <ProjectActivity :displayProject="displayProject" v-if="false" />
    </div>
</template>

<script>
import * as utils from "../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "../api/api";
import ProjectStations from "../components/ProjectStations";
import ProjectActivity from "../components/ProjectActivity";
import CommonComponents from "@/views/shared";

export default {
    name: "ProjectPublic",
    components: {
        ...CommonComponents,
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
            return this.$getters.projectsById[this.displayProject.id].stations;
        },
        projectModules() {
            return this.$getters.projectsById[this.displayProject.id].modules.map((m) => {
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
            return this.$config.baseUrl + project.photo;
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
.project-public {
    display: flex;
    flex-direction: column;
    max-width: 1000px;
    padding-bottom: 100px;
}
.header {
    display: flex;
    flex-direction: row;
}
.header .left {
    margin-right: auto;
    display: flex;
    flex-direction: column;
}
.header .right {
    margin-left: auto;
}
.header .project-name {
    font-size: 24px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
}
.header .project-dashboard {
    font-size: 20px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
    margin-top: 10px;
    margin-bottom: 20px;
}

.details {
    display: flex;
    flex-direction: row;
}
.details > .left {
    flex: 1;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    margin-right: 20px;
    background-color: white;
    padding: 20px;
}
.details > .right {
    flex: 2;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
}
.project-stations {
}

.details .details-heading {
    font-weight: bold;
    padding-bottom: 20px;
}
.details .details-top {
    display: flex;
    flex-direction: row;
    padding-bottom: 20px;
}
.details .details-left {
    padding-right: 20px;
    flex-grow: 1;
}
.details .details-right {
    flex-grow: 1;
}
.details .details-bottom {
    border-top: 1px solid #d8dce0;
    padding-top: 20px;
    display: flex;
    flex-direction: row;
}
.details-bottom .details-team {
    flex: 1;
}
.details-bottom .details-modules {
    flex: 1;
}
.details-bottom .title {
    font-weight: bold;
}

.module-icon {
    width: 35px;
    height: 35px;
}

.project-image {
    max-width: 400px;
    max-height: 400px;
}
</style>
