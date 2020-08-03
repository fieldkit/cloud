<template>
    <div class="project-thumbnails">
        <ProjectThumbnail v-for="project in projects" v-bind:key="project.id" :project="project" :invited="invited" />
    </div>
</template>

<script lang="ts">
import ProjectThumbnail from "./ProjectThumbnail.vue";

export default {
    name: "ProjectThumbnails",
    components: {
        ProjectThumbnail,
    },
    props: {
        projects: {
            type: Array,
            required: true,
        },
        invited: {
            type: Boolean,
            default: false,
        },
    },
};
</script>

<style scoped>
.project-thumbnails {
    display: flex;
    flex-wrap: wrap;
    text-align: left;
}
</style>
