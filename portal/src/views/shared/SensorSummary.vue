<template>
    <div id="sensor-summary-container">
        <div v-for="sensor in sensors" v-bind:key="sensor.chartId">
            <div id="sensor-title">{{ sensor.label }} Statistics</div>
            <div class="sensor-summary-block" v-for="block in sensor.blocks" v-bind:key="block.label">
                <div class="block-heading">{{ block.label }}</div>
                <div class="block-reading">
                    <span class="left">Low</span>
                    <span class="right">{{ block.stats.min }}</span>
                </div>
                <div class="block-reading">
                    <span class="left">Median</span>
                    <span class="right">{{ block.stats.median }}</span>
                </div>
                <div class="block-reading">
                    <span class="left">High</span>
                    <span class="right">{{ block.stats.max }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import * as d3 from "d3";
import * as utils from "../utilities";

const DAY = 1000 * 60 * 60 * 24;
const WEEK = DAY * 7;
const MONTH = DAY * 30;

export default {
    name: "SensorSummary",
    props: ["allSensors"],
    data: () => {
        return {
            sensors: [],
        };
    },
    methods: {
        update(chart) {
            if (!chart.sensor) {
                return;
            }
            const meta = this.allSensors.find((s) => {
                return s.key == chart.sensor.key;
            });
            const sensor = this.sensors.find((s) => {
                return s.chartId == chart.id;
            });
            if (sensor) {
                // update sensor
                sensor.blocks = this.blocks(chart);
                sensor.label = meta ? meta.label : "";
            } else {
                // add sensor
                this.sensors.push({
                    chartId: chart.id,
                    blocks: this.blocks(chart),
                    label: meta ? meta.label : "",
                });
            }
        },
        remove(chartId) {
            const index = this.sensors.findIndex((s) => {
                return s.chartId == chartId;
            });
            if (index > -1) {
                this.sensors.splice(index, 1);
            }
        },
        blocks(dataSet) {
            return [
                { label: "CURRENT VIEW", stats: this.current(dataSet) },
                { label: "LAST 7 DAYS", stats: this.week(dataSet) },
                { label: "LAST 30 DAYS", stats: this.month(dataSet) },
                { label: "OVERALL", stats: this.overall(dataSet) },
            ];
        },
        current(dataSet) {
            if (dataSet.current && dataSet.current.length > 0) {
                return this.computeStats(dataSet.current, dataSet.sensor.key);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        week(dataSet) {
            if (dataSet.overall && dataSet.overall.length > 0) {
                const end = dataSet.overall[dataSet.overall.length - 1].date;
                const start = new Date(end.getTime() - WEEK);
                const filtered = dataSet.overall.filter((d) => {
                    return d.date >= start && d.date <= end;
                });
                return this.computeStats(filtered, dataSet.sensor.key);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        month(dataSet) {
            if (dataSet.overall && dataSet.overall.length > 0) {
                const end = dataSet.overall[dataSet.overall.length - 1].date;
                const start = new Date(end.getTime() - MONTH);
                const filtered = dataSet.overall.filter((d) => {
                    return d.date >= start && d.date <= end;
                });
                return this.computeStats(filtered, dataSet.sensor.key);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        overall(dataSet) {
            if (dataSet.overall && dataSet.overall.length > 0) {
                return this.computeStats(dataSet.overall, dataSet.sensor.key);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        computeStats(data, key) {
            if (data.length == 0) {
                return { min: "--", max: "--", median: "--" };
            }

            const extent = d3.extent(data, (d) => {
                return d[key];
            });
            if (!extent || (!extent[0] && extent[0] != 0)) {
                return { min: "--", max: "--", median: "--" };
            }
            const median = d3.median(data, (d) => {
                return d[key];
            });
            return { min: extent[0].toFixed(2), max: extent[1].toFixed(2), median: median.toFixed(2) };
        },
        getSensorName(module, sensor) {
            const newName = utils.convertOldFirmwareResponse(module);
            return this.$t(newName + ".sensors." + sensor.name);
        },
    },
};
</script>

<style scoped>
#sensor-summary-container {
    float: left;
    width: 360px;
    margin: 10px;
    padding: 10px;
    border: 1px solid rgb(230, 230, 230);
}
#sensor-title {
    font-weight: bold;
    margin: 10px 0;
}
.sensor-summary-block {
    background-color: rgb(244, 245, 247);
    padding: 4px 8px;
    float: left;
    width: 150px;
    margin: 6px;
}
.block-heading {
    font-size: 12px;
    margin-bottom: 10px;
}
.block-reading {
    display: inline-block;
    font-size: 14px;
    width: 100%;
    margin-bottom: 5px;
}
.left {
    float: left;
}
.right {
    float: right;
}
</style>
