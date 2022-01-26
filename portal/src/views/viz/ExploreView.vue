<template>
    <StandardLayout @show-station="showStation" :defaultShowStation="false" :disableScrolling="exportsVisible">
        <ExportPanel v-if="exportsVisible" containerClass="exports-floating" :bookmark="bookmark" @close="closeExports" />
        <div class="container-wrap explore-view">
            <div class="explore-header">
                <DoubleHeader
                    title="Data View"
                    :backTitle="backRoute === 'viewProject' ? $t('layout.backProjectDashboard') : $t('layout.backToStations')"
                    :backRoute="backRoute"
                    :backRouteParams="backRouteParams"
                >
                    <div class="button" @click="openExports">Export</div>
                </DoubleHeader>
            </div>

            <div v-if="showNoSensors" class="notification">Oh snap, this station doesn't appear to have any sensors to show you.</div>

            <div v-if="!workspace && !bookmark">Nothing selected to visualize, please choose a station or project from the left.</div>

            <div class="workspace-container">
                <VizWorkspace v-if="workspace && !workspace.empty" :workspace="workspace" @change="onChange" />
                <div class="busy" v-else><Spinner /></div>

                <Comments :parentData="bookmark" :user="user" @viewDataClicked="onChange" v-if="user"></Comments>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import Promise from "bluebird";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";
import ExportPanel from "./ExportPanel.vue";

import { mapState, mapGetters } from "vuex";
import { GlobalState } from "@/store/modules/global";

import { SensorsResponse } from "./api";
import { Workspace, Bookmark } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

import Comments from "../comments/Comments.vue";

export default Vue.extend({
    name: "ExploreView",
    components: {
        ...CommonComponents,
        StandardLayout,
        VizWorkspace,
        ExportPanel,
        Comments,
    },
    props: {
        bookmark: {
            type: Bookmark,
            required: false,
        },
        exportsVisible: {
            type: Boolean,
            default: false,
        },
    },
    beforeRouteEnter(to, from, next) {
        next((vm: any) => {
            if (from.name === "exploreBookmark" || from.name === "exportBookmark") {
                return;
            }
            vm.backRoute = from.name ? from.name : "mapAllStations";
            vm.backRouteParams = from.params;
        });
    },
    data(): {
        workspace: Workspace | null;
        showNoSensors: boolean;
        stationId: number | null;
        backRoute: string | null;
        backRouteParams: { bounds: string | null; id: number | null };
    } {
        return {
            workspace: null,
            showNoSensors: false,
            stationId: null,
            backRoute: null,
            backRouteParams: { bounds: null, id: null },
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
        }),
    },
    watch: {
        async bookmark(newValue: Bookmark, oldValue: Bookmark): Promise<void> {
            console.log(`viz: route:`, newValue);
            if (this.workspace) {
                await this.workspace.updateFromBookmark(newValue);
            }
        },
    },
    async beforeMount(): Promise<void> {
        if (this.bookmark) {
            await this.$services.api.getAllSensors().then(async (sensorKeys) => {
                // Check for a bookmark that is just to a station with no groups.
                if (this.bookmark.s.length > 0) {
                    console.log("viz-before-show-station", this.bookmark);
                    return this.showStation(this.bookmark.s[0]);
                }
                console.log("viz-before-create-workspace", this.bookmark);
                await this.createWorkspaceIfNecessary();
            });
        }
    },
    methods: {
        async onChange(bookmark: Bookmark): Promise<void> {
            if (Bookmark.sameAs(this.bookmark, bookmark)) {
                return Promise.resolve(this.workspace);
            }
            await this.$router.push({ name: "exploreBookmark", params: { bookmark: JSON.stringify(bookmark) } }).then(() => this.workspace);
        },
        async openExports(): Promise<void> {
            const encoded = JSON.stringify(this.bookmark);
            await this.$router.push({ name: "exportBookmark", params: { bookmark: encoded } });
        },
        async closeExports(): Promise<void> {
            const encoded = JSON.stringify(this.bookmark);
            await this.$router.push({ name: "exploreBookmark", params: { bookmark: encoded } });
        },
        async createWorkspaceIfNecessary(): Promise<Workspace> {
            if (this.workspace) {
                return this.workspace;
            }

            const allSensors: SensorsResponse = await this.$services.api.getAllSensors();
            if (this.bookmark) {
                this.workspace = Workspace.fromBookmark(allSensors, this.bookmark);
                void this.includeAssociatedStations(this.workspace);
            } else {
                this.workspace = new Workspace(allSensors);
            }

            return this.workspace;
        },
        async includeAssociatedStations(ws: Workspace): Promise<void> {
            console.log(`include-associated`, ws.allStationIds);
            if (ws.allStationIds.length != 1) {
                return;
            }

            return this.$services.api.getAssociatedStations(ws.allStationIds[0]).then((associated) => {
                console.log("quick-sensors-associated", associated);
                const ids = associated.stations.map((s) => s.id);
                return ws.addStationIds(ids);
            });
        },
        async showStation(stationId: number): Promise<void> {
            console.log("viz: show-station", stationId);
            this.stationId = stationId;

            return this.createWorkspaceIfNecessary().then((workspace) => {
                return this.$services.api.getQuickSensors([stationId]).then((quickSensors) => {
                    console.log("viz: quick-sensors", quickSensors);
                    if (quickSensors.stations[stationId].length == 0) {
                        console.log("viz: no sensors");
                        this.showNoSensors = true;
                        return Promise.delay(5000).then(() => {
                            this.showNoSensors = false;
                        });
                    }

                    console.log("quick-sensors", quickSensors);

                    const stations = [stationId];
                    const sensors = [[quickSensors.stations[stationId][0].moduleId, quickSensors.stations[stationId][0].sensorId]];

                    return workspace
                        .addStandardGraph(stations, sensors)
                        .eventually((ws) => this.onChange(ws.bookmark()))
                        .then((ws) => Promise.all([ws.query(), this.includeAssociatedStations(workspace)]));
                });
            });
        },
    },
});
</script>

<style lang="scss">
@import "../../scss/layout";

.graph .x-axis {
    color: #7f7f7f;
}

.graph .y-axis {
    color: #7f7f7f;
}

.graph .x-axis text {
    font-size: 7pt;
}

.graph .y-axis text {
    font-size: 7pt;
}

.explore-view {
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
}
.icons-container .unlink-icon {
    background-position: center;
    background-image: url(../../assets/link.png);
    background-size: 30px;
    background-repeat: no-repeat;
    width: 40px;
    height: 40px;
}
.icons-container .link-icon {
    background-position: center;
    background-image: url(../../assets/open_link.png);
    background-size: 30px;
    background-repeat: no-repeat;
    width: 40px;
    height: 40px;
}
.icons-container .remove-icon {
    background-position: center;
    background-image: url(../../assets/Icon_Close_Circle.png);
    background-size: 20px;
    background-repeat: no-repeat;
    width: 20px;
    height: 20px;
}

/* HACK d3 Real talk, no idea how to do this elsewhere. -jlewallen */
.brush-container .selection {
    opacity: 0.3;
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

.controls-container .tree {
    flex-basis: 50%;
}

.controls-container .tree div {
    flex-basis: 50%;
}

.controls-container .left {
    display: flex;
    align-items: center;
}

.controls-container .right {
    display: flex;
    justify-content: center;
    align-items: baseline;
}

.controls-container .row-2 .right {
    flex-basis: 10%;
}

.controls-container .right {
    font-size: 12px;
}

.controls-container .view-by {
    margin: 35px 10px 0 10px;
}

.controls-container .fast-time {
    margin: 35px 10px 0 10px;
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
</style>
