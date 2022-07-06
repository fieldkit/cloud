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
import Vue, { PropType } from "vue";
import { FKApi, TailSensorDataResponse, SensorsResponse, Module } from "@/api";

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

export class SensorDataQuerier {
    private everything: Promise<TailSensorDataResponse> | null = null;

    constructor(private readonly api: FKApi, private readonly stationIds: number[]) {
        console.log("sdq:ctor", stationIds);
    }

    public async query(stationId: number): Promise<[TailSensorDataResponse, SensorsResponse]> {
        if (this.everything == null) {
            const params = new URLSearchParams();
            params.append("stations", this.stationIds.join(","));
            params.append("tail", "1");
            this.everything = this.api.tailSensorData(params);
        }

        const dataQuery = this.everything.then((response) => {
            return {
                data: response.data.filter((row) => row.stationId == stationId),
            };
        });

        return await Promise.all([dataQuery, this.api.getAllSensorsMemoized()()]);
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
                .then(([data, meta]) => {
                    const sensorsToModule = _.fromPairs(
                        _.flatten(
                            meta.modules.map((module) => {
                                return module.sensors.map((sensor) => [sensor.fullKey, module]);
                            })
                        )
                    );

                    const idsToKey = _.mapValues(
                        _.keyBy(meta.sensors, (k) => k.id),
                        (v) => v.key
                    );

                    const idsToValue = _.mapValues(
                        _.keyBy(data.data, (r) => r.sensorId),
                        (r) => r.value
                    );

                    const keysToValue = _(idsToValue)
                        .map((value, id) => [idsToKey[id], value])
                        .fromPairs()
                        .value();

                    const sensorsByKey = _(meta.modules)
                        .map((m) => m.sensors)
                        .flatten()
                        .keyBy((s) => s.fullKey)
                        .value();

                    const readings = _(keysToValue)
                        .map((value, key) => {
                            const sensor = sensorsByKey[key];
                            if (!sensor) throw new Error("no sensor meta");
                            const classes = [key.replaceAll(".", "-")];
                            if (sensor.unitOfMeasure == "°") {
                                classes.push("degrees");
                            }
                            const sensorModule = sensorsToModule[key];
                            if (!sensorModule) throw new Error("no sensor module");
                            // console.log(`sensor:`, sensor);

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
                    // console.log(`sensors`, sensors);
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
