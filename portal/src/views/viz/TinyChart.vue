<template>
    <div class="tiny-chart">
        <LineChart :series="series" :settings="chartSettings" />
        <div v-if="noSensorData" class="no-data">{{ $tc("noData") }}</div>
        <div v-if="busy" class="no-data">
            <Spinner class="spinner" />
        </div>
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
import Spinner from "@/views/shared/Spinner.vue";

export default Vue.extend({
    name: "TinyChart",
    components: {
        Spinner,
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
        moduleKey: {
            type: String,
            required: false,
        },
    },
    data(): {
        chartSettings: ChartSettings;
        series: SeriesData[];
        // can be used with BookmarkFactory.forSensor to create a bookmark for the currently displayed data
        vizData: { vizSensor: VizSensor; timeRange: [number, number] } | null;
        noSensorData: boolean;
        busy: boolean;
    } {
        return {
            chartSettings: ChartSettings.Tiny,
            series: [],
            vizData: null,
            noSensorData: false,
            busy: true,
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
            this.series = [];
            void this.load();
        });
    },
    methods: {
        async load(): Promise<void> {
            this.busy = true;
            const [stationData, quickSensors, meta] = await this.querier.queryTinyChartData(this.stationId);
            const selectedModuleKey = this.moduleKey;
            // eslint-disable-next-line @typescript-eslint/no-this-alias
            const thisComp = this;
            thisComp.noSensorData = false;

            function getQuickSensor(quickSensors): VizSensor {
                let stationId, sensorModuleId, sensorId;
                if (selectedModuleKey) {
                    const module = quickSensors.station.find((station) => station.moduleKey === selectedModuleKey);
                    if (module) {
                        stationId = module.stationId;
                        sensorModuleId = module.moduleId;
                        sensorId = module.sensorId;
                    } else {
                        thisComp.handleNoDataCase();
                    }
                } else {
                    stationId = quickSensors.station[0].stationId;
                    sensorModuleId = quickSensors.station[0].moduleId;
                    sensorId = quickSensors.station[0].sensorId;
                }
                return [stationId, [sensorModuleId, sensorId]];
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
                    return {
                        vizSensor: [station.id, ["", 0]],
                        sensor,
                        sdr: { data: [], bucketSamples: 0, bucketSize: 0, dataEnd: null },
                    };
                }
                const vizSensor = getQuickSensor(quickSensors);
                const sensor = meta.findSensor(vizSensor);
                const data = stationData.data.filter((datum) => datum.sensorId == vizSensor[1][1]); // TODO VizSensor

                if (data.length === 0) {
                    thisComp.handleNoDataCase();
                    return null;
                }

                console.log("station.id", station.id);
                const sdr: SensorDataResponse = {
                    data: data.map((row: TailSensorDataRow): DataRow => {
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
                    }),
                    bucketSize: stationData.stations[station.id].bucketSize,
                    bucketSamples: 0, // Maybe wrong
                    dataEnd: null, // Maybe wrong
                };

                return { vizSensor, sensor, sdr };
            }

            const maybeSensorMetaAndData = getSensorMetaAndData(this.station);
            if (!maybeSensorMetaAndData) {
                console.log("tiny-chart:empty", { quickSensors, station: this.station });
                thisComp.handleNoDataCase();
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

            this.vizData = { vizSensor, timeRange: [queried.timeRange[0], queried.timeRange[1]] };
            this.busy = false;
        },
        handleNoDataCase() {
            this.noSensorData = true;
            this.busy = false;
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
    position: relative;

    .no-data {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        height: 100%;
        width: 100%;
        background-color: #e2e4e6;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 800;
    }
}

.tiny-chart .fit-x {
    width: 100%;
}

.tiny-chart .fit-y {
    height: 100%;
}
</style>
