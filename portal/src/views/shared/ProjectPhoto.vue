<template>
    <img v-if="photo" :src="photo" class="project-photo project-image photo" alt="Project Image" />
    <img v-else-if="missing" src="@/assets/fieldkit_project.png" class="project-photo project-image photo" alt="FieldKit Project" />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

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
                return this.$services.api.loadMedia(this.project.photo, { size: 800 }).then((photo) => {
                    this.photo = photo;
                });
            } else {
                this.photo = null;
            }
        },
    },
});
</script>

<style scoped></style>
