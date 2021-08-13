<template>
    <VueAudio :file="loaded" v-if="loaded" />
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

<style scoped lang="scss">
@import "../../scss/mixins";

.vue-sound-wrapper {
    flex: 1;
}

::v-deep .vue-sound__player {
    @include flex(center);
}

::v-deep .vue-sound__playback-time-wrapper {
    display: flex !important;
    flex: 1;
    width: auto !important;
}
</style>
