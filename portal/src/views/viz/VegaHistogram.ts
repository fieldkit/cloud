import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

import HistogramChart from "./vega/Histogram.vue";

export const VegaHistogram = Vue.extend({
    name: "VegaHistogram",
    components: {
        HistogramChart,
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
        <div class="viz">
            <div class="chart" @dblclick="onDouble" v-if="data">
                <HistogramChart :data="{ data: data.data }" :label="label" v-bind:key="data.key" />
            </div>
        </div>
    `,
});
