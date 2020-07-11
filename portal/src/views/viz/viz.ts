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

    public abstract bookmark(): VizBookmark;
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

type VizBookmark = [Stations, Sensors, [number, number], ChartType, FastTime] | [];
type GroupBookmark = [boolean, VizBookmark[]];

export class Bookmark {
    static Version = 1;

    constructor(public readonly v: number, public readonly g: GroupBookmark[]) {}
}

export class TimeZoom {
    constructor(public readonly fast: FastTime | null, public readonly range: TimeRange | null) {}
}

export class Graph extends Viz {
    private graphing: QueriedData | null = null;
    public all: QueriedData | null = null;
    public visible: TimeRange = TimeRange.eternity;
    public chartType: ChartType = ChartType.TimeSeries;
    public fastTime: FastTime = FastTime.All;

    constructor(public info: VizInfo, public chartParams: DataQueryParams) {
        super(info);
    }

    public clone(): Viz {
        const c = new Graph(this.info, this.chartParams);
        c.all = this.all;
        c.visible = this.visible;
        c.chartType = this.chartType;
        c.fastTime = this.fastTime;
        return c;
    }

    public get scrubberParams(): DataQueryParams {
        return new DataQueryParams(TimeRange.eternity, this.chartParams.stations, this.chartParams.sensors);
    }

    public zoomed(zoom: TimeZoom): TimeRange {
        if (zoom.range) {
            this.visible = zoom.range;
            this.fastTime = FastTime.Custom;
        } else {
            this.visible = this.getFastRange(zoom.fast);
            this.fastTime = zoom.fast;
        }

        this.chartParams = new DataQueryParams(this.visible, this.chartParams.stations, this.chartParams.sensors);

        return this.visible;
    }

    private getFastRange(fastTime: FastTime) {
        if (!this.all) throw new Error("fast time missing all data");

        const sensorRange = this.all.timeRange;
        if (fastTime === FastTime.All) {
            return TimeRange.eternity;
        } else {
            const days = fastTime as number;
            const start = new Date(sensorRange[1]);
            start.setDate(start.getDate() - days);
            return new TimeRange(start.getTime(), sensorRange[1]);
        }
    }

    public changeChart(chartType: ChartType) {
        this.chartType = chartType;
    }

    public changeSensors(option: TreeOption) {
        const sensorParams = option.sensorParams;
        this.chartParams = new DataQueryParams(this.chartParams.when, sensorParams.stations, sensorParams.sensors);
        this.all = null;
    }

    public set data(qd: QueriedData) {
        if (this.shouldUpdateAllWith(qd)) {
            this.all = qd;
        }
        this.graphing = qd;
    }

    public get data(): QueriedData {
        return this.graphing;
    }

    private shouldUpdateAllWith(qd: QueriedData): boolean {
        if (this.all === null) {
            return true;
        }
        const allArray = this.all.timeRange;
        const allRange = new TimeRange(allArray[0], allArray[1]);
        return !allRange.contains(qd.timeRangeQueried);
    }

    public bookmark(): VizBookmark {
        return [this.chartParams.stations, this.chartParams.sensors, [this.visible.start, this.visible.end], this.chartType, this.fastTime];
    }

    public static fromBookmark(bm: VizBookmark): Viz {
        const visible = new TimeRange(bm[2][0], bm[2][1]);
        const chartParams = new DataQueryParams(visible, bm[0], bm[1]);
        const info = new VizInfo("TODO");
        const graph = new Graph(info, chartParams);
        graph.chartType = bm[3];
        graph.fastTime = bm[4];
        return graph;
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

    public zoomed(zoom: TimeZoom) {
        this.vizes.forEach((viz) => {
            if (viz instanceof Graph) {
                // Yeah this is kind of weird... the idea though is
                // that they'll all return the same thing.
                this.visible_ = viz.zoomed(zoom);
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

    public bookmark(): GroupBookmark {
        return [true, this.vizes.map((v) => v.bookmark())];
    }

    public static fromBookmark(bm: GroupBookmark): Group {
        return new Group(bm[1].map((vm) => Graph.fromBookmark(vm)));
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

type ResolveVizData = (qd: QueriedData) => void;

class VizQuery {
    constructor(public readonly params: DataQueryParams, public readonly resolve: ResolveVizData) {}
}

export class Workspace {
    public readonly querier = new Querier();
    public readonly stations: StationMeta[] = [];

    constructor(private readonly meta: SensorsResponse, public groups: Group[] = []) {}

    public addSensor(sensor: SensorMeta, stations: Stations) {
        const info = VizInfo.fromSensor(sensor);
        const graph = new Graph(info, new DataQueryParams(TimeRange.eternity, stations, [sensor.id]));
        return this.addGraph(graph);
    }

    public addGraph(graph: Graph) {
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
        const allGraphs = this.allVizes.map((viz) => viz as Graph);

        // First step is to query to fill in any required scrubbers. I
        // have tried in previous iterations to be clever about this
        // and just being very explicit is the best way.
        const scrubberQueries = allGraphs
            .filter((viz) => viz && !viz.all)
            .map((viz: Graph) => new VizQuery(viz.scrubberParams, (qd) => (viz.all = qd)));

        // Now build the queries for the data being viewed.
        const visibleQueries = allGraphs.filter((viz) => viz).map((viz: Graph) => new VizQuery(viz.chartParams, (qd) => (viz.data = qd)));

        // Combine and make them unique to avoid obvious
        // duplicates. Eventually we can also merge stations/sensors
        // with the same date range and other parameters into one
        // query. Lots of room here.
        const allQueries = [...scrubberQueries, ...visibleQueries];
        const uniqueQueries = _(allQueries)
            .groupBy((q) => q.params.queryString())
            .map((p) => new VizQuery(p[0].params, (qd: QueriedData) => p.map((p) => p.resolve(qd))))
            .value();

        console.log("workspace: querying", uniqueQueries.length, "queries");

        // Make all the queries and then give the queried data to the
        // resolve call for that query. This will end up calling the
        // above mapped resolve to set the appropriate data.
        return Promise.all(uniqueQueries.map((vq) => this.querier.query(vq.params).then((data) => vq.resolve(data))));
    }

    public graphZoomed(viz: Viz, zoom: TimeZoom) {
        this.findGroup(viz).zoomed(zoom);
        return this;
    }

    public groupZoomed(group: Group, zoom: TimeZoom) {
        group.zoomed(zoom);
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

    public addSensorAll(id: number) {
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
            .groupBy((r) => r.sensor.id)
            .mapValues((r) => r.map((s) => s.station.id))
            .value();

        console.log("workspace: adding:", allSensors);

        return allSensors[id].map((stationId) => {
            return this.addGraph(new Graph(new VizInfo("nothing"), new DataQueryParams(TimeRange.eternity, [stationId], [id])));
        });
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

    public bookmark(): Bookmark {
        return new Bookmark(
            Bookmark.Version,
            this.groups.map((group) => group.bookmark())
        );
    }

    public static fromBookmark(meta: SensorsResponse, bm: Bookmark): Workspace {
        if (bm.v !== 1) {
            throw new Error("unexpected bookmark version");
        }
        return new Workspace(
            meta,
            bm.g.map((gm) => Group.fromBookmark(gm))
        );
    }

    public with(callback: (ws: Workspace) => Workspace) {
        return callback(this);
    }
}
