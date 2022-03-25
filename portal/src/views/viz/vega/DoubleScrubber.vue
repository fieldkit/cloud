<template>
    <div>
        <div class="viz scrubber"></div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { TimeRange } from "../common";
import { TimeZoom, SeriesData } from "../viz";
import { ScrubberSpecFactory } from "./ScrubberSpecFactory";

export default {
    name: "Scrubber",
    props: {
        series: {
            type: Array as PropType<SeriesData[]>,
            required: true,
        },
        visible: {
            type: Array,
            required: true,
        },
    },
    data(): { vega: unknown | null } {
        return { vega: null };
    },
    async mounted(): Promise<void> {
        await this.refresh();
    },
    watch: {
        async series(): Promise<void> {
            await this.refresh();
        },
        async visible() {
            console.log("vega-scrubber:visible", this.data, this.visible);
            this.pickRange(this.visible);
        },
    },
    methods: {
        async refresh(): Promise<void> {
            const factory = new ScrubberSpecFactory(this.series);

            const spec = factory.create();

            // Some styling overrides. The height of the scrubber can be set with scrubberSpec.height
            // scrubberSpec.config.axisX.tickSize = 20;
            // scrubberSpec.config.view = { fill: "#f4f5f7", stroke: "transparent" };
            // scrubberSpec.data = { values: this.data.data };
            // scrubberSpec.layer[2].data = { values: [] };
            // scrubberSpec.width = "container";

            const vegaInfo = await vegaEmbed(this.$el, spec, {
                renderer: "svg",
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;

            let scrubbed = [];
            vegaInfo.view.addSignalListener("brush", (_, value) => {
                if (value.time) {
                    scrubbed = value.time;
                } else {
                    scrubbed = this.data.timeRange;
                }
            });
            vegaInfo.view.addEventListener("mouseup", () => {
                console.log("vega-scrubber-brush", scrubbed);
                if (scrubbed.length == 2) {
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });

            this.pickRange(this.visible);
        },
        brush(times) {
            if (!this.vega) {
                return;
            }
            const x = times.map((v) => this.vega.view.scale("x")(v));
            this.vega.view
                .signal("brush_x", x)
                .signal("brush_tuple", {
                    unit: "layer_0",
                    fields: [
                        {
                            field: "time",
                            channel: "x",
                            type: "R",
                        },
                    ],
                    values: times,
                })
                .runAsync();
        },
        pickRange(timeRange) {
            /*
            if (_.isEqual(timeRange, this.data.timeRange)) {
                console.log("vega-scrubber:pick:ignore", timeRange, this.data.timeRange);
                this.brush([]);
            } else {
                console.log("vega-scrubber:pick", timeRange, this.data.timeRange);
                this.brush(timeRange);
            }
            */
        },
    },
};
</script>

<style lang="scss">
@import "src/scss/variables";
.viz {
    width: 100%;
}
</style>
