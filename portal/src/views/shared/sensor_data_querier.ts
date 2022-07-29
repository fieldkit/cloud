import {
    FKApi,
    VizSensor,
    StationInfoResponse,
    ModuleSensorMeta,
    SensorInfoResponse,
    SensorsResponse,
    TailSensorDataResponse,
    QueryRecentlyResponse,
} from "@/api";

import { promiseAfter } from "@/utilities";

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

export class SensorDataQuerier {
    private data: Promise<TailSensorDataResponse> | null = null;
    private quickSensors: Promise<SensorInfoResponse> | null = null;

    constructor(private readonly api: FKApi, private readonly stationIds: number[]) {
        console.log("sdq:ctor", stationIds);
    }

    private getBackend(): string | null {
        return window.localStorage["fk:backend"] || "tsdb";
    }

    private _queue: Promise<
        [Promise<TailSensorDataResponse>, Promise<SensorInfoResponse>, Promise<SensorMeta>, Promise<QueryRecentlyResponse>]
    > | null = null;
    private _queued: number[] = [];

    public async query(stationId: number): Promise<[TailSensorDataResponse, StationQuickSensors, SensorMeta, QueryRecentlyResponse]> {
        // TODO Check if we already have the data.

        this._queued.push(stationId);

        if (this._queue == null) {
            this._queue = promiseAfter(50).then(() => {
                const ids = this._queued;
                this._queued = [];
                this._queue = null;

                console.log("sdq:querying", ids);

                const data = this.api.tailSensorData(ids);
                const quickSensors = this.api.getQuickSensors(ids);

                const sensorMeta = this.api
                    .getAllSensorsMemoized()()
                    .then((meta) => new SensorMeta(meta));

                const recently = this.api.queryStationsRecently(ids);

                return [data, quickSensors, sensorMeta, recently];
            });
        }

        return this._queue
            .then(([data, quickSensors, sensorMeta, recently]) => {
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

                const stationRecently = recently.then((response) => {
                    return _.mapValues(response, (window, hours) => {
                        return window.filter((row) => row.stationId == stationId);
                    });
                });

                return Promise.all([dataQuery, quickSensorsQuery, sensorMeta, stationRecently]);
            })
            .then(([data, quickSensors, meta, recently]) => {
                return [data, quickSensors, meta, recently];
            });
    }
}
