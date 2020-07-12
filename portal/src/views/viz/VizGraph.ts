import _ from "lodash";
import Vue from "vue";

import { TimeRange } from "./common";
import { Graph, QueriedData, Workspace, ChartType } from "./viz";

import { ViewingControls } from "./ViewingControls";
import { D3TimeSeriesGraph } from "./D3TimeSeriesGraph";
import { D3Histogram } from "./D3Histogram";
import { D3Range } from "./D3Range";
import { D3Map } from "./D3Map";
import { D3Scrubber } from "./D3Scrubber";
import { DebuggingPanel } from "./DebuggingPanel";

export const VizGraph = Vue.extend({
    name: "VizGraph",
    components: {
        ViewingControls,
        DebuggingPanel,
        D3TimeSeriesGraph,
        D3Histogram,
        D3Range,
        D3Map,
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
    computed: {
        debug() {
            return false;
        },
    },
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
            switch (graph.chartType) {
                case ChartType.TimeSeries:
                    return "D3TimeSeriesGraph";
                case ChartType.Histogram:
                    return "D3Histogram";
                case ChartType.Range:
                    return "D3Range";
                case ChartType.Map:
                    return "D3Map";
            }
            this.viz.log("unknown chart type");
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

            <DebuggingPanel :viz="viz" :workspace="workspace" v-if="debug" />
		</div>
	`,
});
