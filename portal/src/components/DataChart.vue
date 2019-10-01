<template>
    <div id="data-chart-container">
        <div v-if="this.stationData.length == 0" class="no-data-message">
            <p>No data yet.</p>
        </div>
        <div class="no-data-message">
            <p>{{ this.message }}</p>
        </div>
        <div v-if="this.station">
            <img alt="Data viz image" src="../assets/viz_placeholder.jpg" />
        </div>
    </div>
</template>

<script>
export default {
    name: "DataChart",
    data: () => {
        return {
            message: ""
        };
    },
    props: ["summary", "stationData", "station"],
    watch: {
        summary: function(newVal, oldVal) {
            // if(newVal.provisions.length == 0) {
            //     this.message = "No data files uploaded for " + this.station.name + " yet.";
            // }
            console.log("summary", newVal.provisions, "was", oldVal);
        },
        stationData: function(newVal, oldVal) {
            console.log("stationData", newVal, "was", oldVal);
        }
    }
};
</script>

<style scoped>
#data-chart-container {
    float: left;
    margin: 10px 0;
}
.no-data-message {
    font-size: 20px;
}
</style>
