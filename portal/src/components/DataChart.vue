<template>
    <div id="data-chart-container">
        <div v-if="this.stationData.length == 0" class="no-data-message">
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
                            v-if="sensor.name == 'Temperature'"
                            alt="temperature icon"
                            src="../assets/Temp_icon.png"
                        />
                        {{ sensor.name }}
                    </div>
                    <div class="right">
                        <span class="reading-value">{{ sensor.currentReading.toFixed(1) }} </span>
                        <span class="reading-unit">{{ sensor.unit ? sensor.unit : "" }}</span>
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
        </div>
        <div v-if="this.station">
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
                    <div class="time-btn">View By:</div>
                    <div class="time-btn">Day</div>
                    <div class="time-btn">Week</div>
                    <div class="time-btn">2 Week</div>
                    <div class="time-btn">Month</div>
                    <div class="time-btn">Year</div>
                    <div class="time-btn">All</div>
                    <div class="time-btn">Custom</div>
                </div>
                <div class="spacer"></div>
            </div>
            <div id="selected-sensor-label" v-if="this.selectedSensor">{{ selectedSensor.name }}</div>
            <div id="selected-sensor-graph"></div>
        </div>
    </div>
</template>

<script>
export default {
    name: "DataChart",
    data: () => {
        return {
            message: "",
            allSensors: []
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
        stationData: function(newVal, oldVal) {
            console.log("stationData", newVal, "was", oldVal);
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
            let sensor = this.allSensors.find(s => {
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
        }
    }
};
</script>

<style scoped>
#data-chart-container {
    width: 1200px;
    float: left;
    margin: 0;
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
    margin-top: -120px;
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
    padding: 18px 10px 18px 18px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
}
.reading.active {
    border-bottom: 3px solid #1b80c9;
    padding-bottom: 16px;
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
}
.right {
    float: right;
}
#selected-sensor-controls {
    background-color: #ffffff;
    float: left;
    width: 1110px;
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
.time-btn {
    font-size: 12px;
    float: left;
    margin: 25px 10px;
}
.spacer {
    float: left;
    width: 1020px;
    margin: 0 0 20px 70px;
    border-top: 1px solid rgba(230, 230, 230);
}
#selected-sensor-label {
    clear: both;
    margin-left: 70px;
}
#selected-sensor-graph {
    width: 1110px;
    height: 500px;
    float: left;
    background-color: #fbeef0;
}
</style>
