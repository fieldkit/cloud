import _ from "lodash";
import Vue from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { Workspace, Viz, Querier } from "./viz";
import { D3Scrubber } from "./D3Scrubber";
import { D3TimeSeriesGraph } from "./D3TimeSeriesGraph";

export const VizWorkspace = Vue.extend({
    components: {
        D3Scrubber,
        D3TimeSeriesGraph,
        Treeselect,
    },
    props: {
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    data() {
        return {
            selected: null,
        };
    },
    mounted() {
        return this.workspace.query();
    },
    methods: {
        uiNameOf(viz: Viz) {
            return "D3" + viz.constructor.name;
        },
        onVizTimeZoomed(viz, times) {
            viz.log("zooming", times);
            return this.workspace.zoomed(viz, times);
        },
        onVizRemove(viz) {
            viz.log("removing");
            return this.workspace.remove(viz);
        },
        onSelected(options) {
            console.log("selecting", options);
            return this.workspace.selected(options);
        },
    },
    template: `
		<div class="workspace">
			<div class="tree-container">
				<treeselect v-model="selected" :options="workspace.options" open-direction="bottom" @select="onSelected" />
			</div>
			<div class="groups-container">
				<div v-for="group in workspace.groups" :key="group.id">
					<component v-for="viz in group.vizes" :key="viz.id" v-bind:is="uiNameOf(viz)" :viz="viz" @viz-time-zoomed="(range) => onVizTimeZoomed(viz, range)" @viz-remove="() => onVizRemove(viz)"></component>
				</div>
			</div>
		</div>
	`,
});
