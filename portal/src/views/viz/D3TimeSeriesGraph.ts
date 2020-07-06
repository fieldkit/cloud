import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData } from "./viz";

export const D3TimeSeriesGraph = Vue.extend({
    name: "D3TimeSeriesGraph",
    data() {
        return {};
    },
    props: {
        viz: {
            type: Graph,
            required: true,
        },
    },
    computed: {
        data(): QueriedData {
            if (this.viz.data && !this.viz.data.empty) {
                return this.viz.data;
            }
            return null;
        },
    },
    watch: {
        data(newValue, oldValue) {
            this.viz.log("graphing");

            const layout = new ChartLayout(1050, 340, new Margins({ top: 5, bottom: 50, left: 50, right: 0 }));
            const data = newValue;
            const timeRange = data.timeRange;
            const dataRange = data.dataRange;
            const charts = [
                {
                    layout: layout,
                },
            ];

            const x = d3
                .scaleTime()
                .domain(data.timeRange)
                .range([layout.margins.left, layout.width - (layout.margins.right + layout.margins.left)]);

            const y = d3
                .scaleLinear()
                .domain(data.dataRange)
                .range([layout.height - (layout.margins.bottom + layout.margins.top), layout.margins.top]);

            const xAxis = d3.axisBottom(x).ticks(10);

            const yAxis = d3.axisLeft(y).ticks(6);

            const svg = d3
                .select(this.$el)
                .select(".graph")
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const adding = enter
                        .append("svg")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("width", (c) => c.layout.width)
                        .attr("height", (c) => c.layout.height);

                    adding
                        .append("defs")
                        .append("svg:clipPath")
                        .attr("id", "clip-" + this.viz.id)
                        .append("svg:rect")
                        .attr("width", (c) => c.layout.width - c.layout.margins.left * 2 - c.layout.margins.right)
                        .attr("height", (c) => c.layout.height)
                        .attr("x", (c) => c.layout.margins.left)
                        .attr("y", (c) => c.layout.margins.bottom);

                    return adding;
                });

            svg.selectAll(".x-axis")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "x-axis")
                        .attr("transform", "translate(" + 0 + "," + (layout.height - (layout.margins.bottom + layout.margins.top)) + ")")
                )
                .call(xAxis);

            svg.selectAll(".y-axis")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "y-axis")
                        .attr("transform", "translate(" + layout.margins.left + ",0)")
                )
                .call(yAxis);

            const brush = d3
                .brushX()
                .extent([
                    [0, 0],
                    [layout.width, layout.height - layout.margins.bottom],
                ])
                .on("end", (_undefined, _zero, ev, ...args) => {
                    if (!d3.event.selection) {
                        return;
                    }

                    svg.select(".brush").call(brush.move, null);

                    const range = d3.event.selection;
                    const start = x.invert(range[0]);
                    const end = x.invert(range[1]);
                    const newRange = new TimeRange(new Time(start), new Time(end));
                    this.raiseTimeZoomed(newRange);
                });

            const colors = d3
                .scaleSequential()
                .domain([0, 1])
                .interpolator(() => {
                    return "#000000";
                });

            const line = svg
                .selectAll(".d3-line")
                .data(charts)
                .join((enter) => {
                    const adding = enter
                        .append("g")
                        .attr("clip-path", "url(#clip-" + this.viz.id + ")")
                        .attr("class", "d3-line");

                    adding
                        .append("linearGradient")
                        .attr("id", this.viz.id + "-line-style")
                        .attr("gradientUnits", "userSpaceOnUse")
                        .attr("x1", 0)
                        .attr("y1", y(dataRange[0]))
                        .attr("x2", 0)
                        .attr("y2", y(dataRange[1]))
                        .selectAll("stop")
                        .data([
                            { offset: "0%", color: colors(dataRange[0]) },
                            {
                                offset: "20%",
                                color: colors(dataRange[0] + 0.2 * (dataRange[1] - dataRange[0])),
                            },
                            {
                                offset: "40%",
                                color: colors(dataRange[0] + 0.4 * (dataRange[1] - dataRange[0])),
                            },
                            {
                                offset: "60%",
                                color: colors(dataRange[0] + 0.6 * (dataRange[1] - dataRange[0])),
                            },
                            {
                                offset: "80%",
                                color: colors(dataRange[0] + 0.8 * (dataRange[1] - dataRange[0])),
                            },
                            {
                                offset: "100%",
                                color: colors(dataRange[1]),
                            },
                        ])
                        .join((enter) =>
                            enter
                                .append("stop")
                                .attr("offset", (d) => d.offset)
                                .attr("stop-color", (d) => d.color)
                        );

                    return adding;
                });

            const lineFn = d3
                .line()
                .defined((d) => true)
                .x((d) => x(d.time))
                .y((d) => y(d.value))
                .curve(d3.curveBasis);

            line.selectAll(".bkgd-line")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "bkgd-line")
                        .attr("stroke", "#BBBBBB")
                        .attr("stroke-dasharray", "4,4")
                        .attr("fill", "none")
                )
                .attr("d", lineFn(data.data));

            line.selectAll(".data-line")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "data-line")
                        .attr("stroke", "url(#" + this.viz.id + "-line-style)")
                        .attr("stroke-width", "3")
                        .attr("fill", "none")
                )
                .attr("d", lineFn(data.data));

            svg.selectAll(".brush")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "brush")
                        .attr("data-chart", this.viz.id)
                        .call(brush)
                );

            line.selectAll(".circle")
                .data(data.data)
                .join(
                    (enter) =>
                        enter
                            .append("circle")
                            .attr("class", "circle")
                            .attr("r", 2),
                    (updating) => updating,
                    (exiting) => exiting.remove()
                )
                .attr("cx", (d) => x(d.time))
                .attr("cy", (d) => y(d.value));
        },
    },
    methods: {
        onDouble() {
            return this.raiseTimeZoomed(TimeRange.eternity);
        },
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        onRemove() {
            return this.$emit("viz-remove");
        },
    },
    template: `
		<div class="viz">
			<h3 v-if="viz.info.title">{{ viz.info.title }}</h3>
			<div class="btn" @click="onRemove">Remove</div>
			<div class="graph" @dblclick="onDouble"></div>
		</div>
	`,
});
