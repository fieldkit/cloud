<template>
    <div id="data-example-container">
        Hello
    </div>
</template>

<script>
import _ from "lodash";
import FKApi from "../api/api";

export default {
    name: "DataExample",
    data: () => {
        return {
            data: []
        };
    },
    props: {
        stationId: Number
    },
    async beforeCreate() {
        const api = new FKApi();
        const stations = await api.getStations();
        const id = _(stations.stations).first().device_id;
        const summary = await api.getStationDataSummaryByDeviceId(id);
        const data = await api.getJSONDataByDeviceId(id, 0, 1000);
        console.log(summary);
        console.log(data);
    }
};
</script>

<style scoped>
#data-example-container {
    width: 700px;
    margin: 60px 0;
}
</style>
