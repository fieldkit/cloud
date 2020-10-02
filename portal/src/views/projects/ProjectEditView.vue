<template>
    <StandardLayout :viewingProjects="true" :viewingProject="activeProject">
        <div class="main-panel" v-if="!loading">
            <ProjectForm :project="activeProject" @updating="onProjectUpdate" />
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";
import StandardLayout from "../StandardLayout.vue";
import ProjectForm from "./ProjectForm.vue";

import FKApi from "@/api/api";

export default Vue.extend({
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
    mounted(this: any) {
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
            return this.$services.api
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
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/layout";

.small-arrow {
    font-size: 11px;
    float: left;
    margin: 2px 5px 0 0;
}
.projects-link {
    font-size: 14px;
    cursor: pointer;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.show-link {
    text-decoration: underline;
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
