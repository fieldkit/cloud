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
    // This type might be extended with other customizations, it was found at https://github.com/nathanreyes/v-calendar/issues/531
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
            const originalStationId = this.ds.vizSensor[0]; // TODO VizSensor
            return originalStationId;
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
        searchable(): boolean {
            return !(window.screen.availWidth <= 768);
        },
    },
    methods: {
        raiseChangeStation(node: StationTreeOption): void {
            vueTickHack(() => {
                const newSeries = this.workspace.makeSeries(Number(node.id), null);
                console.log("raising viz-change-series", newSeries);
                this.$emit("viz-change-series", newSeries);
            });
        },
        raiseChangeSensor(node: SensorTreeOption): void {
            vueTickHack(() => {
                if (!node.moduleId) throw new Error();
                if (!node.sensorId) throw new Error();
                const newSeries = this.workspace.makeSeries(node.stationId || this.ds.stationId, [node.moduleId, node.sensorId]);
                console.log("raising viz-change-series", newSeries);
                this.$emit("viz-change-series", newSeries);
            });
        },
    },
    template: `
		<div class="tree-pair">
            <treeselect :disabled="disabled" :value="selectedStation" :options="stationOptions" open-direction="bottom" @select="raiseChangeStation" :clearable="false" :searchable="searchable" :disable-branch-nodes="true" />
            <treeselect :disabled="disabled" :value="selectedSensor" :options="sensorOptions" open-direction="bottom" @select="raiseChangeSensor" :default-expand-level="3" :clearable="false" :searchable="false" :disable-branch-nodes="true" />
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
            const originalStationId = vizSensor[0]; // TODO VizSensor
            const stationId = originalStationId;
            const maybePartner = getPartnerCustomization();
            const flatten = maybePartner != null && maybePartner.projectId != null && this.workspace.projects[0] === maybePartner.projectId;
            const sensorOptions = this.workspace.sensorOptions(stationId, flatten);
            /*
                this.viz.log("viz-sensor-options", {
                    vizSensor: vizSensor,
                    originalStationId: originalStationId,
                    stationId: stationId,
                    options: sensorOptions,
                });
            */
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
            const allTypes = [
                {
                    label: "Time Series",
                    id: ChartType.TimeSeries,
                },
                {
                    label: "Bar",
                    id: ChartType.Bar,
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
            ];
            const availableTypes = this.workspace.availableChartTypes(this.viz);
            return allTypes.filter((type) => _.some(availableTypes, (row) => row == type.id));
        },
        selectedChartType(): ChartType {
            return this.viz.chartType;
        },
        manualRangeValue(): { start: Date; end: Date } | null {
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

            // For calendar range selection
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
        raiseManualTime(fromPicker, pickerType): void {
            if (fromPicker) {
                // When the user picks a fast time this gets raised when
                // the viz changes the visible time, which we're bound to
                // so we do this to avoid raising a duplicate and querying
                // twice. I dunno if there's a better way.
                if (!this.manualRangeValue) {
                    return;
                }

                const rangeViz = this.viz.visibleTimeRange;

                let pickerStart: Date = new Date();
                let pickerEnd: Date = new Date();

                if (pickerType === "start") {
                    if (fromPicker.getTime() === rangeViz.start) {
                        console.log("viz: swallow(same-start)");
                        return;
                    }

                    pickerStart = new Date(fromPicker.getFullYear(), fromPicker.getMonth(), fromPicker.getDate(), 0, 0, 0, 0);
                    pickerEnd = this.manualRangeValue.end;
                }

                if (pickerType === "end") {
                    if (fromPicker.getTime() == rangeViz.end) {
                        console.log("viz: swallow(same-end)");
                        return;
                    }

                    pickerStart = this.manualRangeValue.start;
                    pickerEnd = new Date(fromPicker.getFullYear(), fromPicker.getMonth(), fromPicker.getDate(), 23, 59, 59, 999);
                }

                const rangePicker = new TimeRange(pickerStart.getTime(), pickerEnd.getTime());

                if (rangeViz.start != rangePicker.start || rangeViz.end != rangePicker.end) {
                    console.log("viz: raising viz-time-zoomed", "viz", rangeViz.describe(), "picked", rangePicker.describe());
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
                    <div class="fast-time-container">
                        <span class="view-by">View By:</span>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 1)" v-bind:class="{ selected: viz.fastTime == 1 }">Day</div>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 7)" v-bind:class="{ selected: viz.fastTime == 7 }">Week</div>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 14)" v-bind:class="{ selected: viz.fastTime == 14 }">2 Week</div>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 30)" v-bind:class="{ selected: viz.fastTime == 30 }">Month</div>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 365)" v-bind:class="{ selected: viz.fastTime == 365 }">Year</div>
                        <div class="fast-time" @click="ev => raiseFastTime(ev, 0)" v-bind:class="{ selected: viz.fastTime == 0 }">All</div>
                    </div>
					<div class="date-picker flex" v-if="manualRangeValue">
						<v-date-picker :value="manualRangeValue.start" @input="raiseManualTime($event, 'start')" mode="date" :masks="{ input: 'MM/DD/YY' }" :select-attribute="datepickerStyles"
                                       :drag-attribute="datepickerStyles" :maxDate="manualRangeValue.end" is-inline
                        >
                            <template v-slot="{ inputValue, inputEvents }">
                                <div class="flex justify-center items-center">
                                    <input
                                        :value="inputValue"
                                        v-on="inputEvents"
                                        placeholder="Start date"
                                        class="border px-2 py-1 w-32 rounded focus:outline-none focus:border-indigo-300"
                                    />
                                </div>
                            </template>
                        </v-date-picker>
                        <v-date-picker :value="manualRangeValue.end" @input="raiseManualTime($event, 'end')" mode="date" :masks="{ input: 'MM/DD/YY' }" :select-attribute="datepickerStyles"
                                     :drag-attribute="datepickerStyles" :minDate="manualRangeValue.start">
                        <template v-slot="{ inputValue, inputEvents }">
                          <div class="flex justify-center items-center">
                            <input
                                :value="inputValue"
                                v-on="inputEvents"
                                placeholder="End date"
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
                        <treeselect :disabled="viz.busy" :options="chartTypes" :value="selectedChartType" open-direction="bottom" @select="raiseChangeChartType" :clearable="false" />
                    </div>
				</div>
			</div>
		</div>
	`,
});
