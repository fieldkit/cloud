import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { DataRow, SensorRange } from "./api";
import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom, VizInfo, SeriesData } from "./viz";

import LineChart from "./vega/LineChart.vue";
import DoubleLineChart from "./vega/DoubleLineChart.vue";

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
                const vizInfo = this.workspace.vizInfo(this.viz, ds);
                return new SeriesData(ds.graphing.key, ds, ds.graphing.data, vizInfo);
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
                <DoubleLineChart :series="allSeries" v-bind:key="allSeries[0].key + allSeries[1].key" @time-zoomed="raiseTimeZoomed" />
            </div>
            <div class="chart" @dblclick="onDouble" v-if="allSeries.length == 1">
                <LineChart :series="allSeries" v-bind:key="allSeries[0].key" @time-zoomed="raiseTimeZoomed" />
            </div>
        </div>
    `,
});
