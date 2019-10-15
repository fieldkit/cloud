<template>
    <div>
        <svg
            :view-box.camel="viewBox"
            preserveAspectRatio="xMidYMid meet"
            :width="outerWidth"
            :height="outerHeight"
        >
            <g :style="stageStyle">
                <g ref="d3Stage"></g>
            </g>
        </svg>
        <D3LineChart
            :chart="chart"
            :processedData="processedData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3LineChart"
            @timeZoomed="onTimeZoom"
        />
        <D3HistoChart
            :chart="chart"
            :processedData="processedData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3HistoChart"
        />
        <D3RangeChart
            :chart="chart"
            :processedData="processedData"
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

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "D3Chart",
    components: {
        D3LineChart,
        D3HistoChart,
        D3RangeChart
    },
    props: ["stationData", "selectedSensor", "chartType"],
    data: () => {
        return {
            chart: {
                svg: Object,
                extent: [],
                panelID: "",
                start: 0,
                end: 0
            },
            processedData: [],
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
        stationData: function() {
            if (this.stationData.versions.length > 0) {
                this.processedData = this.processData();
                this.initSVG();
                this.initChart();
            } else {
                // TODO: handle lack of data more completely
                document.getElementById("loading").style.display = "none";
            }
        },
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
        processData() {
            let processed = [];
            this.stationData.versions.forEach(v => {
                let station = v.meta.station;
                this.chart.panelID = station.id;
                v.data.forEach(d => {
                    d.d.date = new Date(d.time * 1000);
                    processed.push(d.d);
                });
            });
            //sort data by date
            processed.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });

            let d3Chart = this;
            this.chart.extent = d3.extent(processed, d => {
                return d[d3Chart.selectedSensor.name];
            });

            if (this.$route.query.start && this.$route.query.end) {
                this.chart.start = new Date(parseInt(this.$route.query.start));
                this.chart.end = new Date(parseInt(this.$route.query.end));
            } else {
                this.chart.end = processed[processed.length - 1].date;
                this.chart.start = processed[0].date;
            }

            return processed;
        },
        initSVG() {
            this.chart.svg = d3.select(this.$refs.d3Stage);
            this.$refs.d3LineChart.init();
        },
        initChart() {
            const type = this.$route.query.type;
            if (type) {
                switch (type) {
                    case "Line":
                        this.$refs.d3LineChart.setStatus(true);
                        break;
                    case "Histogram":
                        this.$refs.d3HistoChart.setStatus(true);
                        break;
                    case "Range":
                        this.$refs.d3RangeChart.setStatus(true);
                        break;
                }
            } else {
                // draw line chart by default
                this.$refs.d3LineChart.setStatus(true);
            }
        },
        usePresetTimeRange(time) {
            this.chart.start = this.processedData[0].date;
            this.chart.end = new Date(this.chart.start.getTime() + time * DAY);
            if (time == 0) {
                this.chart.end = this.processedData[this.processedData.length - 1].date;
            }
            this.$emit("timeChanged", { start: this.chart.start, end: this.chart.end });
        },
        onTimeZoom(range) {
            this.$emit("timeChanged", range);
        },
        chartTypeChange() {
            this.$refs.d3LineChart.setStatus(false);
            this.$refs.d3HistoChart.setStatus(false);
            this.$refs.d3RangeChart.setStatus(false);
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
        refresh() {
            let newStart = this.processedData[0].date;
            let newEnd = this.processedData[this.processedData.length - 1].date;
            if (this.$route.query.start && this.$route.query.end) {
                newStart = new Date(parseInt(this.$route.query.start));
                newEnd = new Date(parseInt(this.$route.query.end));
            }
            if (newStart != this.chart.start || newEnd != this.chart.end) {
                this.chart.start = newStart;
                this.chart.end = newEnd;
            }
        }
    }
};
</script>

<style scoped></style>
