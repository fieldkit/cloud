<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

export default {
    name: "D3LineChart",
    props: ["chart", "processedData", "layout", "selectedSensor"],
    data: () => {
        return {
            activeMode: false,
            drawn: false
        };
    },
    watch: {
        chart: {
            handler: function() {
                if (this.drawn && this.activeMode) {
                    this.updateChart();
                }
            },
            deep: true
        },
        selectedSensor: function() {
            if (this.drawn && this.activeMode) {
                this.sensorChange();
            }
        },
        processedData: function() {
            if (this.activeMode) {
                this.makeLine();
            }
        }
    },
    methods: {
        init() {
            this.line = this.chart.svg
                .append("g")
                .attr("clip-path", "url(#clip" + this.chart.panelID + ")")
                .attr("class", "d3line");

            let d3Chart = this;
            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.chart.y(d[d3Chart.selectedSensor.name]);
                })
                .curve(d3.curveBasis);
        },
        setStatus(status) {
            this.activeMode = status;
        },
        makeLine() {
            // Add the gradient area
            this.line
                .append("linearGradient")
                .attr("id", "area-gradient")
                .attr("gradientUnits", "userSpaceOnUse")
                .attr("x1", 0)
                .attr("y1", this.chart.y(this.chart.extent[0]))
                .attr("x2", 0)
                .attr("y2", this.chart.y(this.chart.extent[1]))
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

            // Add the line
            this.line
                .append("path")
                .data([this.processedData])
                .attr("class", "area")
                .attr("fill", "url(#area-gradient)")
                .attr("stroke", "none")
                .transition()
                .duration(1000)
                .attr("d", this.area);

            // Add the brushing
            this.$emit("addBrush");

            //Add dots
            let d3Chart = this;
            this.line
                .selectAll(".circles")
                .data(this.processedData)
                .enter()
                .append("circle")
                .attr("class", "dot")
                .attr("cx", d => {
                    return d3Chart.chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.chart.y(d[d3Chart.selectedSensor.name]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.selectedSensor.name]));
            // tooltip will be added back

            //Add x axis
            this.xAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "x axis")
                .attr(
                    "transform",
                    "translate(" +
                        0 +
                        "," +
                        (this.layout.height - (this.layout.marginBottom + this.layout.marginTop)) +
                        ")"
                )
                .call(this.chart.xAxis);

            //Add y axis
            this.yAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.chart.yAxis);

            this.drawn = true;
        },
        updateChart() {
            // Update x axis and line position
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisBottom(this.chart.x));
            this.line
                .select(".area")
                .transition()
                .duration(1000)
                .attr("d", this.area(this.processedData));

            let d3Chart = this;
            this.line
                .selectAll(".dot")
                .transition()
                .duration(1000)
                .attr("fill", d => d3Chart.chart.colors(d[d3Chart.selectedSensor.name]))
                .attr("cx", d => {
                    return d3Chart.chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.chart.y(d[d3Chart.selectedSensor.name]);
                });
            // setting location in url will be added back
        },
        sensorChange() {
            let d3Chart = this;
            this.chart.extent = d3.extent(this.processedData, d => {
                return d[d3Chart.selectedSensor.name];
            });

            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.chart.y(d[d3Chart.selectedSensor.name]);
                })
                .curve(d3.curveBasis);

            // Update y axis
            this.chart.y.domain(this.chart.extent);
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisLeft(this.chart.y));

            this.chart.colors = d3
                .scaleSequential()
                .domain(this.chart.extent)
                .interpolator(d3.interpolatePlasma);
        }
    }
};
</script>

<style scoped></style>
