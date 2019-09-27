<template>
    <div id="data-example-container">
        Hello
    </div>
</template>

<script>
export default {
    name: "DataExample",
    data: () => {
        return {};
    },
    props: ["summary", "stationData"],
    watch: {
        summary: function(newVal, oldVal) {
            console.log("summary", newVal, "was", oldVal);
        },
        stationData: function(newVal, oldVal) {
            console.log("stationData", newVal, "was", oldVal);
        }
    }
};
</script>

<style scoped>
#data-example-container {
    width: 700px;
    margin: 60px 0;
}
</style>
