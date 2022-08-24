import _ from "lodash";
import type { VizSensor, DataRow, SensorDataResponse, ModuleSensorMeta, SensorsResponse } from "../common";
import { getString, getSeriesThresholds, getAxisLabel } from "./customizations";
import {
    SeriesData,
    DataSetSeries,
    TimeRange,
    VizInfo,
    QueriedData
} from "../common";

export {
    getAxisLabel,
    getString,
    getSeriesThresholds,
    SeriesData,
    TimeRange,
    DataSetSeries,
    VizInfo,
    QueriedData,
    DataRow,
    SensorDataResponse,
    SensorsResponse,
    ModuleSensorMeta,
    VizSensor,
};

export type MapFunction<T> = (series: SeriesData, i: number) => T;

export class ChartSettings {
    constructor(
        public readonly timeRange = TimeRange.eternity,
        public readonly w: number = 0,
        public readonly h: number = 0,
        public readonly auto = false,
        public readonly tiny = false,
        public readonly mobile = false
    ) {}

    public apply(spec: unknown): unknown {
        if (this.auto) {
            const autoSize = {
                autosize: {
                    type: "fit",
                    contains: "padding",
                    resize: this.tiny,
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

    public static DefaultDesktop = new ChartSettings(TimeRange.eternity, 0, 0, true, false, false);
    public static DefaultMobile = new ChartSettings(TimeRange.eternity, 0, 0, true, false, true);
    public static Container = new ChartSettings(TimeRange.eternity, 0, 0, false, false, false);
    public static Tiny = new ChartSettings(TimeRange.eternity, 0, 0, true, true, false);
}
