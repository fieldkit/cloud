<template>
    <ExploreWorkspace
        v-if="visibleBookmark"
        :bookmark="visibleBookmark"
        :exportsVisible="exportsVisible"
        :shareVisible="shareVisible"
        @open-bookmark="openBookmark"
        @export-bookmark="exportBookmark"
        @share-bookmark="shareBookmark"
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
        resolved: Bookmark | null;
    } {
        return {
            resolved: null,
        };
    },
    computed: {
        visibleBookmark(): null | Bookmark {
            if (this.bookmark) {
                return this.bookmark;
            }
            if (this.resolved) {
                return this.resolved;
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
            console.log(`viz: bookmark-resolving`, this.token);
            const savedBookmark = await this.$services.api.resolveBookmark(this.token);
            console.log(`viz: bookmark-resolved`, savedBookmark);
            this.resolved = deserializeBookmark(savedBookmark.bookmark);
        },
        async openBookmark(bookmark: Bookmark): Promise<void> {
            console.log(`viz: open-bookmark-saving`, bookmark);
            const savedBookmark = await this.$services.api.saveBookmark(serializeBookmark(bookmark));
            console.log(`viz: open-bookmark-saved`, savedBookmark.token);
            await this.$router.push({ name: "exploreShortBookmark", query: { v: savedBookmark.token } });
        },
        async exportBookmark(bookmark: Bookmark): Promise<void> {
            const encoded = serializeBookmark(this.bookmark);
            await this.$router.push({ name: "exportBookmark", query: { bookmark: encoded } });
        },
        async shareBookmark(bookmark: Bookmark): Promise<void> {
            const encoded = serializeBookmark(this.bookmark);
            await this.$router.push({ name: "shareBookmark", query: { bookmark: encoded } });
        },
    },
});
</script>

<style lang="scss"></style>
