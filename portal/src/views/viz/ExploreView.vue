<template>
    <ExploreWorkspace
        v-if="visibleBookmark"
        :token="token"
        :bookmark="visibleBookmark"
        :exportsVisible="exportsVisible"
        :shareVisible="shareVisible"
        @open-bookmark="openBookmark"
        @export="exportWorkspace"
        @share="shareWorkspace"
    />
</template>

<script lang="ts">
import { Bookmark, serializeBookmark, deserializeBookmark } from "./viz";

import Vue from "vue";
import ExploreWorkspace from "./ExploreWorkspace.vue";

export default Vue.extend({
    name: "ExploreView",
    components: {
        ExploreWorkspace,
    },
    props: {
        token: {
            type: String,
            required: false,
        },
        bookmark: {
            type: Bookmark,
            required: false,
        },
        exportsVisible: {
            type: Boolean,
            default: false,
        },
        shareVisible: {
            type: Boolean,
            default: false,
        },
    },
    data(): {
        resolved: { [index: string]: Bookmark };
        bookmarkToToken: { [bookmark: string]: string };
    } {
        return {
            resolved: {},
            bookmarkToToken: {},
        };
    },
    computed: {
        visibleBookmark(): null | Bookmark {
            if (this.bookmark) {
                return this.bookmark;
            }
            if (this.token) {
                if (this.resolved[this.token]) {
                    return this.resolved[this.token];
                }
            }
            return null;
        },
    },
    watch: {
        async token(newValue: Bookmark, oldValue: Bookmark): Promise<void> {
            console.log(`viz: bookmark-route(token):`, newValue);
            await this.refreshBookmarkFromToken();
        },
        async bookmark(newValue: Bookmark, oldValue: Bookmark): Promise<void> {
            console.log(`viz: bookmark-route(bookmark):`, newValue);
        },
    },
    async beforeMount(): Promise<void> {
        if (this.token) {
            await this.refreshBookmarkFromToken();
        }
    },
    methods: {
        async refreshBookmarkFromToken(): Promise<void> {
            const token = this.token;
            console.log(`viz: bookmark-resolving`, token);

            try {
                if (!this.resolved[token]) {
                    const savedBookmark = await this.$services.api.resolveBookmark(token);
                    console.log(`viz: bookmark-resolved`, savedBookmark);
                    Vue.set(this.resolved, token, deserializeBookmark(savedBookmark.bookmark));
                }
            } catch (error) {
                console.log("viz: bad-token", error);
            }
        },
        async openBookmark(bookmark: Bookmark): Promise<void> {
            const encoded = serializeBookmark(bookmark);
            if (!this.bookmarkToToken[encoded]) {
                console.log(`viz: open-bookmark-saving`, encoded);
                const savedBookmark = await this.$services.api.saveBookmark(encoded);
                Vue.set(this.bookmarkToToken, encoded, savedBookmark.token);
                Vue.set(this.resolved, savedBookmark.token, bookmark);
                console.log(`viz: open-bookmark-saved`, savedBookmark.token);
            }
            await this.$router.push({ name: "exploreShortBookmark", query: { v: this.bookmarkToToken[encoded] } });
        },
        async exportWorkspace(): Promise<void> {
            await this.$router.push({ name: "exportWorkspace", query: { v: this.token } });
        },
        async shareWorkspace(): Promise<void> {
            await this.$router.push({ name: "shareWorkspace", query: { v: this.token } });
        },
    },
});
</script>

<style lang="scss"></style>
