import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

import LineChart from "./vega/LineChart.vue";

export const VegaTimeSeriesGraph = Vue.extend({
    name: "VegaTimeSeriesGraph",
    components: {
        LineChart,
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
        data(): QueriedData | null {
            if (this.viz.graphing) {
                return this.viz.graphing;
            }
            return null;
        },
        label(): string {
            const vizInfo = this.workspace.vizInfo(this.viz);
            if (vizInfo.unitOfMeasure) {
                return i18n.tc(vizInfo.firmwareKey) + " (" + _.capitalize(vizInfo.unitOfMeasure) + ")";
            }
            return i18n.tc(vizInfo.firmwareKey);
        },
        valueSuffix(): string | null {
            const vizInfo = this.workspace.vizInfo(this.viz);
            return vizInfo.unitOfMeasure;
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
            <div class="chart" @dblclick="onDouble" v-if="data">
                <LineChart :data="{ data: data.data }" :label="label" :valueSuffix="valueSuffix" v-bind:key="data.key" />
            </div>
        </div>
    `,
});
