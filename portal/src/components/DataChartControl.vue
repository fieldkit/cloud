<template>
    <div id="data-chart-container">
        <div v-if="!this.station" class="no-data-message">
            <p>No data yet.</p>
        </div>
        <div id="readings-container" class="section" v-if="this.station">
            <div id="readings-label">
                Latest Reading <span class="synced">Last synced {{ getSyncedDate() }}</span>
            </div>
            <div id="reading-btns-container">
                <div
                    v-for="sensor in this.allSensors"
                    v-bind:key="sensor.id"
                    :class="'reading' + (selectedSensor && selectedSensor.id == sensor.id ? ' active' : '')"
                    :data-id="sensor.id"
                    v-on:click="switchSensor"
                >
                    <div class="left">
                        <img
                            v-if="labels[sensor.name] == 'Temperature'"
                            alt="temperature icon"
                            src="../assets/Temp_icon.png"
                        />
                        {{ labels[sensor.name] ? labels[sensor.name] : sensor.name }}
                    </div>
                    <div class="right">
                        <span class="reading-value">
                            {{ sensor.currentReading ? sensor.currentReading.toFixed(1) : "ncR" }}
                        </span>
                        <span class="reading-unit">{{ sensor.unit ? sensor.unit : "" }}</span>
                    </div>
                </div>
            </div>
            <div id="left-arrow-container">
                <img
                    v-if="this.allSensors.length > 5"
                    v-on:click="showPrevSensor"
                    alt="left arrow"
                    src="../assets/left_arrow.png"
                    class="left-arrow"
                />
            </div>
            <div id="right-arrow-container">
                <img
                    v-if="this.allSensors.length > 5"
                    v-on:click="showNextSensor"
                    alt="right arrow"
                    src="../assets/right_arrow.png"
                    class="right-arrow"
                />
            </div>
        </div>
        <div class="white-bkgd" v-if="this.station">
            <div id="selected-sensor-controls">
                <div id="control-btn-container">
                    <div id="" class="control-btn">
                        <img alt="" src="../assets/Export_icon.png" />
                        <span>Export</span>
                    </div>
                    <div id="" class="control-btn">
                        <img alt="" src="../assets/Share_icon.png" />
                        <span>Share</span>
                    </div>
                    <div id="" class="control-btn">
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
                        v-on:click="updateTime"
                    >
                        {{ btn.label }}
                    </div>
                </div>
                <div class="spacer"></div>
            </div>
            <div id="selected-sensor-label" v-if="this.selectedSensor">
                {{ labels[selectedSensor.name] ? labels[selectedSensor.name] : selectedSensor.name }}
            </div>
            <div id="chart-type">
                <select v-model="selected" v-on:change="chartTypeChanged">
                    <option v-for="option in options" v-bind:value="option.value" v-bind:key="option.value">
                        {{ option.text }}
                    </option>
                </select>
            </div>
            <D3Chart
                ref="d3Chart"
                :stationData="stationData"
                :selectedSensor="selectedSensor"
                :timeRange="timeRange"
                :chartType="chartType"
            />
        </div>
    </div>
</template>

<script>
import D3Chart from "./D3Chart";

export default {
    name: "DataChartControl",
    components: {
        D3Chart
    },
    data: () => {
        return {
            message: "",
            selected: "Line",
            chartType: "Line",
            options: [
                { text: "Line", value: "Line" },
                { text: "Histogram", value: "Histogram" }
                // { text: "Range", value: "Range" }
            ],
            allSensors: [],
            timeRange: 0,
            // temporary label system
            labels: {
                ph: "pH",
                do: "Dissolved Oxygen",
                ec: "Electrical Conductivity",
                tds: "Total Dissolved Solids",
                salinity: "Salinity",
                temp: "Temperature"
            },
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
    props: ["summary", "stationData", "station", "selectedSensor"],

    watch: {
        summary: function(newVal, oldVal) {
            // if(newVal.provisions.length == 0) {
            //     this.message = "No data files uploaded for " + this.station.name + " yet.";
            // }
            console.log("summary", newVal.provisions, "was", oldVal);
        },
        station: function(_station) {
            _station.status_json.moduleObjects.forEach(m => {
                m.sensorObjects.forEach(s => {
                    this.allSensors.push(s);
                });
            });
        }
    },
    methods: {
        getSyncedDate() {
            const date = this.station.status_json.updated;
            const d = new Date(date);
            return d.toLocaleDateString("en-US");
        },
        switchSensor(event) {
            const id = event.target.getAttribute("data-id");
            const sensor = this.allSensors.find(s => {
                return s.id == id;
            });
            this.$emit("switchedSensor", sensor);
        },
        showNextSensor() {
            const first = this.allSensors[0];
            this.allSensors.splice(0, 1);
            this.allSensors.push(first);
        },
        showPrevSensor() {
            const last = this.allSensors[this.allSensors.length - 1];
            this.allSensors.splice(this.allSensors.length - 1, 1);
            this.allSensors.unshift(last);
        },
        updateTime(event) {
            const time = event.target.getAttribute("data-time");
            this.timeButtons.forEach(b => {
                b.active = false;
                if (b.value == time) {
                    b.active = true;
                }
            });
            // HACK setting to -1 first to always trigger a change event
            this.timeRange = -1;
            this.timeRange = time;
        },
        chartTypeChanged() {
            this.chartType = this.selected;
        }
    }
};
</script>

<style scoped>
#data-chart-container {
    width: 1100px;
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
    background-color: #ffffff;
}
#selected-sensor-controls {
    background-color: #ffffff;
    float: left;
    clear: both;
}
#control-btn-container {
    margin-left: 60px;
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
#selected-sensor-label {
    float: left;
    clear: both;
    margin-bottom: 10px;
    margin-left: 70px;
}
#chart-type {
    float: right;
}
#chart-type select {
    font-size: 16px;
    border: 1px solid lightgray;
    border-radius: 4px;
    padding: 2px 4px;
}
</style>
