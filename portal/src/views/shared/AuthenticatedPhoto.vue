<template>
    <img v-if="photo && !processing" :src="photo" class="authenticated-photo photo" :class="{ processing: processing }" alt="Image" />
    <Spinner v-else class="spinner" />
</template>

<script lang="ts">
import Vue from "vue";
import Spinner from "@/views/shared/Spinner.vue";

export default Vue.extend({
    name: "AuthenticatedPhoto",
    components: {
        Spinner,
    },
    props: {
        url: {
            type: String,
            required: true,
        },
        processing: {
            type: Boolean,
            default: false,
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
            this.loading = true;
            return this.$services.api.loadMedia(this.url).then((photo) => {
                this.photo = photo;
            });
        },
    },
});
</script>

<style scoped lang="scss">
.photo {
    transition: all 0.25s;
}

.photo.processing {
    opacity: 0.5;
    filter: blur(3px);
}

.spinner {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}
</style>
