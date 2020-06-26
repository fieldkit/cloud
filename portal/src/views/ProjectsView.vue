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
            <div id="inner-container">
                <div id="projects-container">
                    <div class="container">
                        <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                            <img alt="Add project" src="../assets/add.png" />
                            Add Project
                        </div>
                        <h1>{{ projectsTitle }}</h1>
                        <!-- display user's projects -->
                        <ProjectThumbnails :projects="userProjects" />
                    </div>
                    <div class="container">
                        <h1>Community Projects</h1>
                        <!-- display community projects -->
                        <ProjectThumbnails :projects="publicProjects" />
                    </div>
                </div>
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
import ProjectThumbnails from "../components/ProjectThumbnails";
import SidebarNav from "../components/SidebarNav";

export default {
    name: "ProjectsView",
    components: {
        HeaderBar,
        ProjectThumbnails,
        SidebarNav,
    },
    props: ["id"],
    data: () => {
        return {
            user: {},
            userProjects: [],
            projectsTitle: "Projects",
            publicProjects: [],
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
                this.projectsTitle = "My Projects";
                this.getPublicProjects();
                this.getUserProjects();
                this.getStations();
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
        getPublicProjects() {
            this.api
                .getPublicProjects()
                .then(projects => {
                    this.publicProjects = projects;
                    this.loading = false;
                })
                .catch(() => {
                    this.loading = false;
                });
        },
        getUserProjects() {
            this.api.getUserProjects().then(projects => {
                if (projects && projects.projects.length > 0) {
                    this.userProjects = projects.projects;
                } else {
                    // create default project for user
                    this.api.addDefaultProject().then(project => {
                        this.userProjects = [project];
                        this.loading = false;
                    });
                }
            });
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
