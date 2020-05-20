<template>
    <div class="readings-container">
        <div class="heading">Latest Readings</div>
        <div class="station-name">
            {{ currentStation.name }}
        </div>
        <div v-if="!currentStation.modules || currentStation.modules.length == 0" class="no-readings">
            No readings uploaded yet
        </div>
        <div class="sensors-container">
            <div v-for="(module, moduleIndex) in currentStation.modules" v-bind:key="module.id">
                <template v-if="module.position < 5">
                    <div
                        v-for="(sensor, sensorIndex) in module.sensors"
                        v-bind:key="sensor.id"
                        :class="getCounter(moduleIndex, sensorIndex) % 2 == 1 ? 'left-reading' : 'right-reading'"
                    >
                        <div class="left sensor-name">{{ getSensorName(module, sensor) }}</div>
                        <div class="right sensor-unit">
                            {{ sensor.unit_of_measure }}
                        </div>
                        <div class="right sensor-reading">
                            {{ getReading(sensor) }}
                        </div>
                    </div>
                </template>
            </div>
        </div>

        <Pagination
            class="pagination-section"
            :currentPage="currentIndex"
            :pageCount="pageCount"
            :numVisiblePages="stationsPerPage"
            @nextPage="onPageChange('next')"
            @previousPage="onPageChange('previous')"
            @loadPage="onPageChange"
        />
    </div>
</template>

<script>
import FKApi from "../api/api";
import * as utils from "../utilities";
import Pagination from "../components/Pagination";

export default {
    name: "StationsReadings",
    components: {
        Pagination,
    },
    data: () => {
        return {
            stations: [],
            currentStation: {},
            currentIndex: 1,
            pageCount: 0,
            stationsPerPage: 1,
            moduleSensorCounter: 0,
            modulesSensors: {},
        };
    },
    props: ["project"],
    async beforeCreate() {
        this.api = new FKApi();
    },
    mounted() {
        this.fetchStations();
    },
    methods: {
        onPageChange(value) {
            switch (value) {
                case "next":
                    this.currentIndex += 1;
                    break;
                case "previous":
                    this.currentIndex -= 1;
                    break;
                default:
                    this.currentIndex = value;
            }
            this.currentStation = this.stations[this.currentIndex - 1];
            if (!this.currentStation) {
                this.currentStation = {};
            }
        },

        fetchStations() {
            this.api.getStationsByProject(this.project.id).then(result => {
                this.stations = result.stations;
                this.pageCount = this.stations.length;
                this.currentStation = this.stations[this.currentIndex - 1];
                if (!this.currentStation) {
                    this.currentStation = {};
                }
            });
        },

        getCounter(moduleIndex, sensorIndex) {
            if (this.modulesSensors[moduleIndex]) {
                if (!this.modulesSensors[moduleIndex][sensorIndex]) {
                    this.moduleSensorCounter += 1;
                    this.modulesSensors[moduleIndex][sensorIndex] = this.moduleSensorCounter;
                }
            } else {
                this.moduleSensorCounter += 1;
                this.modulesSensors[moduleIndex] = {};
                this.modulesSensors[moduleIndex][sensorIndex] = this.moduleSensorCounter;
            }
            return this.modulesSensors[moduleIndex][sensorIndex];
        },

        getReading(sensor) {
            if (!sensor.reading) {
                return "--";
            }
            return sensor.reading.last || sensor.reading.last == 0 ? sensor.reading.last.toFixed(1) : "--";
        },

        getSensorName(module, sensor) {
            const newName = utils.convertOldFirmwareResponse(module);
            return this.$t(newName + ".sensors." + sensor.name);
            // return this.$t(module.name + "." + sensor.name);
        },
    },
};
</script>

<style scoped>
.readings-container {
    position: relative;
    float: left;
    width: 415px;
    height: 300px;
    margin: 22px 0 0 0;
    border: 1px solid #d8dce0;
}
.heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 25px 25px 16px 25px;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    float: left;
    clear: both;
    margin: 0 0 12px 25px;
}
.sensors-container {
    width: 100%;
    height: 150px;
    overflow: scroll;
}
.left {
    float: left;
}
.right {
    float: right;
}
.left-reading,
.right-reading {
    width: 42%;
    margin: 8px 0;
    padding: 5px 10px;
    background-color: #f4f5f7;
    border-radius: 2px;
}
.left-reading {
    float: left;
    clear: both;
    margin-left: 10px;
}
.right-reading {
    float: right;
    margin-right: 3px;
}
.sensor-name {
    font-size: 12px;
    font-weight: 600;
    display: inline-block;
}
.sensor-reading {
    font-size: 16px;
    display: inline-block;
}
.sensor-unit {
    font-size: 10px;
    margin: 6px 0 0 4px;
    display: inline-block;
}
.no-readings {
    float: left;
    clear: both;
    margin: 25px;
}
.pagination-section {
    width: 100%;
    position: absolute;
    bottom: 25px;
    text-align: center;
}
</style>
