<template>
    <StandardLayout @show-station="showStation" :defaultShowStation="false" :disableScrolling="exportsVisible || shareVisible">
        <ExportPanel v-if="exportsVisible" containerClass="exports-floating" :bookmark="bookmark" @close="closePanel" />

        <SharePanel v-if="shareVisible" containerClass="share-floating" :token="token" :bookmark="bookmark" @close="closePanel" />

        <div class="explore-view">
            <div class="explore-header">
                <DoubleHeader :backTitle="$t(backLabelKey)" @back="onBack">
                    <template v-slot:title>
                        <div class="one">
                            Data View
                            <div class="button compare" alt="Add Chart" @click="addChart">
                                <img :src="addIcon" />
                                <div>Add Chart</div>
                            </div>
                        </div>
                    </template>
                    <template v-slot:default>
                        <div class="button" @click="openShare">Share</div>
                        <div class="button" @click="openExports">Export</div>
                    </template>
                </DoubleHeader>
            </div>

            <div v-if="showNoSensors" class="notification">Oh snap, this station doesn't appear to have any sensors to show you.</div>

            <div v-if="!workspace && !bookmark">Nothing selected to visualize, please choose a station or project from the left.</div>

            <div v-bind:class="{ 'workspace-container': true, busy: busy }">
                <div class="busy-panel">&nbsp;</div>

                <VizWorkspace v-if="workspace && !workspace.empty" :workspace="workspace" @change="onChange" />

                <Comments :parentData="bookmark" :user="user" @viewDataClicked="onChange" v-if="bookmark"></Comments>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Promise from "bluebird";

import Vue from "vue";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";
import ExportPanel from "./ExportPanel.vue";
import SharePanel from "./SharePanel.vue";
import { callStationsStations } from "../shared/partners/StationOrSensor.vue";
import { mapState, mapGetters } from "vuex";
import { GlobalState } from "@/store/modules/global";

import { SensorsResponse } from "./api";
import { Workspace, Bookmark, Time, VizSensor, TimeRange, ChartType, FastTime, serializeBookmark } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

import Comments from "../comments/Comments.vue";

export default Vue.extend({
    name: "ExploreWorkspace",
    components: {
        ...CommonComponents,
        StandardLayout,
        VizWorkspace,
        SharePanel,
        ExportPanel,
        Comments,
    },
    props: {
        token: {
            type: String,
            required: false,
        },
        bookmark: {
            type: Bookmark,
            required: true,
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
        workspace: Workspace | null;
        showNoSensors: boolean;
    } {
        return {
            workspace: null,
            showNoSensors: false,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
        }),
        addIcon(): unknown {
            return this.$loadAsset("icon-compare.svg");
        },
        busy(): boolean {
            return !this.workspace || this.workspace.busy;
        },
        backLabelKey(): string {
            if (this.bookmark && this.bookmark.p.length > 0) {
                return "layout.backProjectDashboard";
            }
            return callStationsStations() ? "layout.backToStations" : "layout.backToSensors";
        },
    },
    watch: {
        async bookmark(newValue: Bookmark, oldValue: Bookmark): Promise<void> {
            console.log(`viz: bookmark-route(ew):`, newValue);
            if (this.workspace) {
                await this.workspace.updateFromBookmark(newValue);
            }
        },
    },
    async beforeMount(): Promise<void> {
        if (this.bookmark) {
            await this.$services.api
                .getAllSensors()
                .then(async (sensorKeys) => {
                    // Check for a bookmark that is just to a station with no groups.
                    if (this.bookmark.s.length > 0 && this.bookmark.g.length == 0) {
                        console.log("viz: before-show-station", this.bookmark);
                        return this.showStation(this.bookmark.s[0]);
                    }
                    console.log("viz: before-create-workspace", this.bookmark);
                    await this.createWorkspaceIfNecessary();
                })
                .catch(async (e) => {
                    if (e.name === "ForbiddenError") {
                        await this.$router.push({ name: "login", params: { errorMessage: String(this.$t("login.privateStation")) } });
                    }
                });
        }
    },
    methods: {
        async onBack() {
            console.log("viz:back", this.bookmark);
            if (this.bookmark.p && this.bookmark.p.length > 0) {
                await this.$router.push({ name: "viewProject", params: { id: this.bookmark.p[0] } });
            } else {
                await this.$router.push({ name: "mapAllStations" });
            }
        },
        async addChart() {
            console.log("viz: add");
            if (!this.workspace) throw new Error("viz-add: no workspace");
            return this.workspace.addChart().query();
        },
        async onChange(bookmark: Bookmark): Promise<void> {
            console.log("viz: bookmark-change", bookmark);
            if (Bookmark.sameAs(this.bookmark, bookmark)) {
                return Promise.resolve(this.workspace);
            }
            await this.openBookmark(bookmark);
        },
        async openBookmark(bookmark: Bookmark): Promise<void> {
            this.$emit("open-bookmark", bookmark);
        },
        async openExports(): Promise<void> {
            this.$emit("export");
        },
        async openShare(): Promise<void> {
            this.$emit("share");
        },
        async closePanel(): Promise<void> {
            return await this.openBookmark(this.bookmark);
        },
        async createWorkspaceIfNecessary(): Promise<Workspace> {
            if (this.workspace) {
                return this.workspace;
            }

            const allSensors: SensorsResponse = await this.$services.api.getAllSensors();
            if (this.bookmark) {
                this.workspace = Workspace.fromBookmark(allSensors, this.bookmark);
            } else {
                this.workspace = new Workspace(allSensors);
            }

            console.log(`viz: workspace-created`);

            return this.workspace;
        },
        async showStation(stationId: number): Promise<void> {
            console.log("viz: show-station", stationId);

            return await this.createWorkspaceIfNecessary()
                .then(async (workspace) => {
                    return await this.$services.api.getQuickSensors([stationId]).then(async (quickSensors) => {
                        console.log("viz: quick-sensors", quickSensors);
                        if (quickSensors.stations[stationId].length == 0) {
                            console.log("viz: no sensors TODO: FIX");
                            this.showNoSensors = true;
                            return Promise.delay(5000).then(() => {
                                this.showNoSensors = false;
                            });
                        }

                        const vizSensor: VizSensor = [
                            stationId,
                            [quickSensors.stations[stationId][0].moduleId, quickSensors.stations[stationId][0].sensorId],
                        ];

                        const associated = await this.$services.api.getAssociatedStations(stationId);
                        const stationIds = associated.stations.map((station) => station.id);
                        console.log(`viz: show-station-associated`, associated, stationIds);

                        const bookmark = new Bookmark(
                            this.bookmark.v,
                            [[[[[vizSensor], [Time.Min, Time.Max], [], ChartType.TimeSeries, FastTime.All]]]],
                            stationIds,
                            this.bookmark.p
                        );

                        this.$emit("open-bookmark", bookmark);
                    });
                })
                .catch(async (e) => {
                    if (e.name === "ForbiddenError") {
                        await this.$router.push({ name: "login", params: { errorMessage: String(this.$t("login.privateStation")) } });
                    }
                });
        },
    },
});
</script>

<style lang="scss">
@import "../../scss/layout";

#vg-tooltip-element {
    background-color: #f4f5f7;
    border-radius: 1px;
    box-shadow: none;
    border: none;
}
#vg-tooltip-element .key {
    display: none;
}
#vg-tooltip-element table tr:first-of-type td.value {
    text-align: center;
    font-family: "Avenir", sans-serif;
    font-size: 16px;
    color: #2c3e50;
}
#vg-tooltip-element table tr:nth-of-type(2) td.value {
    text-align: center;
    font-family: "Avenir", sans-serif;
    font-size: 13px;
    color: #2c3e50;
}

.explore-view {
    text-align: left;
    background-color: #fcfcfc;
    padding: 40px;
    flex-grow: 1;
}
.explore-header {
    margin-bottom: 1em;
}
.explore-header .button {
    margin-left: 20px;
    font-size: 12px;
    padding: 5px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}

.workspace-container {
    display: flex;
    flex-direction: column;
    border-radius: 1px;
    box-shadow: 0 0 4px 0 rgba(0, 0, 0, 0.07);
    border: solid 1px #f4f5f7;
}
.tree-container {
    flex: 0;
}
.loading-scrubber {
    padding: 20px;
}
.groups-container {
    background-color: white;
}

.icons-container {
    margin-top: 30px;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: row;
    padding-left: 20px;
    padding-right: 20px;
}

.icons-container.linked {
    background-color: #efefef;
    height: 1px;
}

.icons-container.unlinked {
    border-top: 1px solid #efefef;
    border-bottom: 1px solid #efefef;
    background-color: #fcfcfc;
    height: 8px;
}

.icons-container div:first-child {
    margin-right: auto;
    visibility: hidden;
}
.icons-container div:last-child {
    margin-left: auto;
}
.icons-container .icon {
    background-color: #fcfcfc;
    box-shadow: inset 0 1px 3px 0 rgba(0, 0, 0, 0.11);
    border: 1px solid var(--color-border);
    border-radius: 50%;
    cursor: pointer;
    font-size: 28px;
    padding: 5px;

    &:before {
        color: var(--color-dark);
    }
}

.icons-container .remove-icon {
    background-position: center;
    background-image: url(../../assets/Icon_Close_Circle.png);
    background-size: 20px;
    background-repeat: no-repeat;
    width: 20px;
    height: 20px;
}

.vega-embed {
    width: 100%;

    summary {
        z-index: 0 !important;
        margin-left: 0.25em;
        margin-right: 0.5em;
    }
}
.graph .vega-embed {
    height: 340px;
}
.scrubber .vega-embed {
    height: 40px;

    summary {
        display: none;
    }
}

.workspace-container {
    position: relative;

    .busy-panel {
        position: absolute;
        width: 100%;
        height: 100%;
        display: none;
        z-index: 5;
        opacity: 0.5;
    }

    &.busy .busy-panel {
        display: block;
        background-color: #efefef;
    }

    .viz-loading {
        height: 300px;
        display: flex;
        align-items: center;
    }
}

.controls-container {
    margin-left: 40px;
    margin-right: 40px;
    margin-bottom: 10px;
}

.controls-container .row {
    display: flex;
    justify-content: space-between;
}

.controls-container .row-1 {
    padding: 10px;
    border-bottom: 1px solid #efefef;
    margin-bottom: 5px;
    align-items: baseline;
    min-height: 80px;
}

.controls-container .row-2 {
    margin-top: 5px;
    padding: 10px;
}

.controls-container .left {
    display: flex;
    align-items: center;
    flex-direction: column;
}

.controls-container .left .row {
    align-items: center;
    display: flex;

    .actions {
        margin-left: 1em;
        display: flex;
        align-items: center;

        .button {
            margin-bottom: 0;
        }
    }
}

.controls-container .half {
    flex-basis: 50%;
}

.controls-container .tree-pair {
    display: flex;
    align-items: center;
    width: 100%;
}

.controls-container .tree-pair div {
    flex-basis: 50%;
}

.controls-container .right {
    display: flex;
    justify-content: flex-end;
    align-items: center;
}

.controls-container .right .chart-type {
    flex-basis: 25%;
}

.controls-container .right {
    font-size: 12px;
}

.controls-container .right.half {
    align-items: flex-start;
}

.controls-container .view-by {
    margin: 0px 10px 0 10px;
}

.controls-container .fast-time {
    margin: 0px 10px 0 10px;
    cursor: pointer;
}

.controls-container .date-picker {
    margin-left: 20px;
}

.controls-container .date-picker input {
    padding: 5px;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}

.controls-container .fast-time.selected {
    text-decoration: underline;
    font-weight: bold;
}

.controls-container .left .button {
    margin-right: 20px;
    font-size: 12px;
    padding: 5px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}

.debug-panel {
    font-size: 8px;
}

.notification {
    margin-top: 20px;
    margin-bottom: 20px;
    padding: 20px;
    background-color: #f8d7da;
    border: 2px;
    border-radius: 4px;
}

.svg-container {
    display: inline-block;
    position: relative;
    width: 100%;
    vertical-align: top;
    overflow: hidden;
}

.svg-content-responsive {
    display: inline-block;
    position: absolute;
    top: 10px;
    left: 0;
}

.viz.map .viz-map {
    height: 400px;
    margin-bottom: 20px;
}

.share-floating,
.exports-floating {
    position: absolute;
    right: 0;
    top: 70px;
    bottom: 0;
    background-color: #fcfcfc;
    border-left: 2px solid var(--color-border);
    z-index: 10;
    overflow-y: scroll;
    width: 30em;
}

.loading-options {
    text-align: center;
    color: #afafaf;
}

.button.compare {
    display: flex;
    align-items: center;

    div {
        padding-left: 1em;
    }
}

.brush_brush_bg path {
    body.floodnet & {
        fill: var(--color-border);
        fill-opacity: 1;
    }
}

.layer_1_marks path {
    fill: var(--color-primary);
}

.one {
    display: flex;
    flex-direction: row;
}
</style>
