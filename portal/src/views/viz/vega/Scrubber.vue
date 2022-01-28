<template>
    <div>
        <button v-on:click="pickRange([1629956332434, 1630344275971])" v-if="false">
            Pick date range
        </button>
        <div class="viz scrubber"></div>
    </div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";
import scrubberSpec from "./scrubber.vl.json";
import chartConfig from "./chartConfig.json";

import { TimeRange } from "../common";
import { TimeZoom } from "../viz";

export default {
    name: "Scrubber",
    props: {
        data: {
            type: Object,
            required: true,
        },
        visible: {
            type: Array,
            required: true,
        },
    },
    data() {
        return { vegaView: null };
    },
    watch: {
        visible() {
            console.log("vega-scrubber:visible", this.data, this.visible);
            this.pickRange(this.visible);
        },
    },
    async mounted() {
        scrubberSpec.config = JSON.parse(JSON.stringify(chartConfig));
        // Some styling overrides. The height of the scrubber can be set with scrubberSpec.height
        scrubberSpec.config.axisX.tickSize = 20;
        scrubberSpec.config.view = { fill: "#f4f5f7", stroke: "transparent" };
        scrubberSpec.data = { values: this.data.data };
        scrubberSpec.layer[2].data = { values: [] };

        console.log("vega-scrubber", this.data, this.visible);

        await vegaEmbed(".scrubber", scrubberSpec, {
            renderer: "svg",
            actions: false, // { source: false, editor: false, compiled: false },
        }).then((view) => {
            this.vegaView = view;
            let scrubbed = [];
            this.vegaView.view.addSignalListener("brush", (_, value) => {
                if (value.time) {
                    scrubbed = value.time;
                } else {
                    scrubbed = this.data.timeRange;
                }
            });
            this.vegaView.view.addEventListener("mouseup", () => {
                console.log("vega-scrubber-brush", scrubbed);
                if (scrubbed.length == 2) {
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });
        });

        this.pickRange(this.visible);
    },
    methods: {
        pickRange(timeRange) {
            console.log("vega-scrubber:pick", timeRange, this.data.timeRange);
            if (_.isEqual(timeRange, this.data.timeRange)) {
                return;
            }
            this.vegaView.view
                .signal("brush_x", [this.vegaView.view.scale("x")(timeRange[0]), this.vegaView.view.scale("x")(timeRange[1])])
                .signal("brush_tuple", {
                    unit: "layer_0",
                    fields: [
                        {
                            field: "time",
                            channel: "x",
                            type: "R",
                        },
                    ],
                    values: [timeRange],
                })
                .runAsync();
        },
    },
};
</script>

<style lang="scss">
.viz {
    width: 100%;
}
</style>
