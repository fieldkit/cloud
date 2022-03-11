<template>
    <StandardLayout v-if="station">
        <div class="container-wrap notes-view">
            <DoubleHeader
                :title="station.name"
                :subtitle="headerSubtitle"
                :backTitle="$tc('layout.backProjectDashboard')"
                backRoute="viewProject"
                :backRouteParams="{ id: station.id }"
            />
            <section class="flex station-section">
                <div class="container-box">
                    <div class="flex flex-al-center">
                        <StationPhoto :station="station" />
                        <div>
                            <div class="station-name">{{ station.name }}</div>
                            <div>
                                {{ $tc("station.lastSynced") }}
                                <span class="small-light">{{ station.updatedAt | prettyDate }}</span>
                            </div>
                            <div class="station-battery" v-if="station.battery">
                                <img class="battery" alt="Battery Level" :src="getBatteryIcon()" />
                                <span class="small-light">{{ station.battery.toFixed() }} %</span>
                            </div>
                        </div>
                    </div>

                    <div>
                        <div v-if="station.location" class="station-row flex flex-space-between">
                            <div class="flex flex-al-center station-location">
                                <i class="icon icon-location"></i>
                                <!--
                                <span>{{ station.location }}</span>
-->
                                <span>Los Angeles</span>
                            </div>
                            <div class="flex">
                                <div class="station-coordinate">
                                    <div>{{ station.location.latitude | prettyCoordinate }}</div>
                                    <div>{{ $tc("station.latitude") }}</div>
                                </div>
                                <div class="station-coordinate">
                                    <div>{{ station.location.longitude | prettyCoordinate }}</div>
                                    <div>{{ $tc("station.longitude") }}</div>
                                </div>
                            </div>
                        </div>

                        <div v-if="station.placeNameNative" class="station-row">
                            <i class="icon icon-location"></i>
                            <span>{{ $tc("station.nativeLand") }} {{ station.placeNameNative }}</span>
                        </div>

                        <div class="station-row">
                            <span>{{ $tc("station.firmwareVersion") }}</span>
                            <span class="ml-10">2.00.1</span>
                        </div>

                        <div class="station-row">
                            <span>{{ $tc("station.modules") }}</span>
                            <div class="station-modules ml-10">
                                <img
                                    v-for="module in station.modules"
                                    v-bind:key="module.name"
                                    alt="Module icon"
                                    :src="getModuleImg(module)"
                                />
                            </div>
                        </div>
                    </div>
                </div>
                <div>
                    <!--                    <div class="photo-container" v-for="photo in station.photos" v-bind:key="photo.key">
                        <AuthenticatedPhoto :url="photo.url" />
                    </div>-->
                </div>
            </section>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "@/views/StandardLayout.vue";
import DoubleHeader from "@/views/shared/DoubleHeader.vue";
import StationPhoto from "@/views/shared/StationPhoto.vue";
// import AuthenticatedPhoto from "@/views/shared/AuthenticatedPhoto.vue";
import { DisplayStation, ProjectModule } from "@/store";
import * as utils from "@/utilities";

export default Vue.extend({
    name: "StationView",
    components: {
        StandardLayout,
        DoubleHeader,
        StationPhoto,
        //    AuthenticatedPhoto,
    },
    mounted() {
        this.getStation(parseInt(this.$route.params.id));
    },
    data(): {
        station: DisplayStation | null;
    } {
        return {
            station: null,
        };
    },
    computed: {
        headerSubtitle() {
            let subtitle;
            if (this.station && this.station.deployedAt && this.$options.filters?.prettyDate) {
                if (this.station.deployedAt) {
                    subtitle = this.$tc("station.deployed") + " " + this.$options.filters.prettyDate(this.station.deployedAt);
                } else {
                    subtitle = this.$tc("station.readyToDeploy");
                }
            }
            return subtitle;
        },
    },
    methods: {
        getStation(id: number) {
            this.$services.api.getStation(id).then((station) => {
                console.log("station from api", station);
                this.station = new DisplayStation(station);
                console.log("radoi", this.station);
                // Vue.set(this.notes, stationId, notes);
                // this.loading = false;
            });
        },
        getBatteryIcon(): string {
            if (!this.station || !this.station.battery) {
                return "";
            }
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
    },
});
</script>

<style scoped lang="scss">
@import "../scss/mixins";
@import "../scss/layout";

.container-box {
    border: 1px solid var(--color-border);
    border-radius: 1px;
    background-color: #fff;
    padding: 20px;
    font-size: 14px;
}

.station-section {
    margin-top: 30px;

    > div {
        flex-basis: calc(50% - 10px);
    }

    ::v-deep .station-photo {
        margin-right: 20px;
        object-fit: contain;
        width: 90px;
        height: 90px;
        border-radius: 5px;
    }
}

.station {
    &-name {
        font-size: 18px;
        font-weight: 900;
        margin-bottom: 2px;
    }
    &-battery {
        margin-top: 5px;
        @include flex(flex-start);

        span {
            margin-left: 5px;
        }
    }
    &-modules {
        margin-left: 18px;
        flex-wrap: wrap;
        @include flex;

        img {
            margin-right: 10px;
            width: 30px;
            height: 30px;
        }
    }
    &-coordinate {
        padding: 0 15px;

        &:last-of-type {
            padding-right: 0;
        }

        div:nth-of-type(2) {
            font-size: 12px;
            margin-top: 2px;
        }
    }
    &-location {
        align-self: flex-start;
    }
    &-row {
        padding: 15px 0;
        max-width: 350px;
        @include flex(center);

        &:not(:last-of-type) {
            border-bottom: solid 1px var(--color-border);
        }

        .icon {
            margin-right: 7px;
            align-self: flex-start;
        }
    }
}

.small-light {
    font-size: 12px;
    color: #6a6d71;
}
</style>
