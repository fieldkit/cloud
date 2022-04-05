import _, { map } from "lodash";
import moment, { Moment } from "moment";
import {
    SensorsResponse,
    SensorDataResponse,
    SensorInfoResponse,
    ModuleID,
    SensorSpec,
    Ids,
    TimeRange,
    StationID,
    Stations,
    Sensors,
    SensorParams,
    DataQueryParams,
    VizSensor,
    VizInfo,
    DataSetSeries,
    SeriesData,
    QueriedData,
    ExploreContext,
} from "./common";
import i18n from "@/i18n";
import FKApi from "@/api/api";

export * from "./common";

import { promiseAfter } from "@/utilities";
import { createSensorColorScale } from "./d3-helpers";

type SensorReadAtType = string;

function getString(d) {
    return d["enUS"] || d["enUs"] || d["en-US"]; // HACK
}

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
    constructor(
        public readonly id: number,
        public readonly name: string,
        public readonly location: [number, number],
        public readonly sensors: SensorMeta[]
    ) {
        if (!this.id) throw new Error("id is required");
        if (!this.name) throw new Error("name is required");
        if (!this.sensors) throw new Error("sensors is required");
    }
}

export interface HasSensorParams {
    readonly sensorParams: SensorParams;
}

export class StationTreeOption {
    constructor(public readonly id: string | number, public readonly label: string, public readonly isDisabled: boolean) {}
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

type VizBookmark = [VizSensor[], [number, number], [[number, number], [number, number]] | [], ChartType, FastTime];

type GroupBookmark = VizBookmark[][];

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

export interface VizInfoFactory {
    vizInfo(viz: Viz, ds: DataSetSeries): VizInfo;
}

type ResolveData = (qd: QueriedData) => void;

class VizQuery {
    constructor(public readonly params: DataQueryParams, public readonly vizes: Viz[], private readonly r: ResolveData) {}

    public howBusy(d: number): any {
        return this.vizes.map((v) => v.howBusy(d));
    }

    public resolve(qd: QueriedData) {
        return this.r(qd);
    }
}

class InfoQuery {
    constructor(public readonly params: Stations, public readonly vizes: Viz[]) {}

    public howBusy(d: number): any {
        return this.vizes.map((v) => v.howBusy(d));
    }
}

export class Scrubber {
    constructor(public readonly index: number, public readonly data: QueriedData, public readonly viz: Viz) {}
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

export class Bookmark {
    static Version = 1;

    constructor(
        public readonly v: number,
        public readonly g: GroupBookmark[],
        public readonly s: number[] = [],
        public readonly p: number[] = []
    ) {}

    private get allVizes(): VizBookmark[] {
        return _.flatten(_.flatten(this.g.map((group) => group.map((vizes) => vizes))));
    }

    public get allTimeRange(): TimeRange {
        const times: [number, number][] = this.allVizes.map((viz) => viz[1]).filter((times) => times !== undefined) as [number, number][];
        const start = _.min(_.flatten(times.map((r) => r[0])));
        const end = _.max(_.flatten(times.map((r) => r[1])));
        if (!start || !end) throw new Error(`no time range in bookmark`);
        return new TimeRange(start, end);
    }

    private vizStations(vizBookmark: VizBookmark): StationID[] {
        return [];
    }

    private vizSensors(vizBookmark: VizBookmark): SensorSpec[] {
        return [];
    }

    public get allStations(): number[] {
        return _.uniq(_.flatten(this.allVizes.map((viz) => this.vizStations(viz))));
    }

    public get allSensors(): SensorSpec[] {
        return _.uniq(_.flatten(this.allVizes.map((viz) => this.vizSensors(viz))));
    }

    public static sameAs(a: Bookmark, b: Bookmark): boolean {
        const aEncoded = JSON.stringify(a);
        const bEncoded = JSON.stringify(b);
        return aEncoded == bEncoded;
    }
}

type LegacyVizBookmark = [Stations, Sensors, [number, number], [[number, number], [number, number]] | [], ChartType, FastTime];
type LegacyGroupBookmark = LegacyVizBookmark[][];
type PossibleBookmarks = {
    v: number;
    s: number[];
    p: number[] | undefined;
    g: LegacyGroupBookmark[] | GroupBookmark[];
};

function migrateBookmark(raw: PossibleBookmarks): { v: number; s: number[]; g: GroupBookmark[]; p: number[] } {
    return {
        v: raw.v,
        s: raw.s,
        p: raw.p || [],
        g: raw.g.map((g1) =>
            g1.map((g2) =>
                g2.map((v) => {
                    if (v.length == 6) {
                        const stations = v[0];
                        const sensors = v[1];
                        const fixed = _.concat([[[stations[0], sensors[0]]]], _.drop(v, 2));
                        console.log("viz: migrate:b", v);
                        console.log("viz: migrate:a", fixed);
                        return fixed;
                    }
                    console.log("viz: migrate", v);
                    return v;
                })
            )
        ),
    };
}

export function deserializeBookmark(s: string): Bookmark {
    const migrated = migrateBookmark(JSON.parse(s));
    return Object.assign(new Bookmark(1, []), migrated);
}

export function serializeBookmark(b: Bookmark): string {
    return JSON.stringify(b);
}

export class GeoZoom {
    constructor(public readonly bounds: [[number, number], [number, number]]) {}
}

export class TimeZoom {
    constructor(public readonly fast: FastTime | null, public readonly range: TimeRange | null) {}
}

export class NewParams implements HasSensorParams {
    constructor(public readonly sensorParams: SensorParams) {}
}

export class Graph extends Viz {
    public all: QueriedData | null = null;
    public visible: TimeRange = TimeRange.eternity;
    public queried: TimeRange = TimeRange.eternity;
    public chartType: ChartType = ChartType.TimeSeries;
    public fastTime: FastTime = FastTime.All;
    public geo: GeoZoom | null = null;

    constructor(public readonly when: TimeRange, public dataSets: DataSetSeries[]) {
        super();
    }

    public get loadedDataSets(): DataSetSeries[] {
        if (this.dataSets.length) {
            const all = _.flatten(
                this.dataSets.map((ds) => {
                    if (ds && ds.graphing) {
                        return [ds];
                    }
                    return [];
                })
            );

            if (all.length == this.dataSets.length) {
                return all;
            }
        }
        return [];
    }

    public allSeries(vizInfoFactory: VizInfoFactory): SeriesData[] {
        return this.loadedDataSets.map((ds) => {
            if (!ds.graphing) throw new Error(`viz: No data`);
            const vizInfo = vizInfoFactory.vizInfo(this, ds);
            return new SeriesData(ds.graphing.key, ds, ds.graphing, vizInfo);
        });
    }

    public get visibleTimeRange(): TimeRange {
        if (this.visible.isExtreme()) {
            if (this.all) {
                return new TimeRange(this.all.timeRange[0], this.all.timeRange[1]);
            }
        }
        return this.visible;
    }

    public get allStationIds(): StationID[] {
        return this.dataSets.map((ds) => ds.stationId);
    }

    public allQueries(): VizQuery[] {
        return this.dataSets.map((ds) => {
            const params = new DataQueryParams(this.visible, [ds.vizSensor]);
            return new VizQuery(params, [this], (qd) => {
                ds.graphing = qd;
                return;
            });
        });
    }

    public clone(): Viz {
        const c = new Graph(this.when, this.dataSets);
        c.all = this.all;
        c.visible = this.visible;
        c.chartType = this.chartType;
        c.fastTime = this.fastTime;
        return c;
    }

    public get scrubberParams(): DataQueryParams {
        return new DataQueryParams(
            TimeRange.eternity,
            _.take(
                this.dataSets.map((ds) => ds.vizSensor),
                1
            )
        );
    }

    public timeZoomed(zoom: TimeZoom): TimeRange {
        if (zoom.range !== null) {
            this.visible = zoom.range;
            this.fastTime = FastTime.Custom;
        } else if (zoom.fast !== null) {
            this.visible = this.getFastRange(zoom.fast);
            this.fastTime = zoom.fast;
        }

        return this.visible;
    }

    public geoZoomed(zoom: GeoZoom): GeoZoom {
        this.geo = zoom;
        return this.geo;
    }

    private getFastRange(fastTime: FastTime) {
        if (!this.all) throw new Error("viz: Fast time missing all data");

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

    public changeSensors(hasParams: HasSensorParams) {
        const sensorParams = hasParams.sensorParams;
        this.log(`changing-sensors`, hasParams);
        const stationsBefore = this.dataSets.map((ds) => ds.stationId);
        this.dataSets = sensorParams.sensors.map((vizSensor) => new DataSetSeries(vizSensor));
        const stationsAfter = this.dataSets.map((ds) => ds.stationId);

        if (_.difference(stationsBefore, stationsAfter)) {
            if (this.geo) {
                this.log(`changing-sensors`, "clearing-geo");
                this.geo = null;
            }
        }

        this.all = null;
    }

    public bookmark(): VizBookmark {
        return [
            this.dataSets.map((ds) => ds.bookmark()),
            [this.visible.start, this.visible.end],
            this.geo ? this.geo.bounds : [],
            this.chartType,
            this.fastTime,
        ];
    }

    public modifySeries(index: number, vs: VizSensor): NewParams {
        const updated = this.dataSets.map((ds) => ds.vizSensor);
        updated[index] = vs;
        return new NewParams(new SensorParams(updated));
    }

    public addSeries(): NewParams {
        const updated = this.dataSets.map((ds) => ds.vizSensor);
        updated.push(updated[0]);
        return new NewParams(new SensorParams(updated));
    }

    public removeSeries(index: number): NewParams {
        const updated = this.dataSets.map((ds) => ds.vizSensor);
        updated.splice(index, 1);
        return new NewParams(new SensorParams(updated));
    }

    public static fromBookmark(bm: VizBookmark): Viz {
        const visible = new TimeRange(bm[1][0], bm[1][1]);
        const dataSets = bm[0].map((vizSensor) => new DataSetSeries(vizSensor));
        const graph = new Graph(visible, dataSets);
        graph.geo = bm[2].length ? new GeoZoom(bm[2]) : null;
        graph.chartType = bm[3];
        graph.fastTime = bm[4];
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

    public cloneForCompare(): Group {
        return new Group([this.vizes[0]].map((v) => v.clone()));
    }

    public unlinkAt(viz: Viz): Group {
        const index = _.indexOf(this.vizes, viz);
        if (index < 0) throw new Error("viz: Unlinking of mismatched group/viz");
        console.log("unlink-at:index", index);
        const removing = this.vizes.slice(index);
        console.log("unlink-at:removing", removing);
        this.vizes = this.vizes.filter((viz) => {
            return _.indexOf(removing, viz) < 0;
        });
        console.log("unlink-at:after", this.vizes);
        return new Group(removing);
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

    public get busy() {
        return this.vizes.filter((v) => v.busy).length > 0;
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
                return new Scrubber(r.index, all, r.graph);
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

export class Querier {
    private info: { [index: string]: SensorInfoResponse } = {};
    private data: { [index: string]: QueriedData } = {};

    public queryInfo(iq: InfoQuery): Promise<SensorInfoResponse> {
        if (!iq.params) throw new Error("no params");

        const params = iq.params;
        const queryParams = new URLSearchParams();
        queryParams.append("stations", params.join(","));

        const key = queryParams.toString();

        console.log(`vis: query-info`, key);

        if (this.info[key]) {
            // iq.howBusy(1);

            return promiseAfter(1).then(() => {
                // iq.howBusy(-1);
                return this.info[key];
            });
        }

        return Promise.resolve()
            .then(() => {
                iq.howBusy(1);
            })
            .then(() => {
                return new FKApi()
                    .sensorData(queryParams)
                    .then((info: SensorInfoResponse) => {
                        this.info[key] = info;
                        return info;
                    })
                    .finally(() => {
                        iq.howBusy(-1);
                    });
            });
    }

    public queryData(vq: VizQuery): Promise<QueriedData> {
        if (!vq.params) throw new Error("no params");

        const params = vq.params;
        const queryParams = params.queryParams();
        const key = queryParams.toString();

        console.log(`viz: query-data`, key);

        if (this.data[key]) {
            // vq.howBusy(1);

            return promiseAfter(1).then(() => {
                vq.resolve(this.data[key]);
                // vq.howBusy(-1);
                return this.data[key];
            });
        }

        return Promise.resolve()
            .then(() => {
                vq.howBusy(1);
            })
            .then(() => {
                return new FKApi()
                    .sensorData(queryParams)
                    .then((sdr: SensorDataResponse) => {
                        const queried = new QueriedData(key, params.when, sdr);
                        const filtered = queried./*removeMalformed().*/ removeDuplicates();
                        this.data[key] = filtered;
                        return filtered;
                    })
                    .then((data) => {
                        vq.resolve(data);
                        return data;
                    })
                    .finally(() => {
                        vq.howBusy(-1);
                    });
            });
    }
}

export class Workspace implements VizInfoFactory {
    private stationIds: StationID[] = [];
    private readonly querier = new Querier();
    private readonly stations: { [index: number]: StationMeta } = {};
    public version = 0;

    public get empty(): boolean {
        return this.allVizes.length === 0;
    }

    public get busy(): boolean {
        return this.groups.filter((g) => g.busy).length > 0;
    }

    public get allStationIds(): StationID[] {
        return this.stationIds;
    }

    constructor(private readonly meta: SensorsResponse, private groups: Group[] = [], public readonly projects: number[] = []) {
        this.refreshStationIds();
    }

    private get allVizes(): Viz[] {
        return _(this.groups)
            .map((g) => g.vizes)
            .flatten()
            .value();
    }

    public async query(): Promise<any> {
        const allGraphs = this.allVizes.map((viz) => viz as Graph).filter((viz) => viz);

        // First we may need some data for making the UI useful and
        // pretty. Filling in labels, building drop downs, etc... This
        // is especially important to do from here because we may have
        // been instantiated from a Bookmark. Right now we just query
        // for information on all the stations involved.
        const allStationIds = _.uniq(_.flatten(allGraphs.map((viz) => viz.allStationIds)).concat(this.stationIds));
        const infoQueries = allStationIds.length ? [new InfoQuery(allStationIds, allGraphs)] : [];

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
        const graphingQueries = _.flatten(allGraphs.map((viz: Graph) => viz.allQueries()));

        // Combine and make them unique to avoid obvious
        // duplicates. Eventually we can also merge stations/sensors
        // with the same date range and other parameters into one
        // query. Lots of room here.
        const allQueries = [...scrubberQueries, ...graphingQueries];
        const uniqueQueries = _(allQueries)
            .groupBy((q) => q.params.queryParams().toString())
            .map((p) => new VizQuery(p[0].params, _.flatten(p.map((p) => p.vizes)), (qd: QueriedData) => p.map((p) => p.resolve(qd))))
            .value();

        console.log("viz: workspace: querying", uniqueQueries.length, "data", infoQueries.length, "info");

        // Make all the queries and then give the queried data to the
        // resolve call for that query. This will end up calling the
        // above mapped resolve to set the appropriate data.
        const pendingInfo = infoQueries.map(
            (iq) =>
                this.querier.queryInfo(iq).then((info) => {
                    return _.map(info.stations, (info, stationId) => {
                        const stationName = info[0].stationName;
                        const stationLocation = info[0].stationLocation;
                        const sensors = info
                            .filter((row) => row.moduleId != null)
                            .map((row) => new SensorMeta(row.moduleId, row.moduleKey, row.sensorId, row.sensorKey, row.sensorReadAt));
                        const station = new StationMeta(Number(stationId), stationName, stationLocation, sensors);
                        this.stations[station.id] = station;
                        console.log("viz: station-meta", { station, info });
                        return station;
                    });
                }) as Promise<unknown>
        );

        await Promise.all([...pendingInfo]);

        const pendingData = uniqueQueries.map((vq) => this.querier.queryData(vq) as Promise<unknown>);
        return Promise.all([...pendingInfo, ...pendingData]).then(() => {
            // Update options here if doing so lazily.
            console.log("viz: workspace: query done ");
            this.version++;
        });
    }

    public vizInfo(viz: Graph, ds: DataSetSeries | undefined = undefined): VizInfo {
        // if (viz.chartParams.stations.length != 1) throw new Error("expected 1 station per graph, for now");
        // if (viz.chartParams.sensors.length != 1) throw new Error("expected 1 sensor per graph, for now");

        const sensorDetailsByKey = _.keyBy(_.flatten(this.meta.modules.map((m) => m.sensors)), (s) => s.fullKey);
        const keysById = _.keyBy(this.meta.sensors, (s) => s.id);

        if (!ds) {
            if (viz.dataSets.length == 0) throw new Error("viz: No DataSets");
            ds = viz.dataSets[0];
        }

        // console.log(`viz:vizInfo:ds`, ds);

        const stationId = ds.stationId;
        const sensor = ds.sensorAndModule;
        if (!sensor || sensor.length != 2) throw new Error("viz: Malformed SensorAndModule");
        const sensorId = sensor[1];

        const station = this.stations[stationId];
        const key = keysById[sensorId].key;
        const details = sensorDetailsByKey[key];
        const scale = createSensorColorScale(details);

        // console.log(`viz:vizInfo:sensor`, details);

        const strings = getString(details.strings);

        // console.log(`viz:vizInfo:sensor`, strings);

        return new VizInfo(key, scale, station, details.unitOfMeasure, key, strings.label, details.viz || [], details.ranges);
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

    public async addStationIds(ids: number[]): Promise<Workspace> {
        if (_.difference(ids, this.stationIds).length == 0) {
            console.log("viz: workspace-add-station-ids(ignored)", ids);
            return this;
        }
        this.stationIds = [...this.stationIds, ...ids];
        console.log("viz: workspace-add-station-ids", this.stationIds);
        return this;
    }

    private findGroup(viz: Viz): Group {
        const groups = this.groups.filter((g) => g.contains(viz));
        if (groups.length == 1) {
            return groups[0];
        }
        throw new Error("viz: Orphaned viz");
    }

    public addGraph(graph: Graph): Workspace {
        const group = new Group();
        group.add(graph);
        this.groups.unshift(group);
        this.refreshStationIds();
        return this;
    }

    private refreshStationIds() {
        this.stationIds = _.uniq(_.flatten(_.flatten(this.groups.map((g) => g.vizes.map((v) => (v as Graph).allStationIds)))));
    }

    public addStandardGraph(vizSensor: VizSensor): Workspace {
        return this.addGraph(new Graph(TimeRange.eternity, [new DataSetSeries(vizSensor)]));
    }

    public get stationOptions(): StationTreeOption[] {
        return Object.values(this.stations).map((station) => {
            return new StationTreeOption(station.id, station.name, station.sensors.length == 0);
        });
    }

    public sensorOptions(stationId: number): SensorTreeOption[] {
        const station = this.stations[stationId];
        if (!station) throw new Error("viz: No station");
        const allSensors = station.sensors;
        const allModules = _.groupBy(allSensors, (s) => s.moduleId);
        const keysById = _.fromPairs(allSensors.map((row) => [row.moduleId, row.moduleKey]));
        const allModulesByModuleKey = _.keyBy(this.meta.modules, (m) => m.key);

        // console.log("all-sensors", allSensors);
        // console.log("all-modules", allModules);
        // console.log("keys-by-id", keysById);
        // console.log("all-modules-by-module-key", allModulesByModuleKey);
        // console.log(station.sensors);

        const options = _.map(
            allModules,
            (sensors, moduleId: ModuleID): SensorTreeOption => {
                const moduleKey = keysById[moduleId];
                const moduleMeta = allModulesByModuleKey[moduleKey];
                const uniqueSensors = _.uniqBy(sensors, (s) => s.sensorId);
                const children: SensorTreeOption[] = _.flatten(
                    uniqueSensors.map((row) => {
                        const age = moment.utc(row.sensorReadAt);
                        const label = i18n.tc(row.sensorKey) || row.sensorKey;
                        const optionId = `${row.moduleId}-${row.sensorId}`;
                        const sensor = moduleMeta.sensors.filter((s) => s.fullKey == row.sensorKey);
                        if (sensor.length) {
                            if (sensor[0].internal) {
                                return [];
                            }
                        }
                        return [new SensorTreeOption(optionId, label, undefined, row.moduleId, row.sensorId, age)];
                    })
                );
                const moduleAge = _.max(children.map((c) => c.age));
                if (!moduleAge) throw new Error(`viz: Expected module age: no sensors?`);

                const label = i18n.tc(moduleKey); //  + ` (${moduleAge.fromNow()})`;
                return new SensorTreeOption(`${moduleKey}-${moduleId}`, label, children, moduleId, null, moduleAge);
            }
        );

        const sorted = _.sortBy((options as unknown) as SensorTreeOption[], (option) => {
            return -option.age.valueOf();
        });

        return sorted;
    }

    public makeSeries(stationId: number, sensorAndModule: SensorSpec): DataSetSeries {
        const station = this.stations[stationId];
        if (!station) throw new Error(`viz: No station with id: ${stationId}`);
        const sensors = station.sensors;
        if (sensors.length == 0) throw new Error(`viz: No sensors on station with id: ${stationId}`);
        const moduleIds = sensors.map((s) => s.moduleId);
        if (moduleIds.includes(sensorAndModule[0])) {
            return new DataSetSeries([stationId, sensorAndModule]);
        }
        const fallbackSensorAndModule: SensorSpec = [sensors[0].moduleId, sensors[0].sensorId];
        return new DataSetSeries([stationId, fallbackSensorAndModule]);
    }

    public remove(viz: Viz): Workspace {
        this.groups = this.groups.map((g) => g.remove(viz)).filter((g) => !g.empty);
        return this;
    }

    public addChart(): Workspace {
        if (this.groups.length > 0) {
            const adding = this.groups[0].cloneForCompare();
            console.log("viz-compare", adding);
            this.groups.unshift(adding);
        }
        return this;
    }

    public compare(viz: Viz): Workspace {
        return this.addChart();
    }

    public changeChart(viz: Viz, chartType: ChartType): Workspace {
        if (viz instanceof Graph) {
            viz.changeChart(chartType);
        }
        return this;
    }

    public changeSensors(viz: Viz, hasParams: HasSensorParams): Workspace {
        if (viz instanceof Graph) {
            viz.changeSensors(hasParams);
        }
        return this;
    }

    public changeLinkage(viz: Viz, linking: boolean): Workspace {
        const group = this.findGroup(viz);
        if (linking) {
            console.log("linkage-link", viz.id, group.id, viz, group);
            this.groups.reduce((previous: Group | null, iter: Group): Group => {
                if (iter == group) {
                    if (previous == null) throw new Error("viz: Tried linking first group, nice work");
                    this.removeGroup(group);
                    previous.addAll(group);
                }
                return iter;
            }, null);
        } else {
            console.log("linkage-unlink", viz.id, group.id, viz, group);
            const newGroup = group.unlinkAt(viz);
            const groupIndex = _.indexOf(this.groups, group);
            this.groups.splice(groupIndex + 1, 0, newGroup);
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
            this.groups.map((group) => group.bookmark()),
            this.allStationIds,
            this.projects
        );
    }

    public static fromBookmark(meta: SensorsResponse, bm: Bookmark): Workspace {
        if (bm.v !== 1) {
            throw new Error("viz: Unexpected bookmark version");
        }
        return new Workspace(
            meta,
            bm.g.map((gm) => Group.fromBookmark(gm)),
            bm.p
        );
    }

    public async updateFromBookmark(bm: Bookmark): Promise<void> {
        if (Bookmark.sameAs(this.bookmark(), bm)) {
            console.log(`viz: update-from-bookmark:same`, bm);
            return;
        }

        console.log(`viz: update-from-bookmark`, bm);
        await this.addStationIds(bm.s);
        this.groups = bm.g.map((gm) => Group.fromBookmark(gm));
        await this.query();
        return;
    }

    public with(callback: (ws: Workspace) => Workspace) {
        callback(this);
        return this;
    }
}

export class BookmarkFactory {
    public static forStation(stationId: number, context: ExploreContext | null = null): Bookmark {
        if (context && context.project) {
            return new Bookmark(Bookmark.Version, [], [stationId], [context.project]);
        }
        return new Bookmark(Bookmark.Version, [], [stationId], []);
    }
}
