<template>
    <StandardLayout @show-station="showStation" :defaultShowStation="false">
        <div class="explore-view">
            <div class="explore-header">
                <div class="left">
                    <h2>Data Visualization</h2>
                </div>
            </div>

            <div v-if="showNoSensors" class="notification">
                Oh snap, this station doesn't appear to have any sensors to show you.
            </div>

            <div v-if="!workspace">
                Nothing selected to visualize, please choose a station or project from the left.
            </div>

            <VizWorkspace v-if="workspace" :workspace="workspace" @change="onChange" />
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import Promise from "bluebird";
import StandardLayout from "../StandardLayout.vue";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

import FKApi from "@/api/api";

import { Workspace, Bookmark } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

export default Vue.extend({
    name: "ExploreView",
    components: {
        StandardLayout,
        VizWorkspace,
    },
    props: {
        bookmark: {
            type: Bookmark,
            required: false,
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
            stations: (s: GlobalState) => s.stations.stations.user,
            userProjects: (s: GlobalState) => s.stations.projects.user,
        }),
    },
    beforeMount() {
        if (this.bookmark) {
            return new FKApi().getAllSensors().then((sensorKeys) => {
                this.workspace = Workspace.fromBookmark(sensorKeys, this.bookmark);
            });
        }
    },
    methods: {
        onChange(bookmark: Bookmark) {
            const activeEncoded = JSON.stringify(this.bookmark);
            const newEncoded = JSON.stringify(bookmark);

            if (newEncoded == activeEncoded) {
                return Promise.resolve(this.workspace);
            }

            return this.$router.push({ name: "exploreBookmark", params: { bookmark: newEncoded } }).then(() => this.workspace);
        },
        createWorkspaceIfNecessary() {
            if (this.workspace) {
                return Promise.resolve(this.workspace);
            }

            return new FKApi().getAllSensors().then((sensorKeys) => {
                this.workspace = new Workspace(sensorKeys);
                return this.workspace;
            });
        },
        showStation(stationId: number, ...args) {
            console.log("show-station", stationId, ...args);

            const station = this.$store.state.stations.stations.all[stationId];

            return this.createWorkspaceIfNecessary().then((workspace) => {
                return new FKApi().getQuickSensors([stationId]).then((quickSensors) => {
                    console.log("quick-sensors", quickSensors);
                    if (quickSensors.stations[stationId].length == 0) {
                        console.log("no sensors");
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
    display: flex;
    justify-content: space-between;
    align-items: baseline;
}
.explore-header .left {
    display: flex;
}
.explore-header .right {
    display: flex;
}
.explore-header .btn {
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
    background-image: url(../../assets/open_link.png);
    background-size: 30px;
    background-repeat: no-repeat;
    width: 40px;
    height: 40px;
}
.icons-container .link-icon {
    background-position: center;
    background-image: url(../../assets/link.png);
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

.controls-container .left .btn {
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
</style>
