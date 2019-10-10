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
            @addBrush="addBrush"
            ref="d3LineChart"
        />
        <D3HistoChart
            :chart="chart"
            :processedData="processedData"
            :layout="layout"
            :selectedSensor="selectedSensor"
            ref="d3HistoChart"
        />
    </div>
</template>

<script>
import * as d3 from "d3";
import D3LineChart from "./D3LineChart";
import D3HistoChart from "./D3HistoChart";

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "D3Chart",
    components: {
        D3LineChart,
        D3HistoChart
    },
    props: ["stationData", "selectedSensor", "timeRange", "chartType"],
    data: () => {
        return {
            chart: {
                x: Object,
                y: Object,
                svg: Object,
                clip: Object,
                colors: Object,
                xAxis: Object,
                yAxis: Object,
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
                // drawing line chart by default
                this.$refs.d3LineChart.setStatus(true);
            }
        },
        timeRange: function() {
            this.timeChange();
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

            // x and y map functions
            this.chart.end = processed[processed.length - 1].date;
            this.chart.start = processed[0].date;
            this.chart.x = d3
                .scaleTime()
                .domain([this.chart.start, this.chart.end])
                .range([
                    this.layout.marginLeft,
                    this.layout.width - (this.layout.marginRight + this.layout.marginLeft)
                ]);
            this.chart.y = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([
                    this.layout.height - (this.layout.marginBottom + this.layout.marginTop),
                    this.layout.marginTop
                ]);
            this.chart.xAxis = d3.axisBottom(this.chart.x).ticks(10);
            this.chart.yAxis = d3.axisLeft(this.chart.y).ticks(6);

            return processed;
        },
        initSVG() {
            this.chart.svg = d3.select(this.$refs.d3Stage);

            this.chart.svg
                .append("defs")
                .append("svg:clipPath")
                .attr("id", "clip" + this.chart.panelID)
                .append("svg:rect")
                .attr("width", this.layout.width - this.layout.marginLeft * 2 - this.layout.marginRight)
                .attr("height", this.layout.height)
                .attr("x", this.layout.marginLeft)
                .attr("y", 0);

            this.$refs.d3LineChart.init();

            this.chart.colors = d3
                .scaleSequential()
                .domain(this.chart.extent)
                .interpolator(d3.interpolatePlasma);

            this.brush = d3
                .brushX()
                .extent([[0, 0], [this.layout.width, this.layout.height - this.layout.marginBottom]])
                .on("end", this.brushed);
        },
        addBrush() {
            // Add the brushing
            this.chart.svg
                .append("g")
                .attr("class", "brush")
                .attr("data-panel", this.chart.panelID)
                .call(this.brush);
        },
        brushed() {
            if (!d3.event.selection) {
                return;
            }

            let xRange = d3.event.selection;
            this.chart.start = this.chart.x.invert(xRange[0]);
            this.chart.end = this.chart.x.invert(xRange[1]);
            this.chart.x.domain([this.chart.start, this.chart.end]);
            // Remove the grey brush area after selection
            this.chart.svg.select(".brush").call(this.brush.move, null);
        },
        timeChange() {
            this.chart.start = this.processedData[0].date;
            this.chart.end = new Date(this.chart.start.getTime() + this.timeRange * DAY);
            if (this.timeRange == 0) {
                this.chart.end = this.processedData[this.processedData.length - 1].date;
            }
            this.chart.x.domain([this.chart.start, this.chart.end]);
        },
        chartTypeChange() {
            this.$refs.d3HistoChart.setStatus(false);
            this.$refs.d3LineChart.setStatus(false);
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
                    break;
            }
        }
    }
};
</script>

<style scoped></style>
