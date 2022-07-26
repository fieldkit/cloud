<template>
    <div class="readings-simple">
        <template v-if="!loading">
            <div v-for="sensor in sensors" v-bind:key="sensor.labelKey" class="reading-container">
                <div :class="'reading ' + sensor.classes">
                    <div class="name">{{ $t(sensor.labelKey) }}</div>
                    <div class="value">{{ sensor | prettyReading }}</div>
                    <div class="uom">{{ sensor.unitOfMeasure }}</div>
                </div>
            </div>
            <div v-if="sensors.length == 0">No readings yet.</div>
        </template>
        <div class="loading" v-if="loading">Loading</div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Config from "@/secrets";
import { promiseAfter } from "@/utilities";

import Vue, { PropType } from "vue";

import {
    FKApi,
    TailSensorDataResponse,
    VizSensor,
    StationInfoResponse,
    ModuleSensorMeta,
    SensorInfoResponse,
    SensorsResponse,
    Module,
} from "@/api";

export enum TrendType {
    Downward,
    Steady,
    Upward,
}

export class SensorReading {
    constructor(
        public readonly labelKey: string,
        public readonly classes: string,
        public readonly order: number,
        public readonly unitOfMeasure: string,
        public readonly reading: number,
        public readonly trend: TrendType,
        public readonly internal: boolean,
        public readonly sensorModule: Module
    ) {}
}

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
        if (Config.backend) {
            return Config.backend;
        }
        return window.localStorage["fk:backend"] || null;
    }

    private _queue: Promise<[Promise<TailSensorDataResponse>, Promise<SensorInfoResponse>, Promise<SensorMeta>]> | null = null;
    private _queued: number[] = [];

    public async query(stationId: number): Promise<[TailSensorDataResponse, StationQuickSensors, SensorMeta]> {
        // TODO Check if we already have the data.

        this._queued.push(stationId);

        if (this._queue == null) {
            this._queue = promiseAfter(50).then(() => {
                const ids = this._queued;
                this._queued = [];
                this._queue = null;

                console.log("sdq:querying", ids);

                const params = new URLSearchParams();
                params.append("stations", this.stationIds.join(","));
                params.append("tail", "1");
                const backend = this.getBackend();
                if (backend) {
                    params.append("backend", backend);
                }
                const data = this.api.tailSensorData(params);
                const quickSensors = this.api.getQuickSensors(ids);

                const sensorMeta = this.api
                    .getAllSensorsMemoized()()
                    .then((meta) => new SensorMeta(meta));

                return [data, quickSensors, sensorMeta];
            });
        }

        return this._queue
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
}

export default Vue.extend({
    name: "LatestStationReadings",
    components: {},
    props: {
        id: {
            type: Number,
            required: true,
        },
        moduleKey: {
            type: String,
            required: false,
        },
        sensorDataQuerier: {
            type: Object as PropType<SensorDataQuerier>,
            required: false,
        },
    },
    data(): {
        loading: boolean;
        sensors: SensorReading[];
        querier: SensorDataQuerier;
    } {
        return {
            loading: true,
            sensors: [],
            querier: this.sensorDataQuerier || new SensorDataQuerier(this.$services.api, [this.id]),
        };
    },
    watch: {
        async id(): Promise<void> {
            await this.refresh();
        },
        async moduleKey(): Promise<void> {
            await this.refresh();
        },
    },
    async beforeMount(): Promise<void> {
        await this.refresh();
    },
    methods: {
        refresh() {
            this.loading = true;

            return this.querier
                .query(this.id)
                .then(([data, quickSensors, meta]) => {
                    const sensorsToModule = _.fromPairs(
                        _.flatten(
                            meta.modules.map((module) => {
                                return module.sensors.map((sensor) => [sensor.fullKey, module]);
                            })
                        )
                    );

                    const sensorsByKey = _(meta.modules)
                        .map((m) => m.sensors)
                        .flatten()
                        .keyBy((s) => s.fullKey)
                        .value();

                    const idsToKey = _.mapValues(
                        _.keyBy(meta.sensors, (k) => k.id),
                        (v) => v.key
                    );

                    const idsToValue = _.mapValues(
                        _.keyBy(data.data, (r) => r.sensorId),
                        (r) => r
                    );

                    const keysToValue = _(idsToValue)
                        .map((value, id) => [idsToKey[id], value])
                        .fromPairs()
                        .value();

                    const readings = _(keysToValue)
                        .map((reading, key) => {
                            const sensor = sensorsByKey[key];
                            if (!sensor) throw new Error("no sensor meta");

                            const value = reading.value || reading[sensor.aggregationFunction];

                            const classes = [key.replaceAll(".", "-")];
                            if (sensor.unitOfMeasure == "°") {
                                classes.push("degrees");
                            }
                            const sensorModule = sensorsToModule[key];
                            if (!sensorModule) throw new Error("no sensor module");

                            return new SensorReading(
                                key,
                                classes.join(" "),
                                sensor.order,
                                sensor.unitOfMeasure,
                                value,
                                TrendType.Steady,
                                sensor.internal,
                                sensorModule
                            );
                        })
                        .value();

                    return _.orderBy(
                        readings.filter((sr) => !sr.internal),
                        (sr) => sr.order
                    );
                })
                .then((sensors) => {
                    if (this.moduleKey) {
                        this.sensors = sensors.filter((sensor) => sensor.sensorModule.key === this.moduleKey);
                    } else {
                        this.sensors = sensors;
                    }
                    this.loading = false;
                    this.$emit("layoutChange");

                    return this.sensors;
                });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.readings-simple {
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;

    @at-root .station-readings & {
        justify-content: flex-start;
    }
}
.readings-simple .reading-container {
    background-color: white;
    flex: 0 0 calc(50% - 5px);
    margin-bottom: 10px;

    @include bp-down($xs) {
        flex-basis: 100%;
    }
}
.reading {
    display: flex;
    align-items: center;
    padding: 0 9px;
    height: 40px;
    border-radius: 2px;
    font-size: 12px;
    font-weight: 500;
    background-color: #f4f5f7;
}
.reading .name {
    font-size: 12px;
    margin-right: 5px;
}
.reading .value {
    margin-left: auto;
    margin-right: 2px;
    font-size: 16px;
    white-space: nowrap;
}
.reading .uom {
    font-size: 10px;
    margin-bottom: -3px;
}
.reading.degrees .uom {
    align-self: start;
    padding-top: 5px;
}
</style>
