<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="projects" :projects="projects" :stations="stations" @showStation="showStation" />
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <div id="projects-container" v-if="!addingProject && !activeProject">
                <div class="container">
                    <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                        <img alt="Add project" src="../assets/add.png" />
                        Add Project
                    </div>
                    <h1>{{ projectsTitle }}</h1>
                    <div v-for="project in projects" v-bind:key="project.id" class="project-container">
                        <router-link :to="{ name: 'projectById', params: { id: project.id } }">
                            <img alt="Default Fieldkit Project" src="../assets/fieldkit_project.png" />
                            <div class="project-name">{{ project.name }}</div>
                            <div class="project-description">{{ project.description }}</div>
                        </router-link>
                    </div>
                </div>
                <div class="container">
                    <h1>Community</h1>
                </div>
            </div>
            <div v-if="addingProject">
                <ProjectForm @closeProjectForm="closeAddProject" />
            </div>
            <ProjectSummary :project="activeProject" :stations="stations" :user="user" ref="projectSummary" />
        </div>
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
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
            if (to.params.id) {
                if (this.isAuthenticated) {
                    this.api.getProject(to.params.id).then(project => {
                        if (project) {
                            this.activeProject = project;
                            this.$refs.projectSummary.viewSummary();
                        }
                    });
                }
            } else {
                this.activeProject = null;
                this.$refs.projectSummary.closeSummary();
            }
        }
    },
    data: () => {
        return {
            user: {},
            projects: [],
            projectsTitle: "Projects",
            activeProject: null,
            stations: [],
            isAuthenticated: false,
            addingProject: false,
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
                    this.api.getProject(this.id).then(project => {
                        if (project) {
                            this.activeProject = project;
                            this.$refs.projectSummary.viewSummary();
                            this.loading = false;
                        }
                    });
                } else {
                    this.loading = false;
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
        addProject() {
            this.addingProject = true;
        },
        closeAddProject() {
            this.addingProject = false;
        },
        showStation(station) {
            this.$router.push({ name: "stations", params: { id: station.id } });
        }
    }
};
</script>

<style scoped>
#loading {
    float: left;
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.main-panel {
    margin-left: 280px;
}
.no-auth-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 300px;
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
</style>
