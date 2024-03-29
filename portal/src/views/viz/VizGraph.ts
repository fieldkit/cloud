import _ from "lodash";
import Vue from "vue";
import { Graph, Workspace, ChartType } from "./viz";

import Spinner from "@/views/shared/Spinner.vue";

import { ViewingControls } from "./ViewingControls";
import { DebuggingPanel } from "./DebuggingPanel";

import { VegaTimeSeriesGraph as TimeSeriesGraph } from "./VegaTimeSeriesGraph";
import { VegaHistogram as Histogram } from "./VegaHistogram";
import { VegaRange as Range } from "./VegaRange";

import { D3Map as Map } from "./D3Map";

const Loading = Vue.extend({
    name: "Loading",
    components: {
        Spinner,
    },
    template: `<div class="viz-loading"><Spinner /></div>`,
});

export const VizGraph = Vue.extend({
    name: "VizGraph",
    components: {
        ViewingControls,
        Loading,
        DebuggingPanel,
        TimeSeriesGraph,
        Histogram,
        Range,
        // eslint-disable-next-line
        Map,
    },
    data() {
        return {};
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
        this.viz.log("mounted", this.viz);
    },
    updated() {
        this.viz.log("updated", this.viz);
    },
    computed: {
        debug(): boolean {
            return false;
        },
    },
    methods: {
        raiseTimeZoomed(...args: unknown[]): void {
            this.$emit("viz-time-zoomed", ...args);
        },
        raiseTimeDragged(...args: unknown[]): void {
            this.$emit("viz-time-dragged", ...args);
        },
        raiseGeoZoomed(...args: unknown[]): void {
            this.$emit("viz-geo-zoomed", ...args);
        },
        raiseRemove(...args: unknown[]): void {
            this.$emit("viz-remove", ...args);
        },
        raiseCompare(...args: unknown[]): void {
            this.$emit("viz-compare", ...args);
        },
        raiseFastTime(...args: unknown[]): void {
            this.$emit("viz-fast-time", ...args);
        },
        raiseChangeSensors(...args: unknown[]): void {
            this.$emit("viz-change-sensors", ...args);
        },
        raiseChangeChart(...args: unknown[]): void {
            this.$emit("viz-change-chart", ...args);
        },
        uiNameOf(graph: Graph): string {
            if (this.viz.loadedDataSets.length == 0) {
                return "Loading";
            }

            switch (graph.chartType) {
                case ChartType.TimeSeries:
                    return "TimeSeriesGraph";
                case ChartType.Histogram:
                    return "Histogram";
                case ChartType.Range:
                    return "Range";
                case ChartType.Map:
                    return "Map";
                case ChartType.Bar:
                    return "BarChart";
            }

            this.viz.log("unknown chart type");
            return "TimeSeriesGraph";
        },
    },
    template: `
		<div class="viz graph">
			<ViewingControls :viz="viz" :workspace="workspace" v-bind:key="workspace.version"
				@viz-remove="raiseRemove"
				@viz-compare="raiseCompare"
				@viz-fast-time="raiseFastTime"
				@viz-time-zoomed="raiseTimeZoomed"
				@viz-geo-zoomed="raiseGeoZoomed"
				@viz-change-sensors="raiseChangeSensors"
				@viz-change-chart="raiseChangeChart" />

			<component v-bind:is="uiNameOf(viz)" :viz="viz" :workspace="workspace"
				@viz-geo-zoomed="raiseGeoZoomed"
				@viz-time-zoomed="raiseTimeZoomed"
                @viz-time-dragged="raiseTimeDragged" />

            <DebuggingPanel :viz="viz" :workspace="workspace" v-if="debug" />
		</div>
	`,
});
