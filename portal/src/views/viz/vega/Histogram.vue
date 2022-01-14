<template>
    <div class="viz histogram"></div>
</template>

<script>
import { default as vegaEmbed } from "vega-embed";
import histoSpec from "./histo.vl.json";
import fieldkitBatteryData from "./fieldkitBatteryData.json";
import chartConfig from "./chartConfig.json";

histoSpec.config = chartConfig;

histoSpec.data = { values: fieldkitBatteryData.data };

export default {
    name: "Histogram",
    mounted: function() {
        vegaEmbed(".histogram", histoSpec, {
            renderer: "svg",
            tooltip: { offsetX: -50, offsetY: 50 },
            actions: { source: false, editor: false, compiled: false },
        });
    },
};
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>
