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

    get id(): number {
        return this.station.id;
    }

    get sensors(): SensorMeta[] {
        const sensorsByKey = _.keyBy(this.meta.sensors, (s) => s.key);
        return _(this.station.configurations.all)
            .take(1)
            .map((cfg) => cfg.modules)
            .flatten()
            .filter((m) => !m.internal)
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

export abstract class Viz {
    public readonly id = Ids.make();

    constructor(public readonly info: VizInfo) {}

    public needs(viz: Viz) {
        return this.id === viz.id;
    }

    public log(...args: any[]) {
        console.log(...["viz:", this.id, this.constructor.name, ...args]);
    }

    public abstract clone(): Viz;
}

export class Graph extends Viz {
    private graphing: QueriedData | null = null;
    public all: QueriedData | null = null;
    public visible: TimeRange = TimeRange.eternity;

    constructor(public info: VizInfo, public params: DataQueryParams) {
        super(info);
    }

    public zoomed(range: TimeRange) {
        this.visible = range;
        this.params = new DataQueryParams(range, this.params.stations, this.params.sensors);
    }

    public clone(): Viz {
        return new Graph(this.info, this.params);
    }

    public set data(qd: QueriedData) {
        if (this.all == null) {
            this.all = qd;
        }
        this.graphing = qd;
    }

    public get data(): QueriedData {
        return this.graphing;
    }
}

export class Group {
    constructor(public vizes: Viz[] = []) {}

    public add(viz: Viz) {
        this.vizes.push(viz);
        return this;
    }

    public remove(removing: Viz) {
        this.vizes = this.vizes.filter((viz) => {
            return removing !== viz;
        });
        return this;
    }

    public get empty() {
        return this.vizes.length == 0;
    }

    public contains(viz: Viz): boolean {
        return this.vizes.indexOf(viz) >= 0;
    }

    public zoomed(times: TimeRange) {
        this.vizes.forEach((viz) => {
            if (viz instanceof Graph) {
                viz.zoomed(times);
            }
        });
    }

    public clone(): Group {
        return new Group(this.vizes.map((v) => v.clone()));
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

export class Workspace {
    public readonly querier = new Querier();
    public readonly stations: StationMeta[] = [];
    public groups: Group[] = [];

    constructor(private readonly meta: SensorsResponse) {}

    public addSensor(sensor: SensorMeta, stations: Stations): void {
        const info = VizInfo.fromSensor(sensor);
        const graph = new Graph(info, new DataQueryParams(TimeRange.eternity, stations, [sensor.id]));
        const group = new Group();
        group.add(graph);
        this.groups.push(group);
    }

    public addStation(station: DisplayStation): boolean {
        const stationMeta = new StationMeta(this.meta, station);
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
            .map((viz: Graph) =>
                [viz.params].map((params: DataQueryParams) => {
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

        console.log("workspace: querying", vizToParams.length, "queries");

        return Promise.all(
            vizToParams.map((query) => {
                return this.querier.query(query.params).then((data) => {
                    query.vizes.forEach((viz: Graph) => {
                        viz.log("data", query.params.queryString(), data.data.length, data.timeRange, data.dataRange);
                        viz.data = data;
                    });
                });
            })
        );
    }

    public zoomed(viz: Viz, times: TimeRange) {
        this.findGroup(viz).zoomed(times);
        return this.query();
    }

    private findGroup(viz: Viz): Group {
        const groups = this.groups.filter((g) => g.contains(viz));
        if (groups.length == 1) {
            return groups[0];
        }
        throw new Error("oprhaned viz");
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
                        .map(
                            (station) =>
                                new TreeOption("s-" + station.station.id + "-" + sensor.fullKey, station.station.name, sensor, [
                                    station.station.id,
                                ])
                        )
                        .value()
                );
            })
            .value();

        return options;
    }

    public remove(viz: Viz) {
        this.groups = this.groups.map((g) => g.remove(viz)).filter((g) => !g.empty);
    }

    public compare() {
        if (this.stations.length == 0) {
            throw new Error("no stations");
        }

        const allSensors = _(this.stations)
            .map((s) => s.sensors)
            .flatten()
            .groupBy((s) => s.fullKey)
            .value();

        console.log(allSensors);

        if (this.groups.length == 0) {
            const station = this.stations[0];
            const stationSensors = station.sensors;
            console.log("sensors:", stationSensors);
            if (stationSensors.length == 0) {
                return false;
            }
            const sensor = stationSensors[0];
            this.addSensor(sensor, [station.id]);
        } else {
            console.log(this.groups[0].clone());
            this.groups.unshift(this.groups[0].clone());
        }
        return this.query();
    }

    public selected(viz: Viz, option: TreeOption) {
        // this.addSensor(option.sensor, option.stations);
        return this.query();
    }
}
