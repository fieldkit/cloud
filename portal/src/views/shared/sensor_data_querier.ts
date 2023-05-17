import type { QueryRecentlyResponse, RecentlyAggregatedWindows, ModuleSensorMeta } from "@/api"
import {
    FKApi,
    VizSensor,
    StationInfoResponse,
    SensorInfoResponse,
    SensorsResponse,
    TailSensorDataResponse,
} from "@/api";

import { SensorMeta } from "@/store";
import { promiseAfter } from "@/utilities";
import _ from "lodash";

export { ModuleSensorMeta, QueryRecentlyResponse, RecentlyAggregatedWindows };

export interface StationQuickSensors {
    station: StationInfoResponse[];
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
                        stations: { [stationId]: response.stations[stationId] },
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
