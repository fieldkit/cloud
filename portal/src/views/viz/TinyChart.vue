<template>
    <div class="tiny-chart">
        <LineChart :series="series" :settings="chartSettings" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";

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
        const [data, quickSensors, meta] = await this.querier.query(this.stationId);
        if (quickSensors.station.filter((r) => r.moduleId != null).length == 0) {
            return;
        }

        const sensorModuleId = quickSensors.station[0].moduleId;
        const sensorId = quickSensors.station[0].sensorId;
        const vizSensor: VizSensor = [this.stationId, [sensorModuleId, sensorId]];

        const sensorData = data.data.filter((datum) => datum.sensorId == sensorId);

        const sdr: SensorDataResponse = {
            data: _.cloneDeep(sensorData),
        };
        const queried = new QueriedData("key", TimeRange.eternity, sdr);
        const vizInfo = new VizInfo(
            "key",
            [],
            { name: "", location: null },
            "unit",
            "firmwareKey",
            "name",
            [],
            [],
            "chartLabel",
            "axisLabel"
        );

        console.log("tiny-chart", { stationId: this.stationId, vizSensor, quickSensors, meta, sensorData });

        this.series = [new SeriesData("key", new DataSetSeries(vizSensor, queried, queried), queried, vizInfo)];
    },
    computed: {},
    watch: {},
    methods: {},
});
</script>

<style scoped lang="scss">
.tiny-chart {
    height: 150px;
}

.tiny-chart .fit-y {
    height: 100%;
}
</style>
