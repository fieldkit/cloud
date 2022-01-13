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
    data() {
        return {};
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
        visible(): TimeRange {
            return this.scrubbers.visible;
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
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", new TimeZoom(null, newTimes));
        },
        refresh() {
            console.log("scrubber refresh");
        },
    },
    template: `<div class="viz scrubber"><Scrubber /></div>`,
});
