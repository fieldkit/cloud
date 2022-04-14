import _ from "lodash";
import Vue, { PropType } from "vue";

import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

import { TimeRange, VizSensor } from "./common";
import { Graph, StationTreeOption, SensorTreeOption, Workspace, FastTime, TimeZoom, ChartType, DataSetSeries, NewParams } from "./viz";
import { vueTickHack } from "@/utilities";
import chartStyles from "./vega/chartStyles";
import { getPartnerCustomization } from "@/views/shared/partners";

interface VueDatepickerStyles {
    // this type might be extended with other customisations, it was found at https://github.com/nathanreyes/v-calendar/issues/531
    highlight: {
        start: {
            style: {
                backgroundColor: string;
            };
            contentStyle: {
                color: string;
            };
        };
        base: {
            style: {
                backgroundColor: string;
            };
        };
        end: {
            style: {
                backgroundColor: string;
            };
            contentStyle: {
                color: string;
            };
        };
    };
}

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
        ds: {
            type: Object as PropType<DataSetSeries>,
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
            if (this.disabled) {
                return null;
            }
            return this.ds.vizSensor[0]; // TODO VizSensor
        },
        selectedSensor(): string | null {
            if (this.disabled) {
                return null;
            }
            const sensorAndModule = this.ds.vizSensor[1]; // TODO VizSensor
            return `${sensorAndModule[0]}-${sensorAndModule[1]}`;
        },
        disabled(): boolean {
            if (this.stationOptions.length == 0 || this.sensorOptions.length == 0) {
                return true;
            }
            return this.viz.busy;
        },
    },
    methods: {
        raiseChangeStation(node: StationTreeOption): void {
            vueTickHack(() => {
                const newSeries = this.workspace.makeSeries(Number(node.id), this.ds.sensorAndModule);
                console.log("raising viz-change-series", newSeries);
                this.$emit("viz-change-series", newSeries);
            });
        },
        raiseChangeSensor(node: SensorTreeOption): void {
            vueTickHack(() => {
                if (!node.moduleId) throw new Error();
                if (!node.sensorId) throw new Error();
                const newSeries = this.workspace.makeSeries(this.ds.stationId, [node.moduleId, node.sensorId]);
                console.log("raising viz-change-series", newSeries);
                this.$emit("viz-change-series", newSeries);
            });
        },
    },
    template: `
		<div class="tree-pair">
            <treeselect :disabled="disabled" :value="selectedStation" :options="stationOptions" open-direction="bottom" @select="raiseChangeStation" :clearable="false" :searchable="false" />
            <treeselect :disabled="disabled" :value="selectedSensor" :options="sensorOptions" open-direction="bottom" @select="raiseChangeSensor" :default-expand-level="1" :clearable="false" :searchable="false" :disable-branch-nodes="true" />
		</div>
    `,
});

export const SelectionControls = Vue.extend({
    name: "SelectionControls",
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
            if (this.viz.busy) {
                return [];
            }
            // this.viz.log("station-options", { options: this.workspace.stationOptions });
            return this.workspace.stationOptions;
        },
        showRemove(): boolean {
            return this.viz.dataSets.length >= 2;
        },
        showAdd(): boolean {
            return this.viz.dataSets.length <= 1;
        },
    },
    methods: {
        sensorOptions(vizSensor: VizSensor): SensorTreeOption[] {
            if (this.stationOptions.length == 0) {
                return [];
            }
            const stationId = vizSensor[0]; // TODO VizSensor
            const sensorOptions = this.workspace.sensorOptions(stationId);
            // this.viz.log("sensor-options", { options: sensorOptions });
            return sensorOptions;
        },
        raiseChangeSeries(index: number, newSeries: DataSetSeries): void {
            const newParams = this.viz.modifySeries(index, [newSeries.stationId, newSeries.sensorAndModule]);
            this.viz.log("raise viz-change-sensors", index, newSeries);
            this.$emit("viz-change-sensors", newParams);
        },
        addSeries() {
            const newParams = this.viz.addSeries();
            this.viz.log("raise viz-change-sensors", newParams);
            this.$emit("viz-change-sensors", newParams);
        },
        removeSeries(index: number) {
            const newParams = this.viz.removeSeries(index);
            this.viz.log("raise viz-change-sensors", newParams);
            this.$emit("viz-change-sensors", newParams);
        },
        getKeyColor(idx) {
            const color = idx === 0 ? chartStyles.primaryLine.stroke : chartStyles.secondaryLine.stroke;
            return color;
        },
    },
    template: `
		<div class="left half">
            <div class="row" v-for="(ds, index) in viz.dataSets" v-bind:key="index">
                <div class="tree-key" :style="{color: getKeyColor(index)}">&#9632;</div>
                <SensorSelectionRow :viz="viz" :ds="ds" :workspace="workspace" :stationOptions="stationOptions" :sensorOptions="sensorOptions(ds.vizSensor)" @viz-change-series="(newSeries) => raiseChangeSeries(index, newSeries)"/>
                <div class="actions" v-if="showAdd || showRemove">
                    <div class="button" alt="Add" @click="() => addSeries()" v-if="showAdd">Add</div>
                    <div class="button" alt="Remove" @click="() => removeSeries(index)" v-if="showRemove">Remove</div>
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
        //*
        datepickerStyles(): VueDatepickerStyles {
            const primaryColor = getComputedStyle(document.body).getPropertyValue("--color-primary");
            const darkColor = getComputedStyle(document.body).getPropertyValue("--color-dark");

            const addAlpha = (color, opacity) => {
                const _opacity = Math.round(Math.min(Math.max(opacity || 1, 0), 1) * 255);
                return color + _opacity.toString(16).toUpperCase();
            };

            // for calendar range selection
            return {
                highlight: {
                    start: {
                        style: {
                            backgroundColor: primaryColor,
                        },
                        contentStyle: {
                            color: getPartnerCustomization() ? darkColor : "#ffffff",
                        },
                    },
                    base: {
                        style: {
                            backgroundColor: addAlpha(primaryColor, 0.3),
                        },
                    },
                    end: {
                        style: {
                            backgroundColor: primaryColor,
                        },
                        contentStyle: {
                            color: getPartnerCustomization() ? darkColor : "#ffffff",
                        },
                    },
                },
            };
        },
    },
    methods: {
        raiseCompare(): void {
            console.log("viz: raising viz-compare");
            this.$emit("viz-compare");
        },
        raiseRemove(): void {
            console.log("viz: raising viz-remove");
            this.$emit("viz-remove");
        },
        raiseFastTime(ev: never, fast: FastTime): void {
            console.log("viz: raising viz-time-zoomed");
            this.$emit("viz-time-zoomed", new TimeZoom(fast, null));
        },
        raiseChangeChartType(option: { id: ChartType }): void {
            console.log("viz: raising viz-change-chart", option.id);
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
                const rangeViz = this.viz.visibleTimeRange;
                const rangePicker = new TimeRange(fromPicker.start.getTime(), fromPicker.end.getTime());
                if (rangeViz.start != rangePicker.start || rangeViz.end != rangePicker.end) {
                    console.log("viz: raising viz-time-zoomed", rangeViz, rangePicker);
                    this.$emit("viz-time-zoomed", new TimeZoom(null, rangePicker));
                } else {
                    console.log("viz: swallowing viz-time-zoomed");
                }
            }
        },
    },
    template: `
		<div class="controls-container">
			<div class="row row-1">
				<div class="left">
				</div>
				<div class="right time">
					<span class="view-by">View By:</span>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 1)" v-bind:class="{ selected: viz.fastTime == 1 }">Day</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 7)" v-bind:class="{ selected: viz.fastTime == 7 }">Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 14)" v-bind:class="{ selected: viz.fastTime == 14 }">2 Week</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 30)" v-bind:class="{ selected: viz.fastTime == 30 }">Month</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 365)" v-bind:class="{ selected: viz.fastTime == 365 }">Year</div>
					<div class="fast-time" @click="ev => raiseFastTime(ev, 0)" v-bind:class="{ selected: viz.fastTime == 0 }">All</div>
					<div class="date-picker">
						<v-date-picker :value="manualRangeValue" @input="raiseManualTime" mode="date" :masks="{ input: 'MM/DD/YY' }" :select-attribute="datepickerStyles"
                                       :drag-attribute="datepickerStyles" is-range>
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

				<div class="right half" v-if="chartTypes.length > 1">
                    <div class="chart-type">
                        <treeselect :disabled="viz.busy" :options="chartTypes" :value="viz.chartType" open-direction="bottom" @select="raiseChangeChartType" :clearable="false" />
                    </div>
				</div>
			</div>
		</div>
	`,
});
