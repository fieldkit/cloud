import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace, FastTime, TimeZoom } from "./viz";

export const D3TimeSeriesGraph = Vue.extend({
    name: "D3TimeSeriesGraph",
    components: {},
    data() {
        return {};
    },
    props: {
        viz: {
            type: Graph,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    computed: {
        data(): QueriedData {
            if (this.viz.graphing && !this.viz.graphing.empty) {
                return this.viz.graphing;
            }
            return null;
        },
    },
    watch: {
        data(newValue, oldValue) {
            this.viz.log("graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        this.viz.log("mounted");
        this.refresh();
    },
    updated() {
        this.viz.log("updated");
    },
    methods: {
        onDouble() {
            return this.raiseTimeZoomed(new TimeZoom(FastTime.All, null));
        },
        raiseTimeZoomed(newTimes: TimeZoom) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        refresh() {
            if (!this.data) {
                return;
            }

            const vizInfo = this.workspace.vizInfo(this.viz);
            const layout = new ChartLayout(1050, 340, new Margins({ top: 5, bottom: 50, left: 50, right: 0 }));
            const data = this.data;
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
                .select(".chart")
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const svg = enter
                        .append("svg")
                        .attr("class", "svg-container")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("viewBox", "0 0 " + layout.width + " " + layout.height);

                    return svg;
                });

            svg.selectAll(".x-axis")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "svg-container-responsive x-axis")
                        .attr("transform", "translate(" + 0 + "," + (layout.height - (layout.margins.bottom + layout.margins.top)) + ")")
                )
                .call(xAxis);

            svg.selectAll(".y-axis")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "svg-container-responsive y-axis")
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
                    const newRange = new TimeRange(start.getTime(), end.getTime());
                    this.raiseTimeZoomed(new TimeZoom(null, newRange));
                });

            const colors = vizInfo.colorScale;
            const distance = dataRange[1] - dataRange[0];
            const stops = [
                { offset: "0%", color: colors(dataRange[0]) },
                { offset: "20%", color: colors(dataRange[0] + 0.2 * distance) },
                { offset: "40%", color: colors(dataRange[0] + 0.4 * distance) },
                { offset: "60%", color: colors(dataRange[0] + 0.6 * distance) },
                { offset: "80%", color: colors(dataRange[0] + 0.8 * distance) },
                { offset: "100%", color: colors(dataRange[1]) },
            ];

            const line = svg
                .selectAll(".d3-line")
                .data(charts)
                .join((enter) => {
                    const adding = enter.append("g").attr("class", "svg-container-responsive d3-line");

                    adding
                        .append("linearGradient")
                        .attr("id", this.viz.id + "-line-style")
                        .attr("gradientUnits", "userSpaceOnUse")
                        .attr("x1", 0)
                        .attr("y1", y(dataRange[0]))
                        .attr("x2", 0)
                        .attr("y2", y(dataRange[1]))
                        .selectAll("stop")
                        .data(stops)
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
    template: `<div class="viz time-series-graph"><div class="chart" @dblclick="onDouble"></div></div>`,
});
