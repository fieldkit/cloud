<template>
    <StandardLayout>
        <div class="container-wrap" v-if="station">
            <DoubleHeader
                :backRoute="projectId ? 'viewProject' : 'mapStation'"
                :backTitle="projectId ? $tc('layout.backProjectDashboard') : $tc(partnerCustomization().nav.viz.back.map.label)"
                :backRouteParams="{ id: projectId || station.id }"
            >
                <template v-slot:default>
                    <a v-for="link in partnerCustomization().links" v-bind:key="link.url" :href="link.url" target="_blank" class="link">
                        {{ $t(link.text) }} >
                    </a>
                </template>
            </DoubleHeader>

            <section class="section-station">
                <div class="container-box">
                    <div class="flex flex-al-center">
                        <StationPhoto :station="station" />
                        <div>
                            <div class="station-name">{{ station.name }}</div>

                            <div class="flex flex-al-center flex-wrap">
                                <span class="station-deployed-date">{{ deployedDate }}</span>
                                <div class="station-owner">
                                    <UserPhoto :user="station.owner" />
                                    {{ $t("station.belongsTo") }}
                                    <span>{{ station.owner.name }}</span>
                                </div>
                            </div>

                            <div v-if="partnerCustomization().stationLocationName(station)" class="flex station-location">
                                <i class="icon icon-location"></i>
                                <span>{{ partnerCustomization().stationLocationName(station) }}</span>
                            </div>
                            <div v-if="station.placeNameNative" class="station-location">
                                <i class="icon icon-location"></i>
                                <span>{{ $tc("station.nativeLand") }} {{ station.placeNameNative }}</span>
                            </div>

                            <div v-if="station.location" class="flex">
                                <div class="station-coordinate">
                                    <span class="bold">{{ $tc("station.latitude") }}</span>
                                    <span>{{ station.location.latitude | prettyCoordinate }}</span>
                                </div>

                                <div class="station-coordinate">
                                    <span class="bold">{{ $tc("station.longitude") }}</span>
                                    <span>{{ station.location.longitude | prettyCoordinate }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div>
                        <div class="station-row">
                            <span class="bold">{{ $tc("station.modules") }}</span>
                            <div class="station-modules ml-10">
                                <img
                                    v-for="module in station.modules"
                                    v-bind:key="module.name"
                                    alt="Module icon"
                                    :src="getModuleImg(module)"
                                />
                            </div>
                        </div>

                        <div class="station-row">
                            <StationBattery :station="station" />
                        </div>

                        <div v-if="station.firmwareNumber" class="station-row">
                            <span class="bold">{{ $tc("station.firmwareVersion") }}</span>
                            <span class="ml-10 small-light">{{ station.firmwareNumber }}</span>
                        </div>
                    </div>
                </div>
                <div v-if="photos" class="station-photos">
                    <div class="photo-container cursor-pointer" v-for="(n, index) in 4" v-bind:key="index" @click="navigateToPhotos()">
                        <AuthenticatedPhoto v-if="photos[index]" :url="photos[index].url" />
                        <div v-else class="photo-placeholder">
                            <img src="@/assets/image-placeholder-v2.svg" alt="Image placeholder" />
                        </div>
                    </div>
                    <a @click="navigateToPhotos()" class="station-photos-nav cursor-pointer">
                        <i class="icon icon-grid"></i>
                        {{ $t("station.btn.linkToPhotos") }}
                    </a>
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
                    <StationsMap :mapped="mapped" :showStations="true" :mapBounds="mapped.bounds" :visibleReadings="visibleReadings" />
                </div>
            </section>

            <section v-if="notes && !isCustomizationEnabled()" class="section-notes container-box">
                <NotesForm
                    v-bind:key="station.id"
                    :station="station"
                    :notes="{ notes, media }"
                    :readonly="station.readOnly"
                    @change="dirtyNotes = true"
                />
            </section>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "@/views/StandardLayout.vue";
import DoubleHeader from "@/views/shared/DoubleHeader.vue";
import StationPhoto from "@/views/shared/StationPhoto.vue";
import LatestStationReadings from "@/views/shared/LatestStationReadings.vue";
import AuthenticatedPhoto from "@/views/shared/AuthenticatedPhoto.vue";
import {
    ActionTypes,
    AuthenticationRequiredError,
    BoundingRectangle,
    DisplayModule,
    DisplayStation,
    MappedStations,
    ProjectAttribute,
    ProjectModule,
    VisibleReadings,
} from "@/store";
import * as utils from "@/utilities";
import { NoteMedia, PortalNoteMedia, PortalStationNotes } from "@/views/notes/model";
import NotesForm from "@/views/notes/NotesForm.vue";
import StationsMap from "@/views/shared/StationsMap.vue";
import ProjectAttributes from "@/views/projects/ProjectAttributes.vue";
import StationBattery from "@/views/station/StationBattery.vue";
import { getPartnerCustomizationWithDefault, isCustomisationEnabled, PartnerCustomization } from "@/views/shared/partners";
import UserPhoto from "@/views/shared/UserPhoto.vue";

export default Vue.extend({
    name: "StationView",
    components: {
        StationBattery,
        StandardLayout,
        DoubleHeader,
        StationPhoto,
        LatestStationReadings,
        StationsMap,
        NotesForm,
        AuthenticatedPhoto,
        ProjectAttributes,
        UserPhoto,
    },
    data(): {
        selectedModule: DisplayModule | null;
        isMobileView: boolean;
        loading: boolean;
        dirtyNotes: boolean;
    } {
        return {
            selectedModule: null,
            isMobileView: window.screen.availWidth <= 500,
            loading: true,
            dirtyNotes: false,
        };
    },
    watch: {
        station() {
            this.loading = false;
            this.selectedModule = this.station.modules[0];
        },
    },
    computed: {
        visibleReadings(): VisibleReadings {
            return VisibleReadings.Current;
        },
        projectId(): string {
            return this.$route.params.projectId;
        },
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
        deployedDate(): string | null {
            if (this.station) {
                const deploymentDate = this.partnerCustomization().getStationDeploymentDate(this.station);
                if (deploymentDate) {
                    return this.$tc("station.deployed") + " " + deploymentDate;
                }
                return this.$tc("station.readyToDeploy");
            }
            return null;
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
    beforeRouteLeave(to: never, from: never, next: any) {
        if (this.dirtyNotes) {
            this.$confirm({
                message: this.$tc("notes.confirmLeavePopupMessage"),
                button: {
                    no: this.$tc("no"),
                    yes: this.$tc("yes"),
                },
                callback: (confirm) => {
                    if (confirm) {
                        this.dirtyNotes = false;
                        next();
                    }
                },
            });
        } else {
            next();
        }
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
        partnerCustomization(): PartnerCustomization {
            return getPartnerCustomizationWithDefault();
        },
        isCustomizationEnabled(): boolean {
            return isCustomisationEnabled();
        },
        navigateToPhotos(): void {
            this.$router.push({
                name: this.projectId ? "viewProjectStationPhotos" : "viewStationPhotos",
                params: { projectId: this.projectId, stationId: this.station.id },
            });
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
            flex: 0 0 90px;
            height: 90px;
            object-fit: cover;
            margin-right: 20px;
            border-radius: 5px;
        }

        .photo-container {
            flex: 0 0 calc(50% - 5px);
            margin-bottom: 10px;
            height: calc(50% - 5px);
            position: relative;
            border-radius: 2px;
            overflow: hidden;
            background-color: #e2e4e6;

            &:nth-of-type(3),
            &:nth-of-type(4) {
                margin-bottom: 0;
            }

            ::v-deep img {
                width: 100%;
                height: 100%;
                object-fit: cover;
            }

            .photo-placeholder {
                @include flex(center, center);
                height: 100%;

                img {
                    width: 40%;
                    height: 40%;
                    object-fit: contain;
                }
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
        margin-bottom: 8px;
    }
    &-battery {
        margin-top: 5px;
        @include flex(flex-start);

        span {
            margin-left: 5px;
        }
    }
    &-modules {
        margin-left: 10px;
        flex-wrap: wrap;
        @include flex;

        img {
            margin-right: 8px;
            margin-bottom: 5px;
            width: 25px;
            height: 25px;
        }
    }
    &-coordinate {
        font-size: 12px;

        @include bp-down($xs) {
            display: flex;
        }

        &:last-of-type {
            margin-left: 15px;
        }

        span:nth-of-type(2) {
            margin-left: 2px;

            @include bp-down($xs) {
                order: -1;
                margin-right: 10px;
            }
        }
    }
    &-location {
        align-items: flex-start;
        margin-bottom: 7px;
    }
    &-row {
        padding: 15px 0;
        @include flex(center);

        &:not(:last-of-type) {
            border-bottom: solid 1px var(--color-border);
        }

        &:last-of-type {
            padding-bottom: 0;
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
        max-height: 390px;

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
    }

    &-deployed-date {
        color: #6a6d71;
        margin-bottom: 10px;
        margin-right: 5px;
    }

    &-owner {
        color: #6a6d71;
        font-size: 10px;
        margin-bottom: 10px;
        @include flex(center);

        ::v-deep .default-user-icon {
            width: 18px;
            height: 18px;
            margin-top: 0;
            margin-right: 5px;
        }

        span {
            margin-left: 3px;
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
    margin-right: 8px;
}

.icon-grid {
    font-size: 12px;
    margin-right: 5px;
}

.message-container {
    margin: 0 auto;
}

.double-header {
    .link {
        color: $color-primary;
        font-size: 12px;
        letter-spacing: 0.07px;
        text-decoration: initial;

        body.floodnet & {
            color: $color-dark;
        }
    }
}

::v-deep .back {
    margin-bottom: 15px;
}
</style>
