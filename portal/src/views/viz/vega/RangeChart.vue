<template>
  <div class="viz rangechart"></div>
</template>

<script>
import { default as vegaEmbed } from "vega-embed";
import rangeSpec from "../assets/range.vl.json";
import fieldkitBatteryData from "../assets/fieldkitBatteryData.json";
import chartConfig from "../assets/chartConfig.json";

rangeSpec.config = chartConfig;
rangeSpec.data = { values: fieldkitBatteryData.data, name: "rangedata" };

export default {
  name: "RangeChart",
  mounted: function () {
    vegaEmbed(".rangechart", rangeSpec, {
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
