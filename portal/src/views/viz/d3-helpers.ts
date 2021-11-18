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

export function appendYAxisLabel(svg: d3, unitOfMeasure: string, layout: ChartLayout): void {
    const svgId = "d3-uom";
    svg.select("#" + svgId).remove();
    svg.append("text")
        .attr("id", svgId)
        .attr("text-anchor", "middle")
        .attr("transform", "rotate(-90)")
        .attr("fill", "#7F7F7F")
        .style("font-size", "10px")
        .attr("y", 19)
        .attr("x", unitOfMeasure.length / 2 - (layout.height - (layout.margins.bottom + layout.margins.top)) / 2)
        .text(unitOfMeasure);
}

export function appendXAxisLabel(svg: d3, layout: ChartLayout): void {
    const svgId = "d3-time";
    svg.select("#" + svgId).remove();
    svg.append("text")
        .attr("id", svgId)
        .attr("text-anchor", "middle")
        .attr("fill", "#7F7F7F")
        .style("font-size", "10px")
        .attr("y", layout.height - 20)
        .attr("x", layout.width / 2)
        .text("Time");
}

export function getMaxDigitsForData(dataRange: number[]) {
    return Math.ceil(Math.log10(Math.floor(Math.abs(dataRange[1])) + 1));
}
