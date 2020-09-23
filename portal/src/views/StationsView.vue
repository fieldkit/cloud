<template>
    <StandardLayout :viewingStations="true" :viewingStation="activeStation">
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
                <StationsMap @show-summary="showSummary" :mapped="mapped" v-if="mapped" />
            </div>
            <StationSummary
                    v-if="activeStation"
                    class="summary-container"
                    @close="closeSummary"
                    :station="activeStation"
                    v-bind:key="activeStation.id"
            />
            <div id="no-stations" v-if="isAuthenticated && showNoStationsMessage && hasNoStations">
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
            showNoStationsMessage: true,
            viewType: "map",
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

.container-header {
    height: 70px;
}

.container-map {
    position: absolute;
    width: 100%;
    height: calc(100% - 67px);
    left: 0;
    top: 67px;
}

::v-deep .station-hover-summary {
    left: 360px;
    top: 120px;
    border-radius: 3px;
}

::v-deep .summary-container {
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
    position: unset;
    margin: 120px 0 60px 119px;
    max-width: calc(100vw - 20px);

    .explore-button {
        margin-bottom: 6px;
    }

    @include bp-down($sm) {
        margin: 120px auto 60px auto;
    }
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
    z-index: $z-index-top;
}

#no-stations {
    background-color: #ffffff;
    width: 360px;
    padding: 85px 80px 95px 80px;
    margin: 120px auto 60px auto;
    text-align: center;
    border: 1px solid rgb(215, 220, 225);
    z-index: $z-index-top;

    a {
        display: inline-block;
    }

    .heading {
        font-size: 18px;
        font-weight: 900;
        margin-bottom: 2px;
        margin-top: 0;
    }
}


#no-stations .text {
    font-size: 14px;
    margin: 0;
}

#no-stations .app-btn {
    margin: 35px 14px 0;
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
    @include flex();

    &-container {
        z-index: $z-index-top;
        @include position(absolute, 90px 25px null null);

        @include bp-down($xs) {
            @include position(absolute, 80px 10px null null);
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
        border-right: solid 1px #f4f5f7;

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
    width: calc(100% - 40px);
    box-sizing: border-box;

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
