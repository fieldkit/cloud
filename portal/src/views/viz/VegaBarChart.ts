import _ from "lodash";
import Vue from "vue";

import { Graph, Workspace, FastTime, TimeZoom, VizInfo, SeriesData } from "./viz";

import BarChart from "./vega/BarChart.vue";

export const VegaBarChart = Vue.extend({
    name: "VegaBarChart",
    components: {
        BarChart,
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
            // TODO Version on viz.ts Graph?
            return this.viz.loadedDataSets.map((ds) => {
                if (!ds.graphing) throw new Error(`viz: No data`);
                const vizInfo = this.workspace.vizInfo(this.viz, ds);
                return new SeriesData(ds.graphing.key, ds, ds.graphing, vizInfo);
            });
        },
        key(): string {
            if (this.allSeries) {
                return this.allSeries.map((s) => s.key).join(":");
            }
            return "";
        },
    },
    methods: {
        onDouble(): void {
            this.raiseTimeZoomed(new TimeZoom(FastTime.All, null));
        },
        raiseTimeZoomed(newTimes: TimeZoom): void {
            this.$emit("viz-time-zoomed", newTimes);
        },
    },
    template: `
        <div class="viz bar-chart">
            <div class="chart" @dblclick="onDouble">
                <BarChart :series="allSeries" v-bind:key="key" @time-zoomed="raiseTimeZoomed" />
            </div>
        </div>
    `,
});
