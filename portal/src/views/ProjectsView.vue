<template>
    <div>
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <SidebarNav viewing="projects" />
        <div class="main-panel">
            <div id="projects-container" v-if="!addingProject">
                <div class="container">
                    <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                        <img alt="Add project" src="../assets/add.png" />
                        Add Project
                    </div>
                    <h1>{{ projectsTitle }}</h1>
                    <img alt="Default Fieldkit Project" src="../assets/default.png" />
                </div>
                <div class="container">
                    <h1>Community</h1>
                </div>
            </div>
            <div v-if="addingProject">
                <ProjectForm @closeProjectForm="closeAddProject" />
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import ProjectForm from "../components/ProjectForm";
import SidebarNav from "../components/SidebarNav";

export default {
    name: "ProjectsView",
    components: {
        HeaderBar,
        ProjectForm,
        SidebarNav
    },
    data: () => {
        return {
            user: {},
            projectsTitle: "Projects",
            isAuthenticated: false,
            addingProject: false
        };
    },
    async beforeCreate() {
        const api = new FKApi();
        api.getCurrentUser().then(user => {
            this.user = user;
            this.isAuthenticated = true;
            this.projectsTitle = "My Projects";
            console.log("this is the user info", this.user);
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
        }
    }
};
</script>

<style scoped>
.main-panel {
    margin-left: 90px;
}
#add-project {
    float: right;
    padding: 12px;
    cursor: pointer;
}
#add-project img {
    vertical-align: bottom;
}
</style>
