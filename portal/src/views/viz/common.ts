import _ from "lodash";

export class Ids {
    private static c = 0;

    static make(): string {
        return "ids-" + Ids.c++;
    }
}

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
        return new TimeRange(min, max);
    }

    public static mergeRanges(ranges: TimeRange[]): TimeRange {
        return TimeRange.mergeArrays(ranges.map((r) => r.array));
    }

    public contains(o: TimeRange): boolean {
        return this.start <= o.start && this.end >= o.end;
    }
}

export class Sensor {
    constructor(public readonly id: number, public readonly name: string) {}
}

export type Stations = number[];
export type Sensors = number[];

export class SensorParams {
    constructor(public readonly sensors: Sensors, public readonly stations: Stations) {}

    public get id(): string {
        if (this.stations.length == 0) {
            return ["Z", this.sensors.join("-"), "S", "ALL"].join("~");
        }
        return ["Z", this.sensors.join("-"), "S", this.stations.join("-")].join("~");
    }
}

export class DataQueryParams {
    constructor(public readonly when: TimeRange, public readonly stations: Stations, public readonly sensors: Sensors) {}

    public sameAs(o: DataQueryParams): boolean {
        return _.isEqual(this, o);
    }

    public queryParams(): URLSearchParams {
        const queryParams = new URLSearchParams();
        queryParams.append("start", this.when.start.toString());
        queryParams.append("end", this.when.end.toString());
        queryParams.append("stations", this.stations.join(","));
        queryParams.append("sensors", this.sensors.join(","));
        queryParams.append("resolution", "1000");
        return queryParams;
    }

    public get sensorParams(): SensorParams {
        return new SensorParams(this.sensors, this.stations);
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
