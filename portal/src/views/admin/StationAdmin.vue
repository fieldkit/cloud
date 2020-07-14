<template>
    <StandardLayout>
        <div class="container">
            <LatestStationReadings :id="169" />
            <LatestStationReadings :id="159" />
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "StationAdmin",
    components: {
        StandardLayout,
        ...CommonComponents,
    },
    props: {
        id: {
            type: Number,
            required: true,
        },
    },
    data: () => {
        return {};
    },
    methods: {},
});
</script>

<style scoped>
.container {
    width: 400px;
}
</style>
