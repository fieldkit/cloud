<template>
    <img v-if="photo" :src="photo" class="authenticated-photo photo" alt="Image" />
</template>

<script lang="ts">
import Vue from "vue";

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
            return this.$services.api.loadMedia(this.url).then((photo) => {
                this.photo = photo;
            });
        },
    },
});
</script>

<style scoped lang="scss">

</style>
