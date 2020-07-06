import _ from "lodash";

export class Ids {
    private static c = 0;

    static make(): string {
        return "ids-" + Ids.c++;
    }
}

export class Time {
    constructor(public readonly d: Date) {}

    public static Max = new Time(new Date(8640000000000000));
    public static Min = new Time(new Date(-8640000000000000));

    public getTime() {
        return this.d.getTime();
    }
}

export class TimeRange {
    constructor(public readonly start: Time, public readonly end: Time) {}

    public static eternity = new TimeRange(Time.Min, Time.Max);

    public toArray(): number[] {
        return [this.start.d.getTime(), this.end.d.getTime()];
    }

    public isExtreme(): boolean {
        return this.start == Time.Min || this.end == Time.Max;
    }
}

export class Sensor {
    constructor(public readonly id: number, public readonly name: string) {}
}

export type Stations = number[];
export type Sensors = number[];

export class DataQueryParams {
    constructor(public readonly when: TimeRange, public readonly stations: Stations, public readonly sensors: Sensors) {}

    public sameAs(o: DataQueryParams): boolean {
        return _.isEqual(this, o);
    }

    public queryString(): string {
        const queryParams = new URLSearchParams();
        queryParams.append("start", this.when.start.getTime().toString());
        queryParams.append("end", this.when.end.getTime().toString());
        queryParams.append("stations", this.stations.join(","));
        queryParams.append("sensors", this.sensors.join(","));
        return queryParams.toString();
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
