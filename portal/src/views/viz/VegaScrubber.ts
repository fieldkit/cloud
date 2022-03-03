import _ from "lodash";
import Vue from "vue";

import { TimeRange, Margins, ChartLayout } from "./common";
import { QueriedData, Scrubbers, TimeZoom } from "./viz";

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
    },
    computed: {
        data(): QueriedData[] {
            return this.scrubbers.rows.map((s) => s.data);
        },
        visible(): number[] {
            if (this.scrubbers.visible.isExtreme()) {
                return this.scrubbers.timeRange.toArray();
            }
            return this.scrubbers.visible.toArray();
        },
    },
    watch: {
        visible(newValue, oldValue) {
            // console.log("scrubber graphing (visible)");
            this.refresh();
        },
        data(newValue, oldValue) {
            // console.log("scrubber graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        // console.log("scrubber mounted");
        this.refresh();
    },
    updated() {
        // console.log("scrubber updated");
    },
    methods: {
        raiseTimeZoomed(zoom: TimeZoom) {
            return this.$emit("viz-time-zoomed", zoom);
        },
        refresh() {
            // console.log("scrubber refresh", this.scrubbers);
        },
    },
    template: `
        <div class="viz scrubber" v-if="data && data[0].data">
            <Scrubber :data="data[0]" :visible="visible" @time-zoomed="raiseTimeZoomed" />
        </div>
    `,
});
