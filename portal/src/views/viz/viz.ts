import _ from "lodash";
import moment, { Moment } from "moment";
import {
    SensorsResponse,
    SensorDataResponse,
    SensorInfoResponse,
    Time,
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
import FKApi, { Station, AssociatedStation } from "@/api/api";
import i18n from "@/i18n";

export * from "./common";

import { promiseAfter } from "@/utilities";
import { createSensorColorScale } from "./d3-helpers";
import { DisplayStation } from "@/store";
import { getPartnerCustomizationWithDefault } from "../shared/partners";

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
    constructor(
        public readonly id: string | number,
        public readonly label: string,
        public readonly isDisabled: boolean,
        public readonly isDefaultExpanded: boolean = false,
        public readonly children: StationTreeOption[] | undefined = undefined
    ) {}
}

export class SensorTreeOption {
    constructor(
        public readonly id: string | number,
        public readonly label: string,
        public children: SensorTreeOption[] | undefined = undefined, // TODO HACK
        public readonly moduleId: ModuleID | null,
        public readonly sensorId: number | null,
        public readonly stationId: number | null,
        public readonly age: Moment | null
    ) {}

    public reassignStation(stationId: number): SensorTreeOption {
        return new SensorTreeOption(this.id, this.label, this.children, this.moduleId, this.sensorId, stationId, this.age);
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
    Default = FastTime.TwoWeeks,
}

export enum ChartType {
    TimeSeries,
    Histogram,
    Range,
    Map,
    Bar,
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
        public readonly p: number[] = [],
        public readonly c: ExploreContext | null = null
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
    c: ExploreContext | null;
};

function migrateBookmark(raw: PossibleBookmarks): { v: number; s: number[]; g: GroupBookmark[]; p: number[]; c: ExploreContext | null } {
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
        c: raw.c,
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
                    if (ds && ds.graphing && ds.all) {
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

    public get timeRangeOfAll(): TimeRange | null {
        const everyAllRange = this.dataSets.map((ds) => ds.all?.timeRange).filter((range): range is number[] => range != null);
        if (everyAllRange.length > 0) {
            return TimeRange.mergeArrays(everyAllRange);
        }
        return null;
    }

    public get visibleTimeRange(): TimeRange {
        const range = this.timeRangeOfAll;
        const visible = this.visible;
        if (range && visible.isExtreme()) {
            const start = visible.start == Time.Min ? range.start : visible.start;
            const end = visible.end == Time.Max ? range.end : visible.end;
            return new TimeRange(start, end);
        }
        return this.visible;
    }

    public get allStationIds(): StationID[] {
        return this.dataSets.map((ds) => ds.stationId);
    }

    public graphingQueries(): VizQuery[] {
        return this.dataSets.map((ds) => {
            const params = new DataQueryParams(this.visible, [ds.vizSensor]);
            return new VizQuery(params, [this], (qd) => {
                ds.graphing = qd;
                return;
            });
        });
    }

    public scrubberQueries(): VizQuery[] {
        return this.dataSets.map((ds) => {
            const params = new DataQueryParams(TimeRange.eternity, [ds.vizSensor]);
            return new VizQuery(params, [this], (qd) => {
                ds.all = qd;
                return;
            });
        });
    }

    public clone(): Viz {
        const c = new Graph(this.when, this.dataSets);
        c.visible = this.visible;
        c.chartType = this.chartType;
        c.fastTime = this.fastTime;
        return c;
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
        if (fastTime === FastTime.All) {
            return TimeRange.eternity;
        } else {
            const sensorRange = this.timeRangeOfAll;
            if (sensorRange == null) throw new Error(`viz: No timeRangeOfAll`);
            const days = fastTime as number;
            const start = new Date(sensorRange.end);
            start.setDate(start.getDate() - days);
            return new TimeRange(start.getTime(), Time.Max);
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
            .filter((r) => r.graph.dataSets[0].all != null); // TODO Ignoring others

        const childScrubbers = children.map((r) => {
            const all = r.graph.dataSets[0].all; // TODO Ignoring others
            if (!all) throw new Error(`no viz data on Graph`);
            return new Scrubber(r.index, all, r.graph);
        });

        const mergedVisible = TimeRange.mergeArrays(children.map((v) => v.graph.visibleTimeRange.toArray()));
        console.log("viz: scrubbers", this.visible_.toArray(), mergedVisible);
        return new Scrubbers(this.id, mergedVisible, childScrubbers);
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

        // console.log(`vis: query-info`, key);

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

        // console.log(`viz: query-data`, key);

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
    public version = 0;
    private stationIds: StationID[] = [];
    private stationsFull: Station[] = [];
    private associated: AssociatedStation[] = [];
    private readonly querier = new Querier();
    private readonly stations: { [index: number]: StationMeta } = {};
    private readonly showInternalSensors = false;
    private ready = false;

    constructor(
        private readonly meta: SensorsResponse,
        private groups: Group[] = [],
        public readonly projects: number[] = [],
        public readonly bookmarkStations: number[] | null = null,
        public readonly context: ExploreContext | null = null
    ) {
        if (bookmarkStations) {
            this.stationIds = bookmarkStations;
        } else {
            this.refreshStationIds();
        }
    }

    public get selectedStationId(): number {
        const firstVizSensor = _(this.groups)
            .map((g) => g.vizes.map((v) => v.bookmark()[0]))
            .flatten()
            .flatten()
            .first();
        if (!firstVizSensor) {
            throw new Error();
        }
        console.log("viz: first-viz-sensor", firstVizSensor);
        return this.findStationOverride(firstVizSensor) || firstVizSensor[0];
    }

    public findStationOverride(sensor: VizSensor): number | null {
        const fromStations = _.fromPairs(
            _.flatten(Object.values(this.stations).map((row) => row.sensors.map((sensor) => [sensor.moduleId, row.id])))
        );
        const moduleIdToStationId = _(this.associated)
            .map((a) => {
                return a.station.configurations.all[0].modules.map((m) => {
                    return [m.hardwareIdBase64, a.station.id];
                });
            })
            .flatten()
            .fromPairs()
            .merge(fromStations)
            .value();
        const moduleId = sensor[1][0];
        const sensorStationId = moduleIdToStationId[moduleId];
        return this.stationOverrides[sensorStationId];
    }

    private get stationOverrides(): { [index: number]: number } {
        return _.fromPairs(
            this.associated
                .map((associated) => {
                    if (associated.manual) {
                        return [associated.station.id, associated.manual.otherStationID];
                    }
                    return [];
                })
                .filter((l) => l.length)
        );
    }

    public get empty(): boolean {
        return this.allVizes.length === 0;
    }

    public get busy(): boolean {
        return this.groups.filter((g) => g.busy).length > 0;
    }

    public get allStationIds(): StationID[] {
        return this.stationIds;
    }

    public get allStations(): Station[] {
        return this.stationsFull;
    }

    public getStation(id: number): DisplayStation | null {
        if (this.stationsFull) {
            const found = this.stationsFull.filter((d) => d.id === id);
            if (found.length > 0) {
                return new DisplayStation(found[0]);
            }
        }
        return null;
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
        const scrubberQueries = _.flatten(allGraphs.map((viz: Graph) => viz.scrubberQueries()));
        // Now build the queries for the data being viewed.
        const graphingQueries = _.flatten(allGraphs.map((viz: Graph) => viz.graphingQueries()));

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
                        // console.log("viz: station-meta", { station, info });
                        return station;
                    });
                }) as Promise<unknown>
        );

        await Promise.all([...pendingInfo]);

        const pendingData = uniqueQueries.map((vq) => this.querier.queryData(vq) as Promise<unknown>);
        return Promise.all([...pendingInfo, ...pendingData]).then(() => {
            // Update options here if doing so lazily.
            // console.log("viz: workspace: query done ");
            this.version++;
        });
    }

    public availableChartTypes(viz: Graph, ds: DataSetSeries | undefined = undefined): ChartType[] {
        const vizInfo = this.vizInfo(viz, ds);
        const specifiedNames = vizInfo.viz.map((row) => row.name);
        if (specifiedNames.length == 0) {
            return [ChartType.TimeSeries, ChartType.Histogram, ChartType.Range, ChartType.Bar, ChartType.Map];
        }

        const knownNames = {
            TimeSeriesChart: ChartType.TimeSeries,
            HistogramChart: ChartType.Histogram,
            BarChart: ChartType.Bar,
            RangeChart: ChartType.Range,
            Map: ChartType.Map,
        };

        const migratedNames = specifiedNames.map((name) => name.replace("D3", "").replace("Graph", "Chart"));
        return migratedNames.map((row) => knownNames[row]);
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
        const chartLabel = (strings.chartLabel) ? strings.chartLabel : strings.label;
        const axisLabel = (strings.axisLabel) ? strings.axisLabel : strings.label;

        // console.log(`viz:vizInfo:sensor`, strings);

        return new VizInfo(key, scale, station, details.unitOfMeasure, key, strings.label, details.viz || [], details.ranges, chartLabel, axisLabel);
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

    public async addAssociatedStations(associated: AssociatedStation[]): Promise<Workspace> {
        // console.log("viz: associated-add", { associated });
        this.associated = associated;
        return this;
    }

    public async addFullStations(stations: Station[]): Promise<Workspace> {
        this.stationsFull = stations;
        return this;
    }

    public async addStationIds(ids: number[]): Promise<Workspace> {
        if (_.difference(ids, this.stationIds).length == 0) {
            console.log("viz: workspace-add-station-ids(ignored)", { ids });
            return this;
        }
        this.stationIds = [...this.stationIds, ...ids];
        console.log("viz: workspace-add-station-ids", { ids: this.stationIds });
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
        const vizStationIds = _.uniq(_.flatten(_.flatten(this.groups.map((g) => g.vizes.map((v) => (v as Graph).allStationIds)))));
        this.stationIds = vizStationIds;
    }

    public addStandardGraph(vizSensor: VizSensor): Workspace {
        return this.addGraph(new Graph(TimeRange.eternity, [new DataSetSeries(vizSensor)]));
    }

    private nearbyStationOptions(): StationTreeOption[] {
        // Notice we skip over manually associated stations as that's a stronger association.
        const nearby = this.associated.filter((assoc) => this.stations[assoc.station.id] && assoc.location && !assoc.manual);
        console.log("viz: associated-nearby", nearby);
        if (nearby.length == 0) {
            return [];
        }

        const options = nearby.map((associatedStation) => {
            const station = associatedStation.station;
            return new StationTreeOption(station.id, station.name, false);
        });

        return [new StationTreeOption(`nearby`, "Nearby", false, true, options)];
    }

    private manuallyAssociatedStationOptions(): StationTreeOption[] {
        const manuallyAssociated = this.associated.filter((assoc) => this.stations[assoc.station.id] && assoc.manual);
        console.log("viz: associated-manually", manuallyAssociated);
        if (manuallyAssociated.length == 0) {
            return [];
        }

        const options = manuallyAssociated.map((associatedStation) => {
            const station = associatedStation.station;
            return new StationTreeOption(station.id, station.name, false);
        });

        return [new StationTreeOption(`manual`, "Associated", false, false, options)];
    }

    public get stationOptions(): StationTreeOption[] {
        const hiddenById = _.fromPairs(this.associated.map((assoc) => [assoc.station.id, assoc.hidden]));
        const associatedById = _.groupBy(this.associated, (assoc) => assoc.station.id);
        const nearby = this.nearbyStationOptions();
        const manually = []; // this.manuallyAssociatedStationOptions();

        // This is for removing stations that already have an option because of
        // they're associated. Not a fan of this approach.
        const unassociated = Object.values(this.stations).filter((station) => {
            if (hiddenById[station.id]) {
                return false;
            }

            const maybeAssociated = associatedById[station.id];
            if (maybeAssociated && maybeAssociated.length > 0) {
                return !maybeAssociated[0].location && !maybeAssociated[0].manual;
            }

            return true;
        });

        const regular = unassociated.map((station) => {
            return new StationTreeOption(station.id, station.name, station.sensors.length == 0);
        });

        const partnerCustomization = getPartnerCustomizationWithDefault();

        const grouped = _(regular)
            .filter((option) => _.isNumber(option.id))
            .map((stationOption) => {
                const station = this.getStation(Number(stationOption.id));
                if (station) {
                    return [
                        {
                            option: stationOption,
                            station: station,
                        },
                    ];
                }
                return [];
            })
            .flatten()
            .map((row) => {
                return {
                    option: row.option,
                    group: partnerCustomization.viz.groupStation(row.station),
                };
            })
            .value();

        const ungrouped = grouped.filter((row) => row.group == null).map((row) => row.option);

        const groupOptions = _(grouped)
            .filter((row) => row.group != null)
            .groupBy((row) => row.group)
            .map((group, name) => {
                return new StationTreeOption(
                    `group-${name}`,
                    name,
                    false,
                    false,
                    group.map((child) => child.option)
                );
            })
            .value();

        const allOptions = [...groupOptions, ...ungrouped];
        const all = allOptions.length > 0 ? [new StationTreeOption(`all`, "All", false, false, allOptions)] : [];

        // console.log("viz: ungrouped", { ungrouped });
        // console.log("viz: all", { all });

        return [...manually, ...nearby, ...all];
    }

    public sensorOptions(stationId: number, flatten = false, depth = 0): SensorTreeOption[] {
        const station = this.stations[stationId];
        if (!station) {
            throw new Error(`viz: No station: ${stationId}`);
        }

        if (this.associated.length == 0) {
            throw new Error("viz: Associated required for sensor-options");
        }

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
                        let label = i18n.tc(row.sensorKey) || row.sensorKey;
                        if (flatten) {
                            label = moduleMeta.sensors.filter((d) => d.fullKey === row.sensorKey)[0]["strings"]["enUs"]["label"];
                        }
                        const optionId = `${row.moduleId}-${row.sensorId}`;
                        const sensor = moduleMeta.sensors.filter((s) => s.fullKey == row.sensorKey);
                        if (sensor.length > 0) {
                            if (!this.showInternalSensors) {
                                if (sensor[0].internal) {
                                    return [];
                                }
                            }
                        }
                        return [new SensorTreeOption(optionId, label, undefined, row.moduleId, row.sensorId, stationId, age)];
                    })
                );
                const moduleAge = _.max(children.map((c) => c.age));
                if (!moduleAge) {
                    throw new Error(`viz: Expected module age: no sensors?`);
                }

                const label = i18n.tc(moduleKey);

                if (flatten) {
                    return children[0];
                } else {
                    return new SensorTreeOption(`${moduleKey}-${moduleId}`, label, children, moduleId, null, stationId, moduleAge);
                }
            }
        );

        const sorted = _.sortBy((options as unknown) as SensorTreeOption[], (option) => {
            if (option.age) {
                return -option.age.valueOf();
            }
            return 0;
        });

        const associatedWithStation = this.associated.filter((assoc) => assoc.manual && assoc.manual.otherStationID == stationId);
        if (associatedWithStation.length > 0) {
            const associatedSensorOptions = _(
                associatedWithStation.map((associated) => {
                    const moduleOptions = this.sensorOptions(associated.station.id, false, depth + 1);
                    // console.log("viz: debug-associated", depth, moduleOptions);
                    return _.flatten(moduleOptions.map((option) => option.children || []));
                })
            )
                .flatten()
                .sortBy((o) => o.label)
                .value();

            if (associatedSensorOptions.length > 0) {
                const relatedOption = new SensorTreeOption("related-sensors", "Related", associatedSensorOptions, null, null, 0, null);
                sorted.push(relatedOption);
            }
        }

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
            this.projects,
            this.context
        );
    }

    public static fromBookmark(meta: SensorsResponse, bm: Bookmark): Workspace {
        if (bm.v !== 1) {
            throw new Error("viz: Unexpected bookmark version");
        }
        return new Workspace(
            meta,
            bm.g.map((gm) => Group.fromBookmark(gm)),
            bm.p,
            bm.s,
            bm.c
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
            return new Bookmark(Bookmark.Version, [], [stationId], [context.project], context);
        }
        return new Bookmark(Bookmark.Version, [], [stationId], []);
    }
}
