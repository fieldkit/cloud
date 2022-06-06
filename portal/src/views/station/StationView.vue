<template>
    <div v-if="!isBusy">
        <StandardLayout v-if="station">
            <div class="container-wrap">
                <DoubleHeader
                    backRoute="viewProject"
                    :title="station.name"
                    :subtitle="headerSubtitle"
                    :backTitle="$tc('layout.backProjectDashboard')"
                    :backRouteParams="{ id: $route.params.projectId }"
                />

                <section class="section-station">
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
                            <div class="station-row flex flex-space-between">
                                <div v-if="station.locationName" class="flex flex-al-center station-location">
                                    <i class="icon icon-location"></i>
                                    <span>{{ station.locationName }}</span>
                                </div>
                                <div v-if="station.location" class="flex">
                                    <div class="station-coordinate" :style="station.location ? { 'padding-left': 0 } : ''">
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

                            <div v-if="station.firmwareNumber" class="station-row">
                                <span>{{ $tc("station.firmwareVersion") }}</span>
                                <span class="ml-10">{{ station.firmwareNumber }}</span>
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
                    <div v-if="photos" class="station-photos">
                        <div class="photo-container" v-for="(n, index) in 4" v-bind:key="index">
                            <!-- somehow using v-for like so needs the next v-if -->
                            <AuthenticatedPhoto v-if="photos[index]" :url="photos[index].url" />
                            <div v-else class="photo-placeholder">
                                <img src="@/assets/image-placeholder-v2.svg" alt="Image placeholder" />
                            </div>
                        </div>
                        <router-link :to="{ name: 'viewStationPhotos' }" class="station-photos-nav">
                            <i class="icon icon-grid"></i>
                            {{ $t("station.managePhotos") }}
                        </router-link>
                    </div>
                </section>

                <section class="container-box section-readings" v-if="selectedModule">
                    <div class="station-readings">
                        <ul>
                            <li
                                v-for="module in station.modules"
                                v-bind:key="module.name"
                                :class="{ active: module.name === selectedModule.name }"
                                @click="selectedModule = module"
                            >
                                <img v-bind:key="module.name" alt="Module icon" :src="getModuleImg(module)" />
                                <span>{{ $t(getModuleName(module)) }}</span>
                            </li>
                        </ul>
                        <header v-if="isMobileView">
                            <img alt="Module icon" :src="getModuleImg(selectedModule)" />
                            {{ $t(getModuleName(selectedModule)) }}
                        </header>
                        <div class="station-readings-values">
                            <header v-if="!isMobileView">{{ $t(getModuleName(selectedModule)) }}</header>
                            <LatestStationReadings :id="station.id" :moduleKey="getModuleName(selectedModule)" />
                        </div>
                    </div>
                </section>

                <section v-if="attributes.length > 0" class="section-notes container-box">
                    <ProjectAttributes :attributes="attributes" />
                </section>

                <section v-if="showMap">
                    <div class="container-map">
                        <StationsMap :mapped="mapped" :showStations="true" :mapBounds="mapped.bounds" />
                    </div>
                </section>

                <section v-if="notes" class="section-notes container-box">
                    <NotesForm v-bind:key="station.id" :station="station" :notes="{ notes, media }" :readonly="false" />
                </section>
            </div>
        </StandardLayout>
        <div v-else class="forbidden-view-bg">
            <ForbiddenBanner :title="$t('unauthorized')" :subtitle="$t('unauthorizedStation')"></ForbiddenBanner>
            <router-link to="/dashboard">
                {{ $t("layout.backProjects") }}
            </router-link>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "@/views/StandardLayout.vue";
import DoubleHeader from "@/views/shared/DoubleHeader.vue";
import StationPhoto from "@/views/shared/StationPhoto.vue";
import LatestStationReadings from "@/views/shared/LatestStationReadings.vue";
import AuthenticatedPhoto from "@/views/shared/AuthenticatedPhoto.vue";
import ForbiddenBanner from "@/views/shared/ForbiddenBanner.vue";
import {
    ActionTypes,
    AuthenticationRequiredError,
    BoundingRectangle,
    DisplayModule,
    DisplayStation,
    MappedStations,
    ProjectAttribute,
    ProjectModule,
} from "@/store";
import * as utils from "@/utilities";
import { mergeNotes, NoteMedia, Notes, PortalNoteMedia, PortalStationNotes } from "@/views/notes/model";
import NotesForm from "@/views/notes/NotesForm.vue";
import StationsMap from "@/views/shared/StationsMap.vue";
import { mapGetters } from "vuex";
import ProjectAttributes from "@/views/projects/ProjectAttributes.vue";

export default Vue.extend({
    name: "StationView",
    components: {
        StandardLayout,
        DoubleHeader,
        StationPhoto,
        LatestStationReadings,
        StationsMap,
        NotesForm,
        AuthenticatedPhoto,
        ForbiddenBanner,
        ProjectAttributes,
    },
    data(): {
        selectedModule: DisplayModule | null;
        isMobileView: boolean;
        loading: boolean;
    } {
        return {
            selectedModule: null,
            isMobileView: window.screen.availWidth <= 500,
            loading: true,
        };
    },
    watch: {
        station() {
            this.loading = false;
            this.selectedModule = this.station.modules[0];
        },
    },
    computed: {
        ...mapGetters({ isBusy: "isBusy" }),
        station(): DisplayStation {
            return this.$state.stations.stations[this.$route.params.stationId];
        },
        notes(): PortalStationNotes[] {
            return this.$state.notes.notes;
        },
        media(): PortalNoteMedia[] {
            return this.$state.notes.media;
        },
        photos(): NoteMedia[] {
            return NoteMedia.onlyPhotos(this.$state.notes.media);
        },
        attributes(): ProjectAttribute[] {
            const station = this.$state.stations.stations[this.$route.params.stationId];
            return station.attributes;
        },
        headerSubtitle(): string {
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
        mapBounds(): BoundingRectangle {
            return MappedStations.defaultBounds();
        },
        mapped(): MappedStations | null {
            if (!this.station.id) {
                return null;
            }

            const mapped = MappedStations.make([this.station]);

            return mapped.focusOn(this.station.id);
        },
        showMap(): boolean {
            if (this.mapped && this.mapped.features.length > 0) {
                return true;
            }

            return false;
        },
    },
    beforeMount(): Promise<any> {
        this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.$route.params.stationId });

        return this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.$route.params.stationId }).catch((e) => {
            if (AuthenticationRequiredError.isInstance(e)) {
                return this.$router.push({
                    name: "login",
                    params: { errorMessage: this.$t("login.privateStation") },
                    query: { after: this.$route.path },
                });
            }
        });
    },
    methods: {
        getBatteryIcon(): string {
            if (!this.station || !this.station.battery) {
                return "";
            }
            return this.$loadAsset(utils.getBatteryIcon(this.station.battery));
        },
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        getModuleName(module: DisplayModule) {
            return module.name.replace("modules.", "fk.");
        },
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/mixins";
@import "src/scss/layout";
@import "src/scss/forms.scss";

* {
    box-sizing: border-box;
}

.container-box {
    border: 1px solid var(--color-border);
    border-radius: 2px;
    background-color: #fff;
    padding: 20px;
    font-size: 14px;

    @include bp-down($xs) {
        padding: 10px;
    }
}

.section {
    &-notes {
        @include bp-down($xs) {
            padding: 0;
        }
    }

    &-station {
        margin-top: 30px;
        display: flex;
        justify-content: space-between;

        @include bp-down($sm) {
            flex-wrap: wrap;
            margin-top: 20px;
            margin-bottom: 60px;
        }

        @include bp-down($xs) {
            margin-top: -10px;
        }

        > div {
            flex: 0 0 calc(50% - 10px);

            @include bp-down($sm) {
                flex-basis: 100%;
            }
        }

        ::v-deep .station-photo {
            margin-right: 20px;
            width: 90px;
            height: 90px;
            object-fit: cover;
            border-radius: 5px;
        }

        .photo-container {
            flex: 0 0 calc(50% - 5px);
            margin-bottom: 10px;
            height: calc(50% - 5px);
            max-height: 192px;

            &:nth-of-type(3),
            &:nth-of-type(4) {
                margin-bottom: 0;
            }

            > img {
                width: 100%;
                height: 100%;
                object-fit: cover;
                border-radius: 2px;
            }
        }
    }
    &-readings {
        padding: 0;
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
            margin-bottom: 5px;
            width: 30px;
            height: 30px;
        }
    }
    &-coordinate {
        padding: 0 15px;

        @include bp-down($xs) {
            display: flex;
        }

        &:last-of-type {
            padding-right: 0;
        }

        div:nth-of-type(2) {
            font-size: 12px;
            margin-top: 2px;

            @include bp-down($xs) {
                order: -1;
                margin-right: 10px;
            }
        }
    }
    &-location {
        align-self: flex-start;
        padding-right: 10px;

        @include bp-down($xs) {
            flex: 0 0 100%;
            border-top: 1px solid var(--color-border);
            padding-top: 15px;
            margin-top: 15px;
            order: 3;
        }
    }
    &-row {
        padding: 15px 0;
        max-width: 350px;
        @include flex(center);

        &:first-of-type {
            @include bp-down($xs) {
                padding-top: 30px;
            }
        }

        &:not(:last-of-type) {
            border-bottom: solid 1px var(--color-border);
        }

        @include bp-down($sm) {
            max-width: unset;

            &:last-of-type {
                padding-bottom: 5px;
            }
        }

        @include bp-down($xs) {
            flex-wrap: wrap;
        }

        .icon {
            margin-right: 7px;
            align-self: flex-start;
        }
    }

    &-readings {
        font-size: 16px;
        display: flex;
        min-height: 130px;
        position: relative;

        @include bp-down($xs) {
            padding-top: 54px;
        }

        &-values {
            padding: 27px 20px 10px 30px;
            border-left: 1px solid var(--color-border);
            transform: translateX(-1px);
            width: 100%;

            @include bp-down($xs) {
                padding: 20px 25px;
            }
        }

        ul {
            z-index: $z-index-top;
        }

        li {
            @include flex(center);
            min-width: 250px;
            padding: 13px 16px;
            cursor: pointer;
            border-right: 1px solid var(--color-border);
            transition: border-left-width linear 0.25s;
            border-bottom: 1px solid var(--color-border);

            @include bp-down($sm) {
                padding: 10px 20px;
                min-width: unset;
            }

            &.active {
                border-right: 1px solid #fff;
                border-left: solid 4px var(--color-primary);
                padding-left: 12px;
                cursor: initial;

                @include bp-down($sm) {
                    padding-left: 16px;
                }
            }

            img {
                width: 40px;
                height: 40px;
            }

            span {
                margin-left: 10px;

                @include bp-down($sm) {
                    display: none;
                }
            }
        }

        ::v-deep .reading-container {
            flex: 0 0 calc(33% - 10px);
            margin-right: 10px;

            &:last-of-type {
                margin-bottom: 0;
            }

            @include bp-down($lg) {
                flex: 0 0 calc(50% - 10px);
            }

            @media screen and (max-width: 600px) {
                flex: 0 0 100%;
            }
        }

        header {
            font-size: 20px;
            padding-bottom: 15px;
            border-bottom: 1px solid var(--color-border);
            margin-bottom: 20px;
            width: 100%;
            display: flex;

            img {
                width: 20px;
                height: 20px;
                margin-right: 10px;
            }

            @include bp-down($xs) {
                padding: 15px 10px;
                font-size: 18px;
                @include position(absolute, 0 null null 0);
            }
        }
    }

    &-photos {
        @include flex;
        flex-wrap: wrap;
        justify-content: space-between;
        position: relative;

        @include bp-down($sm) {
            margin-top: 20px;
        }

        &-nav {
            height: 28px;
            padding: 0px 20px;
            border-radius: 3px;
            border: solid 1px #cccdcf;
            background-color: #fff;
            font-size: 14px;
            font-weight: 900;
            @include flex(center, center);
            @include position(absolute, null 20px 20px null);

            @include bp-down($sm) {
                position: unset;
                width: 100%;
                margin-top: 5px;
                height: 35px;
            }
        }

        .photo {
            flex: 0 0 calc(50% - 5px);
        }
    }
}

.small-light {
    font-size: 12px;
    color: #6a6d71;
}

.stations-map {
    height: 400px;

    @include bp-down($sm) {
        height: 450px;
    }
}

section {
    margin-bottom: 20px;
}

.loading-container {
    height: 100%;
    @include flex(center);
}
.notes-view .lower .loading-container.empty {
    padding: 20px;
}

.icon-location {
    margin-top: 1px;
}

.photo-placeholder {
    @include flex(center, center);
    height: 100%;
    background-color: var(--color-border);
    opacity: 0.5;

    img {
        width: 40%;
        height: 40%;
    }
}

.icon-grid {
    font-size: 12px;
    margin-right: 5px;
}

.forbidden-view-bg {
    background-image: linear-gradient(#52b5e4, #1b80c9);
    min-height: 100vh;
    padding-top: 50px;

    body.floodnet & {
        background-image: unset;
        background-color: var(--color-dark);
    }

    a {
        color: #fff;
        margin-top: 20px;
        display: block;
    }
}

.message-container {
    margin: 0 auto;
}
</style>
