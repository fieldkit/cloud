import _ from "lodash";
import Vue from "vue";

import { TimeRange, Margins, ChartLayout } from "./common";
import { QueriedData, Scrubbers, TimeZoom } from "./viz";

import Scrubber from "./vega/Scrubber.vue";
import DoubleScrubber from "./vega/DoubleScrubber.vue";

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
    /*
    watch: {
        async visible(newValue, oldValue) {
            await this.refresh();
        },
        async data(newValue, oldValue) {
            await this.refresh();
        },
    },
    async mounted() {
        await this.refresh();
    },
    */
    methods: {
        raiseTimeZoomed(zoom: TimeZoom) {
            return this.$emit("viz-time-zoomed", zoom);
        },
        /*
        async refresh(): Promise<void> {
            return Promise.resolve();
        },
        */
    },
    template: `
        <div class="viz scrubber" v-if="data && data[0].data">
            <Scrubber :data="data[0]" :visible="visible" @time-zoomed="raiseTimeZoomed" v-if="data.length == 1" />
            <DoubleScrubber :data="data" :visible="visible" @time-zoomed="raiseTimeZoomed" v-if="data.length == 2" />
        </div>
    `,
});
