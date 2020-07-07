import _ from "lodash";
import Vue from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { TimeRange } from "./common";
import { Workspace, Graph, Viz, Querier, TreeOption, Group } from "./viz";
import { D3Scrubber } from "./D3Scrubber";
import { D3MultiScrubber } from "./D3MultiScrubber";
import { D3Graph } from "./D3Graph";

export const VizGroup = Vue.extend({
    components: {
        D3Scrubber,
        D3MultiScrubber,
        D3Graph,
    },
    props: {
        group: {
            type: Group,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    data() {
        return {};
    },
    methods: {
        uiNameOf(viz: Viz) {
            return "D3" + viz.constructor.name;
        },
        raiseRemove(viz: Viz) {
            return this.$emit("viz-remove", viz);
        },
        raiseGroupZoomed(newTimes: TimeRange) {
            return this.$emit("group-time-zoomed", newTimes);
        },
        raiseVizZoomed(viz: Viz, newTimes: TimeRange) {
            return this.$emit("viz-time-zoomed", viz, newTimes);
        },
        raiseSelected(viz: Viz, option: TreeOption) {
            return this.$emit("viz-changed", viz, option);
        },
    },
    template: `
		<div>
			<div v-if="group.scrubbers">
				<D3MultiScrubber :scrubbers="group.scrubbers" @viz-time-zoomed="(range) => raiseGroupZoomed(range)" />
			</div>

			<component v-for="viz in group.vizes" :key="viz.id" v-bind:is="uiNameOf(viz)" :viz="viz" :workspace="workspace"
				@viz-time-zoomed="(range) => raiseVizZoomed(viz, range)"
				@viz-remove="() => raiseRemove(viz)"
				@viz-change="(option) => raiseChange(viz, option)">
			</component>
		</div>
	`,
});

export const VizWorkspace = Vue.extend({
    components: {
        D3Scrubber,
        D3MultiScrubber,
        D3Graph,
        VizGroup,
    },
    props: {
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    data() {
        return {};
    },
    mounted() {
        console.log("workspace: mounted", this.workspace);
        return this.workspace.query();
    },
    methods: {
        onGroupTimeZoomed(group: Group, times: TimeRange) {
            group.log("zooming", times.toArray());
            return this.workspace.groupZoomed(group, times).query();
        },
        onVizTimeZoomed(viz: Viz, times: TimeRange) {
            viz.log("zooming", times.toArray());
            return this.workspace.graphZoomed(viz, times).query();
        },
        onVizRemove(viz: Viz) {
            viz.log("removing");
            return this.workspace.remove(viz).query();
        },
        onVizChange(viz: Viz, option: TreeOption) {
            viz.log("option", option);
            return this.workspace.selected(viz, option).query();
        },
        onCompare() {
            console.log("workspace: compare");
            return this.workspace.compare().query();
        },
    },
    template: `
		<div class="workspace">
			<div class="controls-container">
				<button class="btn" @click="onCompare">Compare</button>
			</div>
			<div class="groups-container">
				<VizGroup v-for="group in workspace.groups" :key="group.id"
					:group="group" :workspace="workspace"
					@viz-remove="onVizRemove"
					@viz-time-zoomed="onVizTimeZoomed"
					@group-time-zoomed="(...args) => onGroupTimeZoomed(group, ...args)"
					@viz-change="onVizChange"
				/>
			</div>
		</div>
	`,
});
