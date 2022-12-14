<template>
    <div v-if="loading" class="station-photo loading-photo">
        <Spinner class="spinner" />
    </div>
    <img v-else-if="station.photos && photo" :src="photo" class="station-photo photo" alt="Station Image" />
    <img
        v-else
        :src="$loadAsset(interpolatePartner('station-image-placeholder-') + '.png')"
        class="station-photo photo"
        alt="Default Station Image"
    />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { DisplayStation } from "@/store";
import Spinner from "./Spinner.vue";
import { interpolatePartner } from "@/views/shared/partners";

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
            loading: false,
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
            // console.log(`loading-photo:`, this.station);
            if (this.station.photos) {
                this.loading = true;
                try {
                    const photo = await this.$services.api.loadMedia(this.station.photos.small);
                    this.photo = photo;
                } finally {
                    this.loading = false;
                }
            }
        },
        interpolatePartner(baseString: string): string {
            return interpolatePartner(baseString);
        },
    },
});
</script>

<style scoped lang="scss"></style>
