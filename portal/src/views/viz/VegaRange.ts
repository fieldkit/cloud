import _ from "lodash";
import Vue from "vue";
import i18n from "@/i18n";

import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

import RangeChart from "./vega/RangeChart.vue";

export const VegaRange = Vue.extend({
    name: "VegaRange",
    components: {
        RangeChart,
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
    },
    watch: {
        data(newValue, oldValue) {
            this.viz.log("graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        this.viz.log("mounted");
        this.refresh();
    },
    updated() {
        this.viz.log("updated");
    },
    methods: {
        onDouble() {
            return this.raiseTimeZoomed(new TimeZoom(FastTime.All, null));
        },
        raiseTimeZoomed(newTimes: TimeZoom) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        refresh() {
            if (!this.data) {
                this.viz.log("refresh: nothing");
                return;
            } else {
                this.viz.log("refresh: data");
            }

            const vizInfo = this.workspace.vizInfo(this.viz);
            console.log("viz-info", vizInfo);
        },
    },
    template: `<div class="viz"><div class="chart" @dblclick="onDouble"><RangeChart /></div></div>`,
});
