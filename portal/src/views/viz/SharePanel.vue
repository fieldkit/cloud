<template>
    <div :class="'share-panel ' + containerClass">
        <div class="heading">
            <div class="title">Share</div>
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="../../assets/close.png" />
            </div>
        </div>
        <div class="share-options">
            <a class="twitter-share-button" :href="twitterUrl" target="blank">
                <div class="button">Twitter</div>
            </a>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { mapState, mapGetters } from "vuex";
import { serializeBookmark } from "./viz";

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
        twitterUrl(): string {
            const qs = new URLSearchParams();
            qs.append("url", this.vizUrl);
            return "https://twitter.com/intent/tweet" + "?" + qs.toString();
        },
        vizUrl(): string {
            return "https://portal.fkdev.org/dashboard/explore?bookmark=" + serializeBookmark(this.bookmark);
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
