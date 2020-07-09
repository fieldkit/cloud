import _ from "lodash";
import Vue from "vue";

import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";

import { ViewingControls } from "./ViewingControls";
import { D3TimeSeriesGraph } from "./D3TimeSeriesGraph";
import { D3Scrubber } from "./D3Scrubber";

export const D3Graph = Vue.extend({
    name: "D3Graph",
    components: {
        ViewingControls,
        D3TimeSeriesGraph,
        D3Scrubber,
    },
    data() {
        return {
            selected: null,
        };
    },
    props: {
        viz: {
            type: Graph,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    mounted() {
        this.viz.log("mounted");
    },
    updated() {
        this.viz.log("updated");
    },
    watch: {},
    methods: {
        raiseTimeZoomed(...args) {
            return this.$emit("viz-time-zoomed", ...args);
        },
        raiseRemove(...args) {
            return this.$emit("viz-remove", ...args);
        },
        raiseCompare(...args) {
            return this.$emit("viz-compare", ...args);
        },
        raiseFastTime(...args) {
            return this.$emit("viz-fast-time", ...args);
        },
        raiseChangeSensors(...args) {
            return this.$emit("viz-change-sensors", ...args);
        },
        raiseChangeChart(...args) {
            return this.$emit("viz-change-chart", ...args);
        },
        uiNameOf(graph: Graph) {
            return "D3TimeSeriesGraph";
        },
    },
    template: `
		<div class="viz graph">
			<ViewingControls :viz="viz" :workspace="workspace"
				@viz-remove="raiseRemove"
				@viz-compare="raiseCompare"
				@viz-fast-time="raiseFastTime"
				@viz-time-zoomed="raiseTimeZoomed"
				@viz-change-sensors="raiseChangeSensors"
				@viz-change-chart="raiseChangeChart"
				/>

			<component v-bind:is="uiNameOf(viz)" :viz="viz" :workspace="workspace"
				@viz-time-zoomed="raiseTimeZoomed" />

			<D3Scrubber :viz="viz" @viz-time-zoomed="raiseTimeZoomed" />
		</div>
	`,
});
