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
			<div> debug {{ viz.id }} stations: {{ viz.params.stations }} sensors: {{ viz.params.sensors }} </div>
			<div> visible({{ viz.visible.start | prettyTime }} - {{ viz.visible.end | prettyTime }}) </div>
			<div> query({{ viz.params.when.start | prettyTime }} - {{ viz.params.when.end | prettyTime }}) </div>
		</div>
	`,
});
