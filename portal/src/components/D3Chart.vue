<template>
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
</template>

<script>
import * as d3 from "d3";

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "D3Chart",
    props: ["stationData", "selectedSensor", "timeRange"],
    data: () => {
        return {
            panelID: "",
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
            this.processData();
            this.initSVG();
            this.makeLine();
        },
        selectedSensor: function() {
            this.sensorChange();
        },
        timeRange: function() {
            this.timeChange();
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
            this.processedData = [];
            this.stationData.versions.forEach(v => {
                let station = v.meta.station;
                this.panelID = station.id;
                v.data.forEach(d => {
                    d.d.date = new Date(d.time * 1000);
                    this.processedData.push(d.d);
                });
            });
            //sort data by date
            this.processedData.sort(function(a, b) {
                return a.date.getTime() - b.date.getTime();
            });
            let d3Chart = this;
            this.currentMinMax = [
                d3.min(this.processedData, d => {
                    return d[d3Chart.selectedSensor.name];
                }),
                d3.max(this.processedData, d => {
                    return d[d3Chart.selectedSensor.name];
                })
            ];
            // x and y map functions
            this.end = this.processedData[this.processedData.length - 1].date;
            this.start = this.processedData[0].date;
            this.x = d3
                .scaleTime()
                .domain([this.start, this.end])
                .range([
                    this.layout.marginLeft,
                    this.layout.width - (this.layout.marginRight + this.layout.marginLeft)
                ]);
            this.y = d3
                .scaleLinear()
                .domain(this.currentMinMax)
                .range([
                    this.layout.height - (this.layout.marginBottom + this.layout.marginTop),
                    this.layout.marginTop
                ]);
            this.xAxis = d3.axisBottom(this.x).ticks(10);
            this.yAxis = d3.axisLeft(this.y).ticks(6);
        },
        initSVG() {
            this.svg = d3.select(this.$refs.d3Stage);

            this.svg
                .append("defs")
                .append("svg:clipPath")
                .attr("id", "clip" + this.panelID)
                .append("svg:rect")
                .attr("width", this.layout.width - this.layout.marginLeft * 2 - this.layout.marginRight)
                .attr("height", this.layout.height)
                .attr("x", this.layout.marginLeft)
                .attr("y", 0);

            this.line = this.svg
                .append("g")
                .attr("clip-path", "url(#clip" + this.panelID + ")")
                .attr("class", "d3line");

            this.colors = d3
                .scaleSequential()
                .domain(this.currentMinMax)
                .interpolator(d3.interpolatePlasma);

            let d3Chart = this;
            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.name]);
                })
                .curve(d3.curveBasis);

            this.brush = d3
                .brushX()
                .extent([[0, 0], [this.layout.width, this.layout.height - this.layout.marginBottom]])
                .on("end", this.brushed);
        },
        makeLine() {
            // Add the gradient area
            this.line
                .append("linearGradient")
                .attr("id", "area-gradient")
                .attr("gradientUnits", "userSpaceOnUse")
                .attr("x1", 0)
                .attr("y1", this.y(this.currentMinMax[0]))
                .attr("x2", 0)
                .attr("y2", this.y(this.currentMinMax[1]))
                .selectAll("stop")
                .data([
                    { offset: "0%", color: this.colors(this.currentMinMax[0]) },
                    {
                        offset: "20%",
                        color: this.colors(
                            this.currentMinMax[0] + 0.2 * (this.currentMinMax[1] - this.currentMinMax[0])
                        )
                    },
                    {
                        offset: "40%",
                        color: this.colors(
                            this.currentMinMax[0] + 0.4 * (this.currentMinMax[1] - this.currentMinMax[0])
                        )
                    },
                    {
                        offset: "60%",
                        color: this.colors(
                            this.currentMinMax[0] + 0.6 * (this.currentMinMax[1] - this.currentMinMax[0])
                        )
                    },
                    {
                        offset: "80%",
                        color: this.colors(
                            this.currentMinMax[0] + 0.8 * (this.currentMinMax[1] - this.currentMinMax[0])
                        )
                    },
                    {
                        offset: "100%",
                        color: this.colors(this.currentMinMax[1])
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
            this.line
                .append("g")
                .attr("class", "brush")
                .attr("data-panel", this.panelID)
                .call(this.brush);

            //Add dots
            let d3Chart = this;
            this.line
                .selectAll(".circles")
                .data(this.processedData)
                .enter()
                .append("circle")
                .attr("class", "dot")
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.name]);
                })
                .attr("r", 2)
                .attr("fill", d => d3Chart.colors(d[d3Chart.selectedSensor.name]));
            // tooltip will be added back

            //Add x axis
            this.xAxisGroup = this.svg
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

            //Add y axis
            this.yAxisGroup = this.svg
                .append("g")
                .attr("class", "y axis")
                .attr("transform", "translate(" + this.layout.marginLeft + ",0)")
                .call(this.yAxis);
        },
        brushed() {
            if (!d3.event.selection) {
                return;
            }

            let xRange = d3.event.selection;
            this.start = this.x.invert(xRange[0]);
            this.end = this.x.invert(xRange[1]);
            this.x.domain([this.start, this.end]);
            // Remove the grey brush area after selection
            this.line.select(".brush").call(this.brush.move, null);
            this.updateChart();
        },
        timeChange() {
            this.start = this.processedData[0].date;
            this.end = new Date(this.start.getTime() + this.timeRange * DAY);
            if (this.timeRange == 0) {
                this.end = this.processedData[this.processedData.length - 1].date;
            }
            this.x.domain([this.start, this.end]);
            this.updateChart();
        },
        sensorChange() {
            let d3Chart = this;
            this.currentMinMax = [
                d3.min(this.processedData, d => {
                    return d[d3Chart.selectedSensor.name];
                }),
                d3.max(this.processedData, d => {
                    return d[d3Chart.selectedSensor.name];
                })
            ];

            // Area gradient fill
            this.area = d3
                .area()
                .x(d => {
                    return d3Chart.x(d.date);
                })
                .y0(this.layout.height - (this.layout.marginBottom + this.layout.marginTop))
                .y1(d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.name]);
                })
                .curve(d3.curveBasis);

            // Update y axis
            this.y.domain(this.currentMinMax);
            this.yAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisLeft(this.y));

            this.colors = d3
                .scaleSequential()
                .domain(this.currentMinMax)
                .interpolator(d3.interpolatePlasma);

            this.updateChart();
        },
        updateChart() {
            // Update x axis and line position
            this.xAxisGroup
                .transition()
                .duration(1000)
                .call(d3.axisBottom(this.x));
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
                .attr("fill", d => d3Chart.colors(d[d3Chart.selectedSensor.name]))
                .attr("cx", d => {
                    return d3Chart.x(d.date);
                })
                .attr("cy", d => {
                    return d3Chart.y(d[d3Chart.selectedSensor.name]);
                });
            // setting location in url will be added back
        }
    }
};
</script>

<style scoped></style>
