import _ from "lodash";
// import Vue, { VNode } from "vue";
import { SensorId, SensorsResponse, SensorDataResponse, ModuleSensor } from "./api";
import { Ids, TimeRange, Stations, Sensors, DataQueryParams } from "./common";
import { DisplayStation } from "@/store/modules/stations";
import FKApi from "@/api/api";

export class SensorMeta {
    constructor(private readonly meta: SensorId, private sensor: ModuleSensor) {}

    // TODO: sensor scales

    public get id(): number {
        return this.meta.id;
    }

    public get fullKey(): string {
        return this.sensor.fullKey;
    }

    public get name(): string {
        return this.sensor.fullKey;
    }
}

export class StationMeta {
    constructor(private readonly meta: SensorsResponse, private station: DisplayStation) {}

    get sensors(): SensorMeta[] {
        const sensorsByKey = _.keyBy(this.meta.sensors, (s) => s.key);
        return _(this.station.configurations.all)
            .take(1)
            .map((cfg) => cfg.modules)
            .flatten()
            .map((m) => m.sensors)
            .flatten()
            .map((s) => new SensorMeta(sensorsByKey[s.fullKey], s))
            .filter((s) => {
                if (!s.meta) {
                    console.log("no meta:", s.sensor.fullKey);
                    return false;
                }
                return true;
            })
            .value();
    }
}

export class QueriedData {
    empty = true;
    dataRange: number[] = [];
    timeRangeData: number[] = [];
    timeRange: number[] = [];

    constructor(public readonly timeRangeQueried: TimeRange, private readonly sdr: SensorDataResponse) {
        if (this.sdr.data.length > 0) {
            const values = _(this.sdr.data).map((d) => d.value);
            const times = _(this.sdr.data).map((d) => d.time);
            this.dataRange = [values.min(), values.max()];
            this.timeRangeData = [times.min(), times.max()];

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
}

export class VizInfo {
    constructor(public readonly title: string) {}

    public static fromSensor(sm: SensorMeta): VizInfo {
        return new VizInfo("Sensor: " + sm.fullKey);
    }
}

export class Viz {
    public readonly id = Ids.make();

    constructor(public readonly info: VizInfo) {}

    public get params(): DataQueryParams[] {
        return [];
    }

    public zoomed(times: TimeRange) {
        return undefined;
    }

    public needs(viz: Viz) {
        return this.id === viz.id;
    }

    public log(...args: any[]) {
        console.log(...["viz:", this.constructor.name, ...args]);
    }
}

export class QueriesSensorData extends Viz {
    public data: QueriedData | null = null;

    constructor(
        public readonly info: VizInfo,
        public when: TimeRange,
        public readonly stations: Stations,
        public readonly sensors: Sensors
    ) {
        super(info);
    }

    public get params(): DataQueryParams[] {
        return [new DataQueryParams(this.when, this.stations, this.sensors)];
    }
}

export class Scrubber extends QueriesSensorData {
    constructor(public readonly when: TimeRange, public readonly graph: Graph) {
        super(null, when, graph.stations, graph.sensors);
    }

    public zoomed(times: TimeRange) {
        this.graph.zoomed(times);
    }

    public needs(viz: Viz) {
        return super.needs(viz) || this.graph.needs(viz);
    }
}

export class Graph extends QueriesSensorData {
    constructor(
        public readonly info: VizInfo,
        public when: TimeRange,
        public readonly stations: Stations,
        public readonly sensors: Sensors
    ) {
        super(info, when, stations, sensors);
    }

    public zoomed(times: TimeRange) {
        this.when = times;
    }

    public makeScrubber(): Scrubber {
        return new Scrubber(TimeRange.eternity, this);
    }
}

export class TimeSeriesGraph extends Graph {}

export class RangeGraph extends Graph {}

export class DistributionGraph extends Graph {}

export class Map extends Viz {}

export class Group {
    public vizes: Viz[] = [];

    public add(viz: Viz) {
        this.vizes.push(viz);
        return this;
    }

    public remove(removing: Viz) {
        this.vizes = _.remove(this.vizes, (viz) => {
            return !viz.needs(removing);
        });
        return this;
    }

    public get empty() {
        return this.vizes.length == 0;
    }
}

export class Querier {
    private cache: { [index: string]: QueriedData } = {};

    public query(params: DataQueryParams): Promise<QueriedData> {
        if (!params) {
            throw new Error("no params");
        }
        const key = params.queryString();
        if (this.cache[key]) {
            return Promise.resolve(this.cache[key]);
        }
        return new FKApi().sensorData(params).then((sdr: SensorDataResponse) => {
            const queried = new QueriedData(params.when, sdr);
            this.cache[key] = queried;
            return queried;
        });
    }
}

export class TreeOption {
    constructor(
        public readonly id: string,
        public readonly label: string,
        public readonly sensor: SensorMeta,
        public readonly stations: Stations,
        public readonly children: TreeOption[] | undefined = undefined
    ) {}
}

export class StationSensors {
    constructor(public readonly sensors: Sensors, public readonly stations: Stations) {}
}

export class Workspace {
    public readonly querier = new Querier();
    public readonly stations: StationMeta[] = [];
    public groups: Group[] = [];

    constructor(private readonly meta: SensorsResponse) {
        // ok
    }

    public addSensor(sensor: SensorMeta, stations: Stations): void {
        const info = VizInfo.fromSensor(sensor);
        const graph = new TimeSeriesGraph(info, TimeRange.eternity, stations, [sensor.id]);
        const group = new Group();
        group.add(graph);
        group.add(graph.makeScrubber());
        this.groups.push(group);
    }

    public addStation(station: DisplayStation): boolean {
        const stationMeta = new StationMeta(this.meta, station);
        const stationSensors = stationMeta.sensors;
        console.log("sensors:", stationSensors);
        if (stationSensors.length == 0) {
            return false;
        }
        const sensor = stationSensors[0];
        this.addSensor(sensor, [station.id]);
        this.stations.push(stationMeta);
        return true;
    }

    private get allVizes(): Viz[] {
        return _(this.groups)
            .map((g) => g.vizes)
            .flatten()
            .value();
    }

    public query(): Promise<any> {
        const vizToParams = _(this.allVizes)
            .map((viz: Viz) =>
                viz.params.map((params: DataQueryParams) => {
                    return {
                        viz: viz,
                        params: params,
                    };
                })
            )
            .flatten()
            .groupBy((p) => p.params.queryString())
            .map((p) => {
                return {
                    params: p[0].params, // All in this group are identical.
                    vizes: _(p)
                        .map((p) => p.viz)
                        .value(),
                };
            })
            .value();

        return Promise.all(
            vizToParams.map((query) => {
                return this.querier.query(query.params).then((data) => {
                    query.vizes.forEach((viz: Viz) => {
                        viz.log("data", query.params.queryString(), data.data.length, data.timeRange, data.dataRange);
                        if (viz instanceof QueriesSensorData) {
                            viz.data = data;
                        }
                    });
                });
            })
        );
    }

    public zoomed(viz: Viz, times: TimeRange) {
        viz.zoomed(times);
        return this.query();
    }

    public get options(): TreeOption[] {
        const allSensors = _(this.stations)
            .map((station: StationMeta) =>
                station.sensors.map((sensor) => {
                    return {
                        station: station,
                        sensor: sensor,
                    };
                })
            )
            .flatten()
            .value();

        const options = _(allSensors)
            .groupBy((s) => s.sensor.sensor.fullKey)
            .values()
            .map((value: { station: DisplayStation; sensor: SensorMeta }[]) => {
                const sensor = value[0].sensor;
                const stations: StationMeta[] = _(value)
                    .map((v) => v.station)
                    .value();
                return new TreeOption(
                    sensor.fullKey,
                    sensor.fullKey,
                    sensor,
                    _(stations)
                        .map((station) => station.station.id)
                        .value(),
                    _(stations)
                        .map((station) => {
                            return new TreeOption("s-" + station.station.id + "-" + sensor.fullKey, station.station.name, sensor, [
                                station.station.id,
                            ]);
                        })
                        .value()
                );
            })
            .value();

        return options;
    }

    public remove(viz: Viz) {
        this.groups = this.groups.map((g) => g.remove(viz)).filter((g) => !g.empty);
    }

    public selected(option: TreeOption) {
        this.addSensor(option.sensor, option.stations);
        return this.query();
    }
}
