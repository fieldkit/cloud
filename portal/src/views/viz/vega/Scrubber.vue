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
import fieldkitBatteryData from "./fieldkitBatteryData.json";

export default {
    name: "Scrubber",
    props: {
        data: {
            type: Object,
            required: true,
        },
    },
    data() {
        return { vegaView: null };
    },
    mounted: function() {
        scrubberSpec.config = JSON.parse(JSON.stringify(chartConfig));
        // Some styling overrides. The height of the scrubber can be set with scrubberSpec.height
        scrubberSpec.config.axisX.tickSize = 20;
        scrubberSpec.config.view = { fill: "#f4f5f7", stroke: "transparent" };
        scrubberSpec.data = { values: this.data.data };
        scrubberSpec.layer[2].data = { values: [] };

        console.log("vega-scrubber", this.data);

        vegaEmbed(".scrubber", scrubberSpec, {
            renderer: "svg",
            actions: false, // { source: false, editor: false, compiled: false },
        }).then((result) => {
            this.vegaView = result;
            this.vegaView.view.addSignalListener("brush", function(_, value) {
                console.log("vega-scrubber-brush", value.time);
            });
        });
    },
    methods: {
        pickRange: function(timeRange) {
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
