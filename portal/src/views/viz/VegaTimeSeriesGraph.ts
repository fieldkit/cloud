import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { DataRow } from "./api";
import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

import LineChart from "./vega/LineChart.vue";
import DoubleLineChart from "./vega/DoubleLineChart.vue";

interface SeriesData {
    key: string;
    data: DataRow[];
}

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
            if (this.viz.dataSets.length) {
                const all = _.flatten(
                    this.viz.dataSets.map((ds) => {
                        if (ds && ds.graphing) {
                            return [
                                {
                                    key: ds.graphing.key,
                                    data: ds.graphing.data,
                                },
                            ];
                        }
                        return [];
                    })
                );

                if (all.length == this.viz.dataSets.length) {
                    return all;
                }
            }
            return [];
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
                <DoubleLineChart :data="allSeries" :labels="labels" :valueSuffixes="valueSuffixes" v-bind:key="allSeries[0].key + allSeries[1].key" @time-zoomed="raiseTimeZoomed" />
            </div>
            <div class="chart" @dblclick="onDouble" v-if="allSeries.length == 1">
                <LineChart :data="allSeries[0]" :label="labels[0]" :valueSuffix="valueSuffixes[0]" v-bind:key="allSeries[0].key" @time-zoomed="raiseTimeZoomed" />
            </div>
        </div>
    `,
});
