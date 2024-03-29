import _ from "lodash";

import { StationID, VizSensor, SensorSpec, VizConfig, SensorRange, DataRow, SensorDataResponse } from "./api";

export * from "./api";

export class Ids {
    private static c = 0;

    static make(): string {
        return "ids-" + Ids.c++;
    }
}

export type ColorScale = any;

export class Time {
    public static Max = 8640000000000000;
    public static Min = -8640000000000000;
}

export class TimeRange {
    public static eternity = new TimeRange(Time.Min, Time.Max);

    private readonly array: number[];

    constructor(public readonly start: number, public readonly end: number) {
        this.array = [start, end];
    }

    public toArray(): number[] {
        return this.array;
    }

    public isExtreme(): boolean {
        return this.start == Time.Min || this.end == Time.Max;
    }

    public contains(o: TimeRange): boolean {
        return this.start <= o.start && this.end >= o.end;
    }

    public describe(): string[] {
        if (this.isExtreme()) {
            return ["extreme", "extreme"];
        }
        return this.toArray().map((t) => new Date(t).toISOString());
    }

    public static mergeArrays(ranges: number[][]): TimeRange {
        if (ranges.length == 0) {
            return TimeRange.eternity;
        }
        const min = _(ranges)
            .map((r) => r[0])
            .min();
        const max = _(ranges)
            .map((r) => r[1])
            .max();
        if (!min || !max) {
            return TimeRange.eternity;
        }
        return new TimeRange(min, max);
    }

    public static fromArrayIntersections(ranges: number[][]): TimeRange {
        return TimeRange.mergeArrays(ranges);
    }

    public static mergeArraysIgnoreExtreme(ranges: number[][]): TimeRange {
        return TimeRange.mergeArrays(ranges.filter((r) => r[0] != Time.Min && r[1] != Time.Max));
    }

    public static mergeRanges(ranges: TimeRange[]): TimeRange {
        return TimeRange.mergeArrays(ranges.map((r) => r.array));
    }

    public rewind(delta: number): TimeRange {
        return new TimeRange(this.start - delta, this.end - delta);
    }

    public expand(days: number): TimeRange {
        const newStart = new Date(this.start);
        const newEnd = new Date(this.end);
        newStart.setDate(newStart.getDate() - days);
        newEnd.setDate(newEnd.getDate() + days);
        return new TimeRange(newStart.getTime(), newEnd.getTime());
    }
}

export class Sensor {
    constructor(public readonly id: number, public readonly name: string) {}
}

export class SensorParams {
    constructor(public readonly sensors: VizSensor[]) {}

    public get stations(): StationID[] {
        return this.sensors.map((vs) => vs[0]);
    }

    public get station(): StationID {
        return this.stations[0];
    }

    public get sensor(): SensorSpec {
        return this.sensors.map((vs) => vs[1])[0];
    }

    public get sensorAndModules(): SensorSpec[] {
        return this.sensors.map((vs) => vs[1]);
    }

    public get id(): string {
        return ["Z", this.sensors.map((s) => _.flatten(s)).join("-"), "S"].join("~");
    }
}

export class DataQueryParams {
    constructor(public readonly when: TimeRange, public readonly sensors: VizSensor[]) {}

    public sameAs(o: DataQueryParams): boolean {
        return _.isEqual(this, o);
    }

    public queryParams(backend: string | null): URLSearchParams {
        const queryParams = new URLSearchParams();
        queryParams.append("start", this.when.start.toString());
        queryParams.append("end", this.when.end.toString());
        queryParams.append("stations", this.stations.join(","));
        queryParams.append("sensors", this.sensorAndModules.join(","));
        queryParams.append("resolution", "1000");
        queryParams.append("complete", "true");
        if (backend) {
            queryParams.append("backend", backend);
        }
        return queryParams;
    }

    public get sensorParams(): SensorParams {
        return new SensorParams(this.sensors);
    }

    public get stations(): StationID[] {
        return this.sensorParams.stations;
    }

    public get station(): StationID {
        return this.sensorParams.station;
    }

    public get sensor(): SensorSpec {
        return this.sensorParams.sensor;
    }

    public get sensorAndModules(): SensorSpec[] {
        return this.sensorParams.sensorAndModules;
    }

    public remapStationsFromModules(map: { [index: string]: number }): DataQueryParams {
        return new DataQueryParams(
            this.when,
            this.sensors.map((vizSensor) => {
                const vizSensorModule = vizSensor[1][0];
                if (map[vizSensorModule]) {
                    return [map[vizSensorModule], vizSensor[1]];
                }
                return vizSensor;
            })
        );
    }
}

export interface MarginsLike {
    top: number;
    bottom: number;
    left: number;
    right: number;
}

export class Margins implements MarginsLike {
    top: number;
    bottom: number;
    left: number;
    right: number;

    constructor(copy: MarginsLike) {
        this.top = copy.top;
        this.bottom = copy.bottom;
        this.left = copy.left;
        this.right = copy.right;
    }
}

export class ChartLayout {
    constructor(public readonly width: number, public readonly height: number, public readonly margins: Margins) {}
}

export class VizInfo {
    constructor(
        public readonly key: string,
        public readonly colorScale: ColorScale,
        public readonly station: { name: string; location: [number, number] | null },
        public readonly unitOfMeasure: string,
        public readonly firmwareKey: string,
        public readonly name: string,
        public readonly viz: VizConfig[],
        public readonly ranges: SensorRange[],
        public readonly chartLabel: string,
        public readonly axisLabel: string
    ) {}

    public get minimumGap(): number | null {
        if (this.viz.length == 0) {
            return null;
        }
        return this.viz[0].minimumGap || null;
    }

    public get constrainedRanges(): SensorRange[] {
        return this.ranges.filter((r) => r.constrained === true);
    }

    public get label(): string {
        if (this.unitOfMeasure) {
            if (this.unitOfMeasure.indexOf(" ") > 0) {
                return this.name + " (" + this.unitOfMeasure + ")";
            }
            return this.name + " (" + _.capitalize(this.unitOfMeasure) + ")";
        }
        return this.name;
    }

    public get aggregationFunction() {
        return this.firmwareKey == "wh.floodnet.depth" ? _.max : _.mean; // HACK
    }

    public applyCustomFilter(rows: DataRow[]): DataRow[] {
        return rows;
    }
}

export function makeRange(values: number[]): [number, number] {
    const min = _.min(values);
    const max = _.max(values);
    if (min === undefined) throw new Error(`no min: ${values.length}`);
    if (max === undefined) throw new Error(`no max: ${values.length}`);
    if (min === max) {
        console.warn(`range-warning: min == max ${values.length}`);
    }
    return [min, max];
}

export function addDays(original: Date, days: number): Date {
    const temp = original;
    temp.setDate(temp.getDate() + days);
    return temp;
}

export function addSeconds(original: Date, seconds: number): Date {
    const temp = original;
    temp.setSeconds(temp.getSeconds() + seconds);
    return temp;
}

export function truncateTime(original: Date): Date {
    return new Date(original.getFullYear(), original.getMonth(), original.getDate());
}

export class QueriedData {
    empty = true;
    dataRange: number[] = [];
    timeRangeData: number[] = [];
    timeRange: number[] = [];

    constructor(public readonly key: string, public readonly timeRangeQueried: TimeRange, private readonly sdr: SensorDataResponse) {
        if (this.sdr.data.length > 0) {
            const filtered = this.sdr.data.filter((d) => _.isNumber(d.value));
            const values = filtered.map((d) => d.value);
            const times = filtered.map((d) => d.time);

            if (values.length == 0) throw new Error(`empty data ranges`);
            if (times.length == 0) throw new Error(`empty time ranges`);

            this.dataRange = makeRange(values as number[]);
            this.timeRangeData = makeRange(times);

            if (this.timeRangeQueried.isExtreme()) {
                this.timeRange = this.timeRangeData;
            } else {
                this.timeRange = this.timeRangeQueried.toArray();
            }
            this.empty = false;
        }
    }

    get data() {
        return this.sdr.data;
    }

    get bucketSize(): number {
        return this.sdr.bucketSize;
    }

    get dataEnd(): number | null {
        return this.sdr.dataEnd;
    }

    private getAverageTimeBetweenSample(): number | null {
        if (this.sdr.data.length <= 1) {
            return null;
        }

        const deltas = this.sdr.data
            .map((row) => row.time)
            .reduce((diffs, time, index) => {
                if (diffs.length == 0) {
                    return [time];
                }
                const head = diffs.slice(0, diffs.length - 1);
                const diff = time - diffs[diffs.length - 1];
                const extra = [...head, diff];
                if (index == this.sdr.data.length - 1) {
                    return extra;
                }
                return [...extra, time];
            }, []);

        return _.mean(deltas);
    }

    get shouldIgnoreMissing(): boolean {
        const averageMs = this.getAverageTimeBetweenSample();
        // console.log("viz: average-delta", averageMs);
        if (averageMs) {
            return averageMs - 1000 < 5; // Was this.sdr.aggregate.interval <= 60;
        }
        return true;
    }

    public sorted(): QueriedData {
        const sorted = {
            data: _.sortBy(this.sdr.data, (d) => d.time),
            dataEnd: this.sdr.dataEnd,
            bucketSize: this.sdr.bucketSize,
            bucketSamples: this.sdr.bucketSamples,
        };
        return new QueriedData(this.key, this.timeRangeQueried, sorted);
    }

    public removeMalformed(): QueriedData {
        const filtered = {
            data: this.sdr.data.filter((d) => d.sensorId),
            dataEnd: this.sdr.dataEnd,
            bucketSize: this.sdr.bucketSize,
            bucketSamples: this.sdr.bucketSamples,
        };
        // console.log(`viz: malformed`, this.sdr.data.length, filtered.data.length);
        return new QueriedData(this.key, this.timeRangeQueried, filtered);
    }

    public removeDuplicates(): QueriedData {
        const filtered = {
            data: _.sortedUniqBy(this.sdr.data, (d) => d.time),
            dataEnd: this.sdr.dataEnd,
            bucketSize: this.sdr.bucketSize,
            bucketSamples: this.sdr.bucketSamples,
        };
        // console.log(`viz: duplicates`, this.sdr.data.length, filtered.data.length);
        return new QueriedData(this.key, this.timeRangeQueried, filtered);
    }
}

export class DataSetSeries {
    constructor(public readonly vizSensor: VizSensor, public graphing: QueriedData | null = null, public all: QueriedData | null = null) {}

    public bookmark(): VizSensor {
        return this.vizSensor;
    }

    public get stationId(): StationID {
        return this.vizSensor[0];
    }

    public get sensorAndModule() {
        return this.vizSensor[1];
    }

    public shouldConstrainBy(dataRange: [number, number], range: [number, number]): boolean {
        if (this.graphing == null) {
            // console.log("viz: constrain:no-graphing");
            return false;
        }
        if (dataRange[1] > range[1]) {
            // console.log("viz: constrain:nope", this.graphing.dataRange[1], range[1]);
            return false;
        }
        return true;
    }
}

export class SeriesData {
    constructor(
        public readonly key: string,
        public readonly visible: TimeRange,
        public readonly ds: DataSetSeries,
        public readonly queried: QueriedData,
        public readonly vizInfo: VizInfo
    ) {}
}

export class ExploreContext {
    constructor(public readonly project: number | null = null, public readonly map = false) {}
}
