<template>
    <StandardLayout>
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <router-link :to="{ name: 'projects' }">
                <div class="projects-link">
                    <span class="small-arrow">&lt;</span>
                    Back to Projects
                </div>
            </router-link>
            <div id="inner-container">
                <!-- add or update a project -->
                <ProjectForm :project="activeProject" @updating="onProjectUpdate" />
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
    </StandardLayout>
</template>

<script>
import StandardLayout from "./StandardLayout";
import FKApi from "../api/api";
import ProjectForm from "../components/ProjectForm";

export default {
    name: "ProjectEditView",
    components: {
        StandardLayout,
        ProjectForm,
    },
    props: ["id"],
    data: () => {
        return {
            user: {},
            userProjects: [],
            activeProject: null,
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
            .then((user) => {
                this.user = user;
                this.isAuthenticated = true;
                this.getUserProjects();
                this.getStations();
                if (this.id) {
                    this.getProject(this.id);
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
            this.api.getStations().then((s) => {
                this.stations = s.stations;
            });
        },
        getUserProjects() {
            this.api.getUserProjects().then((projects) => {
                if (projects && projects.projects.length > 0) {
                    this.userProjects = projects.projects;
                }
            });
        },
        getProject(projectId) {
            this.api
                .getProject(projectId)
                .then((project) => {
                    this.activeProject = project;
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
    margin: 40px 0 0 90px;
    font-size: 14px;
    cursor: pointer;
}
#inner-container {
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
