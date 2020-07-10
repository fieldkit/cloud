import _ from "lodash";
import Vue from "vue";

import { TimeRange } from "./common";
import { Workspace, Viz, Group, TimeZoom } from "./viz";
import { VizGraph } from "./VizGraph";
import { D3MultiScrubber } from "./D3MultiScrubber";

export const VizGroup = Vue.extend({
    components: {
        VizGraph,
        D3MultiScrubber,
    },
    props: {
        group: {
            type: Group,
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
    methods: {
        uiNameOf(viz: Viz) {
            return "Viz" + viz.constructor.name;
        },
        raiseGroupZoomed(range: TimeRange, ...args) {
            return this.$emit("group-time-zoomed", new TimeZoom(null, range), ...args);
        },
        raiseVizTimeZoomed(...args) {
            return this.$emit("viz-time-zoomed", ...args);
        },
        raiseRemove(...args) {
            return this.$emit("viz-remove", ...args);
        },
        raiseCompare(...args) {
            return this.$emit("viz-compare", ...args);
        },
        raiseChangeSensors(...args) {
            return this.$emit("viz-change-sensors", ...args);
        },
        raiseChangeChart(...args) {
            return this.$emit("viz-change-chart", ...args);
        },
    },
    template: `
		<div>
			<div v-if="group.scrubbers">
				<D3MultiScrubber :scrubbers="group.scrubbers" @viz-time-zoomed="(...args) => raiseGroupZoomed(...args)" />
			</div>
			<component v-for="viz in group.vizes" :key="viz.id" v-bind:is="uiNameOf(viz)" :viz="viz" :workspace="workspace"
				@viz-time-zoomed="(...args) => raiseVizTimeZoomed(viz, ...args)"
				@viz-remove="(...args) => raiseRemove(viz, ...args)"
				@viz-compare="(...args) => raiseCompare(viz, ...args)"
				@viz-change-sensors="(...args) => raiseChangeSensors(viz, ...args)"
				@viz-change-chart="(...args) => raiseChangeChart(viz, ...args)"
				/>
			</component>
		</div>
	`,
});
