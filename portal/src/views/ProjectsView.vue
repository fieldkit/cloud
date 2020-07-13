<template>
    <StandardLayout :viewingProjects="true">
        <div v-show="!isBusy && isAuthenticated">
            <div class="projects-container">
                <div class="container mine">
                    <div class="header">
                        <h1 v-if="isAuthenticated">My Projects</h1>
                        <h1 v-if="!isAuthenticated">Projects</h1>
                        <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                            <img alt="Add project" src="../assets/add.png" />
                            Add Project
                        </div>
                    </div>

                    <ProjectThumbnails :projects="userProjects" />
                </div>
                <div class="container community">
                    <div class="header">
                        <h1>Community Projects</h1>
                    </div>
                    <ProjectThumbnails :projects="publicProjects" />
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script>
import StandardLayout from "./StandardLayout";
import ProjectThumbnails from "../components/ProjectThumbnails";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
    name: "ProjectsView",
    components: {
        StandardLayout,
        ProjectThumbnails,
    },
    data() {
        return {};
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s) => s.user.user,
            stations: (s) => s.stations.user.stations,
            userProjects: (s) => s.stations.user.projects,
            publicProjects: (s) => s.stations.community.projects,
        }),
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
    },
};
</script>

<style scoped>
.projects-container {
    margin-left: 20px;
    display: flex;
    flex-direction: column;
    max-width: 900px;
}

.container {
}
.container.community {
    border-top: 2px solid #afafaf;
    margin-right: 20px;
}
.container .header {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: baseline;
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
#add-project {
    margin-left: auto;
    margin-right: 40px;
    cursor: pointer;
}
#add-project img {
    vertical-align: bottom;
}
</style>
