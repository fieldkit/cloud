<template>
    <div :class="'share-panel ' + containerClass">
        <div class="heading">
            <div class="title">Share</div>
            <div class="close-button icon icon-close" v-on:click="onClose"></div>
        </div>
        <div class="share-options">
            <a class="share-button" :href="twitterUrl" target="blank">
                <i class="icon icon-twitter" aria-label="Share on Twitter"></i>
                Twitter
            </a>
            <a class="share-button" :href="facebookUrl" target="blank">
                <i class="icon icon-facebook" aria-label="Share on Facebook"></i>
                <span>
                    Facebook
                </span>
            </a>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { getPartnerCustomization } from "@/views/shared/partners";
import { mapState } from "vuex";

function getRelativeUrl(href: string): string {
    const link = document.createElement("a");
    link.href = href;
    return link.href;
}

export default Vue.extend({
    name: "SharePanel",
    components: {
        ...CommonComponents,
    },
    props: {
        token: {
            type: String,
            required: true,
        },
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
            // TODO We should probably have a "non-partner" customization.
            const partnerCustomization = getPartnerCustomization();
            if (partnerCustomization != null) {
                qs.append("text", partnerCustomization.sharing.viz);
            } else {
                qs.append("text", "Check out this data on FieldKit!");
            }
            return `https://twitter.com/intent/tweet?${qs.toString()}`;
        },
        facebookUrl(): string {
            const qs = new URLSearchParams();
            qs.append("u", this.vizUrl);
            return `https://www.facebook.com/sharer/sharer.php?${qs.toString()}`;
        },
        vizUrl(): string {
            const qs = new URLSearchParams();
            qs.append("v", this.token);
            return getRelativeUrl(`/viz?${qs.toString()}`);
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
    padding: 25px 20px;
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
    padding: 20px 0;

    .icon {
        width: 20px;
        font-size: 19px;
        margin-right: 10px;
        margin-top: -2px;
    }

    .share-button {
        padding: 11px 20px 13px 20px;
        cursor: pointer;
        display: flex;
        align-items: center;
        transition: background-color 0.25s;

        &:hover {
            background-color: #f4f5f7;
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
