<template>
    <div class="sps" v-on:click="onClicked" v-bind:class="{ selected: selected }">
        <div class="name">{{ station.name }}</div>
        <div class="location">{{ location }}</div>
        <div class="status">{{ status }}</div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";

export default Vue.extend({
    name: "StationPickerStation",
    props: {
        station: {
            type: Object,
            required: true,
        },
        selected: {
            type: Boolean,
            required: true,
        },
    },
    computed: {
        location(this: any) {
            const gps = [this.station.location?.latitude, this.station.location?.longitude]
                .filter((v) => v)
                .map((v) => v.toFixed(3))
                .join(", ");
            const locations = [gps, "Unknown"];
            return _(locations).compact().first();
        },
        status(this: any) {
            if (this.station.deployedAt) {
                return "Deployed";
            }
            return "Ready to Deploy";
        },
    },
    methods: {
        onClicked() {
            this.$emit("selected");
        },
    },
});
</script>

<style scoped>
.sps {
    display: flex;
    flex-direction: column;
    min-width: 260px;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    padding: 10px;
    text-align: left;
}
.sps.selected {
    border: 2px solid #1b80c9;
}
.sps .name {
    color: #2c3e50;
    font-size: 16px;
    font-weight: 500;
}
.sps .location {
    font-size: 14px;
    color: #2c3e50;
}
.sps .status {
    font-size: 12px;
    color: #6a6d71;
}
</style>
