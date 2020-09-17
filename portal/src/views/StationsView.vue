<template>
    <StandardLayout @sidebar-toggle="onSidebarToggle" :viewingStations="true">
        <div class="view-type-container">
            <div class="view-type">
                <div class="view-type-map" v-bind:class="{ active: viewType === 'map' }" v-on:click="switchView('map')"></div>
                <div class="view-type-list" v-bind:class="{ active: viewType === 'list' }" v-on:click="switchView('list')"></div>
            </div>
        </div>

        <template v-if="viewType === 'list'">
            <div class="stations-list">
                <StationSummary
                    v-for="station in stations"
                    v-bind:key="station.id"
                    class="summary-container"
                    @close="closeSummary"
                    :station="station"
                />
            </div>
        </template>

        <template v-if="viewType === 'map'">
            <div class="container-map">
                <StationsMap @show-summary="showSummary" :mapped="mapped" :layoutChanges="layoutChanges" v-if="mapped" />
                <StationSummary
                    v-if="activeStation"
                    class="summary-container"
                    @close="closeSummary"
                    :station="activeStation"
                    v-bind:key="activeStation.id"
                />
                <div v-if="isAuthenticated && showNoStationsMessage && hasNoStations" id="no-stations">
                    <div id="close-notice-btn" v-on:click="closeNotice">
                        <img alt="Close" src="@/assets/close.png" />
                    </div>
                    <p class="heading">Add a New Station</p>
                    <p class="text">
                        You currently don't have any stations on your account. Download the FieldKit app and connect your station to add
                        them to your account.
                    </p>
                    <a href="https://apps.apple.com/us/app/fieldkit-org/id1463631293?ls=1" target="_blank">
                        <img alt="App store" src="@/assets/appstore.png" class="app-btn" />
                    </a>
                    <a href="https://play.google.com/store/apps/details?id=com.fieldkit&hl=en_US" target="_blank">
                        <img alt="Google Play" src="@/assets/googleplay.png" class="app-btn" />
                    </a>
                </div>
            </div>
        </template>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "./StandardLayout.vue";
import StationSummary from "./shared/StationSummary.vue";
import StationsMap from "./shared/StationsMap.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "StationsView",
    components: {
        StandardLayout,
        StationsMap,
        StationSummary,
    },
    props: {
        id: { type: Number },
    },
    data: () => {
        return {
            layoutChanges: 0,
            showNoStationsMessage: true,
            viewType: "list",
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
        activeStation() {
            return this.$store.state.stations.stations[this.id];
        },
    },
    beforeMount() {
        if (this.id) {
            return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.id });
        }
    },
    watch: {
        id() {
            if (this.id) {
                return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.id });
            }
        },
    },
    methods: {
        goBack() {
            if (window.history.length) {
                return this.$router.go(-1);
            } else {
                return this.$router.push("/");
            }
        },
        showSummary(params: { id: number }) {
            if (this.id != params.id) {
                return this.$router.push({ name: "viewStation", params: params as any });
            }
        },
        closeSummary() {
            return this.$router.push({ name: "stations" });
        },
        onSidebarToggle() {
            this.layoutChanges++;
        },
        closeNotice() {
            this.showNoStationsMessage = false;
        },
        switchView(type) {
            this.viewType = type;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../scss/mixins.scss";

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

.container-map {
    flex-grow: 1;
}

::v-deep .station-hover-summary {
    left: 360px;
    top: 120px;
    width: 415px;
    border-radius: 3px;
}

#stations-view-panel {
    margin: 0;
}

.closeSummary {
    position: absolute;
}

#no-user {
    font-size: 20px;
    background-color: #ffffff;
    width: 400px;
    position: absolute;
    top: 40px;
    left: 260px;
    padding: 0 15px 15px 15px;
    margin: 60px;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
}

#no-stations {
    background-color: #ffffff;
    width: 360px;
    position: absolute;
    top: 70px;
    left: 50%;
    padding: 75px 60px 75px 60px;
    margin: 135px 0 0 -120px;
    text-align: center;
    border: 1px solid rgb(215, 220, 225);
    z-index: 2;
}

#no-stations .heading {
    font-size: 18px;
    font-weight: bold;
}

#no-stations .text {
    font-size: 14px;
}

#no-stations .app-btn {
    margin: 20px 14px;
}

#close-notice-btn {
    position: absolute;
    top: 10px;
    right: 10px;
    cursor: pointer;
}

.show-link {
    text-decoration: underline;
}

.view-type {
    width: 100px;
    height: 39px;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
    border: solid 1px #f4f5f7;
    background-color: #ffffff;
    cursor: pointer;
    margin-left: auto;
    @include flex();

    &-container {
        z-index: $z-index-top;
        width: 100%;
        max-width: 1440px;
        @include position(absolute, 100px null null 0);

        @include bp-down($xs) {
            @include position(absolute, 83px null null 0);
        }
    }

    &-list {
        background: url("../assets/icon-list.svg") no-repeat center center;
        flex-basis: 50%;

        &.active {
            background: url("../assets/icon-list-selected.svg") no-repeat center center;
        }
    }

    &-map {
        background: url("../assets/icon-map.svg") no-repeat center center;
        flex-basis: 50%;

        &.active {
            background: url("../assets/icon-map-selected.svg") no-repeat center center;
        }
    }
}

::v-deep .stations-list {
    @include flex();
    flex-wrap: wrap;
    padding: 100px 70px;
    margin: -20px;
    max-width: 1080px;

    @include bp-down($md) {
        padding: 100px 30px;
    }

    @include bp-down($sm) {
        justify-content: center;
    }

    @include bp-down($xs) {
        padding: 80px 10px;
        margin: -5px 0;
    }

    .summary-container {
        z-index: 0;
        position: unset;
        margin: 20px;
        flex-basis: 389px;
        box-sizing: border-box;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);

        @include bp-down($md) {
            padding: 19px 11px;
            flex-basis: calc(50% - 40px);
        }

        @include bp-down($sm) {
            justify-self: center;
            max-width: 400px;
            flex-basis: 100%;
            margin: 10px 0;
        }

        @include bp-down($xs) {
            margin: 5px 0;
        }

        .close-button {
            display: none;
        }
    }
}
</style>
