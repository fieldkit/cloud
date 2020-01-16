<template>
    <div>
        <SidebarNav viewing="projects" :projects="projects" :stations="stations" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <div id="inner-container">
                <!-- display all projects -->
                <div id="projects-container" v-if="viewingAll">
                    <div class="container">
                        <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                            <img alt="Add project" src="../assets/add.png" />
                            Add Project
                        </div>
                        <h1>{{ projectsTitle }}</h1>
                        <div v-for="project in projects" v-bind:key="project.id" class="project-container">
                            <router-link :to="{ name: 'viewProject', params: { id: project.id } }">
                                <div v-if="project.media_url" class="custom-project-image-container">
                                    <img
                                        alt="Fieldkit Project"
                                        :src="getImageUrl(project)"
                                        class="custom-project-image"
                                    />
                                </div>
                                <div v-else>
                                    <img
                                        alt="Default Fieldkit Project"
                                        src="../assets/fieldkit_project.png"
                                    />
                                </div>
                                <div class="project-name">{{ project.name }}</div>
                                <div class="project-description">{{ project.description }}</div>
                            </router-link>
                        </div>
                    </div>
                    <div class="container">
                        <h1>Community</h1>
                    </div>
                </div>
                <!-- add or update a project -->
                <div v-show="addingOrUpdating">
                    <ProjectForm
                        :project="activeProject"
                        @closeProjectForm="closeProjectForm"
                        @updating="onProjectUpdate"
                    />
                </div>
                <!-- display one project -->
                <ProjectSummary
                    :project="activeProject"
                    :stations="stations"
                    :users="users"
                    ref="projectSummary"
                    @inviteUser="sendInvite"
                    @removeUser="removeUser"
                />
            </div>
        </div>
        <div v-if="failedAuth" class="no-auth-message">
            <p>
                Please
                <router-link :to="{ name: 'login' }" class="show-link">
                    log in
                </router-link>
                to view projects.
            </p>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import { API_HOST } from "../secrets";
import HeaderBar from "../components/HeaderBar";
import ProjectForm from "../components/ProjectForm";
import ProjectSummary from "../components/ProjectSummary";
import SidebarNav from "../components/SidebarNav";

export default {
    name: "ProjectsView",
    components: {
        HeaderBar,
        ProjectForm,
        ProjectSummary,
        SidebarNav
    },
    props: ["id"],
    watch: {
        // watching $route picks up changes that beforeRouteUpdate does not
        $route(to) {
            this.routeTo = to;
            if (to.params.id && this.isAuthenticated) {
                this.getProject(to.params.id);
                // refresh projects list
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
            } else {
                this.viewAllProjects();
            }
        }
    },
    data: () => {
        return {
            baseUrl: API_HOST,
            user: {},
            projects: [],
            projectsTitle: "Projects",
            activeProject: null,
            stations: [],
            users: [],
            isAuthenticated: false,
            viewingAll: false,
            addingOrUpdating: false,
            failedAuth: false,
            loading: true
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
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    } else {
                        // create default project for user
                        this.api.addDefaultProject().then(project => {
                            this.projects = [project];
                            this.loading = false;
                        });
                    }
                });
                if (this.id) {
                    this.routeTo = this.$route;
                    this.getProject(this.id);
                } else {
                    this.viewAllProjects();
                }
                this.api.getStations().then(s => {
                    this.stations = s.stations;
                });
            })
            .catch(() => {
                this.loading = false;
                this.failedAuth = true;
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        getProject(projectId) {
            this.api.getProject(projectId).then(this.handleProject);
            this.api.getUsersByProject(projectId).then(users => {
                this.users = users && users.users ? users.users : [];
            });
        },
        handleProject(project) {
            if (!this.routeTo || this.routeTo.name == "viewProject") {
                this.viewProject(project);
            }
            if (this.routeTo && this.routeTo.name == "editProject") {
                this.editProject(project);
            }
            // this.$forceUpdate();
        },
        addProject() {
            this.resetFlags();
            this.activeProject = null;
            this.addingOrUpdating = true;
        },
        editProject(project) {
            this.resetFlags();
            this.addingOrUpdating = true;
            this.activeProject = project;
            this.$refs.projectSummary.closeSummary();
            this.loading = false;
        },
        onProjectUpdate() {
            this.loading = true;
        },
        viewProject(project) {
            this.resetFlags();
            this.activeProject = project;
            this.$refs.projectSummary.viewSummary();
            this.loading = false;
        },
        viewAllProjects() {
            this.resetFlags();
            this.viewingAll = true;
            this.$refs.projectSummary.closeSummary();
            this.loading = false;
        },
        getImageUrl(project) {
            return this.baseUrl + "/projects/" + project.id + "/media/?t=" + Date.now();
        },
        closeProjectForm() {
            this.activeProject = null;
            this.resetFlags();
            this.viewingAll = true;
        },
        resetFlags() {
            this.viewingAll = false;
            this.addingOrUpdating = false;
        },
        sendInvite(params) {
            this.api.sendInvite(params).then(result => {
                console.log("invite sent?", result);
            });
        },
        removeUser(params) {
            this.api.removeUserFromProject(params).then(result => {
                console.log("user removed", result);
            });
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        }
    }
};
</script>

<style scoped>
#inner-container {
    width: 780px;
    margin: 40px 60px;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-auth-message {
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
    float: right;
    padding: 12px;
    cursor: pointer;
}
#add-project img {
    vertical-align: bottom;
}
.project-container {
    float: left;
    width: 276px;
    height: 225px;
    margin: 0 20px 20px 0;
    border: 1px solid rgb(235, 235, 235);
}
.project-name {
    font-weight: bold;
    font-size: 16px;
    margin: 10px 20px 0 20px;
}
.project-description {
    font-weight: lighter;
    font-size: 14px;
    margin: 0 20px 10px 20px;
}
.custom-project-image-container {
    text-align: center;
}
.custom-project-image {
    max-width: 275px;
    max-height: 135px;
}
</style>
