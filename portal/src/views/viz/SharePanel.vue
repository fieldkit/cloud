<template>
    <div :class="'share-panel ' + containerClass">
        <div class="heading">
            <div class="title">Share</div>
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="../../assets/close.png" />
            </div>
        </div>
        <div class="share-options">
            <div class="button" @click="onShareTwitter">Twitter</div>
            <a class="twitter-share-button" :href="'https://twitter.com/intent/tweet?url=' + vizUrl">
                Tweet
            </a>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { mapState, mapGetters } from "vuex";

export default Vue.extend({
    name: "SharePanel",
    components: {
        ...CommonComponents,
    },
    props: {
        bookmark: {
            type: Object,
            required: true,
        },
        containerClass: {
            type: String,
            required: true,
        },
    },
    data: (): {} => {
        return {};
    },
    computed: {
        ...mapState({}),

        vizUrl(): string {
            return "https://portal.fieldkit.org/dashboard/explore/" + JSON.stringify(this.bookmark);
        },
    },
    methods: {
        onClose() {
            this.$emit("close");
        },
        onShareTwitter() {
            console.log("ok");
        },
    },
});
</script>

<style lang="scss">
.share-panel .heading {
    padding: 1em;
    display: flex;
}
.share-panel .heading .title {
    font-size: 20px;
    font-weight: 500;
}
.share-panel .heading .close-button {
    margin-left: auto;
    cursor: pointer;
}

.share-options {
    padding: 20px;

    .button {
        font-size: 12px;
        padding: 10px;
        background-color: #ffffff;
        border: 1px solid rgb(215, 220, 225);
        border-radius: 4px;
        cursor: pointer;
        margin-bottom: 1em;
    }
}
</style>
