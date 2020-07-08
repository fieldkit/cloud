<template>
    <div class="container-ignored">
        <div class="container-top">
            <div class="container-side">
                <SidebarNav
                    viewingStations="true"
                    :isAuthenticated="isAuthenticated"
                    :stations="stations"
                    :projects="userProjects"
                    @show-station="showStation"
                />
            </div>

            <div class="container-main">
                <div class="container-header">
                    <HeaderBar @toggled="onSidebarToggle" />
                </div>

                <slot></slot>
            </div>
        </div>
    </div>
</template>

<script>
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
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
            user: (s) => s.user.user,
            hasNoStations: (s) => s.stations.hasNoStations,
            stations: (s) => s.stations.stations.user,
            userProjects: (s) => s.stations.projects.user,
            anyStations: (s) => s.stations.stations.user.length > 0,
        }),
    },
    beforeMount() {
        console.log("StandardLayout: beforeMount");
        this.$store.dispatch(ActionTypes.NEED_COMMON);
    },
    methods: {
        onSidebarToggle(...args) {
            console.log("sidebar-toggle", args);
        },
        showStation(station, ...args) {
            console.log("show-station", station, args);
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
    },
};
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
.container-side {
    width: 240px;
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
