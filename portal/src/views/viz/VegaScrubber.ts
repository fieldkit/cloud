import _ from "lodash";
import Vue, { PropType } from "vue";

import { TimeRange, TimeZoom, SeriesData } from "./viz";

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
        dragging: {
            type: Boolean,
            required: true,
        },
    },
    methods: {
        raiseTimeZoomed(zoom: TimeZoom): void {
            this.$emit("viz-time-zoomed", zoom);
        },
        eventClicked(id: number): void {
            this.$emit("event-clicked", id);
        },
    },
    template: `
        <div class="viz scrubber">
            <Scrubber :series="allSeries" :visible="visible" :dragging="dragging" @time-zoomed="raiseTimeZoomed" @event-clicked="eventClicked" v-if="allSeries.length > 0" />
        </div>
    `,
});
