<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

const HOUR = 1000 * 60 * 60;
const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "D3RangeChart",
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
                    this.chartChanged();
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
                this.makeRange();
            }
        }
    },
    methods: {
        setStatus(status) {
            this.activeMode = status;
        },
        prepareRange() {
            this.colors = d3
                .scaleSequential()
                .domain(this.chart.extent)
                .interpolator(d3.interpolatePlasma);

            // set x scale
            this.xHist = d3
                .scaleTime()
                .domain([this.chart.start, this.chart.end])
                .range([
                    this.layout.marginLeft,
                    this.layout.width - this.layout.marginLeft - this.layout.marginRight
                ]);

            let interval = DAY;
            this.timeRange = this.chart.end - this.chart.start;
            if (this.timeRange < DAY) {
                interval = HOUR;
            }

            // use this formula to create thresholds, so that bars will have the same width
            // and none of them will end outside the chart area
            const min = this.chart.start;
            const max = this.chart.end;
            const thresholds = d3.range(min, max, (max - min) / 64);

            // set histogram function
            this.histogram = d3
                .histogram()
                .value(d => {
                    return d.date;
                })
                .domain(this.xHist.domain())
                .thresholds(thresholds);

            // apply histogram function
            let bins = this.histogram(this.processedData);

            // set y scale
            this.yHist = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([
                    this.layout.height - this.layout.marginBottom - this.layout.marginTop,
                    this.layout.marginTop
                ]);

            // set axes
            this.xAxis = d3.axisBottom(this.xHist).ticks(Math.min(16, this.timeRange / interval));
            this.yAxis = d3.axisLeft(this.yHist).ticks(6);

            return bins;
        },
        makeRange() {
            let d3Chart = this;
            let bins = this.prepareRange();

            // color gradient for each bin
            bins.forEach((bin, i) => {
                let gradient = this.chart.svg
                    .append("defs")
                    .attr("class", "range-gradient")
                    .append("linearGradient")
                    .attr("id", "grad-" + i)
                    .attr("x1", "0%")
                    .attr("x2", "0%")
                    .attr("y1", "0%")
                    .attr("y2", "100%");

                gradient
                    .append("stop")
                    .attr("offset", "0%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.max(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);

                gradient
                    .append("stop")
                    .attr("offset", "50%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.median(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);

                gradient
                    .append("stop")
                    .attr("offset", "100%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.min(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);
            });

            // append the bar rectangles
            this.chart.svg
                .selectAll(".rangebar")
                .data(bins)
                .enter()
                .append("rect")
                .attr("class", "rangebar")
                .attr("transform", d => {
                    const mx = d3.max(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const r = mx
                        ? "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(mx) + ")"
                        : "translate(0,0)";
                    return r;
                })
                .attr("width", d => {
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .style("fill", (d, i) => {
                    return d.length > 0 ? "url(#grad-" + i + ")" : "none";
                })
                .transition()
                .duration(1000)
                .attr("height", d => {
                    const extent = d3.extent(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const height = d3Chart.yHist(extent[1]) - d3Chart.yHist(extent[0]);
                    return -(height ? height : 0);
                });

            // add x axis
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
                .call(this.xAxis);

            // add y axis
            this.yAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.yAxis);

            this.drawn = true;
        },
        chartChanged() {
            // chart changes when start and end dates change
            // but possibly at other events, too - TODO: distinguish btw those
            let bins = this.prepareRange();
            this.updateRange(bins);
        },
        sensorChange() {
            let d3Chart = this;
            // define extent for this sensor
            this.chart.extent = d3.extent(this.processedData, d => {
                return d[d3Chart.selectedSensor.name];
            });
            let bins = this.prepareRange();
            this.updateRange(bins);
        },
        updateRange(bins) {
            let d3Chart = this;

            this.chart.svg.selectAll(".range-gradient").remove();
            // color gradient for each bin
            bins.forEach((bin, i) => {
                let gradient = this.chart.svg
                    .append("defs")
                    .attr("class", "range-gradient")
                    .append("linearGradient")
                    .attr("id", "grad-" + i)
                    .attr("x1", "0%")
                    .attr("x2", "0%")
                    .attr("y1", "0%")
                    .attr("y2", "100%");

                gradient
                    .append("stop")
                    .attr("offset", "0%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.max(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);

                gradient
                    .append("stop")
                    .attr("offset", "50%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.median(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);

                gradient
                    .append("stop")
                    .attr("offset", "100%")
                    .style(
                        "stop-color",
                        this.colors(
                            d3.min(bin, d => {
                                return d[d3Chart.selectedSensor.name];
                            })
                        )
                    )
                    .style("stop-opacity", 1);
            });

            let bars = this.chart.svg.selectAll(".rangebar").data(bins);

            // add any new bars
            bars.enter()
                .append("rect")
                .attr("class", "rangebar")
                .attr("transform", d => {
                    const mx = d3.max(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const r = mx
                        ? "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(mx) + ")"
                        : "translate(0,0)";
                    return r;
                })
                .attr("width", d => {
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .style("fill", (d, i) => {
                    return d.length > 0 ? "url(#grad-" + i + ")" : "none";
                })
                .attr("height", d => {
                    const extent = d3.extent(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const height = d3Chart.yHist(extent[1]) - d3Chart.yHist(extent[0]);
                    return -(height ? height : 0);
                });

            // updating any existing bars
            bars.attr("height", 0)
                .attr("transform", d => {
                    const mx = d3.max(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const r = mx
                        ? "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(mx) + ")"
                        : "translate(0,0)";
                    return r;
                })
                .attr("width", d => {
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .style("fill", (d, i) => {
                    return d.length > 0 ? "url(#grad-" + i + ")" : "none";
                })
                .transition()
                .duration(1000)
                .attr("height", d => {
                    const extent = d3.extent(d, b => {
                        return b[d3Chart.selectedSensor.name];
                    });
                    const height = d3Chart.yHist(extent[1]) - d3Chart.yHist(extent[0]);
                    return -(height ? height : 0);
                });

            // remove any extra bars
            bars.exit().remove();

            // update x axis
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(this.xAxis);

            // update y axis
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(this.yAxis);
        }
    }
};
</script>

<style scoped></style>
