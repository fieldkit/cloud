import _ from "lodash";
import { SensorId, SensorsResponse, SensorDataResponse, ModuleSensor } from "./api";
import { Ids, TimeRange, Stations, Sensors, SensorParams, DataQueryParams } from "./common";
import { DisplayStation } from "@/store/modules/stations";
import FKApi from "@/api/api";

export class QuickSensorStation {
    constructor(public readonly sensorId: number, public readonly key: string) {}
}

export class QuickSensors {
    constructor(public readonly stations: { [index: string]: QuickSensorStation[] }) {}
}

export class SensorMeta {
    constructor(public readonly id: number, public readonly fullKey: string, public readonly name: string) {}

    public static makePlain(meta: SensorId) {
        return new SensorMeta(meta.id, meta.key, meta.key);
    }

    public static makePlainFromQuickie(quickie: QuickSensorStation) {
        return new SensorMeta(quickie.sensorId, quickie.key, quickie.key);
    }
}

export class StationMeta {
    constructor(private readonly meta: SensorsResponse, private quickSensors: QuickSensors, private station: DisplayStation) {}

    get id(): number {
        return this.station.id;
    }

    get name(): string {
        return this.station.name;
    }

    get sensors(): SensorMeta[] {
        // This toString is so annoying. -jlewallen
        const sensorsByKey = _.keyBy(this.meta.sensors, (s) => s.key);
        const quickies = this.quickSensors.stations[this.id.toString()].map((q) => SensorMeta.makePlainFromQuickie(q));
        const fromConfiguration = _(this.station.configurations.all)
            .take(1)
            .map((cfg) => cfg.modules)
            .flatten()
            .filter((m) => !m.internal)
            .map((m) => m.sensors)
            .flatten()
            .map((s) => SensorMeta.makePlain(sensorsByKey[s.fullKey]))
            .filter((s) => {
                return true;
            })
            .value();

        return Object.values(
            Object.assign(
                _.keyBy(quickies, (s) => s.id),
                _.keyBy(fromConfiguration, (s) => s.id)
            )
        );
    }
}

export class TreeOption {
    constructor(
        public readonly id: string,
        public readonly label: string,
        public readonly sensorParams: SensorParams,
        public readonly children: TreeOption[] | undefined = undefined
    ) {}
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

    public log(...args: any[]) {
        console.log(...["viz:", this.id, this.constructor.name, ...args]);
    }

    public abstract clone(): Viz;

    public get canRemove(): boolean {
        return true;
    }
}

export enum FastTime {
    Custom = -1,
    Day = 1,
    Week = 7,
    TwoWeeks = 14,
    Month = 30,
    Year = 365,
    All = 0,
}

export enum ChartType {
    TimeSeries,
    Histogram,
    Range,
    Map,
}

export class Graph extends Viz {
    private graphing: QueriedData | null = null;
    public all: QueriedData | null = null;
    public visible: TimeRange = TimeRange.eternity;
    public chartType: ChartType = ChartType.TimeSeries;
    public fastTime: FastTime = FastTime.All;

    constructor(public info: VizInfo, public params: DataQueryParams) {
        super(info);
    }

    public clone(): Viz {
        return new Graph(this.info, this.params);
    }

    public zoomed(range: TimeRange) {
        this.visible = range;
        this.fastTime = FastTime.Custom;
        this.params = new DataQueryParams(range, this.params.stations, this.params.sensors);
    }

    public useFastTime(fastTime: FastTime) {
        if (!this.all) throw new Error("fast time missing all data");

        const calculate = () => {
            const sensorRange = this.all.timeRange;
            if (fastTime == FastTime.All) {
                return new TimeRange(sensorRange[0], sensorRange[1]);
            } else {
                const days = fastTime as number;
                const start = new Date(sensorRange[1]);
                start.setDate(start.getDate() - days);
                return new TimeRange(start.getTime(), sensorRange[1]);
            }
        };

        this.visible = calculate();
        this.params = new DataQueryParams(this.visible, this.params.stations, this.params.sensors);
        this.fastTime = fastTime;
    }

    public changeChart(chartType: ChartType) {
        this.chartType = chartType;
    }

    public changeSensors(option: TreeOption) {
        const sensorParams = option.sensorParams;
        this.params = new DataQueryParams(this.params.when, sensorParams.stations, sensorParams.sensors);
        this.all = null;
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

export class Scrubber {
    constructor(public readonly index: number, public readonly data: QueriedData) {}
}

export class Scrubbers {
    public readonly timeRange: TimeRange;

    constructor(public readonly id: string, public readonly visible: TimeRange, public readonly rows: Scrubber[]) {
        this.timeRange = TimeRange.mergeArrays(rows.map((s) => s.data.timeRange));
    }
}

export class Group {
    public readonly id = Ids.make();
    private visible_: TimeRange = TimeRange.eternity;

    constructor(public vizes: Viz[] = []) {}

    public log(...args: any[]) {
        console.log(...["viz:", this.id, this.constructor.name, ...args]);
    }

    public clone(): Group {
        return new Group(this.vizes.map((v) => v.clone()));
    }

    public add(viz: Viz) {
        this.vizes.push(viz);
        return this;
    }

    public addAll(o: Group) {
        o.vizes.forEach((v) => this.add(v));
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
        this.visible_ = times;
        this.vizes.forEach((viz) => {
            if (viz instanceof Graph) {
                viz.zoomed(times);
            }
        });
    }

    public get scrubbers(): Scrubbers {
        if (this.vizes.length > 1) {
            return new Scrubbers(
                this.id,
                this.visible_,
                this.vizes.filter((viz) => (viz as Graph).all).map((viz, i) => new Scrubber(i, (viz as Graph).all))
            );
        } else {
            return null;
        }
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

export class Workspace {
    public readonly querier = new Querier();
    public readonly stations: StationMeta[] = [];
    public groups: Group[] = [];

    constructor(private readonly meta: SensorsResponse) {}

    public addSensor(sensor: SensorMeta, stations: Stations) {
        const info = VizInfo.fromSensor(sensor);
        const graph = new Graph(info, new DataQueryParams(TimeRange.eternity, stations, [sensor.id]));
        const group = new Group();
        group.add(graph);
        this.groups.push(group);
        return this;
    }

    public addStation(quickSensors: QuickSensors, station: DisplayStation) {
        const stationMeta = new StationMeta(this.meta, quickSensors, station);
        this.stations.push(stationMeta);
        return this;
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

    public graphZoomed(viz: Viz, times: TimeRange) {
        this.findGroup(viz).zoomed(times);
        return this;
    }

    public groupZoomed(group: Group, times: TimeRange) {
        group.zoomed(times);
        return this;
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
            .groupBy((s) => s.sensor.fullKey)
            .values()
            .map((value: { station: DisplayStation; sensor: SensorMeta }[]) => {
                const sensor = value[0].sensor;
                const stations: StationMeta[] = _(value)
                    .map((v) => v.station)
                    .value();

                const label = sensor.fullKey;
                const stationIds = stations.map((station) => station.id);
                const sensorIds = [sensor.id];
                const sensorParams = new SensorParams(sensorIds, stationIds);
                const children = stations
                    .map((station) => {
                        const label = station.name;
                        const sensorParams = new SensorParams(sensorIds, [station.id]);
                        return new TreeOption(sensorParams.id, label, sensorParams);
                    })
                    .filter((child) => {
                        return child.id != sensorParams.id;
                    });

                if (children.length == 0) {
                    return new TreeOption(sensorParams.id, label, sensorParams);
                }

                return new TreeOption(sensorParams.id, label, sensorParams, children);
            })
            .value();

        console.log("workspace: options:", options);

        return options;
    }

    public remove(viz: Viz) {
        this.groups = this.groups.map((g) => g.remove(viz)).filter((g) => !g.empty);
        return this;
    }

    public compare(viz: Viz) {
        if (this.stations.length == 0) {
            throw new Error("no stations");
        }

        const allSensors = _(this.stations)
            .map((s) => s.sensors)
            .flatten()
            .keyBy((s) => s.fullKey)
            .value();

        console.log("workspace: sensors:", allSensors);

        if (this.groups.length == 0) {
            const station = this.stations[0];
            const stationSensors = station.sensors;
            console.log("workspace: station:", stationSensors);
            if (stationSensors.length == 0) {
                return this;
            }
            const sensor = stationSensors[0];
            this.addSensor(sensor, [station.id]);
        } else {
            this.groups.unshift(this.groups[0].clone());
        }

        return this;
    }

    public changeChart(viz: Viz, chartType: ChartType) {
        if (viz instanceof Graph) {
            viz.changeChart(chartType);
        }
        return this;
    }

    public changeSensors(viz: Viz, option: TreeOption) {
        if (viz instanceof Graph) {
            viz.changeSensors(option);
        }
        return this;
    }

    public fastTime(viz: Viz, fastTime: FastTime) {
        if (viz instanceof Graph) {
            viz.useFastTime(fastTime);
        }
        return this;
    }

    /**
     * Pop first Group and move all the Viz children to the new first Group
     */
    public combine() {
        if (this.groups.length <= 1) {
            return this;
        }
        const removing = this.groups.shift();
        this.groups[0].addAll(removing);
        return this;
    }
}
