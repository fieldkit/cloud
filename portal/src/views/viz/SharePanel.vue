<template>
    <div :class="'share-panel ' + containerClass">
        <div class="heading">
            <div class="title">Share</div>
            <div class="close-button icon icon-close" v-on:click="onClose"></div>
        </div>
        <div class="share-options">
            <a class="twitter-share-button" :href="twitterUrl" target="blank">
                <img alt="Share on Twitter" src="../../assets/icon-twitter.svg" />
                Twitter
            </a>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { mapState } from "vuex";
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
            return `https://twitter.com/intent/tweet?${qs.toString()}`;
        },
        vizUrl(): string {
            const qs = new URLSearchParams();
            qs.append("bookmark", serializeBookmark(this.bookmark));
            return `https://portal.fkdev.org/dashboard/explore?${qs.toString()}`;
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
    align-items: center;
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

    .twitter-share-button {
        padding: 10px 10px 10px 0px;
        cursor: pointer;
        display: flex;
        align-items: center;

        img {
            padding: 5px 10px 5px 5px;
        }
    }

    .link {
        font-size: 12px;
        padding: 10px;
        border: 1px solid rgb(215, 220, 225);
        border-radius: 4px;
        cursor: pointer;
        margin-bottom: 1em;
    }
}
</style>
