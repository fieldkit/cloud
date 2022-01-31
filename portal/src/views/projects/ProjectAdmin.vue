<template>
    <div class="project-admin project-container" v-if="project">
        <div class="details">
            <div class="left">
                <div class="photo-container">
                    <ProjectPhoto :project="project" />
                </div>

                <DisplayProjectTags :tags="project.tags" />
                <FollowPanel :project="project" v-bind:key="project.id" />
            </div>

            <div class="right">
                <div class="details-heading">
                    <div class="title">{{ $t("project.details.title") }}</div>
                    <div v-on:click="editProject" class="link">{{ $t("project.edit.link") }}</div>
                </div>
                <div class="details-top">
                    <div class="details-left">
                        <div class="project-detail" v-if="project.goal">{{ $t("project.goal", { goal: project.goal }) }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="details-right">
                        <div class="details-row" v-if="!project.startTime">
                            <div class="details-icon-container">
                                <i class="icon icon-calendar"></i>
                            </div>
                            <template>{{ $t("project.started", { started: project.startTime }) }}</template>
                        </div>
                        <div class="details-row" v-if="!displayProject.duration">
                            <div class="details-icon-container">
                                <i class="icon icon-time" role="img" aria-label="duration"></i>
                            </div>
                            <template>{{ $t("project.duration", { duration: displayProject.duration }) }}</template>
                        </div>
                        <div class="details-row location-name" v-if="project.location" width="12px" height="14px">
                            <div class="details-icon-container">
                                <i class="icon icon-location" role="img" aria-label="location"></i>
                            </div>
                            <template>{{ $t("project.location", { location: project.location }) }}</template>
                        </div>
                        <div class="details-row location-native" v-if="displayProject.places.native">
                            <div class="details-icon-container">
                                <i class="icon icon-location"></i>
                            </div>
                            <template>{{ $t("project.nativeLands", { nativeLands: displayProject.places.native }) }}</template>
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

        <div class="row-section data-readings" v-if="false">
            <div class="project-data">
                <ProjectDataFiles :projectStations="displayProject.stations" />
            </div>
            <div class="project-readings">
                <StationsReadings :project="project" :displayProject="displayProject" />
            </div>
        </div>

        <Comments :parentData="displayProject.id" :user="user"></Comments>

        <TeamManager :displayProject="displayProject" v-bind:key="displayProject.id" />
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { CurrentUser, ProjectModule, DisplayStation, Project, DisplayProject, ProjectUser } from "@/store";
import CommonComponents from "@/views/shared";
import ProjectStations from "./ProjectStations.vue";
import ProjectActivity from "./ProjectActivity.vue";
import ProjectDataFiles from "./ProjectDataFiles.vue";
import StationsReadings from "./StationsReadings.vue";
import Comments from "../comments/Comments.vue";
import TeamManager from "./TeamManager.vue";
import * as utils from "../../utilities";
import { twitterCardMeta } from "@/social";

export default Vue.extend({
    name: "ProjectAdmin",
    components: {
        ...CommonComponents,
        ProjectStations,
        ProjectDataFiles,
        StationsReadings,
        TeamManager,
        Comments,
    },
    metaInfo() {
        return {
            meta: twitterCardMeta(this.displayProject),
            afterNavigation() {
                console.log("hello: after-navigation");
            },
        };
    },
    props: {
        user: {
            type: Object as PropType<CurrentUser>,
            required: true,
        },
        displayProject: {
            type: Object as PropType<DisplayProject>,
            required: true,
        },
        userStations: {
            type: Array as PropType<DisplayStation[]>,
            required: true,
        },
    },
    data: (): { viewingActivityFeed: boolean } => {
        return {
            viewingActivityFeed: false,
        };
    },
    computed: {
        project(): Project {
            return this.displayProject.project;
        },
        projectModules(this: any): { name: string; url: string }[] {
            return this.displayProject.modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
    },
    methods: {
        getProjectUserImage(projectUser: ProjectUser): string | null {
            if (projectUser.user.photo) {
                return this.$config.baseUrl + "/" + projectUser.user.photo.url;
            }
            return null;
        },
        openProjectNotes(this: any): void {
            this.$router.push({ name: "viewProjectNotes", params: { projectId: `${this.project.id}` } });
        },
        editProject(this: any): void {
            this.$router.push({ name: "editProject", params: { id: `${this.project.id}` } });
        },
        addUpdate(this: any): void {
            this.$router.push({ name: "addProjectUpdate", params: { project: `${this.project}` } });
        },
        viewProfile(): void {
            this.$emit("viewProfile");
        },
        closeActivityFeed(): void {
            this.viewingActivityFeed = false;
        },
        openActivityFeed(): void {
            this.viewingActivityFeed = true;
        },
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/project";

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
    font-family: var(--font-family-bold);
    margin: 0 15px 0 0;
    display: inline-block;
}
.header .project-dashboard {
    font-size: 20px;
    font-family: var(--font-family-bold);
    margin: 0 15px 0 0;
    display: inline-block;
    margin-top: 10px;
    margin-bottom: 20px;
}

.details {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
}
.details > .left {
    max-width: 400px;
    border: 1px solid var(--color-border);
    border-radius: 1px;
    margin-right: 20px;
    background-color: #fff;
    padding: 25px;
    flex-direction: column;
    justify-content: space-evenly;
    @include flex();

    @include bp-down($md) {
        max-width: unset;
    }

    @include bp-down($xs) {
        flex: 1;
        margin: 0;
        padding: 15px 10px;
    }
}
.details > .photo {
    display: flex;
    flex-direction: column;
    margin-bottom: 10px;
}

.photo-container {
    width: 288px;
    max-height: 139px;
    align-self: center;

    @include bp-down($sm) {
        max-height: 150px;
    }

    @include bp-down($xs) {
        max-height: 150px;
    }
}

.details > .right {
    border: 1px solid var(--color-border);
    border-radius: 1px;
    background-color: white;
    padding: 20px 30px;

    @include bp-down($md) {
        flex-basis: 100%;
        margin-top: 25px;
    }

    @include bp-down($sm) {
        flex-basis: 100%;
        padding: 20px 20px;
    }

    @include bp-down($xs) {
        padding: 14px 10px 20px;
    }
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

    @include bp-down($sm) {
        padding-boottom: 25px;
    }

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
}
.details .details-heading .link {
    margin-left: auto;
}
.link {
    font-family: var(--font-family-bold);
    font-size: 14px;
    cursor: pointer;

    body.floodnet & {
        font-family: $font-family-floodnet-button;
    }
}

.details .details-bottom {
    border-top: 1px solid var(--color-border);
    padding-top: 20px;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;

    @include bp-down($xs) {
        padding-top: 15px;
    }
}
.details-bottom .details-team {
    flex: 1;

    @include bp-down($xs) {
        flex-basis: 100%;
        margin-bottom: 15px;
    }
}
.details-bottom .details-modules {
    flex: 1;
}
.details-bottom .title {
    font-weight: 500;
    font-size: 14px;

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
}
.details-icon-container {
    width: 20px;
    display: flex;
    flex-shrink: 0;
}

.row-section.data-readings {
    margin-top: 25px;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
}
.project-data,
.project-readings {
    border: 1px solid var(--color-border);
    border-radius: 1px;
    background-color: white;
    padding: 20px 25px;
    display: flex;
    flex-direction: column;

    @include bp-down($xs) {
        padding: 15px 10px 2px;
    }
}

.project-data {
    margin-right: 20px;
    flex: 2;
    padding: 20px 25px;

    @include bp-down($sm) {
        flex-basis: 100%;
        margin: 0 0 25px;
    }
}
.data-readings .project-readings {
    flex: 1;
    min-width: 360px;

    @include bp-down($xs) {
        min-width: unset;
    }
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
    font-family: $font-family-light;
    overflow-wrap: anywhere;

    &:not(:last-of-type) {
        padding-bottom: 6px;
    }
}

::v-deep .default-user-icon {
    margin: 6px 7px 0 0;
    width: 35px;
    height: 35px;
}

.location-name,
.location-native {
    white-space: break-spaces;
    display: flex;
    align-items: baseline;
    overflow-wrap: anywhere;
}
</style>
