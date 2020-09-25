<template>
    <div class="project-public project-container" v-if="project">
        <div class="details">
            <div class="left">
                <div class="photo-container">
                    <ProjectPhoto :project="project" />
                </div>

                <div class="below-photo">
                    <DisplayProjectTags :tags="project.tags" />
                </div>
            </div>

            <div class="right">
                <div class="details-heading">
                    {{ project.name }}
                </div>
                <div class="details-top">
                    <div class="details-row" v-if="project.startTime">
                        <div class="details-icon">
                            <img alt="Location" src="@/assets/icon-calendar.svg" class="icon" width="12px" height="14px" />
                        </div>
                        <template>Started: {{ project.startTime | prettyDate }}</template>
                    </div>
                    <div class="details-row" v-if="displayProject.duration">
                        <div class="details-icon">
                            <img alt="Location" src="@/assets/icon-time.svg" class="icon" width="12px" height="14px"/>
                        </div>
                        <template>{{ displayProject.duration | prettyDuration }}</template>
                    </div>
                    <div class="details-row" v-if="project.location">
                        <div class="details-icon">
                            <img alt="Location" src="@/assets/icon-location.svg" class="icon" width="13px" height="16px" />
                        </div>
                        <template>{{ project.location }}</template>
                    </div>
                    <div class="details-row" v-if="displayProject.places.native">
                        <div class="details-icon">
                            <img alt="Location" src="@/assets/icon-location.svg" class="icon"  width="13px" height="16px" />
                        </div>
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
                <div class="right-actions">
                    <button class="button-social">
                        <img src="@/assets/icon-heart-dark-blue.svg" width="16px" alt="Icon">
                        Follow
                    </button>
                    <button class="button-social">
                        <img src="@/assets/icon-share.svg" width="12px" alt="Icon">
                        Share
                    </button>
                </div>
            </div>
        </div>

        <div class="update">
            <UserPhoto :user="user" />
            <div>
                <h1> Project update title </h1>
                <h2> by Lauren Martin | 01/21/2020 </h2>
                <p>
                    Thank you so much for contributing and supporting us during this project. Lorem ipsum dolor sit amet,
                    consectetur adipiscing elit. Sedo eiusmod tempor incididunt ut labore et dolore magna aliqua.
                </p>
            </div>
            <button class="button-solid"> Get involved </button>
        </div>

        <div class="project-stations" style="display: none">
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
            return this.$loadAsset(utils.getModuleImg(module));
        },
        getTeamHeading() {
            const members = this.displayProject.users.length == 1 ? "member" : "members";
            return "Project Team (" + this.displayProject.users.length + " " + members + ")";
        },
    },
};
</script>

<style scoped lang="scss">
@import '../../scss/project';
@import '../../scss/global';

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
    border-radius: 2px;
    border: solid 1px #d8dce0;
    background-color: white;
    flex-wrap: wrap;
}
.details > .left {
    flex: 1;
    padding: 23px 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    flex-basis: 100%;
}
.details .project-detail {
    margin-bottom: 9px;
    line-height: 1.5;
}

.photo-container {
    width: 380px;
    max-height: 234px;

    @include bp-down($xs) {
        width: 100%;
    }
}

.details .below-photo {
    margin-top: auto;
}
.details > .right {
    padding: 20px 20px 67px 20px;
    position: relative;
}
.details .details-heading {
    font-weight: 600;
    padding-bottom: 10px;
    font-size: 24px;
    color: #2c3e50;
}
.details .details-top {
    display: flex;
    flex-direction: column;
    padding-bottom: 10px;
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
.details-modules {
    margin-top: 4px;
}
.details-bottom .title {
    font-weight: bold;
}

.module-icon {
    width: 40px;
    height: 40px;
    margin-right: 10px;
}

::v-deep .project-image {
    width: 100%;
    height: auto;
}

.project-container {
    margin-top: 20px;
}

.details-icon {
    width: 23px;
}

.right-actions {
    @include position(absolute, null 14px 18px null);
    @include flex();

    button {
        margin-right: 10px;
    }
}

.update {
    padding: 26px 29px 23px;
    border-radius: 2px;
    border: solid 1px #d8dce0;
    background-color: #ffffff;
    margin-top: 16px;
    position: relative;
    @include flex();

    @include bp-down($sm) {
        flex-wrap: wrap;
    }

    @include bp-down($xs) {
        padding: 16px 10px 23px;
    }

    .default-user-icon {
        width: 35px;
        height: 35px;
        margin: 4px 14px 0;

        @include bp-down($md) {
            margin-left: 0;
        }

        @include bp-down($sm) {
            margin: 4px 7px 0 0;

            @include position(absolute, 26px null null 29px);
        }

        @include bp-down($xs) {
            @include position(absolute, 16px null null 10px);
        }
    }

    h1 {
        font-size: 20px;
        font-weight: 500;
        margin: 0;
        line-height: 1.4;

        @include bp-down($sm) {
            padding-left: 42px;
        }
    }

    h2 {
        font-size: 16px;
        line-height: 1.3;
        margin: 0;
        font-weight: 300;

        @include bp-down($sm) {
            padding-left: 42px;
        }
    }

    p {
        line-height: 1.5;
        margin: 9px 54px 0 0;

        @include bp-down($md) {
            margin-right: 20px;
        }

        @include bp-down($sm) {
            flex-basis: 100%;
            margin: 11px 0 15px 0;
        }
    }

    .button-solid {
        padding: 0 45px;
        white-space: nowrap;
        align-self: center;

        @include bp-down($xs) {
            width: 100%;
        }
    }
}

</style>
