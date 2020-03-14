<template>
    <div :id="chart.id">
        <svg :view-box.camel="viewBox" preserveAspectRatio="xMidYMid meet" :width="outerWidth" :height="outerHeight" @dblclick="zoomOut">
            <g :style="stageStyle">
                <g ref="d3Stage"></g>
            </g>
        </svg>
        <D3LineChart :chart="chart" :layout="layout" ref="d3LineChart" @timeZoomed="onTimeZoom" />
        <D3HistoChart :chart="chart" :layout="layout" ref="d3HistoChart" />
        <D3RangeChart :chart="chart" :layout="layout" ref="d3RangeChart" />
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
        D3RangeChart,
    },
    props: ["chartParam", "station"],
    data: () => {
        return {
            chart: {},
            layout: {
                width: 1050,
                height: 350,
                marginTop: 5,
                marginRight: 0,
                marginBottom: 0,
                marginLeft: 50,
            },
        };
    },
    watch: {},
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
                transform: "translate(" + this.layout.marginLeft + "px," + this.layout.marginTop + "px)",
            };
        },
    },
    mounted() {
        this.chart = this.chartParam;
        this.chart.svg = d3.select(this.$refs.d3Stage);
        this.chart.colors = this.chartParam.sensor.colorScale;
        this.chart.panelID = this.station.device_id;
        this.activateChart();
    },
    methods: {
        activateChart() {
            switch (this.chart.type) {
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
        getTimeRange() {
            return { start: this.chart.start, end: this.chart.end };
        },
        setTimeRange(range) {
            this.chart.start = range.start;
            this.chart.end = range.end;
        },
        onTimeZoom(range) {
            this.$emit("timeZoomed", { range: range, parent: this.chart.parent, id: this.chart.id });
        },
        updateData(data, extent, colorScale) {
            this.chart.data = data;
            this.chart.extent = extent;
            this.chart.colors = colorScale;
            this.$refs.d3LineChart.dataChanged();
            this.$refs.d3HistoChart.dataChanged();
            this.$refs.d3RangeChart.dataChanged();
        },
        updateChartType() {
            this.$refs.d3LineChart.setStatus(false);
            this.$refs.d3HistoChart.setStatus(false);
            this.$refs.d3RangeChart.setStatus(false);
            // clear this svg and start fresh
            this.chart.svg.html(null);
            this.activateChart();
            switch (this.chart.type) {
                case "Line":
                    this.$refs.d3LineChart.init();
                    this.$refs.d3LineChart.makeLine();
                    break;
                case "Histogram":
                    this.$refs.d3HistoChart.makeHistogram();
                    break;
                case "Range":
                    this.$refs.d3RangeChart.makeRange();
                    break;
            }
        },
        zoomOut() {
            this.$emit("zoomOut", { parent: this.chart.parent, id: this.chart.id });
        },
    },
};
</script>

<style scoped></style>
