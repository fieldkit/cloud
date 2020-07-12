import _ from "lodash";
import { SensorId, SensorsResponse, SensorDataResponse, SensorInfoResponse, ModuleSensor } from "./api";
import { Ids, TimeRange, Stations, Sensors, SensorParams, DataQueryParams } from "./common";
import { DisplayStation } from "@/store/modules/stations";
import FKApi from "@/api/api";

import { ColorScale, createSensorColorScale } from "./d3-helpers";

export class SensorMeta {
    constructor(public readonly id: number, public readonly fullKey: string, public readonly name: string) {}
}

export class StationMeta {
    constructor(public readonly id: number, public readonly name: string, public readonly sensors: SensorMeta[]) {}
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
    constructor(public readonly colorScale: ColorScale, public readonly station: { name: string }) {}
}

export abstract class Viz {
    public readonly id = Ids.make();

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

type GroupBookmark = [VizBookmark[]];

export class Bookmark {
    static Version = 1;

    constructor(public readonly v: number, public readonly g: GroupBookmark[]) {}
}

export class TimeZoom {
    constructor(public readonly fast: FastTime | null, public readonly range: TimeRange | null) {}
}

export class Graph extends Viz {
    public graphing: QueriedData | null = null;
    public all: QueriedData | null = null;
    public visible: TimeRange = TimeRange.eternity;
    public chartType: ChartType = ChartType.TimeSeries;
    public fastTime: FastTime = FastTime.All;

    constructor(public chartParams: DataQueryParams) {
        super();
    }

    public clone(): Viz {
        const c = new Graph(this.chartParams);
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

    public bookmark(): VizBookmark {
        return [this.chartParams.stations, this.chartParams.sensors, [this.visible.start, this.visible.end], this.chartType, this.fastTime];
    }

    public static fromBookmark(bm: VizBookmark): Viz {
        const visible = new TimeRange(bm[2][0], bm[2][1]);
        const chartParams = new DataQueryParams(visible, bm[0], bm[1]);
        const graph = new Graph(chartParams);
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
        return new Scrubbers(
            this.id,
            this.visible_,
            this.vizes.filter((viz) => (viz as Graph).all).map((viz, i) => new Scrubber(i, (viz as Graph).all))
        );
    }

    public bookmark(): GroupBookmark {
        return [this.vizes.map((v) => v.bookmark())];
    }

    public static fromBookmark(bm: GroupBookmark): Group {
        return new Group(bm[0].map((vm) => Graph.fromBookmark(vm)));
    }
}

type ResolveData = (qd: QueriedData) => void;

class VizQuery {
    constructor(public readonly params: DataQueryParams, public readonly resolve: ResolveData) {}
}

class InfoQuery {
    constructor(public readonly params: Stations) {}
}

export class Querier {
    private info: { [index: string]: SensorInfoResponse } = {};
    private data: { [index: string]: QueriedData } = {};

    public queryInfo(params: Stations): Promise<SensorInfoResponse> {
        if (!params) {
            throw new Error("no params");
        }

        const queryParams = new URLSearchParams();
        queryParams.append("stations", params.join(","));

        const key = queryParams.toString();
        if (this.info[key]) {
            return Promise.resolve(this.info[key]);
        }
        return new FKApi().sensorData(queryParams).then((info: SensorInfoResponse) => {
            this.info[key] = info;
            return info;
        });
    }

    public queryData(params: DataQueryParams): Promise<QueriedData> {
        if (!params) {
            throw new Error("no params");
        }

        const queryParams = params.queryParams();
        const key = queryParams.toString();
        if (this.data[key]) {
            return Promise.resolve(this.data[key]);
        }
        return new FKApi().sensorData(queryParams).then((sdr: SensorDataResponse) => {
            const queried = new QueriedData(params.when, sdr);
            this.data[key] = queried;
            return queried;
        });
    }
}

export class Workspace {
    private readonly querier = new Querier();
    private readonly stations: { [index: number]: StationMeta } = {};
    private _options: TreeOption[] = [];

    public get options(): TreeOption[] {
        return this._options;
    }

    constructor(private readonly meta: SensorsResponse, private groups: Group[] = []) {}

    private get allVizes(): Viz[] {
        return _(this.groups)
            .map((g) => g.vizes)
            .flatten()
            .value();
    }

    public query(): Promise<any> {
        const allGraphs = this.allVizes.map((viz) => viz as Graph).filter((viz) => viz);

        // First we may need some data for making the UI useful and
        // pretty. Filling in labels, building drop downs, etc... This
        // is especially important to do from here because we may have
        // been instantiated from a Bookmark. Right now we just query
        // for information on all the stations involved.
        const allStationIds = _.uniq(_.flatten(allGraphs.map((viz) => viz).map((viz) => viz.chartParams.stations)));
        const infoQueries = allStationIds.length ? [new InfoQuery(allStationIds)] : [];

        // Second step is to query to fill in any required scrubbers. I
        // have tried in previous iterations to be clever about this
        // and just being very explicit is the best way.
        const scrubberQueries = allGraphs
            .filter((viz) => !viz.all)
            .map(
                (viz: Graph) =>
                    new VizQuery(viz.scrubberParams, (qd) => {
                        viz.all = qd;
                        return;
                    })
            );

        // Now build the queries for the data being viewed.
        const graphingQueries = allGraphs.map(
            (viz: Graph) =>
                new VizQuery(viz.chartParams, (qd) => {
                    viz.graphing = qd;
                    return;
                })
        );

        // Combine and make them unique to avoid obvious
        // duplicates. Eventually we can also merge stations/sensors
        // with the same date range and other parameters into one
        // query. Lots of room here.
        const allQueries = [...scrubberQueries, ...graphingQueries];
        const uniqueQueries = _(allQueries)
            .groupBy((q) => q.params.queryParams().toString())
            .map((p) => new VizQuery(p[0].params, (qd: QueriedData) => p.map((p) => p.resolve(qd))))
            .value();

        console.log("workspace: querying", uniqueQueries.length, "data", infoQueries.length, "info");

        // Make all the queries and then give the queried data to the
        // resolve call for that query. This will end up calling the
        // above mapped resolve to set the appropriate data.
        const pendingInfo = infoQueries.map((iq) =>
            this.querier.queryInfo(iq.params).then((info) => {
                return _.map(info.stations, (info, stationId) => {
                    const name = info[0].stationName;
                    const sensors = info.map((row) => new SensorMeta(row.sensorId, row.key, row.key));
                    const station = new StationMeta(Number(stationId), name, sensors);
                    this.stations[station.id] = station;
                    return station;
                });
            })
        );

        const pendingData = uniqueQueries.map((vq) => this.querier.queryData(vq.params).then((data) => vq.resolve(data)));
        return Promise.all([...pendingInfo, ...pendingData]).then(() => {
            this._options = this.updateOptions();
        });
    }

    public vizInfo(viz: Graph): VizInfo {
        const sensorDetailsByKey = _.keyBy(_.flatten(this.meta.modules.map((m) => m.sensors)), (s) => s.fullKey);
        const keysById = _.keyBy(this.meta.sensors, (s) => s.id);

        if (viz.chartParams.stations.length != 1) {
            throw new Error("expected 1 station per graph, for now");
        }

        if (viz.chartParams.sensors.length != 1) {
            throw new Error("expected 1 sensor per graph, for now");
        }

        const stationId = viz.chartParams.stations[0];
        const sensorId = viz.chartParams.sensors[0];

        const station = this.stations[stationId];
        const key = keysById[sensorId].key;
        const details = sensorDetailsByKey[key];
        const scale = createSensorColorScale(details);

        return new VizInfo(scale, station);
    }

    public graphZoomed(viz: Viz, zoom: TimeZoom): Workspace {
        this.findGroup(viz).zoomed(zoom);
        return this;
    }

    public groupZoomed(group: Group, zoom: TimeZoom): Workspace {
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

    public addGraph(graph: Graph): Workspace {
        const group = new Group();
        group.add(graph);
        this.groups.unshift(group);
        return this;
    }

    public addStandardGraph(stations: Stations, sensors: Sensors): Workspace {
        return this.addGraph(new Graph(new DataQueryParams(TimeRange.eternity, stations, sensors)));
    }

    private updateOptions(): TreeOption[] {
        console.log("workspace: stations:", this.stations);

        const allSensors = _(Object.values(this.stations))
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
            .map((value) => {
                const sensor = value[0].sensor;
                const stations = _(value)
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

        console.log("workspace: options:", { options: options });

        return options;
    }

    public remove(viz: Viz): Workspace {
        this.groups = this.groups.map((g) => g.remove(viz)).filter((g) => !g.empty);
        return this;
    }

    public compare(viz: Viz): Workspace {
        if (this.groups.length > 0) {
            this.groups.unshift(this.groups[0].clone());
        }
        return this;
    }

    public changeChart(viz: Viz, chartType: ChartType): Workspace {
        if (viz instanceof Graph) {
            viz.changeChart(chartType);
        }
        return this;
    }

    public changeSensors(viz: Viz, option: TreeOption): Workspace {
        if (viz instanceof Graph) {
            viz.changeSensors(option);
        }
        return this;
    }

    public changeLinkage(viz: Viz): Workspace {
        const group = this.findGroup(viz);
        if (group.vizes.length > 1) {
            group.remove(viz);
            this.groups.unshift(new Group([viz]));
        } else {
            this.groups.reduce((previous, iter) => {
                if (iter == group) {
                    if (previous == null) {
                        throw new Error("tried linking first group, nice work");
                    }
                    this.removeGroup(group);
                    previous.addAll(group);
                }
                return iter;
            }, null);
        }
        return this;
    }

    private removeGroup(group: Group): Workspace {
        this.groups = this.groups.filter((g) => g !== group);
        return this;
    }

    /**
     * Pop first Group and move all the Viz children to the new first Group
     */
    private combine() {
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

    public eventually(callback: (ws: Workspace) => Promise<any>) {
        callback(this);
        return Promise.resolve(this);
    }

    public with(callback: (ws: Workspace) => Workspace) {
        callback(this);
        return this;
    }
}
