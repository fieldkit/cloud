<template>
    <div :id="'scrubber-container-' + chart.id">
        <svg :view-box.camel="viewBox" preserveAspectRatio="xMidYMid meet" :width="outerWidth" :height="scrubberHeight">
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
            <g :style="stageStyle">
                <g ref="d3ScrubberStage"></g>
            </g>
        </svg>
    </div>
</template>

<script>
import * as d3 from "d3";

export default {
    name: "D3LineChart",
    props: ["chart", "layout", "summary"],
    data: () => {
        return {
            activeMode: false,
        };
    },
    watch: {
        chart: function() {
            // make a copy of chart.data because
            // we will add blank points
            this.chartData = this.chart.data.slice();

            if (!this.line) {
                this.init();
                this.initScrubber();
            }

            if (this.activeMode) {
                if (this.chart.svg.selectAll(".data-line, .bkgd-line").empty()) {
                    this.makeLine();
                }
            }
        },
    },
    computed: {
        outerWidth: function() {
            return this.layout.width + this.layout.marginLeft + this.layout.marginRight;
        },
        scrubberHeight: function() {
            return 50;
        },
        viewBox: function() {
            return "0 0 " + this.outerWidth + " " + this.scrubberHeight;
        },
        stageStyle: function() {
            return {
                transform: "translate(" + this.layout.marginLeft + "px," + this.layout.marginTop + "px)",
            };
        },
    },
    mounted() {
        if (!this.activeMode) {
            document.getElementById("scrubber-container-" + this.chart.id).style.display = "none";
        }
    },
    methods: {
        init() {
            let d3Chart = this;
            this.line = this.chart.svg
                .append("g")
                .attr("clip-path", "url(#clip" + this.chart.id + ")")
                .attr("class", "d3line");

            this.chart.svg
                .append("defs")
                .append("svg:clipPath")
                .attr("id", "clip" + this.chart.id)
                .append("svg:rect")
                .attr("width", this.layout.width - this.layout.marginLeft * 2 - this.layout.marginRight)
                .attr("height", this.layout.height)
                .attr("x", this.layout.marginLeft)
                .attr("y", 0);

            this.brush = d3
                .brushX()
                .extent([[0, 0], [this.layout.width, this.layout.height - this.layout.marginBottom]])
                .on("end", this.brushed);

            this.x = d3
                .scaleTime()
                .domain([this.chart.start, this.chart.end])
                .range([this.layout.marginLeft, this.layout.width - (this.layout.marginRight + this.layout.marginLeft)]);

            this.y = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([this.layout.height - (this.layout.marginBottom + this.layout.marginTop), this.layout.marginTop]);

            this.lineFn = d3
                .line()
                .defined(d => !d.blankPoint)
                .x(d => d3Chart.x(d.date))
                .y(d => d3Chart.y(d[d3Chart.chart.sensor.key]))
                .curve(d3.curveBasis);
        },
        initScrubber() {
            let d3Chart = this;
            this.scrubberData = this.summary.filter(d => {
                return d[d3Chart.chart.sensor.key] === 0 || d[d3Chart.chart.sensor.key];
            });
            this.scrubberTimeRange = [];
            this.scrubberTimeRange[0] = this.summary[0].date;
            this.scrubberTimeRange[1] = this.summary[this.summary.length - 1].date;
            this.scrubberExtent = d3.extent(this.scrubberData, d => {
                return d[d3Chart.chart.sensor.key];
            });

            d3.selectAll(".d3scrubber-" + this.chart.id).remove();
            this.scrubberSVG = d3
                .select(this.$refs.d3ScrubberStage)
                .append("g")
                .attr("class", "d3scrubber-" + this.chart.id);

            this.scrubberFn = d3
                .brushX()
                .extent([[0, 0], [this.layout.width, this.scrubberHeight - this.layout.marginBottom]])
                .on("start brush end", this.scrubberMoved);

            this.scrubberX = d3
                .scaleTime()
                .domain(this.scrubberTimeRange)
                .range([this.layout.marginLeft, this.layout.width - (this.layout.marginRight + this.layout.marginLeft)]);

            this.scrubberY = d3
                .scaleLinear()
                .domain(this.scrubberExtent)
                .range([this.scrubberHeight - (this.layout.marginBottom + this.layout.marginTop), this.layout.marginTop]);

            // area fn
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.scrubberX(d.date);
                })
                .y0(this.scrubberHeight - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.scrubberY(d[d3Chart.chart.sensor.key]);
                });

            // add background rect
            this.scrubberSVG
                .append("rect")
                .attr("x", this.layout.marginLeft)
                .attr("y", 0)
                .attr("width", this.layout.width - this.layout.marginLeft * 2 - this.layout.marginRight)
                .attr("height", this.scrubberHeight)
                .attr("fill", "rgb(244, 245, 247)");
            // add gray area
            this.scrubberSVG
                .append("path")
                .data([this.scrubberData])
                .attr("class", "area")
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
                .attr("id", "selected-scrubber")
                .attr("clip-path", "url(#scrubber-clip" + this.chart.id + ")")
                .attr("class", "area")
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
                .attr("class", d => {
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
            this.scrubberUI.call(this.scrubberFn.move, [this.chart.start, this.chart.end].map(this.scrubberX));
        },
        scrubberMoved() {
            if (!d3.event.selection) {
                return;
            }

            const xRange = d3.event.selection;
            // set bounds for scrubber handles
            if (xRange[0] < this.layout.marginLeft) {
                xRange[0] = this.layout.marginLeft;
            }
            if (xRange[1] > this.layout.width - this.layout.marginLeft - this.layout.marginRight) {
                xRange[1] = this.layout.width - this.layout.marginLeft - this.layout.marginRight;
            }
            const start = this.scrubberX.invert(xRange[0]);
            const end = this.scrubberX.invert(xRange[1]);
            // save these so scrubber can snap to them
            // instead of chart.start and chart.end
            this.requestedStart = start;
            this.requestedEnd = end;
            // adjust the selected area and handles to match selection
            this.scrubberClip.attr("x", xRange[0]).attr("width", xRange[1] - xRange[0]);
            this.handles.attr("transform", d => {
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
                this.$emit("timeZoomed", { start: start, end: end });
            }
        },
        setStatus(status) {
            this.activeMode = status;
            if (this.activeMode) {
                document.getElementById("scrubber-container-" + this.chart.id).style.display = "block";
            } else {
                document.getElementById("scrubber-container-" + this.chart.id).style.display = "none";
            }
        },
        setRequestedTime(range) {
            this.requestedStart = range.start;
            this.requestedEnd = range.end;
        },
        getRequestedTime() {
            return [this.requestedStart, this.requestedEnd];
        },
        dataChanged() {
            this.initScrubber();
            this.chartData = this.chart.data.slice();
            if (this.activeMode) {
                this.updateChart();
            }
        },
        makeLine() {
            let d3Chart = this;

            // create blank points for dotted line
            this.createBlankPoints();

            // add the gradient area
            this.createGradient();

            // add the lines
            this.chart.svg.selectAll(".data-line, .bkgd-line").remove();
            // continuous dotted line
            this.line
                .append("path")
                .attr("class", "bkgd-line")
                .attr("stroke", "#BBBBBB")
                .attr("stroke-dasharray", "4,4")
                .attr("fill", "none")
                .attr("d", this.lineFn(this.chartData.filter(this.lineFn.defined())));
            // data line with gaps
            this.line
                .append("path")
                .attr("class", "data-line")
                .attr("stroke", "url(#" + this.chart.id + "-area-gradient)")
                .attr("stroke-width", "3")
                .attr("fill", "none")
                .attr("d", this.lineFn(this.chartData));

            // add the brushing
            this.chart.svg.selectAll(".brush").remove();
            this.chart.svg
                .append("g")
                .attr("class", "brush")
                .attr("data-chart", this.chart.id)
                .call(this.brush);

            // add dots
            this.chart.svg.selectAll(".dot").remove();
            this.line
                .selectAll(".circles")
                .data(
                    this.chartData.filter(d => {
                        return !d.blankPoint;
                    })
                )
                .enter()
                .append("circle")
                .attr("class", "dot")
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.chart.sensor.key]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.chart.sensor.key]));
            // tooltip will be added back

            this.xAxis = d3.axisBottom(this.x).ticks(10);
            this.yAxis = d3.axisLeft(this.y).ticks(6);

            this.chart.svg.selectAll(".x-axis").remove();
            // add x axis
            this.xAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "x-axis")
                .attr("transform", "translate(" + 0 + "," + (this.layout.height - (this.layout.marginBottom + this.layout.marginTop)) + ")")
                .call(this.xAxis);

            this.chart.svg.selectAll(".y-axis").remove();
            // add y axis
            this.yAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "y-axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.yAxis);

            document.getElementById(this.chart.id + "-loading").style.display = "none";
        },
        createBlankPoints() {
            if (this.chartData.length > 1) {
                const start = this.chartData[0].date.getTime();
                const next = this.chartData[1].date.getTime();
                const last = this.chartData[this.chartData.length - 1].date.getTime();
                // TODO: this interval might be too small
                const interval = next - start;
                for (let i = start; i < last; i += interval) {
                    const point = this.chartData.find(d => {
                        return d.date.getTime() >= i - interval && d.date.getTime() <= i;
                    });
                    if (!point) {
                        // insert "blank" point
                        this.chartData.push({
                            date: new Date(i),
                            blankPoint: true,
                        });
                    }
                }
            }
            this.chartData.sort(function(a, b) {
                return a.date - b.date;
            });
        },
        brushed() {
            if (!d3.event.selection) {
                return;
            }

            let xRange = d3.event.selection;
            this.chart.start = this.x.invert(xRange[0]);
            this.chart.end = this.x.invert(xRange[1]);

            this.requestedStart = this.x.invert(xRange[0]);
            this.requestedEnd = this.x.invert(xRange[1]);

            // remove the grey brush area after selection
            this.chart.svg.select(".brush").call(this.brush.move, null);
            // initiate chain of events to get new data
            this.$emit("timeZoomed", { start: this.chart.start, end: this.chart.end });
        },
        updateChart() {
            let d3Chart = this;
            // update x scale
            this.x.domain([this.chart.start, this.chart.end]);
            // update y scale
            this.y.domain(this.chart.extent);
            // create blank points for dotted line
            this.createBlankPoints();
            // update the gradient
            this.createGradient();

            // update the scrubber position
            this.scrubberUI.call(this.scrubberFn.move, [this.requestedStart, this.requestedEnd].map(this.scrubberX));

            // update lines
            // continuous dotted line
            this.line
                .selectAll(".bkgd-line")
                .transition()
                .duration(1000)
                .attr("d", this.lineFn(this.chartData.filter(this.lineFn.defined())));
            // data line with gaps
            this.line
                .selectAll(".data-line")
                .transition()
                .duration(1000)
                .attr("d", this.lineFn(this.chartData));

            // update dots
            let dots = this.line.selectAll(".dot").data(
                this.chartData.filter(d => {
                    return !d.blankPoint;
                })
            );
            // add any new dots
            dots.enter()
                .append("circle")
                .attr("class", "dot")
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.chart.sensor.key]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.chart.sensor.key]));
            // updating any existing dots
            dots.transition()
                .duration(1000)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.chart.sensor.key]))
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.chart.sensor.key]);
                });
            // remove any extra dots
            dots.exit().remove();

            // update x axis
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisBottom(this.x));

            // update y axis
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisLeft(this.y));

            document.getElementById(this.chart.id + "-loading").style.display = "none";
        },
        createGradient() {
            // Add the gradient area
            this.chart.svg.selectAll("#" + this.chart.id + "-area-gradient").remove();
            this.line
                .append("linearGradient")
                .attr("id", this.chart.id + "-area-gradient")
                .attr("gradientUnits", "userSpaceOnUse")
                .attr("x1", 0)
                .attr("y1", this.y(this.chart.extent[0]))
                .attr("x2", 0)
                .attr("y2", this.y(this.chart.extent[1]))
                .selectAll("stop")
                .data([
                    { offset: "0%", color: this.chart.colors(this.chart.extent[0]) },
                    {
                        offset: "20%",
                        color: this.chart.colors(this.chart.extent[0] + 0.2 * (this.chart.extent[1] - this.chart.extent[0])),
                    },
                    {
                        offset: "40%",
                        color: this.chart.colors(this.chart.extent[0] + 0.4 * (this.chart.extent[1] - this.chart.extent[0])),
                    },
                    {
                        offset: "60%",
                        color: this.chart.colors(this.chart.extent[0] + 0.6 * (this.chart.extent[1] - this.chart.extent[0])),
                    },
                    {
                        offset: "80%",
                        color: this.chart.colors(this.chart.extent[0] + 0.8 * (this.chart.extent[1] - this.chart.extent[0])),
                    },
                    {
                        offset: "100%",
                        color: this.chart.colors(this.chart.extent[1]),
                    },
                ])
                .enter()
                .append("stop")
                .attr("offset", d => {
                    return d.offset;
                })
                .attr("stop-color", d => {
                    return d.color;
                });
        },
    },
};
</script>

<style scoped></style>
