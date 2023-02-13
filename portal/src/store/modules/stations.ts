import _ from "lodash";
import Vue from "vue";
import * as MutationTypes from "../mutations";
import * as ActionTypes from "../actions";
import { BoundingRectangle, LngLat, Location } from "../map-types";

import {
    Activity,
    Configurations,
    ModuleSensor,
    OnNoReject,
    Photos,
    Project,
    ProjectAttribute,
    ProjectFollowers,
    ProjectUser,
    Services,
    Station,
    StationModule,
    StationRegion,
    StationStatus,
    SensorsResponse,
    ModuleSensorMeta,
    VizThresholds,
    QueryRecentlyResponse,
    RecentlyAggregatedWindows,
    RecentlyAggregatedLast, Owner,
} from "@/api";

import { VizSensor, VizConfig } from "@/views/viz/viz";

import * as d3 from "d3";

export const NEED_SENSOR_META = "NEED_SENSOR_META";
export const HAVE_USER_STATIONS = "HAVE_USER_STATIONS";
export const HAVE_USER_PROJECTS = "HAVE_USER_PROJECTS";
export const HAVE_COMMUNITY_PROJECTS = "HAVE_COMMUNITY_PROJECTS";
export const PROJECT_USERS = "PROJECT_USERS";
export const PROJECT_FOLLOWS = "PROJECT_FOLLOWS";
export const PROJECT_STATIONS = "PROJECT_STATIONS";
export const PROJECT_ACTIVITY = "PROJECT_ACTIVITY";
export const STATION_UPDATE = "STATION_UPDATE";
export const STATION_CLEAR = "STATION_CLEAR";
export const PROJECT_LOADED = "PROJECT_LOADED";
export const PROJECT_UPDATE = "PROJECT_UPDATE";
export const PROJECT_DELETED = "PROJECT_DELETED";
export const SENSOR_META = "SENSOR_META";

export class SensorMeta {
    constructor(private readonly meta: SensorsResponse) {}

    public get sensors() {
        return this.meta.sensors;
    }

    public get modules() {
        return this.meta.modules;
    }

    public findSensorByKey(sensorKey: string): ModuleSensorMeta {
        const sensors = _(this.meta.modules)
            .map((m) => m.sensors)
            .flatten()
            .groupBy((s) => s.fullKey)
            .value();

        const byKey = sensors[sensorKey];
        if (byKey.length == 0) {
            throw new Error(`viz: Missing sensor meta: ${sensorKey}`);
        }

        return byKey[0];
    }

    public findSensorById(id: number): ModuleSensorMeta {
        const row = this.meta.sensors.find((row) => row.id == id);
        if (!row) {
            throw new Error(`viz: Missing sensor meta: ${id}`);
        }
        return this.findSensorByKey(row.key);
    }

    public findSensor(vizSensor: VizSensor): ModuleSensorMeta {
        const sensorId = vizSensor[1][1];

        const sensorKeysById = _(this.meta.sensors)
            .groupBy((r) => r.id)
            .value();

        if (!sensorKeysById[String(sensorId)]) {
            console.log(`viz: sensors: ${JSON.stringify(_.keys(sensorKeysById))}`);
            throw new Error(`viz: Missing sensor: ${sensorId}`);
        }

        const sensorKey = sensorKeysById[String(sensorId)][0].key;
        return this.findSensorByKey(sensorKey);
    }
}

export class DisplaySensor {
    name: string;
    unitOfMeasure: string;
    reading: number | null;
    time: number | null;
    meta: {
        viz: VizConfig[];
    };

    constructor(sensor: ModuleSensor) {
        this.name = sensor.name;
        this.meta = sensor.meta;
        this.unitOfMeasure = sensor.unitOfMeasure;
        if (sensor.reading) {
            this.reading = sensor.reading.last;
            this.time = sensor.reading.time;
        }
    }
}

export class DisplayModule {
    name: string;
    sensors: DisplaySensor[];

    constructor(module: StationModule) {
        this.name = module.name;
        this.sensors = module.sensors.map((s) => new DisplaySensor(s));
    }
}

export enum VisibleReadings {
    Current,
    Last72h,
}

export class DecoratedReading {
    constructor(public readonly value: number, public readonly color: string, public readonly thresholdLabel: string | null) {}
}

export class StationReadings {
    constructor(
        private readonly stationId: number,
        private readonly meta: SensorMeta,
        private readonly recently: RecentlyAggregatedWindows,
        private readonly last: RecentlyAggregatedLast
    ) {
        // console.log("station-readings:ctor", stationId, recently, last);
    }

    public get hasData(): boolean {
        return this.last?.last != null;
    }

    public hasReading(visible: VisibleReadings): boolean {
        return this.getReading(visible) != null;
    }

    public getReading(visible: VisibleReadings): number | null {
        const readings = this.getDecoratedReadings(visible);
        if (readings && readings.length > 0) {
            return readings[0].value;
        }
        return null;
    }

    private createDecoratedReading(sensorId: number, value: number): DecoratedReading[] {
        const sensor = this.meta.findSensorById(sensorId);

        if (sensor && sensor.viz && sensor.viz.length > 0) {
            if (sensor.viz[0].thresholds) {
                const thresholds = sensor.viz[0].thresholds;
                if (thresholds) {
                    const level = thresholds.levels.find((level) => level.start <= value && level.value > value);
                    if (level) {
                        const color = level.color;
                        const label = level.plainLabel?.enUS || level?.mapKeyLabel?.enUS;
                        return [new DecoratedReading(value, color, label)];
                    }
                }
            }
        }

        return [new DecoratedReading(value, "#00ccff", null)];
    }

    public getDecoratedReadings(visible: VisibleReadings): DecoratedReading[] | null {
        switch (visible) {
            case VisibleReadings.Current: {
                const all = _.sortBy(_.flatten(_.values(this.recently)), (row) => row.time);
                if (all.length > 0) {
                    const last = all[all.length - 1];
                    const sensorId = last.sensorId;
                    const value = last.last;
                    if (value === undefined) {
                        return null;
                    }
                    // console.log("station-readings:curr", value);
                    return this.createDecoratedReading(sensorId, value);
                }
                break;
            }
            case VisibleReadings.Last72h: {
                // TODO Order sensors and pick first.
                if (this.recently[72] && this.recently[72].length > 0) {
                    const sensorId = this.recently[72][0].sensorId;
                    const value = this.recently[72][0].max; // TODO Use aggregate function.
                    if (value === undefined) {
                        return null;
                    }
                    // console.log("station-readings:72h", value);
                    return this.createDecoratedReading(sensorId, value);
                }
                break;
            }
        }
        // console.log("station-readings:null", this.recently);
        return null;
    }
}

type StationSortTuiple = [number, number, string];

export class DisplayStation {
    public readonly id: number;
    public readonly name: string;
    public readonly configurations: Configurations;
    public readonly updatedAt: Date;
    public readonly lastReadingAt: Date | null = null;
    public readonly deployedAt: Date | null = null;
    public readonly totalReadings = 0;
    public readonly location: Location | null = null;
    public readonly locationName: string;
    public readonly photos: Photos;
    public readonly modules: DisplayModule[] = [];
    public readonly placeNameOther: string | null;
    public readonly placeNameNative: string | null;
    public readonly battery: number | null;
    public readonly regions: StationRegion[] | null;
    public readonly firmwareNumber: number | null;
    public readonly primarySensor: ModuleSensor | null;
    public readonly attributes: ProjectAttribute[];
    public readonly readOnly: boolean;
    public readonly status: StationStatus;
    public readonly owner: Owner;
    public readonly description: string;

    public get latestPrimary(): number | null {
        if (!this.readings) {
            console.log("stations: No readings");
            return null;
        }

        // TODO Remove after wiring the map values to the new sensorDataQuerier code
        if (this.station.status === StationStatus.down) {
            return null;
        }

        return this.getVisibleReading(VisibleReadings.Current);
    }

    public getVisibleReading(which: VisibleReadings): number | null {
        if (this.readings) {
            return this.readings.getReading(which);
        }
        return null;
    }

    public getDecoratedReadings(visible: VisibleReadings): DecoratedReading[] | null {
        if (this.readings) {
            return this.readings.getDecoratedReadings(visible);
        }
        return null;
    }

    public getSortOrder(which: VisibleReadings): StationSortTuiple {
        if (this.inactive) {
            return [3, 0, this.name];
        }
        if (!this.hasData) {
            return [2, 0, this.name];
        }
        if (this.readings) {
            const value = this.readings.getReading(which);
            if (value !== null) {
                return [0, value, this.name];
            }
        }
        return [1, 0, this.name];
    }

    private get inactive(): boolean {
        return this.status == StationStatus.down;
    }

    public get hasData(): boolean {
        if (this.readings) {
            return this.readings.hasData;
        }
        return false;
    }

    public get hasRecentData(): boolean {
        if (this.readings) {
            return this.readings.hasReading(VisibleReadings.Last72h);
        }
        return false;
    }

    constructor(private readonly station: Station, public readings: StationReadings | null = null) {
        this.id = station.id;
        this.name = station.name;
        this.configurations = station.configurations;
        this.photos = station.photos;
        this.battery = station.battery;
        this.locationName = station.locationName;
        this.placeNameOther = station.placeNameOther;
        this.placeNameNative = station.placeNameNative;
        this.deployedAt = station.recordingStartedAt;
        if (!station.updatedAt) throw new Error(`station missing updatedAt`);
        this.updatedAt = new Date(station.updatedAt);
        this.lastReadingAt = station.lastReadingAt ? new Date(station.lastReadingAt) : null;
        this.attributes = station.attributes.attributes;
        this.readOnly = station.readOnly;
        this.status = station.status;
        this.readings = readings;
        this.owner = station.owner;
        this.description = station.description;

        if (station.configurations.all.length > 0) {
            const ordered = _.orderBy(station.configurations.all[0].modules, ["position"]);
            if (ordered[0]) {
                const orderedSensor = _.orderBy(ordered[0].sensors, ["order"])[0];
                this.primarySensor = orderedSensor;
            }
        }

        this.modules =
            _(station.configurations.all)
                .map((c) => c.modules.filter((m) => !m.internal).map((m) => new DisplayModule(m)))
                .head() || [];

        if (station.location) {
            if (station.location.precise) {
                this.location = new Location(station.location.precise[1], station.location.precise[0]);
            }
            this.regions = station.location.regions;
        }

        this.firmwareNumber = station.firmwareNumber;

    }
}

export function whenWasStationUpdated(station: Station): Date {
    const uploads = station.uploads.map((u) => u.time);
    if (uploads.length > 0) {
        return new Date(uploads[0]);
    }
    return new Date(station.updatedAt);
}

export class ProjectModule {
    constructor(public readonly name: string, public url: string | null) {}
}

export class MapFeature {
    public readonly type = "Feature";
    public readonly geometry: { type: string; coordinates: LngLat | LngLat[][] } | null = null;
    public readonly properties: {
        icon: string;
        id: number;
        value: number | "-" | null;
        thresholds: object | null;
        color: string;
    } | null = null;

    constructor(private readonly station: DisplayStation, type: string, coordinates: any, public readonly bounds: LngLat[]) {
        this.geometry = {
            type: type,
            coordinates: coordinates,
        };
        let thresholds: VizThresholds | null = null;
        if (station.primarySensor && station.primarySensor.meta && station.primarySensor.meta.viz.length > 0) {
            thresholds = station.primarySensor.meta.viz[0].thresholds;
        }

        let color = "#00CCFF";
        // Marker color scale

        if (thresholds) {
            const markerScale = d3
                .scaleThreshold()
                .domain(thresholds.levels.map((d) => d.value))
                .range(thresholds.levels.map((d) => d.color));

            color = markerScale(station.latestPrimary);
        }
        this.properties = {
            id: station.id,
            value: station.status === StationStatus.down ? null : station.latestPrimary,
            icon: "marker",
            thresholds: thresholds,
            color: color,
        };
    }

    public static makeFeatures(station: DisplayStation): MapFeature[] {
        if (station.location) {
            const precise = station.location.lngLat();
            return [new MapFeature(station, "Point", precise, [precise])];
        }
        if (!station.regions || station.regions.length === 0) {
            throw new Error("unmappable");
        }
        const region = station.regions[0];
        return region.shape.map((polygon) => new MapFeature(station, "Polygon", [polygon], polygon));
    }
}

const DefaultLocation: LngLat = [-118.0730372, 34.3318104]; // TwinPeaks
const DefaultMargin = 10000;

export class MappedStations {
    public static make(stations: DisplayStation[]): MappedStations {
        const located = stations.filter((station) => station.location != null);
        const features = _.flatten(located.map((ds) => MapFeature.makeFeatures(ds)));
        const around: BoundingRectangle = features.reduce(
            (bb: BoundingRectangle, feature: MapFeature) => bb.includeAll(feature.bounds),
            new BoundingRectangle()
        );
        return new MappedStations(stations, features, around.zoomOutOrAround(DefaultLocation, DefaultMargin));
    }

    constructor(
        public readonly stations: DisplayStation[] = [],
        public readonly features: MapFeature[] = [],
        public readonly bounds: BoundingRectangle | null = null
    ) {}

    // Test if all displayed map sensors are of the same type
    public get isSingleType(): boolean {
        const moduleNames = this.stations
            .filter((station) => station.configurations.all.length > 0)
            .map((station) => {
                return station.configurations.all[0].modules.map((mod) => {
                    return mod.name;
                });
            });

        return _.uniq(_.flatten(moduleNames)).length === 1;
    }

    public get valid(): boolean {
        return this.bounds != null;
    }

    public boundsLngLat(): LngLat[] | null {
        if (this.bounds) {
            return this.bounds.lngLat();
        }
        return null;
    }

    public overrideBounds(bounds: [LngLat, LngLat]): MappedStations {
        return new MappedStations(this.stations, this.features, new BoundingRectangle(bounds[0], bounds[1]));
    }

    public focusOn(id: number): MappedStations {
        const station = this.stations.find((s) => s.id == id);
        if (!station || !station.location) {
            return this;
        }
        const centered = new BoundingRectangle().include(station.location.lngLat());
        return new MappedStations(this.stations, this.features, centered.zoomOutOrAround(DefaultLocation, DefaultMargin));
    }

    public static defaultBounds(): BoundingRectangle {
        return new BoundingRectangle().zoomOutOrAround(DefaultLocation, DefaultMargin);
    }
}

export class StationsState {
    sensorMeta: SensorMeta | null = null;
    stations: { [index: number]: DisplayStation } = {};
    hasNoStations: boolean;
    projects: { [index: number]: Project } = {};
    projectUsers: { [index: number]: ProjectUser[] } = {};
    projectFollowers: { [index: number]: ProjectFollowers } = {};
    projectStations: { [index: number]: DisplayStation[] } = {};
    projectActivities: { [index: number]: Activity[] } = {};
    user: {
        stations: { [index: number]: DisplayStation };
        projects: { [index: number]: Project };
    } = {
        stations: {},
        projects: {},
    };
    community: {
        projects: { [index: number]: Project };
    } = {
        projects: {},
    };
    mapped: MappedStations | null = null;
    readings: { [index: number]: StationReadings } = {};
}

export class DisplayProject {
    public readonly id: number;
    public readonly name: string;
    public readonly modules: ProjectModule[];
    public readonly duration: number | null = null;
    public readonly mapped: MappedStations;
    public readonly places: { native: string | null } = { native: null };
    public readonly bounds?: BoundingRectangle;

    constructor(public readonly project: Project, public readonly users: ProjectUser[], public readonly stations: DisplayStation[]) {
        this.id = project.id;
        this.name = project.name;
        this.bounds = project.bounds;
        this.modules = _(stations)
            .filter((station) => station.configurations.all.length > 0)
            .map((station) => station.configurations.all[0].modules.filter((m) => !m.internal))
            .flatten()
            .uniqBy((m) => m.name)
            .map((m) => new ProjectModule(m.name, null))
            .value();
        this.mapped = MappedStations.make(stations);
        if (project.startTime && project.endTime) {
            const start = new Date(project.startTime);
            const end = new Date(project.endTime);
            end.setDate(end.getDate() + 1); // FK-1927
            this.duration = end.getTime() - start.getTime();
        }
        this.places.native = _(stations)
            .map((s) => s.placeNameNative)
            .filter((n) => n != null && n.length > 0)
            .uniq()
            .join(", ");
    }
}

const getters = {
    projectsById(state: StationsState): { [index: number]: DisplayProject } {
        return _(Object.values(state.projects))
            .map((p) => {
                const users = state.projectUsers[p.id];
                const stations = state.projectStations[p.id];
                if (stations && users) {
                    return [p.id, new DisplayProject(p, users, stations)];
                }
                return [p.id, null];
            })
            .fromPairs()
            .value();
    },
    stationsById(state: StationsState): { [index: number]: DisplayStation } {
        return { ...state.stations, ...state.user.stations };
    },
    mapped(state: StationsState): MappedStations | null {
        return state.mapped;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ dispatch }: { dispatch: any }) => {
            await dispatch(ActionTypes.NEED_COMMON);
        },
        [ActionTypes.AUTHENTICATED]: async ({ dispatch }: { dispatch: any }) => {
            await dispatch(ActionTypes.NEED_COMMON);
        },
        [ActionTypes.NEED_COMMON]: async ({ dispatch, commit }: { dispatch: any; commit: any }) => {
            await dispatch(NEED_SENSOR_META);
            await Promise.all([dispatch(ActionTypes.NEED_PROJECTS), dispatch(ActionTypes.NEED_STATIONS)]);
        },
        [NEED_SENSOR_META]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
            const meta = await services.api.getAllSensorsMemoized()(); // TODO  Why?
            const sensorMeta = new SensorMeta(meta);
            commit(SENSOR_META, sensorMeta);
        },
        [ActionTypes.NEED_PROJECTS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
            commit(MutationTypes.LOADING, { projects: true });
            const [communityProjects, userProjects] = await Promise.all([
                services.api.getPublicProjects(),
                services.api.getUserProjects(() => {
                    return Promise.resolve({
                        projects: [],
                    });
                }),
            ]);
            commit(HAVE_COMMUNITY_PROJECTS, communityProjects.projects);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
            commit(MutationTypes.LOADING, { projects: false });
        },
        [ActionTypes.NEED_STATIONS]: async ({ commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState }) => {
            commit(MutationTypes.LOADING, { stations: true });
            const [stations] = await Promise.all([
                services.api.getUserStations(() => {
                    return Promise.resolve({
                        stations: [],
                    });
                }),
            ]);

            const recently = await services.api.queryStationsRecently(stations.stations.map((s: { id: number }) => s.id));

            commit(HAVE_USER_STATIONS, { stations: stations.stations, recently: recently });
            commit(MutationTypes.LOADING, { stations: false });
        },
        [ActionTypes.NEED_PROJECT]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            payload: { id: number }
        ) => {
            commit(MutationTypes.LOADING, { projects: true });

            const [project, users, stations] = await Promise.all([
                services.api.getProject(payload.id),
                services.api.getUsersByProject(payload.id),
                services.api.getStationsByProject(payload.id),
            ]);

            const recently = await services.api.queryStationsRecently(stations.stations.map((s: { id: number }) => s.id));

            commit(PROJECT_LOADED, project);
            commit(PROJECT_USERS, { projectId: payload.id, users: users.users });
            commit(PROJECT_STATIONS, { projectId: payload.id, stations: stations.stations, recently: recently });

            commit(MutationTypes.LOADING, { projects: false });
        },
        [ActionTypes.NEED_STATION]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            payload: { id: number }
        ) => {
            commit(MutationTypes.LOADING, { stations: true });

            const [station, recently] = await Promise.all([
                services.api.getStation(payload.id),
                services.api.queryStationsRecently([payload.id]),
            ]);

            commit(STATION_UPDATE, { station, recently });

            commit(MutationTypes.LOADING, { stations: false });
        },
        [ActionTypes.UPDATE_STATION]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            payload: {id: number, name: string, description: string | null }
        ) => {
            commit(MutationTypes.LOADING, { stations: true });

            services.api.updateStation(payload).then((station) => {
                commit(STATION_UPDATE, { station });
            });

            commit(MutationTypes.LOADING, { stations: false });
        },
        [ActionTypes.CLEAR_STATION]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: StationsState },
            id: number
        ) => {
            commit(STATION_CLEAR, id);
        },
        [ActionTypes.PROJECT_FOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.followProject(payload.projectId);
            commit(PROJECT_LOADED, await services.api.getProject(payload.projectId));
        },
        [ActionTypes.PROJECT_UNFOLLOW]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.unfollowProject(payload.projectId);
            commit(PROJECT_LOADED, await services.api.getProject(payload.projectId));
        },
        [ActionTypes.STATION_PROJECT_ADD]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { stationId: number; projectId: number }
        ) => {
            await services.api.addStationToProject(payload);

            const stations = await services.api.getStationsByProject(payload.projectId);
            commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
        },
        [ActionTypes.STATION_PROJECT_REMOVE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { stationId: number; projectId: number }
        ) => {
            await services.api.removeStationFromProject(payload);

            const stations = await services.api.getStationsByProject(payload.projectId);
            commit(PROJECT_STATIONS, { projectId: payload.projectId, stations: stations.stations });
        },
        [ActionTypes.PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { projectId: number; email: string; role: string }
        ) => {
            await services.api.sendInvite(payload);
            const usersReply = await services.api.getUsersByProject(payload.projectId);
            commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
        },
        [ActionTypes.PROJECT_REMOVE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { projectId: number; email: string }
        ) => {
            await services.api.removeUserFromProject(payload);
            const usersReply = await services.api.getUsersByProject(payload.projectId);
            commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
        },
        [ActionTypes.PROJECT_EDIT_ROLE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { projectId: number; email: string; role: number }
        ) => {
            await services.api.editRole(payload);
            const usersReply = await services.api.getUsersByProject(payload.projectId);
            commit(PROJECT_USERS, { projectId: payload.projectId, users: usersReply.users });
        },
        [ActionTypes.ACCEPT_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.acceptProjectInvite(payload);

            const userProjects = await services.api.getUserProjects(OnNoReject);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
        },
        [ActionTypes.DECLINE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.declineProjectInvite(payload);
        },
        [ActionTypes.ACCEPT_PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { id: number; token: string }
        ) => {
            await services.api.acceptInvite(payload);

            const userProjects = await services.api.getUserProjects(OnNoReject);
            commit(HAVE_USER_PROJECTS, userProjects.projects);
        },
        [ActionTypes.DECLINE_PROJECT_INVITE]: async (
            { commit, dispatch }: { commit: any; dispatch: any },
            payload: { id: number; token: string }
        ) => {
            await services.api.declineInvite(payload);
        },
        [ActionTypes.DELETE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: { projectId: number }) => {
            await services.api.deleteProject(payload);

            commit(PROJECT_DELETED, payload);

            return payload;
        },
        [ActionTypes.ADD_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: any) => {
            const project = await services.api.addProject(payload);

            commit(PROJECT_UPDATE, project);

            return project;
        },
        [ActionTypes.SAVE_PROJECT]: async ({ commit, dispatch }: { commit: any; dispatch: any }, payload: any) => {
            const project = await services.api.updateProject(payload);

            commit(PROJECT_UPDATE, project);

            return project;
        },
    };
};

function makeDisplayStations(sensorMeta: SensorMeta | null, stations: Station[], recently: QueryRecentlyResponse | null): DisplayStation[] {
    if (sensorMeta === null) throw new Error("fatal: Sensor meta load order error");
    return stations.map((station) => {
        if (recently) {
            const windows = _.mapValues(recently.windows, (rows, hours) => {
                return rows.filter((row) => row.stationId == station.id);
            });
            const readings = new StationReadings(station.id, sensorMeta, windows, recently.stations[station.id]);
            return new DisplayStation(station, readings);
        }
        return new DisplayStation(station);
    });
}

const mutations = {
    [HAVE_COMMUNITY_PROJECTS]: (state: StationsState, projects: Project[]) => {
        Vue.set(
            state.community,
            "projects",
            _.keyBy(projects, (p) => p.id)
        );
        Vue.set(state, "projects", { ...state.projects, ..._.keyBy(projects, (p) => p.id) });
    },
    [HAVE_USER_PROJECTS]: (state: StationsState, projects: Project[]) => {
        Vue.set(
            state.user,
            "projects",
            _.keyBy(projects, (p) => p.id)
        );
        Vue.set(state, "projects", { ...state.projects, ..._.keyBy(projects, (p) => p.id) });
    },
    [HAVE_USER_STATIONS]: (state: StationsState, payload: { stations: Station[]; recently: QueryRecentlyResponse | null }) => {
        const userStations = makeDisplayStations(state.sensorMeta, payload.stations, payload.recently);

        Vue.set(
            state.user,
            "stations",
            _.keyBy(userStations, (s) => s.id)
        );
        Vue.set(state, "stations", { ...state.stations, ..._.keyBy(userStations, (s) => s.id) });
        Vue.set(state, "hasNoStations", userStations.length == 0);
        Vue.set(state, "mapped", MappedStations.make(userStations));
    },
    [PROJECT_USERS]: (state: StationsState, payload: { projectId: number; users: ProjectUser[] }) => {
        Vue.set(state.projectUsers, payload.projectId, payload.users);
    },
    [PROJECT_FOLLOWS]: (state: StationsState, payload: { projectId: number; followers: ProjectFollowers }) => {
        Vue.set(state.projectFollowers, payload.projectId, payload.followers);
    },
    [PROJECT_STATIONS]: (
        state: StationsState,
        payload: { projectId: number; stations: Station[]; recently: QueryRecentlyResponse | null }
    ) => {
        const projectStations = makeDisplayStations(state.sensorMeta, payload.stations, payload.recently);

        state.stations = { ...state.stations, ..._.keyBy(projectStations, (s) => s.id) };
        Vue.set(state.projectStations, payload.projectId, projectStations);
    },
    [PROJECT_ACTIVITY]: (state: StationsState, payload: { projectId: number; activities: Activity[] }) => {
        Vue.set(state.projectActivities, payload.projectId, payload.activities);
    },
    [STATION_UPDATE]: (state: StationsState, payload: { station: Station; recently: QueryRecentlyResponse }) => {
        const updated = makeDisplayStations(state.sensorMeta, [payload.station], payload.recently);
        Vue.set(state.stations, payload.station.id, updated[0]);
    },
    [STATION_CLEAR]: (state: StationsState, id: number) => {
        Vue.set(state.stations, id, null);
    },
    [PROJECT_LOADED]: (state: StationsState, project: Project) => {
        Vue.set(state.projects, project.id, project);
        if (state.user.projects[project.id]) {
            Vue.set(state.user.projects, project.id, project);
        }
        if (state.community.projects[project.id]) {
            Vue.set(state.community.projects, project.id, project);
        }
    },
    [PROJECT_UPDATE]: (state: StationsState, project: Project) => {
        Vue.set(state.projects, project.id, project);
        Vue.set(state.user.projects, project.id, project);
        if (state.community.projects[project.id]) {
            Vue.set(state.community.projects, project.id, project);
        }
    },
    [PROJECT_DELETED]: (state: StationsState, payload: { projectId: number }) => {
        delete state.projects[payload.projectId];
        delete state.user.projects[payload.projectId];
        delete state.community.projects[payload.projectId];
    },
    [SENSOR_META]: (state: StationsState, payload: SensorMeta) => {
        state.sensorMeta = payload;
    },
};

export const stations = (services: Services) => {
    const state = () => new StationsState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
