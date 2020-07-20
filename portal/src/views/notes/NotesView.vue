<template>
    <StandardLayout>
        <div class="notes-view">
            <DoubleHeader :title="project.name" subtitle="Field Notes" backTitle="Back to dashboard" backRoute="projects" v-if="project" />
            <DoubleHeader title="My Stations" subtitle="Field Notes" backTitle="Back to dashboard" backRoute="projects" v-if="!project" />

            <div class="lower">
                <div class="side">
                    <StationTabs :stations="visibleStations" :selected="selectedStation" @selected="onSelected" />
                </div>
                <template v-if="loading">
                    <div class="main">
                        <Spinner />
                    </div>
                </template>
                <template v-else>
                    <div class="main" v-if="selectedStation && selectedNotes">
                        <div class="notifications">
                            <div v-if="failed" class="notification failed">
                                Oops, there was a problem.
                            </div>

                            <div v-if="success" class="notification success">
                                Saved.
                            </div>
                        </div>
                        <NotesForm
                            :station="selectedStation"
                            :notes="selectedNotes"
                            @save="saveForm"
                            v-bind:key="stationId"
                            @change="onChange"
                        />
                    </div>
                    <div v-else class="main empty">
                        Please choose a station from the left.
                    </div>
                </template>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import Promise from "bluebird";
import CommonComponents from "@/views/shared";
import StandardLayout from "../StandardLayout.vue";
import StationTabs from "./StationTabs.vue";
import NotesForm from "./NotesForm.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

import { serializePromiseChain } from "@/utilities";

import FKApi from "@/api/api";

import { Notes, mergeNotes } from "./model";

export default Vue.extend({
    name: "NotesView",
    components: {
        ...CommonComponents,
        StandardLayout,
        StationTabs,
        NotesForm,
    },
    props: {
        projectId: {
            type: Number,
            required: false,
        },
        stationId: {
            type: Number,
            required: false,
        },
    },
    data: () => {
        return {
            notes: {},
            dirty: false,
            loading: false,
            success: false,
            failed: false,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
        }),
        project() {
            if (this.projectId) {
                return this.$getters.projectsById[this.projectId];
            }
            return null;
        },
        visibleStations() {
            if (this.projectId) {
                const project = this.$getters.projectsById[this.projectId];
                if (project) {
                    return project.stations;
                }
                return [];
            }
            return this.$store.state.stations.user.stations;
        },
        selectedStation() {
            if (this.stationId) {
                const station = this.$getters.stationsById[this.stationId];
                if (station) {
                    return station;
                }
            }
            return null;
        },
        selectedNotes() {
            if (this.stationId && this.notes) {
                return this.notes[this.stationId];
            }
            return null;
        },
    },
    watch: {
        stationId(this: any) {
            return this.loadNotes(this.stationId);
        },
    },
    mounted(this: any) {
        const pending = [];
        if (this.projectId) {
            pending.push(this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: this.projectId }));
        }
        if (this.stationId) {
            pending.push(this.loadNotes(this.stationId));
        }
        return Promise.all(pending);
    },
    beforeRouteUpdate(this: any, to, from, next) {
        console.log("router: update");
        if (this.maybeConfirmLeave()) {
            next();
        }
    },
    beforeRouteLeave(this: any, to, from, next) {
        console.log("router: leave");
        if (this.maybeConfirmLeave()) {
            next();
        }
    },
    methods: {
        loadNotes(this: any, stationId: number) {
            this.success = false;
            this.failed = false;
            Vue.set(this, "loading", true);
            return new FKApi().getStationNotes(stationId).then((notes) => {
                Vue.set(this.notes, stationId, notes);
                Vue.set(this, "loading", false);
            });
        },
        onSelected(this: any, station) {
            if (this.stationId != station.id) {
                return this.$router.push({
                    name: this.projectId ? "viewProjectStationNotes" : "viewStationNotes",
                    params: {
                        projectId: this.projectId,
                        stationId: station.id,
                    },
                });
            }
        },
        onChange() {
            this.dirty = true;
        },
        maybeConfirmLeave() {
            if (this.dirty) {
                if (window.confirm("You may have unsaved changes, are you sure you'd like to leave?")) {
                    this.dirty = false;
                    return true;
                } else {
                    return false;
                }
            }
            return true;
        },
        saveForm(this: any, formNotes: Notes) {
            this.success = false;
            this.failed = false;

            return serializePromiseChain(formNotes.addedPhotos, (photo) => {
                return new FKApi().uploadStationMedia(this.stationId, photo.key, photo.file).then((media) => {
                    console.log(media);
                    return [];
                });
            }).then(() => {
                const payload = mergeNotes(this.notes[this.stationId], formNotes);
                return new FKApi().patchStationNotes(this.stationId, payload).then(
                    (updated) => {
                        this.dirty = false;
                        this.success = true;
                        console.log("success", updated);
                    },
                    () => {
                        this.failed = true;
                        console.log("failed");
                    }
                );
            });
        },
    },
});
</script>

<style scoped>
.notes-view {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: #fcfcfc;
    text-align: left;
    margin-left: 40px;
    margin-right: 40px;
}
.notes-view .header {
    margin-top: 40px;
    display: flex;
    flex-direction: column;
}
.notes-view .lower {
    display: flex;
    flex-direction: row;
    border: 1px solid #d8dce0;
    background: white;
    margin-top: 20px;
}
.notes-view .lower .side {
    flex-basis: 16rem;
    flex-shrink: 0;
}
.notes-view .lower .main {
    flex-grow: 1;
    text-align: left;
}
.notes-view .lower .main.empty {
    padding: 20px;
}

.notification.success {
    margin-top: 20px;
    margin-bottom: 20px;
    padding: 20px;
    border: 2px;
    border-radius: 4px;
}
.notification.success {
    background-color: #d4edda;
}
.notification.failed {
    background-color: #f8d7da;
}
.notifications {
    padding: 10px;
}

.spinner {
    margin-top: 40px;
    margin-left: auto;
    margin-right: auto;
}
</style>
