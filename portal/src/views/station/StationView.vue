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
                <div v-if="notes.media" class="station-photos">
                    <div class="photo-container" v-for="(n, index) in 4" v-bind:key="index">
                        <!-- somehow using v-for like so needs the next v-if -->
                        <AuthenticatedPhoto v-if="notes.media[index]" :url="notes.media[index].url" />
                        <i v-else class="photo icon icon-image-placeholder"></i>
                    </div>
                    <router-link :to="{ name: 'test' }" class="station-photos-nav">
                        {{ $t("station.viewAllPhotos") }}
                    </router-link>
                </div>
            </section>

            <section class="container-box section-readings">
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
                    <div class="station-readings-values">
                        <!--                        <header>{{ getModuleName(selectedModule) }}</header>-->
                        <LatestStationReadings :id="station.id" :moduleKey="selectedModule.name" />
                    </div>
                </div>
            </section>

            <section>
                <div class="container-map">
                    <!--                    <StationMap :mapped="mapped" v-if="mapped && false" :showStations="true" />-->
                </div>
            </section>

            <section v-if="notes && notes.notes" class="container-box">
                <div class="notifications">
                    <div v-if="notesState.failed" class="notification failed">{{ $tc("notes.failed") }}</div>

                    <div v-if="notesState.success" class="notification success">{{ $tc("notes.success") }}</div>
                </div>
                <NotesForm
                    v-bind:key="station.id"
                    :station="station"
                    :notes="notes"
                    :readonly="false"
                    @save="saveForm"
                    @change="onChange"
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
import { ActionTypes, DisplayStation, MappedStations, ModuleSensor, ProjectModule } from "@/store";
import * as utils from "@/utilities";
import StationMap from "@/views/station/StationMap.vue";
import { mergeNotes, NoteMedia, Notes, PortalStationNotesReply } from "@/views/notes/model";
import NotesForm from "@/views/notes/NotesForm.vue";
import { serializePromiseChain } from "@/utilities";
import { convertOldFirmwareResponse } from "@/utilities";

export default Vue.extend({
    name: "StationView",
    components: {
        StandardLayout,
        DoubleHeader,
        StationPhoto,
        LatestStationReadings,
        // StationMap,
        NotesForm,
        AuthenticatedPhoto,
    },
    mounted() {
        // this.$store.dispatch(ActionTypes.NEED_STATION, { id: this.$route.params.id });

        this.$services.api
            .getStation(this.$route.params.id)
            .then((station) => {
                this.station = new DisplayStation(station);
                this.selectedModule = this.station.modules[0];

                console.log("radoi station", this.station);
                //    console.log("radoi module", this.selectedModule);

                //   console.log("gettinggg radoi station notes by id", this.station.id);
                // get notes & media
                this.$services.api
                    .getStationNotes(this.station.id)
                    .then((notes) => {
                        //        console.log("Radoi notes", notes);
                        this.notes = notes;

                        console.log("GOT NOTES", this.notes);
                    })
                    .catch((e) => {
                        console.log("Radoi e notes", e);
                    });
            })
            .catch((e) => {
                console.log("e radoi station", e);
            });
    },
    data(): {
        station: DisplayStation | null;
        notes: { [stationId: number]: PortalStationNotesReply };
        notesState: {
            dirty: boolean;
            success: boolean;
            failed: boolean;
        };
        selectedModule: ModuleSensor | null;
    } {
        return {
            station: null,
            notes: {},
            notesState: {
                dirty: false,
                success: false,
                failed: false,
            },
            selectedModule: null,
        };
    },
    computed: {
        photos(this: any) {
            return NoteMedia.onlyPhotos(this.notes.media);
        },
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
        /*station(): DisplayStation {
            console.log("rrr", this.$getters.stationsById[this.$route.params.id]);
            return this.$getters.stationsById[this.$route.params.id];
        },*/
        mapped(): any {
            const stationId = parseInt(this.$route.params.id);
            const station = this.$getters.mapped.stations.filter((station) => station.id === stationId);

            if (!this.$getters.mapped) {
                return null;
            }
            /* if (this.bounds) {
                console.log(`focusing bounds: ${this.bounds}`);
                return this.$getters.mapped.overrideBounds(this.bounds);
            }*/

            /* if (this.station.id) {
                console.log(`focusing station: ${this.id}`);
                return this.$getters.mapped.focusOn(this.id);
            }*/

            const mapped = this.$getters.mapped;
            console.log("radoi mapped", mapped);

            const returnValue = {
                features: mapped.features,
                bounds: mapped.bounds,
                stations: station,
                isSingleType: mapped.isSingleType,
            };

            console.log("radoi filtered", returnValue);

            return returnValue;
        },
    },
    beforeRouteUpdate(to: never, from: never, next: any) {
        if (this.confirmLeave()) {
            next();
        }
    },
    beforeRouteLeave(to: never, from: never, next: any) {
        if (this.confirmLeave()) {
            next();
        }
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
        async saveForm(formNotes: Notes): Promise<void> {
            this.notesState.success = false;
            this.notesState.failed = false;

            await serializePromiseChain(formNotes.addedPhotos, (photo) => {
                return this.$services.api.uploadStationMedia(this.station.id, photo.key, photo.file).then((media) => {
                    return [];
                });
            }).then(() => {
                const payload = mergeNotes(this.notes[this.station.id], formNotes);
                return this.$services.api.patchStationNotes(this.station.id, payload).then(
                    (updated) => {
                        this.notesState.dirty = false;
                        this.notesState.success = true;
                    },
                    () => {
                        this.notesState.failed = true;
                    }
                );
            });
        },
        onChange(): void {
            this.notesState.dirty = true;
        },
        confirmLeave(): boolean {
            if (this.notesState.dirty) {
                if (window.confirm("You may have unsaved changes, are you sure you'd like to leave?")) {
                    this.notesState.dirty = false;
                    return true;
                } else {
                    return false;
                }
            }
            return true;
        },
        getModuleName(module) {
            return module.name.replace("modules.", "fk.");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/layout";

* {
    box-sizing: border-box;
}

.container-box {
    border: 1px solid var(--color-border);
    border-radius: 2px;
    background-color: #fff;
    padding: 20px;
    font-size: 14px;
}

.section {
    &-station {
        margin-top: 30px;
        display: flex;
        justify-content: space-between;

        > div {
            flex: 0 0 calc(50% - 10px);
        }

        ::v-deep .station-photo {
            margin-right: 20px;
            object-fit: contain;
            width: 90px;
            height: 90px;
            border-radius: 5px;
        }

        .photo-container {
            flex: 0 0 calc(50% - 5px);
            margin-bottom: 10px;
            height: calc(50% - 5px);

            img {
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

    &-readings {
        font-size: 16px;
        display: flex;
        min-height: 130px;

        &-values {
            padding: 27px 20px 10px 30px;
            border-left: 1px solid var(--color-border);
            transform: translateX(-1px);
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

            @include bp-down($lg) {
                flex: 0 0 calc(50% - 10px);
            }

            @include bp-down($sm) {
                flex: 0 0 100%;
            }
        }

        header {
            font-size: 20px;
            padding-bottom: 15px;
            border-bottom: 1px solid var(--color-border);
            margin-bottom: 20px;
        }
    }

    &-photos {
        @include flex;
        flex-wrap: wrap;
        justify-content: space-between;
        position: relative;

        &-nav {
            height: 28px;
            padding: 5px 20px;
            border-radius: 3px;
            border: solid 1px #cccdcf;
            background-color: #fff;
            font-size: 14px;
            font-weight: 900;
            @include position(absolute, null 20px 20px null);
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

.notes-form {
    padding: 0;

    ::v-deep {
        .header {
            padding-bottom: 15px;
        }
        .note-editor {
            margin-top: 0;
            margin-bottom: 20px;
            padding: 15px 20px 5px;
        }
    }
}
::v-deep .site-notes {
    margin-top: 30px;
}

section {
    margin-bottom: 20px;
}
</style>
