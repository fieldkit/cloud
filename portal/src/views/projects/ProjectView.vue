<template>
    <StandardLayout :viewingProjects="true" :viewingProject="displayProject" :disableScrolling="activityVisible" :sidebarNarrow="bigMap">
        <div class="container-wrap" :class="{ 'big-map': bigMap }">
            <template v-if="displayProject">
                <DoubleHeader
                    :title="isAdministrator ? displayProject.name : null"
                    :subtitle="isAdministrator ? $t('project.dashboard') : null"
                    :backTitle="isAdministrator ? $t('layout.backProjects') : $t('layout.backProjectDashboard')"
                    backRoute="projects"
                    v-if="displayProject && !bigMap"
                >
                    <div class="activity-button" v-if="isAdministrator" v-on:click="onActivityToggle">
                        <i class="icon icon-notification"></i>
                        {{ $t("project.activity") }}
                    </div>
                </DoubleHeader>

                <div v-if="isAdministrator" v-bind:key="id">
                    <ProjectActivity
                        v-if="activityVisible"
                        :user="user"
                        :displayProject="displayProject"
                        containerClass="project-activity-floating"
                        @close="closeActivity"
                    />
                    <ProjectAdmin :user="user" :displayProject="displayProject" :userStations="stations" v-if="user" />
                </div>
                <ProjectPublic
                    v-if="!isAdministrator && displayProject && !bigMap"
                    :user="user"
                    :displayProject="displayProject"
                    :userStations="stations"
                />
                <ProjectBigMap
                    v-if="!isAdministrator && displayProject && bigMap"
                    :user="user"
                    :displayProject="displayProject"
                    :userStations="stations"
                />
            </template>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";
import ProjectPublic from "./ProjectPublic.vue";
import ProjectAdmin from "./ProjectAdmin.vue";
import ProjectActivity from "./ProjectActivity.vue";

import ProjectBigMap from "./ProjectBigMap.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "ProjectView",
    components: {
        ...CommonComponents,
        StandardLayout,
        ProjectPublic,
        ProjectAdmin,
        ProjectActivity,
        ProjectBigMap,
    },
    props: {
        id: {
            required: true,
            type: Number,
        },
        forcePublic: {
            required: true,
            type: Boolean,
        },
        activityVisible: {
            type: Boolean,
            default: false,
        },
        bigMap: {
            type: Boolean,
            default: false,
        },
    },
    data: () => {
        return {};
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => Object.values(s.stations.user.stations),
            userProjects: (s: GlobalState) => Object.values(s.stations.user.projects),
            displayProject() {
                return this.$getters.projectsById[this.id];
            },
            isAdministrator() {
                if (!this.forcePublic) {
                    const p = this.$getters.projectsById[this.id];
                    if (p) {
                        return !p.project.readOnly;
                    }
                }
                return false;
            },
        }),
    },
    watch: {
        id() {
            return this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id });
        },
    },
    beforeMount() {
        return this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id }).catch((e) => {
            if (e.name === "ForbiddenError") {
                return this.$router.push({ name: "login", params: { errorMessage: this.$t("login.privateProject") } });
            }
        });
    },
    methods: {
        goBack() {
            if (window.history.length > 1) {
                this.$router.go(-1);
            } else {
                this.$router.push("/");
            }
        },
        showStation(station) {
            return this.$router.push({ name: "mapStation", params: { id: station.id } });
        },
        onActivityToggle(this: any) {
            if (this.activityVisible) {
                return this.$router.push({ name: "viewProject", params: { id: this.id } });
            } else {
                return this.$router.push({ name: "viewProjectActivity", params: { id: this.id } });
            }
        },
        closeActivity(this: any) {
            return this.$router.push({ name: "viewProject", params: { id: this.id } });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/layout";

.small-arrow {
    font-size: 11px;
    float: left;
    margin: 2px 5px 0 0;
}
.projects-link {
    margin-top: 10px;
    margin-bottom: 10px;
    font-size: 14px;
    cursor: pointer;
}
.inner-container {
}
.view-container {
}
#projects-container {
    width: 890px;
    margin: 20px 0 0 0;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-user-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 40px;
}
.show-link {
    text-decoration: underline;
}
.container {
    float: left;
}
.activity-button {
    border-radius: 3px;
    border: solid 1px #cccdcf;
    background-color: #ffffff;
    cursor: pointer;
    font-family: var(--font-family-bold);
    padding: 10px 22px;
    @include flex(center, center);

    body.floodnet & {
        font-family: $font-family-floodnet-button;
    }

    .icon {
        margin-right: 14px;
        font-size: 17px;

        body.floodnet & {
            &:before {
                color: var(--color-dark);
            }
        }
    }
}
.project-activity-floating {
    position: absolute;
    right: 0;
    top: 70px;
    bottom: 0;
    background-color: #fcfcfc;
    border-left: 2px solid var(--color-border);
    z-index: 10;
    overflow-y: scroll;
    width: 30em;
    paddinig: 1em;

    @include bp-down($sm) {
        @include position(fixed, 55px null null 0);
        border: 0;
        width: 100%;
        padding: 20px 0;
    }
}
::v-deep .project-tag {
    margin-right: 10px;
}

::v-deep .pagination {
    > div {
        @include flex(center, ceenter);
    }

    .button {
        font-size: 13px;
        border: 0;
        margin: 0;
        transform: translateY(2px);
    }
}
</style>
