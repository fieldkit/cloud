import _ from "lodash";
import Vue from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { TimeRange } from "./common";
import { Graph, QueriedData, Workspace, FastTime, ChartType } from "./viz";

export const DebuggingPanel = Vue.extend({
    name: "DebuggingPanel",
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
    data() {
        return {};
    },
    template: `
		<div style="debug-panel" style="padding: 20px;">
			<div> debug {{ viz.id }} stations: {{ viz.chartParams.stations }} sensors: {{ viz.chartParams.sensors }} </div>
			<div> visible({{ viz.visible.start | prettyTime }} - {{ viz.visible.end | prettyTime }}) </div>
			<div> query({{ viz.chartParams.when.start | prettyTime }} - {{ viz.chartParams.when.end | prettyTime }}) </div>
		</div>
	`,
});
