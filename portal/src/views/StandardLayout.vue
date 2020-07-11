<template>
    <div class="container-ignored">
        <div class="container-top">
            <SidebarNav
                viewingStations="true"
                :isAuthenticated="isAuthenticated"
                :stations="stations"
                :projects="userProjects"
                :narrow="sidebar.narrow"
                @show-station="showStation"
            />

            <div class="container-main">
                <div class="container-header">
                    <HeaderBar @toggled="onSidebarToggle" />
                </div>

                <slot></slot>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";
import HeaderBar from "../components/HeaderBar.vue";
import SidebarNav from "../components/SidebarNav.vue";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "StandardLayout",
    components: {
        HeaderBar,
        SidebarNav,
    },
    data: () => {
        return {
            sidebar: {
                narrow: false,
            },
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy", mapped: "mapped" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            hasNoStations: (s: GlobalState) => s.stations.hasNoStations,
            stations: (s: GlobalState) => s.stations.stations.user,
            userProjects: (s: GlobalState) => s.stations.projects.user,
            anyStations: (s: GlobalState) => s.stations.stations.user.length > 0,
        }),
    },
    beforeMount() {
        console.log("StandardLayout: beforeMount");
        // this.$store.dispatch(ActionTypes.NEED_COMMON);
    },
    methods: {
        onSidebarToggle(...args) {
            console.log("sidebar-toggle", args);
            this.sidebar.narrow = !this.sidebar.narrow;
            this.$emit("sidebar-toggle", this.sidebar.narrow);
        },
        showStation(station, ...args) {
            console.log("show-station", station, args);
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
    },
});
</script>

<style scoped>
.container-ignored {
    height: 100%;
}
.container-top {
    display: flex;
    flex-direction: row;
    height: 100%;
}
.container-main {
    flex-grow: 1;
    flex-direction: column;
    display: flex;
}
.container-header {
    height: 70px;
}

/*
.container-map {
    flex-grow: 1;
}
*/
</style>
