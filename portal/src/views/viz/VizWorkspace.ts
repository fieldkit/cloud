import _ from "lodash";
import Vue from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { TimeRange } from "./common";
import { Workspace, Viz, Querier, TreeOption } from "./viz";
import { D3Scrubber } from "./D3Scrubber";
import { D3Graph } from "./D3Graph";

export const VizWorkspace = Vue.extend({
    components: {
        D3Scrubber,
        D3Graph,
        Treeselect,
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
        uiNameOf(viz: Viz) {
            return "D3" + viz.constructor.name;
        },
        onVizTimeZoomed(viz: Viz, times: TimeRange) {
            viz.log("zooming", times.toArray());
            return this.workspace.zoomed(viz, times).query();
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
				<div v-for="group in workspace.groups" :key="group.id">
					<component v-for="viz in group.vizes" :key="viz.id" v-bind:is="uiNameOf(viz)" :viz="viz" :workspace="workspace"
						@viz-time-zoomed="(range) => onVizTimeZoomed(viz, range)"
						@viz-remove="() => onVizRemove(viz)"
						@viz-change="() => onVizChange(viz)">
					</component>
				</div>
			</div>
		</div>
	`,
});
