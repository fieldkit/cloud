import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";

export const D3Histogram = Vue.extend({
    name: "D3Histogram",
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

            const svg = d3
                .select(this.$el)
                .select(".chart")
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const adding = enter
                        .append("svg")
                        .attr("class", "svg-container")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("viewBox", "0 0 " + layout.width + " " + layout.height);

                    return adding;
                });

            const NumberOfBins = 16;
            const thresholds = d3.range(dataRange[0], dataRange[1], (dataRange[1] - dataRange[0]) / NumberOfBins);

            const x = d3
                .scaleLinear()
                .domain(data.dataRange)
                .range([layout.margins.left, layout.width - (layout.margins.right + layout.margins.left)]);

            const histogram = d3
                .histogram()
                .value((d) => d.value)
                .domain(x.domain())
                .thresholds(thresholds);

            const bins = histogram(this.data.data);

            this.viz.log("bins", bins);

            const y = d3
                .scaleLinear()
                .domain([0, d3.max(bins, (d) => d.length)])
                .range([layout.height - (layout.margins.bottom + layout.margins.top), layout.margins.top]);

            const xAxis = d3.axisBottom(x).ticks(NumberOfBins);
            const yAxis = d3.axisLeft(y).ticks(10);

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

            const colors = d3
                .scaleSequential()
                .domain([0, 1])
                .interpolator(() => {
                    return "#000000";
                });

            // NOTE Why 800?
            const MinimumBarWidth = 800 / NumberOfBins;
            svg.selectAll(".histogram-bar")
                .data(bins)
                .join((enter) => {
                    return enter.append("rect").attr("class", "histogram-bar");
                })
                .attr("transform", (d) => {
                    return "translate(" + x(d.x0) + "," + y(d.length) + ")";
                })
                .attr("width", (d) => {
                    const w = x(d.x1) - x(d.x0) - 1;
                    return w <= 0 ? MinimumBarWidth : w;
                })
                .style("fill", (d) => vizInfo.colorScale)
                .attr("height", (d) => {
                    return d.length == 0 ? 0 : layout.height - y(d.length) - layout.margins.bottom - layout.margins.top;
                });
        },
    },
    template: `<div class="viz histogram"><div class="chart"></div></div>`,
});
