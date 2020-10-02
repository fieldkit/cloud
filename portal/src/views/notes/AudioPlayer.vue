<template>
    <div class="">
        <VueAudio :file="loaded" v-if="loaded" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import VueAudio from "vue-audio";

export default Vue.extend({
    name: "AudioPlayer",
    components: {
        VueAudio,
    },
    props: {
        url: {
            type: String,
            required: true,
        },
    },
    data: () => {
        return {
            loaded: null,
        };
    },
    mounted(this: any) {
        return this.$services.api.loadMedia(this.url).then((loaded) => {
            this.loaded = loaded;
        });
    },
});
</script>

<style scoped></style>
