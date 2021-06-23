<template>
    <Spinner v-if="loading" class="spinner" />
    <img v-else-if="station.photos && photo" :src="photo" class="station-photo photo" alt="Station Image" />
    <img v-else src="@/assets/station-image-placeholder.png" class="station-photo photo" alt="Default Station Image" />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { DisplayStation } from "@/store";
import Spinner from "./Spinner.vue";

export default Vue.extend({
    name: "StationPhoto",
    components: {
        Spinner,
    },
    props: {
        station: {
            type: Object as PropType<DisplayStation>,
            required: true,
        },
    },
    data(): {
        photo: unknown | null;
        loading: boolean;
    } {
        return {
            photo: null,
            loading: true,
        };
    },
    watch: {
        async station(): Promise<void> {
            await this.refresh();
        },
    },
    async mounted(): Promise<void> {
        await this.refresh();
    },
    methods: {
        async refresh(): Promise<void> {
            console.log(`loading-photo:`, this.station);
            if (this.station.photos) {
                this.loading = true;
                try {
                    const photo = await this.$services.api.loadMedia(this.station.photos.small);
                    this.photo = photo;
                } finally {
                    this.loading = false;
                }
            } else {
                this.loading = false;
            }
        },
    },
});
</script>

<style scoped>
.spinner {
    margin-left: auto;
    margin-right: auto;
    margin-top: auto;
    margin-bottom: auto;
}
</style>
