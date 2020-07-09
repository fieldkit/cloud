import _ from "lodash";
import Vue from "vue";

import { TimeRange } from "./common";
import { Workspace, Group, Viz, TreeOption, FastTime, ChartType } from "./viz";
import { VizGroup } from "./VizGroup";

export const VizWorkspace = Vue.extend({
    components: {
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
        onGraphTimeZoomed(group: Group, viz: Viz, times: TimeRange) {
            viz.log("zooming", times.toArray());
            return this.workspace.graphZoomed(viz, times).query();
        },
        onRemove(group: Group, viz: Viz) {
            viz.log("removing");
            return this.workspace.remove(viz).query();
        },
        onCompare(group: Group, viz: Viz) {
            viz.log("compare");
            return this.workspace.compare(viz).query();
        },
        onFastTime(group: Group, viz: Viz, fastTime: FastTime) {
            viz.log("fast-time", fastTime);
            return this.workspace.fastTime(viz, fastTime).query();
        },
        onChangeSensors(group: Group, viz: Viz, option: TreeOption) {
            viz.log("sensors", option);
            return this.workspace.changeSensors(viz, option).query();
        },
        onChangeChart(group: Group, viz: Viz, chartType: ChartType) {
            viz.log("chart", chartType);
            return this.workspace.changeChart(viz, chartType).query();
        },
    },
    template: `
		<div class="workspace-container">
			<div class="groups-container">
				<VizGroup v-for="group in workspace.groups" :key="group.id" :group="group" :workspace="workspace"
					@group-time-zoomed="(...args) => onGroupTimeZoomed(group, ...args)"
					@viz-time-zoomed="(...args) => onGraphTimeZoomed(group, ...args)"
					@viz-remove="(...args) => onRemove(group, ...args)"
					@viz-compare="(...args) => onCompare(group, ...args)"
					@viz-fast-time="(...args) => onFastTime(group, ...args)"
					@viz-change-sensors="(...args) => onChangeSensors(group, ...args)"
					@viz-change-chart="(...args) => onChangeChart(group, ...args)"
				/>
			</div>
		</div>
	`,
});
