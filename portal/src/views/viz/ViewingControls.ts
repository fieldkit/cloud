import _ from "lodash";
import Vue, { PropType } from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import Spinner from "@/views/shared/Spinner.vue";

import { TimeRange } from "./common";
import { Graph, StationTreeOption, SensorTreeOption, Workspace, FastTime, TimeZoom, ChartType } from "./viz";
import { vueTickHack } from "@/utilities";

export const SensorSelectionRow = Vue.extend({
    name: "SensorSelectionRow",
    components: {
        Treeselect,
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
        stationOptions: {
            type: Array as PropType<StationTreeOption[]>,
            required: true,
        },
        sensorOptions: {
            type: Array as PropType<StationTreeOption[]>,
            required: true,
        },
    },
    computed: {
        selectedStation(): number | null {
            return this.viz.chartParams.sensorParams.stations[0];
        },
        selectedSensor(): string | null {
            const sensorAndModule = this.viz.chartParams.sensorParams.sensor;
            return `${sensorAndModule[0]}-${sensorAndModule[1]}`;
        },
    },
    methods: {
        raiseChangeStation(node: StationTreeOption): void {
            const sensor = this.viz.chartParams.sensorParams.sensor;
            console.log("raising viz-change-sensors", "sensor", sensor);
            vueTickHack(() => {
                const params = this.workspace.makeParamsForStationChange(Number(node.id), sensor);
                this.$emit("viz-change-sensors", params);
            });
        },
        raiseChangeSensor(node: SensorTreeOption): void {
            const station = this.viz.chartParams.sensorParams.stations[0];
            console.log("raising viz-change-sensors", "station", station);
            vueTickHack(() => {
                if (!node.moduleId) throw new Error();
                if (!node.sensorId) throw new Error();
                const params = this.workspace.makeParamsForSensorChange(station, [node.moduleId, node.sensorId]);
                this.$emit("viz-change-sensors", params);
            });
        },
    },
    template: `
		<div class="tree-pair">
            <treeselect v-if="stationOptions.length" :value="selectedStation" :options="stationOptions" open-direction="bottom" @select="raiseChangeStation" :clearable="false" :searchable="false" />
            <div v-else class="loading-options">Loading Options</div>
            <treeselect v-if="sensorOptions.length" :value="selectedSensor" :options="sensorOptions" open-direction="bottom" @select="raiseChangeSensor" :default-expand-level="1" :clearable="false" :searchable="false" :disable-branch-nodes="true" />
            <div v-else class="loading-options">Loading Options</div>
		</div>
    `,
});

export const SelectionControls = Vue.extend({
    name: "SensorControls",
    components: {
        SensorSelectionRow,
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
            this.viz.log("station-options", { options: this.workspace.stationOptions });
            return this.workspace.stationOptions;
        },
        sensorOptions(): SensorTreeOption[] {
            const stationId = this.viz.chartParams.sensorParams.stations[0]; // this.selectedStation
            if (stationId == null) {
                return [];
            }
            const sensorOptions = this.workspace.sensorOptions(stationId);
            this.viz.log("sensor-options", { options: sensorOptions });
            return sensorOptions;
        },
    },
    methods: {
        raiseChangeSensors(...args: unknown[]): void {
            this.$emit("viz-change-sensors", ...args);
        },
    },
    template: `
		<div class="left half">
            <div class="row">
                <SensorSelectionRow :viz="viz" :workspace="workspace" :stationOptions="stationOptions" :sensorOptions="sensorOptions" @viz-change-sensors="raiseChangeSensors" />
                <div class="actions">
					<div class="button" alt="Add">Add</div>
                </div>
            </div>
        </div>
    `,
});

export const ViewingControls = Vue.extend({
    name: "ViewingControls",
    components: {
        SelectionControls,
        Treeselect,
        Spinner,
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
        compareIcon() {
            return this.$loadAsset("icon-compare.svg");
        },
        chartTypes(): { label: string; id: ChartType }[] {
            const vizInfo = this.workspace.vizInfo(this.viz);
            const allTypes = [
                {
                    label: "Time Series",
                    id: ChartType.TimeSeries,
                    vueName: "D3TimeSeriesGraph",
                },
                {
                    label: "Histogram",
                    id: ChartType.Histogram,
                    vueName: "D3Histogram",
                },
                {
                    label: "Range",
                    id: ChartType.Range,
                    vueName: "D3Range",
                },
                {
                    label: "Map",
                    id: ChartType.Map,
                    vueName: "D3Map",
                },
            ];
            if (vizInfo.viz.length == 0) {
                return allTypes;
            }
            const names = vizInfo.viz.map((vc) => vc.name);
            return allTypes.filter((type) => _.some(names, (name) => name == type.vueName));
        },
        manualRangeValue(): { start: Date; end: Date } | null {
            // console.log(`manual-range-value:`, this.viz.visible, this.viz.visibleTimeRange);
            if (!this.viz.visibleTimeRange || this.viz.visibleTimeRange.isExtreme()) {
                // TODO This happens initially cause we query for
                // eternity... probably best if this isn't set until
                // we get that data back and know the range.
                return null;
            }
            return {
                start: new Date(this.viz.visibleTimeRange.start),
                end: new Date(this.viz.visibleTimeRange.end),
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
        raiseChangeChartType(option: { id: ChartType }): void {
            console.log("raising viz-change-chart", option.id);
            vueTickHack(() => {
                this.$emit("viz-change-chart", Number(option.id));
            });
        },
        raiseChangeSensors(...args: unknown[]): void {
            this.$emit("viz-change-sensors", ...args);
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
				<div class="left buttons" v-if="!viz.busy">
					<div class="button compare" @click="raiseCompare" alt="Compare"> <img :src="compareIcon" /><div>Compare Sensor Graphs</div></div>
				</div>
				<div class="left busy" v-else><Spinner /></div>
				<div class="right time">
					<span class="view-by">View By:</span>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 1)" v-bind:class="{ selected: viz.fastTime == 1 }">Day</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 7)" v-bind:class="{ selected: viz.fastTime == 7 }">Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 14)" v-bind:class="{ selected: viz.fastTime == 14 }">2 Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 30)" v-bind:class="{ selected: viz.fastTime == 30 }">Month</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 365)" v-bind:class="{ selected: viz.fastTime == 365 }">Year</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 0)" v-bind:class="{ selected: viz.fastTime == 0 }">All</div>
					<div class="date-picker">
						<v-date-picker :value="manualRangeValue" @input="raiseManualTime" mode="date" :masks="{ input: 'MM/DD/YY' }" is-range>
                            <template v-slot="{ inputValue, inputEvents }">
                                <div class="flex justify-center items-center">
                                    <input
                                        :value="inputValue.start"
                                        v-on="inputEvents.start"
                                        class="border px-2 py-1 w-32 rounded focus:outline-none focus:border-indigo-300"
                                    />
                                    <input
                                        :value="inputValue.end"
                                        v-on="inputEvents.end"
                                        class="border px-2 py-1 w-32 rounded focus:outline-none focus:border-indigo-300"
                                    />
                                </div>
                            </template>
                        </v-date-picker>
					</div>
				</div>
			</div>
			<div class="row row-2">
                <SelectionControls :viz="viz" :workspace="workspace" @viz-change-sensors="raiseChangeSensors" />

				<div class="right chart-type" v-if="chartTypes.length > 1">
					<treeselect :options="chartTypes" :value="viz.chartType" open-direction="bottom" @select="raiseChangeChartType" :clearable="false" />
				</div>
			</div>
		</div>
	`,
});
