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
            <div id="inner-container">
                <!-- display admin view -->
                <ProjectAdmin
                    v-if="adminView"
                    :user="user"
                    :project="activeProject"
                    :userStations="stations"
                    :users="users"
                    @viewProfile="switchToPublic"
                />
                <!-- display public view -->
                <ProjectPublic v-if="publicView" :user="user" :project="activeProject" :users="users" />
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
import { API_HOST } from "../secrets";
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
    props: ["id"],
    watch: {
        id() {
            if (this.id) {
                this.reset();
            }
        },
    },
    data: () => {
        return {
            baseUrl: API_HOST,
            user: {},
            adminView: false,
            publicView: false,
            previewing: false,
            userProjects: [],
            activeProject: null,
            stations: [],
            users: [],
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
                this.getStations();
                this.init();
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
        reset() {
            this.adminView = false;
            this.publicView = false;
            this.userProjects = [];
            this.activeProject = null;
            this.users = [];
            this.loading = true;
            this.init();
        },
        getStations() {
            this.api.getStations().then(s => {
                this.stations = s.stations;
            });
        },
        init() {
            this.api.getUserProjects().then(projects => {
                if (projects && projects.projects.length > 0) {
                    this.userProjects = projects.projects;
                }
                this.getProject();
            });
        },
        getProjectUsers() {
            this.api.getUsersByProject(this.id).then(result => {
                let users = result && result.users ? result.users : [];
                this.users = users.map(u => {
                    // pending users come back with id = 0, which causes problems
                    if (u.user.id == 0) {
                        u.user.id = "pending-" + Math.random();
                    }
                    if (u.user.media_url) {
                        u.userImage = this.baseUrl + "/user/" + u.user.id + "/media";
                    } else {
                        let imgPath = require.context("../assets/", false, /\.png$/);
                        let img = "new_user.png";
                        u.userImage = imgPath("./" + img);
                    }
                    return u;
                });
            });
        },
        getProject() {
            this.api
                .getProject(this.id)
                .then(this.handleProject)
                .catch(() => {
                    this.$router.push({ name: "projects" });
                });
        },
        handleProject(project) {
            this.publicView = false;
            this.adminView = false;
            // TODO: don't show projects marked as private if user isn't member
            // ~ if (project.private) {
            //     this.$router.push({ name: "projects" });
            // }
            if (project.read_only) {
                this.publicView = true;
            } else {
                this.adminView = true;
            }
            this.getProjectUsers();
            this.activeProject = project;
            this.loading = false;
        },
        switchToAdmin() {
            this.previewing = false;
            this.adminView = true;
            this.publicView = false;
        },
        switchToPublic() {
            this.previewing = true;
            this.publicView = true;
            this.adminView = false;
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
