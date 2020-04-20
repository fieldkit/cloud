<template>
    <div id="data-chart-container">
        <div v-if="noStation" class="size-20">
            <p>No station found.</p>
        </div>

        <div id="main-loading">
            <img alt="" src="../assets/progress.gif" />
        </div>

        <div class="white-bkgd" v-if="!noStation">
            <!-- compare and time window buttons -->
            <div id="chart-controls">
                <div id="compare-btn-container">
                    <div class="compare-btn" v-on:click="addChildChart">
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

            <div v-for="(chart, chartIndex) in noDataCharts" v-bind:key="chart.id" :class="chartIndex > 0 ? 'top-border' : ''">
                <div class="no-data">
                    <p>No data from {{ chart.station.name }} has been uploaded yet.</p>
                </div>
            </div>

            <!-- all data charts and their drop-down menus -->
            <div v-for="(chart, chartIndex) in charts" v-bind:key="chart.id" :class="chartIndex > 0 ? 'top-border' : ''">
                <div class="no-data" v-if="chart.noData">
                    <p>No data from {{ chart.station.name }} has been uploaded yet.</p>
                </div>
                <template v-if="!chart.noData">
                    <div :id="chart.id + '-loading'" class="loading">
                        <img alt="" src="../assets/progress.gif" />
                    </div>
                    <div v-if="chartIndex > 0" v-on:click="toggleLinkage">
                        <img v-if="linkedCharts" class="link-icon" src="../assets/link.png" />
                        <img v-if="!linkedCharts" class="open-link-icon" src="../assets/open_link.png" />
                    </div>
                    <div class="sensor-selection-dropdown">
                        <treeselect
                            :id="chart.id + 'treeselect'"
                            v-model="chart.treeSelectValue"
                            :showCount="true"
                            :options="treeSelectOptions"
                            :instanceId="chart.id"
                            open-direction="bottom"
                            @select="onTreeSelect"
                        >
                            <div slot="value-label" slot-scope="{ node }">{{ node.raw.customLabel }}</div>
                        </treeselect>
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
                        class="d3Chart"
                        :ref="chart.ref"
                        :chartParam="chart"
                        @timeZoomed="onTimeZoomed"
                        @unlinkCharts="unlinkCharts"
                        @zoomOut="setTimeRangeByDays"
                    />
                </template>
            </div>
            <div v-for="(pre, preIndex) in preloading" v-bind:key="pre.id" :class="preIndex > 0 ? 'top-border' : ''">
                <div :id="pre.id + '-loading'" class="loading">
                    <img alt="" src="../assets/progress.gif" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import _ from "lodash";
import * as d3 from "d3";
import D3Chart from "./D3Chart";
import Treeselect from "@riophae/vue-treeselect";
import "@riophae/vue-treeselect/dist/vue-treeselect.css";

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "DataChartControl",
    components: {
        D3Chart,
        Treeselect,
    },
    data: () => {
        return {
            noStation: false,
            charts: [],
            pending: [],
            preloading: [],
            noDataCharts: [],
            linkedCharts: false,
            urlQuery: {},
            prevQuery: {},
            chartOptions: [{ text: "Line", value: "Line" }, { text: "Histogram", value: "Histogram" }, { text: "Range", value: "Range" }],
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
    props: ["treeSelectOptions"],
    methods: {
        initialize() {
            const params = Object.keys(this.$route.query);
            let stationParams = [];
            let chartParams = {};
            params.forEach(p => {
                if (p.indexOf("start") > -1 || p.indexOf("end") > -1) {
                    this.$route.query[p] = parseInt(this.$route.query[p]);
                }
                if (p.indexOf("station") > 0) {
                    this.preloading.push({ id: p });
                }
                if (p != "stationId") {
                    this.urlQuery[p] = this.$route.query[p];
                    this.prevQuery[p] = this.$route.query[p];
                }

                if (p == "stationId") {
                    stationParams.push(p);
                    chartParams[p] = "chart-1";
                } else if (p.indexOf("station") > -1) {
                    stationParams.push(p);
                    const chartId = p.split("station")[0];
                    chartParams[p] = chartId;
                }
            });
            if (stationParams.length > 0) {
                stationParams.forEach(p => {
                    this.$emit("stationAdded", this.$route.query[p], chartParams[p]);
                });
            } else {
                this.noStationFound();
            }
        },
        noStationFound() {
            this.noStation = true;
            this.hideLoading();
        },
        noInitialDataFound(station, chartId) {
            const id = chartId.split("chart-")[1];
            const newChart = {
                id: chartId,
                ref: "d3Chart" + id,
                station: station,
                treeSelectValue: station.name,
            };
            this.noDataCharts.push(newChart);
            this.hideLoading();
        },
        noDataFound(station, chartId) {
            const chart = this.charts.find(c => {
                return c.id == chartId;
            });
            chart.station = station;
            chart.noData = true;
        },
        initializeChart(data, deviceId, chartId) {
            if (data[deviceId].data) {
                this.urlQuery[chartId + "station"] = data[deviceId].station.id.toString();
                this.addChartFromParams(data, deviceId, chartId);
            }
        },
        addChartFromParams(data, deviceId, chartId) {
            const id = chartId.split("chart-")[1];
            const chartData = data[deviceId].data;
            const totalTime = data[deviceId].timeRange;
            const timeCheck = this.checkTimeWindowForChart(chartId, totalTime, deviceId);
            const sensor = this.findSensorForChart(chartId, chartData, data[deviceId].sensors);
            const type = this.getChartType(chartId);
            const extent = d3.extent(chartData, d => {
                return d[sensor.key];
            });
            const filteredData = chartData.filter(d => {
                return d[sensor.key] === 0 || d[sensor.key];
            });
            const newChart = {
                id: chartId,
                ref: "d3Chart" + id,
                sensor: sensor,
                treeSelectValue: data[deviceId].station.name + sensor.key,
                type: type,
                parent: id == 1,
                data: filteredData,
                start: timeCheck.range[0],
                end: timeCheck.range[1],
                extent: extent,
                sensors: data[deviceId].sensors,
                station: data[deviceId].station,
                overall: chartData,
                current: chartData,
                totalTime: totalTime,
                noData: false,
            };
            if (timeCheck.addChart) {
                if (document.getElementById("main-loading")) {
                    document.getElementById("main-loading").style.display = "none";
                }
                this.charts.push(newChart);
            } else {
                this.pending.push(newChart);
            }
            if (document.getElementById(chartId + "station-loading")) {
                document.getElementById(chartId + "station-loading").style.display = "none";
            }
            this.updateRoute();

            // if (this.charts.length == 1) {
            //     this.$emit("switchedSensor", newChart.sensor);
            // }
        },
        addChildChart() {
            const id = this.makeChartId();
            let newChart = {};
            _.assign(newChart, this.charts[0]);
            newChart.id = "chart-" + id;
            newChart.ref = "d3Chart" + id;
            newChart.parent = false;
            const requested = this.$refs[this.charts[0].ref][0].getRequestedTime();
            if (requested) {
                newChart.requestedStart = requested[0];
                newChart.requestedEnd = requested[1];
            }
            this.charts.push(newChart);
            this.urlQuery[newChart.id + "station"] = newChart.station.id.toString();
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
        findSensorForChart(chartId, data, sensors) {
            const sensorKey = this.$route.query[chartId + "sensor"]
                ? this.$route.query[chartId + "sensor"]
                : this.getSensorWithData(data, sensors);
            const sensor = sensors.find(s => {
                return s.key == sensorKey;
            });
            this.urlQuery[chartId + "sensor"] = sensor.key;
            return sensor;
        },
        getSensorWithData(data, sensors) {
            // the first sensor with data
            let summaryKeys;
            for (let i = 0; i < data.length; i++) {
                summaryKeys = Object.keys(data[i]);
                // rule out ones that only have _filters, _resampled, and date:
                if (summaryKeys.length > 3) {
                    i = data.length;
                }
            }
            const actualSensors = _.intersectionBy(summaryKeys, sensors, s => {
                return s.key ? s.key : s;
            });
            if (actualSensors && actualSensors.length > 0) {
                return actualSensors[0];
            } else {
                return sensors[0].key;
            }
        },
        checkTimeWindowForChart(chartId, totalTime, deviceId) {
            let addChart = true;
            let start = totalTime[0];
            let end = totalTime[1];
            if (this.$route.query.start && this.$route.query.end) {
                start = new Date(parseInt(this.$route.query.start));
                end = new Date(parseInt(this.$route.query.end));
            }
            // but use chart's own time range, if defined
            if (this.$route.query[chartId + "start"] && this.$route.query[chartId + "end"]) {
                start = new Date(parseInt(this.$route.query[chartId + "start"]));
                end = new Date(parseInt(this.$route.query[chartId + "end"]));
                // fetch data if needed by triggering a time change
                if (start.getTime() != totalTime[0].getTime() || end.getTime() != totalTime[1].getTime()) {
                    // and flag this chart so it doesn't get drawn yet:
                    addChart = false;
                    const chartParam = { id: chartId, station: { device_id: deviceId } };
                    this.$emit("timeChanged", { start: start, end: end }, chartParam);
                    if (document.getElementById("main-loading")) {
                        document.getElementById("main-loading").style.display = "none";
                    }
                }
            }
            this.urlQuery[chartId + "start"] = start.getTime();
            this.urlQuery[chartId + "end"] = end.getTime();
            return { addChart: addChart, range: [start, end] };
        },
        getChartType(chartId) {
            const type = this.$route.query[chartId + "type"] ? this.$route.query[chartId + "type"] : "Line";
            this.urlQuery[chartId + "type"] = type;
            return type;
        },
        removeChart(event) {
            const id = event.target.getAttribute("data-id");
            const index = this.charts.findIndex(c => {
                return c.id == id;
            });
            if (index > -1) {
                this.charts.splice(index, 1);
                // remove all the associated params
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
                        c.treeSelectValue = this.charts[0].treeSelectValue;
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
        setTimeRangeByDays(event) {
            // method can be called by time buttons,
            // but also emitted by D3Chart, for zooming out
            // if emitted by D3Chart, arg will have 'id' property
            const days = event.id ? 0 : event.target.getAttribute("data-time");
            this.charts.forEach(c => {
                this.showLoading(c.id);
                const endDate = c.totalTime[1];
                let range = {
                    start: new Date(endDate.getTime() - days * DAY),
                    end: endDate,
                };
                if (days == 0) {
                    range.start = c.totalTime[0];
                }
                if (this.$refs[c.ref]) {
                    this.$refs[c.ref][0].setTimeRange(range);
                    this.$refs[c.ref][0].setRequestedTime(range);
                    this.urlQuery[c.id + "start"] = range.start.getTime();
                    this.urlQuery[c.id + "end"] = range.end.getTime();
                    this.$emit("timeChanged", range, c);
                }
            });
            this.updateRoute();
            // display active state for appropriate button
            this.timeButtons.forEach(b => {
                b.active = false;
                if (b.value == days) {
                    b.active = true;
                }
            });
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
                this.$emit("timeChanged", zoomed.range, chart);
            }
        },
        propagateTimeChange(range) {
            this.charts.forEach(c => {
                if (this.$refs[c.ref]) {
                    this.$refs[c.ref][0].setTimeRange(range);
                    this.$refs[c.ref][0].setRequestedTime(range);
                    this.urlQuery[c.id + "start"] = range.start.getTime();
                    this.urlQuery[c.id + "end"] = range.end.getTime();
                    this.$emit("timeChanged", range, c);
                }
            });
            this.updateRoute();
        },
        onTreeSelect(selection, instanceId) {
            this.showLoading(instanceId);
            const chart = this.charts.find(c => {
                return c.id == instanceId;
            });
            chart.noData = false;
            let sensor;
            if (selection.key) {
                sensor = chart.sensors.find(s => {
                    return s.key == selection.key;
                });
                this.urlQuery[chart.id + "sensor"] = selection.key;
                chart.sensor = sensor;
            }
            // if only sensor changed
            if (selection.key && selection.stationId == chart.station.id) {
                this.chartSensorChanged(sensor, chart);
            }
            // if station changed
            if (selection.stationId != chart.station.id) {
                this.$emit("stationChanged", selection.stationId, chart);
                this.urlQuery[chart.id + "station"] = selection.stationId;
            }
            this.updateRoute();
        },
        chartSensorChanged(selected, chart) {
            // the sensor type on a single chart instance has changed
            const filteredData = chart.current.filter(d => {
                return d[selected.key] === 0 || d[selected.key];
            });
            const extent = d3.extent(filteredData, d => {
                return d[selected.key];
            });
            // if it is the parent chart and they are linked, change all
            if (chart.parent && this.linkedCharts) {
                this.propagateSensorChange(selected, filteredData, extent);
            } else {
                // otherwise, change this one
                this.$refs[chart.ref][0].updateData(filteredData, extent, chart.sensor.colorScale);
                this.unlinkCharts();
            }
            this.updateRoute();
        },
        propagateSensorChange(selected, filteredData, extent) {
            this.charts.forEach(c => {
                c.sensor = selected;
                // TODO: some stations might not have this sensor
                c.treeSelectValue = c.station.name + selected.key;
                this.urlQuery[c.id + "sensor"] = selected.key;
                this.$refs[c.ref][0].updateData(filteredData, extent, c.sensor.colorScale);
            });
        },
        showLoading(chartId) {
            if (chartId && document.getElementById(chartId + "-loading")) {
                document.getElementById(chartId + "-loading").style.display = "block";
            } else {
                if (document.getElementById("main-loading")) {
                    document.getElementById("main-loading").style.display = "block";
                }
                const loadings = document.getElementsByClassName("loading");
                for (let i = 0; i < loadings.length; i++) {
                    loadings[i].style.display = "block";
                }
            }
        },
        hideLoading() {
            if (document.getElementById("main-loading")) {
                document.getElementById("main-loading").style.display = "none";
            }
            const loadings = document.getElementsByClassName("loading");
            for (let i = 0; i < loadings.length; i++) {
                loadings[i].style.display = "none";
            }
        },
        resetChart(data, deviceId, chartId) {
            const chart = this.charts.find(c => {
                return c.id == chartId;
            });
            // not checking pending... ?
            chart.sensors = data[deviceId].sensors;
            chart.station = data[deviceId].station;
            chart.overall = data[deviceId].data;
            chart.current = data[deviceId].data;
            chart.totalTime = data[deviceId].timeRange;
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
                    pendingChart.current = data;
                    pendingChart.extent = extent;
                    this.charts.push(pendingChart);
                    // if (this.charts.length == 1) {
                    //     // just for sensor summary in lower right corner
                    //     this.$emit("switchedSensor", pendingChart.sensor);
                    // }
                }
            } else {
                chart.current = data;
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
            // vue throws uncaught expression: Object if the same route is pushed again
            if (!_.isEqual(this.prevQuery, this.urlQuery)) {
                this.$router.push({ name: "viewData", query: this.urlQuery });
                this.prevQuery = {};
                const keys = Object.keys(this.urlQuery);
                keys.forEach(k => {
                    this.prevQuery[k] = this.urlQuery[k];
                });
            }
        },
        refresh() {
            // refresh window (back or forward browser button pressed)
            this.charts = [];
            this.pending = [];
            this.noStation = false;
            this.preloading = [];
            this.noDataCharts = [];
            this.urlQuery = {};
            this.prevQuery = {};
            this.showLoading();
            this.initialize();
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
.size-20 {
    font-size: 20px;
}
.no-data {
    font-size: 20px;
    float: left;
    clear: both;
    margin: 15px;
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
#compare-btn-container {
    float: left;
}
#time-control-container {
    float: right;
    margin-right: 10px;
}
.compare-btn {
    font-size: 12px;
    float: left;
    padding: 5px 10px;
    margin: 20px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
.compare-btn img {
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
    border-top: 1px solid rgb(230, 230, 230);
}
.sensor-selection-dropdown {
    float: left;
    clear: both;
    width: 350px;
    font-size: 14px;
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
