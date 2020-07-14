<template>
    <div class="readings-simple">
        <template v-if="!loading">
            <div v-for="sensor in sensors" v-bind:key="sensor.key" class="reading-container">
                <div class="reading">
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
import Vue from "vue";
import FKApi from "@/api/api";

export enum TrendType {
    Downward,
    Steady,
    Upward,
}

export class SensorReading {
    constructor(
        public readonly labelKey: string,
        public readonly unitOfMeasure: string,
        public readonly reading: number,
        public readonly trend: TrendType
    ) {}
}

export default Vue.extend({
    name: "LatestStationReadings",
    components: {},
    props: {
        id: {
            type: Number,
            required: true,
        },
    },
    data: () => {
        return {
            loading: true,
            sensors: [],
        };
    },
    beforeMount(this: any) {
        return this.refresh();
    },
    watch: {
        id() {
            this.refresh();
        },
    },
    methods: {
        refresh() {
            const data = () => {
                const params = new URLSearchParams();
                params.append("stations", [this.id].join(","));
                params.append("tail", "1");
                return new FKApi().tailSensorData(params);
            };

            const meta = () => new FKApi().getAllSensors();

            this.loading = true;

            return Promise.all([data(), meta()])
                .then(([data, meta]) => {
                    const idsToKey = _.mapValues(
                        _.keyBy(meta.sensors, (k) => k.id),
                        (v) => v.key
                    );
                    const idsToValue = _.mapValues(
                        _.keyBy(data.data, (r) => r.sensorId),
                        (r) => r.value
                    );
                    const keysToValue = _(idsToValue)
                        .map((value, id) => {
                            return [idsToKey[id], value];
                        })
                        .fromPairs()
                        .value();

                    const sensorsByKey = _(meta.modules)
                        .map((m) => m.sensors)
                        .flatten()
                        .keyBy((s) => s.fullKey)
                        .value();

                    return _(keysToValue)
                        .map((value, key) => {
                            const sensor = sensorsByKey[key];
                            if (!sensor) {
                                throw new Error("no sensor meta");
                            }
                            return new SensorReading(key, sensor.unitOfMeasure, value, TrendType.Steady);
                        })
                        .value();
                })
                .then((sensors) => {
                    console.log("sensors", sensors);
                    this.sensors = sensors;
                    this.loading = false;
                    return this.sensors;
                });
        },
    },
});
</script>

<style scoped>
.readings-simple {
    display: flex;
    flex-wrap: wrap;
}
.readings-simple .reading-container {
    background-color: white;
    flex: 50%;
    margin-bottom: 10px;
}
.reading {
    background-color: #efefef;
    display: flex;
    flex-direction: row;
    padding: 5px;
    margin: 5px 5px;
}
.reading .name {
    font-size: 12px;
    line-height: 20px;
}
.reading .value {
    margin-left: auto;
    font-size: 16px;
}
.reading .uom {
    font-size: 10px;
}
</style>
