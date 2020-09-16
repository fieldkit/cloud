<template>
    <div class="sps" v-bind:class="{ selected: selected }">
        <div class="standard" v-on:click="onClicked">
            <div class="name">{{ station.name }}</div>
            <div class="status">{{ status }}</div>
            <div class="location" v-if="!narrow">{{ location }}</div>
            <div class="seen" v-if="!narrow">{{ station.updatedAt | prettyDate }}</div>
        </div>
        <slot></slot>
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
            default: false,
        },
        narrow: {
            type: Boolean,
            default: false,
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

<style scoped lang="scss">
@import '../../scss/mixins';

.sps.selected {
    border: 2px solid #1b80c9;
}
.sps {
    display: flex;
    min-width: 260px;
    border: 1px solid #d8dce0;
    border-radius: 1px;

    @include bp-down($xs) {
        min-width: 220px;
    }
}
.sps .standard {
    text-align: left;
    display: flex;
    flex-direction: column;
    padding: 10px;
    width: 100%;
}
.sps .name {
    color: #2c3e50;
    font-size: 14px;
    font-weight: 500;
}
.sps .location {
    font-size: 12px;
    color: #6a6d71;
}
.sps .status {
    font-size: 12px;
    color: #6a6d71;
}
.sps .seen {
    font-size: 12px;
    color: #6a6d71;
}
</style>
