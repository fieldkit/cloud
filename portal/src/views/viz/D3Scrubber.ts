import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData } from "./viz";

export const D3Scrubber = Vue.extend({
    name: "D3Scrubber",
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
            if (this.viz.all && !this.viz.all.empty) {
                return this.viz.all;
            }
            return null;
        },
        visible(): TimeRange {
            return this.viz.visible;
        },
    },
    watch: {
        visible(newValue, oldValue) {
            this.viz.log("graphing (visible)");
            this.refresh();
        },
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
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        refresh() {
            if (!this.data) {
                return;
            }
            const layout = new ChartLayout(1050, 45, new Margins({ top: 0, bottom: 0, left: 0, right: 0 }));
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
                .domain(timeRange)
                .range([layout.margins.left, layout.width - layout.margins.left - layout.margins.right]);

            const y = d3
                .scaleLinear()
                .domain(dataRange)
                .range([layout.height - layout.margins.bottom, layout.margins.top]);

            const area = d3
                .area()
                .x((d) => x(d.time))
                .y0(layout.height - layout.margins.bottom)
                .y1((d) => y(d.value));

            const svg = d3
                .select(this.$el)
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const svg = enter
                        .append("svg")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("width", (c) => c.layout.width)
                        .attr("height", (c) => c.layout.height);

                    const defs = svg.append("defs");

                    svg.append("rect")
                        .attr("x", layout.margins.left)
                        .attr("y", 0)
                        .attr("width", layout.width)
                        .attr("height", layout.height)
                        .attr("fill", "#f4f5f7");

                    const blurFilter = defs
                        .append("filter")
                        .attr("id", "dropshadow-" + this.viz.id)
                        .attr("height", "150%")
                        .attr("width", "150%");
                    blurFilter
                        .append("feGaussianBlur")
                        .attr("in", "SourceAlpha")
                        .attr("stdDeviation", 2);
                    blurFilter
                        .append("feOffset")
                        .attr("dx", 2)
                        .attr("dy", 2)
                        .attr("result", "offsetblur");
                    blurFilter
                        .append("feComponentTransfer")
                        .append("feFuncA")
                        .attr("type", "linear")
                        .attr("slope", 0.5);
                    const merge = blurFilter.append("feMerge");
                    merge.append("feMergeNode");
                    merge.append("feMergeNode").attr("in", "SourceGraphic");

                    return svg;
                });

            const backgroundArea = svg
                .selectAll(".background-area")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "background-area")
                        .attr("fill", "rgb(220, 222, 223)")
                )
                .data([this.data.data])
                .attr("d", area);

            const foregroundArea = svg
                .selectAll(".foreground-area")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("clip-path", "url(#scrubber-clip-" + this.viz.id + ")")
                        .attr("class", "foreground-area")
                        .attr("fill", "rgb(45, 158, 204)")
                )
                .data([this.data.data])
                .attr("d", area);

            const handles = (g, selection) =>
                g
                    .selectAll(".handle-custom")
                    .data([{ type: "w" }, { type: "e" }])
                    .join((enter) => {
                        const handle = enter.append("g").attr("class", "handle-custom");
                        handle
                            .append("line")
                            .attr("x1", 0)
                            .attr("x2", 0)
                            .attr("y1", -(layout.height - layout.margins.top - layout.margins.bottom) / 2)
                            .attr("y2", (layout.height - layout.margins.top - layout.margins.bottom) / 2)
                            .attr("stroke-width", "2")
                            .attr("stroke", "#4f4f4f");
                        handle
                            .append("circle")
                            .attr("filter", "url(#dropshadow-" + this.viz.id + ")")
                            .attr("fill", "white")
                            .attr("cursor", "ew-resize")
                            .attr("r", 9);
                        return handle;
                    })
                    .attr("display", selection === null ? "none" : null)
                    .attr(
                        "transform",
                        selection === null
                            ? null
                            : (d, i) => `translate(${selection[i]},${(layout.height + layout.margins.top - layout.margins.bottom) / 2})`
                    );

            const raiseZoomed = (newTimes) => this.raiseTimeZoomed(newTimes);

            const clip = svg
                .selectAll(".scrubber-clip")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("defs")
                        .append("clipPath")
                        .attr("id", "scrubber-clip-" + this.viz.id)
                        .attr("class", "scrubber-clip")
                        .append("rect")
                        .attr("width", layout.width)
                        .attr("height", layout.height)
                        .attr("x", 0)
                        .attr("y", 0)
                );

            const brush = d3
                .brushX()
                .extent([
                    [layout.margins.left, layout.margins.top],
                    [layout.width - layout.margins.left - layout.margins.right, layout.height - layout.margins.bottom],
                ])
                .on("start brush end", function(this: any) {
                    const selection = d3.event.selection;
                    if (selection !== null) {
                        const sx = selection.map(x.invert);
                        clip.attr("x", selection[0]).attr("width", selection[1] - selection[0]);
                    }

                    d3.select(this).call(handles, selection);

                    if (d3.event.type == "end" && d3.event.sourceEvent) {
                        if (selection != null) {
                            const start = x.invert(selection[0]);
                            const end = x.invert(selection[1]);
                            raiseZoomed(new TimeRange(start.getTime(), end.getTime()));
                        } else {
                            raiseZoomed(TimeRange.eternity);
                        }
                    }
                });

            const visible = () => {
                if (this.viz.visible.isExtreme()) {
                    return this.data.timeRange;
                }
                return this.viz.visible.toArray();
            };

            const brushTop = svg
                .selectAll(".brush-container")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("g")
                        .attr("class", "brush-container")
                        .call(brush)
                )
                .call(brush.move, visible().map(x));
        },
    },
    template: `<div class="viz scrubber"></div>`,
});
