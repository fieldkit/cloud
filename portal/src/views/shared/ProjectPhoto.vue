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
        imageSize: {
            type: Number,
            default: 800
        }
    },
    data(): { photo: unknown } {
        return {
            photo: null,
        };
    },
    watch: {
        async project(): Promise<void> {
            void this.refresh();
        },
    },
    computed: {
        missing(): boolean {
            return !this.project.photo;
        },
    },
    mounted(): Promise<void> {
        return this.refresh();
    },
    methods: {
        async refresh(): Promise<void> {
            if (this.project.photo) {
                this.photo = await this.$services.api.loadMedia(this.project.photo, { size: this.imageSize }).catch((e) => {
                    this.photo = null;
                });
            } else {
                this.photo = null;
            }
        },
    },
});
</script>

<style scoped></style>
