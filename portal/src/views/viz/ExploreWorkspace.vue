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
                        <div class="button-submit" @click="openShare">
                            <i class="icon icon-share"></i>
                            Share
                        </div>
                        <div v-show="user" class="button-submit" @click="openExports">
                            <i class="icon icon-export"></i>
                            Export
                        </div>
                    </template>
                </DoubleHeader>
            </div>

            <div v-if="showNoSensors" class="notification">Oh snap, this station doesn't appear to have any sensors to show you.</div>

            <div v-if="!workspace && !bookmark">Nothing selected to visualize, please choose a station or project from the left.</div>

            <div v-bind:class="{ 'workspace-container': true, busy: busy }">
                <div class="busy-panel">&nbsp;</div>
                <div class="station-summary" v-if="selectedStation">
                    <StationSummaryContent :station="selectedStation" v-if="workspace && !workspace.empty" class="summary-content" />
                    <div class="pagination" v-if="workspace && !workspace.empty">
                        <PaginationControls
                            :page="selectedIndex"
                            :totalPages="getValidStations().length"
                            @new-page="onNewSummaryStation"
                            textual
                            wrap
                        />
                    </div>
                </div>

                <VizWorkspace v-if="workspace && !workspace.empty" :workspace="workspace" @change="onChange" />

                <Comments :parentData="bookmark" :user="user" @viewDataClicked="onChange" v-if="bookmark"></Comments>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Promise from "bluebird";

import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";
import ExportPanel from "./ExportPanel.vue";
import SharePanel from "./SharePanel.vue";
import StationSummaryContent from "../shared/StationSummaryContent.vue";
import PaginationControls from "@/views/shared/PaginationControls.vue";
import { getPartnerCustomization } from "../shared/partners";
import { mapState, mapGetters } from "vuex";
import { Station, ActionTypes } from "@/store";
import { GlobalState } from "@/store/modules/global";
import { SensorsResponse } from "./api";
import { ForbiddenError } from "@/api";
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
        StationSummaryContent,
        PaginationControls,
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
        selectedIndex: number;
        validStations: number[];
    } {
        return {
            workspace: null,
            showNoSensors: false,
            selectedIndex: 0,
            validStations: [],
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
            const partnerCustomization = getPartnerCustomization();
            if (this.bookmark && this.bookmark.c) {
                if (!this.bookmark.c.map) {
                    return "layout.backProjectDashboard";
                }
            }
            if (partnerCustomization) {
                return partnerCustomization.nav.viz.back.map.label;
            }
            return "layout.backToStations";
        },
        selectedId(): number {
            return Number(_.flattenDeep(this.bookmark.g)[0]);
        },
        selectedStation(): Station | null {
            if (this.workspace) {
                return this.workspace.getStation(this.selectedId);
            }
            return null;
        },
    },
    watch: {
        async bookmark(newValue: Bookmark, oldValue: Bookmark): Promise<void> {
            console.log(`viz: bookmark-route(ew):`, newValue);
            if (this.workspace) {
                await this.workspace.updateFromBookmark(newValue);
            }
        },
        async selectedId(newValue: number, oldValue: number): Promise<void> {
            console.log("viz: selected-changed-associated", newValue);
            if (this.workspace) {
                await this.includeAssociated(this.workspace);
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
                    const ws = await this.createWorkspaceIfNecessary();
                })
                .catch(async (e) => {
                    if (e.name === "ForbiddenError") {
                        await this.$router.push({ name: "login", params: { errorMessage: String(this.$t("login.privateStation")) } });
                    }
                });
        }
    },
    methods: {
        async includeAssociated(ws: Workspace): Promise<Workspace> {
            const allStationIds = ws.allStationIds;
            console.log("viz: include-associated(0)", allStationIds);

            if (allStationIds.length > 0) {
                const associated = await this.$services.api.getAssociatedStations(allStationIds[0]).catch(async (e) => {
                    if (ForbiddenError.isInstance(e)) {
                        await this.$router.push({
                            name: "login",
                            params: { errorMessage: String(this.$t("login.privateStation")) },
                            query: { after: this.$route.path },
                        });
                    }
                });
                if (associated) {
                    const ids = associated.stations.map((s) => s.station.id);
                    console.log(`viz: include-associated(0)`, associated);
                    console.log(
                        `viz: include-associated(0)`,
                        associated.stations
                            .filter((row) => row.manual)
                            .map((row) => {
                                return { station: row.station, manual: row.manual };
                            }),
                        associated.stations
                            .filter((row) => row.location)
                            .map((row) => {
                                return { station: row.station, location: row.location };
                            })
                    );
                    await ws.addStationIds(ids);
                    await ws.addFullStations(associated.stations.map((s) => s.station));
                    await ws.addAssociatedStations(associated.stations);
                }
            }

            return ws;
        },
        async onBack() {
            if (this.bookmark.c) {
                if (this.bookmark.c.map) {
                    await this.$router.push({ name: "viewProjectBigMap", params: { id: this.bookmark.c.project } });
                } else {
                    await this.$router.push({ name: "viewProject", params: { id: this.bookmark.c.project } });
                }
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
            if (Bookmark.sameAs(this.bookmark, bookmark)) {
                console.log("viz: bookmark-no-change", bookmark);
                return Promise.resolve(this.workspace);
            }
            console.log("viz: bookmark-change", bookmark);
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
            const ws = this.bookmark ? Workspace.fromBookmark(allSensors, this.bookmark) : new Workspace(allSensors);

            console.log(`viz: workspace-created`);

            await this.includeAssociated(ws);

            this.workspace = ws;

            return ws;
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

                        const sensorModuleId = quickSensors.stations[stationId][0].moduleId;
                        const sensorId = quickSensors.stations[stationId][0].sensorId;
                        const vizSensor: VizSensor = [stationId, [sensorModuleId, sensorId]];

                        const associated = await this.$services.api.getAssociatedStations(stationId);
                        const stationIds = associated.stations.map((associatedStation) => associatedStation.station.id);
                        console.log(`viz: show-station-associated`, associated, stationIds);

                        const getInitialBookmark = () => {
                            const quickSensor = quickSensors.stations[stationId].filter((qs) => qs.sensorId == sensorId);
                            if (quickSensor.length == 1) {
                                const end = new Date(quickSensor[0].sensorReadAt);
                                const start = new Date(end);

                                start.setDate(end.getDate() - 14); // TODO Use getFastTime

                                return new Bookmark(
                                    this.bookmark.v,
                                    [[[[[vizSensor], [start.getTime(), end.getTime()], [], ChartType.TimeSeries, FastTime.TwoWeeks]]]],
                                    stationIds,
                                    this.bookmark.p,
                                    this.bookmark.c
                                );
                            }

                            console.log("viz: ERROR missing expected quick row, default to FastTime.All");

                            return new Bookmark(
                                this.bookmark.v,
                                [[[[[vizSensor], [Time.Min, Time.Max], [], ChartType.TimeSeries, FastTime.All]]]],
                                stationIds,
                                this.bookmark.p,
                                this.bookmark.c
                            );
                        };

                        this.$emit("open-bookmark", getInitialBookmark());
                    });
                })
                .catch(async (e) => {
                    if (e.name === "ForbiddenError") {
                        await this.$router.push({ name: "login", params: { errorMessage: String(this.$t("login.privateStation")) } });
                    }
                });
        },
        getValidStations(): number[] {
            const validStations = Object.entries(this.workspace.stations)
                .filter(([key, station]) => station.sensors.length > 0)
                .map((d) => +d[0]);

            this.selectedIndex = validStations.indexOf(this.selectedId);

            return validStations;
        },
        onNewSummaryStation(evt) {
            const stations = this.getValidStations();
            this.showStation(stations[evt]);
            this.selectedIndex = evt;
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
    text-align: left;
    font-family: "Avenir", sans-serif;
    color: #2c3e50;

    h3 {
        margin-top: 0;
        margin-bottom: 0;
        margin-right: 5px;
        font-size: 13px;
        display: flex;
        align-items: center;
        line-height: 24px;
    }
    .tooltip-color {
        margin-right: 5px;
        font-size: 2em;
    }
    p {
        margin: 0.3em;
    }
    p.value {
        font-size: 16px;
    }
    p.time {
        font-size: 13px;
    }
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
}

.controls-container .row-1 {
    padding: 10px;
    border-bottom: 1px solid #efefef;
    margin-bottom: 5px;
    align-items: center;
    min-height: 60px;
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

.controls-container .tree-pair {
    display: flex;
    align-items: center;
    width: 100%;
    flex: 0 0 500px;
}

.controls-container .tree-pair > div {
    flex: 0 1 auto;
}

.tree-key {
    flex-basis: 0;
    margin-right: 15px;
    line-height: 35px;
    font-size: 40px;

    body.floodnet & {
        margin-top: -10px;
    }
}

.controls-container .right {
    display: flex;
    justify-content: flex-end;
    align-items: center;

    &.time {
        margin-left: auto;
    }
}

.controls-container .right {
    font-size: 12px;
}

.controls-container .right.half {
    align-items: flex-start;
    flex: 0 0 110px;
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

    span {
        &:nth-of-type(1) {
            margin-right: 5px;
        }
    }

    .vc-day-layer {
        left: -2px;
    }
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
        fill: var(--color-primary);
        fill-opacity: 1;
    }
}

.layer_1_marks path {
    fill: var(--color-primary);

    body.floodnet & {
        fill: #3f5d62;
    }
}

.one {
    display: flex;
    flex-direction: row;
}

.button-submit {
    padding: 0 28px;

    &:nth-child(n + 1) {
        margin-left: 20px;
    }
}
.station-summary {
    background-color: #fff;
    border-bottom: 1px solid var(--color-border);
    padding: 20px;
    display: flex;
    justify-content: space-between;

    .summary-content {
        align-items: center;
    }
    @include bp-down($sm) {
        flex-direction: column;

        .pagination {
            margin-top: 0.5em;
        }
    }
}
.pagination {
    display: flex;
    margin-right: 13px;
    justify-content: center;
}
</style>
