<template>
    <div class="datafiles-container">
        <div class="heading">Data Files</div>
        <div class="readings-stats-container">
            No readings have been uploaded yet.
        </div>
    </div>
</template>

<script>
import _ from "lodash";

export default {
    name: "ProjectDataFiles",
    data: () => {
        return {};
    },
    props: { projectStations: { required: true } },
};
</script>

<style scoped>
.datafiles-container {
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
}
.stat {
    font-weight: 600;
}
</style>
