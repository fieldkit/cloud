<template>
    <div>
        <button v-on:click="pickRange([1629956332434, 1630344275971])">
            Pick date range
        </button>
        <div class="viz scrubber"></div>
    </div>
</template>

<script>
import { default as vegaEmbed } from "vega-embed";
import scrubberSpec from "./scrubber.vl.json";
import fieldkitBatteryData from "./fieldkitBatteryData.json";
import iconsData from "./iconsData.json";
import chartConfig from "./chartConfig.json";

scrubberSpec.config = JSON.parse(JSON.stringify(chartConfig));
// Some styling overrides. The height of the scrubber can be set with scrubberSpec.height
scrubberSpec.config.axisX.tickSize = 20;
scrubberSpec.config.view = { fill: "#f4f5f7", stroke: "transparent" };

scrubberSpec.data = { values: fieldkitBatteryData.data };
scrubberSpec.layer[2].data = { values: iconsData.data };

export default {
    name: "Scrubber",
    data() {
        return { vegaView: null };
    },
    mounted: function() {
        vegaEmbed(".scrubber", scrubberSpec, {
            renderer: "svg",
            actions: { source: false, editor: false, compiled: false },
        }).then((result) => {
            this.vegaView = result;
            this.vegaView.view.addSignalListener("brush", function(_, value) {
                console.log(value.time);
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

<style scoped>
.viz {
    width: 100%;
}
</style>
