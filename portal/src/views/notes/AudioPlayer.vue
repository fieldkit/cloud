<template>
    <div>
        <VueAudio :file="loaded" type="audio/mp4" v-if="loaded" />

        <audio v-if="loaded" controls>
            <source :src="loaded" type="audio/x-caf" />
        </audio>
    </div>

    <!--  <div>
    <audio-player v-if="loaded && url"
        ref="audioPlayer"
        :audio-list="[loaded]"
        theme-color="#ff2929"
    />
  </div>-->
</template>

<script lang="ts">
import Vue from "vue";
import VueAudio from "vue-audio";
import AudioPlayer from "@liripeng/vue-audio-player";

export default Vue.extend({
    name: "AudioPlayer",
    components: {
        VueAudio,
        //   AudioPlayer,
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
            console.log("radoi loaded url", this.url);
            console.log("radoi loaded media", loaded);
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
