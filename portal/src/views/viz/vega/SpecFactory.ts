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

interface WidthAndHeight {
    w: number;
    h: number
};

export class ChartSettings {
    constructor(
        public readonly timeRange = TimeRange.eternity,
        public readonly estimated: WidthAndHeight | undefined = undefined,
        public readonly size: WidthAndHeight  = { w: 0, h:0},
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

        if (this.size.w > 0 && this.size.h > 0) {
            const fixedSize = {
                width: this.size.w,
                height: this.size.h,
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

    public static makeDefaultDesktop(): ChartSettings {
        const estimated = {
            w: window.innerWidth,
            h: window.innerHeight,
        };
        return new ChartSettings(TimeRange.eternity, estimated , { w: 0, h: 0 }, true, false, false);
    }

    public static DefaultMobile = new ChartSettings(TimeRange.eternity, undefined, { w: 0, h: 0 }, true, false, true);
    public static Container = new ChartSettings(TimeRange.eternity, undefined, { w: 0, h: 0 }, false, false, false);
    public static Tiny = new ChartSettings(TimeRange.eternity, undefined, { w: 0, h: 0 }, true, true, false);
}
