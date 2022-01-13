<template>
  <div class="charts">
    <h1>FieldKit Vega charts</h1>
    <h2>Line chart</h2>
    <div class="chart">
      <LineChart />
    </div>
    <h2>Line chart with 2 Y-axes</h2>
    <div class="chart">
      <DoubleLineChart />
    </div>
    <h2>Range chart</h2>
    <div class="chart">
      <RangeChart />
    </div>
    <h2>Histogram</h2>
    <div class="chart">
      <Histogram />
    </div>
    <h2>Scrubber</h2>
    <p>
      The selected time range can be read (see console.log) and set (click the
      button).
    </p>
    <div class="chart">
      <Scrubber />
    </div>
  </div>
</template>

<script>
import LineChart from "./components/LineChart.vue";
import DoubleLineChart from "./components/DoubleLineChart.vue";
import Histogram from "./components/Histogram.vue";
import RangeChart from "./components/RangeChart.vue";
import Scrubber from "./components/Scrubber.vue";

export default {
  name: "App",
  components: {
    LineChart,
    Histogram,
    RangeChart,
    Scrubber,
    DoubleLineChart,
  },
};
</script>

<style>
@font-face {
  font-family: Avenir Light;
  src: url("/assets/fonts/Avenir/OpenType CFF Pro/AvenirLTPro-Light.otf");
}

.charts {
  width: 80%;
  margin: auto;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
}

@media screen and (max-width: 640px) {
  .charts {
    width: 100%;
  }
}

.chart {
  width: 100%;
}
</style>
