<template>
    <div class="datafiles-container">
        <div class="heading">Data Files</div>
        <div class="message">
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
    props: {
        projectStations: { required: true },
    },
};
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.datafiles-container {
}
.heading {
    font-size: 20px;
    font-weight: 500;
}
.message {
    margin-top: 10px;
    font-family: $font-family-light;
}
</style>
