<template>
    <div>
        <button v-on:click="downloadChart('png')" v-if="false">Download chart png</button>
        <button v-on:click="downloadChart('svg')" v-if="false">Download chart svg</button>
        <div class="viz linechart"></div>
    </div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";
import lineSpec from "./line.vl.json";
import fieldkitBatteryData from "./fieldkitBatteryData.json";
import chartConfig from "./chartConfig.json";

lineSpec.config = chartConfig;

lineSpec.data = { values: fieldkitBatteryData.data };

export default {
    name: "LineChart",
    mounted: function() {
        vegaEmbed(".linechart", lineSpec, {
            renderer: "svg",
            tooltip: { offsetX: -50, offsetY: 50 },
            actions: { source: false, editor: false, compiled: false },
        }).then((result) => {
            this.vegaView = result;
            result.view.addSignalListener("brush", function(_, value) {
                console.log(value.time);
            });
        });
    },
    methods: {
        // From https://vega.github.io/vega/docs/api/view/#view_toImageURL
        downloadChart: function(fileFormat) {
            this.vegaView.view
                .toImageURL(fileFormat, 2)
                .then(function(url) {
                    const link = document.createElement("a");
                    link.setAttribute("href", url);
                    link.setAttribute("target", "_blank");
                    link.setAttribute("download", "vega-export." + fileFormat);
                    link.dispatchEvent(new MouseEvent("click"));
                })
                .catch(function(error) {
                    console.log(error);
                });
        },
    },
};
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>
