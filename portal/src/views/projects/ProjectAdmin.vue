<template>
    <div class="project-admin project-container" v-if="project">
        <div class="details">
            <div class="left">
                <div class="photo">
                    <ProjectPhoto :project="project" />
                </div>

                <div class="below-photo">
                    <DisplayProjectTags :tags="project.tags" />
                    <FollowPanel :project="project" v-bind:key="project.id" />
                </div>
            </div>

            <div class="right">
                <div class="details-heading">
                    <div class="title">Project Details</div>
                    <div v-on:click="editProject" class="link">Edit Project</div>
                </div>
                <div class="details-top">
                    <div class="details-left">
                        <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="details-right">
                        <div class="time-container" v-if="project.startTime">
                            <div class="details-icon-container">
                                <img alt="Location" src="@/assets/icon-calendar.png" class="icon" />
                            </div>
                            <template>Started: {{ project.startTime | prettyDate }}</template>
                        </div>
                        <div class="duration-container" v-if="displayProject.duration">
                            <div class="details-icon-container">
                                <img alt="Location" src="@/assets/icon-time.png" class="icon" />
                            </div>
                            <template>{{ displayProject.duration | prettyDuration }}</template>
                        </div>
                        <div class="location-container" v-if="project.location">
                            <div class="details-icon-container">
                                <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            </div>
                            <template>{{ project.location }}</template>
                        </div>
                        <div class="location-container" v-if="displayProject.places.native">
                            <div class="details-icon-container">
                                <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            </div>
                            <template>Native Lands: {{ displayProject.places.native }}</template>
                        </div>
                    </div>
                </div>

                <div class="details-bottom">
                    <div class="details-team">
                        <div class="title">Team</div>
                        <template v-for="projectUser in displayProject.users">
                            <UserPhoto :user="projectUser.user" v-if="!projectUser.invited" v-bind:key="projectUser.user.email" />
                        </template>
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

        <div class="row-section project-stations">
            <ProjectStations :project="project" :admin="true" :userStations="userStations" />
        </div>

        <div class="row-section data-readings">
            <div class="project-data">
                <ProjectDataFiles :projectStations="displayProject.stations" />
            </div>
            <div class="project-readings">
                <StationsReadings :project="project" :displayProject="displayProject" />
            </div>
        </div>

        <TeamManager :displayProject="displayProject" v-bind:key="displayProject.id" />

        <ProjectActivity :displayProject="displayProject" v-if="false" />
    </div>
</template>

<script>
import CommonComponents from "@/views/shared";
import ProjectStations from "./ProjectStations.vue";
import ProjectActivity from "./ProjectActivity.vue";
import ProjectDataFiles from "./ProjectDataFiles.vue";
import StationsReadings from "./StationsReadings.vue";
import TeamManager from "./TeamManager.vue";

import * as utils from "../../utilities";

export default {
    name: "ProjectAdmin",
    components: {
        ...CommonComponents,
        ProjectStations,
        ProjectActivity,
        ProjectDataFiles,
        StationsReadings,
        TeamManager,
    },
    props: {
        displayProject: {
            type: Object,
            required: true,
        },
        userStations: {
            type: Array,
            required: true,
        },
    },
    data: () => {
        return {
            viewingActivityFeed: false,
        };
    },
    computed: {
        project() {
            return this.displayProject.project;
        },
        projectModules() {
            return this.displayProject.modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
    },
    methods: {
        getProjectUserImage(projectUser) {
            if (projectUser.user.photo) {
                return this.$config.baseUrl + "/" + projectUser.user.photo.url;
            }
            return null;
        },
        openProjectNotes() {
            return this.$router.push({ name: "viewProjectNotes", params: { projectId: this.project.id } });
        },
        editProject() {
            return this.$router.push({ name: "editProject", params: { id: this.project.id } });
        },
        addUpdate() {
            return this.$router.push({ name: "addProjectUpdate", params: { project: this.project } });
        },
        viewProfile() {
            return this.$emit("viewProfile");
        },
        closeActivityFeed() {
            this.viewingActivityFeed = false;
        },
        openActivityFeed() {
            this.viewingActivityFeed = true;
        },
        getModuleImg(module) {
            return this.$loadAsset("modules-lg/" + utils.getModuleImg(module));
        },
    },
};
</script>

<style scoped lang="scss">
@import '../../scss/mixins';

.project-admin {
    display: flex;
    flex-direction: column;
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
    max-width: 400px;
    min-height: 335px;
    flex: 1;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    margin-right: 20px;
    background-color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
}
.details .photo {
    display: flex;
    flex-direction: column;
}
.details .below-photo {
    margin-top: auto;
}
.details > .right {
    flex: 2;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
    padding: 20px 30px;
    display: flex;
    flex-direction: column;
}
.project-stations {
}

.row-section {
}

.details .details-heading {
    display: flex;
    flex-direction: row;
}
.details .details-heading .title {
    padding-bottom: 30px;
    font-size: 20px;
    font-weight: 500;
}
.details .details-heading .link {
    margin-left: auto;
}
.link {
    font-weight: bold;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.details .details-top {
    display: flex;
    flex-direction: row;
    padding-bottom: 20px;
    line-height: 1.5;
}
.details .details-left {
    padding-right: 20px;
    flex-grow: 1;
    flex: 2;
}
.details .details-right {
    flex-grow: 1;
    flex: 1;

    > div {
        margin-bottom: 5px;
        @include flex(center);
    }
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
    font-weight: 500;
    font-size: 14px;
}
.details-icon-container {
    width: 20px;
}

.row-section.data-readings {
    margin-top: 20px;
    display: flex;
    flex-direction: row;
}
.project-data {
    margin-right: 20px;
}
.project-data,
.project-readings {
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
}
.data-readings .project-data {
    flex: 2;
}
.data-readings .project-readings {
    flex: 1;
}

.project-container {
    margin-top: 18px;
}
::v-deep .project-image {
    width: 100%;
    height: auto;
}
.module-icon {
    width: 35px;
    height: 35px;
    margin: 6px 7px 0 0;
}
.project-detail {
    &:nth-of-type(1) {
        padding-bottom: 6px;
    }
}

::v-deep .default-user-icon {
    margin: 6px 7px 0 0;
    width: 35px;
    height: 35px;
    border-radius: 50%;
}
</style>
