import _ from "lodash";
import Vue, { PropType } from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import Spinner from "@/views/shared/Spinner.vue";

import { TimeRange } from "./common";
import { Graph, HasSensorParams, SensorParams, StationTreeOption, SensorTreeOption, Workspace, FastTime, TimeZoom, ChartType } from "./viz";
import { vueTickHack } from "@/utilities";

class NewParams implements HasSensorParams {
    public readonly sensorParams: SensorParams;

    constructor(stationId: number, sensorAndModule: [number, number]) {
        this.sensorParams = new SensorParams([stationId], [sensorAndModule]);
    }
}

export const ViewingControls = Vue.extend({
    name: "ViewingControls",
    components: {
        Treeselect,
        Spinner,
    },
    data(): {
        chartTypes: { label: string; id: ChartType }[];
    } {
        return {
            chartTypes: [
                {
                    label: "Time Series",
                    id: ChartType.TimeSeries,
                },
                {
                    label: "Histogram",
                    id: ChartType.Histogram,
                },
                {
                    label: "Range",
                    id: ChartType.Range,
                },
                {
                    label: "Map",
                    id: ChartType.Map,
                },
            ],
        };
    },
    props: {
        viz: {
            type: Object as PropType<Graph>,
            required: true,
        },
        workspace: {
            type: Object as PropType<Workspace>,
            required: true,
        },
    },
    computed: {
        stationOptions(): StationTreeOption[] {
            this.viz.log("station-options", this.workspace.stationOptions);
            return this.workspace.stationOptions;
        },
        sensorOptions(): SensorTreeOption[] {
            this.viz.log("sensor-options", this.workspace.sensorOptions);
            return this.workspace.sensorOptions;
        },
        selectedStation(): number | null {
            return this.viz.chartParams.sensorParams.stations[0];
        },
        selectedSensor(): string | null {
            const sensorAndModule = this.viz.chartParams.sensorParams.sensors[0];
            return `${sensorAndModule[0]}-${sensorAndModule[1]}`;
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
        raiseCompare(): void {
            console.log("raising viz-compare");
            this.$emit("viz-compare");
        },
        raiseRemove(): void {
            console.log("raising viz-remove");
            this.$emit("viz-remove");
        },
        raiseFastTime(ev: never, fast: FastTime): void {
            console.log("raising viz-time-zoomed");
            this.$emit("viz-time-zoomed", new TimeZoom(fast, null));
        },
        raiseChangeStation(node: StationTreeOption): void {
            const sensor = this.viz.chartParams.sensorParams.sensors[0];
            console.log("raising viz-change-sensors", "sensor", sensor);
            vueTickHack(() => {
                this.$emit("viz-change-sensors", new NewParams(Number(node.id), sensor));
            });
        },
        raiseChangeSensor(node: SensorTreeOption): void {
            const station = this.viz.chartParams.sensorParams.stations[0];
            console.log("raising viz-change-sensors", "station", station);
            vueTickHack(() => {
                if (!node.moduleId) throw new Error();
                if (!node.sensorId) throw new Error();
                this.$emit("viz-change-sensors", new NewParams(station, [node.moduleId, node.sensorId]));
            });
        },
        raiseChangeChartType(option: { id: ChartType }): void {
            console.log("raising viz-change-chart", option.id);
            vueTickHack(() => {
                this.$emit("viz-change-chart", Number(option.id));
            });
        },
        raiseManualTime(fromPicker): void {
            if (fromPicker) {
                // When the user picks a fast time this gets raised when
                // the viz changes the visible time, which we're bound to
                // so we do this to avoid raising a duplicate and querying
                // twice. I dunno if there's a better way.
                const rangeViz = this.viz.visible;
                const rangePicker = new TimeRange(fromPicker.start.getTime(), fromPicker.end.getTime());
                if (rangeViz.start != rangePicker.start || rangeViz.end != rangePicker.end) {
                    console.log("raising viz-time-zoomed");
                    this.$emit("viz-time-zoomed", new TimeZoom(null, rangePicker));
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
					<treeselect v-if="stationOptions.length" :value="selectedStation" :options="stationOptions" open-direction="bottom" @select="raiseChangeStation" :clearable="false" :searchable="false" />
					<treeselect v-if="sensorOptions.length" :value="selectedSensor" :options="sensorOptions" open-direction="bottom" @select="raiseChangeSensor" :default-expand-level="1" :clearable="false" :searchable="false" :disable-branch-nodes="true" />
					<div v-if="stationOptions.length == 0 || sensorOptions.length == 0" class="loading-options">Loading Options</div>
				</div>

				<div class="right chart-type">
					<treeselect v-if="stationOptions.length" :options="chartTypes" :value="viz.chartType" open-direction="bottom" @select="raiseChangeChartType" :clearable="false" />
				</div>
			</div>
		</div>
	`,
});
