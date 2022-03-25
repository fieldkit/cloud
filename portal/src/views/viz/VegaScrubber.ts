import _ from "lodash";
import Vue from "vue";

import { Workspace, QueriedData, Scrubbers, TimeZoom, SeriesData, Graph } from "./viz";

import Scrubber from "./vega/Scrubber.vue";
import DoubleScrubber from "./vega/DoubleScrubber.vue";

export const VegaScrubber = Vue.extend({
    name: "VegaScrubber",
    components: {
        Scrubber,
        DoubleScrubber,
    },
    props: {
        scrubbers: {
            type: Scrubbers,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    computed: {
        allSeries(): SeriesData[] | null {
            const all = this.scrubbers.rows.map((scrubber) => {
                if (!scrubber.data) throw new Error(`viz: No data`);
                const graph = scrubber.viz as Graph;
                const graphDataSets = graph.loadedDataSets;
                const vizInfo = this.workspace.vizInfo(graph, graphDataSets[0]);
                return new SeriesData(scrubber.data.key, graphDataSets[0], scrubber.data, vizInfo);
            });
            console.log("viz: scrubber:all", all);
            return all;
        },
        visible(): number[] {
            if (this.scrubbers.visible.isExtreme()) {
                return this.scrubbers.timeRange.toArray();
            }
            return this.scrubbers.visible.toArray();
        },
    },
    methods: {
        raiseTimeZoomed(zoom: TimeZoom): void {
            this.$emit("viz-time-zoomed", zoom);
        },
    },
    template: `
        <div class="viz scrubber">
            <Scrubber :series="allSeries" :visible="visible" @time-zoomed="raiseTimeZoomed" v-if="allSeries.length == 1" />
            <DoubleScrubber :series="allSeries" :visible="visible" @time-zoomed="raiseTimeZoomed" v-if="allSeries.length == 2" />
        </div>
    `,
});
