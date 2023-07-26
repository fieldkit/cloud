<template>
    <div class="tiny-chart">
        <LineChart :series="series" :settings="chartSettings" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import { ModuleSensor, DisplayStation, TailSensorDataRow } from "@/store";
import { SensorDataQuerier, StandardObserver } from "../shared/sensor_data_querier";
import { DataRow } from "@/views/viz/api";

import {
    ChartSettings,
    SeriesData,
    DataSetSeries,
    QueriedData,
    VizInfo,
    VizSensor,
    TimeRange,
    SensorDataResponse,
    ModuleSensorMeta,
} from "./vega/SpecFactory";

import LineChart from "./vega/LineChart.vue";

export default Vue.extend({
    name: "TinyChart",
    components: {
        LineChart,
    },
    props: {
        stationId: {
            type: Number,
            required: true,
        },
        station: {
            type: DisplayStation,
            required: true,
        },
        querier: {
            type: SensorDataQuerier,
            required: true,
        },
    },
    data(): {
        chartSettings: ChartSettings;
        series: SeriesData[];
    } {
        return {
            chartSettings: ChartSettings.Tiny,
            series: [],
        };
    },
    watch: {
        stationId(): void {
            console.log("tiny-chart:", "station-id");
            // TODO Should we also be checking if we're in view?
            void this.load();
        },
    },
    async mounted() {
        new StandardObserver().observe(this.$el, () => {
            console.log("tiny-chart:observed");
            void this.load();
        });
    },
    methods: {
        async load(): Promise<void> {
            // TODO Show busy when loading?
            const [stationData, quickSensors, meta] = await this.querier.queryTinyChartData(this.stationId);

            function getFirstQuickSensor(quickSensors): VizSensor {
                const sensorModuleId = quickSensors.station[0].moduleId;
                const sensorId = quickSensors.station[0].sensorId;
                return [quickSensors.station[0].stationId, [sensorModuleId, sensorId]];
            }

            // Not happy with this, no easy way until we get rid of the 'quick sensors'
            // call, which I'd love to do. Quick sensors won't return a sensor if the sensor
            // has no data and so to properly show a blank chart in that case we gotta
            // resolve the sensor to show. This uses the DisplayStation configuration data
            // to fill that in, by choosing the first one and then making up a fake vizSensor.
            function getSensorMetaAndData(
                station: DisplayStation
            ): { vizSensor: VizSensor; sensor: ModuleSensorMeta; sdr: SensorDataResponse } | null {
                if (quickSensors.station.filter((r) => r.moduleId != null).length == 0) {
                    const moduleSensor: ModuleSensor[] = _.flattenDeep(
                        station.configurations.all.map((c) => c.modules.map((m) => m.sensors))
                    );
                    if (moduleSensor.length == 0) {
                        return null;
                    }

                    const sensor = meta.findSensorByKey(moduleSensor[0].fullKey);
                    return { vizSensor: [station.id, ["", 0]], sensor, sdr: { data: [], bucketSamples: 0, bucketSize: 0, dataEnd: null } };
                }

                const vizSensor = getFirstQuickSensor(quickSensors);
                const sensor = meta.findSensor(vizSensor);
                const data = stationData.data.filter((datum) => datum.sensorId == vizSensor[1][1]); // TODO VizSensor

                const sdr: SensorDataResponse = {
                    data: data.map(
                        (row: TailSensorDataRow): DataRow => {
                            return {
                                time: row.time,
                                stationId: row.stationId,
                                moduleId: row.moduleId,
                                sensorId: row.sensorId,
                                value: row.max !== undefined ? row.max : null, // TODO
                                avg: row.avg,
                                min: row.min,
                                max: row.max,
                                last: row.last,
                                location: null,
                            };
                        }
                    ),
                    bucketSize: stationData.stations[station.id].bucketSize,
                    bucketSamples: 0, // Maybe wrong
                    dataEnd: null, // Maybe wrong
                };

                return { vizSensor, sensor, sdr };
            }

            const maybeSensorMetaAndData = getSensorMetaAndData(this.station);
            if (!maybeSensorMetaAndData) {
                console.log("tiny-chart:empty", { quickSensors, station: this.station });
                return;
            }

            const { vizSensor, sensor, sdr } = maybeSensorMetaAndData;

            const queried = new QueriedData("key", TimeRange.eternity, sdr);
            const strings = sensor.strings["enUs"];
            const name = strings["label"] || "Unknown";
            const axisLabel = strings["axisLabel"] || strings["label"];
            const station = { name: this.station.name, location: null };
            const vizInfo = new VizInfo(
                sensor.key,
                [], // Scale
                station, // Station
                sensor.unitOfMeasure,
                sensor.key,
                name,
                sensor.viz || [],
                sensor.ranges,
                "", // ChartLabel
                axisLabel
            );

            console.log("tiny-chart", { stationId: this.stationId, station: this.station, sensor, quickSensors, meta });

            this.series = [
                new SeriesData(
                    "key",
                    new TimeRange(queried.timeRange[0], queried.timeRange[1]),
                    new DataSetSeries(vizSensor, queried, queried),
                    queried,
                    vizInfo
                ),
            ];
        },
    },
});
</script>

<style scoped lang="scss">
.tiny-chart {
    margin-top: 10px;
    width: 100%;
    height: 150px;
    border: solid 1px #d8dce0;
}

.tiny-chart .fit-x {
    width: 100%;
}

.tiny-chart .fit-y {
    height: 100%;
}
</style>
