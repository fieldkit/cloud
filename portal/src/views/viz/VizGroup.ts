import _ from "lodash";
import Vue from "vue";

import { TimeRange } from "./common";
import { Workspace, Viz, Group } from "./viz";
import { D3MultiScrubber } from "./D3MultiScrubber";
import { D3Graph } from "./D3Graph";

export const VizGroup = Vue.extend({
    components: {
        D3MultiScrubber,
        D3Graph,
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
            return "D3" + viz.constructor.name;
        },
        raiseGroupZoomed(newTimes: TimeRange) {
            return this.$emit("group-time-zoomed", newTimes);
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
        raiseFastTime(...args) {
            return this.$emit("viz-fast-time", ...args);
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
				@viz-fast-time="(...args) => raiseFastTime(viz, ...args)"
				@viz-change-sensors="(...args) => raiseChangeSensors(viz, ...args)"
				@viz-change-chart="(...args) => raiseChangeChart(viz, ...args)"
				/>
			</component>
		</div>
	`,
});
