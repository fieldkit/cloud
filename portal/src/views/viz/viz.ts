import _ from "lodash";
import moment, { Moment } from "moment";
import { SensorsResponse, SensorDataResponse, SensorInfoResponse } from "./api";
import { ModuleID, SensorSpec, Ids, TimeRange, Stations, Sensors, SensorParams, DataQueryParams } from "./common";
import i18n from "@/i18n";
import FKApi from "@/api/api";

export * from "./common";

import { ColorScale, createSensorColorScale } from "./d3-helpers";

type SensorReadAtType = string;

export class SensorMeta {
    constructor(
        public readonly moduleId: ModuleID,
        public readonly moduleKey: string,
        public readonly sensorId: number,
        public readonly sensorKey: string,
        public readonly sensorReadAt: SensorReadAtType
    ) {}
}

export class StationMeta {
    constructor(public readonly id: number, public readonly name: string, public readonly sensors: SensorMeta[]) {
        if (!this.id) throw new Error("id is required");
        if (!this.name) throw new Error("name is required");
        if (!this.sensors) throw new Error("sensors is required");
    }
}

export interface HasSensorParams {
    readonly sensorParams: SensorParams;
}

export class StationTreeOption {
    constructor(public readonly id: string | number, public readonly label: string) {}
}

export class SensorTreeOption {
    constructor(
        public readonly id: string | number,
        public readonly label: string,
        public readonly children: SensorTreeOption[] | undefined = undefined,
        public readonly moduleId: ModuleID,
        public readonly sensorId: number | null,
        public readonly age: Moment
    ) {}
}

function makeRange(values: number[]): [number, number] {
    const min = _.min(values);
    const max = _.max(values);
    if (min === undefined) throw new Error(`no min: ${values.length}`);
    if (max === undefined) throw new Error(`no max: ${values.length}`);
    if (min === max) {
        console.warn(`range-warning: min == max ${values.length}`);
    }
    return [min, max];
}

export class QueriedData {
    empty = true;
    dataRange: number[] = [];
    timeRangeData: number[] = [];
    timeRange: number[] = [];

    constructor(public readonly timeRangeQueried: TimeRange, private readonly sdr: SensorDataResponse) {
        if (this.sdr.data.length > 0) {
            const filtered = this.sdr.data.filter((d) => _.isNumber(d.value));
            const values = filtered.map((d) => d.value);
            const times = filtered.map((d) => d.time);

            if (values.length == 0) throw new Error(`empty data ranges`);
            if (times.length == 0) throw new Error(`empty time ranges`);

            this.dataRange = makeRange(values);
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
}

export class VizInfo {
    constructor(public readonly key: string, public readonly colorScale: ColorScale, public readonly station: { name: string }) {}
}

export abstract class Viz {
    protected busyDepth = 0;
    public readonly id = Ids.make();

    public howBusy(d: number): boolean {
        this.busyDepth += d;
        return this.busy;
    }

    public get busy(): boolean {
        return this.busyDepth > 0;
    }

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

    public get empty(): boolean {
        return this.rows.length == 0;
    }

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

type VizBookmark = [Stations, Sensors, [number, number], [[number, number], [number, number]] | [], ChartType, FastTime] | [];

type GroupBookmark = [VizBookmark[]];

export class Bookmark {
    static Version = 1;

    constructor(public readonly v: number, public readonly g: GroupBookmark[], public readonly s: number[] = []) {}

    private get allVizes(): VizBookmark[] {
        return _.flatten(_.flatten(this.g.map((group) => group.map((vizes) => vizes))));
    }

    public get allTimeRange(): TimeRange {
        const times: [number, number][] = this.allVizes.map((viz) => viz[2]).filter((times) => times !== undefined) as [number, number][];
        const start = _.min(_.flatten(times.map((r) => r[0])));
        const end = _.max(_.flatten(times.map((r) => r[1])));
        if (!start || !end) throw new Error(`no time range in bookmark`);
        return new TimeRange(start, end);
    }

    public get allStations(): number[] {
        return _.uniq(_.flatten(this.allVizes.map((viz) => viz[0] || [])));
    }

    public get allSensors(): SensorSpec[] {
        return _.uniq(_.flatten(this.allVizes.map((viz) => viz[1] || [])).map((sensorAndModule) => sensorAndModule));
    }

    public static sameAs(a: Bookmark, b: Bookmark): boolean {
        const aEncoded = JSON.stringify(a);
        const bEncoded = JSON.stringify(b);
        return aEncoded == bEncoded;
    }
}

export class GeoZoom {
    constructor(public readonly bounds: [[number, number], [number, number]]) {}
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
    public geo: GeoZoom | null = null;

    constructor(public chartParams: DataQueryParams) {
        super();
        if (this.chartParams.sensors.length != 1) throw new Error(`viz: Graph too many sensors`);
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

    public timeZoomed(zoom: TimeZoom): TimeRange {
        if (zoom.range !== null) {
            this.visible = zoom.range;
            this.fastTime = FastTime.Custom;
        } else if (zoom.fast !== null) {
            this.visible = this.getFastRange(zoom.fast);
            this.fastTime = zoom.fast;
        }

        this.chartParams = new DataQueryParams(this.visible, this.chartParams.stations, this.chartParams.sensors);

        return this.visible;
    }

    public geoZoomed(zoom: GeoZoom): GeoZoom {
        this.geo = zoom;
        return this.geo;
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

    public changeSensors(option: HasSensorParams) {
        console.log(`changing-sensors`, option);
        const sensorParams = option.sensorParams;
        this.chartParams = new DataQueryParams(this.chartParams.when, sensorParams.stations, sensorParams.sensors);
        this.all = null;
    }

    public bookmark(): VizBookmark {
        return [
            this.chartParams.stations,
            this.chartParams.sensors,
            [this.visible.start, this.visible.end],
            this.geo ? this.geo.bounds : [],
            this.chartType,
            this.fastTime,
        ];
    }

    public static fromBookmark(bm: VizBookmark): Viz {
        if (bm.length == 0) throw new Error(`empty bookmark`);
        const visible = new TimeRange(bm[2][0], bm[2][1]);
        const chartParams = new DataQueryParams(visible, bm[0], bm[1]);
        const graph = new Graph(chartParams);
        graph.geo = bm[3].length ? new GeoZoom(bm[3]) : null;
        graph.chartType = bm[4];
        graph.fastTime = bm[5];
        graph.visible = visible;
        return graph;
    }
}

export class Group {
    public readonly id = Ids.make();
    private visible_: TimeRange = TimeRange.eternity;

    constructor(public vizes: Viz[] = []) {
        // This returns eternity when merging empty array.
        this.visible_ = TimeRange.mergeRanges(
            vizes
                .map((v) => v as Graph)
                .filter((v) => v)
                .map((v) => v.visible)
        );
    }

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

    public timeZoomed(zoom: TimeZoom) {
        this.vizes.forEach((viz) => {
            if (viz instanceof Graph) {
                // Yeah this is kind of weird... the idea though is
                // that they'll all return the same thing.
                this.visible_ = viz.timeZoomed(zoom);
            }
        });
    }

    public geoZoomed(zoom: GeoZoom) {
        this.vizes.forEach((viz) => {
            if (viz instanceof Graph) {
                viz.geoZoomed(zoom);
            }
        });
    }

    public get scrubbers(): Scrubbers {
        const children = this.vizes
            .map((viz) => viz as Graph)
            .map((graph, index) => {
                return {
                    graph,
                    index,
                };
            })
            .filter((r) => r.graph.all != null)
            .map((r) => {
                const all = r.graph.all;
                if (!all) throw new Error(`no viz data on Graph`);
                return new Scrubber(r.index, all);
            });

        return new Scrubbers(this.id, this.visible_, children);
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
    constructor(public readonly params: DataQueryParams, public readonly vizes: Viz[], public readonly resolve: ResolveData) {}

    public howBusy(d: number): any {
        return this.vizes.map((v) => v.howBusy(d));
    }
}

class InfoQuery {
    constructor(public readonly params: Stations) {}
}

export class Querier {
    private info: { [index: string]: SensorInfoResponse } = {};
    private data: { [index: string]: QueriedData } = {};

    public queryInfo(iq: InfoQuery): Promise<SensorInfoResponse> {
        if (!iq.params) throw new Error("no params");

        const params = iq.params;
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

    public queryData(vq: VizQuery): Promise<QueriedData> {
        if (!vq.params) throw new Error("no params");

        const params = vq.params;
        const queryParams = params.queryParams();
        const key = queryParams.toString();
        if (this.data[key]) {
            vq.resolve(this.data[key]);
            return Promise.resolve(this.data[key]);
        }

        vq.howBusy(1);

        return new FKApi()
            .sensorData(queryParams)
            .then((sdr: SensorDataResponse) => {
                const queried = new QueriedData(params.when, sdr);
                this.data[key] = queried;
                return queried;
            })
            .then((data) => {
                vq.howBusy(-1);
                vq.resolve(data);
                return data;
            });
    }
}

export class Workspace {
    private readonly querier = new Querier();
    private readonly stations: { [index: number]: StationMeta } = {};
    public version = 0;

    public get empty(): boolean {
        return this.allVizes.length === 0;
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
                    new VizQuery(viz.scrubberParams, [viz], (qd) => {
                        viz.all = qd;
                        return;
                    })
            );

        // Now build the queries for the data being viewed.
        const graphingQueries = allGraphs.map(
            (viz: Graph) =>
                new VizQuery(viz.chartParams, [viz], (qd) => {
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
            .map((p) => new VizQuery(p[0].params, _.flatten(p.map((p) => p.vizes)), (qd: QueriedData) => p.map((p) => p.resolve(qd))))
            .value();

        console.log("workspace: querying", uniqueQueries.length, "data", infoQueries.length, "info");

        // Make all the queries and then give the queried data to the
        // resolve call for that query. This will end up calling the
        // above mapped resolve to set the appropriate data.
        const pendingInfo = infoQueries.map(
            (iq) =>
                this.querier.queryInfo(iq).then((info) => {
                    return _.map(info.stations, (info, stationId) => {
                        const stationName = info[0].stationName;
                        const sensors = info.map(
                            (row) => new SensorMeta(row.moduleId, row.moduleKey, row.sensorId, row.sensorKey, row.sensorReadAt)
                        );
                        const station = new StationMeta(Number(stationId), stationName, sensors);
                        this.stations[station.id] = station;
                        console.log("station-meta", { station, info });
                        return station;
                    });
                }) as Promise<unknown>
        );

        const pendingData = uniqueQueries.map((vq) => this.querier.queryData(vq) as Promise<unknown>);
        return Promise.all([...pendingInfo, ...pendingData]).then(() => {
            // Update options here if doing so lazily.
            this.version++;
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
        const sensor = viz.chartParams.sensors[0];
        const sensorId = sensor[1];

        const station = this.stations[stationId];
        const key = keysById[sensorId].key;
        const details = sensorDetailsByKey[key];
        const scale = createSensorColorScale(details);

        return new VizInfo(key, scale, station);
    }

    public graphTimeZoomed(viz: Viz, zoom: TimeZoom): Workspace {
        this.findGroup(viz).timeZoomed(zoom);
        return this;
    }

    public graphGeoZoomed(viz: Viz, zoom: GeoZoom): Workspace {
        this.findGroup(viz).geoZoomed(zoom);
        return this;
    }

    public groupZoomed(group: Group, zoom: TimeZoom): Workspace {
        group.timeZoomed(zoom);
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

    public get stationOptions(): StationTreeOption[] {
        return Object.values(this.stations).map((station) => {
            return new StationTreeOption(station.id, station.name);
        });
    }

    public get sensorOptions(): SensorTreeOption[] {
        const allSensors = _.flatten(Object.values(this.stations).map((station) => station.sensors));
        const allModules = _.groupBy(allSensors, (s) => s.moduleId);
        const keysById = _.fromPairs(allSensors.map((row) => [row.moduleId, row.moduleKey]));
        const allModulesByModuleKey = _.keyBy(this.meta.modules, (m) => m.key);

        // console.log("all-sensors", allSensors);
        // console.log("all-modules", allModules);
        // console.log("keys-by-id", keysById);
        // console.log("all-modules-by-module-key", allModulesByModuleKey);

        const options = _.map(
            allModules,
            (sensors, moduleId: ModuleID): StationTreeOption => {
                const uniqueSensors = _.uniqBy(sensors, (s) => s.sensorId);
                const children: SensorTreeOption[] = _.flatten(
                    uniqueSensors.map((row) => {
                        const age = moment.utc(row.sensorReadAt);
                        const label = i18n.tc(row.sensorKey) || row.sensorKey;
                        const optionId = `${row.moduleId}-${row.sensorId}`;
                        return [new SensorTreeOption(optionId, label, undefined, row.moduleId, row.sensorId, age)];
                    })
                );
                const moduleKey = keysById[moduleId];
                const moduleAge = _.max(children.map((c) => c.age));
                if (!moduleAge) throw new Error(`expected module age: no sensors?`);

                const label = i18n.tc(moduleKey); //  + ` (${moduleAge.fromNow()})`;
                return new SensorTreeOption(`${moduleKey}-${moduleId}`, label, children, moduleId, null, moduleAge);
            }
        );

        const sorted = _.sortBy((options as unknown) as SensorTreeOption[], (option) => {
            return -option.age.valueOf();
        });

        console.log("sensor-tree-options", { sorted: sorted });

        return sorted;
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

    public changeSensors(viz: Viz, hasParams: HasSensorParams): Workspace {
        console.log("changing sensors", viz, hasParams);
        if (viz instanceof Graph) {
            viz.changeSensors(hasParams);
        }
        return this;
    }

    public changeLinkage(viz: Viz): Workspace {
        const group = this.findGroup(viz);
        if (group.vizes.length > 1) {
            group.remove(viz);
            this.groups.unshift(new Group([viz]));
        } else {
            this.groups.reduce((previous: Group | null, iter: Group): Group => {
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
        if (removing) {
            this.groups[0].addAll(removing);
        }
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

    public async updateFromBookmark(bm: Bookmark): Promise<void> {
        if (Bookmark.sameAs(this.bookmark(), bm)) {
            return;
        }
        this.groups = bm.g.map((gm) => Group.fromBookmark(gm));
        await this.query();
        return;
    }

    private eventually(callback: (ws: Workspace) => Promise<any>) {
        callback(this);
        return Promise.resolve(this);
    }

    public with(callback: (ws: Workspace) => Workspace) {
        callback(this);
        return this;
    }
}

export class BookmarkFactory {
    public static forStation(stationId: number): Bookmark {
        return new Bookmark(Bookmark.Version, [], [stationId]);
    }
}
