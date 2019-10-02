<template>
    <div id="data-chart-container">
        <div v-if="this.stationData.length == 0" class="no-data-message">
            <p>No data yet.</p>
        </div>
        <div id="readings-container" class="section" v-if="this.station">
            <div id="readings-label">
                Latest Reading <span class="synced">Last synced {{ getSyncedDate() }}</span>
            </div>
            <div
                v-for="sensor in this.allSensors"
                v-bind:key="sensor.id"
                :class="'reading' + (selectedSensor && selectedSensor.id == sensor.id ? ' active' : '')"
                :data-id="sensor.id"
                v-on:click="switchSensor"
            >
                <div class="left">{{ sensor.name }}</div>
                <div class="right">
                    <span class="reading-value">{{ sensor.current_reading }} </span>
                    <span class="reading-unit">{{ sensor.unit ? sensor.unit : "" }}</span>
                </div>
            </div>
        </div>
        <div v-if="this.station">
            <div id="selected-sensor-graph"></div>
            <!-- <img alt="Data viz image" src="../assets/viz_placeholder.jpg" /> -->
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
.reading {
    font-size: 12px;
    width: 200px;
    float: left;
    line-height: 20px;
    padding: 18px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
}
.reading.active {
    border-bottom: 3px solid #1b80c9;
    padding-bottom: 15px;
}
.reading-value {
    font-size: 16px;
}
.reading-unit {
    font-size: 11px;
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
#selected-sensor-graph {
    width: 1110px;
    height: 500px;
    background-color: #fbeef0;
}
</style>
