import * as d3 from "d3";
import { ChartLayout } from "@/views/viz/common";

export type ColorScale = any;

export interface SensorDetails {
    key: string;
    internal: boolean;
    fullKey: string;
    firmwareKey: string;
    unitOfMeasure: string;
    ranges: { minimum: number; maximum: number }[];
}

export function createSensorColorScale(sensor: SensorDetails | null): ColorScale {
    if (sensor == null || sensor.ranges.length == 0) {
        return d3
            .scaleSequential()
            .domain([0, 1])
            .interpolator(() => "#000000");
    }

    const range = sensor.ranges[0];

    return d3
        .scaleSequential()
        .domain([range.minimum, range.maximum])
        .interpolator(d3.interpolatePlasma);
}

export function appendUnitOfMeasureLabel(svg: d3, unitOfMeasure: string, layout: ChartLayout): void {
    svg.select("#uom").remove();
    svg.append("text")
        .attr("id", "uom")
        .attr("text-anchor", "middle")
        .attr("transform", "rotate(-90)")
        .attr("fill", "#7F7F7F")
        .style("font-size", "10px")
        .attr("y", 19)
        .attr("x", unitOfMeasure.length / 2 - (layout.height - (layout.margins.bottom + layout.margins.top)) / 2)
        .text(unitOfMeasure);
}

export function appendTimeLabel(svg: d3, layout: ChartLayout): void {
    const text = 'Time (Days)'
    svg.select("#time").remove();
    svg.append("text")
        .attr("id", "time")
        .attr("text-anchor", "middle")
        .attr("fill", "#7F7F7F")
        .style("font-size", "10px")
        .attr("y", layout.height - 17)
        .attr("x", (layout.width - 0) / 2)
        .text(text);
}
