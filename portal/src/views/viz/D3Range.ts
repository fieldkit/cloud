import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";
import {appendUnitOfMeasureLabel} from '@/views/viz/d3-helpers';

export const D3Range = Vue.extend({
    name: "D3Range",
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
        data(): QueriedData | null {
            if (this.viz.graphing) {
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

            const NumberOfBins = 64;
            const thresholds = d3.range(timeRange[0], timeRange[1], (timeRange[1] - timeRange[0]) / NumberOfBins);

            const x = d3
                .scaleTime()
                .domain(timeRange)
                .range([layout.margins.left, layout.width - (layout.margins.right + layout.margins.left)]);

            const histogram = d3
                .histogram()
                .value((d) => d.time)
                .domain(x.domain())
                .thresholds(thresholds);

            const y = d3
                .scaleLinear()
                .domain(dataRange)
                .range([layout.height - (layout.margins.bottom + layout.margins.top), layout.margins.top]);

            const bins = histogram(this.data.data.filter((d) => d.value));

            this.viz.log("bins", bins);

            const charts = [
                {
                    layout: layout,
                    bins: bins,
                },
            ];

            const xAxis = d3.axisBottom(x).ticks(10);
            const yAxis = d3.axisLeft(y).ticks(6);

            const svg = d3
                .select(this.$el)
                .select(".chart")
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const adding = enter
                        .append("svg")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("viewBox", "0 0 " + layout.width + " " + layout.height);

                    adding.append("defs");

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

            svg.select("#uom").remove();
            svg.append("text")
                .attr("id", "uom")
                .attr("text-anchor", "middle")
                .attr("transform", "rotate(-90)")
                .attr("fill", "#2C3E50")
                .attr("y", 17)
                .attr("x", vizInfo.unitOfMeasure.length / 2 - (layout.height - (layout.margins.bottom + layout.margins.top)) / 2)
                .text(vizInfo.unitOfMeasure);

            const defs = svg
                .selectAll("defs")
                .data(charts)
                .join((enter) => enter.append("defs"));

            const colors = vizInfo.colorScale;
            const stops = [
                {
                    offset: "0%",
                    color: (bin) => colors(d3.min(bin, (d) => d.value)),
                },
                {
                    offset: "50%",
                    color: (bin) => colors(d3.median(bin, (d) => d.value)),
                },
                {
                    offset: "100%",
                    color: (bin) => colors(d3.max(bin, (d) => d.value)),
                },
            ];

            const styles = defs
                .selectAll(".bar-style")
                .data(bins)
                .join((enter, ...args) => {
                    const style = enter
                        .append("linearGradient")
                        .attr("id", (d, i) => "range-bar-style-" + i)
                        .attr("x1", "0%")
                        .attr("x2", "0%")
                        .attr("y1", "0%")
                        .attr("y2", "100%");

                    style
                        .selectAll("stop")
                        .data(stops)
                        .join((enter) => enter.append("stop").attr("offset", (d) => d.offset))
                        .attr("stop-color", function(this: any, stop: any) {
                            return stop.color(d3.select(this.parentNode).datum());
                        })
                        .style("stop-opacity", 1);

                    return style;
                });

            svg.selectAll(".range-bar")
                .data(bins)
                .join((enter) => enter.append("rect").attr("class", "range-bar"))
                .attr("transform", (d) => {
                    const mx = d3.max(d, (b) => b.value);
                    return mx || mx === 0 ? "translate(" + x(d.x0) + "," + y(mx) + ")" : "translate(0,0)";
                })
                .attr("width", (d) => {
                    return x(d.x1) - x(d.x0) - 1;
                })
                .style("fill", (d, i) => {
                    return d.length > 0 ? "url(#range-bar-style-" + i + ")" : "none";
                })
                .attr("height", (d) => {
                    const MinimumHeight = 2; // NOTE Why 2?
                    const extent = d3.extent(d, (b) => b.value);
                    const height = y(extent[1]) - y(extent[0]);
                    const min = d.length > 0 ? -MinimumHeight : 0;
                    return -(height ? height : min);
                });

            appendUnitOfMeasureLabel(svg, vizInfo.unitOfMeasure, layout);
        },
    },
    template: `<div class="viz histogram"><div class="chart"></div></div>`,
});
