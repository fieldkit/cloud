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
                        <div class="w-100">
                            <div class="station-name">{{ station.name }}</div>

                            <div class="flex flex-al-center flex-wrap">
                                <span class="station-deployed-date">{{ deployedDate }}</span>
                                <div class="station-owner">
                                    <UserPhoto :user="station.owner" />
                                    {{ $t("station.belongsTo") }}
                                    <span>{{ station.owner.name }}</span>
                                </div>
                            </div>

                            <div v-if="!isPartnerCustomisationEnabled" class="station-description">
                                <textarea
                                    v-if="form.description !== undefined"
                                    class="input"
                                    oninput='this.style.height = "";this.style.height = this.scrollHeight + "px"'
                                    v-model="form.description"
                                    :disabled="!editingDescription"
                                />
                                <a
                                    v-if="form.description === undefined"
                                    @click="
                                        form.description = '';
                                        editingDescription = true;
                                    "
                                    class="station-description-add"
                                >
                                    {{ $t("station.addDescription") }}
                                </a>

                                <template>
                                    <a
                                        v-if="form.description && !editingDescription"
                                        @click="editingDescription = true"
                                        class="station-description-edit"
                                    >
                                        {{ $t("edit") }}
                                    </a>
                                    <a
                                        @click="saveStationDescription()"
                                        v-if="editingDescription && form.description"
                                        class="station-description-edit"
                                        style="margin-top: 4px;"
                                    >
                                        {{ $t("save") }}
                                    </a>
                                </template>
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
                                    v-for="(module, moduleIndex) in station.modules"
                                    v-bind:key="moduleIndex"
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
                <div>
                    <div class="station-projects" v-if="stationProjects.length > 0">
                        <template v-if="stationProjects.length === 1">{{ $tc("station.singleProjectTitle") }}&nbsp;</template>
                        <template v-else>{{ $tc("station.multipleProjectsTitle") }} &nbsp;</template>
                        <router-link
                            v-for="(project, index) in stationProjects"
                            v-bind:key="project.id"
                            :to="{ name: 'viewProject', params: { id: project.id } }"
                            target="_blank"
                        >
                            {{ project.name }}
                            <template v-if="stationProjects.length > 1 && index !== stationProjects.length - 1">,&nbsp;</template>
                        </router-link>
                    </div>
                    <div v-if="photos" class="station-photos">
                        <div class="photo-container" v-for="(n, index) in 4" v-bind:key="index" @click="navigateToPhotos()">
                            <AuthenticatedPhoto v-if="photos[index]" :url="photos[index].url" />
                            <div v-else class="photo-placeholder">
                                <img src="@/assets/image-placeholder-v2.svg" alt="Image placeholder" />
                            </div>
                        </div>
                        <a class="station-photos-nav" @click="navigateToPhotos()">
                            <i class="icon icon-grid"></i>
                            {{ $t("station.btn.linkToPhotos") }}
                        </a>
                    </div>
                </div>
            </section>

            <section class="container-box" v-if="station.modules.length > 0">
                <h2>{{ $t("station.data") }}</h2>

                <ul class="flex flex-wrap flex-space-between module-data-container">
                    <li
                        class="module-data-item"
                        v-for="module in station.modules"
                        v-bind:key="module.name"
                        @click="onModuleClick(module.id)"
                    >
                        <h3 class="module-data-title flex flex-al-center">
                            <img alt="Module icon" :src="getModuleImg(module)" />
                            {{ $t(getModuleName(module)) }}
                        </h3>
                        <TinyChart
                            :ref="'tinyChart-' + module.id"
                            :moduleKey="module.name"
                            :station-id="station.id"
                            :station="station"
                            :querier="sensorDataQuerier"
                        />
                    </li>
                </ul>
                <button class="btn module-data-btn" @click="onClickExplore">{{ $t("station.exploreData") }}</button>
            </section>

            <section class="container-box section-readings" v-if="selectedModule">
                <div class="station-readings">
                    <ul>
                        <li
                            v-for="(module, moduleIndex) in station.modules"
                            v-bind:key="moduleIndex"
                            :class="{ active: module.name === selectedModule.name }"
                            @click="selectModule(module)"
                        >
                            <img alt="Module icon" :src="getModuleImg(module)" />
                            <input
                                v-if="editedModule && editedModule.id === module.id"
                                class="input"
                                maxlength="25"
                                :disabled="editedModule.id !== selectedModule.id"
                                :title="editedModule.label"
                                v-model="editedModule.label"
                            />
                            <input
                                v-else
                                class="input"
                                maxlength="25"
                                :title="module.label ? module.label : $t(getModuleName(module))"
                                disabled
                                :value="module.label ? module.label : $t(getModuleName(module))"
                            />
                            <template v-if="!isCustomizationEnabled()">
                                <a
                                    v-if="!editedModule || (editedModule && editedModule.id !== module.id)"
                                    @click="onEditModuleNameClick(module)"
                                    class="module-edit-name"
                                >
                                    {{ $t("edit") }}
                                </a>
                                <a v-if="editedModule && editedModule.id === module.id" @click="saveModuleName()" class="module-edit-name">
                                    {{ $t("save") }}
                                </a>
                            </template>
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
                <NotesForm v-bind:key="station.id" :station="station" :readonly="station.readOnly" @change="dirtyNotes = true" />
            </section>

            <section class="section-notes container-box">
                <FieldNotes :stationName="station.name"></FieldNotes>
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
    GlobalState,
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
import { Project } from "@/api";
import { mapState } from "vuex";
import { SensorDataQuerier } from "@/views/shared/sensor_data_querier";
import TinyChart from "@/views/viz/TinyChart.vue";
import { BookmarkFactory, serializeBookmark } from "@/views/viz/viz";
import { ExploreContext } from "@/views/viz/common";
import FieldNotes from "@/views/fieldNotes/FieldNotes.vue";

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
        TinyChart,
        UserPhoto,
        FieldNotes,
    },
    data(): {
        selectedModule: DisplayModule | null;
        isMobileView: boolean;
        loading: boolean;
        dirtyNotes: boolean;
        sensorDataQuerier: SensorDataQuerier;
        editModuleIndex: number | null;
        editingDescription: boolean;
        form: {
            description: string | null;
        };
        editedModule: DisplayModule | null;
    } {
        return {
            selectedModule: null,
            isMobileView: window.screen.availWidth <= 500,
            loading: true,
            dirtyNotes: false,
            editedModule: null,
            editModuleIndex: null,
            editingDescription: false,
            form: {
                description: null,
            },
            sensorDataQuerier: new SensorDataQuerier(this.$services.api),
        };
    },
    watch: {
        station() {
            this.loading = false;
            this.selectedModule = this.station.modules[0];
            this.form.description = this.station.description;
        },
    },
    computed: {
        ...mapState({
            userStations: (s: GlobalState) => Object.values(s.stations.user.stations),
        }),
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
        photos(): NoteMedia[] | null {
            if (this.$state.notes.media) {
                return NoteMedia.onlyPhotos(this.$state.notes.media);
            }
            return null;
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
        stationProjects(): Project[] {
            return this.$store.getters.stationProjects;
        },
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
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
        const stationId = this.$route.params.stationId;

        this.$store.dispatch(ActionTypes.NEED_NOTES, { id: stationId });
        this.$store.dispatch(ActionTypes.NEED_PROJECTS_FOR_STATION, { id: this.$route.params.stationId });

        return this.$store.dispatch(ActionTypes.NEED_STATION, { id: stationId }).catch((e) => {
            if (AuthenticationRequiredError.isInstance(e)) {
                return this.$router.push({
                    name: "login",
                    params: { errorMessage: this.$t("login.privateStation").toString() },
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
            if (!module.label) {
                return module.name.replace("modules.", "fk.");
            }
            return module.label;
        },
        partnerCustomization(): PartnerCustomization {
            return getPartnerCustomizationWithDefault();
        },
        isCustomizationEnabled(): boolean {
            return isCustomisationEnabled();
        },
        onClickExplore(): void {
            const exploreContext = new ExploreContext();
            const bm = BookmarkFactory.forStation(this.station.id, exploreContext);
            const url = this.$router.resolve({
                name: "exploreBookmark",
                query: { bookmark: serializeBookmark(bm) },
            }).href;
            window.open(url, "_blank");
        },
        saveStationDescription(): void {
            const payload = { id: this.station.id, name: this.station.name, ...this.form };
            this.$store.dispatch(ActionTypes.UPDATE_STATION, payload);
            this.editingDescription = false;
        },
        onEditModuleNameClick(module: DisplayModule): void {
            this.editedModule = JSON.parse(JSON.stringify(module));
            if (this.editedModule) {
                this.editedModule.label = this.$tc(this.getModuleName(module));
            }
        },
        saveModuleName(): void {
            if (!this.editedModule) {
                return;
            }
            const payload = { stationId: this.station.id, moduleId: this.editedModule.id, label: this.editedModule.label };
            this.$store.dispatch(ActionTypes.UPDATE_STATION_MODULE, payload).then(() => {
                this.editedModule = null;
            });
        },
        selectModule(module: DisplayModule) {
            this.selectedModule = module;
        },
        navigateToPhotos(): void {
            this.$router.push({
                name: this.projectId ? "viewProjectStationPhotos" : "viewStationPhotos",
                params: { projectId: this.projectId, stationId: String(this.station.id) },
            });
        },

        onModuleClick(moduleId: number) {
            const tinyChartComp = this.$refs["tinyChart-" + moduleId];
            if (tinyChartComp && tinyChartComp[0]) {
                const vizData = tinyChartComp[0].vizData;
                if (vizData) {
                    const bm = BookmarkFactory.forSensor(this.station.id, vizData.vizSensor, vizData.timeRange);
                    const url = this.$router.resolve({
                        name: "exploreBookmark",
                        query: { bookmark: serializeBookmark(bm) },
                    }).href;
                    window.open(url, "_blank");
                }
            }
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
    padding: 15px 20px;
    font-size: 14px;

    @include bp-down($xs) {
        padding: 10px;
    }
}

.container-title {
    font-size: 20px;
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
            align-self: start;
        }

        .photo-container {
            flex: 0 0 calc(50% - 5px);
            margin-bottom: 10px;
            height: calc(50% - 5px);
            position: relative;
            border-radius: 2px;
            overflow: hidden;
            background-color: #e2e4e6;
            cursor: pointer;

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
            width: 300px;
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

                input {
                    cursor: initial;
                }
            }

            img {
                width: 40px;
                height: 40px;
            }

            input {
                margin-left: 5px;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                width: 100%;
                cursor: pointer;
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
        height: 390px;

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

    &-description {
        font-size: 12px;
        color: #6a6d71;
        margin-top: -2px;
        margin-bottom: 10px;
        display: flex;
        align-items: flex-start;

        a {
            cursor: pointer;
        }

        .input {
            font-size: 12px;
            width: 100%;
            resize: none;
            overflow: hidden;

            &:disabled {
                padding: 0;
            }
        }
    }

    &-description-edit {
        margin-left: auto;
        padding-left: 5px;
        cursor: pointer;
    }

    &-description-add {
        opacity: 0.5;
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

    &-projects {
        font-size: 16px;
        color: #6a6d71;
        margin: 30px 0;

        @include bp-down($xs) {
            margin: 20px 0;
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

.module-data-container {
    gap: 20px;

    @include bp-down($sm) {
        gap: 10px;
    }
}

.module-data-item {
    flex: 1 1 calc(50% - 10px);
    min-width: 0;

    @include bp-down($sm) {
        flex: 0 0 100%;
    }
}

.module-data-title {
    color: $color-primary;
    font-size: 12px;
    margin-bottom: 10px;
    cursor: pointer;

    img {
        margin-right: 7px;
        width: 19px;
        height: 19px;
    }
}

.module-data-btn {
    margin: 30px auto 8px auto;
    display: block;
}

.module-edit-name {
    opacity: 0.4;
    font-size: 12px;
    margin-left: 7px;
    margin-bottom: -1px;
}

::v-deep .back {
    margin-bottom: 15px;
}

::v-deep .back {
    margin-bottom: 15px;
}

.module-edit-name {
    opacity: 0.4;
    font-size: 12px;
    margin-bottom: -1px;
    cursor: pointer;
}
</style>
