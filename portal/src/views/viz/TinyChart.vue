<template>
    <div class="tiny-chart">
        <LineChart :series="series" :settings="chartSettings" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import { ModuleSensor, DisplayStation } from "@/store";

import {
    ChartSettings,
    SeriesData,
    DataSetSeries,
    QueriedData,
    VizInfo,
    VizSensor,
    TimeRange,
    SensorDataResponse,
} from "./vega/SpecFactory";

import LineChart from "./vega/LineChart.vue";
import { SensorDataQuerier } from "../shared/LatestStationReadings.vue";

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
    async mounted() {
        const [stationData, quickSensors, meta] = await this.querier.query(this.stationId);

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
        function getSensorMetaAndData(station: DisplayStation) {
            if (quickSensors.station.filter((r) => r.moduleId != null).length == 0) {
                const moduleSensor: ModuleSensor | undefined = _.first(
                    _.flattenDeep(station.configurations.all.map((c) => c.modules.map((m) => m.sensors)))
                );
                if (!moduleSensor) {
                    return null;
                }

                const sensor = meta.findSensorByKey(moduleSensor.fullKey);
                return { vizSensor: [station.id, ["", 0]], sensor, data: [] };
            }

            const vizSensor = getFirstQuickSensor(quickSensors);
            const sensor = meta.findSensor(vizSensor);
            const data = stationData.data.filter((datum) => datum.sensorId == vizSensor[1][1]); // TODO VizSensor

            return { vizSensor, sensor, data };
        }

        const maybeSensorMetaAndData = getSensorMetaAndData(this.station);
        if (!maybeSensorMetaAndData) {
            console.log("tiny-chart:empty", { quickSensors, station: this.station });
            return;
        }

        const { vizSensor, sensor, data } = maybeSensorMetaAndData;

        const sdr: SensorDataResponse = {
            data: [], // _.cloneDeep(data),
        };
        const queried = new QueriedData("key", TimeRange.eternity, sdr);
        const name = sensor.strings["enUs"]["label"] || "Unknown";
        const vizInfo = new VizInfo(
            sensor.key,
            [], // Scale
            { name: "", location: null }, // Station
            sensor.unitOfMeasure,
            sensor.key,
            name,
            sensor.viz || [],
            sensor.ranges,
            "", // ChartLabel
            "" // AxisLabel
        );

        console.log("tiny-chart", { stationId: this.stationId, station: this.station, sensor, quickSensors, meta });

        this.series = [new SeriesData("key", new DataSetSeries(vizSensor, queried, queried), queried, vizInfo)];
    },
    computed: {},
    watch: {},
    methods: {},
});
</script>

<style scoped lang="scss">
.tiny-chart {
    margin-top: 10px;
    width: 100%;
    height: 200px;
    border: solid 1px #d8dce0;
}

.tiny-chart .fit-x {
    width: 100%;
}

.tiny-chart .fit-y {
    height: 100%;
}
</style>
