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
import { getPartnerCustomization, getPartnerCustomizationWithDefault } from "@/views/shared/partners";
import { mapState } from "vuex";
import { isCustomisationEnabled } from "@/views/shared/partners";

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
    data: (): {
        showCopiedLink: boolean;
    } => {
        return {
            showCopiedLink: false,
        };
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
        getCurrentURL(): string {
            return window.location.href;
        },
        copyUrlToClipboard(): void {
            const inputEl = this.$refs["url"] as HTMLInputElement;
            navigator.clipboard.writeText(inputEl.value).then(() => {
                this.showCopiedLink = true;
                setTimeout(() => {
                    this.showCopiedLink = false;
                }, 3000);
            });
        },
        openMailClient(): void {
            const partnerCustomization = getPartnerCustomizationWithDefault();
            const body = this.getCurrentURL();
            const subject = this.$t(partnerCustomization.email.subject);
            window.location.href = "mailto:?subject=" + subject + "&body=" + body;
        },
    },
});
</script>

<style lang="scss">
@import "src/scss/mixins";

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

        &-mail {
            font-size: 14px;
        }
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

        &--url {
            cursor: initial;

            &:hover {
                background-color: initial;
            }
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

    .url {
        margin: 0 20px;
        display: flex;
        align-items: center;
        position: relative;

        input {
            text-overflow: ellipsis;
            overflow: hidden;
            white-space: nowrap;
            font-size: 14px;
            margin-left: 30px;
            border: 0;
            width: 100%;
            color: var(--color-dark);
            font-family: var(--font-family-medium);
            background-color: transparent;

            &:focus {
                outline: none;
            }
        }

        button {
            margin-left: 10px;
            font-size: 12px;
            padding: 5px 10px;
            background-color: #ffffff;
            border: 1px solid #d7dce1;
            border-radius: 4px;
            font-family: var(--font-family-bold);
        }
    }

    .url-copied {
        background-color: rgba(160, 219, 225, 0.1);
        opacity: 0;
        visibility: hidden;
        font-size: 14px;
        padding: 5px 10px;
        border-radius: 4px;
        transition: opacity 0.25s;
        box-shadow: 0 0px 4px 0 rgba(0, 0, 0, 0.12);
        @include position(absolute, null 0 -30px null);

        &.visible {
            opacity: 1;
            visibility: visible;
        }
    }
}
</style>
