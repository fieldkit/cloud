<template>
    <div id="sensor-summary-container" v-if="this.selectedSensor">
        <div id="sensor-title">{{ this.labels[this.selectedSensor.key] }} Statistics</div>
        <div class="sensor-summary-block" v-for="block in this.blocks" v-bind:key="block.label">
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
</template>

<script>
import * as d3 from "d3";

const DAY = 1000 * 60 * 60 * 24;
const WEEK = DAY * 7;
const MONTH = DAY * 30;

export default {
    name: "SensorSummary",
    props: ["selectedSensor", "stationData", "timeRange", "labels"],
    computed: {
        current: function() {
            if (this.stationData && this.stationData.length > 0) {
                let start = this.stationData[0].date;
                let end = this.stationData[this.stationData.length - 1].date;
                if (this.timeRange) {
                    start = this.timeRange.start;
                    end = this.timeRange.end;
                }
                let filtered = this.stationData.filter(d => {
                    return d.date >= start && d.date <= end;
                });
                return this.computeStats(filtered);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        week: function() {
            if (this.stationData && this.stationData.length > 0) {
                let end = this.stationData[this.stationData.length - 1].date;
                let start = new Date(end.getTime() - WEEK);
                let filtered = this.stationData.filter(d => {
                    return d.date >= start && d.date <= end;
                });
                return this.computeStats(filtered);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        month: function() {
            if (this.stationData && this.stationData.length > 0) {
                let end = this.stationData[this.stationData.length - 1].date;
                let start = new Date(end.getTime() - MONTH);
                let filtered = this.stationData.filter(d => {
                    return d.date >= start && d.date <= end;
                });
                return this.computeStats(filtered);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        overall: function() {
            if (this.stationData && this.stationData.length > 0) {
                return this.computeStats(this.stationData);
            } else {
                return { min: "--", max: "--", median: "--" };
            }
        },
        blocks: function() {
            return [
                { label: "CURRENT VIEW", stats: this.current },
                { label: "LAST 7 DAYS", stats: this.week },
                { label: "LAST 30 DAYS", stats: this.month },
                { label: "OVERALL", stats: this.overall },
            ];
        },
    },
    data: () => {
        return {};
    },
    methods: {
        computeStats: function(data) {
            if (data.length == 0) {
                return { min: "--", max: "--", median: "--" };
            }

            let sensorSummary = this;
            let extent = d3.extent(data, d => {
                return d[sensorSummary.selectedSensor.key];
            });
            if (!extent || (!extent[0] && extent[0] != 0)) {
                return { min: "--", max: "--", median: "--" };
            }
            let median = d3.median(data, d => {
                return d[sensorSummary.selectedSensor.key];
            });
            return { min: extent[0].toFixed(2), max: extent[1].toFixed(2), median: median.toFixed(2) };
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
