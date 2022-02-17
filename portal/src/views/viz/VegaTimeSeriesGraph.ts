import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { DataRow, SensorRange } from "./api";
import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

import LineChart from "./vega/LineChart.vue";
import DoubleLineChart from "./vega/DoubleLineChart.vue";

interface SeriesData {
    key: string;
    data: DataRow[];
}

type AxisRange = SensorRange;

export const VegaTimeSeriesGraph = Vue.extend({
    name: "VegaTimeSeriesGraph",
    components: {
        LineChart,
        DoubleLineChart,
    },
    data() {
        return {};
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
        allSeries(): SeriesData[] | null {
            return this.viz.loadedDataSets.map((ds) => {
                if (!ds.graphing) throw new Error(`viz: No data`);
                return {
                    key: ds.graphing.key,
                    data: ds.graphing.data,
                };
            });
        },
        labels(): string[] {
            return this.viz.dataSets.map((ds) => {
                const vizInfo = this.workspace.vizInfo(this.viz, ds);
                if (vizInfo.unitOfMeasure) {
                    return i18n.tc(vizInfo.firmwareKey) + " (" + _.capitalize(vizInfo.unitOfMeasure) + ")";
                }
                return i18n.tc(vizInfo.firmwareKey);
            });
        },
        valueSuffixes(): string[] | null {
            return this.viz.dataSets.map((ds) => {
                const vizInfo = this.workspace.vizInfo(this.viz, ds);
                return vizInfo.unitOfMeasure, vizInfo.unitOfMeasure;
            });
        },
        constrainDataAxis(): AxisRange[] {
            return _.flatten(
                this.viz.dataSets.map((ds) => {
                    if (ds.constrainDataAxis) {
                        const vizInfo = this.workspace.vizInfo(this.viz, ds);
                        if (vizInfo.viz.length == 0) {
                            return [];
                        }
                        return vizInfo.constrainedRanges;
                    }
                    return [];
                })
            );
        },
        thresholds(): any[] {
            return this.viz.dataSets.map((ds) => {
                const vizInfo = this.workspace.vizInfo(this.viz, ds);
                if (vizInfo.viz.length == 0) {
                    return [];
                }

                const vizConfig = vizInfo.viz[0]; // TODO Pick for TimeSeries
                if (!vizConfig.thresholds) {
                    return [];
                }

                const EnglishLocale = "enUS";
                const thresholds = vizConfig.thresholds.levels;
                const thresholdLayers = thresholds
                    .map((d, i) => {
                        return {
                            transform: [
                                {
                                    calculate: "datum.value <= " + d.value + " ? datum.value : null",
                                    as: "layerValue" + i,
                                },
                                {
                                    calculate: "datum.layerValue" + i + " <= " + d.value + " ? '" + d.label[EnglishLocale] + "' : null",
                                    as: vizConfig.thresholds.label[EnglishLocale],
                                },
                            ],
                            encoding: {
                                y: { field: "layerValue" + i },
                                stroke: {
                                    field: vizConfig.thresholds.label[EnglishLocale],
                                    legend: {
                                        orient: "top",
                                    },
                                    scale: {
                                        domain: thresholds.map((d) => d.label[EnglishLocale]),
                                        range: thresholds.map((d) => d.color),
                                    },
                                },
                            },
                            mark: {
                                type: "line",
                                interpolate: "monotone",
                                tension: 1,
                            },
                        };
                    })
                    .reverse();

                return thresholdLayers;
            });
        },
    },
    methods: {
        onDouble() {
            return this.raiseTimeZoomed(new TimeZoom(FastTime.All, null));
        },
        raiseTimeZoomed(newTimes: TimeZoom) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
    },
    template: `
        <div class="viz time-series-graph">
            <div class="chart" @dblclick="onDouble" v-if="allSeries.length == 2">
                <DoubleLineChart :data="allSeries" :labels="labels" :valueSuffixes="valueSuffixes" :thresholds="thresholds" :constrainDataAxis="constrainDataAxis" v-bind:key="allSeries[0].key + allSeries[1].key" @time-zoomed="raiseTimeZoomed" />
            </div>
            <div class="chart" @dblclick="onDouble" v-if="allSeries.length == 1">
                <LineChart :data="allSeries[0]" :label="labels[0]" :valueSuffix="valueSuffixes[0]" :thresholds="thresholds[0]" :constrainDataAxis="constrainDataAxis" v-bind:key="allSeries[0].key" @time-zoomed="raiseTimeZoomed" />
            </div>
        </div>
    `,
});
