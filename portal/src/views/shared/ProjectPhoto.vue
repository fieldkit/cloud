<template>
    <div v-if="loading" class="station-photo loading-container">
        <Spinner class="spinner" />
    </div>
    <img v-else-if="photo" :src="photo" class="project-photo project-image photo" alt="Project Image" />
    <img v-else src="@/assets/fieldkit_project.png" class="project-photo project-image photo" alt="FieldKit Project" />
</template>

<script lang="ts">
import Vue from "vue";
import Spinner from "@/views/shared/Spinner.vue";

export default Vue.extend({
    name: "ProjectPhoto",
    components: {
        Spinner,
    },
    props: {
        project: {
            type: Object,
            required: true,
        },
        imageSize: {
            type: Number,
            default: 800,
        },
    },
    data(): { photo: unknown; loading: boolean } {
        return {
            photo: null,
            loading: false,
        };
    },
    watch: {
        async project(): Promise<void> {
            void this.refresh();
        },
    },
    mounted(): Promise<void> {
        return this.refresh();
    },
    methods: {
        async refresh(): Promise<void> {
            if (this.project.photo) {
                this.loading = true;
                this.photo = await this.$services.api
                    .loadMedia(this.project.photo, { size: this.imageSize })
                    .catch(() => {
                        this.photo = null;
                    })
                    .finally(() => {
                        this.loading = false;
                    });
            } else {
                this.photo = null;
            }
        },
    },
});
</script>

<style scoped></style>
