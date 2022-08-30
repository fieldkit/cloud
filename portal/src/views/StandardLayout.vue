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
                :narrow="sidebarNarrow"
                :clipStations="clipStations"
                @show-station="showStation"
                @sidebar-toggle="onSidebarToggle"
            />

            <div class="container-main">
                <HeaderBar />

                <slot></slot>

                <template v-if="isPartnerCustomisationEnabled()">
                    <div class="app-logo">
                        <span>{{ $t("layout.side.madeBy") }}</span>
                        <br />
                        <a href="https://www.fieldkit.org" target="_blank">
                            <i role="img" class="icon icon-logo-fieldkit" :aria-label="$tc('layout.logo.fieldkit')"></i>
                        </a>
                    </div>
                </template>
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
import { GlobalState, DisplayStation } from "@/store";
import { getPartnerCustomizationWithDefault, isCustomisationEnabled } from "@/views/shared/partners";

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
        sidebarNarrow: {
            type: Boolean,
            default: getPartnerCustomizationWithDefault().sidebarNarrow,
        },
        clipStations: {
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
                this.$router.push({ name: "mapStation", params: { id: String(station.id) } });
            }
            return Promise.resolve();
        },
        onSidebarToggle(): void {
            this.$emit("sidebar-toggle");
        },
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/variables";
@import "src/scss/mixins";

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
    width: 100%;
    min-height: 100vh;

    body.floodnet & {
        background: #{lighten($color-floodnet-border, 5%)};
    }
}
.scrolling-disabled {
    overflow-y: hidden;
}

.app-logo {
    font-size: 11px;
    padding-top: 10px;
    position: fixed;
    bottom: 8px;
    left: 6px;
    text-align: left;
    display: none;

    @include bp-down($sm) {
        position: initial;
        padding: 30px;
        display: block;
    }

    @include bp-down($xs) {
        padding: 10px;
    }

    span {
        font-family: var(--font-family-bold);
        margin-bottom: 5px;
        font-size: 10px;
        min-width: 50px;
    }

    i:before {
        color: var(--color-dark);
    }
}
</style>
