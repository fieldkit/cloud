<template>
    <StandardLayout @show-station="showStation" :defaultShowStation="false" :disableScrolling="exportsVisible">
        <ExportPanel v-if="exportsVisible" containerClass="exports-floating" :bookmark="bookmark" @close="closeExports" />
        <div class="explore-view">
            <div class="explore-header">
                <DoubleHeader title="Data View">
                    <div class="button" @click="openExports">Export</div>
                </DoubleHeader>
            </div>

            <div v-if="showNoSensors" class="notification">Oh snap, this station doesn't appear to have any sensors to show you.</div>

            <div v-if="!workspace && !bookmark">Nothing selected to visualize, please choose a station or project from the left.</div>

            <VizWorkspace v-if="workspace && !workspace.empty" :workspace="workspace" @change="onChange" />
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
import * as ActionTypes from "@/store/actions";
import { ExportDataAction } from "@/store/typed-actions";
import { GlobalState } from "@/store/modules/global";

import { Workspace, Bookmark } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

export default Vue.extend({
    name: "ExploreView",
    components: {
        ...CommonComponents,
        StandardLayout,
        VizWorkspace,
        ExportPanel,
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
    data: () => {
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
    },
    watch: {
        $route(to, from) {
            console.log("viz: route");
        },
    },
    beforeMount(this: any) {
        if (this.bookmark) {
            return this.$services.api.getAllSensors().then((sensorKeys) => {
                if (this.bookmark.s.length > 0) {
                    return this.showStation(this.bookmark.s[0]);
                }
                if (this.bookmark.g.length > 0) {
                    this.workspace = Workspace.fromBookmark(sensorKeys, this.bookmark);
                }
            });
        }
    },
    methods: {
        onChange(bookmark: Bookmark) {
            const activeEncoded = JSON.stringify(this.bookmark);
            const newEncoded = JSON.stringify(bookmark);

            console.log("viz: new", newEncoded);
            console.log("viz: old", activeEncoded);

            if (newEncoded == activeEncoded) {
                return Promise.resolve(this.workspace);
            }

            return this.$router.push({ name: "exploreBookmark", params: { bookmark: newEncoded } }).then(() => this.workspace);
        },
        openExports() {
            const encoded = JSON.stringify(this.bookmark);
            return this.$router.push({ name: "exportBookmark", params: { bookmark: encoded } });
        },
        closeExports() {
            const encoded = JSON.stringify(this.bookmark);
            return this.$router.push({ name: "exploreBookmark", params: { bookmark: encoded } });
        },
        createWorkspaceIfNecessary() {
            if (this.workspace) {
                return Promise.resolve(this.workspace);
            }

            return this.$services.api.getAllSensors().then((sensorKeys) => {
                this.workspace = new Workspace(sensorKeys);
                return this.workspace;
            });
        },
        showStation(stationId: number) {
            console.log("viz: show-station", stationId);

            const station = this.$store.state.stations.stations[stationId];

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

                    const stations = [stationId];
                    const sensors = [quickSensors.stations[stationId][0].sensorId];

                    return workspace
                        .addStandardGraph(stations, sensors)
                        .eventually((ws) => this.onChange(ws.bookmark()))
                        .then((ws) => ws.query());
                });
            });
        },
    },
});
</script>

<style>
.graph .x-axis {
    color: #7f7f7f;
}
.graph .y-axis {
    color: #7f7f7f;
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
    height: 100%;
    display: flex;
    flex-direction: column;
}
.tree-container {
    flex: 0;
}
.loading-scrubber {
    padding: 20px;
}
.groups-container {
    border: 2px solid #efefef;
    border-radius: 4px;
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
    border: 1px solid #d8dce0;
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
}

.controls-container .row-2 {
    margin-top: 5px;
    padding: 10px;
    align-items: top;
}

.controls-container .tree {
    flex-basis: 30%;
}

.controls-container .left {
    display: flex;
    align-items: baseline;
}

.controls-container .right {
    display: flex;
    justify-content: center;
    align-items: baseline;
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

.controls-container .time-picker {
    margin-left: 20px;
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

.viz.map {
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
    border-left: 2px solid #d8dce0;
    z-index: 10;
    overflow-y: scroll;
    width: 30em;
}
</style>
