import { getString, getSeriesThresholds } from "./customizations";
import { SeriesData, TimeRange } from "../common";
import _ from "lodash";

export { getString, getSeriesThresholds, SeriesData };

export type MapFunction<T> = (series: SeriesData, i: number) => T;

export class ChartSettings {
    constructor(
        public readonly timeRange = TimeRange.eternity,
        public readonly w: number = 0,
        public readonly h: number = 0,
        public readonly auto = false
    ) {}

    public apply(spec: unknown): unknown {
        if (this.auto) {
            const autoSize = {
                autosize: {
                    type: "fit",
                    contains: "padding",
                },
            };
            return _.extend(spec, autoSize);
        }

        if (this.w > 0 && this.h > 0) {
            const fixedSize = {
                width: this.w,
                height: this.h,
                autosize: "pad",
            };
            return _.extend(spec, fixedSize);
        }

        const containerSize = {
            width: "container",
            height: "container",
        };
        return _.extend(spec, containerSize);
    }

    public static Container = new ChartSettings(TimeRange.eternity, 0, 0, false);
    public static Auto = new ChartSettings(TimeRange.eternity, 0, 0, true);
}
