<template>
    <div class="datafiles-container">
        <div class="heading">Data Files</div>
        <div class="readings-stats-container" v-if="projectStations.length > 0">
            A total of
            <span class="stat">{{ totalReadings.toLocaleString() }}</span>
            readings from
            <span class="stat">{{ projectStations.length }}</span>
            stations have been uploaded between
            <span class="stat">{{ first.toLocaleDateString() }}</span>
            and
            <span class="stat">{{ last.toLocaleDateString() }}.</span>
        </div>
        <div class="readings-stats-container" v-else>
            No readings have been uploaded yet.
        </div>
    </div>
</template>

<script>
import _ from "lodash";

export default {
    name: "ProjectDataFiles",
    data: () => {
        return {
            totalReadings: 0,
            first: new Date(),
            last: new Date("1/1/1970"),
        };
    },
    props: ["projectStations"],
    watch: {
        projectStations() {
            if (this.projectStations) {
                this.init();
            }
        },
    },
    methods: {
        init() {
            this.totalReadings = 0;
            this.projectStations.forEach(s => {
                const highest = _.maxBy(s.uploads, u => {
                    return u.blocks[1];
                });
                this.totalReadings += highest.blocks[1];
                const earliest = _.minBy(s.uploads, u => {
                    return u.time;
                });
                const latest = _.maxBy(s.uploads, u => {
                    return u.time;
                });
                if (earliest.time < this.first.getTime()) {
                    this.first = new Date(earliest.time);
                }
                if (latest.time > this.last.getTime()) {
                    this.last = new Date(latest.time);
                }
            });
        },
    },
};
</script>

<style scoped>
.datafiles-container {
    position: relative;
    float: left;
    width: 582px;
    height: 300px;
    margin: 22px 22px 0 0;
    border: 1px solid #d8dce0;
}
.heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 25px;
}
.readings-stats-container {
    float: left;
    clear: both;
    margin: 0 25px 25px 25px;
    font-size: 20px;
}
.stat {
    font-weight: 600;
}
</style>
