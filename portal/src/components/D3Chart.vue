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
        <D3MapChart :chart="chart" :layout="layout" ref="d3MapChart" />
        <div :id="'scrubber-container-' + chart.id">
            <svg
                :id="'scrubber-svg-' + chart.id"
                :view-box.camel="scrubberViewBox"
                preserveAspectRatio="xMidYMid meet"
                :width="outerWidth"
                :height="scrubberHeight"
            >
                <filter :id="'dropshadow-' + chart.id" height="150%" width="150%">
                    <feGaussianBlur in="SourceAlpha" stdDeviation="2" />
                    <!-- stdDeviation is how much to blur -->
                    <feOffset dx="2" dy="2" result="offsetblur" />
                    <!-- how much to offset -->
                    <feComponentTransfer>
                        <feFuncA type="linear" slope="0.5" />
                        <!-- slope is the opacity of the shadow -->
                    </feComponentTransfer>
                    <feMerge>
                        <feMergeNode />
                        <!-- this contains the offset blurred image -->
                        <feMergeNode in="SourceGraphic" />
                        <!-- this contains the element that the filter is applied to -->
                    </feMerge>
                </filter>
                <g :style="scrubberStageStyle">
                    <g ref="d3ScrubberStage"></g>
                </g>
            </svg>
        </div>
    </div>
</template>

<script>
import * as d3 from "d3";
import D3LineChart from "./D3LineChart";
import D3HistoChart from "./D3HistoChart";
import D3RangeChart from "./D3RangeChart";
import D3MapChart from "./D3MapChart";

export default {
    name: "D3Chart",
    components: {
        D3LineChart,
        D3HistoChart,
        D3RangeChart,
        D3MapChart,
    },
    props: ["chartParam"],
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
            scrubberMarginLeft: 5,
            chartRefs: ["d3LineChart", "d3HistoChart", "d3RangeChart", "d3MapChart"],
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
        scrubberHeight: function() {
            return 50;
        },
        scrubberViewBox: function() {
            return "0 0 " + this.outerWidth + " " + this.scrubberHeight;
        },
        scrubberStageStyle: function() {
            return {
                transform: "translate(" + this.scrubberMarginLeft + "px," + this.layout.marginTop + "px)",
            };
        },
    },
    mounted() {
        this.chart = this.chartParam;
        this.chart.svg = d3.select(this.$refs.d3Stage);
        this.chart.colors = this.chartParam.sensor.colorScale;
        this.chart.panelID = this.chart.station.deviceId;
        this.activateChart();
        this.initScrubber();
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
                case "Map":
                    this.$refs.d3MapChart.setStatus(true);
                    break;
                default:
                    this.$refs.d3LineChart.setStatus(true);
                    break;
            }
        },
        deactivateCharts() {
            this.chartRefs.forEach((r) => {
                if (this.$refs[r]) {
                    this.$refs[r].setStatus(false);
                }
            });
        },
        getTimeRange() {
            return { start: this.chart.start, end: this.chart.end };
        },
        setTimeRange(range) {
            this.chart.start = range.start;
            this.chart.end = range.end;
        },
        setRequestedTime(range) {
            this.requestedStart = range.start;
            this.requestedEnd = range.end;
        },
        getRequestedTime() {
            return [this.requestedStart, this.requestedEnd];
        },
        onTimeZoom(range) {
            this.requestedStart = range.start;
            this.requestedEnd = range.end;
            this.$emit("timeZoomed", { range: range, parent: this.chart.parent, id: this.chart.id });
        },
        updateData(data, extent, colorScale) {
            this.chart.data = data;
            this.chart.extent = extent;
            this.chart.colors = colorScale;
            this.chartRefs.forEach((r) => {
                if (this.$refs[r]) {
                    this.$refs[r].dataChanged();
                }
            });
            this.updateScrubber();
        },
        updateChartType() {
            this.deactivateCharts();
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
                case "Map":
                    this.$refs.d3MapChart.makeMap();
                    break;
            }
        },
        zoomOut() {
            this.$emit("allDays", { parent: this.chart.parent, id: this.chart.id });
        },
        initScrubber() {
            this.scrubberData = this.chart.overall.filter((d) => {
                return d[this.chart.sensor.key] === 0 || d[this.chart.sensor.key];
            });
            this.scrubberTimeRange = [];
            this.scrubberTimeRange[0] = this.chart.totalTime[0];
            this.scrubberTimeRange[1] = this.chart.totalTime[1];
            this.scrubberExtent = d3.extent(this.scrubberData, (d) => {
                return d[this.chart.sensor.key];
            });

            d3.selectAll(".d3scrubber-" + this.chart.id).remove();
            this.scrubberSVG = d3
                .select(this.$refs.d3ScrubberStage)
                .append("g")
                .attr("class", "d3scrubber-" + this.chart.id);

            this.scrubberFn = d3
                .brushX()
                .extent([
                    [0, 0],
                    [this.layout.width, this.scrubberHeight - this.layout.marginBottom],
                ])
                .on("start brush end", this.scrubberMoved);

            this.scrubberX = d3
                .scaleTime()
                .domain(this.scrubberTimeRange)
                .range([this.scrubberMarginLeft, this.layout.width - (this.layout.marginRight + this.scrubberMarginLeft)]);

            this.scrubberY = d3
                .scaleLinear()
                .domain(this.scrubberExtent)
                .range([this.scrubberHeight - (this.layout.marginBottom + this.layout.marginTop), this.layout.marginTop]);

            // area fn
            this.area = d3
                .area()
                .x((d) => {
                    return this.scrubberX(d.date);
                })
                .y0(this.scrubberHeight - (this.layout.marginBottom + this.layout.marginTop))
                .y1((d) => {
                    return this.scrubberY(d[this.chart.sensor.key]);
                });

            // add background rect
            this.scrubberSVG
                .append("rect")
                .attr("x", this.scrubberMarginLeft)
                .attr("y", 0)
                .attr("width", this.layout.width - this.scrubberMarginLeft * 2 - this.layout.marginRight)
                .attr("height", this.scrubberHeight)
                .attr("fill", "rgb(244, 245, 247)");
            // add gray area
            this.scrubberSVG
                .append("path")
                .data([this.scrubberData])
                .attr("class", "background-area")
                .attr("fill", "rgb(220, 222, 223)")
                .attr("d", this.area);
            // add the clip-path
            this.scrubberClip = this.scrubberSVG
                .append("defs")
                .append("svg:clipPath")
                .attr("id", "scrubber-clip" + this.chart.id)
                .append("svg:rect")
                .attr("width", this.scrubberX(this.chart.end) - this.scrubberX(this.chart.start))
                .attr("height", this.scrubberHeight)
                .attr("x", this.scrubberX(this.chart.start))
                .attr("y", 0);
            // add clipped selected area
            this.scrubberSVG
                .append("path")
                .data([this.scrubberData])
                .attr("class", "selected-scrubber")
                .attr("clip-path", "url(#scrubber-clip" + this.chart.id + ")")
                .attr("fill", "rgb(45, 158, 204)")
                .attr("d", this.area);
            // add the brushing
            this.scrubberUI = this.scrubberSVG
                .append("g")
                .attr("class", "scrubbrush")
                .call(this.scrubberFn);
            // add the handles
            this.handles = this.scrubberUI
                .selectAll(".handle--custom")
                .data([{ type: "w" }, { type: "e" }])
                .enter()
                .append("g")
                .attr("class", (d) => {
                    return d.type == "w" ? "left-handle" : "right-handle";
                })
                .attr("cursor", "ew-resize");
            this.handles
                .append("line")
                .attr("x1", 0)
                .attr("x2", 0)
                .attr("y1", 0)
                .attr("y2", this.scrubberHeight)
                .attr("stroke", "black");
            this.handles
                .append("circle")
                .attr("filter", "url(#dropshadow-" + this.chart.id + ")")
                .attr("cx", 0)
                .attr("cy", this.scrubberHeight / 2.0)
                .attr("r", 9)
                .attr("fill", "white");
            // set the initial position of the scrubber
            const start = this.chart.requestedStart ? this.chart.requestedStart : this.chart.start;
            const end = this.chart.requestedEnd ? this.chart.requestedEnd : this.chart.end;
            this.scrubberUI.call(this.scrubberFn.move, [start, end].map(this.scrubberX));
        },
        scrubberMoved() {
            if (!d3.event.selection) {
                return;
            }
            const xRange = d3.event.selection;
            // set bounds for scrubber handles
            if (xRange[0] < this.scrubberMarginLeft) {
                xRange[0] = this.scrubberMarginLeft;
            }
            if (xRange[1] > this.layout.width - this.scrubberMarginLeft - this.layout.marginRight) {
                xRange[1] = this.layout.width - this.scrubberMarginLeft - this.layout.marginRight;
            }
            const start = this.scrubberX.invert(xRange[0]);
            const end = this.scrubberX.invert(xRange[1]);
            // save these so scrubber can snap to them
            // instead of chart.start and chart.end
            this.requestedStart = start;
            this.requestedEnd = end;
            // adjust the selected area and handles to match selection
            this.scrubberClip.attr("x", xRange[0]).attr("width", xRange[1] - xRange[0]);
            this.handles.attr("transform", (d) => {
                const x = d.type == "w" ? xRange[0] : xRange[1];
                return "translate(" + x + ",0)";
            });
            // don't display the default d3 brush selection rect
            this.scrubberSVG
                .select(".selection")
                .attr("fill-opacity", 0)
                .attr("stroke", "none");

            // initiate chain of events to get new data
            if (d3.event.type == "end" && d3.event.sourceEvent) {
                this.onTimeZoom({ start: start, end: end });
            }
        },
        updateScrubber() {
            this.scrubberTimeRange = [];
            this.scrubberTimeRange[0] = this.chart.totalTime[0];
            this.scrubberTimeRange[1] = this.chart.totalTime[1];
            this.scrubberX.domain(this.scrubberTimeRange);

            this.scrubberData = this.chart.overall.filter((d) => {
                return d[this.chart.sensor.key] === 0 || d[this.chart.sensor.key];
            });
            this.scrubberExtent = d3.extent(this.scrubberData, (d) => {
                return d[this.chart.sensor.key];
            });
            this.scrubberY = d3
                .scaleLinear()
                .domain(this.scrubberExtent)
                .range([this.scrubberHeight - (this.layout.marginBottom + this.layout.marginTop), this.layout.marginTop]);

            // update the gray background
            this.scrubberSVG
                .selectAll(".background-area")
                .data([this.scrubberData])
                .transition()
                .attr("d", this.area);
            // updated the clipped selected area
            this.scrubberSVG
                .selectAll(".selected-scrubber")
                .data([this.scrubberData])
                .transition()
                .attr("d", this.area);

            // update the scrubber position
            const start = this.requestedStart ? this.requestedStart : this.chart.start;
            const end = this.requestedEnd ? this.requestedEnd : this.chart.end;
            this.scrubberUI.call(this.scrubberFn.move, [start, end].map(this.scrubberX));
        },
    },
};
</script>

<style scoped></style>
