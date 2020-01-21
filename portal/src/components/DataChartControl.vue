<template>
    <div id="data-chart-container">
        <div v-if="this.noStation" class="no-data-message">
            <p>No station found.</p>
        </div>
        <div v-if="this.station && this.foundNoData" class="no-data-message">
            <p>No data from {{ station.name }} has been uploaded yet.</p>
        </div>
        <div id="readings-container" class="section" v-if="this.station && !this.loading">
            <div id="readings-label">
                Latest Reading <span class="synced">{{ getSyncedDate() }}</span>
            </div>
            <!-- top row sensor selection buttons -->
            <div id="reading-btns-container">
                <div
                    v-for="sensor in this.displaySensors"
                    v-bind:key="sensor.key"
                    :class="'reading' + (selectedSensor && selectedSensor.key == sensor.key ? ' active' : '')"
                    :data-key="sensor.key"
                    v-on:click="switchSelectedSensor"
                >
                    <div class="left">
                        <img
                            v-if="labels[sensor.key] == 'Temperature'"
                            alt="temperature icon"
                            src="../assets/Temp_icon.png"
                        />
                        {{ labels[sensor.key] ? labels[sensor.key] : sensor.key }}
                    </div>
                    <div class="right">
                        <span class="reading-value">
                            {{ sensor.currentReading ? sensor.currentReading.toFixed(1) : "--" }}
                        </span>
                        <span class="reading-unit">{{ sensor.units ? sensor.units : "" }}</span>
                    </div>
                </div>
            </div>
            <div id="left-arrow-container" v-if="this.displaySensors.length > 5">
                <img
                    v-on:click="showPrevSensor"
                    alt="left arrow"
                    src="../assets/left_arrow.png"
                    class="left-arrow"
                />
            </div>
            <div id="right-arrow-container" v-if="this.displaySensors.length > 5">
                <img
                    v-on:click="showNextSensor"
                    alt="right arrow"
                    src="../assets/right_arrow.png"
                    class="right-arrow"
                />
            </div>
        </div>
        <div id="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>

        <div class="white-bkgd" v-if="this.station && !this.loading">
            <!-- export/share/compare and time window buttons -->
            <div id="selected-sensor-controls">
                <div id="control-btn-container">
                    <div class="control-btn">
                        <img alt="" src="../assets/Export_icon.png" />
                        <span>Export</span>
                    </div>
                    <div class="control-btn">
                        <img alt="" src="../assets/Share_icon.png" />
                        <span>Share</span>
                    </div>
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
            <!-- all data charts and their drop-down menus -->
            <div v-for="(chart, chartIndex) in charts" v-bind:key="chart.id" class="top-border">
                <div v-if="chartIndex > 0" v-on:click="toggleLinkage">
                    <img v-if="linkedCharts" class="link-icon" src="../assets/link.png" />
                    <img v-if="!linkedCharts" class="open-link-icon" src="../assets/open_link.png" />
                </div>
                <div class="sensor-selection-dropdown">
                    <select v-model="chart.sensorOption" v-on:change="chartSensorChanged" :data-id="chart.id">
                        <option
                            v-for="sensor in displaySensors"
                            v-bind:value="sensor.key"
                            v-bind:key="sensor.key"
                        >
                            {{ labels[sensor.key] ? labels[sensor.key] : sensor.key }}
                        </option>
                    </select>
                </div>
                <div class="chart-type">
                    <select v-model="chart.typeOption" v-on:change="chartTypeChanged" :data-id="chart.id">
                        <option
                            v-for="option in chartOptions"
                            v-bind:value="option.value"
                            v-bind:key="option.value"
                        >
                            {{ option.text }}
                        </option>
                    </select>
                    <div class="remove-chart" v-if="chartIndex > 0">
                        <img
                            alt="Remove"
                            src="../assets/close.png"
                            :data-id="chart.id"
                            v-on:click="removeChart"
                        />
                    </div>
                </div>
                <D3Chart
                    :id="chart.id"
                    :ref="chart.ref"
                    :station="station"
                    :stationData="stationData"
                    :selectedSensor="chart.sensor"
                    :chartType="chart.type"
                    :parent="chart.parent"
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
import * as utils from "../utilities";
import D3Chart from "./D3Chart";

const DAY = 1000 * 60 * 60 * 24;

export default {
    name: "DataChartControl",
    components: {
        D3Chart
    },
    data: () => {
        return {
            foundNoData: false,
            loading: true,
            charts: [],
            stationData: [],
            linkedCharts: true,
            urlQuery: {},
            prevQuery: {},
            chartOptions: [
                { text: "Line", value: "Line" },
                { text: "Histogram", value: "Histogram" },
                { text: "Range", value: "Range" }
            ],
            timeRange: { start: 0, end: 0 },
            displaySensors: [],
            timeButtons: [
                {
                    active: false,
                    label: "Day",
                    value: 1
                },
                {
                    active: false,
                    label: "Week",
                    value: 7
                },
                {
                    active: false,
                    label: "2 Weeks",
                    value: 14
                },
                {
                    active: false,
                    label: "Month",
                    value: 31
                },
                {
                    active: false,
                    label: "Year",
                    value: 365
                },
                {
                    active: false,
                    label: "All",
                    value: 0
                }
            ]
        };
    },
    props: ["combinedStationInfo", "station", "labels", "noStation"],
    watch: {
        combinedStationInfo() {
            if (this.combinedStationInfo.stationData) {
                // this.fullData = this.combinedStationInfo.stationData;
                this.stationData = this.combinedStationInfo.stationData;
                this.displaySensors = this.combinedStationInfo.sensors;
                this.initialRange = [
                    this.stationData[0].date,
                    this.stationData[this.stationData.length - 1].date
                ];
                if (this.stationData.length > 0) {
                    this.initSelectedSensor();
                    this.initTimeWindow();
                    this.initChartType();
                    this.initCharts();
                    this.loading = false;
                } else {
                    this.foundNoData = true;
                    document.getElementById("loading").style.display = "none";
                }
            }
        }
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
            let id = this.charts.length + 1;
            if (this.charts.length > 0) {
                const lastChart = this.charts[this.charts.length - 1];
                const lastChartId = parseInt(lastChart.id.split("chart-")[1]);
                id = lastChartId + 1;
            }
            const chartId = "chart-" + id;
            const sensorOption = this.$route.query[chartId + "sensor"]
                ? this.$route.query[chartId + "sensor"]
                : this.selectedSensor.key;
            const sensor = this.displaySensors.find(s => {
                return s.key == sensorOption;
            });
            const type = this.$route.query[chartId + "type"] ? this.$route.query[chartId + "type"] : "Line";
            let range = this.timeRange;
            if (this.$route.query[chartId + "start"] && this.$route.query[chartId + "end"]) {
                range = {
                    start: new Date(parseInt(this.$route.query[chartId + "start"])),
                    end: new Date(parseInt(this.$route.query[chartId + "end"]))
                };
            }
            const newChart = {
                id: "chart-" + id,
                ref: "d3Chart" + id,
                sensor: sensor,
                sensorOption: sensorOption,
                type: type,
                typeOption: type,
                parent: id == 1
            };
            this.charts.push(newChart);
            // takes a moment for ref to become defined
            let interval = setInterval(() => {
                if (this.$refs[newChart.ref]) {
                    clearInterval(interval);
                    this.$refs[newChart.ref][0].initChild(range);
                }
            }, 250);
        },
        addChildChart() {
            const lastChart = this.charts[this.charts.length - 1];
            const lastChartId = parseInt(lastChart.id.split("chart-")[1]);
            const id = lastChartId + 1;
            const newChart = {
                id: "chart-" + id,
                ref: "d3Chart" + id,
                sensor: this.charts[0].sensor,
                sensorOption: this.charts[0].sensorOption,
                type: this.charts[0].type,
                typeOption: this.charts[0].type,
                parent: false
            };
            this.charts.push(newChart);
            this.urlQuery.numCharts = this.charts.length;
            this.updateRoute();
            // takes a moment for ref to become defined
            let interval = setInterval(() => {
                if (this.$refs[newChart.ref]) {
                    clearInterval(interval);
                    this.$refs[newChart.ref][0].initChild(this.timeRange);
                    document.getElementById(newChart.id).scrollIntoView({
                        behavior: "smooth"
                    });
                }
            }, 250);
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
                this.charts.forEach((c, i) => {
                    if (i > 0) {
                        c.type = this.charts[0].type;
                        c.typeOption = this.charts[0].type;
                        this.urlQuery[c.id + "type"] = this.charts[0].type;
                        c.sensor = this.charts[0].sensor;
                        c.sensorOption = this.charts[0].sensorOption;
                        this.urlQuery[c.id + "sensor"] = this.charts[0].sensorOption;
                        this.$refs[c.ref][0].setTimeRange(this.timeRange);
                        this.urlQuery[c.id + "start"] = this.timeRange.start.getTime();
                        this.urlQuery[c.id + "end"] = this.timeRange.end.getTime();
                    }
                });
                this.updateRoute();
            }
        },
        unlinkCharts() {
            this.linkedCharts = false;
        },
        initSelectedSensor() {
            // get selected sensor from url, if present
            if (this.$route.query.sensor) {
                this.displaySensors.forEach(s => {
                    if (s.key == this.$route.query.sensor) {
                        this.selectedSensor = s;
                    }
                });
            }
            // or set the first sensor to be selected sensor
            if (!this.$route.query.sensor || !this.selectedSensor) {
                this.selectedSensor = this.displaySensors[0];
            }
            // set chart sensors
            this.charts.forEach(c => {
                // use their own sensor choice, if defined
                if (this.$route.query[c.id + "sensor"]) {
                    const sensorKey = this.$route.query[c.id + "sensor"];
                    const sensor = this.displaySensors.find(s => {
                        return s.key == sensorKey;
                    });
                    c.sensor = sensor;
                    c.sensorOption = sensorKey;
                } else {
                    c.sensor = this.selectedSensor;
                    c.sensorOption = this.selectedSensor.key;
                }
            });
            this.$emit("switchedSensor", this.selectedSensor);
        },
        switchSelectedSensor(event) {
            const key = event.target.getAttribute("data-key");
            const sensor = this.displaySensors.find(s => {
                return s.key == key;
            });
            this.selectedSensor = sensor;
            // all charts get updated by main sensor
            this.charts.forEach(c => {
                c.sensor = this.selectedSensor;
                c.sensorOption = this.selectedSensor.key;
                this.urlQuery[c.id + "sensor"] = this.selectedSensor.key;
            });
            this.urlQuery.sensor = this.selectedSensor.key;
            this.updateRoute();
            this.$emit("switchedSensor", this.selectedSensor);
        },
        showNextSensor() {
            const first = this.displaySensors[0];
            this.displaySensors.splice(0, 1);
            this.displaySensors.push(first);
        },
        showPrevSensor() {
            const last = this.displaySensors[this.displaySensors.length - 1];
            this.displaySensors.splice(this.displaySensors.length - 1, 1);
            this.displaySensors.unshift(last);
        },
        chartSensorChanged(event) {
            // the sensor type on a single chart instance has changed
            const selected = this.displaySensors[event.target.selectedIndex];
            const id = event.target.getAttribute("data-id");
            const chart = this.charts.find(c => {
                return c.id == id;
            });
            chart.sensor = selected;
            this.urlQuery[chart.id + "sensor"] = selected.key;
            // if it is the parent chart and they are linked, change all
            if (chart.parent && this.linkedCharts) {
                this.charts.forEach((c, i) => {
                    if (i > 0) {
                        c.sensor = selected;
                        c.sensorOption = selected.key;
                        this.urlQuery[c.id + "sensor"] = selected.key;
                    }
                });
            } else {
                this.unlinkCharts();
            }
            this.updateRoute();
        },
        initTimeWindow() {
            let newStart = this.stationData[0].date;
            let newEnd = this.stationData[this.stationData.length - 1].date;
            if (this.$route.query.start && this.$route.query.end) {
                newStart = new Date(parseInt(this.$route.query.start));
                newEnd = new Date(parseInt(this.$route.query.end));
            }
            if (newStart != this.timeRange.start || newEnd != this.timeRange.end) {
                this.timeRange.start = newStart;
                this.timeRange.end = newEnd;
            }
            // use charts' own time range, if defined, otherwise, this.timeRange
            this.charts.forEach(c => {
                let range = this.timeRange;
                if (this.$route.query[c.id + "start"] && this.$route.query[c.id + "end"]) {
                    range = {
                        start: new Date(parseInt(this.$route.query[c.id + "start"])),
                        end: new Date(parseInt(this.$route.query[c.id + "end"]))
                    };
                }
                if (this.$refs[c.ref]) {
                    this.$refs[c.ref][0].setTimeRange(range);
                }
            });
            // check to see if we have enough data
            // and fetch more if needed by triggering a time change
            if (this.timeRange.start < this.initialRange[0]) {
                this.$emit("timeChanged", this.timeRange);
            }
        },
        setTimeRangeByDays(event) {
            // method can be called by time buttons,
            // but also emitted by D3Chart, for zooming out
            // if emitted by D3Chart, arg will have 'id' property
            const days = event.id ? 0 : event.target.getAttribute("data-time");
            const endDate = this.initialRange[1];
            this.timeRange = {
                start: new Date(endDate.getTime() - days * DAY),
                end: endDate
            };
            if (days == 0) {
                this.timeRange.start = this.initialRange[0];
            }
            if (event.id && !event.parent) {
                // if a D3Chart emitted this and they are not parent, only change them
                const chart = this.charts.find(c => {
                    return c.id == event.id;
                });
                if (this.$refs[chart.ref]) {
                    this.$refs[chart.ref][0].setTimeRange(this.timeRange);
                    this.urlQuery[chart.id + "start"] = this.timeRange.start.getTime();
                    this.urlQuery[chart.id + "end"] = this.timeRange.end.getTime();
                    this.updateRoute();
                }
                this.unlinkCharts();
                return;
            }
            // display active state for appropriate button
            this.timeButtons.forEach(b => {
                b.active = false;
                if (b.value == days) {
                    b.active = true;
                }
            });
            this.propagateTimeChange(true);
        },
        onTimeZoomed(zoomed) {
            if (zoomed.parent) {
                this.timeRange = zoomed.range;
                this.propagateTimeChange(this.linkedCharts);
            } else {
                // only update url for the chart that emitted this zoom
                const chart = this.charts.find(c => {
                    return c.id == zoomed.id;
                });
                this.urlQuery[chart.id + "start"] = zoomed.range.start.getTime();
                this.urlQuery[chart.id + "end"] = zoomed.range.end.getTime();
                this.updateRoute();
                this.unlinkCharts();
            }
        },
        propagateTimeChange(setForAll) {
            if (setForAll) {
                this.charts.forEach(c => {
                    if (this.$refs[c.ref]) {
                        this.$refs[c.ref][0].setTimeRange(this.timeRange);
                        this.urlQuery[c.id + "start"] = this.timeRange.start.getTime();
                        this.urlQuery[c.id + "end"] = this.timeRange.end.getTime();
                    }
                });
            }
            this.$emit("timeChanged", this.timeRange);
            this.urlQuery.start = this.timeRange.start.getTime();
            this.urlQuery.end = this.timeRange.end.getTime();
            this.updateRoute();
        },
        initChartType() {
            let type = "Line";
            // set chart type from url, if present
            if (this.$route.query.type) {
                type = this.$route.query.type;
            }
            this.charts.forEach(c => {
                // use their own chart type choice, if defined
                let urlType = this.$route.query[c.id + "type"];
                c.typeOption = urlType ? urlType : type;
                c.type = c.typeOption;
            });
        },
        chartTypeChanged() {
            // the chart type on a single chart instance has changed
            const selected = this.chartOptions[event.target.selectedIndex].text;
            const id = event.target.getAttribute("data-id");
            const chart = this.charts.find(c => {
                return c.id == id;
            });
            chart.type = selected;
            this.urlQuery[chart.id + "type"] = selected;
            // if it is the parent chart and they are linked, change all
            if (chart.parent && this.linkedCharts) {
                this.charts.forEach((c, i) => {
                    if (i > 0) {
                        c.type = selected;
                        c.typeOption = selected;
                        this.urlQuery[c.id + "type"] = selected;
                    }
                });
            } else {
                this.unlinkCharts();
            }
            this.updateRoute();
        },
        showLoading() {
            document.getElementById("loading").style.display = "block";
        },
        updateData(data) {
            this.stationData = data;
            // allow expansion of total time range
            if (this.stationData[0].date < this.initialRange[0]) {
                this.initialRange[0] = this.stationData[0].date;
            }
            if (this.stationData[1].date > this.initialRange[1]) {
                this.initialRange[1] = this.stationData[this.stationData.length - 1].date;
            }
            document.getElementById("loading").style.display = "none";
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

        updateNumberOfCharts() {
            return new Promise(resolve => {
                let extra = 0;
                let missing = 0;
                if (this.$route.query.numCharts) {
                    if (this.$route.query.numCharts != this.charts.length) {
                        if (this.$route.query.numCharts < this.charts.length) {
                            extra = this.charts.length - this.$route.query.numCharts;
                        } else if (this.$route.query.numCharts > this.charts.length) {
                            missing = this.$route.query.numCharts - this.charts.length;
                        }
                    }
                } else if (this.charts.length > 1) {
                    // no numCharts in params, but we have more than one
                    extra = this.charts.length - 1;
                }
                for (let i = 0; i < extra; i++) {
                    this.charts.splice(this.charts.length - 1, 1);
                }
                for (let i = 0; i < missing; i++) {
                    this.addChartFromParams();
                }
                if (missing > 0) {
                    // takes a moment for refs to become defined
                    setTimeout(resolve, 500);
                } else {
                    resolve();
                }
            });
        },
        refresh(data) {
            this.updateData(data);
            // refresh window (back or forward browser button pressed)
            this.updateNumberOfCharts().then(() => {
                this.initSelectedSensor();
                this.initTimeWindow();
                this.initChartType();
            });
        },
        prepareNewStation() {
            document.getElementById("loading").style.display = "block";
            this.foundNoData = false;
            this.loading = true;
            this.charts = [];
            this.linkedCharts = true;
            this.urlQuery = {};
            this.prevQuery = {};
        },
        getSyncedDate() {
            return "Last synced " + utils.getUpdatedDate(this.station);
        }
    }
};
</script>

<style scoped>
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
    position: absolute;
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
#readings-label {
    font-size: 20px;
    margin-bottom: 10px;
}
#reading-btns-container {
    max-width: 1110px;
    height: 60px;
    float: left;
    overflow: hidden;
}
#left-arrow-container,
#right-arrow-container {
    width: 40px;
    height: 50px;
    clear: both;
    margin-top: -60px;
    padding-top: 10px;
    cursor: pointer;
}
#left-arrow-container {
    float: left;
    margin-left: -5px;
    background: rgb(255, 255, 255);
    background: linear-gradient(90deg, rgba(255, 255, 255, 1) 0%, rgba(255, 255, 255, 0) 100%);
}
#right-arrow-container {
    float: right;
    background: rgb(255, 255, 255);
    background: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 1) 100%);
}
.right-arrow {
    float: right;
}
.reading {
    font-size: 12px;
    width: 190px;
    float: left;
    line-height: 20px;
    padding: 16px 10px 16px 18px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    cursor: pointer;
}
.reading.active {
    border-bottom: 3px solid #1b80c9;
    padding-bottom: 14px;
}
.reading-value {
    font-size: 16px;
}
.reading-unit {
    font-size: 11px;
}
.reading img {
    vertical-align: middle;
}
.left,
.right {
    pointer-events: none;
}
.left {
    float: left;
    width: 100px;
    height: 24px;
    line-height: 12px;
}
.right {
    float: right;
}
.white-bkgd {
    float: left;
    background-color: #ffffff;
}
#selected-sensor-controls {
    background-color: #ffffff;
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
    width: 1020px;
    margin: 0 0 20px 70px;
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
