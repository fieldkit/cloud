<template>
    <div class="container-ignored" v-bind:class="{ 'scrolling-disabled': disableScrolling }">
        <div class="container-top">
            <SidebarNav
                :viewingStations="viewingStations"
                :viewingProjects="viewingProjects"
                :viewingProject="viewingProject"
                :isAuthenticated="isAuthenticated"
                :stations="stations"
                :projects="userProjects"
                :narrow="sidebar.narrow"
                @show-station="showStation"
                @toggle-menu="onSidebarToggle"
            />

            <div class="container-main">
                <HeaderBar @toggled="onSidebarToggle" :isMenuNarrow="sidebar.narrow" />

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
import { GlobalState } from "@/store/modules/global";

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
        defaultShowStation: {
            type: Boolean,
            default: true,
        },
        disableScrolling: {
            type: Boolean,
            default: false,
        },
    },
    data: () => {
        return {
            sidebar: {
                narrow: window.screen.availWidth > 1040 ? false : true,
            },
        };
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
    beforeMount() {
        console.log("StandardLayout: beforeMount");
    },
    mounted() {
        const desktopBreakpoint = 1040;

        const windowAny: any = window;
        const resizeObserver = new windowAny.ResizeObserver((entries) => {
            if (entries[0].contentRect.width < desktopBreakpoint) {
                if (!this.sidebar.narrow) {
                    this.onSidebarToggle();
                }
            }
        });
        resizeObserver.observe(document.querySelector("body"));
    },
    methods: {
        onSidebarToggle(...args) {
            this.sidebar.narrow = !this.sidebar.narrow;
            this.$emit("sidebar-toggle", this.sidebar.narrow);
        },
        showStation(station, ...args) {
            this.$emit("show-station", station.id);
            if (this.defaultShowStation) {
                this.$router.push({ name: "viewStation", params: { id: station.id } });
            }
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
}
.container-main {
    flex-grow: 1;
    flex-direction: column;
    display: flex;
    background: #fcfcfc;
}
.scrolling-disabled {
    overflow-y: hidden;
}
</style>
