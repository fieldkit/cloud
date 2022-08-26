import _ from "lodash";
import Vue, { PropType } from "vue";

import { Workspace, TimeRange, TimeZoom, SeriesData } from "./viz";

import Scrubber from "./vega/Scrubber.vue";

export const VegaScrubber = Vue.extend({
    name: "VegaScrubber",
    components: {
        Scrubber,
    },
    props: {
        allSeries: {
            type: Array as PropType<SeriesData[]>,
            required: true,
        },
        visible: {
            type: Object as PropType<TimeRange>,
            required: true,
        },
        workspace: {
            type: Object as PropType<Workspace>,
            required: true,
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
