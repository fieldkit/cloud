<template>
    <img :src="photo" class="authenticated-photo photo" v-if="photo" alt="Image" />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "AuthenticatedPhoto",
    props: {
        url: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            photo: null,
        };
    },
    watch: {
        url(this: any) {
            return this.refresh();
        },
    },
    created(this: any) {
        return this.refresh();
    },
    methods: {
        refresh(this: any) {
            return new FKApi().loadMedia(this.url).then((photo) => {
                this.photo = photo;
            });
        },
    },
});
</script>

<style scoped></style>
