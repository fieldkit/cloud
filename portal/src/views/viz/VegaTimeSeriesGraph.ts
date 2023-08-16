import _ from "lodash";
import Vue from "vue";

import { Graph, Workspace, FastTime, TimeZoom, VizInfo, SeriesData } from "./viz";

import LineChart from "./vega/LineChart.vue";
import InfoTooltip from '@/views/shared/InfoTooltip.vue';

export const VegaTimeSeriesGraph = Vue.extend({
    name: "VegaTimeSeriesGraph",
    components: {
        LineChart,
        InfoTooltip,
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
                return new SeriesData(ds.graphing.key, this.viz.visible, ds, ds.graphing, vizInfo);
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
        raiseTimeDragged(newTimes: TimeZoom): void {
            this.$emit("viz-time-dragged", newTimes);
        },
    },
    template: `
        <div class="viz time-series-graph">
            <InfoTooltip :message="$tc('dataView.swipeTip')"></InfoTooltip>
            <div class="chart" @dblclick="onDouble">
                <LineChart :series="allSeries" v-bind:key="key" @time-zoomed="raiseTimeZoomed" @time-dragged="raiseTimeDragged" />
            </div>
        </div>
    `,
});
