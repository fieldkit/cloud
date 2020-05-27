<template>
    <div>
        <SidebarNav
            :isAuthenticated="isAuthenticated"
            viewing="projects"
            :projects="userProjects"
            :stations="stations"
            @showStation="showStation"
        />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <router-link :to="{ name: 'viewProject', params: { id: project.id } }" v-if="project">
                <div class="projects-link">
                    <span class="small-arrow">&lt;</span>
                    Back to {{ project.name }}
                </div>
            </router-link>
            <div id="inner-container">
                <!-- add or update a project update -->
                <ProjectUpdateForm :projectUpdate="activeUpdate" :project="project" @updating="onProjectUpdate" />
            </div>
        </div>
        <div v-if="noCurrentUser" class="no-user-message">
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
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import ProjectUpdateForm from "../components/ProjectUpdateForm";
import SidebarNav from "../components/SidebarNav";

export default {
    name: "ProjectEditView",
    components: {
        HeaderBar,
        ProjectUpdateForm,
        SidebarNav,
    },
    props: ["id", "project"],
    data: () => {
        return {
            user: {},
            userProjects: [],
            activeUpdate: null,
            stations: [],
            isAuthenticated: false,
            noCurrentUser: false,
            loading: true,
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                this.getUserProjects();
                this.getStations();
                if (this.id) {
                    this.getProjectUpdate(this.id);
                } else {
                    // adding project
                    this.loading = false;
                }
            })
            .catch(() => {
                this.loading = false;
                this.noCurrentUser = true;
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        getStations() {
            this.api.getStations().then(s => {
                this.stations = s.stations;
            });
        },
        getUserProjects() {
            this.api.getUserProjects().then(projects => {
                if (projects && projects.projects.length > 0) {
                    this.userProjects = projects.projects;
                }
            });
        },
        getProjectUpdate(updateId) {
            this.api
                .getProjectUpdate(updateId)
                .then(update => {
                    this.activeUpdate = update;
                    this.loading = false;
                })
                .catch(() => {
                    this.$router.push({ name: "projects" });
                });
        },
        onProjectUpdate() {
            this.loading = true;
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
    margin: 40px 0 0 60px;
    font-size: 14px;
    cursor: pointer;
}
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
