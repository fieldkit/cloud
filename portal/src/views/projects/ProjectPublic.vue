<template>
    <div class="project-public project-container" v-if="project">
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
                    {{ project.name }}
                </div>
                <div class="details-top">
                    <div class="time-container" v-if="project.startTime">
                        <img alt="Location" src="@/assets/icon-calendar.png" class="icon" />
                        <template>Started: {{ project.startTime | prettyDate }}</template>
                    </div>
                    <div class="duration-container" v-if="displayProject.duration">
                        <img alt="Location" src="@/assets/icon-time.png" class="icon" />
                        <template>{{ displayProject.duration | prettyDuration }}</template>
                    </div>
                    <div class="location-container" v-if="project.location">
                        <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                        <template>{{ project.location }}</template>
                    </div>
                    <div class="location-container" v-if="displayProject.places.native">
                        <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                        <template>Native Lands: {{ displayProject.places.native }}</template>
                    </div>
                </div>
                <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                <div class="project-detail">{{ project.description }}</div>
                <div class="details-modules">
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

        <div class="project-stations">
            <ProjectStations :project="project" :admin="false" :userStations="userStations" />
        </div>

        <ProjectActivity :displayProject="displayProject" v-if="false" />
    </div>
</template>

<script>
import * as utils from "../../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";
import ProjectStations from "./ProjectStations";
import ProjectActivity from "./ProjectActivity";
import CommonComponents from "@/views/shared";

export default {
    name: "ProjectPublic",
    components: {
        ...CommonComponents,
        ProjectStations,
        ProjectActivity,
    },
    data: () => {
        return {};
    },
    props: {
        user: {
            required: true,
        },
        userStations: {
            required: true,
        },
        displayProject: {
            required: true,
        },
    },
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
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
}
.details > .left {
    max-width: 400px;
    min-height: 335px;
    flex: 1;
    padding: 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
}
.details .project-detail {
    margin-top: 10px;
    margin-bottom: 10px;
}
.details .below-photo {
    margin-top: auto;
}
.details > .right {
    flex: 2;
    padding: 20px;
    display: flex;
    flex-direction: column;
}
.project-stations {
}

.details .details-heading {
    font-weight: 900;
    padding-bottom: 20px;
    font-size: 24px;
    color: #2c3e50;
}
.details .details-top {
    display: flex;
    flex-direction: column;
}
.details .icon {
    padding-right: 0.2em;
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
    margin-right: 10px;
}

/deep/ .project-image {
    width: 100%;
    height: auto;
}

.project-container {
    margin-top: 20px;
}
</style>
