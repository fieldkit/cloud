<template>
    <ExploreWorkspace />
</template>

<script lang="ts">
import { Workspace, Bookmark, serializeBookmark } from "./viz";

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
    data(): {} {
        return {};
    },
    async beforeMount(): Promise<void> {
        if (this.token) {
            console.log(`viz: resolve`, this.token);
            const resolved = await this.$services.api.resolveBookmark(this.token);
            console.log(`viz: resolve`, resolved);
        }
        if (this.bookmark) {
            /*
            await this.$services.api.getAllSensors().then(async (sensorKeys) => {
                // Check for a bookmark that is just to a station with no groups.
                if (this.bookmark.s.length > 0 && this.bookmark.g.length == 0) {
                    console.log("viz: before-show-station", this.bookmark);
                    return this.showStation(this.bookmark.s[0]);
                }
                console.log("viz: before-create-workspace", this.bookmark);
                await this.createWorkspaceIfNecessary();
            });
            */
        }
    },
    methods: {},
});
</script>

<style lang="scss"></style>
