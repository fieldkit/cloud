import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { QueriedData } from "./viz";

export class Scrubber {
    constructor(public readonly data: QueriedData, public readonly index: number) {}
}

export class ScrubberSet {
    public readonly timeRange: TimeRange;

    constructor(public readonly id: string, public readonly visible: TimeRange, public readonly rows: Scrubber[]) {
        this.timeRange = TimeRange.mergeArrays(rows.map((s) => s.data.timeRange));
    }
}

export const D3MultiScrubber = Vue.extend({
    name: "D3MultiScrubber",
    data() {
        return {};
    },
    props: {
        scrubbers: {
            type: ScrubberSet,
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
            console.log("d3multi graphing (visible)");
            this.refresh();
        },
        data(newValue, oldValue) {
            console.log("d3multi graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        console.log("d3multi mounted");
        this.refresh();
    },
    updated() {
        console.log("d3multi updated");
    },
    methods: {
        raiseTimeZoomed(newTimes) {
            return this.$emit("viz-time-zoomed", newTimes);
        },
        refresh() {
            if (!this.data || this.data.length == 0) {
                return;
            }

            const scrubberSize = 40;
            const layout = new ChartLayout(
                1050,
                scrubberSize * this.scrubbers.rows.length + 30,
                new Margins({ top: 10, bottom: 20, left: 20, right: 20 })
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
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("width", (c) => c.layout.width)
                        .attr("height", (c) => c.layout.height);

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

                    return svg;
                });

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
                    .y1((d) => y(d.value));
            };

            const backgroundAreas = svg
                .selectAll(".background-area")
                .data((d) => d.scrubbers.rows)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("class", "background-area")
                        .attr("fill", "rgb(220, 222, 223)")
                        .attr("d", (scrubber) => area(scrubber)(scrubber.data.sdr.data))
                );

            const foregroundAreas = svg
                .selectAll(".foreground-area")
                .data((d) => d.scrubbers.rows)
                .join((enter) =>
                    enter
                        .append("path")
                        .attr("clip-path", "url(#scrubber-clip-" + this.scrubbers.id + ")")
                        .attr("class", "foreground-area")
                        .attr("fill", "rgb(45, 158, 204)")
                        .attr("d", (scrubber) => area(scrubber)(scrubber.data.sdr.data))
                );

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
                            .attr("stroke", "black");
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

            const raiseZoomed = (newTimes) => this.raiseTimeZoomed(newTimes);

            const clip = svg
                .selectAll(".scrubber-clip")
                .data(charts)
                .join((enter) =>
                    enter
                        .append("defs")
                        .append("clipPath")
                        .attr("id", "scrubber-clip-" + this.scrubbers.id)
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
                        const start = x.invert(selection[0]);
                        const end = x.invert(selection[1]);
                        raiseZoomed(new TimeRange(start, end));
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
