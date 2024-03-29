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
            <a class="share-button" @click="openMailClient()">
                <i class="icon icon-mail"></i>
                <span>
                    {{ $t("sharePanel.emailUrl") }}
                </span>
            </a>
            <a class="share-button share-button--url">
                <i class="icon icon-copy" :aria-label="$t('sharePanel.copyUrl')"></i>
                <span>
                    {{ $t("sharePanel.copyUrl") }}
                </span>
                <div class="url" target="blank">
                    <input readonly ref="url" :value="getCurrentURL()" />
                    <button @click="copyUrlToClipboard()">
                        {{ $t("sharePanel.copyBtn") }}
                    </button>
                    <div class="url-copied-wrap" :class="{ visible: showCopiedLink }">
                        <span class="url-copied">
                            {{ $t("sharePanel.linkCopied") }}
                        </span>
                    </div>
                </div>
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
            const path = this.$route.path.split("/")[2];
            const fullPath = window.location.href.replace("/" + path, "");
            return fullPath;
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
            flex-wrap: wrap;
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
        margin-left: 32px;
        margin-top: 10px;
        display: flex;
        align-items: center;
        position: relative;
        flex: 1 1 100%;

        input {
            text-overflow: ellipsis;
            overflow: hidden;
            white-space: nowrap;
            font-size: 14px;
            width: 100%;
            height: 32px;
            padding-left: 5px;
            color: var(--color-dark);
            border: 1px solid var(--color-border);
            font-family: var(--font-family-medium);

            &:focus {
                outline: none;
            }
        }

        button {
            margin-left: 10px;
            font-size: 12px;
            padding: 6px 10px;
            background-color: #ffffff;
            border: 1px solid #d7dce1;
            border-radius: 4px;
            font-family: var(--font-family-bold);
        }
    }

    .url-copied-wrap {
        background-color: #fff;
        box-shadow: 0 0px 4px 0 rgba(0, 0, 0, 0.12);
        opacity: 0;
        visibility: hidden;
        transition: opacity 0.25s;
        border-radius: 4px;
        display: flex;
        @include position(absolute, null 0 -35px null);

        &.visible {
            opacity: 1;
            visibility: visible;
        }
    }

    .url-copied {
        background-color: rgba(160, 219, 225, 0.1);
        border-radius: 4px;
        font-size: 14px;
        padding: 5px 10px;
    }
}
</style>
