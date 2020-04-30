<template>
    <div>
        <SidebarNav
            :isAuthenticated="isAuthenticated"
            viewing="projects"
            :projects="projects"
            :stations="stations"
            @showStation="showStation"
        />
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
                                <div class="project-image-container">
                                    <img
                                        alt="Fieldkit Project"
                                        v-if="project.media_url"
                                        :src="getImageUrl(project)"
                                        class="project-image"
                                    />
                                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                                </div>
                                <img v-if="project.private" alt="Private project" src="../assets/private.png" class="private-icon" />

                                <div class="project-name">{{ project.name }}</div>
                                <div class="project-description">{{ project.description }}</div>
                                <div class="stats-icon-container">
                                    <div class="stat follows">
                                        <img alt="Follows" src="../assets/heart.png" class="follow-icon" />
                                        <span>12</span>
                                    </div>
                                    <div class="stat notifications">
                                        <img alt="Notifications" src="../assets/notification.png" class="notify-icon" />
                                        <span>2</span>
                                    </div>
                                    <div class="stat comments">
                                        <img alt="Comments" src="../assets/comment.png" class="comment-icon" />
                                        <span>3</span>
                                    </div>
                                </div>
                            </router-link>
                        </div>
                    </div>
                    <div class="container">
                        <h1>Community Projects</h1>
                    </div>
                </div>
                <!-- add or update a project -->
                <div v-show="addingOrUpdating">
                    <ProjectForm :project="activeProject" @closeProjectForm="closeProjectForm" @updating="onProjectUpdate" />
                </div>
                <!-- display one project -->
                <ProjectSummary
                    :project="activeProject"
                    :userStations="stations"
                    :users="users"
                    ref="projectSummary"
                    @inviteUser="sendInvite"
                    @removeUser="removeUser"
                    @showStation="showStation"
                />
            </div>
        </div>
        <div v-if="failedAuth" class="no-auth-message">
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
        SidebarNav,
    },
    props: ["id"],
    watch: {
        // watching $route picks up changes that beforeRouteUpdate does not
        $route(to) {
            this.routeTo = to;
            if (this.isAuthenticated) {
                // refresh projects list
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
            }
            if (to.params.id && this.isAuthenticated) {
                this.getProject(to.params.id);
            } else {
                this.viewAllProjects();
            }
        },
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
            this.api
                .getProject(projectId)
                .then(this.handleProject)
                .catch(() => {
                    this.$router.push({ name: "projects" });
                });
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
            this.$router.push({ name: "projects" });
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
        },
    },
};
</script>

<style scoped>
#inner-container {
    width: 1080px;
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
    margin: 40px 0 0 0;
    float: right;
    padding: 12px;
    cursor: pointer;
}
#add-project img {
    vertical-align: bottom;
}
.project-container {
    position: relative;
    float: left;
    width: 270px;
    height: 265px;
    margin: 0 25px 25px 0;
    border: 1px solid rgb(235, 235, 235);
}
.project-name {
    font-weight: bold;
    font-size: 16px;
    margin: 10px 15px 0 15px;
}
.project-description {
    font-weight: lighter;
    font-size: 14px;
    margin: 0 15px 10px 15px;
}
.project-image-container {
    height: 138px;
    text-align: center;
    border-bottom: 1px solid rgb(235, 235, 235);
}
.project-image {
    max-width: 270px;
    max-height: 138px;
}
.private-icon {
    float: right;
    margin: -14px 14px 0 0;
    position: relative;
    z-index: 10;
}
.stats-icon-container {
    position: absolute;
    bottom: 20px;
}
.stat {
    display: inline-block;
    font-size: 14px;
    font-weight: 600;
    margin: 0 14px 0 15px;
}
.stat img {
    float: left;
    margin: 2px 4px 0 0;
}
</style>
