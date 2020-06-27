<template>
    <div class="datafiles-container">
        <div class="heading">Data Files</div>
        <div class="readings-stats-container" v-if="statistics.totalReadings > 0">
            A total of
            <span class="stat">{{ statistics.totalReadings.toLocaleString() }}</span>
            readings from
            <span class="stat">{{ projectStations.length }}</span>
            stations have been uploaded between
            <span class="stat">{{ statistics.first.toLocaleDateString() }}</span>
            and
            <span class="stat">{{ statistics.last.toLocaleDateString() }}.</span>
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
    props: { projectStations: { required: true } },
    computed: {
        statistics() {
            const totalReadings = _(this.projectStations)
                .map(ps => ps.totalReadings)
                .sum();
            const first = _(this.projectStations)
                .map(ps => ps.receivedAt)
                .min();
            const last = _(this.projectStations)
                .map(ps => ps.receivedAt)
                .max();
            return {
                totalReadings: totalReadings,
                first: first,
                last: last,
            };
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
