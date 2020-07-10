<template>
    <StandardLayout>
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div class="main-panel" v-show="!loading">
            <router-link :to="{ name: 'viewProject', params: { id: this.id } }">
                <div class="projects-link">
                    <span class="small-arrow">&lt;</span>
                    Back to Project
                </div>
            </router-link>
            <div id="inner-container" v-if="!loading">
                <ProjectForm :project="activeProject" @updating="onProjectUpdate" />
            </div>
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
    props: {
        id: {
            type: Number,
        },
    },
    data: () => {
        return {
            user: {},
            activeProject: null,
            loading: true,
        };
    },
    mounted() {
        if (this.id) {
            return this.getProject(this.id);
        } else {
            this.loading = false;
        }
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        getProject(projectId) {
            this.loading = true;
            return new FKApi()
                .getProject(projectId)
                .then((project) => {
                    this.activeProject = project;
                    this.loading = false;
                })
                .catch(() => {
                    return this.$router.push({ name: "projects" });
                });
        },
        onProjectUpdate() {
            this.loading = true;
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
