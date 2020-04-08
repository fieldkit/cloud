<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

const binCount = 16;
const MIN_WIDTH = 800 / 16;

export default {
    name: "D3HistoChart",
    props: ["chart", "layout"],
    data: () => {
        return {
            activeMode: false,
        };
    },
    watch: {
        chart: function() {
            if (this.activeMode) {
                this.makeHistogram();
            }
        },
    },
    methods: {
        setStatus(status) {
            this.activeMode = status;
        },
        dataChanged() {
            if (this.activeMode) {
                this.makeHistogram();
            }
        },
        prepareHistogram() {
            let d3Chart = this;
            // set x scale
            this.xHist = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([this.layout.marginLeft, this.layout.width - this.layout.marginRight - this.layout.marginLeft]);

            // use this formula to create thresholds, so that bars will have the same width
            // and none of them will end outside the chart area
            const min = this.chart.extent[0];
            const max = this.chart.extent[1];
            const thresholds = d3.range(min, max, (max - min) / binCount);

            // set histogram function
            this.histogram = d3
                .histogram()
                .value(d => {
                    return d[d3Chart.chart.sensor.key];
                })
                .domain(this.xHist.domain())
                .thresholds(thresholds);

            // apply histogram function
            let bins = this.histogram(this.chart.data);

            // set y scale
            this.yHist = d3
                .scaleLinear()
                .domain([
                    0,
                    d3.max(bins, d => {
                        return d.length;
                    }),
                ])
                .range([this.layout.height - this.layout.marginTop - this.layout.marginBottom, this.layout.marginTop]);

            // set axes
            this.xAxis = d3.axisBottom(this.xHist).ticks(binCount);
            this.yAxis = d3.axisLeft(this.yHist).ticks(10);

            return bins;
        },
        makeHistogram() {
            let d3Chart = this;

            let bins = this.prepareHistogram();

            // append the bar rectangles
            this.chart.svg.selectAll(".histobar").remove();
            this.chart.svg
                .selectAll(".histobar")
                .data(bins)
                .enter()
                .append("rect")
                .attr("class", "histobar")
                .attr("transform", d => {
                    return "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(d.length) + ")";
                })
                .attr("width", d => {
                    const w = d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                    return w <= 0 ? MIN_WIDTH : w;
                })
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .transition()
                .duration(1000)
                .attr("height", d => {
                    return d.length == 0
                        ? 0
                        : d3Chart.layout.height - d3Chart.yHist(d.length) - d3Chart.layout.marginBottom - d3Chart.layout.marginTop;
                });

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
        updateHistogram(bins) {
            let d3Chart = this;

            let bars = this.chart.svg.selectAll(".histobar").data(bins);

            // add any new bars
            bars.enter()
                .append("rect")
                .attr("class", "histobar")
                .attr("transform", d => {
                    return "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(d.length) + ")";
                })
                .attr("width", d => {
                    const w = d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                    return w <= 0 ? MIN_WIDTH : w;
                })
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .attr("height", d => {
                    return d.length == 0
                        ? 0
                        : d3Chart.layout.height - d3Chart.yHist(d.length) - d3Chart.layout.marginBottom - d3Chart.layout.marginTop;
                });

            // updating any existing bars
            bars.transition()
                .duration(1000)
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .attr("transform", d => {
                    return "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(d.length) + ")";
                })
                .attr("width", d => {
                    const w = d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                    return w <= 0 ? MIN_WIDTH : w;
                })
                .attr("height", d => {
                    return d.length == 0
                        ? 0
                        : d3Chart.layout.height - d3Chart.yHist(d.length) - d3Chart.layout.marginBottom - d3Chart.layout.marginTop;
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
