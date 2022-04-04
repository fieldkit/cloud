import _ from "lodash";
import Vue from "vue";

import { Graph, QueriedData, Workspace, FastTime, TimeZoom, SeriesData } from "./viz";

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
        allSeries(): SeriesData[] | null {
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
        <div class="viz">
            <div class="chart" @dblclick="onDouble">
                <RangeChart :series="allSeries" v-bind:key="key" />
            </div>
        </div>
    `,
});
