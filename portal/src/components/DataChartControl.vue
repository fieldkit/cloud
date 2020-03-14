<template>
    <div id="data-chart-container">
        <div v-if="this.noStation" class="no-data-message">
            <p>No station found.</p>
        </div>
        <div v-if="this.station && this.foundNoData" class="no-data-message">
            <p>No data from {{ station.name }} has been uploaded yet.</p>
        </div>

        <div id="main-loading">
            <img alt="" src="../assets/progress.gif" />
        </div>

        <div class="white-bkgd" v-if="this.station">
            <!-- compare and time window buttons -->
            <div id="chart-controls">
                <div id="control-btn-container">
                    <div class="control-btn" v-on:click="addChildChart">
                        <img alt="" src="../assets/Compare_icon.png" />
                        <span>Compare</span>
                    </div>
                </div>
                <div id="time-control-container">
                    <div class="time-btn-label">View By:</div>
                    <div
                        v-for="btn in timeButtons"
                        v-bind:key="btn.value"
                        :class="'time-btn' + (btn.active ? ' active' : '')"
                        :data-time="btn.value"
                        v-on:click="setTimeRangeByDays"
                    >
                        {{ btn.label }}
                    </div>
                </div>
                <div class="spacer"></div>
            </div>

            <div v-for="(chart, chartIndex) in pending" v-bind:key="chart.id" :class="chartIndex > 0 ? 'top-border' : ''">
                <div :id="chart.id + '-loading'" class="loading">
                    <img alt="" src="../assets/progress.gif" />
                </div>
            </div>

            <!-- all data charts and their drop-down menus -->
            <div v-for="(chart, chartIndex) in charts" v-bind:key="chart.id" :class="chartIndex > 0 ? 'top-border' : ''">
                <div :id="chart.id + '-loading'" class="loading">
                    <img alt="" src="../assets/progress.gif" />
                </div>
                <div v-if="chartIndex > 0" v-on:click="toggleLinkage">
                    <img v-if="linkedCharts" class="link-icon" src="../assets/link.png" />
                    <img v-if="!linkedCharts" class="open-link-icon" src="../assets/open_link.png" />
                </div>
                <div class="sensor-selection-dropdown">
                    <select v-model="chart.sensorOption" v-on:change="chartSensorChanged" :data-id="chart.id">
                        <option v-for="sensor in stationSensors" v-bind:value="sensor.key" v-bind:key="sensor.key">
                            {{ labels[sensor.key] ? labels[sensor.key] : sensor.key }}
                        </option>
                    </select>
                </div>
                <div class="chart-type">
                    <select v-model="chart.type" v-on:change="chartTypeChanged" :data-id="chart.id">
                        <option v-for="option in chartOptions" v-bind:value="option.value" v-bind:key="option.value">
                            {{ option.text }}
                        </option>
                    </select>
                    <div class="remove-chart" v-if="chartIndex > 0">
                        <img alt="Remove" src="../assets/close.png" :data-id="chart.id" v-on:click="removeChart" />
                    </div>
                </div>
                <D3Chart
                    :ref="chart.ref"
                    :chartParam="chart"
                    :station="station"
                    :stationSummary="stationSummary"
                    @timeZoomed="onTimeZoomed"
                    @unlinkCharts="unlinkCharts"
                    @zoomOut="setTimeRangeByDays"
                />
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as d3 from "d3";
import * as utils from "../utilities";
import D3Chart from "./D3Chart";

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "DataChartControl",
    components: {
        D3Chart,
    },
    data: () => {
        return {
            foundNoData: false,
            charts: [],
            pending: [],
            stationSummary: [],
            currentSummary: [],
            linkedCharts: true,
            urlQuery: {},
            prevQuery: {},
            chartOptions: [{ text: "Line", value: "Line" }, { text: "Histogram", value: "Histogram" }, { text: "Range", value: "Range" }],
            stationSensors: [],
            timeButtons: [
                {
                    active: false,
                    label: "Day",
                    value: 1,
                },
                {
                    active: false,
                    label: "Week",
                    value: 7,
                },
                {
                    active: false,
                    label: "2 Weeks",
                    value: 14,
                },
                {
                    active: false,
                    label: "Month",
                    value: 31,
                },
                {
                    active: false,
                    label: "Year",
                    value: 365,
                },
                {
                    active: false,
                    label: "All",
                    value: 0,
                },
            ],
        };
    },
    props: ["combinedStationInfo", "station", "labels", "noStation", "totalTime"],
    watch: {
        combinedStationInfo() {
            if (this.combinedStationInfo.stationData) {
                this.stationSummary = this.combinedStationInfo.stationData;
                this.currentSummary = this.combinedStationInfo.stationData;
                this.stationSensors = this.combinedStationInfo.sensors;
                if (this.stationSummary.length > 0) {
                    this.initCharts();
                } else {
                    this.foundNoData = true;
                    this.hideLoading();
                }
            }
        },
    },
    mounted() {
        const keys = Object.keys(this.$route.query);
        keys.forEach(k => {
            if (k.indexOf("start") > -1 || k.indexOf("end") > -1) {
                this.$route.query[k] = parseInt(this.$route.query[k]);
            }
            this.urlQuery[k] = this.$route.query[k];
            this.prevQuery[k] = this.$route.query[k];
        });
    },
    methods: {
        initCharts() {
            if (this.$route.query.numCharts) {
                for (let i = 0; i < this.$route.query.numCharts; i++) {
                    // unlinked by default when adding via params
                    this.linkedCharts = false;
                    this.addChartFromParams();
                }
            } else {
                // only one chart
                this.linkedCharts = true;
                this.addChartFromParams();
            }
        },
        addChartFromParams() {
            const id = this.makeChartId();
            const chartId = "chart-" + id;
            const timeCheck = this.checkTimeWindowForChart(chartId);
            const sensor = this.findSensorForChart(chartId);
            const type = this.$route.query[chartId + "type"] ? this.$route.query[chartId + "type"] : "Line";
            const extent = d3.extent(this.stationSummary, d => {
                return d[sensor.key];
            });
            const filteredData = this.stationSummary.filter(d => {
                return d[sensor.key] === 0 || d[sensor.key];
            });
            const newChart = {
                id: chartId,
                ref: "d3Chart" + id,
                sensor: sensor,
                sensorOption: sensor.key,
                type: type,
                parent: id == 1,
                data: filteredData,
                start: timeCheck.range[0],
                end: timeCheck.range[1],
                extent: extent,
            };
            if (timeCheck.addChart) {
                document.getElementById("main-loading").style.display = "none";
                this.charts.push(newChart);
            } else {
                this.pending.push(newChart);
            }
            if (this.charts.length == 1) {
                this.$emit("switchedSensor", newChart.sensor);
            }
        },
        addChildChart() {
            const id = this.makeChartId();
            const newChart = {
                id: "chart-" + id,
                ref: "d3Chart" + id,
                sensor: this.charts[0].sensor,
                sensorOption: this.charts[0].sensorOption,
                type: this.charts[0].type,
                parent: false,
                data: this.charts[0].data,
                start: this.charts[0].start,
                end: this.charts[0].end,
                extent: this.charts[0].extent,
            };
            const requested = this.$refs[this.charts[0].ref][0].getRequestedTime();
            if (requested) {
                newChart.requestedStart = requested[0];
                newChart.requestedEnd = requested[1];
            }
            this.charts.push(newChart);
            this.urlQuery.numCharts = this.charts.length;
            this.urlQuery[newChart.id + "type"] = newChart.type;
            this.urlQuery[newChart.id + "sensor"] = newChart.sensor.key;
            this.urlQuery[newChart.id + "start"] = newChart.start.getTime();
            this.urlQuery[newChart.id + "end"] = newChart.end.getTime();
            this.updateRoute();
        },
        makeChartId() {
            let id = this.charts.length + this.pending.length + 1;
            if (this.charts.length > 0 || this.pending.length > 0) {
                const lastChart = this.charts[this.charts.length - 1];
                const lastChartId = lastChart ? parseInt(lastChart.id.split("chart-")[1]) : 0;
                const lastPending = this.pending[this.pending.length - 1];
                const lastPendingId = lastPending ? parseInt(lastPending.id.split("chart-")[1]) : 0;
                if (lastChartId > lastPendingId && lastChartId > id) {
                    id = lastChartId + 1;
                } else if (lastPendingId > lastChartId && lastPendingId > id) {
                    id = lastPendingId + 1;
                }
            }
            return id;
        },
        findSensorForChart(chartId) {
            const sensorKey = this.$route.query[chartId + "sensor"] ? this.$route.query[chartId + "sensor"] : this.getSensorWithData();
            const sensor = this.stationSensors.find(s => {
                return s.key == sensorKey;
            });
            return sensor;
        },
        getSensorWithData() {
            // the first sensor with data
            let summaryKeys;
            for (let i = 0; i < this.stationSummary.length; i++) {
                summaryKeys = Object.keys(this.stationSummary[i]);
                // rule out ones that only have _filters, _resampled, and date:
                if (summaryKeys.length > 3) {
                    i = this.stationSummary.length;
                }
            }
            const actualSensors = _.intersectionBy(summaryKeys, this.stationSensors, s => {
                return s.key ? s.key : s;
            });
            if (actualSensors && actualSensors.length > 0) {
                return actualSensors[0];
            } else {
                return this.stationSensors[0];
            }
        },
        checkTimeWindowForChart(chartId) {
            let addChart = true;
            let start = this.totalTime[0];
            let end = this.totalTime[1];
            if (this.$route.query.start && this.$route.query.end) {
                start = new Date(parseInt(this.$route.query.start));
                end = new Date(parseInt(this.$route.query.end));
            }
            // but use chart's own time range, if defined
            if (this.$route.query[chartId + "start"] && this.$route.query[chartId + "end"]) {
                start = new Date(parseInt(this.$route.query[chartId + "start"]));
                end = new Date(parseInt(this.$route.query[chartId + "end"]));
                // fetch data if needed by triggering a time change
                if (start.getTime() != this.totalTime[0].getTime() || end.getTime() != this.totalTime[1].getTime()) {
                    // and flag this chart so it doesn't get drawn yet:
                    addChart = false;
                    this.$emit("chartTimeChanged", { start: start, end: end }, chartId);
                    document.getElementById("main-loading").style.display = "none";
                }
            }
            return { addChart: addChart, range: [start, end] };
        },
        removeChart(event) {
            const id = event.target.getAttribute("data-id");
            const index = this.charts.findIndex(c => {
                return c.id == id;
            });
            if (index > -1) {
                this.charts.splice(index, 1);
                this.urlQuery.numCharts = this.charts.length;
                // also remove all the associated params
                let keys = Object.keys(this.urlQuery);
                keys.forEach(k => {
                    if (k.indexOf(id) == 0) {
                        delete this.urlQuery[k];
                    }
                });
                this.updateRoute();
            }
        },
        toggleLinkage() {
            this.linkedCharts = !this.linkedCharts;
            if (this.linkedCharts) {
                // now re-linked, reset all to parent
                const parentTime = this.$refs[this.charts[0].ref][0].getTimeRange();
                const requested = this.$refs[this.charts[0].ref][0].getRequestedTime();
                this.charts.forEach((c, i) => {
                    if (i > 0) {
                        c.type = this.charts[0].type;
                        c.sensor = this.charts[0].sensor;
                        c.sensorOption = this.charts[0].sensorOption;
                        this.$refs[c.ref][0].setTimeRange(parentTime);
                        if (requested) {
                            const requestedRange = { start: requested[0], end: requested[1] };
                            this.$refs[c.ref][0].setRequestedTime(requestedRange);
                        }
                        this.$refs[c.ref][0].updateChartType();
                        this.$refs[c.ref][0].updateData(this.charts[0].data, this.charts[0].extent, this.charts[0].sensor.colorScale);
                        this.urlQuery[c.id + "type"] = this.charts[0].type;
                        this.urlQuery[c.id + "sensor"] = this.charts[0].sensor.key;
                        this.urlQuery[c.id + "start"] = parentTime.start.getTime();
                        this.urlQuery[c.id + "end"] = parentTime.end.getTime();
                    }
                });
                this.updateRoute();
            }
        },
        unlinkCharts() {
            this.linkedCharts = false;
        },
        chartTypeChanged() {
            // the chart type on a single chart instance has changed
            const selected = this.chartOptions[event.target.selectedIndex].text;
            const id = event.target.getAttribute("data-id");
            const chart = this.charts.find(c => {
                return c.id == id;
            });
            chart.type = selected;
            this.$refs[chart.ref][0].updateChartType();
            this.urlQuery[chart.id + "type"] = selected;
            // if it is the parent chart and they are linked, change all
            if (chart.parent && this.linkedCharts) {
                this.charts.forEach((c, i) => {
                    if (i > 0) {
                        c.type = selected;
                        this.$refs[c.ref][0].updateChartType();
                        this.urlQuery[c.id + "type"] = selected;
                    }
                });
            } else {
                this.unlinkCharts();
            }
            this.updateRoute();
        },
        chartSensorChanged(event) {
            // the sensor type on a single chart instance has changed
            const selected = this.stationSensors[event.target.selectedIndex];
            const id = event.target.getAttribute("data-id");
            const chart = this.charts.find(c => {
                return c.id == id;
            });
            const filteredData = this.currentSummary.filter(d => {
                return d[selected.key] === 0 || d[selected.key];
            });
            const extent = d3.extent(filteredData, d => {
                return d[selected.key];
            });
            const chartTime = this.$refs[chart.ref][0].getTimeRange();

            // set up some time range checks to see if we need to refetch data
            const summaryStart = this.currentSummary[0].date;
            const summaryEnd = this.currentSummary[this.currentSummary.length - 1].date;
            const summaryStartMatch = chart.start.getTime() == summaryStart.getTime();
            const summaryEndMatch = chart.end.getTime() == summaryEnd.getTime();
            const changeTime = !summaryStartMatch || !summaryEndMatch;

            // if it is the parent chart and they are linked, change all
            if (chart.parent && this.linkedCharts) {
                this.charts.forEach(c => {
                    c.sensor = selected;
                    c.sensorOption = selected.key;
                    this.urlQuery[c.id + "sensor"] = selected.key;
                    if (!changeTime) {
                        this.$refs[c.ref][0].updateData(filteredData, extent, c.sensor.colorScale);
                    }
                });
                if (changeTime) {
                    this.showLoading();
                    this.$emit("timeChanged", chartTime);
                }
            } else {
                // otherwise, change this one, including time if needed
                chart.sensor = selected;
                chart.sensorOption = selected.key;
                this.urlQuery[chart.id + "sensor"] = selected.key;
                if (!changeTime) {
                    this.$refs[chart.ref][0].updateData(filteredData, extent, chart.sensor.colorScale);
                } else {
                    this.showLoading(chart.id);
                    this.$emit("chartTimeChanged", chartTime, chart.id);
                }
                this.unlinkCharts();
            }
            this.updateRoute();
        },
        setTimeRangeByDays(event) {
            this.showLoading();
            // method can be called by time buttons,
            // but also emitted by D3Chart, for zooming out
            // if emitted by D3Chart, arg will have 'id' property
            const days = event.id ? 0 : event.target.getAttribute("data-time");
            // TODO: consider expanding endDate to capture any new data since page loaded
            const endDate = this.totalTime[1];
            let range = {
                start: new Date(endDate.getTime() - days * DAY),
                end: endDate,
            };
            if (days == 0) {
                range.start = this.totalTime[0];
            }
            if (event.id && !event.parent) {
                // if a D3Chart emitted this and they are not parent,
                // consider changing this to only affect them
                const chart = this.charts.find(c => {
                    return c.id == event.id;
                });
                if (this.$refs[chart.ref]) {
                    this.$refs[chart.ref][0].setTimeRange(range);
                    this.urlQuery[chart.id + "start"] = range.start.getTime();
                    this.urlQuery[chart.id + "end"] = range.end.getTime();
                    this.updateRoute();
                }
                // allow ranges to affect everything for now
                // this.unlinkCharts();
                // return;
            }
            // display active state for appropriate button
            this.timeButtons.forEach(b => {
                b.active = false;
                if (b.value == days) {
                    b.active = true;
                }
            });
            this.propagateTimeChange(range);
        },
        onTimeZoomed(zoomed) {
            if (zoomed.parent && this.linkedCharts) {
                this.showLoading();
                this.propagateTimeChange(zoomed.range);
            } else {
                // only update url for the chart that emitted this zoom
                this.showLoading(zoomed.id);
                const chart = this.charts.find(c => {
                    return c.id == zoomed.id;
                });
                this.$refs[chart.ref][0].setTimeRange(zoomed.range);
                this.urlQuery[chart.id + "start"] = zoomed.range.start.getTime();
                this.urlQuery[chart.id + "end"] = zoomed.range.end.getTime();
                this.updateRoute();
                this.unlinkCharts();
                this.$emit("chartTimeChanged", zoomed.range, zoomed.id);
            }
        },
        propagateTimeChange(range) {
            this.charts.forEach(c => {
                if (this.$refs[c.ref]) {
                    this.$refs[c.ref][0].setTimeRange(range);
                    this.$refs[c.ref][0].setRequestedTime(range);
                    this.urlQuery[c.id + "start"] = range.start.getTime();
                    this.urlQuery[c.id + "end"] = range.end.getTime();
                }
            });
            this.$emit("timeChanged", range);
            this.urlQuery.start = range.start.getTime();
            this.urlQuery.end = range.end.getTime();
            this.updateRoute();
        },
        showLoading(chartId) {
            if (chartId) {
                document.getElementById(chartId + "-loading").style.display = "block";
            } else {
                const loadings = document.getElementsByClassName("loading");
                for (let i = 0; i < loadings.length; i++) {
                    loadings[i].style.display = "block";
                }
            }
        },
        hideLoading() {
            const loadings = document.getElementsByClassName("loading");
            for (let i = 0; i < loadings.length; i++) {
                loadings[i].style.display = "none";
            }
        },
        updateAll(data) {
            this.currentSummary = data;
            const range =
                data.length == 0
                    ? this.$refs[this.charts[0].ref][0].getTimeRange()
                    : {
                          start: data[0].date,
                          end: data[data.length - 1].date,
                      };
            // respond to a global time change event
            this.charts.forEach(c => {
                const filteredData = data.filter(d => {
                    return d[c.sensor.key] === 0 || d[c.sensor.key];
                });
                const extent = d3.extent(filteredData, d => {
                    return d[c.sensor.key];
                });
                this.$refs[c.ref][0].setTimeRange(range);
                this.$refs[c.ref][0].updateData(filteredData, extent, c.sensor.colorScale);
            });
        },
        updateChartData(data, chartId) {
            // only update url for the chart that emitted this zoom
            const chart = this.charts.find(c => {
                return c.id == chartId;
            });
            if (!chart) {
                // check pending
                const pendingIndex = this.pending.findIndex(c => {
                    return c.id == chartId;
                });
                if (pendingIndex > -1) {
                    const pendingChart = this.pending.splice(pendingIndex, 1)[0];
                    const range =
                        data.length == 0
                            ? { start: pendingChart.start, end: pendingChart.end }
                            : {
                                  start: data[0].date,
                                  end: data[data.length - 1].date,
                              };
                    const filteredData = data.filter(d => {
                        return d[pendingChart.sensor.key] === 0 || d[pendingChart.sensor.key];
                    });
                    const extent = d3.extent(filteredData, d => {
                        return d[pendingChart.sensor.key];
                    });
                    pendingChart.start = range.start;
                    pendingChart.end = range.end;
                    pendingChart.data = filteredData;
                    pendingChart.extent = extent;
                    this.charts.push(pendingChart);
                    if (this.charts.length == 1) {
                        // just for sensor summary in lower right corner
                        this.$emit("switchedSensor", pendingChart.sensor);
                    }
                }
            } else {
                const range =
                    data.length == 0
                        ? { start: chart.start, end: chart.end }
                        : {
                              start: data[0].date,
                              end: data[data.length - 1].date,
                          };
                const filteredData = data.filter(d => {
                    return d[chart.sensor.key] === 0 || d[chart.sensor.key];
                });
                const extent = d3.extent(filteredData, d => {
                    return d[chart.sensor.key];
                });
                this.$refs[chart.ref][0].setTimeRange(range);
                this.$refs[chart.ref][0].updateData(filteredData, extent, chart.sensor.colorScale);
                this.unlinkCharts();
            }
        },
        updateRoute() {
            // temp temp temp
            // just for demo of Ancient Goose, keep id = 0
            // temp temp temp
            let stationId = this.station.id;
            if (stationId == 159) {
                stationId = 0;
            }
            // vue throws uncaught expression: Object if the same route is pushed again
            if (!_.isEqual(this.prevQuery, this.urlQuery)) {
                this.$router.push({ name: "viewData", params: { id: stationId }, query: this.urlQuery });
                this.prevQuery = {};
                const keys = Object.keys(this.urlQuery);
                keys.forEach(k => {
                    this.prevQuery[k] = this.urlQuery[k];
                });
            }
        },
        refresh(data) {
            // refresh window (back or forward browser button pressed)
            this.stationSummary = data;
            this.currentSummary = data;
            this.charts = [];
            this.pending = [];
            if (this.stationSummary.length > 0) {
                this.initCharts();
            } else {
                this.foundNoData = true;
                this.hideLoading();
            }
        },
        prepareNewStation() {
            this.showLoading();
            this.foundNoData = false;
            this.charts = [];
            this.pending = [];
            this.stationSensors = [];
            this.stationSummary = [];
            this.currentSummary = [];
            this.linkedCharts = true;
            this.urlQuery = {};
            this.prevQuery = {};
        },
        getSyncedDate() {
            return "Last synced " + utils.getUpdatedDate(this.station);
        },
    },
};
</script>

<style scoped>
#main-loading,
.loading {
    width: 1115px;
    height: 590px;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
    position: absolute;
}
#main-loading {
    margin-top: 20px;
}
#data-chart-container {
    float: left;
    margin-right: 20px;
}
.no-data-message {
    font-size: 20px;
}
.synced {
    margin-left: 10px;
    font-size: 14px;
}
.white-bkgd {
    float: left;
    margin-top: 20px;
    background-color: #ffffff;
    border-radius: 4px;
    border: 1px solid rgb(235, 235, 235);
}
#chart-controls {
    float: left;
    clear: both;
}
#control-btn-container {
    float: left;
}
#time-control-container {
    float: right;
    margin-right: 10px;
}
.control-btn {
    font-size: 12px;
    float: left;
    padding: 5px 10px;
    margin: 20px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
.control-btn img {
    vertical-align: middle;
    margin-right: 10px;
}
.time-btn-label,
.time-btn {
    font-size: 12px;
    float: left;
    margin: 25px 10px;
}
.time-btn {
    cursor: pointer;
}
.time-btn.active {
    text-decoration: underline;
}
.spacer {
    float: left;
    width: 1070px;
    margin: 0 0 20px 10px;
    border-top: 1px solid rgba(230, 230, 230);
}
.sensor-selection-dropdown {
    float: left;
    clear: both;
    margin-bottom: 10px;
    margin-left: 10px;
}
.chart-type {
    float: right;
    margin-right: 15px;
}
.sensor-selection-dropdown select,
.chart-type select {
    font-size: 16px;
    border: 1px solid lightgray;
    border-radius: 4px;
    padding: 2px 4px;
}
.top-border {
    border-top: 1px solid rgb(215, 220, 225);
}
.link-icon,
.open-link-icon {
    width: 50px;
    height: 50px;
    opacity: 0.25;
    margin-left: 525px;
    margin-top: -25px;
    margin-bottom: -5px;
    cursor: pointer;
}
.remove-chart {
    float: right;
    margin: -20px 60px 0 20px;
    cursor: pointer;
    padding: 2px;
}
</style>
