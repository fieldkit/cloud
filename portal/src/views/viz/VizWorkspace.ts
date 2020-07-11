import _ from "lodash";
import Vue from "vue";

import { TimeRange } from "./common";
import { Workspace, Group, Viz, TreeOption, FastTime, ChartType, TimeZoom, Bookmark } from "./viz";
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
        signal(ws: Workspace): Workspace {
            const bookmark = ws.bookmark();
            this.$emit("change", bookmark);
            console.log("SIGNAL", bookmark);
            return ws;
        },
        onGroupTimeZoomed(group: Group, zoom: TimeZoom) {
            group.log("zooming", zoom);
            return this.workspace
                .groupZoomed(group, zoom)
                .with((this as any).signal)
                .query();
        },
        onGraphTimeZoomed(group: Group, viz: Viz, zoom: TimeZoom) {
            viz.log("zooming", zoom);
            return this.workspace
                .graphZoomed(viz, zoom)
                .with((this as any).signal)
                .query();
        },
        onRemove(group: Group, viz: Viz) {
            viz.log("removing");
            return this.workspace
                .remove(viz)
                .with((this as any).signal)
                .query();
        },
        onCompare(group: Group, viz: Viz) {
            viz.log("compare");
            return this.workspace
                .compare(viz)
                .with((this as any).signal)
                .query();
        },
        onChangeSensors(group: Group, viz: Viz, option: TreeOption) {
            viz.log("sensors", option);
            return this.workspace
                .changeSensors(viz, option)
                .with((this as any).signal)
                .query();
        },
        onChangeChart(group: Group, viz: Viz, chartType: ChartType) {
            viz.log("chart", chartType);
            return this.workspace
                .changeChart(viz, chartType)
                .with((this as any).signal)
                .query();
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
					@viz-change-sensors="(...args) => onChangeSensors(group, ...args)"
					@viz-change-chart="(...args) => onChangeChart(group, ...args)"
				/>
			</div>
		</div>
	`,
});
