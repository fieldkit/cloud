<template>
    <div class="readings-container">
        <div class="heading">Latest Readings</div>

        <template v-if="!stations.length">
            <div class="empty-project">
                There are no stations in this project.
            </div>
        </template>
        <template v-else>
            <div class="station-name">
                {{ visibleStation.name }}
            </div>
            <div class="sensors-container">
                <LatestStationReadings :id="visibleStation.id" />
            </div>
            <div class="pagination">
                <PaginationControls :page="index" :totalPages="stations.length" @new-page="onNewStation" />
            </div>
        </template>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

import LatestStationReadings from "@/views/shared/LatestStationReadings.vue";
import PaginationControls from "@/views/shared/PaginationControls.vue";

export default Vue.extend({
    name: "StationsReadings",
    components: {
        LatestStationReadings,
        PaginationControls,
    },
    data: () => {
        return {
            index: 0,
        };
    },
    props: {
        project: {
            type: Object,
            required: true,
        },
        displayProject: {
            type: Object,
            required: true,
        },
    },
    computed: {
        stations() {
            return this.displayProject.stations;
        },
        visibleStation() {
            if (this.displayProject.stations.length > 0) {
                return this.displayProject.stations[this.index];
            }
            return null;
        },
    },
    mounted(this: any) {
        // console.log("display", this.displayProject);
        // console.log("project", this.project);
    },
    methods: {
        onNewStation(this: any, index: number) {
            this.index = index;
        },
    },
});
</script>

<style scoped>
.readings-container {
    display: flex;
    flex-direction: column;
    min-height: 300px;
}
.heading {
    font-size: 20px;
    font-weight: 600;
}
.station-name {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 10px;
    margin-top: 10px;
}
.sensors-container {
}

.pagination {
    margin-top: auto;
}
</style>
