<template>
    <StandardLayout class="explore-view">
        <VizWorkspace v-if="workspace" :workspace="workspace"></VizWorkspace>
    </StandardLayout>
</template>

<script>
import StandardLayout from "../StandardLayout";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

import FKApi from "@/api/api";

import { Workspace } from "./viz";
import { VizWorkspace } from "./VizWorkspace";

export default {
    name: "ExploreView",
    components: {
        StandardLayout,
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
                    return this.workspace.compare().compare().combine();
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
