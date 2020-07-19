<template>
    <StandardLayout :viewingProjects="true">
        <div id="loading" v-if="isBusy">
            <img alt="" src="@/assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!isBusy && isAuthenticated">
            <div class="inner-container">
                <router-link :to="{ name: 'projects' }">
                    <div class="projects-link">
                        <span class="small-arrow">&lt;</span>
                        Back to Projects
                    </div>
                </router-link>
                <div class="view-container">
                    <ProjectAdmin
                        v-if="isAdministrator && displayProject"
                        :user="user"
                        :displayProject="displayProject"
                        :userStations="stations"
                    />
                    <ProjectPublic
                        v-if="!isAdministrator && displayProject"
                        :user="user"
                        :displayProject="displayProject"
                        :userStations="stations"
                    />
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";
import StandardLayout from "../StandardLayout.vue";
import ProjectPublic from "./ProjectPublic.vue";
import ProjectAdmin from "./ProjectAdmin.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "ProjectView",
    components: {
        StandardLayout,
        ProjectPublic,
        ProjectAdmin,
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
    },
    data: () => {
        return {};
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
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
        return this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id });
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
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
    },
});
</script>

<style scoped>
.main-panel {
    height: 100%;
    text-align: left;
    background-color: #fcfcfc;
    padding: 40px;
}
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
</style>
