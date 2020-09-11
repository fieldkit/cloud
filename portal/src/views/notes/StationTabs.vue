<template>
    <div class="station-tabs">
        <div
            class="tab"
            v-for="station in stations"
            v-bind:key="station.id"
            v-on:click="(ev) => clicked(ev, station)"
            v-bind:class="{ selected: selected && station.id == selected.id }"
        >
            <div class="name">
                {{ station.name }}
            </div>
            <div v-if="station.deployedAt" class="deployed">Deployed</div>
            <div v-else class="undeployed">Not Deployed</div>
        </div>
        <div class="vertical">&nbsp;</div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

import FKApi from "@/api/api";

export default Vue.extend({
    name: "StationTabs",
    components: {},
    props: {
        stations: {
            type: Array,
            required: true,
        },
        selected: {
            type: Object,
            required: false,
        },
    },
    data: () => {
        return {};
    },
    computed: {},
    mounted(this: any) {
        //
    },
    methods: {
        clicked(ev, station) {
            return this.$emit("selected", station);
        },
    },
});
</script>

<style scoped lang="scss">
.station-tabs {
    text-align: left;
    display: flex;
    flex-direction: column;
    height: 100%;
}
.tab {
    padding: 16px 13px;
    border-right: 1px solid #d8dce0;
    border-bottom: 1px solid #d8dce0;
    border-left: 4px solid white;
    cursor: pointer;

    &.selected {
        border-right: none;
    }
}

.vertical {
    margin-top: auto;
    border-right: 1px solid #d8dce0;
    height: 100%;
}
.name {
    font-size: 16px;
    font-weight: 500;
    color: #2c3e50;
    margin-bottom: 1px;
    font-weight: 500;
}
.undeployed,
.deployed {
    font-size: 13px;
    color: #6a6d71;
    font-weight: 500;
}
.selected {
    border-left: 4px solid #1b80c9;
}
</style>
