<template>
    <div class="explore-view">
        <SidebarNav :isAuthenticated="isAuthenticated" :stations="stations" :projects="userProjects" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" />
        <VizWorkspace v-if="workspace" :workspace="workspace"></VizWorkspace>
    </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import HeaderBar from "../../components/HeaderBar";
import SidebarNav from "../../components/SidebarNav";

import FKApi from "@/api/api";

import { Workspace } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

export default {
    name: "ExploreView",
    components: {
        HeaderBar,
        SidebarNav,
        VizWorkspace,
    },
    props: {},
    data: () => {
        return {
            workspace: null,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s) => s.user.user,
            stations: (s) => s.stations.stations.user,
            userProjects: (s) => s.stations.projects.user,
        }),
    },
    beforeMount() {
        this.$store.watch(
            (state, getters) => state.stations.stations.all[12],
            (station, _old) => {
                return new FKApi().getAllSensors().then((sensors) => {
                    this.workspace = new Workspace(sensors);
                    this.workspace.addStation(station);
                    console.log("workspace", this.workspace);
                });
            }
        );

        return Promise.all([this.$store.dispatch(ActionTypes.NEED_PROJECTS), this.$store.dispatch(ActionTypes.NEED_STATIONS)]);
    },
    methods: {
        showStation(station) {
            console.log("station");
        },
    },
};
</script>

<style>
.explore-view {
    height: 100%;
    width: auto;
    text-align: left;
}
.workspace {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: stretch;
}
.tree-container {
    flex: 0;
}
.groups-container {
    flex: 1;
}
</style>
