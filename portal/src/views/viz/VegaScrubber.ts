import _ from "lodash";
import Vue from "vue";

import { Workspace, QueriedData, Scrubbers, TimeZoom, SeriesData, Graph } from "./viz";

import Scrubber from "./vega/Scrubber.vue";

export const VegaScrubber = Vue.extend({
    name: "VegaScrubber",
    components: {
        Scrubber,
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
            const all = _.flatten(
                this.scrubbers.rows.map((scrubber) => {
                    if (!scrubber.data) throw new Error(`viz: No data`);
                    const graph = scrubber.viz as Graph;
                    return graph.loadedDataSets.map((ds) => {
                        const vizInfo = this.workspace.vizInfo(graph, ds);
                        return new SeriesData(scrubber.data.key, ds, scrubber.data, vizInfo);
                    });
                })
            );
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
            <Scrubber :series="allSeries" :visible="visible" @time-zoomed="raiseTimeZoomed" v-if="allSeries.length > 0" />
        </div>
    `,
});
