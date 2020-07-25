import _ from "lodash";
import Vue from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import Spinner from "@/views/shared/Spinner.vue";

import { TimeRange } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom, ChartType } from "./viz";

export const ViewingControls = Vue.extend({
    name: "ViewingControls",
    components: {
        Spinner,
        Treeselect,
    },
    data() {
        return {
            chartTypes: [
                {
                    label: "Time Series",
                    value: ChartType.TimeSeries,
                },
                {
                    label: "Histogram",
                    value: ChartType.Histogram,
                },
                {
                    label: "Range",
                    value: ChartType.Range,
                },
                {
                    label: "Map",
                    value: ChartType.Map,
                },
            ],
            range: null,
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
    computed: {
        selectedTreeOption() {
            return this.viz.chartParams.sensorParams.id;
        },
        manualRangeValue() {
            if (this.viz.visible.isExtreme()) {
                // TODO This happens initially cause we query for
                // eternity... probably best if this isn't set until
                // we get that data back and know the range.
                return null;
            }
            return {
                start: new Date(this.viz.visible.start),
                end: new Date(this.viz.visible.end),
            };
        },
    },
    methods: {
        raiseCompare(ev) {
            console.log("raising viz-compare");
            return this.$emit("viz-compare");
        },
        raiseRemove(ev) {
            console.log("raising viz-remove");
            return this.$emit("viz-remove");
        },
        raiseFastTime(ev, fast: FastTime, ...args) {
            console.log("raising viz-time-zoomed");
            return this.$emit("viz-time-zoomed", new TimeZoom(fast, null));
        },
        raiseChangeSensors(node, ...args) {
            console.log("raising viz-change-sensors");
            return this.$emit("viz-change-sensors", node);
        },
        raiseChangeChartType(chartType, ...args) {
            console.log("raising viz-change-chart");
            return this.$emit("viz-change-chart", Number(chartType));
        },
        raiseManualTime(fromPicker) {
            if (fromPicker) {
                // When the user picks a fast time this gets raised when
                // the viz changes the visible time, which we're bound to
                // so we do this to avoid raising a duplicate and querying
                // twice. I dunno if there's a better way.
                const rangeViz = this.viz.visible;
                const rangePicker = new TimeRange(fromPicker.start.getTime(), fromPicker.end.getTime());
                if (rangeViz.start != rangePicker.start || rangeViz.end != rangePicker.end) {
                    console.log("raising viz-time-zoomed");
                    return this.$emit("viz-time-zoomed", rangePicker);
                } else {
                    console.log("swallowing viz-time-zoomed");
                }
            }
        },
    },
    template: `
		<div class="controls-container">
			<div class="row row-1">
				<div class="left buttons">
					<div class="button" @click="raiseCompare">Compare</div>
					<div class="button" @click="raiseRemove" v-if="false">Remove</div>
					<div class="busy" v-if="viz.busy"><Spinner /></div>
				</div>
				<div class="right time">
					<span class="view-by">View By:</span>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 1)" v-bind:class="{ selected: viz.fastTime == 1 }">Day</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 7)" v-bind:class="{ selected: viz.fastTime == 7 }">Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 14)" v-bind:class="{ selected: viz.fastTime == 14 }">2 Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 30)" v-bind:class="{ selected: viz.fastTime == 30 }">Month</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 365)" v-bind:class="{ selected: viz.fastTime == 365 }">Year</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 0)" v-bind:class="{ selected: viz.fastTime == 0 }">All</div>
					<div class="time-picker">
						<v-date-picker class="vc-calendar" :value="manualRangeValue" @input="raiseManualTime" mode="range" :popover="{ placement: 'bottom left' }" :masks="{ input: 'MM/DD/YY' }" />
					</div>
				</div>
			</div>
			<div class="row row-2">
				<div class="left tree">
					<treeselect v-if="workspace.options.length" :value="selectedTreeOption" :options="workspace.options" open-direction="bottom" @select="raiseChangeSensors" />
					<div v-if="workspace.options.length == 0" class="loading-options">Loading Options</div>
				</div>

				<div class="right chart-type">
					<select v-on:change="(ev) => raiseChangeChartType(ev.target.value)" :value="viz.chartType">
						<option v-for="chartOption in chartTypes" v-bind:value="chartOption.value" v-bind:key="chartOption.value">
							{{ chartOption.label }}
						</option>
					</select>
				</div>
			</div>
		</div>
	`,
});
