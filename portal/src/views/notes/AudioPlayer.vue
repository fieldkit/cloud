<template>
    <div>
        <audio controls v-if="loaded">
            <source :src="loaded" />
        </audio>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "AudioPlayer",
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
            this.loaded = loaded.replace("m4a", "mp4");
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
