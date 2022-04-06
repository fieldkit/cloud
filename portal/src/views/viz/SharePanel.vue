<template>
    <div :class="'share-panel ' + containerClass">
        <div class="heading">
            <div class="title">Share</div>
            <div class="close-button icon icon-close" v-on:click="onClose"></div>
        </div>
        <div class="share-options">
            <div>
                <a class="share-button" :href="twitterUrl" target="blank">
                    <div>
                        <img alt="Share on Twitter" src="../../assets/icon-twitter.svg" />
                    </div>
                    Twitter
                </a>
            </div>
            <div>
                <a class="share-button" :href="facebookUrl" target="blank">
                    <div>
                        <img alt="Share on Facebook" src="../../assets/icon-facebook.svg" />
                    </div>
                    Facebook
                </a>
            </div>
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

    a div {
        width: 40px;
    }

    .share-button {
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
