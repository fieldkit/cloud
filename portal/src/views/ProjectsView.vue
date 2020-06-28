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
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <div id="inner-container">
                <div id="projects-container">
                    <div class="container">
                        <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                            <img alt="Add project" src="../assets/add.png" />
                            Add Project
                        </div>
                        <h1 v-if="isAuthenticated">My Projects</h1>
                        <h1 v-if="!isAuthenticated">Projects</h1>
                        <ProjectThumbnails :projects="userProjects" />
                    </div>
                    <div class="container">
                        <h1>Community Projects</h1>
                        <ProjectThumbnails :projects="publicProjects" />
                    </div>
                </div>
            </div>
        </div>
        <div v-if="user" class="no-user-message">
            <p>
                Please
                <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="show-link">
                    log in
                </router-link>
                to view projects.
            </p>
        </div>
    </div>
</template>

<script>
import HeaderBar from "../components/HeaderBar";
import ProjectThumbnails from "../components/ProjectThumbnails";
import SidebarNav from "../components/SidebarNav";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
    name: "ProjectsView",
    components: {
        HeaderBar,
        ProjectThumbnails,
        SidebarNav,
    },
    data() {
        return {};
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({
            user: s => s.user.user,
            stations: s => s.stations.stations.user,
            userProjects: s => s.stations.projects.user,
            publicProjects: s => s.stations.projects.community,
            loading: s => s.stations.loading.stations || s.stations.loading.projects,
        }),
    },
    beforeMount() {
        this.$store.dispatch(ActionTypes.NEED_PROJECTS);
        this.$store.dispatch(ActionTypes.NEED_STATIONS);
    },
    methods: {
        goBack() {
            if (window.history.length > 1) {
                this.$router.go(-1);
            } else {
                this.$router.push("/");
            }
        },
        addProject() {
            this.$router.push({ name: "addProject" });
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
    },
};
</script>

<style scoped>
#inner-container {
    margin: 20px 60px;
    overflow: scroll;
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
#add-project {
    margin: 40px 0 0 0;
    float: right;
    padding: 12px;
    cursor: pointer;
}
#add-project img {
    vertical-align: bottom;
}
</style>
