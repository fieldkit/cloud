import * as d3 from "d3";

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
