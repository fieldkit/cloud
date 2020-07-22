<template>
    <div>
        <img v-if="photo" :src="photo" class="project-photo project-image photo" alt="Project Image" />
        <img v-if="missing" src="@/assets/fieldkit_project.png" class="project-photo project-image photo" alt="FieldKit Project" />
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "ProjectPhoto",
    props: {
        project: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            photo: null,
        };
    },
    watch: {
        project(this: any) {
            return this.refresh();
        },
    },
    computed: {
        missing(this: any) {
            return !this.project.photo;
        },
    },
    mounted(this: any) {
        return this.refresh();
    },
    methods: {
        refresh(this: any) {
            if (this.project.photo) {
                return new FKApi().loadMedia("/projects/" + this.project.id + "/media").then((photo) => {
                    this.photo = photo;
                });
            }
        },
    },
});
</script>

<style scoped></style>
