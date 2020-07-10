<template>
    <StandardLayout>
        <div class="explore-view">
            <div class="explore-header">
                <div class="left">
                    <h2>Default FieldKit Project</h2>
                </div>
                <div class="right">
                    <div class="btn">Share</div>
                    <div class="btn">Export</div>
                </div>
            </div>
            <VizWorkspace v-if="workspace" :workspace="workspace"></VizWorkspace>
        </div>
    </StandardLayout>
</template>

<script>
import StandardLayout from "../StandardLayout";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

import FKApi from "@/api/api";

import { Workspace } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

export default {
    name: "ExploreView",
    components: {
        StandardLayout,
        VizWorkspace,
    },
    props: {},
    data: () => {
        return {
            workspace: null,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s) => s.user.user,
            stations: (s) => s.stations.stations.user,
            userProjects: (s) => s.stations.projects.user,
        }),
    },
    beforeMount() {
        this.$store.watch(
            (state, getters) => state.stations.stations.all[12],
            (station, _old) => {
                return new FKApi().getQuickSensors([station.id]).then((quickSensors) => {
                    return new FKApi().getAllSensors().then((sensors) => {
                        this.workspace = new Workspace(sensors);
                        this.workspace.addStation(quickSensors, station);
                        // return this.workspace.compare().compare().combine();
                        // return this.workspace.compare();
                        return this.workspace.addSensorAll(36);
                    });
                });
            }
        );
    },
    methods: {
        showStation(station) {
            console.log("station");
        },
    },
};
</script>

<style>
.graph .x-axis {
    color: #7f7f7f;
}
.graph .y-axis {
    color: #7f7f7f;
}

.explore-view {
    height: 100%;
    text-align: left;
    background-color: #fcfcfc;
    padding: 40px;
    max-width: 1054px; /* HACK fixed width of svg */
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
    padding-bottom: 20px;
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
</style>
