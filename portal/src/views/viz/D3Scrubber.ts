import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { TimeRange, Margins, ChartLayout } from "./common";
import { QueriedData, Scrubbers, TimeZoom } from "./viz";

export const D3Scrubber = Vue.extend({
    name: "D3Scrubber",
    data() {
        return {};
    },
    props: {
        scrubbers: {
            type: Scrubbers,
            required: true,
        },
    },
    computed: {
        data(): QueriedData[] {
            return this.scrubbers.rows.map((s) => s.data);
        },
        visible(): TimeRange {
            return this.scrubbers.visible;
        },
    },
    watch: {
        visible(newValue, oldValue) {
            // console.log("d3scrubber graphing (visible)");
            this.refresh();
        },
        data(newValue, oldValue) {
            // console.log("d3scrubber graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        // console.log("d3scrubber mounted");
        this.refresh();
    },
    updated() {
        // console.log("d3scrubber updated");
    },
    methods: {
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", new TimeZoom(null, newTimes));
        },
        refresh() {
            if (!this.data || this.data.length == 0) {
                return;
            }

            const scrubberSize = this.scrubbers.rows.length > 1 ? 30 : 40;
            const layout = new ChartLayout(
                1050,
                scrubberSize * this.scrubbers.rows.length,
                new Margins({ top: 0, bottom: 0, left: 0, right: 0 })
            );
            const timeRange = this.scrubbers.timeRange;
            const charts = [
                {
                    layout: layout,
                    scrubbers: this.scrubbers,
                },
            ];

            const x = d3
                .scaleTime()
                .domain(timeRange.toArray())
                .range([layout.margins.left, layout.width - layout.margins.left - layout.margins.right]);

            const svg = d3
                .select(this.$el)
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const svg = enter
                        .append("svg")
                        .attr("class", "svg-container")
                        .attr("preserveAspectRatio", "xMidYMid meet");

                    const defs = svg.append("defs");

                    const blurFilter = defs
                        .append("filter")
                        .attr("id", "dropshadow-" + this.scrubbers.id)
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

                    svg.append("g").attr("class", "background-areas");

                    svg.append("g").attr("class", "foreground-areas");

                    return svg;
                })
                .attr("viewBox", "0 0 " + layout.width + " " + layout.height);

            const area = (scrubber) => {
                const top = layout.margins.top + scrubberSize * scrubber.index;
                const bottom = top + scrubberSize;
                const y = d3
                    .scaleLinear()
                    .domain(scrubber.data.dataRange)
                    .range([bottom, top]);

                return d3
                    .area()
                    .x((d) => x(d.time))
                    .y0(bottom)
                    .y1((d) => {
                        if (d.value) {
                            return y(d.value);
                        }
                        return bottom;
                    });
            };

            const renderedArea = (scrubber) => area(scrubber)(scrubber.data.sdr.data);

            const backgroundAreas = svg
                .select(".background-areas")
                .selectAll(".background-area")
                .data((d) => d.scrubbers.rows)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "background-area")
                        .attr("fill", "rgb(220, 222, 223)")
                )
                .attr("d", renderedArea);

            const foregroundAreas = svg
                .select(".foreground-areas")
                .selectAll(".foreground-area")
                .data((d) => d.scrubbers.rows)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "foreground-area")
                        .attr("fill", "rgb(45, 158, 204)")
                )
                .attr("clip-path", "url(#scrubber-clip-" + this.scrubbers.id + ")")
                .attr("d", renderedArea);

            const handles = (g, selection) => {
                g.selectAll(".handle-custom")
                    .data([{ type: "w" }, { type: "e" }])
                    .join((enter) => {
                        const handle = enter.append("g").attr("class", "handle-custom");
                        handle.append("line").attr("stroke", "black");
                        handle
                            .append("circle")
                            .attr("filter", "url(#dropshadow-" + this.scrubbers.id + ")")
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

                g.selectAll("line")
                    .attr("x1", 0)
                    .attr("x2", 0)
                    .attr("y1", -(layout.height - layout.margins.top - layout.margins.bottom) / 2)
                    .attr("y2", (layout.height - layout.margins.top - layout.margins.bottom) / 2);
            };

            const raiseZoomed = (newTimes) => this.raiseTimeZoomed(newTimes);

            const clip = svg
                .selectAll(".scrubber-clip")
                .data(charts)
                .join((enter) => {
                    const clip = enter
                        .select("defs")
                        .append("clipPath")
                        .attr("class", "scrubber-clip");

                    clip.append("rect")
                        .attr("x", 0)
                        .attr("y", 0);

                    return clip;
                })
                .attr("id", "scrubber-clip-" + this.scrubbers.id)
                .select("rect")
                .attr("width", layout.width)
                .attr("height", layout.height);

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
                        const start = x.invert(selection[0]);
                        const end = x.invert(selection[1]);
                        raiseZoomed(new TimeRange(start.getTime(), end.getTime()));
                    }
                });

            const visible = () => {
                if (this.scrubbers.visible.isExtreme()) {
                    return this.scrubbers.timeRange.toArray();
                }
                return this.scrubbers.visible.toArray();
            };

            const brushTop = svg
                .selectAll(".brush-container")
                .data(charts)
                .join((enter) => enter.append("g").attr("class", "brush-container"))
                .call(brush)
                .call(brush.move, visible().map(x));
        },
    },
    template: `<div class="viz scrubber"></div>`,
});
