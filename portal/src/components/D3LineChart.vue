<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

export default {
    name: "D3LineChart",
    props: ["chart", "stationData", "layout", "selectedSensor"],
    data: () => {
        return {
            activeMode: false,
            drawn: false
        };
    },
    watch: {
        selectedSensor: function() {
            if (this.drawn && this.activeMode) {
                this.sensorChange();
            }
        },
        stationData: function() {
            if (this.activeMode) {
                if (this.chart.svg.selectAll(".data-line").empty()) {
                    this.makeLine();
                } else {
                    let d3Chart = this;
                    this.filteredData = this.stationData.filter(d => {
                        return d[d3Chart.selectedSensor.key];
                    });
                    // this.chart.extent = d3.extent(this.filteredData, d => {
                    //     return d[d3Chart.selectedSensor.key];
                    // });

                    this.timeChanged();
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

            // this.colors = d3
            //     .scaleSequential()
            //     .domain(this.chart.extent)
            //     .interpolator(d3.interpolatePlasma);

            let d3Chart = this;
            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.key]);
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
                .defined(d => {
                    return !d.blankPoint;
                })
                .x(d => {
                    return d3Chart.x(d.date);
                })
                .y(d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.key]);
                })
                .curve(d3.curveBasis);
        },
        setStatus(status) {
            this.activeMode = status;
        },
        makeLine() {
            let d3Chart = this;
            this.filteredData = this.stationData.filter(d => {
                return d[d3Chart.selectedSensor.key];
            });
            // this.createBlankPoints();

            // Add the gradient area
            this.createGradient();

            // Add the line
            // this.line
            //     .append("path")
            //     .data([this.filteredData])
            //     .attr("class", "area")
            //     .attr("stroke-width", "3")
            //     .attr("stroke", "url(#" + this.chart.id + "-area-gradient)")
            //     .attr("fill", "none")
            //     .transition()
            //     .duration(1000)
            //     .attr("d", this.area);

            this.chart.svg.selectAll(".data-line").remove();
            this.line
                .append("path")
                .attr("class", "data-line")
                .attr("d", this.lineFn(this.filteredData))
                .attr("stroke-width", "3")
                .attr("stroke", "url(#" + this.chart.id + "-area-gradient)")
                .attr("fill", "none");

            // Add the brushing
            this.chart.svg.selectAll(".brush").remove();
            this.chart.svg
                .append("g")
                .attr("class", "brush")
                .attr("data-panel", this.chart.panelID)
                .call(this.brush);

            // Add dots
            this.chart.svg.selectAll(".dot").remove();
            this.line
                .selectAll(".circles")
                .data(
                    this.filteredData.filter(d => {
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
                    return d3Chart.y(d[d3Chart.selectedSensor.key]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.selectedSensor.key]));
            // tooltip will be added back

            this.xAxis = d3.axisBottom(this.x).ticks(10);
            this.yAxis = d3.axisLeft(this.y).ticks(6);

            this.chart.svg.selectAll(".x-axis").remove();
            //Add x axis
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
            //Add y axis
            this.yAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "y-axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.yAxis);

            this.drawn = true;
            document.getElementById("loading").style.display = "none";
        },
        createBlankPoints() {
            // TODO: relying on span between first two points is not reliable
            // use better method for creating "blank" points
            if (this.filteredData.length > 1) {
                const start = this.filteredData[0].date.getTime();
                const next = this.filteredData[1].date.getTime();
                const last = this.filteredData[this.filteredData.length - 1].date.getTime();
                const interval = next - start;
                for (let i = start; i < last; i += interval) {
                    const point = this.filteredData.find(d => {
                        return d.date.getTime() >= i - interval && d.date.getTime() <= i;
                    });
                    if (!point) {
                        // insert "blank" point
                        this.filteredData.push({
                            date: new Date(i),
                            blankPoint: true
                        });
                    }
                }
            }
            this.filteredData.sort(function(a, b) {
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
            // Remove the grey brush area after selection
            this.chart.svg.select(".brush").call(this.brush.move, null);
            this.timeChanged();
            this.$emit("timeZoomed", { start: this.chart.start, end: this.chart.end });
        },
        timeChanged() {
            // update x scale
            this.x.domain([this.chart.start, this.chart.end]);

            // Update x axis and line position
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisBottom(this.x));

            this.updateChart();
        },
        sensorChange() {
            let d3Chart = this;
            this.filteredData = this.stationData.filter(d => {
                return d[d3Chart.selectedSensor.key];
            });
            this.chart.extent = d3.extent(this.filteredData, d => {
                return d[d3Chart.selectedSensor.key];
            });
            this.updateChart();
        },
        updateChart() {
            let d3Chart = this;

            // update y domain with new extent
            this.y.domain(this.chart.extent);
            // this.createBlankPoints();

            this.createGradient();

            // update line area
            // this.line
            //     .select(".area")
            //     .transition()
            //     .duration(1000)
            //     .attr("d", this.area(this.filteredData));
            this.line
                .selectAll(".data-line")
                .transition()
                .duration(1000)
                .attr("d", this.lineFn(this.filteredData));

            // update dots
            let dots = this.line.selectAll(".dot").data(
                this.filteredData.filter(d => {
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
                    return d3Chart.y(d[d3Chart.selectedSensor.key]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.selectedSensor.key]));

            // updating any existing dots
            dots.transition()
                .duration(1000)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.selectedSensor.key]))
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.key]);
                });

            // remove any extra dots
            dots.exit().remove();

            // update y axis
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisLeft(this.y));
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
