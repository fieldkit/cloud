<template>
    <div class="project-public project-container" v-if="project">
        <div class="details">
            <div class="left">
                <div class="photo-container">
                    <ProjectPhoto :project="project" />
                </div>

                <DisplayProjectTags :tags="project.tags" />
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
                            <img alt="Location" src="@/assets/icon-time.svg" class="icon" width="12px" height="14px" />
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
                            <img alt="Location" src="@/assets/icon-location.svg" class="icon" width="13px" height="16px" />
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
                        <img src="@/assets/icon-heart-dark-blue.svg" width="16px" alt="Icon" />
                        Follow
                    </button>
                    <button class="button-social">
                        <img src="@/assets/icon-share.svg" width="12px" alt="Icon" />
                        Share
                    </button>
                </div>
            </div>
        </div>

        <div class="update todo-disabled">
            <UserPhoto :user="user" />
            <div>
                <h1>Project update title</h1>
                <h2>by Lauren Martin | 01/21/2020</h2>
                <p>
                    Thank you so much for contributing and supporting us during this project. Lorem ipsum dolor sit amet, consectetur
                    adipiscing elit. Sedo eiusmod tempor incididunt ut labore et dolore magna aliqua.
                </p>
            </div>
            <button class="button-solid">Get involved</button>
        </div>

        <div class="project-stations">
            <ProjectStations :project="project" :admin="false" :userStations="userStations" />
        </div>

        <div class="project-team-activity todo-disabled">
            <div class="project-team">
                <h1>Project Team (4 members)</h1>
                <ul>
                    <li>
                        <img src="@/assets/image-placeholder.svg" alt="Team Member Avatar" />
                        <div>
                            <div>Lauren Martin</div>
                            <div class="project-team-role">Project Lead</div>
                        </div>
                    </li>
                    <li>
                        <img src="@/assets/image-placeholder.svg" alt="Team Member Avatar" />
                        <div>
                            <div>Lauren Martin</div>
                            <div class="project-team-role">Project Lead</div>
                        </div>
                    </li>
                </ul>
            </div>

            <div class="recent-activity">
                <h1>Recent Activity</h1>
                <ul>
                    <li>
                        <img src="@/assets/icon-compass.svg" width="30" />
                        <div>
                            <div>
                                <h2>Downloaded New Data</h2>
                                <span>12/12/2020</span>
                            </div>
                            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. sed do eiusmod tempor</p>
                        </div>
                    </li>
                    <li>
                        <img src="@/assets/icon-compass.svg" width="30" />
                        <div>
                            <div>
                                <h2>Downloaded New Data</h2>
                                <span>12/12/2020</span>
                            </div>
                            <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. sed do eiusmod tempor</p>
                        </div>
                    </li>
                </ul>
            </div>
        </div>

        <Comments :parentData="displayProject.id" :user="user"> </Comments>

    </div>
</template>

<script>
import * as utils from "../../utilities";
import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";
import ProjectStations from "./ProjectStations";
import ProjectActivity from "./ProjectActivity";
import CommonComponents from "@/views/shared";
import Comments from "../shared/Comments";

export default {
    name: "ProjectPublic",
    components: {
        Comments,
        ...CommonComponents,
        ProjectStations,
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
@import "../../scss/project";
@import "../../scss/global";

.project-public {
    display: flex;
    flex-direction: column;

    @include bp-down($sm) {
        padding-bottom: 20px;
    }
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
    font-family: $font-family-bold;
    margin: 0 15px 0 0;
    display: inline-block;
}
.header .project-dashboard {
    font-size: 20px;
    font-family: $font-family-bold;
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

    @include bp-down($sm) {
        flex-wrap: wrap;
    }
}
.details > .left {
    flex: 1;
    padding: 23px 20px;
    display: flex;
    flex-direction: column;

    @include bp-down($sm) {
        flex-basis: 100%;
        padding-bottom: 0;
    }

    @include bp-down($xs) {
        padding: 10px 10px 0;
    }
}
.details .project-detail {
    font-family: $font-family-light;
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

.details > .right {
    padding: 20px 20px 67px 20px;
    position: relative;

    @include bp-down($xs) {
        padding-top: 14px;
    }

    @include bp-down($xs) {
        padding: 10px 10px 64px;
    }
}
.details .details-heading {
    font-family: $font-family-bold;
    padding-bottom: 10px;
    font-size: 24px;
    color: #2c3e50;

    @include bp-down($sm) {
        font-size: 22px;
    }
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
    font-family: $font-family-bold;
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
    margin-top: -10px;

    @include bp-down($sm) {
        margin-top: -28px;
    }
}

.details-icon {
    width: 23px;
}

.right-actions {
    @include position(absolute, null 14px 18px null);
    @include flex();

    @include bp-down($xs) {
        right: 0;
        bottom: 12px;
    }

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

    @include bp-down($md) {
        padding: 26px 23px 23px;
    }

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

        @include bp-down($md) {
            padding-left: 52px;
        }

        @include bp-down($xs) {
            padding-left: 42px;
        }
    }

    h2 {
        font-size: 16px;
        line-height: 1.3;
        margin: 0;
        font-weight: 300;

        @include bp-down($md) {
            padding-left: 52px;
        }

        @include bp-down($xs) {
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

.recent-activity {
    border-radius: 2px;
    border: solid 1px #d8dce0;
    background-color: #ffffff;
    padding: 17px 23px;
    flex: 1;

    @include bp-down($sm) {
        padding: 19px 10px;
    }

    h1 {
        margin: 0 0 23px;
        font-size: 20px;
        font-weight: 500;

        @include bp-down($sm) {
            font-size: 18px;
        }
    }

    h2 {
        margin: 0;
        font-size: 16px;
        font-weight: 500;

        + span {
            margin-left: auto;
            padding-left: 10px;
            font-weight: 300;
            font-size: 14px;
            color: #6a6d71;
        }
    }

    li {
        @include flex(flex-start);
        font-size: 16px;
        margin-bottom: 30px;

        > div {
            display: flex;
            flex-direction: column;
            flex: 1;
            line-height: 1.3;

            > div {
                display: flex;
            }
        }

        img {
            margin-right: 16px;
            margin-top: 2px;
        }

        p {
            margin-top: 5px;
            font-size: 14px;
        }
    }
}

.project-team-activity {
    @include flex();
    margin-top: 21px;

    @include bp-down($sm) {
        flex-wrap: wrap;
    }
}

.project-team {
    flex-basis: 349px;
    padding: 17px 20px;
    border-radius: 2px;
    border: solid 1px #d8dce0;
    background-color: #ffffff;
    margin-right: 32px;

    @include bp-down($sm) {
        flex-basis: 100%;
        margin-right: 0;
        margin-bottom: 26px;
    }

    @include bp-down($sm) {
        padding: 19px 10px;
    }

    h1 {
        font-size: 20px;
        font-weight: 500;
        margin: 0 0 23px;

        @include bp-down($sm) {
            font-size: 18px;
        }
    }

    li {
        padding-left: 4px;
        font-size: 14px;
        line-height: 1.5;
        margin-bottom: 15px;
        @include flex(center);

        img {
            width: 35px;
            height: 35px;
            border-radius: 50%;
            object-fit: cover;
            margin-right: 14px;
        }
    }

    &-role {
        color: #818181;
        font-size: 13px;
        line-height: 1.2;
    }
}

.todo-disabled {
    display: none;
}
</style>
