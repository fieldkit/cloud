<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

export default {
    name: "D3LineChart",
    props: ["chart", "layout"],
    data: () => {
        return {
            activeMode: false
        };
    },
    watch: {
        chart: function() {
            // make a copy of chart.data because
            // we will add blank points
            this.chartData = this.chart.data.slice();

            if (!this.line) {
                this.init();
            }

            if (this.activeMode) {
                if (this.chart.svg.selectAll(".data-line, .bkgd-line").empty()) {
                    this.makeLine();
                }
            }
        }
    },
    methods: {
        init() {
            this.line = this.chart.svg
                .append("g")
                .attr("clip-path", "url(#clip" + this.chart.panelID + ")")
                .attr("class", "d3line");

            this.chart.svg
                .append("defs")
                .append("svg:clipPath")
                .attr("id", "clip" + this.chart.panelID)
                .append("svg:rect")
                .attr("width", this.layout.width - this.layout.marginLeft * 2 - this.layout.marginRight)
                .attr("height", this.layout.height)
                .attr("x", this.layout.marginLeft)
                .attr("y", 0);

            this.brush = d3
                .brushX()
                .extent([[0, 0], [this.layout.width, this.layout.height - this.layout.marginBottom]])
                .on("end", this.brushed);

            let d3Chart = this;
            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.y(d3Chart.y(d[d3Chart.chart.sensor.key]));
                })
                .curve(d3.curveBasis);

            this.x = d3
                .scaleTime()
                .domain([this.chart.start, this.chart.end])
                .range([
                    this.layout.marginLeft,
                    this.layout.width - (this.layout.marginRight + this.layout.marginLeft)
                ]);

            this.y = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([
                    this.layout.height - (this.layout.marginBottom + this.layout.marginTop),
                    this.layout.marginTop
                ]);

            this.lineFn = d3
                .line()
                .defined(d => !d.blankPoint)
                .x(d => d3Chart.x(d.date))
                .y(d => d3Chart.y(d[d3Chart.chart.sensor.key]))
                .curve(d3.curveBasis);
        },
        setStatus(status) {
            this.activeMode = status;
        },
        dataChanged() {
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
                .attr("data-panel", this.chart.panelID)
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
                .attr(
                    "transform",
                    "translate(" +
                        0 +
                        "," +
                        (this.layout.height - (this.layout.marginBottom + this.layout.marginTop)) +
                        ")"
                )
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
                            blankPoint: true
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
                        color: this.chart.colors(
                            this.chart.extent[0] + 0.2 * (this.chart.extent[1] - this.chart.extent[0])
                        )
                    },
                    {
                        offset: "40%",
                        color: this.chart.colors(
                            this.chart.extent[0] + 0.4 * (this.chart.extent[1] - this.chart.extent[0])
                        )
                    },
                    {
                        offset: "60%",
                        color: this.chart.colors(
                            this.chart.extent[0] + 0.6 * (this.chart.extent[1] - this.chart.extent[0])
                        )
                    },
                    {
                        offset: "80%",
                        color: this.chart.colors(
                            this.chart.extent[0] + 0.8 * (this.chart.extent[1] - this.chart.extent[0])
                        )
                    },
                    {
                        offset: "100%",
                        color: this.chart.colors(this.chart.extent[1])
                    }
                ])
                .enter()
                .append("stop")
                .attr("offset", d => {
                    return d.offset;
                })
                .attr("stop-color", d => {
                    return d.color;
                });
        }
    }
};
</script>

<style scoped></style>
