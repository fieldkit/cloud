import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";

import { D3TimeSeriesGraph } from "./D3TimeSeriesGraph";
import { D3Scrubber } from "./D3Scrubber";

export const D3Graph = Vue.extend({
    name: "D3Graph",
    components: {
        D3TimeSeriesGraph,
        D3Scrubber,
        Treeselect,
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
        raiseRemove() {
            return this.$emit("viz-remove");
        },
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        raiseSelected() {
            return this.$emit("viz-changed");
        },
    },
    template: `
		<div class="viz graph">
			<h3 v-if="viz.info.title">{{ viz.info.title }}</h3>
			<div class="controls-container">
				<treeselect v-model="selected" :options="workspace.options" open-direction="bottom" @select="raiseSelected" />
				<div class="btn" @click="raiseRemove">Remove</div>
			</div>
			<D3TimeSeriesGraph :viz="viz" :workspace="workspace" @viz-time-zoomed="raiseTimeZoomed" />
			<D3Scrubber :viz="viz" @viz-time-zoomed="raiseTimeZoomed" />
		</div>
	`,
});
