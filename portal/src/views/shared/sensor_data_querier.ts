import {
    FKApi,
    VizSensor,
    StationInfoResponse,
    ModuleSensorMeta,
    SensorInfoResponse,
    SensorsResponse,
    TailSensorDataResponse,
    QueryRecentlyResponse,
    RecentlyAggregatedWindows,
} from "@/api";

import { promiseAfter } from "@/utilities";
import _ from "lodash";

export { ModuleSensorMeta, QueryRecentlyResponse, RecentlyAggregatedWindows };

export interface StationQuickSensors {
    station: StationInfoResponse[];
}

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

type HandlerType<PromisedData> = (ids: number[]) => PromisedData;

class Batcher<T> {
    private queued: number[] = [];
    private queue: Promise<T> | null = null;

    constructor(private readonly handler: HandlerType<T>) {}

    public async query(id: number): Promise<T> {
        this.queued.push(id);

        if (this.queue == null) {
            this.queue = promiseAfter(50).then(() => {
                const ids = this.queued;
                this.queued = [];
                this.queue = null;
                return this.handler(ids);
            });
        }

        return this.queue;
    }
}

export class SensorDataQuerier {
    constructor(private readonly api: FKApi) {}

    private tinyChartData = new Batcher<[Promise<TailSensorDataResponse>, Promise<SensorInfoResponse>, Promise<SensorMeta>]>(
        (ids: number[]) => {
            console.log("tcd:querying", ids);
            const data = this.api.tailSensorData(ids);
            const quickSensors = this.api.getQuickSensors(ids);
            const sensorMeta = this.querySensorMeta();
            return [data, quickSensors, sensorMeta];
        }
    );

    private recently = new Batcher<Promise<QueryRecentlyResponse>>((ids: number[]) => {
        console.log("qrd:querying", ids);
        return this.api.queryStationsRecently(ids);
    });

    public async querySensorMeta(): Promise<SensorMeta> {
        return this.api
            .getAllSensorsMemoized()()
            .then((meta) => new SensorMeta(meta));
    }

    public async queryTinyChartData(stationId: number): Promise<[TailSensorDataResponse, StationQuickSensors, SensorMeta]> {
        return this.tinyChartData
            .query(stationId)
            .then(([data, quickSensors, sensorMeta]) => {
                const dataQuery = data.then((response) => {
                    return {
                        data: response.data.filter((row) => row.stationId == stationId),
                    };
                });

                const quickSensorsQuery = quickSensors.then((response) => {
                    return {
                        station: response.stations[stationId],
                    };
                });

                return Promise.all([dataQuery, quickSensorsQuery, sensorMeta]);
            })
            .then(([data, quickSensors, meta]) => {
                return [data, quickSensors, meta];
            });
    }

    public async queryRecently(stationId: number): Promise<QueryRecentlyResponse> {
        return this.recently
            .query(stationId)
            .then((response) => response)
            .then((response) => {
                const filteredWindows = _.mapValues(response.windows, (rows, hours) => {
                    return rows.filter((row) => row.stationId == stationId);
                });

                return {
                    windows: filteredWindows,
                    stations: response.stations,
                };
            });
    }
}

export class StandardObserver {
    observe(el, handler) {
        if (!("IntersectionObserver" in window)) {
            console.log("tiny-chart:warning", "no-intersection-observer");
        } else {
            const observer = new IntersectionObserver((entries) => {
                // Use `intersectionRatio` because of Edge 15's lack of support for
                // `isIntersecting`.  See: // https://github.com/w3c/IntersectionObserver/issues/211
                if (entries[0].intersectionRatio <= 0) return;

                // Cleanup
                observer.unobserve(el);

                handler();
            });

            // We observe the root `$el` of the mounted loading component to detect
            // when it becomes visible.
            observer.observe(el);
        }
    }
}
