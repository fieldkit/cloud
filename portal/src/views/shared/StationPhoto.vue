<template>
    <img :src="photo" class="station-photo photo" v-if="photo" alt="Station Image" />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "StationPhoto",
    props: {
        station: {
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
        station(this: any) {
            return this.refresh();
        },
    },
    mounted(this: any) {
        return this.refresh();
    },
    methods: {
        refresh(this: any) {
            return new FKApi().loadMedia(this.station.photos.small).then((photo) => {
                this.photo = photo;
            });
        },
    },
});
</script>

<style scoped></style>
