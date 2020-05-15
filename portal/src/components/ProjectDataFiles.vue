<template>
    <div class="datafiles-container">
        <div class="heading">Data Files</div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "ProjectDataFiles",
    data: () => {
        return {};
    },
    props: ["project"],
    async beforeCreate() {
        this.api = new FKApi();
    },
    mounted() {},
    methods: {},
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
</style>
