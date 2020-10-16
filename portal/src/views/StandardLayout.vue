<template>
    <div class="container-ignored" v-bind:class="{ 'scrolling-disabled': disableScrolling }">
        <div class="container-top">
            <SidebarNav
                :viewingStations="viewingStations"
                :viewingStation="viewingStation"
                :viewingProjects="viewingProjects"
                :viewingProject="viewingProject"
                :isAuthenticated="isAuthenticated"
                :stations="stations"
                :projects="userProjects"
                @show-station="showStation"
                @sidebar-toggle="onSidebarToggle"
            />

            <div class="container-main">
                <HeaderBar />

                <slot></slot>
            </div>
        </div>
        <Zoho />
    </div>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";
import HeaderBar from "./shared/HeaderBar.vue";
import SidebarNav from "./shared/SidebarNav.vue";
import Zoho from "./shared/Zoho.vue";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState, DisplayStation } from "@/store";

export default Vue.extend({
    name: "StandardLayout",
    components: {
        HeaderBar,
        SidebarNav,
        Zoho,
    },
    props: {
        viewingProjects: {
            type: Boolean,
            default: false,
        },
        viewingProject: {
            type: Object,
            default: null,
        },
        viewingStations: {
            type: Boolean,
            default: false,
        },
        viewingStation: {
            type: Object,
            default: null,
        },
        defaultShowStation: {
            type: Boolean,
            default: true,
        },
        disableScrolling: {
            type: Boolean,
            default: false,
        },
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy", mapped: "mapped" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            hasNoStations: (s: GlobalState) => s.stations.hasNoStations,
            stations: (s: GlobalState) => Object.values(s.stations.user.stations),
            userProjects: (s: GlobalState) => Object.values(s.stations.user.projects),
            anyStations: (s: GlobalState) => Object.values(s.stations.user.stations).length > 0,
        }),
    },
    beforeMount(): void {
        console.log("StandardLayout: beforeMount");
    },
    methods: {
        showStation(station: DisplayStation): Promise<any> {
            this.$emit("show-station", station.id);
            if (this.defaultShowStation) {
                // All parameters are strings.
                this.$router.push({ name: "viewStation", params: { id: String(station.id) } });
            }
            return Promise.resolve();
        },
        onSidebarToggle(): void {
            this.$emit("sidebar-toggle");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../scss/variables";
.container-ignored {
    height: 100%;
}
.container-top {
    display: flex;
    flex-direction: row;
    min-height: 100vh;
    position: relative;
}
.container-main {
    flex-grow: 1;
    flex-direction: column;
    display: flex;
    background: #fcfcfc;
    position: relative;
}
.scrolling-disabled {
    overflow-y: hidden;
}
</style>
