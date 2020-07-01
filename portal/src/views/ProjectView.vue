<template>
    <div>
        <SidebarNav
            :isAuthenticated="isAuthenticated"
            :viewingProjects="true"
            :projects="userProjects"
            :stations="stations"
            @showStation="showStation"
        />
        <HeaderBar />
        <div id="loading" v-if="isBusy">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!isBusy && isAuthenticated">
            <div class="inner-container">
                <router-link :to="{ name: 'projects' }" v-if="!previewing">
                    <div class="projects-link">
                        <span class="small-arrow">&lt;</span>
                        Back to Projects
                    </div>
                </router-link>
                <div class="projects-link" v-if="previewing" v-on:click="switchToAdmin">
                    <span class="small-arrow">&lt;</span>
                    Back to Project Dashboard
                </div>
                <div class="view-container">
                    <ProjectAdmin
                        v-if="isAdministrator && displayProject"
                        :user="user"
                        :displayProject="displayProject"
                        :userStations="stations"
                        @viewProfile="switchToPublic"
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
    </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import HeaderBar from "../components/HeaderBar";
import ProjectPublic from "../components/ProjectPublic";
import SidebarNav from "../components/SidebarNav";
import ProjectAdmin from "../components/ProjectAdmin";

export default {
    name: "ProjectView",
    components: {
        HeaderBar,
        ProjectPublic,
        ProjectAdmin,
        SidebarNav,
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
        return {
            previewing: false,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s) => s.user.user,
            stations: (s) => s.stations.stations.user,
            userProjects: (s) => s.stations.projects.user,
            displayProject() {
                return this.$store.getters.projectsById[this.id];
            },
            isAdministrator() {
                if (!this.forcePublic) {
                    const p = this.$store.getters.projectsById[this.id];
                    if (p) {
                        return !p.readOnly;
                    }
                }
                return false;
            },
        }),
    },
    beforeMount() {
        this.$store.dispatch(ActionTypes.NEED_COMMON);
        this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.id });
    },
    methods: {
        goBack() {
            if (window.history.length > 1) {
                this.$router.go(-1);
            } else {
                this.$router.push("/");
            }
        },
        switchToAdmin() {
            this.previewing = false;
            // this.adminView = true;
            // this.publicView = false;
        },
        switchToPublic() {
            this.previewing = true;
            // this.publicView = true;
            // this.adminView = false;
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
    },
};
</script>

<style scoped>
.small-arrow {
    font-size: 11px;
    float: left;
    margin: 2px 5px 0 0;
}
.projects-link {
    margin: 40px 0 0 90px;
    font-size: 14px;
    cursor: pointer;
}
.inner-container {
}
.view-container {
    margin: 20px 60px;
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
