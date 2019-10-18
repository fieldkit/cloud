<template>
    <div>
        <svg
            :view-box.camel="viewBox"
            preserveAspectRatio="xMidYMid meet"
            :width="outerWidth"
            :height="outerHeight"
            @dblclick="zoomOut"
        >
            <g :style="stageStyle">
                <g ref="d3Stage"></g>
            </g>
        </svg>
        <D3LineChart
            :chart="chart"
            :stationData="stationData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3LineChart"
            @timeZoomed="onTimeZoom"
        />
        <D3HistoChart
            :chart="chart"
            :stationData="stationData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3HistoChart"
        />
        <D3RangeChart
            :chart="chart"
            :stationData="stationData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3RangeChart"
        />
    </div>
</template>

<script>
import * as d3 from "d3";
import D3LineChart from "./D3LineChart";
import D3HistoChart from "./D3HistoChart";
import D3RangeChart from "./D3RangeChart";

export default {
    name: "D3Chart",
    components: {
        D3LineChart,
        D3HistoChart,
        D3RangeChart
    },
    props: ["id", "station", "stationData", "selectedSensor", "chartType", "parent"],
    data: () => {
        return {
            chart: {
                svg: Object,
                extent: [],
                panelID: "",
                start: 0,
                end: 0
            },
            layout: {
                width: 1050,
                height: 350,
                marginTop: 5,
                marginRight: 0,
                marginBottom: 0,
                marginLeft: 50
            }
        };
    },
    watch: {
        chartType: function() {
            this.chartTypeChange();
        }
    },
    computed: {
        outerWidth: function() {
            return this.layout.width + this.layout.marginLeft + this.layout.marginRight;
        },
        outerHeight: function() {
            return this.layout.height + 50; //this.layout.marginTop + this.layout.marginBottom;
        },
        viewBox: function() {
            return "0 0 " + this.outerWidth + " " + this.outerHeight;
        },
        stageStyle: function() {
            return {
                transform: "translate(" + this.layout.marginLeft + "px," + this.layout.marginTop + "px)"
            };
        }
    },
    methods: {
        initChild(range) {
            this.chart.start = range.start;
            this.chart.end = range.end;
            this.initChart();
            this.initSVG();
            switch (this.chartType) {
                case "Line":
                    this.$refs.d3LineChart.makeLine();
                    break;
                case "Histogram":
                    this.$refs.d3HistoChart.makeHistogram();
                    break;
                case "Range":
                    this.$refs.d3RangeChart.makeRange();
                    break;
                default:
                    this.$refs.d3LineChart.makeLine();
                    break;
            }
        },
        initSVG() {
            this.chart.svg = d3.select(this.$refs.d3Stage);
            this.$refs.d3LineChart.init();
        },
        initChart() {
            let d3Chart = this;
            this.chart.panelID = this.station.device_id;
            this.chart.extent = d3.extent(this.stationData, d => {
                return d[d3Chart.selectedSensor.key];
            });
            switch (this.chartType) {
                case "Line":
                    this.$refs.d3LineChart.setStatus(true);
                    break;
                case "Histogram":
                    this.$refs.d3HistoChart.setStatus(true);
                    break;
                case "Range":
                    this.$refs.d3RangeChart.setStatus(true);
                    break;
                default:
                    this.$refs.d3LineChart.setStatus(true);
                    break;
            }
        },
        setTimeRange(range) {
            this.chart.start = range.start;
            this.chart.end = range.end;
            switch (this.chartType) {
                case "Line":
                    this.$refs.d3LineChart.timeChanged();
                    break;
                case "Histogram":
                    this.$refs.d3HistoChart.timeChanged();
                    break;
                case "Range":
                    this.$refs.d3RangeChart.timeChanged();
                    break;
            }
        },
        onTimeZoom(range) {
            this.$emit("timeZoomed", { range: range, parent: this.parent, id: this.id });
        },
        chartTypeChange() {
            this.$refs.d3LineChart.setStatus(false);
            this.$refs.d3HistoChart.setStatus(false);
            this.$refs.d3RangeChart.setStatus(false);
            // clear this svg and start fresh
            this.chart.svg.html(null);
            this.initSVG();
            switch (this.chartType) {
                case "Line":
                    this.$refs.d3LineChart.setStatus(true);
                    this.$refs.d3LineChart.makeLine();
                    break;
                case "Histogram":
                    this.$refs.d3HistoChart.setStatus(true);
                    this.$refs.d3HistoChart.makeHistogram();
                    break;
                case "Range":
                    this.$refs.d3RangeChart.setStatus(true);
                    this.$refs.d3RangeChart.makeRange();
                    break;
            }
        },
        zoomOut() {
            this.$emit("zoomOut", { parent: this.parent, id: this.id });
        }
    }
};
</script>

<style scoped></style>
