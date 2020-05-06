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
            <router-link :to="{ name: 'projects' }" v-if="!viewingAll && !previewingProfile">
                <div class="projects-link">
                    <span class="small-arrow">&lt;</span>
                    Back to Projects
                </div>
            </router-link>
            <div class="projects-link" v-if="!viewingAll && previewingProfile" v-on:click="switchToDashboard">
                <span class="small-arrow">&lt;</span>
                Back to Project Dashboard
            </div>
            <div id="inner-container">
                <!-- display user's projects -->
                <div id="projects-container" v-if="viewingAll">
                    <div class="container">
                        <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                            <img alt="Add project" src="../assets/add.png" />
                            Add Project
                        </div>
                        <h1>{{ projectsTitle }}</h1>
                        <ProjectThumbnails :projects="userProjects" />
                    </div>
                    <div class="container">
                        <h1>Community Projects</h1>
                        <ProjectThumbnails :projects="publicProjects" />
                    </div>
                </div>
                <!-- add or update a project -->
                <div v-show="addingOrUpdating">
                    <ProjectForm :project="activeProject" @closeProjectForm="closeProjectForm" @updating="onProjectUpdate" />
                </div>
                <!-- display project dashboard -->
                <ProjectDashboard
                    :user="user"
                    :project="activeProject"
                    :userStations="stations"
                    :users="users"
                    ref="projectDashboard"
                    @inviteUser="sendInvite"
                    @removeUser="removeUser"
                    @viewProfile="switchToProfile"
                    @showStation="showStation"
                />
                <!-- display project profile -->
                <ProjectProfile :user="user" :project="activeProject" :users="users" ref="projectProfile" @showStation="showStation" />
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
import ProjectDashboard from "../components/ProjectDashboard";
import ProjectProfile from "../components/ProjectProfile";
import ProjectThumbnails from "../components/ProjectThumbnails";
import SidebarNav from "../components/SidebarNav";

export default {
    name: "ProjectsView",
    components: {
        HeaderBar,
        ProjectForm,
        ProjectDashboard,
        ProjectProfile,
        ProjectThumbnails,
        SidebarNav,
    },
    props: ["id"],
    watch: {
        // watching $route picks up changes that beforeRouteUpdate does not
        $route(to) {
            this.routeTo = to;
            if (this.isAuthenticated) {
                // refresh projects list
                this.api.getUserProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.userProjects = projects.projects;
                    }
                });
            }
            if (to.params.id && this.isAuthenticated) {
                this.getProject(to.params.id);
            } else {
                this.viewAllProjects();
            }
            this.api.getPublicProjects().then(projects => {
                this.publicProjects = projects;
            });
        },
    },
    data: () => {
        return {
            baseUrl: API_HOST,
            user: {},
            userProjects: [],
            projectsTitle: "Projects",
            activeProject: null,
            publicProjects: [],
            stations: [],
            users: [],
            isAuthenticated: false,
            viewingAll: false,
            previewingProfile: false,
            addingOrUpdating: false,
            failedAuth: false,
            loading: true,
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api.getPublicProjects().then(projects => {
            this.publicProjects = projects;
        });
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                this.projectsTitle = "My Projects";
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
            if (this.routeTo && this.routeTo.name == "editProject") {
                this.editProject(project);
            } else {
                const userProject = this.userProjects.find(p => {
                    return p.id == project.id;
                });
                if (userProject) {
                    this.viewProjectDashboard(project);
                } else if (!project.private) {
                    this.viewProjectProfile(project);
                } else {
                    this.$router.push({ name: "projects" });
                }
            }
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
            this.$refs.projectProfile.closeSummary();
            this.$refs.projectDashboard.closeSummary();
            this.loading = false;
        },
        onProjectUpdate() {
            this.loading = true;
        },
        viewProjectDashboard(project) {
            this.resetFlags();
            this.activeProject = project;
            this.$refs.projectProfile.closeSummary();
            this.$refs.projectDashboard.viewSummary();
            this.loading = false;
        },
        switchToDashboard() {
            this.previewingProfile = false;
            this.$refs.projectProfile.closeSummary();
            this.$refs.projectDashboard.viewSummary();
        },
        viewProjectProfile(project) {
            this.resetFlags();
            this.activeProject = project;
            this.$refs.projectDashboard.closeSummary();
            this.$refs.projectProfile.viewSummary();
            this.loading = false;
        },
        switchToProfile() {
            this.previewingProfile = true;
            this.$refs.projectDashboard.closeSummary();
            this.$refs.projectProfile.viewSummary();
        },
        viewAllProjects() {
            this.resetFlags();
            this.viewingAll = true;
            this.$refs.projectDashboard.closeSummary();
            this.$refs.projectProfile.closeSummary();
            this.loading = false;
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
</style>
