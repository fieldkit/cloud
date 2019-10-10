<template>
    <div></div>
</template>

<script>
import * as d3 from "d3";

const binCount = 16;

export default {
    name: "D3HistoChart",
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
                this.makeHistogram();
            }
        }
    },
    methods: {
        init() {},
        setStatus(status) {
            this.activeMode = status;
        },
        prepareHistogram() {
            let d3Chart = this;
            this.xHist = d3
                .scaleLinear()
                .domain(this.chart.extent)
                .range([
                    this.layout.marginLeft,
                    this.layout.width - this.layout.marginRight - this.layout.marginLeft
                ]);

            // use this formula to create thresholds, so that bars will have the same width
            // and none of them will end outside the chart area
            const min = this.chart.extent[0];
            const max = this.chart.extent[1];
            const thresholds = d3.range(min, max, (max - min) / binCount);

            this.histogram = d3
                .histogram()
                .value(d => {
                    return d[d3Chart.selectedSensor.name];
                })
                .domain(this.xHist.domain())
                .thresholds(thresholds);

            let filteredData = this.processedData.filter(d => {
                return d.date > d3Chart.chart.x.domain()[0] && d.date < d3Chart.chart.x.domain()[1];
            });
            let bins = this.histogram(filteredData);

            this.yHist = d3
                .scaleLinear()
                .range([
                    this.layout.height - this.layout.marginTop - this.layout.marginBottom,
                    this.layout.marginTop
                ]);
            this.yHist.domain([
                0,
                d3.max(bins, d => {
                    return d.length;
                })
            ]);

            return bins;
        },
        makeHistogram() {
            let d3Chart = this;

            let bins = this.prepareHistogram();

            // append the bar rectangles
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
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .transition()
                .duration(1000)
                .attr("height", d => {
                    return (
                        d3Chart.layout.height -
                        d3Chart.yHist(d.length) -
                        d3Chart.layout.marginBottom -
                        d3Chart.layout.marginTop
                    );
                });

            // Axes
            this.xAxis = d3.axisBottom(this.xHist).ticks(binCount);
            this.yAxis = d3.axisLeft(this.yHist).ticks(10);

            // Add x axis
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

            // Add y axis
            this.yAxisGroup = this.chart.svg
                .append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.yAxis);

            this.drawn = true;
        },
        updateChart() {
            // Coming soon (need more data across time to test)
        },
        sensorChange() {
            let d3Chart = this;
            this.chart.extent = d3.extent(this.processedData, d => {
                return d[d3Chart.selectedSensor.name];
            });
            let bins = this.prepareHistogram();

            this.chart.colors = d3
                .scaleSequential()
                .domain(this.chart.extent)
                .interpolator(d3.interpolatePlasma);

            let bars = this.chart.svg.selectAll(".histobar").data(bins);

            // add any new bars
            bars.enter()
                .append("rect")
                .attr("class", "histobar")
                .attr("transform", d => {
                    return "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(d.length) + ")";
                })
                .attr("width", d => {
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .attr("height", d => {
                    return (
                        d3Chart.layout.height -
                        d3Chart.yHist(d.length) -
                        d3Chart.layout.marginBottom -
                        d3Chart.layout.marginTop
                    );
                });

            // updating any existing bars
            bars.transition()
                .duration(1000)
                .style("fill", d => d3Chart.chart.colors(d.x0))
                .attr("transform", d => {
                    return "translate(" + d3Chart.xHist(d.x0) + "," + d3Chart.yHist(d.length) + ")";
                })
                .attr("width", d => {
                    return d3Chart.xHist(d.x1) - d3Chart.xHist(d.x0) - 1;
                })
                .attr("height", d => {
                    return (
                        d3Chart.layout.height -
                        d3Chart.yHist(d.length) -
                        d3Chart.layout.marginBottom -
                        d3Chart.layout.marginTop
                    );
                });

            // remove any extra bars
            bars.exit().remove();

            // Axes
            this.xAxis = d3.axisBottom(this.xHist).ticks(binCount);
            this.yAxis = d3.axisLeft(this.yHist).ticks(10);

            // update the x axis
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(this.xAxis);

            // update y axis
            this.chart.y.domain(this.chart.extent);
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(this.yAxis);
        }
    }
};
</script>

<style scoped></style>
