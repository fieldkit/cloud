<template>
    <StandardLayout>
        <div class="notes-view">
            <div class="header">
                <DoubleHeader />
            </div>
            <div class="lower">
                <div class="side">
                    <StationTabs :stations="stations" :selected="selectedStation" @selected="onSelected" />
                </div>
                <div class="main" v-if="selectedStation && selectedNotes">
                    <NotesForm :station="selectedStation" :notes="selectedNotes" @save="saveForm" v-bind:key="id" @change="onChange" />
                </div>
                <div v-else class="main empty">
                    Please choose a station from the left.
                </div>
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
        id: {
            type: Number,
            required: false,
        },
    },
    data: () => {
        return {
            notes: {},
            dirty: false,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
        }),
        selectedStation() {
            if (this.id && this.stations) {
                return _.first(this.stations.filter((station) => station.id == this.id));
            }
            return null;
        },
        selectedNotes() {
            if (this.id && this.notes) {
                return this.notes[this.id];
            }
            return null;
        },
    },
    mounted(this: any) {
        return this.load(this.$store.state.stations.user.stations).then(() => {
            return this.$store.watch(
                (state) => state.stations.user.stations,
                (newStations, oldStations) => {
                    return this.load(newStations);
                }
            );
        });
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
        load(this: any, stations: any[]) {
            return Promise.all(
                stations.map((station) => {
                    return new FKApi().getStationNotes(station.id).then((notes) => {
                        Vue.set(this.notes, station.id, notes);
                    });
                })
            );
        },
        onSelected(this: any, station) {
            if (this.id != station.id) {
                return this.$router.push({
                    name: "viewStationNotes",
                    params: {
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
                if (window.confirm("Do you really want to leave? you have unsaved changes!")) {
                    this.dirty = false;
                    return true;
                } else {
                    return false;
                }
            }
            return true;
        },
        saveForm(this: any, formNotes: Notes) {
            return serializePromiseChain(formNotes.addedPhotos, (photo) => {
                return new FKApi().uploadStationMedia(this.id, photo.key, photo.file).then((media) => {
                    console.log(media);
                    return [];
                });
            }).then(() => {
                const payload = mergeNotes(this.notes[this.id], formNotes);
                return new FKApi().patchStationNotes(this.id, payload).then(
                    (updated) => {
                        this.dirty = false;
                        console.log("success", updated);
                    },
                    () => {
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
    margin-left: 40px;
    margin-right: 40px;
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
</style>
